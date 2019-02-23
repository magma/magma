/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import "github.com/prometheus/client_golang/prometheus"

// Prometheus counters are monotonically increasing
// Counters reset to zero on service restart
var (
	//These metrics are intended as an example
	TotalRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests",
	})
	RequestFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failures_total",
		Help: "Total number of request failures",
	})
)

func init() {
	prometheus.MustRegister(TotalRequests, RequestFailures)
}
