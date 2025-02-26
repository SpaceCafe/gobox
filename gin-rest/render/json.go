package render

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/spacecafe/gobox/gin-rest/types"
)

// Ensure JSON implements types.IRender interface.
var _ types.IRender = (*JSON)(nil)

// JSON represents a paginated response with generic data type T.
type JSON struct {
	// Page is the current page number.
	Page int `json:"page"`

	// PageSize is the number of items per page.
	PageSize int `json:"page_size"`

	// Total is the total number of items.
	Total int `json:"total"`

	// TotalPages is the total number of pages.
	TotalPages int `json:"total_pages"`

	// Data holds the items for the current page.
	Data any `json:"data"`
}

// Render writes the YAML as JSON to the provided http.ResponseWriter, implements render.IRender interface.
func (r *JSON) Render(writer http.ResponseWriter) error {
	r.WriteContentType(writer)
	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return err
	}
	return nil
}

// WriteContentType sets the Content-Type header to application/json with UTF-8 charset, implements render.IRender interface.
func (r *JSON) WriteContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", binding.MIMEJSON+"; charset=utf-8")
}

func (r *JSON) MimeType() string {
	return binding.MIMEJSON
}
