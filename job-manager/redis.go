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
)

const (
	RedisMonitorInterval     = time.Second * 1
	RedisQueuePendingJobs    = "pending"
	RedisQueueCompletedJobs  = "completed"
	RedisStreamJobProgress   = "progress"
	RedisChannelJobCompleted = "completed"
)

var (
	ErrTimeoutExceeded      = errors.New("timeout exceeded")
	ErrJobManagerTerminated = errors.New("job manager was terminated")
	ErrNoJobPointer         = errors.New("job must be a pointer")
)

// RedisJobManager manages jobs using Redis as the backend.
// It handles job creation, monitoring, and retrieval.
type RedisJobManager struct {
	// pendingJobsQueue holds the name of the queue used for pushing new jobs.
	pendingJobsQueue string

	// completedJobsQueue holds the name of the queue used for fetching completed jobs.
	completedJobsQueue string

	// processingJobsQueue holds the name of the queue used by workers to move jobs from pending to processing state.
	processingJobsQueue string

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

	// hasJobHooks indicates whether the job type implements JobHooks interface.
	hasJobHooks bool
}

// NewRedisJobManager initializes a new RedisJobManager with the given configuration.
// It returns the job manager instance and any error encountered during initialization.
func NewRedisJobManager(config *Config, jobFactory func() Job) (jobManager *RedisJobManager, err error) {
	jobManager = &RedisJobManager{
		pendingJobsQueue:    config.RedisNamespace + ":" + RedisQueuePendingJobs,
		completedJobsQueue:  config.RedisNamespace + ":" + RedisQueueCompletedJobs,
		processingJobsQueue: config.RedisNamespace + ":" + config.WorkerName,
		jobFactory:          jobFactory,
		config:              config,
		ready:               false,
		hasJobHooks:         false,
	}
	jobManager.readyCond = sync.NewCond(&jobManager.readyMutex)
	if _, ok := jobFactory().(JobHooks); ok {
		jobManager.hasJobHooks = true
	}
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
	if r.hasJobHooks {
		go r.watchCompletedJobsQueue()
	}

	return
}

// StartWorker initializes a worker for the RedisJobManager.
// It starts by calling Start and then begins draining the worker queue and watching the jobs queue.
func (r *RedisJobManager) StartWorker(ctx context.Context, done func()) (err error) {
	if err = r.Start(ctx, done); err != nil {
		return
	}
	go func() {
		r.drainProcessingJobsQueue()
		r.watchPendingJobsQueue()
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

// SetJob stores a job to the Redis store.
func (r *RedisJobManager) SetJob(jobID string, job Job) (err error) {
	r.config.Logger.Debugf("adding job '%s': %+v", jobID, job)
	_, err = r.client.JSONSet(r.ctx, r.config.RedisNamespace+":"+jobID, "$", job).Result()
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

	jobString, err := r.client.JSONGet(r.ctx, r.config.RedisNamespace+":"+jobID).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to get job '%s': %v", jobID, err)
		return
	}
	return json.Unmarshal([]byte(jobString), job)
}

// GetJobProgress retrieves the current state and progress of a job identified by jobID.
// It uses the lastArtefact to determine where to start reading the stream from.
// The function returns the current state, progress, and the last message ID as an artefact.
func (r *RedisJobManager) GetJobProgress(jobID string, lastArtefact any) (state string, progress uint64, artefact any) {
	lastMessageID, ok := lastArtefact.(string)
	if !ok {
		lastMessageID = "0"
	}

	stream, err := r.client.XRead(r.ctx, &redis.XReadArgs{
		Streams: []string{r.config.RedisNamespace + ":" + jobID + ":" + RedisStreamJobProgress, lastMessageID},
		Block:   0,
	}).Result()

	// If the stream is empty or there's an error, return the default values.
	if err != nil || len(stream) == 0 {
		return StatePending, 0, "0"
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
func (r *RedisJobManager) SetJobProgress(jobID, state string, progress uint64) {
	_, err := r.client.XAdd(r.ctx, &redis.XAddArgs{
		Stream: r.config.RedisNamespace + ":" + jobID + ":" + RedisStreamJobProgress,
		Values: map[string]interface{}{
			"state":    state,
			"progress": progress,
		},
	}).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to set job '%s' progress: %v", jobID, err)
	}
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

	subscription := r.client.Subscribe(r.ctx, r.config.RedisNamespace+":"+jobID.String()+":"+RedisChannelJobCompleted)
	defer func(subscription *redis.PubSub) {
		_ = subscription.Close()
	}(subscription)

	err = r.addJob(jobID.String(), job)
	if err != nil {
		return
	}

	// Receive job completion message
	message := subscription.Channel()
	select {
	case <-message:
		r.config.Logger.Debugf("job '%s' was completed", jobID)
		return r.GetJob(jobID.String(), job)
	case <-time.After(r.config.Timeout):
		r.config.Logger.Infof("job '%s' was timed out", jobID)
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

	_, err = r.client.Expire(r.ctx, r.config.RedisNamespace+":"+jobID, r.config.RedisTTL).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add ttl to job '%s': %v", jobID, err)
		return
	}

	_, err = r.client.LPush(r.ctx, r.pendingJobsQueue, jobID).Result()
	if err != nil {
		r.config.Logger.Warnf("failed to add job '%s' to queue: %v", jobID, err)
	}
	return
}

// sendJobCompletionMessage publishes a job completion message to the global completion queue to initiate
// post-processing tasks and job-specific completion channel to notify any synchronously waiting subscribers.
func (r *RedisJobManager) sendJobCompletionMessage(jobID string) {
	var subscribers int64
	var err error
	subscribers, err = r.client.Publish(r.ctx, r.config.RedisNamespace+":"+jobID+":"+RedisChannelJobCompleted, "1").Result()
	if err != nil {
		r.config.Logger.Warnf("failed to send completion message to job '%s': %v", jobID, err)
	} else {
		r.config.Logger.Debugf("send completion message of job '%s' with %d subscribers", jobID, subscribers)
	}

	if r.hasJobHooks {
		_, err = r.client.LPush(r.ctx, r.completedJobsQueue, jobID).Result()
		if err != nil {
			r.config.Logger.Warnf("failed to add completion message to job '%s' to completed jobs queue: %v", jobID, err)
		}
	}
	return
}

// drainProcessingJobsQueue processes all remaining jobs in the worker queue that were left over
// from the previous shutdown or restart.
func (r *RedisJobManager) drainProcessingJobsQueue() {
	r.config.Logger.Infof("drain processing jobs queue '%s'", r.processingJobsQueue)
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
func (r *RedisJobManager) watchPendingJobsQueue() {
	r.config.Logger.Infof("watch pending jobs queue '%s'", r.pendingJobsQueue)
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
				r.config.Logger.Warnf("failed to watch pending jobs queue '%s': %v", r.pendingJobsQueue, err)
				continue
			}
			if jobID != "" {
				r.processJob(jobID)
			}
		}
	}
}

