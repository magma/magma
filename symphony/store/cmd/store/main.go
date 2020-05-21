// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	stdlog "log"
	"os"
	"syscall"

	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/oc"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/store/sign/s3"
	"go.uber.org/zap"

	"github.com/jessevdk/go-flags"
)

type cliFlags struct {
	Addr   string     `env:"ADDR" long:"addr" default:":http" description:"the address to listen on"`
	S3     s3.Config  `group:"s3" namespace:"s3" env-namespace:"S3"`
	Log    log.Config `group:"log" namespace:"log" env-namespace:"LOG"`
	Census oc.Options `group:"oc" namespace:"oc" env-namespace:"OC"`
}

func main() {
	var cf cliFlags
	if _, err := flags.Parse(&cf); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	app, cleanup, err := newApplication(&cf)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer cleanup()

	app.Info("starting application",
		zap.String("address", cf.Addr),
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
	addr   string
}

func (app *application) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	g := ctxgroup.WithContext(ctx)
	g.Go(func(context.Context) error {
		err := app.server.ListenAndServe(app.addr)
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
