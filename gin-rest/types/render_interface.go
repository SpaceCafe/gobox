package types

import (
	"github.com/gin-gonic/gin/render"
)

// IRender defines an interface for rendering HTTP responses with a specific MIME type.
type IRender interface {
	// Render is embedded from gin package, providing the basic render functionality.
	render.Render

	// MimeType returns the MIME type of the rendered content.
	MimeType() string
}
