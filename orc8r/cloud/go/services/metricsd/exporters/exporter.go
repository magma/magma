/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// The Exporter is an interface for converting protobuf Metrics into timeseries
// datapoints. It also handles writing these datapoints into storage
package exporters

import (
	"fmt"
	"strings"

	"magma/orc8r/cloud/go/protos"

	dto "github.com/prometheus/client_model/go"
	prometheus_proto "github.com/prometheus/client_model/go"
)

// Exporter exports metrics received by the metricsd servicer to a datasink.
type Exporter interface {
	// This method has to be thread-safe
	// Submit submits metrics to the exporter
	Submit(metrics []MetricAndContext) error

	Start()
}

// MetricAndContext wraps a metric family and metric context
type MetricAndContext struct {
	Family  *dto.MetricFamily
	Context MetricsContext
}

// MetricsContext provides information to the exporter about where this metric
// comes from.
// OriginatingEntity - unique identifier for the originator of a metric
// DecodedName       - name of the metric family
type MetricsContext struct {
	NetworkID, GatewayID, HardwareID, OriginatingEntity, DecodedName, MetricName string
}

func NewMetricsContext(family *prometheus_proto.MetricFamily, networkID, hardwareID, gatewayID string) MetricsContext {
	return MetricsContext{
		MetricName:        protos.GetDecodedName(family),
		NetworkID:         networkID,
		HardwareID:        hardwareID,
		GatewayID:         gatewayID,
		OriginatingEntity: networkID + "." + gatewayID,
		DecodedName:       protos.GetDecodedName(family),
	}
}

// UnpackCloudMetricName takes a "cloud" metric name and attempts to parse out
// the networkID and gatewayID from the name. Returns an error if either do not
// exist.
func UnpackCloudMetricName(metricName string) (string, string, error) {
	const (
		networkLabel = "networkId"
		gatewayLabel = "gatewayId"
	)
	networkLabelIndex := strings.Index(metricName, networkLabel)
	gatewayLabelIndex := strings.Index(metricName, gatewayLabel)
	if networkLabelIndex == -1 || gatewayLabelIndex == -1 {
		return "", "", fmt.Errorf("no gateway or network label in cloud metric: %s", metricName)
	}
	networkStart := networkLabelIndex + len(networkLabel) + 1
	gatewayStart := gatewayLabelIndex + len(gatewayLabel) + 1

	gatewayID := metricName[gatewayStart : networkLabelIndex-1]
	networkID := metricName[networkStart:]

	return networkID, gatewayID, nil
}
