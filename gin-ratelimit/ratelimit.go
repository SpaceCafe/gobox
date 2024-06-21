package ratelimit

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
)

// RateLimit holds the channels for rate limiting.
type RateLimit struct {
	// config holds the configuration settings for rate limiting.
	config *Config

	// burstChannel is used to control the burst limit.
	burstChannel chan struct{}

	// concurrentChannel is used to control the maximum number of concurrent requests.
	concurrentChannel chan struct{}

	// queueChannel is used to manage the waiting queue for incoming requests.
	queueChannel chan struct{}
}

// New is a middleware function that enforces rate limiting based on the provided configuration.
// It initializes the rate limiting channels and starts a goroutine to periodically drain the burst channel.
// The middleware controls the burst limit, maximum number of concurrent requests,
// and manages a waiting queue for incoming requests.
func New(config *Config) gin.HandlerFunc {
	rl := &RateLimit{
		config:            config,
		burstChannel:      make(chan struct{}, config.MaxBurstRequests),
		concurrentChannel: make(chan struct{}, config.MaxConcurrentRequests),
		queueChannel:      make(chan struct{}, config.RequestQueueSize),
	}

	go rl.drainBurstChannel()

	return func(ctx *gin.Context) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), rl.config.RequestTimeout)
		defer cancel()

		select {
		case rl.queueChannel <- struct{}{}:
			defer func() { <-rl.queueChannel }()

			select {
			case rl.burstChannel <- struct{}{}:
				select {
				case rl.concurrentChannel <- struct{}{}:
					defer func() { <-rl.concurrentChannel }()
					ctx.Next()

				case <-ctxWithTimeout.Done():
					_ = ctx.Error(problems.ProblemRequestTimeout)
					ctx.Abort()
				}
			case <-ctxWithTimeout.Done():
				_ = ctx.Error(problems.ProblemRequestTimeout)
				ctx.Abort()
			}
		default:
			_ = ctx.Error(problems.ProblemQueueFull)
			ctx.Abort()
		}
	}
}

// drainBurstChannel periodically drains the burstChannel based on the burst duration and maximum burst requests.
// It creates a ticker that ticks at intervals determined by dividing the burst duration by the maximum burst requests.
// At each tick, it attempts to remove an item from the burstChannel, ensuring that the burst limit is respected.
func (rl *RateLimit) drainBurstChannel() {
	ticker := time.NewTicker(rl.config.BurstDuration / time.Duration(rl.config.MaxBurstRequests))
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-rl.burstChannel:
		default:
		}
	}
}
