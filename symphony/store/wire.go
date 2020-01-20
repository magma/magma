// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package main

import (
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/server/xserver"
	"github.com/facebookincubator/symphony/store/handler"
	"github.com/facebookincubator/symphony/store/sign/s3"

	"github.com/google/wire"
	"gocloud.dev/server/health"
)

// NewServer provider a store http service from cli flags.
func NewServer(flags *cliFlags) (*server.Server, func(), error) {
	wire.Build(
		xserver.ServiceSet,
		xserver.DefaultViews,
		log.Set,
		wire.FieldsOf(new(*cliFlags), "Census", "Log", "S3"),
		wire.Value([]health.Checker(nil)),
		s3.Set,
		handler.Set,
	)
	return nil, nil, nil
}
