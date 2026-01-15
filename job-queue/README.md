# Job Queue

A simple asynchronous job processing system built in Go, demonstrating concurrency, retry mechanisms, and basic observability.

## Features

- Asynchronous job processing with worker pools
- Automatic retry with exponential backoff
- HTTP API for job submission and monitoring
- Metrics collection for successes, failures, and in-progress jobs
- Graceful shutdown support

## Installation

Ensure you have Go 1.22.2 or later installed.

Clone the repository and navigate to the project directory:

```bash
cd /home/keerthan/go-projects/job-queue
go mod tidy
go build
```

## Usage

Run the application:

```bash
./job-queue
```

The server will start on port 8080.

### Submitting a Job

POST to `/submit` with JSON payload:

```bash
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d '{"type": "sleep"}'
```

Response: `{"job_id": "uuid-here"}`

### Checking Job Status

GET `/status/{job_id}`:

```bash
curl http://localhost:8080/status/uuid-here
```

### Metrics

GET `/metrics`:

```bash
curl http://localhost:8080/metrics
```

Returns JSON with counters for jobs in progress, successes, failures, and dead jobs.

### Health Check

GET `/health`:

```bash
curl http://localhost:8080/health
```

## Architecture

- **Main**: Initializes channels, store, workers, and HTTP server.
- **Workers**: Process jobs from queue, handle retries on failure.
- **Scheduler**: Manages delayed retries (note: currently not started in main).
- **Store**: In-memory job storage with mutex.
- **Backoff**: Exponential backoff calculation.
- **Metrics**: Atomic counters for monitoring.

Jobs flow: Submit -> Queue -> Worker -> Success/Failure -> Retry if failed.

## Configuration

Constants in `config.go`:
- MaxAttempts: 3
- BaseBackoff: 500ms
- MaxBackoff: 5s

## Dependencies

- github.com/google/uuid v1.6.0

## Known Issues

- Scheduler for retries is not started in main.go; add `go retryScheduler(delayed, queue, shutdown, store)` to enable retries.
- In-memory only; no persistence.
- Fixed queue sizes may block on high load.

## Contributing

Fix bugs, add tests, or enhance features. Ensure `go run .` compiles all files.