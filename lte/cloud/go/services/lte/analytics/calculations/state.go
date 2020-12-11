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
	models1 "magma/orc8r/cloud/go/models"
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

	// NetworkCountMetric - provides the count of different network types present in the current instance of orc8r
	NetworkCountMetric = "network_count"

	// EnodebConnectedMetric - provides the count of enodebs connected in the network along with their gateway labels
	EnodebConnectedMetric = "enodeb_connected_count"

	// EnodeConfigTypeMetric - provides the count of unmanaged and managed enodebs
	EnodeConfigTypeMetric = "enodeb_configtype_count"

	// GatewayMagmaVersionMetric - provides count of gateways running with specific magma versions
	GatewayMagmaVersionMetric = "gateway_version_count"

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
	GatewayMagmaVersionLabel = "magmaVersion"
)

//GeneralMetricsCalculation ...
type GeneralMetricsCalculation struct {
	calculations.CalculationParams
}

//Calculate this mainly computes the number of LTE, FEG_LTE and FEG networks
// in a particular deployment
func (x *GeneralMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(10).Info("Calculate Generic Metrics")

	results := []*protos.CalculationResult{}
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}
	networkMap := make(map[models1.NetworkType]float64)
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
		if _, ok := networkMap[ret.Type]; !ok {
			networkMap[ret.Type] = 1
		} else {
			networkMap[ret.Type]++
		}
	}
	for networkType, numNetworks := range networkMap {
		labels := prometheus.Labels{}
		networkTypeStr := string(networkType)
		labels[NetworkTypeLabel] = networkTypeStr
		results = append(results, calculations.NewResult(numNetworks, NetworkCountMetric, labels))
	}
	glog.V(10).Info("Generic Metrics Results ", results)
	return results, nil
}

//UserMetricsCalculation ...
type UserMetricsCalculation struct {
	calculations.CalculationParams
}

//Calculate site metrics
func (x *UserMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(10).Info("Calculate User Metrics")

	results := []*protos.CalculationResult{}
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	for _, networkID := range networks {
		subscriberEnts, err := configurator.LoadAllEntitiesOfType(networkID, lte.SubscriberEntityType,
			configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsFromThis: true}, serdes.Entity)
		if err != nil {
			return nil, err
		}
		labels := prometheus.Labels{metrics.NetworkLabelName: networkID}
		results = append(results, calculations.NewResult(float64(len(subscriberEnts)), ConfiguredSubscribersMetric, labels))

		// get all subscribers state across (configured and federated subscribers)
		subscriberStateTypes := []string{lte.SubscriberStateType}
		states, err := state.SearchStates(networkID, subscriberStateTypes, nil, nil, serdes.State)
		if err != nil {
			return nil, err
		}
		users := make(map[string]struct{})
		var exists = struct{}{}
		activeSessionsPerAPN := make(map[string]float64)
		for stateID, st := range states {
			if st.ReportedState != nil {
				reportedSt, ok := st.ReportedState.(*state.ArbitraryJSON)
				if !ok {
					glog.Errorf("reported state for session state having unexpected type %T", st.ReportedState)
					continue
				}
				users[stateID.DeviceID] = exists
				for apnID, apnSessionStatesIntf := range *reportedSt {
					apnSessionStates, ok := apnSessionStatesIntf.([]interface{})
					if !ok {
						glog.Errorf("apnSession states got unexpected type %T", apnSessionStatesIntf)
						continue
					}
					for _, apnSessionStateIntf := range apnSessionStates {
						sessionState, ok := apnSessionStateIntf.(map[string]interface{})
						if !ok {
							glog.Errorf("apnSession state got unexpected type %T", apnSessionStateIntf)
							continue
						}
						if sessionStateIntf, ok := (sessionState)["lifecycle_state"]; ok {
							if sessionState, ok := sessionStateIntf.(string); ok {
								if sessionState == "SESSION_ACTIVE" {
									activeSessionsPerAPN[apnID]++
								}
							}
						}
					}
				}
			}
		}
		results = append(results, calculations.NewResult(float64(len(users)), ActualSubscribersMetric, labels))
		for apnID, numActiveSessionsPerAPN := range activeSessionsPerAPN {
			labels := prometheus.Labels{metrics.NetworkLabelName: networkID, APNLabel: apnID}
			results = append(results, calculations.NewResult(float64(numActiveSessionsPerAPN), ActiveSessionAPNMetric, labels))
		}
	}
	glog.V(10).Info("User Metrics Results ", results)
	return results, nil
}

