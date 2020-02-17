// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/oc"
	"github.com/facebookincubator/symphony/pkg/orc8r"
	"github.com/facebookincubator/symphony/pkg/server"

	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"
	_ "gocloud.dev/pubsub/awssnssqs"
	_ "gocloud.dev/pubsub/mempubsub"
)

type cliFlags struct {
	HTTPAddress string       `env:"HTTP_ADDRESS" long:"http-address" default:":http" description:"the http address to listen on"`
	GRPCAddress string       `env:"GRPC_ADDRESS" long:"grpc-address" default:":https" description:"the grpc address to listen on"`
	MySQL       string       `env:"MYSQL_DSN" long:"mysql-dsn" description:"connection string to mysql"`
	PubSubURL   string       `env:"PUBSUB_URL" long:"pubsub-url" default:"mem://events" description:"events pubsub topic"`
	Log         log.Config   `group:"log" namespace:"log" env-namespace:"LOG"`
	Census      oc.Options   `group:"oc" namespace:"oc" env-namespace:"OC"`
	Orc8r       orc8r.Config `group:"orc8r" namespace:"orc8r" env-namespace:"ORC8R"`
}

func main() {
	var cf cliFlags
	if _, err := flags.Parse(&cf); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	ctx := ctxutil.WithSignal(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	app, _, err := NewApplication(ctx, &cf)
	if err != nil {
		stdlog.Fatal(err)
	}

	app.Info("starting application",
		zap.String("http", cf.HTTPAddress),
		zap.String("grpc", cf.GRPCAddress),
	)
	err = app.run(ctx)
	app.Info("terminating application", zap.Error(err))
}

type application struct {
	*zap.Logger
	http struct {
		*server.Server
		addr string
	}
	grpc struct {
		*grpc.Server
		addr string
	}
}

func (app *application) run(ctx context.Context) error {
	g := ctxgroup.WithContext(ctx)
	g.Go(func(context.Context) error {
		err := app.http.ListenAndServe(app.http.addr)
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("starting http server: %w", err)
		}
		return nil
	})
	g.Go(func(context.Context) error {
		lis, err := net.Listen("tcp", app.grpc.addr)
		if err != nil {
			return fmt.Errorf("creating grpc listener: %w", err)
		}
		if err = app.grpc.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			return fmt.Errorf("starting grpc server: %w", err)
		}
		return nil
	})
	<-ctx.Done()

	g.Go(func(context.Context) error {
		app.grpc.GracefulStop()
		return nil
	})
	g.Go(func(context.Context) error {
		_ = app.http.Shutdown(context.Background())
		return nil
	})
	return g.Wait()
}
