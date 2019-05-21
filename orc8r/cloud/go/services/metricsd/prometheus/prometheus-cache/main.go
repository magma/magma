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
	defaultPort          = "9091"
	defaultQueueCapacity = 5
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	queueCapacity := flag.Int("queue-capacity", defaultQueueCapacity, fmt.Sprintf("Maximum number of datapoints per unique series stored in cache. Default is %d\n", defaultQueueCapacity))
	flag.Parse()

	metricCache := cache.NewMetricCache(*queueCapacity)
	e := echo.New()

	e.POST("/metrics", metricCache.Receive)
	e.GET("/metrics", metricCache.Scrape)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
