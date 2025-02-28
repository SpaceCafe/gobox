package job_manager

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	//nolint:gochecknoglobals // Used for testing OnCompletion method.
	onCompletionHookCalled = ""
)

type TestJob struct {
	ExitCode int
	StdIn    string
	StdOut   string
	StdErr   string
}

func (r *TestJob) Start() error {
	r.ExitCode = 0
	r.StdOut = r.StdIn
	return nil
}

func (r *TestJob) OnCompletion(ctx IJobHookContext) {
	onCompletionHookCalled = ctx.Get("scope").(string)
}

func TestNewRedisJobManager(t *testing.T) {
	config := NewConfig(nil)
	config.RedisHost = "127.0.0.1"
	config.RedisNamespace = "test"
	config.RedisTTL = time.Minute
	jm, err := NewRedisJobManager[TestJob](config)
	assert.NoError(t, err)
	assert.NotNil(t, jm)
	jm.SetHookContext("scope", "FALSE")
	jm.SetHookContext("scope", "TRUE")

	err = jm.StartWorker(context.Background(), func() {})
	assert.NoError(t, err)

	entity := &TestJob{
		ExitCode: -1,
		StdIn:    "hello world",
	}

	jobID, err := jm.AddJob(entity)
	assert.NoError(t, err)
	assert.NotNil(t, jobID)

	time.Sleep(1 * time.Second)
	state, progress, msgID := jm.GetJobProgress(jobID, nil)
	assert.Equal(t, StateCompleted, state)
	assert.Equal(t, uint64(0), progress)
	assert.NotNil(t, msgID)

	err = jm.GetJob(jobID, entity)
	assert.NoError(t, err)
	assert.Equal(t, 0, entity.ExitCode)
	assert.Equal(t, "hello world", entity.StdOut)
	assert.Equal(t, "TRUE", onCompletionHookCalled)

	entity = &TestJob{
		ExitCode: -1,
		StdIn:    "hello universe",
	}

	err = jm.AddJobAndWait(entity)
	assert.NoError(t, err)
	assert.Equal(t, 0, entity.ExitCode)
	assert.Equal(t, "hello universe", entity.StdOut)
}
