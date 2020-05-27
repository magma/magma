// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package main

import (
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server/xserver"
	"github.com/facebookincubator/symphony/store/handler"
	"github.com/facebookincubator/symphony/store/sign/s3"
	"github.com/google/wire"
	"gocloud.dev/server/health"
)

func newApplication(flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.Struct(new(application), "*"),
		wire.FieldsOf(new(*cliFlags),
			"ListenAddress",
			"S3Config",
			"LogConfig",
			"TelemetryConfig",
		),
		log.Provider,
		xserver.ServiceSet,
		xserver.DefaultViews,
		wire.Value([]health.Checker(nil)),
		s3.Provider,
		handler.Set,
	)
	return nil, nil, nil
}
