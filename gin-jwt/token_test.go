package jwt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	config := &Config{secretKey: []byte("secret"), Signer: DefaultSigner, TokenExpiration: time.Minute}
	type args struct {
		config *Config
		claims *Claims
	}
	tests := []struct {
		name        string
		args        args
		wantSubject string
		wantErr     error
	}{
		{"valid", args{config, NewClaims(config, "test_subject")}, "test_subject", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(tt.args.config, tt.args.claims)
			if tt.wantErr != nil {
				assert.ErrorAs(t, err, &tt.wantErr)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, tt.wantSubject, token.Claims().Subject)
			}
		})
	}
}

func TestNewTokenFromRequest(t *testing.T) {
	config := &Config{secretKey: []byte("secret"), Signer: DefaultSigner, CookieName: "test_token"}
	requestValidCookie := httptest.NewRequest(http.MethodGet, "/", nil)
	requestValidCookie.AddCookie(&http.Cookie{
		Name:  "test_token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YmplY3QiLCJleHAiOjQ4OTMxMjA2NDAsIm5iZiI6MTczOTUyMDY0MCwiaWF0IjoxNzM5NTIwNjQwLCJyb2xlcyI6bnVsbCwiZ3JvdXBzIjpudWxsLCJlbnRpdGxlbWVudHMiOm51bGwsImF1dGhvcml6YXRpb25fZGV0YWlscyI6bnVsbCwicHJlZmVycmVkX3VzZXJuYW1lIjoiIiwibmFtZSI6IiIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiIn0.iWUAvNJVwWbMbqjDlnUmJt70cobM1oJlgOCxGfKrek4",
	})

	requestValidHeader := httptest.NewRequest(http.MethodGet, "/", nil)
	requestValidHeader.Header.Set(
		"Authorization",
		"Bearer "+
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YmplY3QiLCJleHAiOjQ4OTMxMjA2NDAsIm5iZiI6MTczOTUyMDY0MCwiaWF0IjoxNzM5NTIwNjQwLCJyb2xlcyI6bnVsbCwiZ3JvdXBzIjpudWxsLCJlbnRpdGxlbWVudHMiOm51bGwsImF1dGhvcml6YXRpb25fZGV0YWlscyI6bnVsbCwicHJlZmVycmVkX3VzZXJuYW1lIjoiIiwibmFtZSI6IiIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiIn0.iWUAvNJVwWbMbqjDlnUmJt70cobM1oJlgOCxGfKrek4",
	)
	requestExpired := httptest.NewRequest(http.MethodGet, "/", nil)
	requestExpired.Header.Set(
		"Authorization",
		"Bearer "+
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YmplY3QiLCJleHAiOjE3Mzk1MjA2MDEsIm5iZiI6MTczOTUyMDYwMCwiaWF0IjoxNzM5NTIwNjAwLCJyb2xlcyI6bnVsbCwiZ3JvdXBzIjpudWxsLCJlbnRpdGxlbWVudHMiOm51bGwsImF1dGhvcml6YXRpb25fZGV0YWlscyI6bnVsbCwicHJlZmVycmVkX3VzZXJuYW1lIjoiIiwibmFtZSI6IiIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiIn0.Sy-YVJ7N9DQnKbNJZrVdx0lqcsrR5TN-BHqZYw9UIB4",
	)

	type args struct {
		config *Config
		req    *http.Request
	}
	tests := []struct {
		name        string
		args        args
		wantSubject string
		wantErr     error
	}{
		{"valid cookie", args{config, requestValidCookie}, "test_subject", nil},
		{"valid header", args{config, requestValidHeader}, "test_subject", nil},
		{"expired", args{&Config{secretKey: []byte("secret"), Signer: DefaultSigner}, requestExpired}, "", jwt.ErrTokenExpired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = tt.args.req
			token, err := NewTokenFromRequest(tt.args.config, ctx)
			if tt.wantErr != nil {
				assert.ErrorAs(t, err, &tt.wantErr)
				assert.Nil(t, token)
			} else {
				fmt.Printf("%+v", err)
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, tt.wantSubject, token.Claims().Subject)
			}
		})
	}
}
