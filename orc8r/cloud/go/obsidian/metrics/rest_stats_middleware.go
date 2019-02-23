/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import (
	"strconv"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

// CollectStats is the middleware function
func CollectStats(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		requestCount.Inc()
		status := strconv.Itoa(c.Response().Status)
		respStatuses.WithLabelValues(status, c.Request().Method).Inc()
		glog.V(2).Infof(
			"REST API code: %v, method: %v, url: %v\n",
			status,
			c.Request().Method,
			c.Request().URL,
		)
		return nil
	}
}
