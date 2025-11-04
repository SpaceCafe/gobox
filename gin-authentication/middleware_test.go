package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/spacecafe/gobox/gin-authentication/jwt"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestConfig creates a test configuration with all necessary defaults.
func setupTestConfig() *Config {
	cfg := &Config{
		Tokens:     []string{"valid-token", "another-token"},
		Principals: map[string]string{"user1": "password1", "user2": "password2"},
		JWT: &jwt.Config{
			Secret:            jwt.Secret("this-is-a-test-secret-min-32-chars"),
			RefreshSecret:     jwt.Secret("this-is-a-test-refresh-secret-min-32-chars"),
			Audience:          []string{"test-audience"},
			Issuer:            "test-issuer",
			CookieName:        "__Host-access_token",
			RefreshCookieName: "__Host-refresh_token",
			Signer:            jwt2.SigningMethodHS256,
			AccessTokenTTL:    time.Minute,
			RefreshTokenTTL:   time.Hour,
		},
	}
	cfg.Repository = NewConfigRepository(cfg)
	return cfg
}

// createValidJWT creates a valid JWT token for testing.
func createValidJWT(cfg *Config, subject string) string {
	claims := jwt.NewClaims(subject)
	claims.IdentityClaims.Name = "Test User"
	token := jwt.New(cfg.JWT, claims, jwt.AccessToken)
	signedToken, _ := token.SignedString()
	return signedToken
}

// testHandler is a simple handler that returns 200 OK with the principal ID.
func testHandler(c *gin.Context) {
	principal, ok := PrincipalFromContext(c)
	if ok {
		c.JSON(http.StatusOK, gin.H{"principal_id": principal.ID(), "principal_name": principal.Name()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "no principal"})
	}
}

