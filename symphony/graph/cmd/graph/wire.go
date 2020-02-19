// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build wireinject

package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/graphhttp"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/server"

	"github.com/google/wire"
	"gocloud.dev/pubsub"
	"google.golang.org/grpc"
)

// NewApplication creates a new graph application.
func NewApplication(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(*cliFlags), "Log", "Census", "MySQL", "Orc8r"),
		log.Set,
		newApplication,
		newTenancy,
		newAuthURL,
		newTopic,
		newSubscribeFunc,
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

func newAuthURL(flags *cliFlags) (*url.URL, error) {
	u, err := url.Parse(flags.AuthURL)
	if err != nil {
		return nil, fmt.Errorf("parsing auth url: %w", err)
	}
	return u, nil
}

func newTopic(ctx context.Context, flags *cliFlags) (*pubsub.Topic, func(), error) {
	topic, err := pubsub.OpenTopic(ctx, flags.PubSubURL)
	if err != nil {
		return nil, nil, fmt.Errorf("opening events topic: %w", err)
	}
	return topic, func() { topic.Shutdown(ctx) }, nil
}

func newSubscribeFunc(flags *cliFlags) func(context.Context) (*pubsub.Subscription, error) {
	return func(ctx context.Context) (*pubsub.Subscription, error) {
		subscription, err := pubsub.OpenSubscription(ctx, flags.PubSubURL)
		if err != nil {
			return nil, fmt.Errorf("opening events subscription: %w", err)
		}
		return subscription, nil
	}
}
