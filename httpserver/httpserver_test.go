package httpserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServer_Start(t *testing.T) {
	log := logger.New()

	tests := []struct {
		name   string
		schema string
		config *Config
	}{
		{"starts HTTPServer without TLS", "http", &Config{
			Host:              "127.0.0.1",
			CertFile:          "",
			KeyFile:           "",
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			Port:              58080,
		}},
		{"starts HTTPServer with TLS", "https", &Config{
			Host:              "127.0.0.1",
			CertFile:          "testdata/cert.pem",
			KeyFile:           "testdata/key.pem",
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			Port:              58080,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(tt.config, log)

			// Register a simple route for testing.
			server.Engine.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			// Start the server in a new goroutine.
			err := server.Start(context.Background(), func() {})
			assert.NoError(t, err)

			// Wait for the server to start.
			time.Sleep(1 * time.Second)

			// Make a request to check if the server is running.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, "GET", tt.schema+"://127.0.0.1:58080/ping", http.NoBody)
			assert.NoError(t, err)
			// #nosec G402
			client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer func(resp *http.Response) {
				_ = resp.Body.Close()
			}(resp)

			// Check if we get a 200 OK response.
			assert.Equal(t, 200, resp.StatusCode)

			server.Stop()
		})
	}
}

func TestHTTPServer_Stop(t *testing.T) {
	t.Run("stops HTTPServer", func(t *testing.T) {
		config := &Config{}
		config.SetDefaults()
		config.Port = 58080

		server := New(config, logger.New())

		// Start the server in a separate goroutine.
		ctx, cancel := context.WithCancel(context.Background())
		err := server.Start(ctx, func() {})
		assert.NoError(t, err)

		// Wait for the server to start
		time.Sleep(100 * time.Millisecond)

		// Call the Stop method and wait for it to finish.
		cancel()

		// Try to access the server after stopping
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:58080", http.NoBody)
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.Error(t, err)
		if resp != nil {
			_ = resp.Body.Close()
		}
	})
}

func TestNewHTTPServer(t *testing.T) {
	t.Run("creates new HTTPServer", func(t *testing.T) {
		config := &Config{}
		config.SetDefaults()
		server := New(config, logger.New())

		assert.NotNil(t, server)
		assert.Equal(t, config, server.config)
		assert.NotNil(t, server.server)
		assert.NotNil(t, server.Engine)
		assert.Equal(t, server.Engine, server.server.Handler)
	})
}
