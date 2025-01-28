package job_manager

import (
	"context"
)

// IJobManager is an interface that defines methods for managing jobs.
// It provides functionality to add jobs, add jobs and wait for their completion, and retrieve jobs by their unique identifier.
type IJobManager interface {
	// Start begins the tracking goroutine with a given context and a done callback function.
	Start(ctx context.Context, done func()) (err error)

	// StartWorker initiates a tracking goroutine using the provided context and done callback function.
	// It starts processing tasks by draining the worker queue and then monitors the jobs queue for new tasks.
	StartWorker(ctx context.Context, done func()) (err error)

	// Stop halts the tracking goroutine.
	Stop()

	// IsReady checks if the JobManager is fully initialized and ready to perform its tasks.
	IsReady() bool

	// WaitUntilReady blocks the execution until the JobManager is ready.
	// This function ensures that any dependent operations are only executed once the JobManager is prepared.
	WaitUntilReady()

	// AddJob adds a new job to the manager and returns a unique identifier for the job.
	AddJob(job Job) (jobID string, err error)

	// AddJobAndWait adds a new job to the manager and waits for its completion, returning the completed job.
	AddJobAndWait(job Job) (err error)

	// GetJob retrieves a job from the manager using its unique identifier.
	GetJob(jobID string, job Job) (err error)
}

// Job represents a task or work item that provides additional
// functionality related to job management and status tracking.
type Job interface {
	// IsPending returns a boolean indicating whether the job is in a pending state,
	// waiting to be started.
	IsPending() bool

	// IsActive returns a boolean indicating whether the job is currently active.
	IsActive() bool

	// IsDone returns a boolean indicating whether the job has been completed or not.
	IsDone() bool

	// Progress returns a percentage (0 to 100) representing the current progress of the job.
	Progress() int

	// Start begins the job processing.
	Start() error
}
