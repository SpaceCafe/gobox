package terminator

import (
	"context"
)

// FullTracking is an interface that defines methods for starting and stopping a tracking goroutine.
type FullTracking interface {

	// Start begins the tracking goroutine with a given context and a done callback function.
	Start(ctx context.Context, done func()) (err error)

	// Stop halts the tracking goroutine.
	Stop()
}

// ContextTracking is an interface that defines methods for starting and stopping a tracking goroutine.
type ContextTracking interface {

	// Start begins the tracking goroutine with a given context.
	Start(ctx context.Context) (err error)

	// Stop halts the tracking goroutine.
	Stop()
}
