package main

import (
	"context"
	"log"
	"os"

	"magma/dp/cloud/go/active_mode_controller/config"
	"magma/dp/cloud/go/active_mode_controller/internal/app"
	"magma/dp/cloud/go/active_mode_controller/internal/signal"
	"magma/dp/cloud/go/active_mode_controller/internal/time"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Printf("failed to read config: %s", err)
		os.Exit(1)
	}
	a := app.NewApp(
		app.WithConfig(cfg),
		app.WithClock(&time.Clock{}),
	)
	ctx := context.Background()
	if err := signal.Run(ctx, a); err != nil && err != context.Canceled {
		log.Printf("failed to stop app: %s", err)
		os.Exit(1)
	}
}
