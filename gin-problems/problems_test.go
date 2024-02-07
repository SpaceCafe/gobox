package problems

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type wants struct {
		status int
		body   string
	}
	tests := []struct {
		name  string
		arg   string
		wants wants
	}{
		{
			name: "Test1",
			arg:  "/test-error",
			wants: wants{
				status: 400,
				body:   "{\"detail\":\"This is a test error\", \"instance\":\"/test-error\", \"status\":400, \"title\":\"Test Title\", \"type\":\"/errors/test-title\"}",
			},
		},
		{
			name: "Test2",
			arg:  "/test-ok",
			wants: wants{
				status: 200,
				body:   "{\"text\":\"This is a test\"}",
			},
		},
	}

	r := gin.Default()
	r.Use(New())
	r.GET("/test-error", func(c *gin.Context) {
		_ = c.Error(NewProblem("", "Test Title", 400, "This is a test error"))
	})
	r.GET("/test-ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			Text string `json:"text"`
		}{"This is a test"})
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, httptest.NewRequest("GET", tt.arg, nil))
			response := recorder.Result()

			assert.Equal(t, tt.wants.status, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			_ = response.Body.Close()
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wants.body, string(body))
		})
	}
}
