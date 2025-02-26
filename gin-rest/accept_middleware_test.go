package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAcceptHeader(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []string
	}{
		{
			name: "Single media type",
			arg:  "application/json",
			want: []string{"application/json"},
		},
		{
			name: "Multiple media types with q-values",
			arg:  "text/html, application/xml;q=0.9, application/xhtml+xml, */*;q=0.8",
			want: []string{"text/html", "application/xhtml+xml", "application/xml", "*/*"},
		},
		{
			name: "Single media type with q-value",
			arg:  "image/png;q=1.0",
			want: []string{"image/png"},
		},
		{
			name: "Multiple media types without q-values",
			arg:  "audio/mpeg, audio/ogg",
			want: []string{"audio/mpeg", "audio/ogg"},
		},
		{
			name: "Empty accept header",
			arg:  "",
			want: []string{},
		},
		{
			name: "Single wildcard media type",
			arg:  "*/*",
			want: []string{"*/*"},
		},
		{
			name: "Invalid media type with space",
			arg:  "application/ json",
			want: []string{},
		},
		{
			name: "Invalid q-value greater than 1",
			arg:  "text/plain;q=2.0",
			want: []string{"text/plain"},
		},
		{
			name: "Invalid media type with invalid character",
			arg:  "application/<invalid>",
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseAcceptHeader(tt.arg))
		})
	}
}