func TestNew_BearerAuthenticator(t *testing.T) {
	tests := []struct {
		name              string
		setupHeader       func(*Config) string
		expectedStatus    int
		expectedPrincipal bool
		expectedWWWAuth   string
	}{
		{
			name: "valid bearer token",
			setupHeader: func(cfg *Config) string {
				return "Bearer " + createValidJWT(cfg, "test-user")
			},
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
		{
			name:              "empty bearer token",
			setupHeader:       func(cfg *Config) string { return "Bearer " },
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   "Bearer",
		},
		{
			name:              "invalid bearer token",
			setupHeader:       func(cfg *Config) string { return "Bearer invalid-token" },
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   "Bearer",
		},
		{
			name:              "missing bearer prefix",
			setupHeader:       func(cfg *Config) string { return "invalid-token" },
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
		},
		{
			name:              "no authorization header",
			setupHeader:       func(cfg *Config) string { return "" },
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
		},
		{
			name:              "case insensitive bearer",
			setupHeader:       func(cfg *Config) string { return "BEARER " + createValidJWT(cfg, "test-user") },
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			cfg := setupTestConfig()
			cfg.Authenticators = []Authenticator{NewBearerAuthenticator(cfg)}

			r := gin.New()
			r.Use(problems.New())
			r.Use(New(cfg))
			r.GET("/", testHandler)

			req := httptest.NewRequest("GET", "/", http.NoBody)
			header := tt.setupHeader(cfg)
			if header != "" {
				req.Header.Set("Authorization", header)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedWWWAuth != "" {
				assert.Equal(t, tt.expectedWWWAuth, w.Header().Get("WWW-Authenticate"))
			}
		})
	}
}

func TestNew_JWTAuthenticator(t *testing.T) {
	tests := []struct {
		name              string
		setupCookie       func(*Config) *http.Cookie
		expectedStatus    int
		expectedPrincipal bool
		expectedWWWAuth   string
	}{
		{
			name: "valid jwt cookie",
			setupCookie: func(cfg *Config) *http.Cookie {
				token := createValidJWT(cfg, "test-user")
				return &http.Cookie{
					Name:  cfg.JWT.CookieName,
					Value: token,
				}
			},
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
		{
			name: "invalid jwt cookie",
			setupCookie: func(cfg *Config) *http.Cookie {
				return &http.Cookie{
					Name:  cfg.JWT.CookieName,
					Value: "invalid-token",
				}
			},
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   "JWT",
		},
		{
			name: "missing cookie",
			setupCookie: func(cfg *Config) *http.Cookie {
				return nil
			},
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   "JWT",
		},
		{
			name: "wrong cookie name",
			setupCookie: func(cfg *Config) *http.Cookie {
				token := createValidJWT(cfg, "test-user")
				return &http.Cookie{
					Name:  "wrong-cookie-name",
					Value: token,
				}
			},
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   "JWT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			cfg := setupTestConfig()
			cfg.Authenticators = []Authenticator{NewJWTAuthenticator(cfg)}

			r := gin.New()
			r.Use(problems.New())
			r.Use(New(cfg))
			r.GET("/", testHandler)

			req := httptest.NewRequest("GET", "/", http.NoBody)
			cookie := tt.setupCookie(cfg)
			if cookie != nil {
				req.AddCookie(cookie)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedWWWAuth != "" {
				assert.Equal(t, tt.expectedWWWAuth, w.Header().Get("WWW-Authenticate"))
			}
		})
	}
}

func TestNew_TokenAuthenticator(t *testing.T) {
	tests := []struct {
		name              string
		authHeader        string
		expectedStatus    int
		expectedPrincipal bool
		expectedWWWAuth   string
	}{
		{
			name:              "valid token",
			authHeader:        "Token valid-token",
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
		{
			name:              "another valid token",
			authHeader:        "Token another-token",
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
		{
			name:              "invalid token",
			authHeader:        "Token invalid-token",
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   "Token",
		},
		{
			name:              "missing token prefix",
			authHeader:        "valid-token",
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
		},
		{
			name:              "missing token",
			authHeader:        "Token ",
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
		},
		{
			name:              "no authorization header",
			authHeader:        "",
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
		},
		{
			name:              "case insensitive token",
			authHeader:        "token valid-token",
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			cfg := setupTestConfig()
			cfg.Authenticators = []Authenticator{NewTokenAuthenticator(cfg)}

			r := gin.New()
			r.Use(problems.New())
			r.Use(New(cfg))
			r.GET("/", testHandler)

			req := httptest.NewRequest("GET", "/", http.NoBody)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedWWWAuth != "" {
				assert.Equal(t, tt.expectedWWWAuth, w.Header().Get("WWW-Authenticate"))
			}
		})
	}
}

func TestNew_BasicAuthenticator(t *testing.T) {
	tests := []struct {
		name              string
		username          string
		password          string
		setAuth           bool
		expectedStatus    int
		expectedPrincipal bool
		expectedWWWAuth   string
	}{
		{
			name:              "valid credentials",
			username:          "user1",
			password:          "password1",
			setAuth:           true,
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
		{
			name:              "another valid user",
			username:          "user2",
			password:          "password2",
			setAuth:           true,
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
		},
		{
			name:              "invalid password",
			username:          "user1",
			password:          "wrong-password",
			setAuth:           true,
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   `Basic realm="restricted", charset="UTF-8"`,
		},
		{
			name:              "invalid username",
			username:          "unknown",
			password:          "password1",
			setAuth:           true,
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   `Basic realm="restricted", charset="UTF-8"`,
		},
		{
			name:              "no authorization header",
			setAuth:           false,
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			expectedWWWAuth:   `Basic realm="restricted", charset="UTF-8"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			cfg := setupTestConfig()
			cfg.Authenticators = []Authenticator{NewBasicAuthenticator(cfg)}

			r := gin.New()
			r.Use(problems.New())
			r.Use(New(cfg))
			r.GET("/", testHandler)

			req := httptest.NewRequest("GET", "/", http.NoBody)
			if tt.setAuth {
				req.SetBasicAuth(tt.username, tt.password)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedWWWAuth != "" {
				assert.Equal(t, tt.expectedWWWAuth, w.Header().Get("WWW-Authenticate"))
			}
		})
	}
}

func TestNew_MultipleAuthenticators(t *testing.T) {
	tests := []struct {
		name              string
		setupRequest      func(*Config, *http.Request)
		expectedStatus    int
		expectedPrincipal bool
		description       string
	}{
		{
			name: "bearer succeeds first",
			setupRequest: func(cfg *Config, req *http.Request) {
				token := createValidJWT(cfg, "jwt-user")
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
			description:       "Bearer auth should succeed and stop chain",
		},
		{
			name: "token succeeds when bearer not present",
			setupRequest: func(cfg *Config, req *http.Request) {
				req.Header.Set("Authorization", "Token valid-token")
			},
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
			description:       "Token auth should succeed when bearer not used",
		},
		{
			name: "basic succeeds when others not present",
			setupRequest: func(cfg *Config, req *http.Request) {
				req.SetBasicAuth("user1", "password1")
			},
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
			description:       "Basic auth should succeed when others not used",
		},
		{
			name: "jwt cookie succeeds when authorization header is token",
			setupRequest: func(cfg *Config, req *http.Request) {
				token := createValidJWT(cfg, "jwt-user")
				req.AddCookie(&http.Cookie{
					Name:  cfg.JWT.CookieName,
					Value: token,
				})
			},
			expectedStatus:    http.StatusOK,
			expectedPrincipal: true,
			description:       "JWT cookie auth should succeed",
		},
		{
			name: "all methods invalid",
			setupRequest: func(cfg *Config, req *http.Request) {
				// Don't set any valid auth
			},
			expectedStatus:    http.StatusUnauthorized,
			expectedPrincipal: false,
			description:       "Should fail when no valid auth provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			cfg := setupTestConfig()
			cfg.Authenticators = []Authenticator{
				NewBearerAuthenticator(cfg),
				NewJWTAuthenticator(cfg),
				NewTokenAuthenticator(cfg),
				NewBasicAuthenticator(cfg),
			}

			r := gin.New()
			r.Use(problems.New())
			r.Use(New(cfg))
			r.GET("/", testHandler)

			req := httptest.NewRequest("GET", "/", http.NoBody)
			tt.setupRequest(cfg, req)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, tt.description)
		})
	}
}

func TestPrincipalFromContext(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func(*gin.Context)
		expectOk     bool
		expectID     string
	}{
		{
			name: "principal exists in context",
			setupContext: func(c *gin.Context) {
				c.Set(PrincipalContextKey, &DefaultPrincipal{id: "test-id", name: "Test User"})
			},
			expectOk: true,
			expectID: "test-id",
		},
		{
			name: "principal does not exist in context",
			setupContext: func(c *gin.Context) {
				// Don't set anything
			},
			expectOk: false,
		},
		{
			name: "wrong type in context",
			setupContext: func(c *gin.Context) {
				c.Set(PrincipalContextKey, "not-a-principal")
			},
			expectOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupContext(c)

			principal, ok := PrincipalFromContext(c)
			assert.Equal(t, tt.expectOk, ok)

			if tt.expectOk {
				require.NotNil(t, principal)
				assert.Equal(t, tt.expectID, principal.ID())
			} else {
				assert.Nil(t, principal)
			}
		})
	}
}

func TestNew_PrincipalSetInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := setupTestConfig()
	cfg.Authenticators = []Authenticator{NewTokenAuthenticator(cfg)}

	r := gin.New()
	r.Use(problems.New())
	r.Use(New(cfg))
	r.GET("/", func(c *gin.Context) {
		principal, ok := PrincipalFromContext(c)
		require.True(t, ok, "principal should be set in context")
		require.NotNil(t, principal, "principal should not be nil")
		assert.Equal(t, "token", principal.ID())
		assert.Equal(t, "token", principal.Name())
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("Authorization", "Token valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestToken_Renew(t *testing.T) {
	tests := []struct {
		name      string
		tokenType jwt.TokenType
		expectTTL time.Duration
	}{
		{
			name:      "renew access token",
			tokenType: jwt.AccessToken,
			expectTTL: time.Minute,
		},
		{
			name:      "renew refresh token",
			tokenType: jwt.RefreshToken,
			expectTTL: time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := setupTestConfig()
			claims := jwt.NewClaims("test-user")
			claims.IdentityClaims.Name = "Test User"
			token := jwt.New(cfg.JWT, claims, tt.tokenType)
			originalSessionID := token.Claims().SessionID()
			originalExpiresAt := token.Claims().ExpiresAt()
			originalIssuedAt := token.Claims().IssuedAt.Time
			originalNotBefore := token.Claims().NotBefore.Time

			time.Sleep(2 * time.Second)

			err := token.Renew()
			require.NoError(t, err)

			renewedClaims := token.Claims()
			require.NotNil(t, renewedClaims)

			newExpiresAt := renewedClaims.ExpiresAt()
			assert.True(t, newExpiresAt.After(originalExpiresAt))

			// Verify the new expiry is approximately correct (within 1-second tolerance)
			expectedExpiry := time.Now().Add(tt.expectTTL)
			timeDiff := newExpiresAt.Sub(expectedExpiry).Abs()
			assert.Less(t, timeDiff, time.Second)

			assert.NotEqual(t, originalSessionID, renewedClaims.SessionID())
			assert.Equal(t, "test-user", renewedClaims.Subject)
			assert.Equal(t, "Test User", renewedClaims.IdentityClaims.Name)
			assert.Equal(t, originalIssuedAt, renewedClaims.IssuedAt.Time)
			assert.Equal(t, originalNotBefore, renewedClaims.NotBefore.Time)

			// Verify the token can be signed after renewal
			signedToken, err := token.SignedString()
			require.NoError(t, err)
			assert.NotEmpty(t, signedToken)

			// Verify the renewed token can be parsed and validated
			parsedToken, err := jwt.NewFromString(cfg.JWT, signedToken, tt.tokenType)
			require.NoError(t, err)
			assert.True(t, parsedToken.Valid)
			assert.Equal(t, "test-user", parsedToken.Claims().Subject)
			assert.Equal(t, "Test User", parsedToken.Claims().IdentityClaims.Name)
		})
	}
}
