package job_manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"github.com/spacecafe/gobox/logger"
)

const (
	RedisMonitorInterval     = time.Second * 1
	RedisQueuePendingJobs    = "pending"
	RedisQueueCompletedJobs  = "completed"
	RedisStreamJobProgress   = "progress"
	RedisChannelJobCompleted = "completed"
)

var (
	_ Manager = (*RedisManager[any])(nil)
)

// RedisManager manages jobs using Redis as the backend.
// It handles job creation, monitoring, and retrieval.
type RedisManager[T any] struct {
	// pendingJobsQueue holds the name of the queue used for pushing new jobs.
	pendingJobsQueue string

	// completedJobsQueue holds the name of the queue used for fetching completed jobs.
	completedJobsQueue string

	// processingJobsQueue holds the name of the queue used by workers to move jobs from pending to processing state.
	processingJobsQueue string

	// ctx is the context for managing the lifecycle of jobs.
	ctx context.Context

	// log is a custom logger instance used to output messages.
	log logger.Logger

	// done is a function to signal completion or cancellation of jobs.
	done func()

	// cfg holds the configuration settings for the RedisManager.
	cfg *Config

	// client is the Redis client used to interact with the Redis server.
	client *redis.Client

	// readyCond is a condition variable used to signal readiness.
	readyCond *sync.Cond

	// hookContext is passed to Job hooks as a parameter.
	hookContext HookContext

	// readyMutex is a mutex used to protect access to the ready condition.
	readyMutex sync.Mutex

	// ready indicates whether the job manager is ready to process jobs.
	ready bool

	// hasCompletionHooks indicates whether the job type implements CompletionHook interface.
	hasCompletionHooks bool
}

// NewRedisManager initializes a new RedisManager with the given configuration.
// It returns the job manager instance and any error encountered during initialization.
func NewRedisManager[T any](cfg *Config, log logger.Logger) (manager *RedisManager[T], err error) {

	manager = &RedisManager[T]{
		pendingJobsQueue:    cfg.RedisNamespace + ":" + RedisQueuePendingJobs,
		completedJobsQueue:  cfg.RedisNamespace + ":" + RedisQueueCompletedJobs,
		processingJobsQueue: cfg.RedisNamespace + ":" + cfg.WorkerName,
		log:                 log,
		cfg:                 cfg,
		hookContext:         make(KVStorage),
		ready:               false,
		hasCompletionHooks:  false,
	}
	manager.readyCond = sync.NewCond(&manager.readyMutex)
	if _, ok := any((*T)(nil)).(CompletionHook); ok {
		manager.hasCompletionHooks = true
	}
	return
}

// SetHookContext sets a key-value pair in the hook context.
func (r *RedisManager[T]) SetHookContext(key string, value any) {
	r.hookContext.Set(key, value)
}

// Start begins the operation of the RedisManager.
// It requires a context and a done function to handle graceful shutdown.
func (r *RedisManager[T]) Start(ctx context.Context, done func()) (err error) {
	r.log.Info("starting redis job-manager")
	r.ctx = ctx
	r.done = done

	go func() {
		<-ctx.Done()
		r.Stop()
	}()

	r.createClient()
	go r.monitorClientConnection()
	if r.hasCompletionHooks {
		go r.watchCompletedJobsQueue()
	}

	return
}

// StartWorker initializes a worker for the RedisManager.
// It starts by calling Start and then begins draining the worker queue and watching the job queue.
func (r *RedisManager[T]) StartWorker(ctx context.Context, done func()) (err error) {
	if err = r.Start(ctx, done); err != nil {
		return
	}
	go func() {
		r.drainProcessingJobsQueue()
		r.watchPendingJobsQueue()
	}()
	return
}

// Stop gracefully stops the RedisManager and closes the Redis client connection.
func (r *RedisManager[T]) Stop() {
	defer r.done()
	r.log.Info("stopping redis job-manager")
	_ = r.client.Close()
}

