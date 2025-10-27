package job_manager

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/spacecafe/gobox/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

func (r *TestJob) OnCompletion(ctx HookContext) {
	onCompletionHookCalled = ctx.Get("scope").(string)
}

func setupRedisWithJSON(t *testing.T) (host string, port int, term func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis/redis-stack-server:latest", // Includes RedisJSON
		ExposedPorts: []string{"6379/tcp"},
		Env: map[string]string{
			"REDIS_ARGS": "--notify-keyspace-events AKE",
		},
		WaitingFor: wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, _ = redisC.Host(ctx)
	port_, _ := redisC.MappedPort(ctx, "6379")

	return host, port_.Int(), func() {
		err := redisC.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestNewRedisManager(t *testing.T) {
	var term func()
	cfg := &Config{}
	cfg.SetDefaults()
	cfg.WorkerName = "worker-1"
	cfg.RedisHost, cfg.RedisPort, term = setupRedisWithJSON(t)
	cfg.RedisNamespace = "test"
	cfg.Timeout = 5 * time.Second

	defer term()

	manager, err := NewRedisManager[TestJob](cfg, logger.Default())
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	manager.SetHookContext("scope", "FALSE")
	manager.SetHookContext("scope", "TRUE")

	err = manager.StartWorker(context.Background(), func() {})
	assert.NoError(t, err)

	// Wait a bit for goroutines to start
	time.Sleep(100 * time.Millisecond)

	// Prepare test job
	entity := &TestJob{
		ExitCode: -1,
		StdIn:    "hello world",
	}

	jobID, err := manager.AddJob(entity)
	assert.NoError(t, err)
	assert.NotNil(t, jobID)

	time.Sleep(2 * time.Second)
	state, progress, msgID := manager.GetJobProgress(jobID, nil, time.Second)
	assert.Equal(t, StateCompleted, state)
	assert.Equal(t, uint64(0), progress)
	assert.NotNil(t, msgID)
	assert.Equal(t, "TRUE", onCompletionHookCalled)

	err = manager.GetJob(jobID, entity)
	assert.NoError(t, err)
	assert.Equal(t, 0, entity.ExitCode)
	assert.Equal(t, "hello world", entity.StdOut)
}
