package rest

import (
	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// BindID sets the ID of the entity if it implements types.IModel and the ResourceID is present in the URL parameters.
func BindID(ctx *gin.Context, entity any) (aborted bool) {
	id := ctx.Param(types.ResourceID)
	if entity, ok := entity.(types.IModel); id != "" && ok {
		entity.SetID(id)
		return false
	}
	problems.ProblemMissingRequiredParameter.Abort(ctx)
	return true
}

// BindJSON binds the JSON body of an HTTP request to the provided entity.
// It also handles before and after bind hooks, and error handling.
func BindJSON(ctx *gin.Context, entity any) (aborted bool) {
	if BeforeBindHook(ctx, entity) {
		return true
	}
	if HandleError(ctx, ctx.ShouldBindJSON(entity)) {
		return true
	}
	if AfterBindHook(ctx, entity) {
		return true
	}
	return false
}

// BeforeBindHook is a helper function that is executed by a controller before binding data to the entity.
// It checks if the entity implements the IModelBeforeBind interface and calls its BeforeBind method.
// If an error occurs, it handles the error using HandleError and returns true to abort further processing.
func BeforeBindHook(ctx *gin.Context, entity any) (aborted bool) {
	if entity, ok := entity.(types.IModelBeforeBind); ok {
		if HandleError(ctx, entity.BeforeBind(ctx)) {
			return true
		}
	}
	return false
}

// AfterBindHook is a helper function that is executed by a controller after binding data to the entity.
// It checks if the entity implements the IModelAfterBind interface and calls its AfterBind method.
// If an error occurs, it handles the error using HandleError and returns true to abort further processing.
func AfterBindHook(ctx *gin.Context, entity any) (aborted bool) {
	if entity, ok := entity.(types.IModelAfterBind); ok {
		if HandleError(ctx, entity.AfterBind(ctx)) {
			return true
		}
	}
	return false
}

// BeforeRenderHook is a helper function that is executed by a controller before rendering the response.
// It checks if the entity implements the IModelBeforeRender interface and calls its BeforeRender method.
// If an error occurs, it handles the error using HandleError and returns true to abort further processing.
func BeforeRenderHook(ctx *gin.Context, entity any) (aborted bool) {
	if entity, ok := entity.(types.IModelBeforeRender); ok {
		if HandleError(ctx, entity.BeforeRender(ctx)) {
			return true
		}
	}
	return false
}
