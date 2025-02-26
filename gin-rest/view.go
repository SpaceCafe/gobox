package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/gin-rest/render"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// Ensure View implements IView interface.
var _ types.IView = (*View)(nil)

// View is a struct that handles rendering of resources in different formats.
type View struct {
	// renderer holds functions to render data into specific MIME types.
	renderer map[string]func(page, pageSize, total, totalPages int, data any) types.IRender

	// IResourceGetter is an interface for getting resources.
	types.IResourceGetter

	// supportedMimeTypes contains a set of supported MIME types.
	supportedMimeTypes map[string]any
}

// SetResource sets the resource for this view.
func (r *View) SetResource(resource types.IResource) {
	r.IResourceGetter = resource
	r.init()
}

// SupportedMimeTypes returns a pointer to the set of supported MIME types.
func (r *View) SupportedMimeTypes() map[string]any {
	return r.supportedMimeTypes
}

// Create handles the creation of a new entity and returns a render object.
func (r *View) Create(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// Read handles reading an entity and returns a render object.
func (r *View) Read(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// List handles listing entities and returns a render object.
func (r *View) List(ctx *gin.Context, entities any) types.IRender {
	return r.newRender(ctx, entities)
}

// Update handles updating an entity and returns a render object.
func (r *View) Update(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// Delete handles deleting an entity and returns a render object.
func (r *View) Delete(ctx *gin.Context, entity any) types.IRender {
	return r.newRender(ctx, entity)
}

// init initializes the renderer and supportedMimeTypes for the View.
func (r *View) init() {
	r.renderer = map[string]func(page, pageSize, total, totalPages int, data any) types.IRender{
		(&render.JSON{}).MimeType(): newJSONRender,
		(&render.YAML{}).MimeType(): newYAMLRender,
	}
	r.supportedMimeTypes = map[string]any{
		(&render.JSON{}).MimeType(): nil,
		(&render.YAML{}).MimeType(): nil,
	}
}

// newRender creates a render object based on the context and data.
func (r *View) newRender(ctx *gin.Context, data any) types.IRender {
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
