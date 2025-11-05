package csrf_test

import (
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	authentication "github.com/spacecafe/gobox/gin-authentication"
	csrf "github.com/spacecafe/gobox/gin-csrf"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockPrincipal struct{}

func (r *MockPrincipal) ID() string {
	return ""
}

func (r *MockPrincipal) Name() string {
	return ""
}

type MockSession struct {
	MockPrincipal

	sessionID string
}

func (r *MockSession) CreatedAt() time.Time {
	return time.Now()
}

func (r *MockSession) ExpiresAt() time.Time {
	return time.Now().Add(time.Hour)
}

func (r *MockSession) SessionID() string {
	return r.sessionID
}

func setupTestConfig() *csrf.Config {
	cfg := &csrf.Config{
		Secret:     csrf.Secret("this-is-a-test-secret-min-32-chars"),
		CookieName: "csrf-token",
		HeaderName: "X-CSRF-Token",
		Signer:     sha256.New,
	}

	return cfg
}

func TestNew(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		setupRequest   func(*csrf.Config, *gin.Context, *http.Request)
		expectedStatus int
	}{
		{
			name:           "GET method",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "POST method with missing token headers",
			method: http.MethodPost,
			setupRequest: func(_ *csrf.Config, _ *gin.Context, _ *http.Request) {
				// Intentionally not setting any CSRF tokens
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "POST method with invalid tokens",
			method: http.MethodPost,
			setupRequest: func(cfg *csrf.Config, ctx *gin.Context, req *http.Request) {
				ctx.Set(
					authentication.PrincipalContextKey,
					&MockSession{sessionID: "this-is-a-test-session-id"},
				)
				req.AddCookie(&http.Cookie{
					Name:  cfg.CookieName,
					Value: "invalid-token",
				})
				req.Header.Set(cfg.HeaderName, "invalid-token")
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "POST method with another invalid tokens",
			method: http.MethodPost,
			setupRequest: func(cfg *csrf.Config, ctx *gin.Context, req *http.Request) {
				ctx.Set(
					authentication.PrincipalContextKey,
					&MockSession{sessionID: "this-is-a-test-session-id"},
				)
				req.AddCookie(&http.Cookie{
					Name:  cfg.CookieName,
					Value: "aW52YWxpZA.aW52YWxpZA",
				})
				req.Header.Set(cfg.HeaderName, "aW52YWxpZA.aW52YWxpZA")
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "POST method with mismatched token headers",
			method: http.MethodPost,
			setupRequest: func(cfg *csrf.Config, ctx *gin.Context, req *http.Request) {
				ctx.Set(
					authentication.PrincipalContextKey,
					&MockSession{sessionID: "this-is-a-test-session-id"},
				)

				cookieToken, err := csrf.NewToken(cfg, ctx)
				require.NoError(t, err)
				require.NotEmpty(t, cookieToken)

				headerToken, err := csrf.NewToken(cfg, ctx)
				require.NoError(t, err)
				require.NotEmpty(t, headerToken)

				req.AddCookie(cookieToken.Cookie())
				req.Header.Set(cfg.HeaderName, headerToken.String())
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "POST method with valid token headers",
			method: http.MethodPost,
			setupRequest: func(cfg *csrf.Config, ctx *gin.Context, req *http.Request) {
				ctx.Set(
					authentication.PrincipalContextKey,
					&MockSession{sessionID: "this-is-a-test-session-id"},
				)

				token, err := csrf.NewToken(cfg, ctx)
				require.NoError(t, err)
				require.NotEmpty(t, token)

				req.AddCookie(token.Cookie())
				req.Header.Set(cfg.HeaderName, token.String())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "PUT method with too short session ID",
			method: http.MethodPut,
			setupRequest: func(cfg *csrf.Config, ctx *gin.Context, _ *http.Request) {
				ctx.Set(authentication.PrincipalContextKey, &MockSession{sessionID: "too-short"})
				token, err := csrf.NewToken(cfg, ctx)
				require.Error(t, err)
				assert.Nil(t, token)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "PATCH method without session ID",
			method: http.MethodPatch,
			setupRequest: func(cfg *csrf.Config, ctx *gin.Context, _ *http.Request) {
				ctx.Set(authentication.PrincipalContextKey, &MockPrincipal{})
				token, err := csrf.NewToken(cfg, ctx)
				require.Error(t, err)
				assert.Nil(t, token)
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := setupTestConfig()
			response := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(response)

			router.Use(problems.New())
			router.Use(func(ctx *gin.Context) {
				ctx.Set(
					authentication.PrincipalContextKey,
					&MockSession{sessionID: "this-is-a-test-session-id"},
				)
				ctx.Next()
			})
			router.Use(csrf.New(cfg))
			router.Any("/", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(tt.method, "/", http.NoBody)
			if tt.setupRequest != nil {
				tt.setupRequest(cfg, ctx, req)
			}

			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedStatus, response.Code)
		})
	}
}
