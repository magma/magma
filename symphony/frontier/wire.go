// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package main

import (
	"errors"

	"github.com/facebookincubator/symphony/frontier/handler"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/server/xserver"

	"github.com/google/wire"
	"go.opencensus.io/stats/view"
	"gocloud.dev/server/health"
)

// NewServer provider a frontier http service from cli flags.
func NewServer(flags *cliFlags) (*server.Server, func(), error) {
	wire.Build(
		xserver.ServiceSet,
		defaultViews,
		log.Set,
		wire.FieldsOf(new(*cliFlags), "KeyPairs", "Census", "Log"),
		wire.Value([]health.Checker(nil)),
		handler.Set,
		proxyTarget,
		staticTarget,
		authKey,
	)
	return nil, nil, nil
}

func proxyTarget(flags *cliFlags) handler.ProxyTarget {
	return handler.ProxyTarget(flags.ProxyTarget.URL)
}

func staticTarget(flags *cliFlags) handler.StaticTarget {
	return handler.StaticTarget(flags.InventoryTarget.URL)
}

func authKey(keys []key) ([]byte, error) {
	if len(keys) > 0 {
		return keys[0], nil
	}
	return nil, errors.New("empty key set")
}

func defaultViews() []*view.View {
	return append(xserver.DefaultViews(), handler.Views()...)
}
