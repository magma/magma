/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package proxy

import (
	"context"
	"errors"

	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// Config configuration structure for proxy module
type Config struct {
	Target string
}

// ModuleCtx ...
type ModuleCtx struct {
	target string
}

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var proxyConfig Config
	err := mapstructure.Decode(config, &proxyConfig)
	if err != nil {
		return nil, err
	}

	if proxyConfig.Target == "" {
		return nil, errors.New("proxy module cannot be initialize with empty Target value")
	}

	return ModuleCtx{target: proxyConfig.Target}, nil
}

// Handle module interface implementation
func Handle(m modules.Context, _ *modules.RequestContext, r *radius.Request, _ modules.Middleware) (*modules.Response, error) {
	mCtx := m.(ModuleCtx)
	res, err := radius.Exchange(context.Background(), r.Packet, mCtx.target)
	if err != nil {
		return nil, err
	}

	return &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
	}, nil
}
