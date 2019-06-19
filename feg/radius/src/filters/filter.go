package filters

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
)

type (
	// Filter represents a request filter action
	Filter interface {
		Init(c *config.ServerConfig) error
		Process(c *modules.RequestContext, l string, r *radius.Request) error
	}

	// FilterInitFunc type for filter's Init function
	FilterInitFunc func(c *config.ServerConfig) error

	// FilterProcessFunc type for filter's Process function
	FilterProcessFunc func(c *modules.RequestContext, l string, r *radius.Request) error
)
