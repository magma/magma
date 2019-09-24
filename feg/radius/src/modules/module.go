/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate protoc -I../../../../magma/feg/gateway/services/aaa/protos --go_out=plugins=grpc,paths=source_relative:./protos context.proto
//go:generate protoc -I../../../../magma/feg/gateway/services/aaa/protos --go_out=plugins=grpc,paths=source_relative:./protos eap.proto
//go:generate protoc -I../../../../magma/feg/gateway/services/aaa/protos --go_out=plugins=grpc,paths=source_relative:./protos authorization.proto
//go:generate protoc -I../../../../magma/feg/gateway/services/aaa/protos --go_out=plugins=grpc,paths=source_relative:./protos accounting.proto

package modules

import (
	"context"

	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"

	"go.uber.org/zap"
)

type (
	// ModuleConfig represents a module configuration (free form)
	ModuleConfig map[string]interface{}

	// RequestContext Info about the request and utils for the handler
	RequestContext struct {
		context.Context
		RequestID      int64
		Logger         *zap.Logger
		SessionID      string
		SessionStorage session.Storage
	}

	// Response the response of a plugin handler
	Response struct {
		Code       radius.Code
		Attributes radius.Attributes
		Raw        []byte // Optional raw version of the packet
	}

	// Middleware a middleware method. A module may "decide" not to call the
	// next middleware and just return
	Middleware func(c *RequestContext, r *radius.Request) (*Response, error)

	// Context is an instance that holds module-specific parameters
	Context interface{}

	// Module a pluggable RADIUS request handler
	Module interface {
		Init(loggert *zap.Logger, config ModuleConfig) (Context, error)
		Handle(m Context, c *RequestContext, r *radius.Request, next Middleware) (*Response, error)
	}

	// ModuleInitFunc type for module's Init function
	ModuleInitFunc func(loggert *zap.Logger, config ModuleConfig) (Context, error)

	// ModuleHandleFunc type for module's Handle function
	ModuleHandleFunc func(m Context, c *RequestContext, r *radius.Request, next Middleware) (*Response, error)
)
