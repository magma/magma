/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package calculations

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state/wrappers"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/metrics"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

type NetworkMetricsCalculation struct {
	calculations.BaseCalculation
}

func (x *NetworkMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate Network Metrics")

	var results []*protos.CalculationResult
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	metricConfig, ok := x.AnalyticsConfig.Metrics[metrics.NetworkTypeMetric]
	if !ok {
		glog.Errorf("%s metric not found in metric config", metrics.NetworkTypeMetric)
		return results, err
	}

	for _, networkID := range networks {
		network, err := configurator.LoadNetwork(networkID, true, true, serdes.Network)
		if err == merrors.ErrNotFound {
			glog.Errorf("Network %s not found", networkID)
			continue
		}
		if err != nil {
			glog.Errorf("Failed %v loading network %s", err, networkID)
			continue
		}
		ret := (&models.Network{}).FromConfiguratorNetwork(network)
		labels := prometheus.Labels{
			metrics.NetworkLabelName: networkID,
			metrics.NetworkTypeLabel: string(ret.Type),
		}
		labels = calculations.CombineLabels(labels, metricConfig.Labels)
		result := calculations.NewResult(1, metrics.NetworkTypeMetric, labels)
		results = append(results, result)
		glog.V(1).Info(result)
	}

	return results, nil
}

type SiteMetricsCalculation struct {
	calculations.BaseCalculation
}

// Calculate computes site specific calculations based on gateway state present in the orc8r
func (x *SiteMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate Site Metrics")
	var results []*protos.CalculationResult

	gatewayVersionCfg, gatewayVersionCfgOk := x.AnalyticsConfig.Metrics[metrics.GatewayMagmaVersionMetric]
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil || !gatewayVersionCfgOk {
		return results, err
	}

	for _, networkID := range networks {
		gatewayEnts, _, err := configurator.LoadEntities(
			networkID,
			swag.String(orc8r.MagmadGatewayType),
			nil,
			nil,
			nil,
			configurator.EntityLoadCriteria{},
			serdes.Entity,
		)
		if err != nil {
			continue
		}
		for _, ent := range gatewayEnts {
			status, err := wrappers.GetGatewayStatus(networkID, ent.PhysicalID)
			if err != nil ||
				status == nil ||
				status.PlatformInfo == nil ||
				len(status.PlatformInfo.Packages) == 0 {
				glog.V(2).Infof("gateway %s, err %v or version not available",
					ent.PhysicalID, err)
				continue
			}

			gatewayVersion := ""
			for _, pkg := range status.PlatformInfo.Packages {
				if pkg.Name != "magma" {
					continue
				}
				gatewayVersion = pkg.Version
				break
			}

			if gatewayVersion == "" {
				glog.V(2).Infof("gateway %s, version not found", ent.PhysicalID)
				continue
			}

			labels := prometheus.Labels{
				metrics.NetworkLabelName:         networkID,
				metrics.GatewayLabelName:         ent.PhysicalID,
				metrics.GatewayMagmaVersionLabel: gatewayVersion,
			}
			labels = calculations.CombineLabels(labels, gatewayVersionCfg.Labels)
			result := calculations.NewResult(1, metrics.GatewayMagmaVersionMetric, labels)
			results = append(results, result)
			glog.V(1).Info(result)
		}
	}

	return results, nil
}
