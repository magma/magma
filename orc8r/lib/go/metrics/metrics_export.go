/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	prometheus_proto "github.com/prometheus/client_model/go"
)

const (
	NetworkLabelName   = "networkID"
	GatewayLabelName   = "gatewayID"
	CloudHostLabelName = "cloudHost"
)

// GetMetrics gathers metrics from Prometheus' default registry,
// and adds a timestamp to each metric. This method is called
// in Service303 Server's GetMetrics rpc implementation.
// All servicers register their metrics with the default registry
// by calling prometheus.MustRegister().
func GetMetrics() ([]*prometheus_proto.MetricFamily, error) {

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return []*prometheus_proto.MetricFamily{},
			fmt.Errorf("err gathering from registry: %v\n", err)
	}
	// timeStamp in milliseconds
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	for _, metric_family := range families {
		for _, sample := range metric_family.Metric {
			sample.TimestampMs = &timeStamp
		}
	}
	return families, nil
}
