package http_server

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.New()
	logfile := filepath.Join(t.TempDir(), "test.log")
	err := log.SetOutput(logfile)
	require.NoError(t, err)

	tests := []struct {
		name   string
		method string
		url    string
		status int
	}{
		{"Test1", http.MethodGet, "/testpath?id=42", http.StatusOK},
		{"Test2", http.MethodHead, "", http.StatusBadRequest},
		{"Test3", http.MethodPost, "/testpath", http.StatusMultipleChoices},
		{"Test4", http.MethodPut, "", http.StatusInternalServerError},
		{"Test5", http.MethodPatch, "", http.StatusOK},
		{"Test6", http.MethodDelete, "", http.StatusOK},
		{"Test7", http.MethodTrace, "", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = os.Truncate(logfile, 0)
			require.NoError(t, err)

			ctx, _ := gin.CreateTestContext(nil)
			ctx.Request, _ = http.NewRequestWithContext(ctx, tt.method, tt.url, http.NoBody)
			ctx.Status(tt.status)
			NewGinLogger(log)(ctx)

			content, err := os.ReadFile(logfile)
			assert.NoError(t, err)
			assert.NotEmpty(t, content)
			assert.Contains(t, string(content), strconv.Itoa(tt.status))
			assert.Contains(t, string(content), tt.method)
			assert.Contains(t, string(content), tt.url)
		})
	}
}
