// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build wireinject

package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphevents"
	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/graphhttp"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/server"
	"gocloud.dev/server/health"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

// NewApplication creates a new graph application.
func NewApplication(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(*cliFlags), "Log", "Census", "MySQL", "Event", "Orc8r"),
		log.Provider,
		newApplication,
		newTenancy,
		newHealthChecks,
		newMySQLTenancy,
		newAuthURL,
		mysql.Open,
		event.Set,
		graphhttp.NewServer,
		wire.Struct(new(graphhttp.Config), "*"),
		graphgrpc.NewServer,
		wire.Struct(new(graphgrpc.Config), "*"),
		graphevents.NewServer,
		wire.Struct(new(graphevents.Config), "*"),
	)
	return nil, nil, nil
}

func newApplication(logger log.Logger, httpServer *server.Server, grpcServer *grpc.Server, eventServer *graphevents.Server, flags *cliFlags) *application {
	var app application
	app.Logger = logger.Background()
	app.http.Server = httpServer
	app.http.addr = flags.HTTPAddress
	app.grpc.Server = grpcServer
	app.grpc.addr = flags.GRPCAddress
	app.event = eventServer
	return &app
}

func newTenancy(tenancy *viewer.MySQLTenancy, logger log.Logger, emitter event.Emitter) (viewer.Tenancy, error) {
	eventer := event.Eventer{Logger: logger, Emitter: emitter}
	return viewer.NewCacheTenancy(tenancy, eventer.HookTo), nil
}

func newHealthChecks(tenancy *viewer.MySQLTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newMySQLTenancy(dsn string, logger log.Logger) (*viewer.MySQLTenancy, error) {
	tenancy, err := viewer.NewMySQLTenancy(dsn)
	if err != nil {
		return nil, fmt.Errorf("creating mysql tenancy: %w", err)
	}
	tenancy.SetLogger(logger)
	mysql.SetLogger(logger)
	return tenancy, nil
}

func newAuthURL(flags *cliFlags) (*url.URL, error) {
	u, err := url.Parse(flags.AuthURL)
	if err != nil {
		return nil, fmt.Errorf("parsing auth url: %w", err)
	}
	return u, nil
}
