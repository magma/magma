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

type Exporter interface {
	// This method has to be thread-safe
	// Submit submits a metric to Exporter
	Submit(family *dto.MetricFamily, context MetricsContext) error

	Start()
}

// MetricsContext provides information to the exporter about where this metric
// comes from.
// OriginatingEntity - unique identifier for the originator of a metric
// DecodedName       - name of the metric family
type MetricsContext struct {
	NetworkID, GatewayID, HardwareID, OriginatingEntity, DecodedName, MetricName string
}
