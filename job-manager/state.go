package job_manager

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
