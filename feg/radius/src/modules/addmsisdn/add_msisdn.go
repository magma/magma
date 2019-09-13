/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package addmsisdn

import (
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	expresswifi "fbc/lib/go/radius/expresswifi"

	"go.uber.org/zap"
)

// Init module interface implementation
func Init(loggert *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	return nil, nil
}

// Handle module interface implementation
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	// Load session state
	state, err := c.SessionStorage.Get()
	if err != nil {
		c.Logger.Error(
			"Error loading session state, skipping attachment of MSISDN",
			zap.Error(err),
		)
		return nil, err
	}

	// Add MSISDN to request
	err = expresswifi.XWFMSISDN_Add(r.Packet, []byte(state.MSISDN))
	if err != nil {
		return nil, errors.New("Failed encoding MSISDN attribute: " + err.Error())
	}

	return next(c, r)
}
