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
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/wrappers"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// subscriberStateLifeStateKey key exported by sessiond for identifying
	// lifecycle state of a subscriber session
	subscriberStateLifeStateKey = "lifecycle_state"

	// sessionActive string literal identifying active subscriber session
	sessionActive = "SESSION_ACTIVE"
)

type UserMetricsCalculation struct {
	calculations.BaseCalculation
}

func (x *UserMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate User Metrics")

	var results []*protos.CalculationResult
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	for _, networkID := range networks {
		subscriberEnts, _, err := configurator.LoadAllEntitiesOfType(
			networkID,
			lte.SubscriberEntityType,
			configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsFromThis: true},
			serdes.Entity)
		if err != nil {
			continue
		}

		// get all subscribers state across (configured and federated subscribers)
		subscriberStateTypes := []string{lte.SubscriberStateType}
		states, err := state.SearchStates(networkID, subscriberStateTypes, nil, nil, serdes.State)
		if err != nil {
			continue
		}
		users := make(map[string]struct{})
		var exists = struct{}{}
		activeSessionsPerAPN := make(map[string]int)
		for stateID, st := range states {
			users[stateID.DeviceID] = exists
			for k, v := range getActiveSessionsPerApn(st.ReportedState) {
				activeSessionsPerAPN[k] += v
			}
		}
		labels := prometheus.Labels{metrics.NetworkLabelName: networkID}

		if cfg, ok := x.AnalyticsConfig.Metrics[metrics.ConfiguredSubscribersMetric]; ok {
			results = append(results, calculations.NewResult(
				float64(len(subscriberEnts)),
				metrics.ConfiguredSubscribersMetric,
				calculations.CombineLabels(labels, cfg.Labels)))
		} else {
			glog.Errorf("%s metric not found in metric config", metrics.ConfiguredSubscribersMetric)
		}

		// verify if the total connected users is greater than min user threshold
		if cfg, ok := x.AnalyticsConfig.Metrics[metrics.ActualSubscribersMetric]; ok {
			results = append(results, calculations.NewResult(
				float64(len(users)),
				metrics.ActualSubscribersMetric,
				calculations.CombineLabels(labels, cfg.Labels)))
		} else {
			glog.Errorf("%s metric not found in metric config", metrics.ActualSubscribersMetric)
		}

		if cfg, ok := x.AnalyticsConfig.Metrics[metrics.ActiveSessionAPNMetric]; ok {
			for apnID, numActiveSessionsPerAPN := range activeSessionsPerAPN {
				labels := prometheus.Labels{metrics.NetworkLabelName: networkID, metrics.APNLabel: apnID}
				labels = calculations.CombineLabels(labels, cfg.Labels)
				result := calculations.NewResult(float64(numActiveSessionsPerAPN), metrics.ActiveSessionAPNMetric, labels)
				results = append(results, result)
			}
		} else {
			glog.Errorf("%s metric not found in metric config", metrics.ActiveSessionAPNMetric)
		}
	}
	glog.V(1).Info("User Metrics Results ", results)
	return results, nil
}

type SiteMetricsCalculation struct {
	calculations.BaseCalculation
}

// Calculate computes site specific calculations based on gateway and eNodeB
// state present in the orc8r
func (x *SiteMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate Site Metrics")
	var results []*protos.CalculationResult
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	gatewayVersionCfg, gatewayVersionCfgOk := x.AnalyticsConfig.Metrics[metrics.GatewayMagmaVersionMetric]
	enbConnectedCfg, enbConnectedOk := x.AnalyticsConfig.Metrics[metrics.EnodebConnectedMetric]

	for _, networkID := range networks {
		if gatewayVersionCfgOk {
			gatewayEnts, _, err := configurator.LoadAllEntitiesOfType(
				networkID,
				lte.CellularGatewayEntityType,
				configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsToThis: true},
				serdes.Entity,
			)
			if err != nil {
				continue
			}
			for _, ent := range gatewayEnts {
				status, err := wrappers.GetGatewayStatus(networkID, ent.PhysicalID)
				if err != nil || status == nil || status.PlatformInfo == nil || len(status.PlatformInfo.Packages) == 0 {
					glog.V(2).Infof("gateway %s, err %v or version not available", ent.PhysicalID, err)
					continue
				}

				gatewayVersion := ""
				for _, pkg := range status.PlatformInfo.Packages {
					if pkg.Name == "magma" {
						gatewayVersion = pkg.Version
						glog.V(2).Infof("gateway %s version %s", ent.PhysicalID, gatewayVersion)
						break
					}
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
				results = append(results, calculations.NewResult(1, metrics.GatewayMagmaVersionMetric, labels))
			}
		}

		if enbConnectedOk {
			ents, _, err := configurator.LoadAllEntitiesOfType(
				networkID,
				lte.CellularEnodebEntityType,
				configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsToThis: true},
				serdes.Entity,
			)
			if err != nil {
				continue
			}

			for _, ent := range ents {
				enodeb := (&lte_models.Enodeb{}).FromBackendModels(ent)
				if enodeb.EnodebConfig == nil || enodeb.EnodebConfig.ConfigType == "" {
					continue
				}

				st, err := state.GetState(networkID, lte.EnodebStateType, ent.Key, serdes.State)
				if err != nil {
					continue
				}
				enodebState := st.ReportedState.(*lte_models.EnodebState)
				ent, err := configurator.LoadEntityForPhysicalID(st.ReporterID, configurator.EntityLoadCriteria{}, serdes.Entity)
				if err != nil {
					continue
				}

				enodebState.ReportingGatewayID = ent.Key
				labels := prometheus.Labels{
					metrics.NetworkLabelName:     networkID,
					metrics.GatewayLabelName:     enodebState.ReportingGatewayID,
					metrics.EnodebLabelName:      enodeb.Serial,
					metrics.EnodeConfigTypeLabel: enodeb.EnodebConfig.ConfigType,
				}
				labels = calculations.CombineLabels(labels, enbConnectedCfg.Labels)
				results = append(results, calculations.NewResult(float64(enodebState.UesConnected), metrics.EnodebConnectedMetric, labels))
			}
		}
	}
	glog.V(1).Info("Site Metrics Results ", results)
	return results, nil
}

func getActiveSessionsPerApn(subscriberStateIntf interface{}) map[string]int {
	activeSessionsPerAPN := make(map[string]int)
	if subscriberStateIntf == nil {
		return activeSessionsPerAPN
	}
	subscriberState, ok := subscriberStateIntf.(*state.ArbitraryJSON)
	if !ok {
		glog.Errorf("reported state for session state having unexpected type %T", subscriberStateIntf)
		return activeSessionsPerAPN
	}

	for apnID, apnSessionStatesIntf := range *subscriberState {
		apnSessionStates, ok := apnSessionStatesIntf.([]interface{})
		if !ok {
			glog.Errorf("apnSession states got unexpected type %T", apnSessionStatesIntf)
			continue
		}
		for _, apnSessionStateIntf := range apnSessionStates {
			apnSessionState, ok := apnSessionStateIntf.(map[string]interface{})
			if !ok {
				glog.Errorf("ApnSession state got unexpected type %T", apnSessionStateIntf)
				continue
			}

			sessionStateIntf, ok := (apnSessionState)[subscriberStateLifeStateKey]
			if !ok {
				glog.Errorf("Lifecycle_state key not found in session state %v", sessionStateIntf)
				continue
			}

			sessionState, ok := sessionStateIntf.(string)
			if !ok {
				continue
			}
			if sessionState == sessionActive {
				activeSessionsPerAPN[apnID]++
			}
		}
	}
	return activeSessionsPerAPN
}
