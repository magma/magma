// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package graphevents

import (
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/google/wire"
)

// Config defines the events server config.
type Config struct {
	Tenancy    viewer.Tenancy
	Subscriber event.Subscriber
	Logger     log.Logger
}

// NewServer creates a server from config.
func NewServer(cfg Config) (*Server, func(), error) {
	wire.Build(
		wire.FieldsOf(new(Config), "Logger", "Tenancy", "Subscriber"),
		newServerConfig,
		newServer,
	)
	return nil, nil, nil
}

func newServerConfig(tenancy viewer.Tenancy, logger log.Logger, subscriber event.Subscriber) (cfg serverConfig, err error) {
	cfg = serverConfig{
		tenancy:    tenancy,
		logger:     logger,
		subscriber: subscriber,
	}
	return cfg, nil
}
