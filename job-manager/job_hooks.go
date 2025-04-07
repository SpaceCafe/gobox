package job_manager

// OnCompletionHook is a helper function that is executed by a job manager after completing a job.
// It checks if the entity implements the IJobOnCompletionHook interface and calls its OnCompletion method.
func OnCompletionHook(ctx IJobHookContext, entity any) {
	if entity, ok := entity.(IJobOnCompletionHook); ok {
		entity.OnCompletion(ctx)
	}
}

// OnCreationHook is a helper function that is executed by a job manager after creating a job entity.
// It checks if the entity implements the IJobOnCreationHook interface and calls its OnCreation method.
func OnCreationHook(ctx IJobHookContext, entity any) {
	if entity, ok := entity.(IJobOnCreationHook); ok {
		entity.OnCreation(ctx)
	}
}
