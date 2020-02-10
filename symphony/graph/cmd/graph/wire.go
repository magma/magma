// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build wireinject

package main

import (
	"fmt"

	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/graphhttp"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/server"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

// NewApplication creates a new graph application.
func NewApplication(flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(*cliFlags), "Log", "Census", "MySQL", "Orc8r"),
		log.Set,
		newApplication,
		newTenancy,
		mysql.Open,
		graphhttp.NewServer,
		wire.Struct(new(graphhttp.Config), "*"),
		graphgrpc.NewServer,
		wire.Struct(new(graphgrpc.Config), "*"),
	)
	return nil, nil, nil
}

func newApplication(logger log.Logger, httpServer *server.Server, grpcServer *grpc.Server, flags *cliFlags) *application {
	var app application
	app.Logger = logger.Background()
	app.http.Server = httpServer
	app.http.addr = flags.HTTPAddress
	app.grpc.Server = grpcServer
	app.grpc.addr = flags.GRPCAddress
	return &app
}

func newTenancy(logger log.Logger, dsn string) (*viewer.MySQLTenancy, error) {
	tenancy, err := viewer.NewMySQLTenancy(dsn)
	if err != nil {
		return nil, fmt.Errorf("creating mysql tenancy: %w", err)
	}
	mysql.SetLogger(logger)
	return tenancy, nil
}
