package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/logger"
)

// HTTPServer encapsulates an HTTP server with some additional features.
type HTTPServer struct {
	// Config contains configuration settings for the HTTP server.
	config *Config

	// server is an HTTP server instance from the standard library's net/http package.
	// This could be used to manage HTTP connections and routes.
	server *http.Server

	// Engine is an instance from Gin web framework for Go to handle HTTP requests.
	Engine *gin.Engine
}

// NewHTTPServer creates a new instance of HTTPServer with given configuration.
func NewHTTPServer(config *Config) (httpServer *HTTPServer) {
	httpServer = &HTTPServer{
		config: config,

		// Initializes a new http server with given host and port from config,
		// read timeout and read header timeout from config as well.
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", config.Host, config.Port),
			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
		},

		// Initializes a new Gin engine for handling HTTP requests and responses.
		Engine: gin.Default(),
	}

	// Set Gin engine as HTTP server handler.
	httpServer.server.Handler = httpServer.Engine

	// Enables the server to handle 'Method Not Allowed' errors by returning 405 status code.
	httpServer.Engine.HandleMethodNotAllowed = true

	// Registers a handler function that will be called when a request is made with an unsupported HTTP method.
	httpServer.Engine.NoMethod(ProblemMethodNotAllowed.Abort)

	// Registers a handler function that will be called when no route matches for the requested path and method.
	httpServer.Engine.NoRoute(ProblemNoSuchAccessPoint.Abort)

	return
}

// Start function starts the HTTP server in a separate goroutine.
func (r *HTTPServer) Start() {
	logger.Infof("starting web server and listen to %s", r.server.Addr)

	go func() {
		var err error

		if len(r.config.CertFile) > 0 {
			// Starts with TLS.
			r.server.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
			err = r.server.ListenAndServeTLS(r.config.CertFile, r.config.KeyFile)
		} else {
			// Starts without TLS.
			err = r.server.ListenAndServe()
		}

		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Info(err)
		case err == nil:
			logger.Info("http server was stopped")
		default:
			logger.Fatal(err)
		}
	}()
}

// Stop function stops the HTTP server gracefully within a second time limit.
func (r *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	logger.Infof("stopping http server at '%s'", r.server.Addr)

	if err := r.server.Shutdown(ctx); err != nil {
		logger.Warnf("shutdown of http server was unsuccessful: %s", err)
	}
}
