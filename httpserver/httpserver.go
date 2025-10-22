package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/logger/types"
)

var (
	ErrNoContext = errors.New("context can not be empty")
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

	// Router is a router group from Gin that allows setting a base path for all routes.
	Router *gin.RouterGroup

	done func()
}

// NewHTTPServer creates a new instance of HTTPServer with given configuration.
func NewHTTPServer(config *Config) *HTTPServer {
	var httpServer = &HTTPServer{
		config: config,

		// Initializes a new http server with given host and port from config,
		// read timeout and read header timeout from config as well.
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", config.Host, config.Port),
			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
		},
	} // Set mode of gin dependent on logging level.
	if config.Logger.Level() == types.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initializes a new Gin engine for handling HTTP requests and responses.
	engine := gin.New()
	engine.Use(Logger(config.Logger), gin.Recovery())
	engine.Use(problems.New())
	httpServer.SetEngine(engine)

	// Enables the server to handle 'Method Not Allowed' errors by returning 405 status code.
	httpServer.Engine.HandleMethodNotAllowed = true

	// Registers a handler function that will be called when a request is made with an unsupported HTTP method.
	httpServer.Engine.NoMethod(func(ctx *gin.Context) {
		_ = ctx.Error(problems.ProblemMethodNotAllowed)
		ctx.Abort()
	})

	// Registers a handler function that will be called when no route matches for the requested path and method.
	httpServer.Engine.NoRoute(func(ctx *gin.Context) {
		_ = ctx.Error(problems.ProblemNoSuchAccessPoint)
		ctx.Abort()
	})

	// Sets the base path for all routes using the Router group.
	if config.BasePath == "" {
		httpServer.Router = &engine.RouterGroup
	} else {
		httpServer.Router = engine.Group(config.BasePath)
	}

	return httpServer
}

func (r *HTTPServer) SetEngine(engine *gin.Engine) {
	r.Engine = engine
	r.server.Handler = engine
}

// Start function starts the HTTP server in a separate goroutine.
func (r *HTTPServer) Start(ctx context.Context, done func()) error {
	r.config.Logger.Infof("starting web server and listen to %s", r.server.Addr)

	if ctx == nil {
		return ErrNoContext
	}

	r.done = done

	go func() {
		var err error

		if r.config.CertFile != "" {
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
			r.config.Logger.Info(err)
		default:
			r.config.Logger.Fatal(err)
		}
	}()

	go func() {
		<-ctx.Done()
		r.Stop()
	}()

	return nil
}

// Stop function stops the HTTP server gracefully within a second time limit.
func (r *HTTPServer) Stop() {
	defer r.done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r.config.Logger.Infof("stopping http server at '%s'", r.server.Addr)

	if err := r.server.Shutdown(ctx); err != nil {
		r.config.Logger.Warnf("shutdown of http server was unsuccessful: %s", err)
	}
}
