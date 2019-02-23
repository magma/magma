/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	requestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests obsidian receives",
		},
	)
	respStatuses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Number of obsidian response by status",
		},
		[]string{"code", "method"},
	)
)

func init() {
	prometheus.MustRegister(requestCount, respStatuses)
}
