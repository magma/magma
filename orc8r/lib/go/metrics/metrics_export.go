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

package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	prometheus_proto "github.com/prometheus/client_model/go"
)

const (
	// CertExpiryMetric
	CertExpiresInHoursMetric = "cert_expires_in_hours"

	// NetworkTypeMetric racks different network types in a deployment.
	NetworkTypeMetric = "network_type"

	// EnodebConnectedMetric tracks the number of subscribers connected to
	// eNodeB.
	EnodebConnectedMetric = "enodeb_connected"

	// GatewayMagmaVersionMetric tracks installed magma versions on the
	// gateway.
	GatewayMagmaVersionMetric = "gateway_version"

	// ConfiguredSubscribersMetric tracks the number of subscribers explicitly
	// configured on the network
	ConfiguredSubscribersMetric = "configured_subscribers"

	// ActualSubscribersMetric tracks the actual number of subscribers which
	// have an active session on the network
	ActualSubscribersMetric = "actual_subscriber"

	// ActiveSessionAPNMetric tracks the number of active user sessions per
	// apn
	ActiveSessionAPNMetric = "active_sessions_apn"

	SubscriberThroughputMetric = "subscriber_throughput"

	/* Labels */
	NetworkLabelName         = "networkID"
	GatewayLabelName         = "gatewayID"
	EnodebLabelName          = "enodebID"
	CloudHostLabelName       = "cloudHost"
	NetworkTypeLabel         = "networkType"
	EnodeConfigTypeLabel     = "configType"
	APNLabel                 = "apnType"
	GatewayMagmaVersionLabel = "version"
	QuantileLabel            = "quantile"
	ImsiLabelName            = "IMSI"
	CertNameLabel            = "certName"
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
