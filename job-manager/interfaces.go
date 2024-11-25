package job_manager

import (
	"context"

	"github.com/google/uuid"
)

// IJobManager is an interface that defines methods for managing jobs.
// It provides functionality to add jobs, add jobs and wait for their completion, and retrieve jobs by their unique identifier.
type IJobManager interface {

	// Start begins the tracking goroutine with a given context and a done callback function.
	Start(ctx context.Context, done func()) (err error)

	// Stop halts the tracking goroutine.
	Stop()

	// IsReady checks if the JobManager is fully initialized and ready to perform its tasks.
	IsReady() bool

	// WaitUntilReady blocks the execution until the JobManager is ready.
	// This function ensures that any dependent operations are only executed once the JobManager is prepared.
	WaitUntilReady()

	// AddJob adds a new job to the manager and returns a unique identifier for the job.
	AddJob(job Job) (jobID uuid.UUID, err error)

	// AddJobAndWait adds a new job to the manager and waits for its completion, returning the completed job.
	AddJobAndWait(job Job) (err error)

	// GetJob retrieves a job from the manager using its unique identifier.
	GetJob(jobID uuid.UUID, job Job) (err error)
}
