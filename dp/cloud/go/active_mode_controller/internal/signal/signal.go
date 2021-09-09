package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type app interface {
	Run(ctx context.Context) error
}

func Run(ctx context.Context, app app) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	appCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()
	return app.Run(appCtx)
}
