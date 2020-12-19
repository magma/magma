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
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/wrappers"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	/* newly computed metrics */

	// NetworkTypeMetric - provides information on different network types in a deployment
	NetworkTypeMetric = "network_type"

	// EnodebConnectedMetric - number of subscribers connected to the eNodeB
	EnodebConnectedMetric = "enodeb_connected"

	// GatewayMagmaVersionMetric - provides information on gateway versions installed in a deployment
	GatewayMagmaVersionMetric = "gateway_version"

	// ConfiguredSubscribersMetric - provides the count of configured subscribers in the network
	ConfiguredSubscribersMetric = "configured_subscribers_count"

	// ActualSubscribersMetric - Number of subscribers have some session state
	ActualSubscribersMetric = "actual_subscriber_count"

	// ActiveSessionAPNMetric - Number of active user sessions in a apn
	ActiveSessionAPNMetric = "active_sessions_apn_count"

	/* labels */

	// NetworkTypeLabel - label identifying if the network type is LTE, FEG_LTE, FEG
	NetworkTypeLabel = "networkType"

	// EnodeConfigTypeLabel - label identifying if the enode is a managed or unmanaged enodeb
	EnodeConfigTypeLabel = "configType"

	// APNLabel - label identifying APN
	APNLabel = "apnType"

	// GatewayMagmaVersionLabel - label identifying the current running magma version on the gateway
	GatewayMagmaVersionLabel = "version"

	/* misc string literals */

	// SubscriberStateLifeStateKey key exported by sessiond for identifying lifecycle state of a subscriber session
	SubscriberStateLifeStateKey = "lifecycle_state"

	// SessionActive string literal identifying active subscriber session
	SessionActive = "SESSION_ACTIVE"
)

// GeneralMetricsCalculation ...
type GeneralMetricsCalculation struct {
	calculations.BaseCalculation
}

// Calculate this mainly computes the number of LTE, FEG_LTE and FEG networks
// in a particular deployment
func (x *GeneralMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate Generic Metrics")

	results := []*protos.CalculationResult{}
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	metricConfig, ok := x.AnalyticsConfig.Metrics[NetworkTypeMetric]
	if !ok {
		glog.Errorf("%s metric not found in metric config", NetworkTypeMetric)
		return results, err
	}

	for _, networkID := range networks {
		network, err := configurator.LoadNetwork(networkID, true, true, serdes.Network)
		if err == merrors.ErrNotFound {
			glog.Errorf("network %s not found", networkID)
			continue
		}
		if err != nil {
			glog.Errorf("Failed %v loading network %s", err, networkID)
			continue
		}
		ret := (&models.Network{}).FromConfiguratorNetwork(network)
		labels := prometheus.Labels{
			metrics.NetworkLabelName: networkID,
			NetworkTypeLabel:         string(ret.Type),
		}
		results = append(results,
			calculations.NewResult(1,
				NetworkTypeMetric,
				calculations.CombineLabels(labels, metricConfig.Labels)))
	}

	glog.V(1).Info("Generic Metrics Results ", results)
	return results, nil
}

// UserMetricsCalculation ...
type UserMetricsCalculation struct {
	calculations.BaseCalculation
}

// Calculate site metrics
func (x *UserMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate User Metrics")

	results := []*protos.CalculationResult{}
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	for _, networkID := range networks {
		subscriberEnts, err := configurator.LoadAllEntitiesOfType(
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

		if cfg, ok := x.AnalyticsConfig.Metrics[ConfiguredSubscribersMetric]; ok {
			results = append(results, calculations.NewResult(
				float64(len(subscriberEnts)),
				ConfiguredSubscribersMetric,
				calculations.CombineLabels(labels, cfg.Labels)))
		} else {
			glog.Errorf("%s metric not found in metric config", ConfiguredSubscribersMetric)
		}

		// verify if the total connected users is greater than min user threshold
		if cfg, ok := x.AnalyticsConfig.Metrics[ActualSubscribersMetric]; ok {
			results = append(results, calculations.NewResult(
				float64(len(users)),
				ActualSubscribersMetric,
				calculations.CombineLabels(labels, cfg.Labels)))
		} else {
			glog.Errorf("%s metric not found in metric config", ActualSubscribersMetric)
		}

		if cfg, ok := x.AnalyticsConfig.Metrics[ActiveSessionAPNMetric]; ok {
			for apnID, numActiveSessionsPerAPN := range activeSessionsPerAPN {
				labels := prometheus.Labels{metrics.NetworkLabelName: networkID, APNLabel: apnID}
				labels = calculations.CombineLabels(labels, cfg.Labels)
				result := calculations.NewResult(float64(numActiveSessionsPerAPN), ActiveSessionAPNMetric, labels)
				results = append(results, result)
			}
		} else {
			glog.Errorf("%s metric not found in metric config", ActiveSessionAPNMetric)
		}
	}
	glog.V(1).Info("User Metrics Results ", results)
	return results, nil
}

// SiteMetricsCalculation ...
type SiteMetricsCalculation struct {
	calculations.BaseCalculation
}

// Calculate site metrics
func (x *SiteMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate Site Metrics")
	results := []*protos.CalculationResult{}
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	gatewayVersionCfg, gatewayVersionCfgOk := x.AnalyticsConfig.Metrics[GatewayMagmaVersionMetric]
	enbConnectedCfg, enbConnectedOk := x.AnalyticsConfig.Metrics[EnodebConnectedMetric]

	for _, networkID := range networks {
		// load entities of gateway type
		if gatewayVersionCfgOk {
			gatewayEnts, err := configurator.LoadAllEntitiesOfType(
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
					continue
				}

				gatewayVersion := ""
				for _, pkg := range status.PlatformInfo.Packages {
					if pkg.Name == "magma" {
						gatewayVersion = pkg.Version
						break
					}
				}
				if gatewayVersion == "" {
					continue
				}

				labels := prometheus.Labels{
					metrics.NetworkLabelName: networkID,
					metrics.GatewayLabelName: ent.PhysicalID,
					GatewayMagmaVersionLabel: gatewayVersion,
				}
				labels = calculations.CombineLabels(labels, gatewayVersionCfg.Labels)
				results = append(results, calculations.NewResult(1, GatewayMagmaVersionMetric, labels))
			}
		}

		if enbConnectedOk {
			ents, err := configurator.LoadAllEntitiesOfType(
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
					metrics.NetworkLabelName: networkID,
					metrics.GatewayLabelName: enodebState.ReportingGatewayID,
					metrics.EnodebLabelName:  enodeb.Serial,
					EnodeConfigTypeLabel:     enodeb.EnodebConfig.ConfigType,
				}
				labels = calculations.CombineLabels(labels, enbConnectedCfg.Labels)

				// TODO replace 1 -> enodebState.UesConnected
				results = append(results, calculations.NewResult(1, EnodebConnectedMetric, labels))
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

			sessionStateIntf, ok := (apnSessionState)[SubscriberStateLifeStateKey]
			if !ok {
				glog.Errorf("Lifecycle_state key not found in session state %v", sessionStateIntf)
				continue
			}

			sessionState, ok := sessionStateIntf.(string)
			if !ok {
				continue
			}
			if sessionState == SessionActive {
				activeSessionsPerAPN[apnID]++
			}
		}
	}
	return activeSessionsPerAPN
}
