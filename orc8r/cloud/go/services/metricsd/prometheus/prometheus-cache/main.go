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
	"net/http"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache/cache"

	"github.com/labstack/echo"
)

const (
	defaultPort          = "9091"
	defaultLimit         = -1
	defaultScrapeTimeout = 10 // seconds
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	totalMetricsLimit := flag.Int("limit", defaultLimit, fmt.Sprintf("Limit the total metrics in the cache at one time. Will reject a push if cache is full. Default is %d which is no limit.", defaultLimit))
	scrapeTimeout := flag.Int("scrapeTimeout", defaultScrapeTimeout, fmt.Sprintf("Timeout for scrape calls. Default is %d", defaultScrapeTimeout))
	flag.Parse()

	metricCache := cache.NewMetricCache(*totalMetricsLimit, *scrapeTimeout)
	e := echo.New()

	e.POST("/metrics", metricCache.Receive)
	e.GET("/metrics", metricCache.Scrape)

	e.GET("/debug", metricCache.Debug)

	// For liveness probe
	e.GET("/", func(ctx echo.Context) error { return ctx.NoContent(http.StatusOK) })

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
