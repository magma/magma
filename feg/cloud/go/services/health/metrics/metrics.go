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
	"magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ActiveGatewayChanged = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "active_gateway_changed_total",
			Help: "increases everytime the active gateway for a network is updated",
		},
		[]string{metrics.NetworkLabelName},
	)
	TotalGatewayCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_total_count",
			Help: "Total number of gateways that are in the network"},
		[]string{metrics.NetworkLabelName},
	)
	HealthyGatewayCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_health_count",
			Help: "Number of gateways that are healthy in the network"},
		[]string{metrics.NetworkLabelName},
	)
)

func init() {
	prometheus.MustRegister(ActiveGatewayChanged, TotalGatewayCount, HealthyGatewayCount)
}

// SetHealthyGatewayMetric takes the current health of both active and standby gateways
// in a network and sets the prometheus gauge metric for number of healthy gateways accordingly.
// Note: Prometheus gauge metric Set's are done with the atomic operation StoreUint64
func SetHealthyGatewayMetric(networkID string, gwHealth1, gwHealth2 protos.HealthStatus_HealthState) {
	if gwHealth1 == protos.HealthStatus_HEALTHY && gwHealth2 == protos.HealthStatus_HEALTHY {
		HealthyGatewayCount.WithLabelValues(networkID).Set(2)
	} else if gwHealth1 == protos.HealthStatus_UNHEALTHY && gwHealth2 == protos.HealthStatus_HEALTHY {
		HealthyGatewayCount.WithLabelValues(networkID).Set(1)
	} else if gwHealth1 == protos.HealthStatus_HEALTHY && gwHealth2 == protos.HealthStatus_UNHEALTHY {
		HealthyGatewayCount.WithLabelValues(networkID).Set(1)
	} else {
		glog.Infof("Both gateways are unhealthy in network: %s", networkID)
		HealthyGatewayCount.WithLabelValues(networkID).Set(0)
	}
}
