package csrf

import (
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/httpserver"
	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type fields struct {
		cookieName  string
		headerName  string
		cookieValue string
		headerValue string
		method      string
	}
	type args struct {
		config       *Config
		serverConfig *httpserver.Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"valid",
			fields{
				"csrf_token",
				"X-CSRF-Token",
				"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA",
				"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA",
				"POST",
			},
			args{&Config{
				ExcludedRoutes: []string{},
				SecretKey:      []byte("secret"),
				CookieName:     "csrf_token",
				HeaderName:     "X-CSRF-Token",
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{}},
			http.StatusOK,
		},
		{"invalid",
			fields{
				"csrf_token",
				"X-CSRF-Token",
				"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898koj",
				"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898koj",
				"POST",
			},
			args{&Config{
				ExcludedRoutes: []string{},
				SecretKey:      []byte("secret"),
				CookieName:     "csrf_token",
				HeaderName:     "X-CSRF-Token",
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{}},
			http.StatusForbidden,
		},
		{"no cookie",
			fields{
				"csrf_token",
				"X-CSRF-Token",
				"",
				"",
				"POST",
			},
			args{&Config{
				ExcludedRoutes: []string{},
				SecretKey:      []byte("secret"),
				CookieName:     "csrf_token",
				HeaderName:     "X-CSRF-Token",
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{}},
			http.StatusForbidden,
		},
		{"no header",
			fields{
				"csrf_token",
				"X-CSRF-Token",
				"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA",
				"",
				"POST",
			},
			args{&Config{
				ExcludedRoutes: []string{},
				SecretKey:      []byte("secret"),
				CookieName:     "csrf_token",
				HeaderName:     "X-CSRF-Token",
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{}},
			http.StatusForbidden,
		},
		{"excluded route",
			fields{
				"csrf_token",
				"X-CSRF-Token",
				"",
				"",
				"POST",
			},
			args{&Config{
				ExcludedRoutes: []string{"/test"},
				SecretKey:      []byte("secret"),
				CookieName:     "csrf_token",
				HeaderName:     "X-CSRF-Token",
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{}},
			http.StatusOK,
		},
		{"get no cookie",
			fields{
				"csrf_token",
				"X-CSRF-Token",
				"",
				"",
				"GET",
			},
			args{&Config{
				ExcludedRoutes: []string{},
				SecretKey:      []byte("secret"),
				CookieName:     "csrf_token",
				HeaderName:     "X-CSRF-Token",
				TokenLength:    32,
				Signer:         sha256.New,
				sameSite:       http.SameSiteStrictMode,
				Logger:         logger.Default(),
			}, &httpserver.Config{BasePath: "/"}},
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(problems.New())
			r.Use(New(tt.args.config, tt.args.serverConfig))
			r.GET("/test", func(_ *gin.Context) {})
			r.POST("/test", func(_ *gin.Context) {})

			request := httptest.NewRequest(tt.fields.method, "/test", nil)
			if len(tt.fields.cookieValue) > 0 {
				request.AddCookie(&http.Cookie{Name: tt.fields.cookieName, Value: tt.fields.cookieValue})
			}
			if len(tt.fields.headerValue) > 0 {
				request.Header.Add(tt.fields.headerName, tt.fields.headerValue)
			}
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, request)

			response := recorder.Result()
			_ = response.Body.Close()
			assert.Equal(t, tt.want, response.StatusCode)

			if tt.fields.method == "GET" {
				assert.Len(t, response.Cookies(), 1)
				assert.Equal(t, tt.fields.cookieName, response.Cookies()[0].Name)
				assert.NotEmpty(t, response.Cookies()[0].Value)
			}
		})
	}
}
