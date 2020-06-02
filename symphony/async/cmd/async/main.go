// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"github.com/facebookincubator/symphony/async/handler"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/telemetry"
	"gopkg.in/alecthomas/kingpin.v2"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	_ "gocloud.dev/pubsub/mempubsub"
	_ "gocloud.dev/pubsub/natspubsub"
)

type cliFlags struct {
	MySQLConfig     mysql.Config
	EventConfig     pubsub.Config
	LogConfig       log.Config
	TelemetryConfig telemetry.Config
}

func main() {
	var cf cliFlags
	kingpin.HelpFlag.Short('h')
	kingpin.Flag(
		"mysql.dsn",
		"mysql connection string",
	).
		Envar("MYSQL_DSN").
		Required().
		SetValue(&cf.MySQLConfig)
	kingpin.Flag(
		"event.pubsub-url",
		"events pubsub url",
	).
		Envar("EVENT_PUBSUB_URL").
		Default("mem://events").
		SetValue(&cf.EventConfig)
	log.AddFlagsVar(kingpin.CommandLine, &cf.LogConfig)
	telemetry.AddFlagsVar(kingpin.CommandLine, &cf.TelemetryConfig)
	kingpin.Parse()

	ctx := ctxutil.WithSignal(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	app, cleanup, err := NewApplication(ctx, &cf)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer cleanup()

	app.logger.Info("starting application")
	err = app.run(ctx)
	app.logger.Info("terminating application", zap.Error(err))
}

type application struct {
	logger *zap.Logger
	http   *server.Server
	server *handler.Server
}

func (app *application) run(ctx context.Context) error {
	g := ctxgroup.WithContext(ctx)
	g.Go(func(context.Context) error {
		err := app.http.ListenAndServe(":80")
		app.logger.Debug("server terminated", zap.Error(err))
		return err
	})
	g.Go(func(context.Context) error {
		listener, err := app.server.Subscribe(ctx)
		if err != nil {
			return fmt.Errorf("creating event listener: %w", err)
		}
		defer listener.Shutdown(ctx)
		return listener.Listen(ctx)
	})
	<-ctx.Done()

	g.Go(func(context.Context) error {
		app.logger.Debug("start server termination")
		err := app.http.Shutdown(context.Background())
		app.logger.Debug("end server termination", zap.Error(err))
		return err
	})
	return g.Wait()
}
