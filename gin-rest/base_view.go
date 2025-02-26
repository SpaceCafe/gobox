package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/render"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// Ensure BaseView implements IView interface.
var _ types.IView = (*BaseView)(nil)

// BaseView is a struct that handles rendering of resources in different formats.
type BaseView struct {
	// renderer holds functions to render data into specific MIME types.
	renderer map[string]func(page, pageSize, total, totalPages int, data any) types.IRender

	// IResourceGetter is an interface for getting resources.
	types.IResourceGetter

	// supportedMimeTypes contains a set of supported MIME types.
	supportedMimeTypes map[string]any
}

// SetResource sets the resource for this view.
func (r *BaseView) SetResource(resource types.IResource) {
	r.IResourceGetter = resource
	r.init()
}

// SupportedMimeTypes returns a pointer to the set of supported MIME types.
func (r *BaseView) SupportedMimeTypes() map[string]any {
	return r.supportedMimeTypes
}

// Create handles the creation of a new entity and returns a render object.
func (r *BaseView) Create(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// Read handles reading an entity and returns a render object.
func (r *BaseView) Read(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// List handles listing entities and returns a render object.
func (r *BaseView) List(ctx *gin.Context, entities any) types.IRender {
	return r.newRender(ctx, entities)
}

// Update handles updating an entity and returns a render object.
func (r *BaseView) Update(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// Delete handles deleting an entity and returns a render object.
func (r *BaseView) Delete(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// init initializes the renderer and supportedMimeTypes for the BaseView.
func (r *BaseView) init() {
	r.renderer = map[string]func(page, pageSize, total, totalPages int, data any) types.IRender{
		(&render.JSON{}).MimeType(): newJSONRender,
		(&render.YAML{}).MimeType(): newYAMLRender,
		"*/*":                       newJSONRender,
	}
	r.supportedMimeTypes = map[string]any{
		(&render.JSON{}).MimeType(): nil,
		(&render.YAML{}).MimeType(): nil,
		"*/*":                       nil,
	}
}

// newRender creates a render object based on the context and data.
func (r *BaseView) newRender(ctx *gin.Context, data any) types.IRender {
	listOptions := GetListOptions(ctx)
	page := listOptions.Page
	pageSize := listOptions.PageSize
	total := ctx.GetInt(types.ContextDataTotal)
	if total == 0 {
		total = 1
	}
	totalPages := (total + pageSize - 1) / pageSize

	if renderFunc, ok := r.renderer[ctx.GetString(types.ContextDataRenderMimetype)]; ok {
		return renderFunc(page, pageSize, total, totalPages, data)
	}

	return nil
}

// newJSONRender creates a JSON render object.
func newJSONRender(page, pageSize, total, totalPages int, data any) types.IRender {
	return &render.JSON{Page: page, PageSize: pageSize, Total: total, TotalPages: totalPages, Data: data}
}

// newYAMLRender creates a YAML render object.
func newYAMLRender(page, pageSize, total, totalPages int, data any) types.IRender {
	return &render.YAML{Page: page, PageSize: pageSize, Total: total, TotalPages: totalPages, Data: data}
}
