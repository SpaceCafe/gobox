package job_manager

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
}
