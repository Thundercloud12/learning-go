// worker.go
package main

import (
	"context"
	"fmt"
	_ "fmt"
	"log"
	"sync"
	"time"
)

type Status string

const (
	Waiting    Status = "waiting"
	InProgress Status = "in-progress"
	Completed  Status = "completed"
	Failed     Status = "failed"
)

type Job struct {
	ID          string
	Type        string
	payload     []byte
	Status      Status
	Attempts    int
	MaxAttempts int

	NextRetryAt time.Time
}

type Jobs struct {
	mu   sync.Mutex
	jobs map[string]*Job
}

type Handler func(ctx context.Context, payload []byte) error

func sleepHandler(ctx context.Context, payload []byte) error {
	time.Sleep(2 * time.Second)
	return nil
}

func failHandler(ctx context.Context, payload []byte) error {
	return fmt.Errorf("job failed")
}

func (s *Jobs) Get(id string) (*Job, bool) {

	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	return job, ok

}

func (s *Jobs) SaveJob(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[job.ID] = job
}

func NewJobStore() *Jobs {
	return &Jobs{
		jobs: make(map[string]*Job),
	}
}

type App struct {
	store    *Jobs
	queue    chan string
	delayed  chan string
	handlers map[string]Handler
}

func (a *App) RegisterHandler(name string, h Handler) {
	a.handlers[name] = h
}

func (J *App) processJob(
	workerId int,
	jobId string,
	
) {
	
	job, ok := J.store.Get(jobId)
	if !ok {
		log.Printf("Job not found soraha hun")
		return
	}

	job.Status = InProgress
	J.store.SaveJob(job)
	

	log.Printf("Zesty job running")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	handler, ok := J.handlers[job.Type]

	if !ok {
		log.Printf("[worker-%d] unknown job type: %s\n", workerId, job.Type)
		job.Status = Failed
		J.store.SaveJob(job)
		return
	}

	err := handler(ctx, job.payload)
	if err == nil {
		job.Status = Completed
		J.store.SaveJob(job)
		log.Printf("[worker-%d] job completed: %s\n", workerId, job.ID)
		return
	}

	// Failure path
	job.Attempts++
	metricsJobFailure.Inc()

	if job.Attempts <= job.MaxAttempts {
		backoff := calculateBackoff(job.Attempts)
		job.NextRetryAt = time.Now().Add(backoff)
		job.Status = Waiting

		J.store.SaveJob(job)
		J.delayed <- job.ID
		log.Printf("[worker-%d] retrying job %s\n", workerId, job.ID)
	} else {
		job.Status = Failed
		J.store.SaveJob(job)
		log.Printf("[worker-%d] job dead: %s\n", workerId, job.ID)
	}

}

func (J *App) worker(
	id int,
	wg *sync.WaitGroup,
	shutdown <-chan struct{},
) {
	defer wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-shutdown:
			log.Printf("Worker shutting down")
			return
		case jobId, ok := <-J.queue:
			if !ok {
				return
			}
			J.processJob(id, jobId)
		}
	}
}

func NewApp(
	queueSize int,
	delayedSize int,
	customHandlers map[string]Handler,
) *App {

	app := &App{
		store:    NewJobStore(),
		queue:    make(chan string, queueSize),
		delayed:  make(chan string, delayedSize),
		handlers: make(map[string]Handler),
	}

	// Add predefined handlers
	app.handlers["sleep"] = sleepHandler
	app.handlers["fail"] = failHandler

	// Add custom handlers if provided
	for name, h := range customHandlers {
		app.handlers[name] = h
	}

	return app

}
