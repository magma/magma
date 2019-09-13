/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package alwaysaccept

import (
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"

	"go.uber.org/zap"
)

// Init module interface implementation
func Init(loggert *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	return nil, nil
}

// Handle module interface implementation
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	if r.Code != radius.CodeAccessRequest {
		return nil, errors.New("module cannot handle anything other than Access-Request messages")
	}

	return &modules.Response{
		Code:       radius.CodeAccessAccept,
		Attributes: radius.Attributes{},
	}, nil
}