// createClient initializes the Redis client with the configured options.
func (r *RedisManager[T]) createClient() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.cfg.RedisHost, r.cfg.RedisPort),
		Username: r.cfg.RedisUsername,
		Password: r.cfg.RedisPassword,
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})
}

// monitorClientConnection continuously checks the connection to Redis.
// It updates the ready state based on the connection status.
func (r *RedisManager[T]) monitorClientConnection() {
	for {
		if _, err := r.client.Ping(r.ctx).Result(); err != nil {
			r.setReady(false)
			r.log.Warnf("job-manager failed to ping redis: %v", err)
		} else {
			r.setReady(true)
			r.log.Debug("job-manager successfully pinged redis")
		}
		select {
		case <-r.ctx.Done():
			return
		case <-time.After(RedisMonitorInterval):
		}
	}
}

// IsReady returns the current readiness state of the RedisManager.
func (r *RedisManager[T]) IsReady() bool {
	return r.ready
}

// WaitUntilReady blocks until the RedisManager is ready to process jobs.
func (r *RedisManager[T]) WaitUntilReady() {
	r.readyMutex.Lock()
	defer r.readyMutex.Unlock()
	for !r.IsReady() {
		r.readyCond.Wait()
	}
}

// setReady updates the readiness state and notifies any waiting goroutines.
func (r *RedisManager[T]) setReady(ready bool) {
	r.readyMutex.Lock()
	defer r.readyMutex.Unlock()
	r.ready = ready
	r.readyCond.Broadcast()
}

// SetJob stores a job to the Redis store.
func (r *RedisManager[T]) SetJob(jobID string, entity any) (err error) {
	r.log.Debugf("job-manager sets job '%s': %+v", jobID, entity)
	_, err = r.client.JSONSet(r.ctx, r.cfg.RedisNamespace+":"+jobID, "$", entity).Result()
	if err != nil {
		r.log.Warnf("job-manager failed to set job '%s': %v", jobID, err)
	}
	return
}

// GetJob retrieves a job from Redis using the provided jobID and populates the job parameter.
func (r *RedisManager[T]) GetJob(jobID string, entity Job) (err error) {
	if reflect.TypeOf(entity).Kind() != reflect.Ptr {
		return ErrNoJobPointer
	}

	r.WaitUntilReady()

	jobString, err := r.client.JSONGet(r.ctx, r.cfg.RedisNamespace+":"+jobID).Result()
	if err != nil {
		r.log.Warnf("job-manager failed to get job '%s': %v", jobID, err)
		return
	}
	err = json.Unmarshal([]byte(jobString), entity)
	if err != nil {
		return
	}
	ProcessCreationHook(r.hookContext, entity)
	return
}

// GetJobProgress retrieves the current state and progress of a job identified by jobID.
// It uses the lastArtefact to determine where to start reading the stream from.
// The function returns the current state, progress, and the last message ID as an artefact.
func (r *RedisManager[T]) GetJobProgress(jobID string, lastArtefact any, timeout time.Duration) (state string, progress uint64, artefact any) {
	lastMessageID, ok := lastArtefact.(string)
	if !ok {
		lastMessageID = "0"
	}

	stream, err := r.client.XRead(r.ctx, &redis.XReadArgs{
		Streams: []string{r.cfg.RedisNamespace + ":" + jobID + ":" + RedisStreamJobProgress, lastMessageID},
		Block:   timeout,
	}).Result()

	// If the stream is empty or there's an error, return the default values.
	if err != nil || len(stream) == 0 {
		return "", 0, lastArtefact
	}

	lastMessage := stream[0].Messages[len(stream[0].Messages)-1]
	artefact = lastMessage.ID

	// Parse the message data to extract the state and progress.
	if state, ok = lastMessage.Values["state"].(string); !ok || state == "" {
		state = StateRunning
	}

	// Convert progress from string to uint.
	progressStr, _ := lastMessage.Values["progress"].(string)
	progress, _ = strconv.ParseUint(progressStr, 10, 8)

	return
}

