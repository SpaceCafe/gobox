package terminator

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func SendSigTerm() error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGTERM)
}

func TestWithoutTracking(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_ = New(NewConfig())

		// Mock os.Exit to prevent the test from exiting.
		exitCh := make(chan int)
		osExit = func(code int) {
			exitCh <- code
		}

		SendSigTerm()

		// Wait for the osExit to be called.
		select {
		case code := <-exitCh:
			assert.Equal(t, code, ExitCodeSigTerm)
		case <-time.After(4 * time.Second):
			t.Fatal("Timeout waiting for osExit to be called")
		}
	})
}

func TestWithTracking(t *testing.T) {
	t.Run("", func(t *testing.T) {
		terminator := New(NewConfig())

		// Mock os.Exit to prevent the test from exiting.
		exitCh := make(chan int)
		osExit = func(code int) {
			exitCh <- code
		}

		go func(ctx context.Context, done func()) {
			<-ctx.Done()
			<-time.After(time.Second)
			done()
		}(terminator.FullTracking())

		SendSigTerm()

		// Wait for the osExit to be called.
		select {
		case code := <-exitCh:
			assert.Equal(t, code, ExitCodeSigTerm)
		case <-time.After(4 * time.Second):
			t.Fatal("Timeout waiting for osExit to be called")
		}
	})
}
