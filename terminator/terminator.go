package terminator

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	// ExitCodeSigTerm is the exit status code for SIGTERM,
	// indicating the container received a SIGTERM by the underlying operating system.
	ExitCodeSigTerm = 143
)

var (
	// osExit is a variable for testing purposes.
	//nolint:gochecknoglobals // Allow mocking os.Exit during tests.
	osExit = os.Exit
)

// Terminator is a struct that manages context cancellation and synchronization.
type Terminator struct {

	// ctx is the context for managing cancellation.
	ctx context.Context

	// waitGroup is used to wait for goroutines to finish.
	waitGroup sync.WaitGroup

	// config holds configuration settings.
	config *Config

	// cancelFn is the function to cancel the context.
	cancelFn context.CancelFunc

	signalCh chan os.Signal
}

// New creates a new Terminator instance with the provided configuration.
func New(config *Config) *Terminator {
	ctx, cancel := context.WithCancel(context.Background())
	terminator := &Terminator{
		ctx:      ctx,
		config:   config,
		cancelFn: cancel,
		signalCh: make(chan os.Signal, 1),
	}

	// Listen to interrupt and termination signals.
	signal.Notify(terminator.signalCh, os.Interrupt, syscall.SIGTERM)
	go terminator.listen()

	return terminator
}

// FullTracking returns the context and done function for full goroutine tracking.
func (r *Terminator) FullTracking() (ctx context.Context, doneFn func()) {
	r.waitGroup.Add(1)
	return r.ctx, r.waitGroup.Done
}

// ContextTracking returns the context and adds to the waitGroup.
// Therefore, the application is terminated after Config.Timeout.
func (r *Terminator) ContextTracking() (ctx context.Context) {
	r.waitGroup.Add(1)
	return r.ctx
}

// NoTracking adds to the waitGroup without returning the context.
// Therefore, the application is terminated after Config.Timeout.
func (r *Terminator) NoTracking() {
	r.waitGroup.Add(1)
}

// Wait blocks until all tracked goroutines have finished.
// Can be used in main function.
func (r *Terminator) Wait() {
	r.waitGroup.Wait()
}

// listen waits for interrupt or termination signals and handles them.
func (r *Terminator) listen() {
	<-r.signalCh
	r.cancelFn()

	// Guarantee termination after specified timeout.
	doneCh := make(chan struct{})
	go func() {
		r.waitGroup.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
	case <-time.After(r.config.Timeout):
	}
	osExit(ExitCodeSigTerm)
}
