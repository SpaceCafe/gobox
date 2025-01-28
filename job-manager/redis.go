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
	// jobsQueue holds the name of the queue used for storing jobs.
	jobsQueue string

	// workerQueue holds the name of the queue used for managing workers.
	workerQueue string

	// ctx is the context for managing the lifecycle of jobs.
	ctx context.Context

	// done is a function to signal completion or cancellation of jobs.
	done func()

	// jobFactory is a function that creates new Job instances. This is necessary
	// because we need to know what type of job to unmarshal into when retrieving from Redis.
	jobFactory func() Job

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
func NewRedisJobManager(config *Config, jobFactory func() Job) (jobManager *RedisJobManager, err error) {
	jobManager = &RedisJobManager{
		jobsQueue:   fmt.Sprintf("%s:%s", config.RedisNamespace, RedisQueue),
		workerQueue: fmt.Sprintf("%s:%s", config.RedisNamespace, config.WorkerName),
		jobFactory:  jobFactory,
		config:      config,
		ready:       false,
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

// StartWorker initializes a worker for the RedisJobManager.
// It starts by calling Start and then begins draining the worker queue and watching the jobs queue.
func (r *RedisJobManager) StartWorker(ctx context.Context, done func()) (err error) {
	if err = r.Start(ctx, done); err != nil {
		return
	}
	go func() {
		r.drainWorkerQueue()
		r.watchJobsQueue()
	}()
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

// SetJob stores a (new) job to the Redis store.
func (r *RedisJobManager) SetJob(jobID string, job Job) (err error) {
	r.config.Logger.Debugf("adding job '%s': %+v", jobID, job)
	_, err = r.client.JSONSet(r.ctx, fmt.Sprintf("%s:%s", r.config.RedisNamespace, jobID), "$", job).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add job '%s': %v", jobID, err)
	}
	return
}

// GetJob retrieves a job from Redis using the provided jobID and populates the job parameter.
func (r *RedisJobManager) GetJob(jobID string, job Job) (err error) {
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

// AddJob adds a job to the Redis store and returns the job ID.
func (r *RedisJobManager) AddJob(job Job) (jobID string, err error) {
	jobUUID, err := uuid.NewV7()
	if err != nil {
		return
	}
	jobID = jobUUID.String()

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

	err = r.addJob(jobID.String(), job)
	if err != nil {
		return
	}

	messages := subscription.Channel()
	select {
	case <-messages:
		r.config.Logger.Debugf("job '%s' is done", jobID)
		return r.GetJob(jobID.String(), job)
	case <-time.After(r.config.Timeout):
		r.config.Logger.Infof("job '%s' is timed out", jobID)
		return ErrTimeoutExceeded
	case <-r.ctx.Done():
		return ErrJobManagerTerminated
	}
}

// addJob adds a job to the Redis store and sets its expiration.
// It returns any error encountered during the operation.
func (r *RedisJobManager) addJob(jobID string, job Job) (err error) {
	err = r.SetJob(jobID, job)
	if err != nil {
		return
	}

	_, err = r.client.Expire(r.ctx, fmt.Sprintf("%s:%s", r.config.RedisNamespace, jobID), r.config.RedisTTL).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add ttl to job '%s': %v", jobID, err)
		return
	}

	_, err = r.client.LPush(r.ctx, fmt.Sprintf("%s:%s", r.config.RedisNamespace, RedisQueue), jobID).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add job '%s' to queue: %v", jobID, err)
	}
	return
}

// publishJob publishes a job completion message to the Redis channel.
func (r *RedisJobManager) publishJob(jobID string) (err error) {
	var subscribers int64
	subscribers, err = r.client.Publish(r.ctx, fmt.Sprintf("%s:%s:%s", r.config.RedisNamespace, jobID, RedisChannelDone), "1").Result()
	if err != nil {
		return
	}
	r.config.Logger.Debugf("published job '%s' with %d subscribers", jobID, subscribers)
	return
}

// drainWorkerQueue processes all remaining jobs in the worker queue that were left over
// from the previous shutdown or restart.
func (r *RedisJobManager) drainWorkerQueue() {
	r.config.Logger.Infof("draining worker queue '%s'", r.workerQueue)
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			jobID, err := r.client.RPop(r.ctx, r.workerQueue).Result()
			if err != nil {
				return
			}
			r.processJob(jobID)
		}
	}
}

// watchJobsQueue monitors the jobs queue and transfers new jobs to the worker queue
// to ensure they are processed even after a system shutdown or restart.
func (r *RedisJobManager) watchJobsQueue() {
	r.config.Logger.Infof("watching jobs queue '%s'", r.jobsQueue)
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			jobID, err := r.client.BLMove(r.ctx, r.jobsQueue, r.workerQueue, "RIGHT", "LEFT", time.Second).Result()
			if errors.Is(err, redis.Nil) {
				continue
			}
			if err != nil {
				r.config.Logger.Warnf("failed to watch jobs queue '%s': %v", r.jobsQueue, err)
				continue
			}
			if jobID != "" {
				r.processJob(jobID)
			}
		}
	}
}

// processJob handles the execution of a specific job.
func (r *RedisJobManager) processJob(jobID string) {
	job := r.jobFactory()
	r.config.Logger.Infof("processing job '%s'", jobID)

	defer func() {
		select {
		case <-r.ctx.Done():
			return
		default:
			_ = r.publishJob(jobID)
			_, err := r.client.LRem(r.ctx, r.workerQueue, 1, jobID).Result()
			if err != nil {
				r.config.Logger.Warnf("failed to remove job '%s' from worker queue '%s': %v", jobID, r.workerQueue, err)
			}
		}
	}()

	err := r.GetJob(jobID, job)
	if err != nil {
		return
	}

	err = job.Start()
	if err != nil {
		r.config.Logger.Warnf("failed to start job '%s': %v", jobID, err)
		return
	}
	select {
	case <-r.ctx.Done():
		return
	default:
		if err = r.SetJob(jobID, job); err != nil {
			return
		}
		r.config.Logger.Infof("finished job '%s'", jobID)
	}
}
