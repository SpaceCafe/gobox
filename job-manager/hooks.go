package job_manager

// CompletionHook extends the Job interface with an optional OnCompletion hook.
type CompletionHook interface {
	Job
	// OnCompletion is a hook called by the JobManager when the job is completed.
	// This method can be used for any post-processing tasks, such as cleanup, logging,
	// notifying other systems, or persisting job results into a database.
	OnCompletion(ctx HookContext)
}

// CreationHook extends the Job interface with an optional OnCreation hook.
type CreationHook interface {
	Job
	// OnCreation is a hook called by the JobManager when a job is created.
	// This method can be used as a factory method to prepare the struct.
	OnCreation(ctx HookContext)
}

// HookContext provides an interface for accessing and modifying context information related to job hooks.
type HookContext interface {
	// Get retrieves the value associated with the given key from the context.
	Get(key string) any

	// Set stores a value in the context under the specified key.
	Set(key string, value any)
}

// ProcessCompletionHook is a helper function executed by a job manager after completing a job.
// It checks if the entity implements the CompletionHook interface and calls its OnCompletion method.
func ProcessCompletionHook(ctx HookContext, entity any) {
	if entity, ok := entity.(CompletionHook); ok {
		entity.OnCompletion(ctx)
	}
}

// ProcessCreationHook is a helper function that is executed by a job manager after creating a job entity.
// It checks if the entity implements the CreationHook interface and calls its OnCreation method.
func ProcessCreationHook(ctx HookContext, entity any) {
	if entity, ok := entity.(CreationHook); ok {
		entity.OnCreation(ctx)
	}
}
