package saml

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/logger"
	"github.com/spacecafe/gosaml"
	"github.com/spacecafe/gosaml/samlsp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		arg  *Config
		want assert.ErrorAssertionFunc
	}{
		{
			name: "invalid certificate",
			arg:  &Config{Logger: logger.Default(), CertFile: "testdata/cert.pem", KeyFile: "testdata/invalid"},
			want: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrInvalidCertificates)
			},
		},
		{
			name: "invalid metadata url",
			arg:  &Config{Logger: logger.Default(), CertFile: "testdata/cert.pem", KeyFile: "testdata/key.pem", IDPMetadataURL: "\x7f.invalid"},
			want: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrInvalidMetadataURL)
			},
		},
		{
			name: "invalid metadata",
			arg:  &Config{Logger: logger.Default(), CertFile: "testdata/cert.pem", KeyFile: "testdata/key.pem", IDPMetadataURL: "http://example.invalid"},
			want: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrFetchMetadata)
			},
		},
		{
			name: "invalid uri",
			arg:  &Config{Logger: logger.Default(), CertFile: "testdata/cert.pem", KeyFile: "testdata/key.pem", IDPMetadataURL: "https://mocksaml.com/api/saml/metadata", URI: "\x7f.invalid"},
			want: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrInvalidApplicationURL)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.arg, &gin.RouterGroup{})
			if !tt.want(t, err) {
				return
			}
		})
	}
}

func TestSAML_RequireAccount(t *testing.T) {
	tests := []struct {
		name      string
		setCookie bool
		want      int
	}{
		{"with session", true, http.StatusOK},
		{"without session", false, http.StatusFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			rg := r.Group("/auth")

			middleware, err := New(&Config{
				CertFile:       "testdata/cert.pem",
				KeyFile:        "testdata/key.pem",
				IDPMetadataURL: "https://mocksaml.com/api/saml/metadata",
				Logger:         logger.Default(),
				CookieName:     "token",
			}, rg)
			assert.NoError(t, err)

			r.GET("/", middleware.RequireAccount(), func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "OK")
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)
			if tt.setCookie {
				session, _ := middleware.middleware.Session.(samlsp.CookieSessionProvider).Codec.New(&saml.Assertion{})
				cookie, err := middleware.middleware.Session.(samlsp.CookieSessionProvider).Codec.Encode(session)
				assert.NoError(t, err)
				request.Header.Set("Cookie", "token="+cookie+"; Path=/; Max-Age=7200")
			}

			r.ServeHTTP(recorder, request)
			response := recorder.Result()
			_ = response.Body.Close()
			assert.Equal(t, tt.want, response.StatusCode)
		})
	}
}

func TestSAML_SetOnError(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	middleware := &SAML{middleware: &samlsp.Middleware{}}
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		http.Redirect(w, r, "http://example.com/set-on-error", http.StatusBadRequest)
	}
	middleware.SetOnError(errFunc)

	middleware.middleware.OnError(recorder, request, &saml.InvalidResponseError{})
	response := recorder.Result()
	_ = response.Body.Close()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "http://example.com/set-on-error", response.Header.Get("Location"))
}

func TestSAML_onError(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	middleware := &SAML{config: &Config{DefaultErrorURI: "http://example.com/on-error", Logger: logger.Default()}}

	middleware.onError(recorder, request, &saml.InvalidResponseError{})
	response := recorder.Result()
	_ = response.Body.Close()
	assert.Equal(t, http.StatusSeeOther, response.StatusCode)
	assert.Equal(t, "http://example.com/on-error", response.Header.Get("Location"))
}

func Test_parseCookieSameSite(t *testing.T) {
	tests := []struct {
		arg  string
		want http.SameSite
	}{
		{"Lax", http.SameSiteLaxMode},
		{"STRICT", http.SameSiteStrictMode},
		{"None", http.SameSiteNoneMode},
		{"InvalidValue", http.SameSiteDefaultMode},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			assert.Equal(t, tt.want, parseCookieSameSite(tt.arg))
		})
	}
}
