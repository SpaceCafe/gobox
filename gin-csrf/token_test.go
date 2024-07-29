package csrf

import (
	"crypto/sha256"
	"crypto/sha512"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
	}{
		{"sha256 hash", args{&Config{SecretKey: []byte("secret"), TokenLength: 32, Signer: sha256.New}}},
		{"sha512 hash", args{&Config{SecretKey: []byte("secret"), TokenLength: 32, Signer: sha512.New}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(tt.args.config)
			println(token.String())
			assert.NoError(t, err)
			assert.NotNil(t, token)
			assert.Equal(t, tt.args.config.TokenLength, len(token.Message))
			assert.NotEmpty(t, token.Signature)
			assert.NotEmpty(t, token.String())
		})
	}
}

func TestNewTokenFromCookie(t *testing.T) {
	type args struct {
		config *Config
		cookie string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"good", args{
			&Config{SecretKey: []byte("secret"), CookieName: "csrf_token", Signer: sha256.New},
			"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA",
		}, nil},
		{"invalid hash", args{
			&Config{SecretKey: []byte("secret"), CookieName: "csrf_token", Signer: sha256.New},
			"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUa",
		}, ErrInvalidTokenSignature},
		{"invalid cookie name", args{
			&Config{SecretKey: []byte("secret"), CookieName: "", Signer: sha256.New},
			"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA",
		}, ErrCookieRetrieval},
		{"invalid secret", args{
			&Config{SecretKey: []byte("foobar"), CookieName: "csrf_token", Signer: sha256.New},
			"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA",
		}, ErrInvalidTokenSignature},
		{"invalid encoding", args{
			&Config{SecretKey: []byte("foobar"), CookieName: "csrf_token", Signer: sha256.New},
			"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898koj==",
		}, ErrTokenDecoding},
		{"invalid encoding", args{
			&Config{SecretKey: []byte("foobar"), CookieName: "csrf_token", Signer: sha256.New},
			"c2VjcmV0",
		}, ErrInvalidSubmittedTokenLength},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(nil)
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.AddCookie(&http.Cookie{Name: tt.args.config.CookieName, Value: tt.args.cookie})
			ctx.Request = req

			token, err := NewTokenFromCookie(tt.args.config, ctx)
			assert.ErrorIs(t, tt.wantErr, err)

			if tt.wantErr == nil {
				assert.NotNil(t, token)
				assert.NotEmpty(t, token.Signature)
				assert.NotEmpty(t, token.String())
			}
		})
	}
}

func TestToken_Compare(t *testing.T) {
	type fields struct {
		encodedToken string
	}
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"valid token",
			fields{"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA"},
			args{"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA"},
			true,
		},
		{
			"invalid token",
			fields{"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898kojUA"},
			args{"GzYwHkqMu75ths44OoiQ8jqjrX03x_Qhd5j3IDolLc0h9l9wNREwg5_BjpFVDwDlFBaY_EFphKq8NN898ko"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &Token{encodedToken: tt.fields.encodedToken}
			assert.Equal(t, tt.want, token.Compare(tt.args.token))
		})
	}
}
