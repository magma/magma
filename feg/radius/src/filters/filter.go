/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

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
