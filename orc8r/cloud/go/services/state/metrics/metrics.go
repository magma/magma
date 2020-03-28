/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import (
	"magma/orc8r/lib/go/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	gwCheckinStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_checkin_status",
			Help: "1 for checkin success, 0 for checkin failure",
		},
		[]string{metrics.NetworkLabelName, metrics.GatewayLabelName},
	)
	upGwCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_up_count",
			Help: "Number of gateways that are up in the network"},
		[]string{metrics.NetworkLabelName},
	)
	totalGwCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_total_count",
			Help: "Total number of gateways that are in the network"},
		[]string{metrics.NetworkLabelName},
	)
	gwMconfigAge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_mconfig_age",
			Help: "Age of the mconfig in the gateway in seconds",
		},
		[]string{metrics.NetworkLabelName, metrics.GatewayLabelName},
	)
)

func init() {
	prometheus.MustRegister(
		gwCheckinStatus,
		upGwCount,
		totalGwCount,
		gwMconfigAge,
	)
}
