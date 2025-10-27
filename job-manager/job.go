package job_manager

// Job represents a task or work item that provides additional
// functionality related to job management and status tracking.
type Job interface {
	// Start initiates the execution of the job.
	Start() error
}
