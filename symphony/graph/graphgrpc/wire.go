// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package graphgrpc

import (
	"net/http"

	"database/sql"

	"github.com/facebookincubator/symphony/pkg/actions/action/magmarebootnode"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/magmaalert"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/orc8r"
	"github.com/facebookincubator/symphony/pkg/viewer"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

// Config defines the grpc server config.
type Config struct {
	DB      *sql.DB
	Logger  log.Logger
	Orc8r   orc8r.Config
	Tenancy viewer.Tenancy
}

// NewServer creates a server from config.
func NewServer(cfg Config) (*grpc.Server, func(), error) {
	wire.Build(
		wire.FieldsOf(new(Config), "Tenancy", "DB", "Logger", "Orc8r"),
		newOrc8rClient,
		newActionsRegistry,
		newServer,
	)
	return nil, nil, nil
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