// watchCompletedJobsQueue initiate post-processing tasks of completed jobs.
func (r *RedisJobManager) watchCompletedJobsQueue() {
	r.config.Logger.Infof("watch completed jobs queue '%s'", r.completedJobsQueue)
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			jobIDs, err := r.client.BRPop(r.ctx, time.Second, r.processingJobsQueue).Result()
			if errors.Is(err, redis.Nil) {
				continue
			}
			if err != nil {
				r.config.Logger.Warnf("failed to watch completed jobs queue '%s': %v", r.pendingJobsQueue, err)
				continue
			}
			for _, jobID := range jobIDs {
				if job, ok := r.jobFactory().(JobHooks); ok {
					err = r.GetJob(jobID, job)
					if err != nil {
						r.config.Logger.Warnf("failed to initiate post-processing tasks for job '%s': %v", jobID, err)
						continue
					}
					job.OnCompletion()
				}
			}
		}
	}
}

// processJob handles the execution of a specific job.
func (r *RedisJobManager) processJob(jobID string) {
	job := r.jobFactory()
	r.config.Logger.Infof("processing job '%s'", jobID)

	defer func() {
		r.sendJobCompletionMessage(jobID)
		_, err := r.client.LRem(r.ctx, r.processingJobsQueue, 1, jobID).Result()
		if err != nil {
			r.config.Logger.Warnf("failed to remove job '%s' from processing jobs queue '%s': %v", jobID, r.processingJobsQueue, err)
		}
	}()

	r.SetJobProgress(jobID, StateRunning, 0)
	err := r.GetJob(jobID, job)
	if err != nil {
		r.SetJobProgress(jobID, StateFailed, 0)
		return
	}

	err = job.Start()
	if err != nil {
		r.config.Logger.Warnf("failed to start job '%s': %v", jobID, err)
		r.SetJobProgress(jobID, StateFailed, 0)
		return
	}

	if err = r.SetJob(jobID, job); err != nil {
		r.SetJobProgress(jobID, StateFailed, 0)
		return
	}
	r.config.Logger.Infof("completed job '%s'", jobID)
	r.SetJobProgress(jobID, StateCompleted, 0)
}
