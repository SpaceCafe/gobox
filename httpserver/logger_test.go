package httpserver

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	log := logger.New()
	err := log.SetOutput("testdata/test.log")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		_ = os.Truncate("testdata/test.log", 0)
	}()

	type args struct {
		method string
		url    string
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test1", args{http.MethodGet, "/testpath?id=42", http.StatusOK}},
		{"Test2", args{http.MethodHead, "", http.StatusBadRequest}},
		{"Test3", args{http.MethodPost, "/testpath", http.StatusMultipleChoices}},
		{"Test4", args{http.MethodPut, "", http.StatusInternalServerError}},
		{"Test5", args{http.MethodPatch, "", http.StatusOK}},
		{"Test6", args{http.MethodDelete, "", http.StatusOK}},
		{"Test7", args{http.MethodTrace, "", http.StatusOK}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Truncate("testdata/test.log", 0)

			c, _ := gin.CreateTestContext(nil)
			c.Request, _ = http.NewRequest(tt.args.method, tt.args.url, nil)
			c.Status(tt.args.status)
			Logger(log)(c)

			content, err := os.ReadFile("testdata/test.log")
			assert.NoError(t, err)
			assert.NotEmpty(t, content)
			assert.Contains(t, string(content), strconv.Itoa(tt.args.status))
			assert.Contains(t, string(content), tt.args.method)
			assert.Contains(t, string(content), tt.args.url)
		})
	}
}
