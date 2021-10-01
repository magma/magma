// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package metrics provides a generic metrics abstraction for Magma.
//
// Goals:
//
//   1. Allow packages to send metrics against an abstraction so that they can be
//      deployed in different environments and metrics implementations can be
//      injected. While Prometheus is used currently as the standard Magma
//			time-series database, this can be switched out for other cloud-native
//			platforms.
//
//	 2. Provide abstractions that allow for publishing metrics to a number of
//			different time-series database platforms. Initially, this will be
//			targeted against Prometheus, Azure, Google Cloud, and Amazon Timestream.
//
//	 3. Remove any implementation and imports at the AGW level of any specific
//			database platform. This should be moved to the Orc8r level.
//
//
// Examples:
//
//	// Create your metrics definitions
//	gauge := metrics.NewGauge(clock.New(), "cpu_usage", "CPU Usage %", []string{"service_name"})
//
//	// Create a metrics sender, in this case, a Prometheus metrics sender
//	// Create your gRPC client first, metricsdGrpcClient
//	sender := metrics.NewPrometheusMetricSender(metricsdGrpcClient, "agw_id")

//  // Create a metrics collector
//	collector := metrics.NewCollector(clock.New(), metrics.DefaultCollectionPeriod, sender)
//
//  // Register your metrics with the collector
//	collector.RegisterGauge(gauge)
//
//	// And now start metrics collection
//	collector.Start()
//
//	gauge.Set(1, map[string]string{"service_name": "policydb"})
//
//	// That's it! your metrics will be regularly collected and propagated to
//	// metricsd, and then to Prometheus!

package metrics

import "time"

const DefaultCollectionPeriod = 60 * time.Second