// SetJobProgress updates the progress of a job identified by jobID in the Redis stream.
// It sets the state and progress values for the specified job.
func (r *RedisManager[T]) SetJobProgress(jobID, state string, progress uint64) {
	_, err := r.client.XAdd(r.ctx, &redis.XAddArgs{
		Stream: r.cfg.RedisNamespace + ":" + jobID + ":" + RedisStreamJobProgress,
		Values: map[string]any{
			"state":    state,
			"progress": progress,
		},
	}).Result()
	if err != nil {
		r.log.Warnf("job-manager failed to set job '%s' progress: %v", jobID, err)
	}
}

// AddJob adds a job to the Redis store and returns the job ID.
func (r *RedisManager[T]) AddJob(entity Job) (jobID string, err error) {
	jobUUID, err := uuid.NewV7()
	if err != nil {
		return
	}
	jobID = jobUUID.String()

	r.WaitUntilReady()

	return jobID, r.addJob(jobID, entity)
}

// AddJobAndWait adds a job and waits for its completion.
// It subscribes to a Redis channel to receive completion notifications.
func (r *RedisManager[T]) AddJobAndWait(entity Job) (err error) {
	jobID, err := uuid.NewV7()
	if err != nil {
		return
	}

	r.WaitUntilReady()

	subscription := r.client.Subscribe(r.ctx, r.cfg.RedisNamespace+":"+jobID.String()+":"+RedisChannelJobCompleted)
	defer func(subscription *redis.PubSub) {
		_ = subscription.Close()
	}(subscription)

	err = r.addJob(jobID.String(), entity)
	if err != nil {
		return
	}

	// Receive a job completion message.
	message := subscription.Channel()
	select {
	case <-message:
		r.log.Debugf("job '%s' was completed", jobID)
		return r.GetJob(jobID.String(), entity)
	case <-time.After(r.cfg.Timeout):
		r.log.Infof("job '%s' was timed out", jobID)
		return ErrTimeoutExceeded
	case <-r.ctx.Done():
		return ErrJobManagerTerminated
	}
}

// addJob adds a job to the Redis store and sets its expiration.
// It returns any error encountered during the operation.
func (r *RedisManager[T]) addJob(jobID string, entity Job) (err error) {
	err = r.SetJob(jobID, entity)
	if err != nil {
		return
	}

	r.SetJobProgress(jobID, StatePending, 0)

	for _, key := range []string{
		r.cfg.RedisNamespace + ":" + jobID,
		r.cfg.RedisNamespace + ":" + jobID + ":" + RedisStreamJobProgress,
	} {
		_, err = r.client.Expire(r.ctx, key, r.cfg.RedisTTL).Result()
		if err != nil {
			r.log.Warnf("job-manager failed to add ttl to key '%s': %v", jobID, err)
		}
	}

	_, err = r.client.LPush(r.ctx, r.pendingJobsQueue, jobID).Result()
	if err != nil {
		r.log.Warnf("job-manager failed to add job '%s' to queue: %v", jobID, err)
	}
	r.SetJobProgress(jobID, StatePending, 0)
	return
}

// sendJobCompletionMessage publishes a job completion message to the global completion queue to initiate
// post-processing tasks and job-specific completion channel to notify any synchronously waiting subscribers.
func (r *RedisManager[T]) sendJobCompletionMessage(jobID string) {
	var subscribers int64
	var err error
	subscribers, err = r.client.Publish(r.ctx, r.cfg.RedisNamespace+":"+jobID+":"+RedisChannelJobCompleted, "1").Result()
	if err != nil {
		r.log.Warnf("job-manager failed to send completion message to job '%s': %v", jobID, err)
	} else {
		r.log.Debugf("job-manager sends completion message of job '%s' with %d subscribers", jobID, subscribers)
	}

	if r.hasCompletionHooks {
		_, err = r.client.LPush(r.ctx, r.completedJobsQueue, jobID).Result()
		if err != nil {
			r.log.Warnf("job-manager failed to add completion message to job '%s' to completed queue: %v", jobID, err)
		}
	}
}

