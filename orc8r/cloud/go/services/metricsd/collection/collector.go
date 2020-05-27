/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package collection

import "github.com/prometheus/client_model/go"

// MetricCollector provides an API to query for metrics.
type MetricCollector interface {
	// GetMetrics returns a collection of prometheus MetricFamily structures
	// which contain collected metrics.
	GetMetrics() ([]*io_prometheus_client.MetricFamily, error)
}
