package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	minBatchSiz = 1
	maxBatchSiz = 10
	minInterva  = 1
	maxInterva  = 5
)

func main() {
	stop := make(chan struct{})

	// Start stress test
	startStressTes(stop)

	// Wait for Ctrl+C or SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("Stress test running. Press Ctrl+C to stop.")
	<-sigChan

	log.Println("Shutdown signal received...")
	close(stop)

	// Give goroutines time to exit cleanly
	time.Sleep(2 * time.Second)
	log.Println("Exited cleanly.")
}

func startStressTes(stop <-chan struct{}) {
	rand.Seed(time.Now().UnixNano())

	go logMetric(stop)

	go func() {
		batch := 1
		for {
			select {
			case <-stop:
				log.Println("Stress test stopped")
				return
			default:
			}

			batchSize := rand.Intn(maxBatchSiz-minBatchSiz+1) + minBatchSiz
			interval := time.Duration(
				rand.Intn(maxInterva-minInterva+1)+minInterva,
			) * time.Second

			log.Printf("Starting batch %d: %d jobs over %v", batch, batchSize, interval)

			var wg sync.WaitGroup
			jobSpacing := interval / time.Duration(batchSize)

			for i := 0; i < batchSize; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					jobType := "sleep"
					if rand.Intn(2) == 1 {
						jobType = "fail"
					}

					submitJo(jobType)

					jitter := time.Duration(rand.Int63n(int64(jobSpacing / 2)))
					time.Sleep(jobSpacing/2 + jitter)
				}()
			}

			wg.Wait()
			time.Sleep(interval)
			batch++
		}
	}()
}

func submitJo(jobType string) {
	payload := map[string]string{"type": jobType}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}

	resp, err := http.Post(
		"http://localhost:8080/jobs",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Printf("Failed to submit job: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Job rejected [%d]: %s", resp.StatusCode, body)
		return
	}

	log.Printf("Submitted job: %s", jobType)
}

func logMetric(stop <-chan struct{}) {
	file, err := os.OpenFile(
		"stress_metrics.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			log.Println("Metrics logging stopped")
			return
		case <-ticker.C:
			resp, err := http.Get("http://localhost:8080/metrics")
			if err != nil {
				log.Printf("Failed to fetch metrics: %v", err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Failed to read metrics response: %v", err)
				continue
			}

			entry := time.Now().Format(time.RFC3339) + ": " + string(body) + "\n"
			if _, err := file.WriteString(entry); err != nil {
				log.Printf("Failed to write metrics: %v", err)
			}
		}
	}
}