// drainProcessingJobsQueue processes all remaining jobs in the worker queue that were left over
// from the previous shutdown or restart.
func (r *RedisManager[T]) drainProcessingJobsQueue() {
	r.log.Infof("job-manager drains processing queue '%s'", r.processingJobsQueue)
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			jobID, err := r.client.RPop(r.ctx, r.processingJobsQueue).Result()
			if err != nil {
				return
			}
			r.processJob(jobID)
		}
	}
}

// watchPendingJobsQueue monitors the pending jobs queue and transfers new jobs to the processing jobs queue
// to ensure they are processed even after a system shutdown or restart.
func (r *RedisManager[T]) watchPendingJobsQueue() {
	r.log.Infof("job-manager watches pending jobs queue '%s'", r.pendingJobsQueue)
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			jobID, err := r.client.BLMove(r.ctx, r.pendingJobsQueue, r.processingJobsQueue, "RIGHT", "LEFT", time.Second).Result()
			if errors.Is(err, redis.Nil) {
				continue
			}
			if err != nil {
				r.log.Warnf("job-manager failed to watch pending queue '%s': %v", r.pendingJobsQueue, err)
				continue
			}
			if jobID != "" {
				r.processJob(jobID)
			}
		}
	}
}

// watchCompletedJobsQueue initiate post-processing tasks of completed jobs.
func (r *RedisManager[T]) watchCompletedJobsQueue() {
	r.log.Infof("job-manager watches completed queue '%s'", r.completedJobsQueue)
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			response, err := r.client.BRPop(r.ctx, time.Second, r.completedJobsQueue).Result()
			if errors.Is(err, redis.Nil) {
				continue
			}
			if err != nil {
				r.log.Warnf("job-manager failed to watch completed queue '%s': %v", r.completedJobsQueue, err)
				continue
			}
			//nolint:mnd // Redis returns nil or an array of size 2.
			if len(response) == 2 {
				var entity T
				entityRef := any(&entity).(Job)
				err = r.GetJob(response[1], entityRef)
				if err != nil {
					r.log.Warnf("job-manager failed to initiate post-processing tasks for job '%s': %v", response[1], err)
					continue
				}
				ProcessCompletionHook(r.hookContext, entityRef)
			}
		}
	}
}

// processJob handles the execution of a specific job.
func (r *RedisManager[T]) processJob(jobID string) {
	var entity T
	entityRef := any(&entity).(Job)
	r.log.Infof("job-manager processes job '%s'", jobID)

	defer func() {
		r.sendJobCompletionMessage(jobID)
		_, err := r.client.LRem(r.ctx, r.processingJobsQueue, 1, jobID).Result()
		if err != nil {
			r.log.Warnf("job-manager failed to remove job '%s' from processing queue '%s': %v", jobID, r.processingJobsQueue, err)
		}
	}()

	r.SetJobProgress(jobID, StateRunning, 0)
	err := r.GetJob(jobID, entityRef)
	if err != nil {
		r.SetJobProgress(jobID, StateFailed, 0)
		return
	}

	err = entityRef.Start()
	if err != nil {
		r.log.Warnf("job-manager failed to start job '%s': %v", jobID, err)
		r.SetJobProgress(jobID, StateFailed, 0)
		return
	}

	if err = r.SetJob(jobID, entityRef); err != nil {
		r.SetJobProgress(jobID, StateFailed, 0)
		return
	}
	r.log.Infof("job-manager completed job '%s'", jobID)
	r.SetJobProgress(jobID, StateCompleted, 0)
}
