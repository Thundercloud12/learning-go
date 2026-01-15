package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func registerRoutes(queue chan string, store *Jobs) {

	http.HandleFunc("/jobs", postJobHandler(queue, store))
	http.HandleFunc("/jobs/", getJobHandler(store))
	http.HandleFunc("/health", healthHanler)
	http.HandleFunc("/metrics", metricsHandler)

}

func healthHanler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	type metricsResponse struct {
		JobsSuccess    int64 `json:"jobs_success"`
		JobsFailure    int64 `json:"jobs_failure"`
		JobsInProgress int64 `json:"jobs_in_progress"`
		JobsDead       int64 `json:"jobs_dead"`
	}

	resp := metricsResponse{
		JobsSuccess:    metricsJobSuccess.Get(),
		JobsFailure:    metricsJobFailure.Get(),
		JobsInProgress: metricsJobInProgress.Get(),
		JobsDead:       metricsJobsDead.Get(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func postJobHandler(queue chan string, store *Jobs) http.HandlerFunc {

	type request struct {
		Type string `json:"type"`
	}

	type response struct {
		JobId  string `json:"job_id"`
		Status Status `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		job := &Job{
			ID:          uuid.NewString(),
			Type:        req.Type,
			Status:      Waiting,
			MaxAttempts: 3,
		}

		store.SaveJob(job)

		select {
		case queue <- job.ID:
			log.Printf("[api] job submitted %s", job.ID)
		default:
			http.Error(w, "queue full", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response{
			JobId:  job.ID,
			Status: job.Status,
		})
	}

}

func getJobHandler(store *Jobs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/jobs/")
		job, ok := store.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	}
}
