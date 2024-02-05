package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthentication(t *testing.T) {
	type args struct {
		config   *AuthenticationConfig
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"basic authentication, same username, same blank password",
			args{
				&AuthenticationConfig{Users: map[string]string{"user": "secret"}},
				"user",
				"secret",
			},
			http.StatusNotFound,
		},
		{
			"basic authentication, different username, same blank password",
			args{
				&AuthenticationConfig{Users: map[string]string{"user": "secret"}},
				"another user",
				"secret",
			},
			http.StatusUnauthorized,
		},
		{
			"basic authentication, no username, same blank password",
			args{
				&AuthenticationConfig{Users: map[string]string{"user": "secret"}},
				"",
				"secret",
			},
			http.StatusUnauthorized,
		},
		{
			"basic authentication, no usernames, same blank password",
			args{
				&AuthenticationConfig{Users: map[string]string{"": "secret"}},
				"",
				"secret",
			},
			http.StatusUnauthorized,
		},
		{
			"basic authentication, same username, no passwords",
			args{
				&AuthenticationConfig{Users: map[string]string{"user": ""}},
				"user",
				"",
			},
			http.StatusUnauthorized,
		},
		{
			"basic authentication, same username, same bcrypt hashed passwords",
			args{
				&AuthenticationConfig{Users: map[string]string{"user": "secret", "another user": string(hashBcryptPassword([]byte("another secret")))}},
				"another user",
				"another secret",
			},
			http.StatusNotFound,
		},
		{
			"api-key authentication, same blank password",
			args{
				&AuthenticationConfig{APIKeys: []string{"secret", "another secret"}, HeaderName: "API-Key"},
				"",
				"another secret",
			},
			http.StatusNotFound,
		},
		{
			"api-key authentication, different blank password",
			args{
				&AuthenticationConfig{APIKeys: []string{"secret", "another secret"}, HeaderName: "API-Key"},
				"",
				"unknown secret",
			},
			http.StatusUnauthorized,
		},
		{
			"api-key authentication, no password",
			args{
				&AuthenticationConfig{APIKeys: []string{"secret", "another secret"}, HeaderName: "API-Key"},
				"",
				"",
			},
			http.StatusUnauthorized,
		},
		{
			"api-key authentication, same bcrypt hashed passwords",
			args{
				&AuthenticationConfig{APIKeys: []string{"secret", string(hashBcryptPassword([]byte("another secret")))}, HeaderName: "API-Key"},
				"",
				"another secret",
			},
			http.StatusNotFound,
		},
	}

	config := NewConfig(nil)
	config.Host = "127.0.0.1"
	config.Port = 8888
	server := NewHTTPServer(config)
	server.Start()
	time.Sleep(100 * time.Millisecond)
	defer server.Stop()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := gin.Default()
			handler.Use(Authentication(tt.args.config))
			server.SetEngine(handler)

			// Try to access the server after stopping
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8888", nil)
			assert.NoError(t, err)

			// Set basic authentication, if username is set.
			if len(tt.args.username) > 0 {
				req.SetBasicAuth(tt.args.username, tt.args.password)
			}

			// Set api-key authentication, if HeaderName is set.
			if len(tt.args.config.HeaderName) > 0 {
				req.Header.Set(tt.args.config.HeaderName, tt.args.password)
			}

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			if resp != nil {
				resp.Body.Close()
			}
			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
}

func Test_compareBlankPasswords(t *testing.T) {
	type args struct {
		hashedPassword []byte
		password       []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"different passwords", args{[]byte("another secret"), []byte("secret")}, ErrInvalidPassword},
		{"same passwords", args{[]byte("secret"), []byte("secret")}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := compareBlankPasswords(tt.args.hashedPassword, tt.args.password)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}

func Test_comparePasswords(t *testing.T) {
	type args struct {
		hashedPassword []byte
		password       []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty hashedPassword and password", args{nil, nil}, false},
		{"empty hashedPassword", args{nil, []byte("secret")}, false},
		{"empty password", args{[]byte("secret"), nil}, false},
		{"different blank passwords", args{[]byte("another secret"), []byte("secret")}, false},
		{"same blank passwords", args{[]byte("secret"), []byte("secret")}, true},
		{"same blank unicode password", args{[]byte("🔐secret"), []byte("🔐secret")}, true},
		{"different bcrypt hashed password", args{hashBcryptPassword([]byte("another secret")), []byte("secret")}, false},
		{"same bcrypt hashed password", args{hashBcryptPassword([]byte("secret")), []byte("secret")}, true},
		{"same bcrypt hashed unicode password", args{hashBcryptPassword([]byte("🔐secret")), []byte("🔐secret")}, true},
		{"bcrypt '$2a$' hashed password", args{[]byte("$2a$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC"), []byte("secret")}, true},
		{"bcrypt '$2b$' hashed password", args{[]byte("$2b$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC"), []byte("secret")}, true},
		{"bcrypt '$2x$' hashed password", args{[]byte("$2x$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC"), []byte("secret")}, true},
		{"bcrypt '$2y$' hashed password", args{[]byte("$2y$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC"), []byte("secret")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := comparePasswords(tt.args.hashedPassword, tt.args.password)
			assert.Equalf(t, tt.want, got, "hashedPassword: %s, password: %s", tt.args.hashedPassword, tt.args.password)
		})
	}
}

func hashBcryptPassword(password []byte) []byte {
	h, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(h))
	return h
}
