package job_manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	RedisMonitorInterval = time.Second * 1
	RedisQueue           = "new"
	RedisChannelDone     = "done"
)

var (
	ErrTimeoutExceeded      = errors.New("timeout exceeded")
	ErrJobManagerTerminated = errors.New("job manager was terminated")
	ErrNoJobPointer         = errors.New("job must be a pointer")
)

// RedisJobManager manages jobs using Redis as the backend.
// It handles job creation, monitoring, and retrieval.
type RedisJobManager struct {

	// ctx is the context for managing the lifecycle of jobs.
	ctx context.Context

	// done is a function to signal completion or cancellation of jobs.
	done func()

	// config holds the configuration settings for the RedisJobManager.
	config *Config

	// client is the Redis client used to interact with the Redis server.
	client *redis.Client

	// readyCond is a condition variable used to signal readiness.
	readyCond *sync.Cond

	// readyMutex is a mutex used to protect access to the ready condition.
	readyMutex sync.Mutex

	// ready indicates whether the job manager is ready to process jobs.
	ready bool
}

// NewRedisJobManager initializes a new RedisJobManager with the given configuration.
// It returns the job manager instance and any error encountered during initialization.
func NewRedisJobManager(config *Config) (jobManager *RedisJobManager, err error) {
	jobManager = &RedisJobManager{
		config: config,
		ready:  false,
	}
	jobManager.readyCond = sync.NewCond(&jobManager.readyMutex)
	return
}

// Start begins the operation of the RedisJobManager.
// It requires a context and a done function to handle graceful shutdown.
func (r *RedisJobManager) Start(ctx context.Context, done func()) (err error) {
	r.config.Logger.Info("starting redis job manager")
	r.ctx = ctx
	r.done = done

	go func() {
		<-ctx.Done()
		r.Stop()
	}()

	r.createClient()
	go r.monitorClientConnection()
	return
}

// Stop gracefully stops the RedisJobManager and closes the Redis client connection.
func (r *RedisJobManager) Stop() {
	defer r.done()
	r.config.Logger.Info("stopping redis job manager")
	_ = r.client.Close()
}

// createClient initializes the Redis client with the configured options.
func (r *RedisJobManager) createClient() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.config.RedisHost, r.config.RedisPort),
		Username: r.config.RedisUsername,
		Password: r.config.RedisPassword,
		DB:       0,
	})
}

// monitorClientConnection continuously checks the connection to Redis.
// It updates the ready state based on the connection status.
func (r *RedisJobManager) monitorClientConnection() {
	for {
		if _, err := r.client.Ping(r.ctx).Result(); err != nil {
			r.setReady(false)
			r.config.Logger.Warnf("failed to ping redis: %v", err)
		} else {
			r.setReady(true)
			r.config.Logger.Debug("successfully pinged redis")
		}
		select {
		case <-r.ctx.Done():
			return
		case <-time.After(RedisMonitorInterval):
		}
	}
}

// IsReady returns the current readiness state of the RedisJobManager.
func (r *RedisJobManager) IsReady() bool {
	return r.ready
}

// WaitUntilReady blocks until the RedisJobManager is ready to process jobs.
func (r *RedisJobManager) WaitUntilReady() {
	r.readyMutex.Lock()
	defer r.readyMutex.Unlock()
	for !r.IsReady() {
		r.readyCond.Wait()
	}
}

// setReady updates the readiness state and notifies any waiting goroutines.
func (r *RedisJobManager) setReady(ready bool) {
	r.readyMutex.Lock()
	defer r.readyMutex.Unlock()
	r.ready = ready
	r.readyCond.Broadcast()
}

// addJob adds a job to the Redis store and sets its expiration.
// It returns any error encountered during the operation.
func (r *RedisJobManager) addJob(jobID uuid.UUID, job Job) (err error) {
	r.config.Logger.Debugf("adding job '%s': %+v", jobID, job)
	_, err = r.client.JSONSet(r.ctx, fmt.Sprintf("%s:%s", r.config.RedisNamespace, jobID), "$", job).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add job '%s': %v", jobID, err)
		return
	}

	_, err = r.client.Expire(r.ctx, fmt.Sprintf("%s:%s", r.config.RedisNamespace, jobID), r.config.RedisTTL).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add ttl to job '%s': %v", jobID, err)
		return
	}

	_, err = r.client.LPush(r.ctx, fmt.Sprintf("%s:%s", r.config.RedisNamespace, RedisQueue), jobID.String()).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add job '%s' to queue: %v", jobID, err)
	}
	return
}

// AddJob adds a job to the Redis store and returns the job ID.
func (r *RedisJobManager) AddJob(job Job) (jobID uuid.UUID, err error) {
	jobID, err = uuid.NewV7()
	if err != nil {
		return
	}

	r.WaitUntilReady()

	return jobID, r.addJob(jobID, job)
}

// AddJobAndWait adds a job and waits for its completion.
// It subscribes to a Redis channel to receive completion notifications.
func (r *RedisJobManager) AddJobAndWait(job Job) (err error) {
	jobID, err := uuid.NewV7()
	if err != nil {
		return
	}

	r.WaitUntilReady()

	subscription := r.client.Subscribe(r.ctx, fmt.Sprintf("%s:%s:%s", r.config.RedisNamespace, jobID, RedisChannelDone))
	defer func(subscription *redis.PubSub) {
		_ = subscription.Close()
	}(subscription)

	err = r.addJob(jobID, job)
	if err != nil {
		return
	}

	messages := subscription.Channel()
	select {
	case <-messages:
		r.config.Logger.Debugf("job '%s' is done", jobID)
		return r.GetJob(jobID, job)
	case <-time.After(r.config.Timeout):
		r.config.Logger.Infof("job '%s' is timed out", jobID)
		return ErrTimeoutExceeded
	case <-r.ctx.Done():
		return ErrJobManagerTerminated
	}
}

// GetJob retrieves a job from Redis using the provided jobID and populates the job parameter.
func (r *RedisJobManager) GetJob(jobID uuid.UUID, job Job) (err error) {
	if reflect.TypeOf(job).Kind() != reflect.Ptr {
		return ErrNoJobPointer
	}

	r.WaitUntilReady()

	jobString, err := r.client.JSONGet(r.ctx, fmt.Sprintf("%s:%s", r.config.RedisNamespace, jobID)).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to get job '%s': %v", jobID, err)
		return
	}
	return json.Unmarshal([]byte(jobString), job)
}
