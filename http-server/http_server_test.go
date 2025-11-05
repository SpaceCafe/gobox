package httpserver_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	httpserver "github.com/spacecafe/gobox/http-server"
	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer_Start(t *testing.T) {
	t.Parallel()

	log := logger.New()

	tests := []struct {
		name   string
		schema string
		config *httpserver.Config
	}{
		{"starts HTTPServer without TLS", "http", &httpserver.Config{
			Host:              "127.0.0.1",
			CertFile:          "",
			KeyFile:           "",
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			Port:              50001,
		}},
		{"starts HTTPServer with TLS", "https", &httpserver.Config{
			Host:              "127.0.0.1",
			CertFile:          "testdata/cert.pem",
			KeyFile:           "testdata/key.pem",
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			Port:              50002,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httpserver.New(tt.config, log)

			// Register a simple route for testing.
			server.Engine.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			// Start the server in a new goroutine.
			err := server.Start(context.Background(), func() {})
			require.NoError(t, err)

			// Wait for the server to start.
			time.Sleep(1 * time.Second)

			// Make a request to check if the server is running.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodGet,
				fmt.Sprintf("%s://127.0.0.1:%d/ping", tt.schema, tt.config.Port),
				http.NoBody,
			)
			require.NoError(t, err)
			// #nosec G402
			client := &http.Client{
				Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			}
			resp, err := client.Do(req)
			require.NoError(t, err)

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
	t.Parallel()

	t.Run("stops HTTPServer", func(t *testing.T) {
		t.Parallel()

		config := &httpserver.Config{}
		config.SetDefaults()
		config.Port = 50003

		server := httpserver.New(config, logger.New())

		// Start the server in a separate goroutine.
		ctx, cancel := context.WithCancel(context.Background())
		err := server.Start(ctx, func() {})
		require.NoError(t, err)

		// Wait for the server to start
		time.Sleep(100 * time.Millisecond)

		// Call the Stop method and wait for it to finish.
		cancel()

		// Try to access the server after stopping
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://127.0.0.1:50003",
			http.NoBody,
		)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.Error(t, err)

		if resp != nil {
			_ = resp.Body.Close()
		}
	})
}
