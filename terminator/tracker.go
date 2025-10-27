package terminator

import (
	"context"
)

// Tracker defines an interface for starting and stopping a trackable goroutine.
type Tracker interface {

	// Start begins the trackable goroutine with a given context.
	Start(ctx context.Context) (err error)

	// Stop halts the tracked goroutine.
	Stop()
}

// CallbackTracker defines an interface for managing the lifecycle of a trackable goroutine with a notification callback.
type CallbackTracker interface {

	// Start begins the trackable goroutine with a given context and a done callback function.
	Start(ctx context.Context, done func()) (err error)

	// Stop halts the tracked goroutine.
	Stop()
}
