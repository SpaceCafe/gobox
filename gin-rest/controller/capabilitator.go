package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// CapabilityFn is a struct that implements the types.ControllerCapabilitator interface.
type CapabilityFn struct{}

// Capability is a method of CapabilityFn that sets the Allow and Accept headers.
// It sets the headers and renders an empty response with status OK.
func (r *CapabilityFn) Capability(allow, accept string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Allow", allow)
		ctx.Header("Accept", accept)
		ctx.Render(http.StatusOK, render.Data{Data: []byte{}})
	}
}
