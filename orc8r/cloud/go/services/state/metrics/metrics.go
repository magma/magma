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
	"magma/orc8r/lib/go/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	gwCheckinStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_checkin_status",
			Help: "1 for checkin success, 0 for checkin failure",
		},
		[]string{metrics.NetworkLabelName, metrics.GatewayLabelName},
	)
	upGwCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_up_count",
			Help: "Number of gateways that are up in the network"},
		[]string{metrics.NetworkLabelName},
	)
	totalGwCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_total_count",
			Help: "Total number of gateways that are in the network"},
		[]string{metrics.NetworkLabelName},
	)
	gwMconfigAge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_mconfig_age",
			Help: "Age of the mconfig in the gateway in seconds",
		},
		[]string{metrics.NetworkLabelName, metrics.GatewayLabelName},
	)
)
