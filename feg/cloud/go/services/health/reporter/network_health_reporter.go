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

package reporter

import (
	"time"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/health"
	"magma/feg/cloud/go/services/health/metrics"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
)

type NetworkHealthStatusReporter struct {
}

func (reporter *NetworkHealthStatusReporter) ReportHealthStatus(dur time.Duration) {
	for range time.Tick(dur) {
		err := reporter.reportHealthStatus()
		if err != nil {
			glog.Errorf("err in reportHealthStatus: %v\n", err)
		}
	}
}

func (reporter *NetworkHealthStatusReporter) reportHealthStatus() error {
	networks, err := configurator.ListNetworkIDs()
	if err != nil {
		return err
	}
	for _, networkID := range networks {
		config, err := configurator.LoadNetworkConfig(networkID, feg.FegNetworkType, serdes.Network)
		// Consider a FeG network to be only those that have FeG Network configs defined
		if err != nil || config == nil {
			continue
		}
		gateways, _, err := configurator.LoadEntities(
			networkID, swag.String(orc8r.MagmadGatewayType), nil, nil, nil,
			configurator.EntityLoadCriteria{},
			serdes.Entity,
		)
		if err != nil {
			glog.Errorf("error getting gateways for network %v: %v\n", networkID, err)
			continue
		}
		healthyGateways := 0
		for _, gw := range gateways {
			healthStatus, err := health.GetHealth(networkID, gw.Key)
			if err != nil {
				glog.V(2).Infof("error getting health for network %s, gateway %s: %v\n", networkID, gw.Key, err)
				continue
			}
			status, _, err := servicers.AnalyzeHealthStats(healthStatus, networkID)
			if err != nil {
				glog.V(2).Infof("error analyzing health stats for network %s, gateway %s: %v", networkID, gw.Key, err)
			}
			if status == protos.HealthStatus_HEALTHY {
				healthyGateways++
			}
		}
		metrics.TotalGatewayCount.WithLabelValues(networkID).Set(float64(len(gateways)))
		metrics.HealthyGatewayCount.WithLabelValues(networkID).Set(float64(healthyGateways))
	}
	return nil
}
