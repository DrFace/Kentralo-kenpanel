package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// JobState represents the execution lifecycle of a background task.
type JobState string

const (
	StatePending    JobState = "pending"
	StateRunning    JobState = "running"
	StateCompleted  JobState = "completed"
	StateFailed     JobState = "failed"
	StateRolledBack JobState = "rolled_back"
)

// Job defines a background work unit.
type Job struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Payload    map[string]interface{} `json:"payload"`
	State      JobState               `json:"state"`
	MaxRetries int                    `json:"max_retries"`
	RetryCount int                    `json:"retry_count"`
	ErrorMsg   string                 `json:"error_msg,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// Handler defines the function signature for processing a specific job type.
type Handler func(ctx context.Context, job *Job) error

// Queue manages an in-memory or persisted concurrent job worker pool.
type Queue struct {
	mu       sync.RWMutex
	jobs     map[string]*Job
	handlers map[string]Handler
	workChan chan *Job
}

// NewQueue creates and initializes a task execution queue.
func NewQueue(bufferSize int) *Queue {
	return &Queue{
		jobs:     make(map[string]*Job),
		handlers: make(map[string]Handler),
		workChan: make(chan *Job, bufferSize),
	}
}

// RegisterHandler binds a handler function to a job type string.
func (q *Queue) RegisterHandler(jobType string, handler Handler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[jobType] = handler
}

// Enqueue submits a new job to the queue.
func (q *Queue) Enqueue(job *Job) {
	q.mu.Lock()
	job.State = StatePending
	job.CreatedAt = time.Now().UTC()
	job.UpdatedAt = job.CreatedAt
	q.jobs[job.ID] = job
	q.mu.Unlock()

	q.workChan <- job
}

// StartWorkers launches n concurrent worker goroutines.
func (q *Queue) StartWorkers(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					return
				case job := <-q.workChan:
					q.processJob(ctx, job)
				}
			}
		}(i)
	}
}

func (q *Queue) processJob(ctx context.Context, job *Job) {
	q.mu.Lock()
	handler, exists := q.handlers[job.Type]
	job.State = StateRunning
	job.UpdatedAt = time.Now().UTC()
	q.mu.Unlock()

	if !exists {
		q.mu.Lock()
		job.State = StateFailed
		job.ErrorMsg = fmt.Sprintf("no handler registered for job type: %s", job.Type)
		job.UpdatedAt = time.Now().UTC()
		q.mu.Unlock()
		return
	}

	err := handler(ctx, job)
	q.mu.Lock()
	defer q.mu.Unlock()

	if err != nil {
		job.RetryCount++
		if job.RetryCount <= job.MaxRetries {
			job.State = StatePending
			job.UpdatedAt = time.Now().UTC()
			go func(j *Job) {
				time.Sleep(time.Duration(j.RetryCount*2) * time.Second)
				q.workChan <- j
			}(job)
			return
		}
		job.State = StateFailed
		job.ErrorMsg = err.Error()
	} else {
		job.State = StateCompleted
	}
	job.UpdatedAt = time.Now().UTC()
}
