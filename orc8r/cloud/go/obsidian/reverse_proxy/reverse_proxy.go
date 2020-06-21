/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package reverse_proxy

import (
	"fmt"
	"net"
	"net/url"

	"magma/orc8r/cloud/go/obsidian"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// ReverseProxy is a middleware function to properly route requests based off
// of the request's remote address
func ReverseProxy(next echo.HandlerFunc) echo.HandlerFunc {
	urlString := fmt.Sprintf("http://localhost:%d", obsidian.Port)
	localObsidianURL, err := url.Parse(urlString)
	if err != nil {
		glog.Error(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: localObsidianURL,
		},
	}
	reverseProxyMiddleware := middleware.Proxy(middleware.NewRoundRobinBalancer(targets))
	return func(c echo.Context) error {
		if err != nil {
			return err
		}
		if c.RealIP() == "127.0.0.1" {
			return next(c)
		}
		// Re-route non-localhost requests to localhost:<obsidian_port>
		// This is an intermediate implementation state needed to de-compose
		// the Orchestrator plugin
		remoteAddr := c.Request().RemoteAddr
		_, p, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			return err
		}
		updatedRemoteAddr := fmt.Sprintf("127.0.0.1:%s", p)
		c.Request().RemoteAddr = updatedRemoteAddr
		reverseProxyHandler := reverseProxyMiddleware(next)
		return reverseProxyHandler(c)
	}
}
