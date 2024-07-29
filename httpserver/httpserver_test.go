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
			Port:              8888,
			Logger:            logger.Default(),
		}},
		{"starts HTTPServer with TLS", "https", &Config{
			Host:              "127.0.0.1",
			CertFile:          "testdata/cert.pem",
			KeyFile:           "testdata/key.pem",
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			Port:              8888,
			Logger:            logger.Default(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewHTTPServer(tt.config)

			// Register a simple route for testing.
			server.Engine.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			// Start the server in a new goroutine.
			server.Start(context.Background(), func() {})

			// Wait for the server to start.
			time.Sleep(1 * time.Second)

			// Make a request to check if the server is running.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, "GET", tt.schema+"://127.0.0.1:8888/ping", http.NoBody)
			assert.NoError(t, err)
			// #nosec G402
			client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check if we get a 200 OK response.
			assert.Equal(t, 200, resp.StatusCode)

			server.Stop()
		})
	}
}

func TestHTTPServer_Stop(t *testing.T) {
	t.Run("stops HTTPServer", func(t *testing.T) {
		config := NewConfig(nil)
		config.Host = "127.0.0.1"
		config.Port = 8888

		server := NewHTTPServer(config)

		// Start the server in a separate goroutine.
		ctx, cancel := context.WithCancel(context.Background())
		server.Start(ctx, func() {})

		// Wait for the server to start
		time.Sleep(100 * time.Millisecond)

		// Call Stop method and wait for it to finish.
		cancel()

		// Try to access the server after stopping
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8888", http.NoBody)
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.Error(t, err)
		if resp != nil {
			resp.Body.Close()
		}
	})
}

func TestNewHTTPServer(t *testing.T) {
	t.Run("creates new HTTPServer", func(t *testing.T) {
		config := NewConfig(nil)

		server := NewHTTPServer(config)

		assert.NotNil(t, server)
		assert.Equal(t, config, server.config)
		assert.NotNil(t, server.server)
		assert.NotNil(t, server.Engine)
		assert.Equal(t, server.Engine, server.server.Handler)
	})
}
