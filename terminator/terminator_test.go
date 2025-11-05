package terminator_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/spacecafe/gobox/terminator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sendSigTerm(t *testing.T) {
	t.Helper()

	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)

	err = p.Signal(syscall.SIGTERM)
	require.NoError(t, err)
}

//nolint:paralleltest // This test is not safe to run in parallel.
func TestWithoutTracking(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_ = terminator.New(&terminator.Config{
			Timeout: time.Second,
			Force:   true,
		})

		// Mock os.Exit to prevent the test from exiting.
		exitCh := make(chan int)
		terminator.OsExit = func(code int) {
			exitCh <- code
		}

		sendSigTerm(t)

		// Wait for the osExit to be called.
		select {
		case code := <-exitCh:
			assert.Equal(t, terminator.ExitCodeSigTerm, code)
		case <-time.After(4 * time.Second):
			t.Fatal("Timeout waiting for os.Exit to be called")
		}
	})
}

//nolint:paralleltest // This test is not safe to run in parallel.
func TestWithTracking(t *testing.T) {
	cfg := &terminator.Config{}
	cfg.SetDefaults()
	t.Run("", func(t *testing.T) {
		term := terminator.New(&terminator.Config{
			Timeout: time.Second,
			Force:   true,
		})

		// Mock os.Exit to prevent the test from exiting.
		exitCh := make(chan int)
		terminator.OsExit = func(code int) {
			exitCh <- code
		}

		go func(ctx context.Context, done func()) {
			<-ctx.Done()
			<-time.After(time.Second)
			done()
		}(term.TrackWithDone())

		sendSigTerm(t)

		// Wait for the os.Exit to be called.
		select {
		case code := <-exitCh:
			assert.Equal(t, terminator.ExitCodeSigTerm, code)
		case <-time.After(4 * time.Second):
			t.Fatal("Timeout waiting for os.Exit to be called")
		}
	})
}
