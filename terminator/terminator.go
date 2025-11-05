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
	ExitCodeSigTerm = 128 + int(syscall.SIGTERM) // equals 143
)

// OsExit is a variable for testing purposes.
//
//nolint:gochecknoglobals // This is a mock for os.Exit used in tests to prevent actual program termination
var OsExit = os.Exit

// Terminator is a struct that manages context cancellation and synchronization.
type Terminator struct {
	// ctx is the context for managing cancellation.
	//nolint:containedctx // Terminator is an extension of context.Context to provide additional functionality.
	ctx context.Context

	// waitGroup is used to wait for goroutines to finish.
	waitGroup sync.WaitGroup

	// cfg holds configuration settings.
	cfg *Config

	// cancelFn is the function to cancel the context.
	cancelFn context.CancelFunc

	signalCh chan os.Signal

	doneCh chan struct{}
}

// New creates a new Terminator instance with the provided configuration.
func New(cfg *Config) *Terminator {
	ctx, cancel := context.WithCancel(context.Background())
	terminator := &Terminator{
		ctx:      ctx,
		cfg:      cfg,
		cancelFn: cancel,
		signalCh: make(chan os.Signal, 1),
		doneCh:   make(chan struct{}),
	}

	// Listen to interrupt and termination signals.
	signal.Notify(terminator.signalCh, os.Interrupt, syscall.SIGTERM)

	go terminator.awaitSignal()

	return terminator
}

// Context returns the context but does not track the goroutine.
// This is useful when you need the context outside the termination flow.
func (r *Terminator) Context() (ctx context.Context) {
	return r.ctx
}

// Go calls the given task in a new goroutine and adds that task to the waitGroup.
// When the task returns, it's removed from the waitGroup.
func (r *Terminator) Go(task func()) {
	r.waitGroup.Add(1)

	go func() {
		defer r.waitGroup.Done()

		task()
	}()
}

// Track increments the waitGroup by one without returning the context.
// Therefore, the application is terminated after Config.Timeout.
func (r *Terminator) Track() {
	r.waitGroup.Add(1)
}

// TrackWithContext returns the context and increments the waitGroup by one.
// Therefore, the application is terminated after Config.Timeout.
func (r *Terminator) TrackWithContext() (ctx context.Context) {
	r.waitGroup.Add(1)

	return r.ctx
}

// TrackWithDone returns the context and done function for full goroutine tracking.
func (r *Terminator) TrackWithDone() (ctx context.Context, doneFn func()) {
	r.waitGroup.Add(1)

	return r.ctx, r.waitGroup.Done
}

// Wait blocks until all tracked goroutines have finished.
// If the `Track()` method is used, it'll never return.
// Use this function at the end of the main function.
func (r *Terminator) Wait() {
	<-r.doneCh
}

// awaitSignal waits for interrupt or termination signals and handles them.
func (r *Terminator) awaitSignal() {
	<-r.signalCh
	r.cancelFn()

	// Guarantee termination after the specified timeout.
	go func() {
		r.waitGroup.Wait()
		close(r.doneCh)
	}()

	select {
	case <-r.doneCh:
	case <-time.After(r.cfg.Timeout):
	}

	if r.cfg.Force {
		OsExit(ExitCodeSigTerm)
	}
}
