package job_manager

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestJob struct {
	ExitCode int
	StdIn    string
	StdOut   string
	StdErr   string
}

func NewTestJob() Job {
	return &TestJob{}
}

func (r *TestJob) IsPending() bool {
	return true
}

func (r *TestJob) IsActive() bool {
	return true
}

func (r *TestJob) IsDone() bool {
	return r.ExitCode >= 0
}

func (r *TestJob) Progress() int {
	return 0
}

func (r *TestJob) Start() error {
	r.ExitCode = 0
	r.StdOut = r.StdIn
	return nil
}

func TestNewRedisJobManager(t *testing.T) {
	config := NewConfig(nil)
	config.RedisHost = "127.0.0.1"
	config.RedisNamespace = "test"
	config.RedisTTL = time.Minute
	jm, err := NewRedisJobManager(config, NewTestJob)
	assert.NoError(t, err)
	assert.NotNil(t, jm)

	err = jm.StartWorker(context.Background(), func() {})
	assert.NoError(t, err)

	job := &TestJob{
		ExitCode: -1,
		StdIn:    "hello world",
	}

	jobID, err := jm.AddJob(job)
	assert.NoError(t, err)
	assert.NotNil(t, jobID)

	time.Sleep(1 * time.Second)
	err = jm.GetJob(jobID, job)
	assert.NoError(t, err)
	assert.Equal(t, 0, job.ExitCode)
	assert.Equal(t, "hello world", job.StdOut)

	job = &TestJob{
		ExitCode: -1,
		StdIn:    "hello universe",
	}

	err = jm.AddJobAndWait(job)
	assert.NoError(t, err)
	assert.Equal(t, 0, job.ExitCode)
	assert.Equal(t, "hello universe", job.StdOut)
}
