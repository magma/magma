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
	MARRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mar_requests_total",
		Help: "Total number of MAR requests sent to HSS",
	})
	MARSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mar_send_failures_total",
		Help: "Total number of MAR requests that failed to send to HSS",
	})
	SARRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sar_requests_total",
		Help: "Total number of SAR requests sent to HSS",
	})
	SARSendFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sar_send_failures_total",
		Help: "Total number of SAR requests that failed to send to HSS",
	})
	SwxTimeouts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_timeouts_total",
		Help: "Total number of swx timeouts",
	})
	SwxUnparseableMsg = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_unparseable_msg_total",
		Help: "Total number of swx messages received that cannot be parsed",
	})
	SwxInvalidSessions = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "swx_invalid_sessions_total",
		Help: "Total number of swx responses received with invalid sids",
	})
	SwxResultCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "swx_result_codes",
			Help: "swx accumulated result codes",
		},
		[]string{"code"},
	)
	SwxExperimentalResultCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "swx_experimental_result_codes",
			Help: "swx accumulated experimental result codes",
		},
		[]string{"code"},
	)
)

func init() {
	prometheus.MustRegister(MARRequests, MARSendFailures, SARRequests,
		SARSendFailures, SwxTimeouts, SwxUnparseableMsg, SwxInvalidSessions,
		SwxResultCodes, SwxExperimentalResultCodes)
}
