// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package graphhttp

import (
	"context"
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
	"gocloud.dev/pubsub"
	"gocloud.dev/server/health"
)

// Config defines the http server config.
type Config struct {
	Tenancy   *viewer.MySQLTenancy
	Topic     *pubsub.Topic
	Subscribe func(context.Context) (*pubsub.Subscription, error)
	Logger    log.Logger
	Census    oc.Options
	Orc8r     orc8r.Config
}

// NewServer creates a server from config.
func NewServer(cfg Config) (*server.Server, func(), error) {
	wire.Build(
		xserver.ServiceSet,
		xserver.DefaultViews,
		newHealthChecker,
		wire.FieldsOf(new(Config), "Tenancy", "Logger", "Census"),
		newRouterConfig,
		newRouter,
		wire.Bind(new(http.Handler), new(*mux.Router)),
	)
	return nil, nil, nil
}

func newHealthChecker(tenancy *viewer.MySQLTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newRouterConfig(config Config) (cfg routerConfig, err error) {
	client, _ := orc8r.NewClient(config.Orc8r)
	registry := executor.NewRegistry()
	if err = registry.RegisterTrigger(magmaalert.New()); err != nil {
		return
	}
	if err = registry.RegisterAction(magmarebootnode.New(client)); err != nil {
		return
	}
	cfg = routerConfig{
		tenancy: config.Tenancy,
		logger:  config.Logger,
	}
	cfg.events.topic = config.Topic
	cfg.events.subscribe = config.Subscribe
	cfg.orc8r.client = client
	cfg.actions.registry = registry
	return cfg, nil
}
