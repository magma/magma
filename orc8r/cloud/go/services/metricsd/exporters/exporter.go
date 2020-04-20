/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package exporters provides an interface for converting protobuf metrics to
// timeseries datapoints and writing these datapoints to storage.
package exporters

import (
	prometheus_models "github.com/prometheus/client_model/go"
)

// Exporter exports metrics received by the metricsd servicer to a datasink.
type Exporter interface {
	// Submit metrics to the exporter.
	// This method must be thread-safe.
	Submit(metrics []MetricAndContext) error
}

// MetricAndContext wraps a metric family and metric context
type MetricAndContext struct {
	Family  *prometheus_models.MetricFamily
	Context MetricContext
}

// MetricContext provides information to the exporter about where this metric
// comes from.
type MetricContext struct {
	MetricName        string
	AdditionalContext AdditionalMetricContext
}

type AdditionalMetricContext interface {
	isExtraMetricContext()
}

type CloudMetricContext struct {
	// CloudHost is the hostname of the cloud host which originated the metric.
	CloudHost string
}

type GatewayMetricContext struct {
	NetworkID string
	GatewayID string
}

type PushedMetricContext struct {
	NetworkID string
}

func (c *CloudMetricContext) isExtraMetricContext()   {}
func (c *GatewayMetricContext) isExtraMetricContext() {}
func (c *PushedMetricContext) isExtraMetricContext()  {}
