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
	dto "github.com/prometheus/client_model/go"
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