//SiteMetricsCalculation ...
type SiteMetricsCalculation struct {
	calculations.CalculationParams
}

//Calculate site metrics
func (x *SiteMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(10).Info("Calculate Site Metrics")
	results := []*protos.CalculationResult{}
	networks, err := configurator.ListNetworkIDs()
	if err != nil || networks == nil {
		return results, err
	}

	for _, networkID := range networks {
		// load entities of gateway type
		gatewayEnts, err := configurator.LoadAllEntitiesOfType(
			networkID, lte.CellularGatewayEntityType,
			configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsToThis: true},
			serdes.Entity,
		)
		if err != nil {
			continue
		}

		magmaVersionMap := make(map[string]float64)
		for _, ent := range gatewayEnts {
			status, err := wrappers.GetGatewayStatus(networkID, ent.PhysicalID)
			if err != nil || status == nil || status.PlatformInfo == nil || len(status.PlatformInfo.Packages) == 0 {
				continue
			}

			for _, pkg := range status.PlatformInfo.Packages {
				if pkg.Name == "magma" {
					if pkg.Version != "" {
						if _, ok := magmaVersionMap[pkg.Version]; !ok {
							magmaVersionMap[pkg.Version] = 1
						} else {
							magmaVersionMap[pkg.Version]++
						}
					}
					break
				}
			}
		}

		// load entities of enodeb type
		ents, err := configurator.LoadAllEntitiesOfType(
			networkID, lte.CellularEnodebEntityType,
			configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsToThis: true},
			serdes.Entity,
		)
		if err != nil {
			continue
		}

		enodebConfigTypeMap := make(map[string]float64)
		gatewayEnodebMap := make(map[string]float64)
		for _, ent := range ents {
			enodeb := (&lte_models.Enodeb{}).FromBackendModels(ent)
			if enodeb.EnodebConfig == nil || enodeb.EnodebConfig.ConfigType == "" {
				continue
			}
			if _, ok := enodebConfigTypeMap[enodeb.EnodebConfig.ConfigType]; !ok {
				enodebConfigTypeMap[enodeb.EnodebConfig.ConfigType] = 1
			} else {
				enodebConfigTypeMap[enodeb.EnodebConfig.ConfigType]++
			}

			st, err := state.GetState(networkID, lte.EnodebStateType, ent.Key, serdes.State)
			if err != nil {
				continue
			}
			enodebState := st.ReportedState.(*lte_models.EnodebState)
			ent, err := configurator.LoadEntityForPhysicalID(st.ReporterID, configurator.EntityLoadCriteria{}, serdes.Entity)
			if err == nil {
				enodebState.ReportingGatewayID = ent.Key
			}

			if enodebState.ReportingGatewayID != "" {
				if _, ok := gatewayEnodebMap[enodebState.ReportingGatewayID]; !ok {
					gatewayEnodebMap[enodebState.ReportingGatewayID] = 1
				} else {
					gatewayEnodebMap[enodebState.ReportingGatewayID]++
				}
			}
		}

		// metric identifying managed, unmanaged enodebs
		for configType, numEnode := range enodebConfigTypeMap {
			labels := prometheus.Labels{}
			configTypeStr := string(configType)
			labels[EnodeConfigTypeLabel] = configTypeStr
			labels[metrics.NetworkLabelName] = networkID
			results = append(results, calculations.NewResult(numEnode, EnodeConfigTypeMetric, labels))
		}

		// metric identifying gateways behind enodeBs
		for gatewayID, numEnode := range gatewayEnodebMap {
			labels := prometheus.Labels{}
			labels[metrics.GatewayLabelName] = gatewayID
			labels[metrics.NetworkLabelName] = networkID
			results = append(results, calculations.NewResult(numEnode, EnodebConnectedMetric, labels))
		}

		// metrics giving information on magma version across the gateways
		for version, versionCount := range magmaVersionMap {
			labels := prometheus.Labels{}
			labels[metrics.NetworkLabelName] = networkID
			labels[GatewayMagmaVersionLabel] = version
			results = append(results, calculations.NewResult(versionCount, GatewayMagmaVersionMetric, labels))
		}
	}
	glog.V(10).Info("Site Metrics Results ", results)
	return results, nil
}
