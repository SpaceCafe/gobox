package job_manager

import (
	"context"
	"time"
)

// Manager is an interface that defines methods for managing jobs.
// It provides functionality to add jobs, add jobs and wait for their completion, and retrieve jobs by their unique identifier.
type Manager interface {
	// Start begins the tracking goroutine with a given context and a done callback function.
	Start(ctx context.Context, done func()) (err error)

	// StartWorker initiates a tracking goroutine using the provided context and done callback function.
	// It starts processing tasks by draining the worker queue and then monitors the job queue for new tasks.
	StartWorker(ctx context.Context, done func()) (err error)

	// Stop halts the tracking goroutine.
	Stop()

	// IsReady checks if the Manager is fully initialized and ready to perform its tasks.
	IsReady() bool

	// WaitUntilReady blocks the execution until the Manager is ready.
	// This function ensures that any dependent operations are only executed once the Manager is prepared.
	WaitUntilReady()

	// AddJob adds a new job to the Manager and returns a unique identifier for the job.
	AddJob(entity Job) (jobID string, err error)

	// AddJobAndWait adds a new job to the Manager and waits for its completion, returning the completed job.
	AddJobAndWait(entity Job) (err error)

	// GetJob retrieves a job from the Manager using its unique identifier.
	GetJob(jobID string, entity Job) (err error)

	// GetJobProgress retrieves the current state and progress of a job.
	// Depending on the implementation, this method may also return an optional artifact,
	// such as a message ID that can be used in a later request.
	GetJobProgress(jobID string, lastArtefact any, timeout time.Duration) (state string, progress uint64, artefact any)

	// SetHookContext stores additional context information related to the job hooks.
	SetHookContext(key string, value any)

	// SetJob assigns a job entity to the specified job ID.
	SetJob(jobID string, entity any) (err error)

	// SetJobProgress updates the state and progress of a job within the Manager.
	SetJobProgress(jobID, state string, progress uint64)
}
