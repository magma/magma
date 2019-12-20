// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package graphgrpc

import (
	"database/sql"

	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

// Config defines the grpc server config.
type Config struct {
	DB     *sql.DB
	Logger log.Logger
}

// NewServer creates a server from config.
func NewServer(cfg Config) (*grpc.Server, func(), error) {
	wire.Build(
		wire.FieldsOf(new(Config), "DB", "Logger"),
		newServer,
	)
	return nil, nil, nil
}
