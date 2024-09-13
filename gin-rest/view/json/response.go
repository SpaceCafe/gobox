package json

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin/binding"
)

// Response represents a paginated response with generic data type T.
type Response[T any] struct {

	// Page is the current page number.
	Page int `json:"page"`

	// PageSize is the number of items per page.
	PageSize int `json:"page_size"`

	// Total is the total number of items.
	Total int `json:"total"`

	// TotalPages is the total number of pages.
	TotalPages int `json:"total_pages"`

	// Data holds the items for the current page.
	Data *[]T `json:"data"`
}

// Render writes the Response as JSON to the provided http.ResponseWriter, implements render.Render interface.
func (r *Response[T]) Render(writer http.ResponseWriter) error {
	r.WriteContentType(writer)
	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return err
	}
	return nil
}

// WriteContentType sets the Content-Type header to application/json with UTF-8 charset, implements render.Render interface.
func (r *Response[T]) WriteContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", binding.MIMEJSON+"; charset=utf-8")
}
