/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package coafixedip

import (
	"context"
	"net"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"go.uber.org/zap"
)

// Config config has only one parameter which is the ip to forward the request
type Config struct {
	Target string
}

var target string

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) error {
	var coaConfig Config
	err := mapstructure.Decode(config, &coaConfig)
	if err != nil {
		return err
	}

	if coaConfig.Target == "" {
		return errors.New("coa module cannot be initialized with empty target value")
	}

	// Validating the correctness of Target
	var host string
	host, _, err = net.SplitHostPort(coaConfig.Target)
	if err != nil {
		return err
	}

	if nil == net.ParseIP(host) {
		return errors.Wrap(err, "Invalid ip address specified")
	}

	target = coaConfig.Target
	return nil
}

// Handle module interface implementation
func Handle(c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	c.Logger.Debug("Starting to handle coa radius request")
	requestCode := r.Code
	// Checking that we have received a coa request
	if requestCode != radius.CodeDisconnectRequest && requestCode != radius.CodeCoARequest {
		return next(c, r)
	}

	// Handling the coa request
	res, err := radius.Exchange(context.Background(), r.Packet, target)

	if err != nil {
		return nil, err
	}

	return &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
	}, nil
}
