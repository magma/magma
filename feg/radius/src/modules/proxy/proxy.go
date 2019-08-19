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

var target string

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) error {
	var proxyConfig Config
	err := mapstructure.Decode(config, &proxyConfig)
	if err != nil {
		return err
	}

	if proxyConfig.Target == "" {
		return errors.New("proxy module cannot be initialize with empty Target value")
	}

	target = proxyConfig.Target
	return nil
}

// Handle module interface implementation
func Handle(_ *modules.RequestContext, r *radius.Request, _ modules.Middleware) (*modules.Response, error) {
	res, err := radius.Exchange(context.Background(), r.Packet, target)
	if err != nil {
		return nil, err
	}

	return &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
	}, nil
}
