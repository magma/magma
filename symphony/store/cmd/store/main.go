// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	stdlog "log"
	"net"
	"os"
	"syscall"

	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/facebookincubator/symphony/store/sign/s3"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cliFlags struct {
	ListenAddress   *net.TCPAddr
	S3Config        s3.Config
	LogConfig       log.Config
	TelemetryConfig telemetry.Config
}

func main() {
	var cf cliFlags
	kingpin.HelpFlag.Short('h')
	kingpin.Flag(
		"web.listen-address",
		"Address to listen on.",
	).
		Default(":http").
		TCPVar(&cf.ListenAddress)
	s3.AddFlagsVar(kingpin.CommandLine, &cf.S3Config)
	log.AddFlagsVar(kingpin.CommandLine, &cf.LogConfig)
	telemetry.AddFlagsVar(kingpin.CommandLine, &cf.TelemetryConfig)
	kingpin.Parse()

	app, cleanup, err := newApplication(&cf)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer cleanup()

	app.Info("starting application",
		zap.Stringer("address", cf.ListenAddress),
	)
	err = app.run(
		ctxutil.WithSignal(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		),
	)
	app.Info("terminating application", zap.Error(err))
}

type application struct {
	*zap.Logger
	server *server.Server
	addr   *net.TCPAddr
}

func (app *application) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	g := ctxgroup.WithContext(ctx)
	g.Go(func(context.Context) error {
		err := app.server.ListenAndServe(app.addr.String())
		app.Debug("server terminated", zap.Error(err))
		return err
	})
	g.Go(func(ctx context.Context) error {
		defer cancel()
		<-ctx.Done()
		return nil
	})
	<-ctx.Done()

	app.Warn("start application termination",
		zap.NamedError("reason", ctx.Err()),
	)
	defer app.Debug("end application termination")

	g.Go(func(context.Context) error {
		app.Debug("start server termination")
		err := app.server.Shutdown(context.Background())
		app.Debug("end server termination", zap.Error(err))
		return err
	})
	return g.Wait()
}
