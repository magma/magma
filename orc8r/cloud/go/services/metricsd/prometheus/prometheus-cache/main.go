/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"flag"
	"fmt"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache/cache"

	"github.com/labstack/echo"
)

const (
	defaultPort = "9091"
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	flag.Parse()

	metricCache := cache.NewMetricCache()
	e := echo.New()

	e.POST("/metrics", metricCache.Receive)
	e.GET("/metrics", metricCache.Scrape)

	e.GET("/debug", metricCache.Debug)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
