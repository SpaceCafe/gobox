package csrf

import (
	"crypto/sha256"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/httpserver"
	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
)

func TestHeaderToken(t *testing.T) {
	type fields struct {
		config       *Config
		serverConfig *httpserver.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"valid",
			fields{&Config{
				ExcludedRoutes: []string{},
				SecretKey:      []byte("secret"),
				CookieName:     "csrf_token",
				HeaderName:     "X-CSRF-Token",
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{BasePath: "/"}},
			"X-CSRF-Token,",
		},
		{"invalid",
			fields{&Config{
				ExcludedRoutes: []string{},
				SecretKey:      []byte("secret"),
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{BasePath: "/"}},
			",",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(problems.New())
			r.Use(New(tt.fields.config, tt.fields.serverConfig))
			r.GET("/test", func(c *gin.Context) {
				name, value := HeaderToken(c)
				c.String(http.StatusOK, "%s,%s", name, value)
			})

			request := httptest.NewRequest("GET", "/test", nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, request)

			response := recorder.Result()
			body, err := io.ReadAll(response.Body)
			assert.NoError(t, err)
			_ = response.Body.Close()

			assert.True(t, strings.HasPrefix(string(body), tt.want))
		})
	}
}
