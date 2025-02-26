package render

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/spacecafe/gobox/gin-rest/types"
	"gopkg.in/yaml.v3"
)

// Ensure YAML implements types.IRender interface.
var _ types.IRender = (*YAML)(nil)

// YAML represents a paginated response with generic data type T.
type YAML struct {
	// Page is the current page number.
	Page int `yaml:"page"`

	// PageSize is the number of items per page.
	PageSize int `yaml:"page_size"`

	// Total is the total number of items.
	Total int `yaml:"total"`

	// TotalPages is the total number of pages.
	TotalPages int `yaml:"total_pages"`

	// Data holds the items for the current page.
	Data any `yaml:"data"`
}

// Render writes the YAML as YAML to the provided http.ResponseWriter, implements render.Render interface.
func (r *YAML) Render(writer http.ResponseWriter) error {
	r.WriteContentType(writer)
	enc := yaml.NewEncoder(writer)
	if err := enc.Encode(r); err != nil {
		return err
	}
	return nil
}

// WriteContentType sets the Content-Type header to application/x-yaml with UTF-8 charset, implements render.Render interface.
func (r *YAML) WriteContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", binding.MIMEYAML+"; charset=utf-8")
}

func (r *YAML) MimeType() string {
	return binding.MIMEYAML
}
