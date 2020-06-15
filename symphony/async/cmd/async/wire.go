// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/async/handler"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/server/xserver"
	"github.com/facebookincubator/symphony/pkg/viewer"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
	"gocloud.dev/server/health"
)

// NewApplication creates a new graph application.
func NewApplication(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(*cliFlags),
			"MySQLConfig",
			"EventConfig",
			"LogConfig",
			"TelemetryConfig",
			"TenancyConfig",
		),
		log.Provider,
		newTenancy,
		newMySQLTenancy,
		newHealthChecks,
		mux.NewRouter,
		wire.Bind(new(http.Handler), new(*mux.Router)),
		xserver.ServiceSet,
		provideViews,
		pubsub.Set,
		wire.Struct(new(handler.Config), "*"),
		handler.NewServer,
		newApplication,
	)
	return nil, nil, nil
}

func newApplication(server *handler.Server, http *server.Server, logger *zap.Logger) *application {
	var app application
	app.logger = logger
	app.server = server
	app.http = http
	return &app
}

func newTenancy(tenancy *viewer.MySQLTenancy, logger log.Logger, emitter pubsub.Emitter) (viewer.Tenancy, error) {
	return viewer.NewCacheTenancy(tenancy, nil), nil
}

func newHealthChecks(tenancy *viewer.MySQLTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newMySQLTenancy(mySQLConfig mysql.Config, tenancyConfig viewer.Config, logger log.Logger) (*viewer.MySQLTenancy, error) {
	tenancy, err := viewer.NewMySQLTenancy(mySQLConfig.String(), tenancyConfig.TenantMaxConn)
	if err != nil {
		return nil, fmt.Errorf("creating mysql tenancy: %w", err)
	}
	tenancy.SetLogger(logger)
	mysql.SetLogger(logger)
	return tenancy, nil
}

func provideViews() []*view.View {
	views := xserver.DefaultViews()
	views = append(views, mysql.DefaultViews...)
	views = append(views, pubsub.DefaultViews...)
	return views
}
