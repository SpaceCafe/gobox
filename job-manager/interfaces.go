package job_manager

import (
	"context"
)

const (
	// StatePending indicates that a job is waiting for an available worker to be started.
	// It's the default state of newly created jobs.
	StatePending = "pending"

	// StateRunning signifies that a job is actively being processed. This state can be assigned
	// by a worker when it begins working on a job, or it may be set if GetJobProcess receives
	// a message without an explicitly defined state.
	StateRunning = "running"

	// StateCompleted signifies that a job has been successfully completed. A worker assigns
	// this state upon finishing the job without errors.
	StateCompleted = "completed"

	// StateFailed signifies that a job has been successfully completed. This state can be assigned
	// by A worker upon finishing the job with errors.
	StateFailed = "failed"
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
	AddJob(entity IJob) (jobID string, err error)

	// AddJobAndWait adds a new job to the manager and waits for its completion, returning the completed job.
	AddJobAndWait(entity IJob) (err error)

	// GetJob retrieves a job from the manager using its unique identifier.
	GetJob(jobID string, entity IJob) (err error)

	// GetJobProgress retrieves the current state and progress of a job.
	// Depending on the implementation, this method may also return an optional artifact,
	// such as a message ID that can be utilized in a subsequent request.
	GetJobProgress(jobID string, lastArtefact any) (state string, progress uint64, artefact any)

	// SetJobProgress updates the state and progress of a job within the manager.
	SetJobProgress(jobID, state string, progress uint64)

	// SetHookContext stores additional context information related to the job hooks.
	SetHookContext(key string, value any)
}

// IJob represents a task or work item that provides additional
// functionality related to job management and status tracking.
type IJob interface {
	// Start initiates the execution of the job.
	Start() error
}

// IJobHooks extends the IJob interface with optional hooks, processed by the job manager after completion.
type IJobHooks interface {
	IJob
	// OnCompletion is a hook called by the JobManager when the job is completed.
	// This method can be used for any post-processing tasks, such as cleanup, logging,
	// notifying other systems, or persisting job results into a database.
	OnCompletion(ctx IJobHookContext)
}

// IJobHookContext provides an interface for accessing and modifying context information related to job hooks.
type IJobHookContext interface {
	// Get retrieves the value associated with the given key from the context.
	Get(key string) any

	// Set stores a value in the context under the specified key.
	Set(key string, value any)
}
