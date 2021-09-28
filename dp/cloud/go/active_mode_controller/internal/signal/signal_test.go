package signal_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"magma/dp/cloud/go/active_mode_controller/internal/signal"

	"github.com/stretchr/testify/assert"
)

const (
	timeout = time.Millisecond * 10
)

func TestRun(t *testing.T) {
	app, errs := whenAppWasStarted()
	thenAppWasStarted(t, app)

	whenSignalWasSent(t)
	thenAppWasShutdown(t, errs)
}

func whenAppWasStarted() (*stubApp, chan error) {
	errs := make(chan error)
	app := &stubApp{running: make(chan struct{})}
	go func() {
		errs <- signal.Run(context.Background(), app)
	}()
	return app, errs
}

func thenAppWasStarted(t *testing.T, s *stubApp) {
	select {
	case <-s.running:
		return
	case <-time.After(timeout):
		assert.Fail(t, "app failed to start")
	}
}

func whenSignalWasSent(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)
	assert.NoError(t, proc.Signal(syscall.SIGTERM))
}

func thenAppWasShutdown(t *testing.T, errs chan error) {
	select {
	case err := <-errs:
		assert.NoError(t, err)
	case <-time.After(timeout):
		assert.Fail(t, "app failed to stop")
	}
}

type stubApp struct {
	running chan struct{}
}

func (s *stubApp) Run(ctx context.Context) error {
	close(s.running)
	<-ctx.Done()
	return nil
}
