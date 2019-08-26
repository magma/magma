/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package testloopback

import (
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"

	"go.uber.org/zap"
)

/*
 * This module is used for testing - it reflects all attributes from the request into the response.
 * This allows a test to check the changes made by the modules within the chain
 * its also possible to add some new attributes if a test requires them & they dont interfere with generic testing
 */

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) error {
	return nil
}

// Handle module interface implementation
func Handle(c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	logger := c.Logger.With(zap.String("module_name", "testloopback"))
	// create a response with all attributes copied from request
	resp := modules.Response{
		Code:       radius.CodeAccessAccept,
		Attributes: r.Packet.Attributes,
	}
	// add some attributes so test code can verify this module processed the packet
	a, err := radius.NewString("Eitan")
	if err != nil {
		logger.Error("failed to create RADIUS attribute to response")
	}
	resp.Attributes.Add(rfc2865.UserName_Type, a)
	logger.Debug("generating dummy response", zap.Any("dummy response", resp))
	return &resp, nil
}
