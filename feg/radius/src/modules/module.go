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
	}

	// Middleware a middleware method. A module may "decide" not to call the
	// next middleware and just return
	Middleware func(c *RequestContext, r *radius.Request) (*Response, error)

	// Module a pluggable RADIUS request handler
	Module interface {
		Init(loggert *zap.Logger, config ModuleConfig) error
		Handle(c *RequestContext, r *radius.Request, next Middleware) (*Response, error)
	}

	// ModuleInitFunc type for module's Init function
	ModuleInitFunc func(loggert *zap.Logger, config ModuleConfig) error

	// ModuleHandleFunc type for module's Handle function
	ModuleHandleFunc func(c *RequestContext, r *radius.Request, next Middleware) (*Response, error)
)
