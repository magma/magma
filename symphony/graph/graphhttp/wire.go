// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package graphhttp

import (
	"net/http"

	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/actions/action/magmarebootnode"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/magmaalert"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/oc"
	"github.com/facebookincubator/symphony/pkg/orc8r"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/server/xserver"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"gocloud.dev/server/health"
)

// Config defines the http server config.
type Config struct {
	Tenancy *viewer.MySQLTenancy
	Logger  log.Logger
	Census  oc.Options
	Orc8r   orc8r.Config
}

// NewServer creates a server from config.
func NewServer(cfg Config) (*server.Server, func(), error) {
	wire.Build(
		xserver.ServiceSet,
		xserver.DefaultViews,
		newHealthChecker,
		newOrc8rClient,
		newActionsRegistry,
		wire.FieldsOf(new(Config), "Tenancy", "Logger", "Census", "Orc8r"),
		newRouter,
		wire.Bind(new(http.Handler), new(*mux.Router)),
		wire.Bind(new(viewer.Tenancy), new(*viewer.MySQLTenancy)),
	)
	return nil, nil, nil
}

func newHealthChecker(tenancy *viewer.MySQLTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newOrc8rClient(config orc8r.Config) *http.Client {
	client, _ := orc8r.NewClient(config)
	return client
}

func newActionsRegistry(orc8rClient *http.Client) *executor.Registry {
	registry := executor.NewRegistry()
	registry.MustRegisterTrigger(magmaalert.New())
	registry.MustRegisterAction(magmarebootnode.New(orc8rClient))
	return registry
}
