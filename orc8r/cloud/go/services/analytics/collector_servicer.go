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

package analytics

import (
	"context"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
)

// CollectorService implements the operations of collecting the metrics from CWF service
type CollectorService struct {
	calculations     []calculations.Calculation
	promAPIClient    query_api.PrometheusAPI
	UserStateManager calculations.UserStateManager
	analyticsConfig  *calculations.AnalyticsConfig
}

// NewCollectorService constructs new collector service
func NewCollectorService(
	analyticsConfig *calculations.AnalyticsConfig,
	promAPIClient query_api.PrometheusAPI,
	calculations []calculations.Calculation,
	UserStateManager calculations.UserStateManager) *CollectorService {
	return &CollectorService{
		promAPIClient:    promAPIClient,
		calculations:     calculations,
		analyticsConfig:  analyticsConfig,
		UserStateManager: UserStateManager,
	}
}

// filterResults filters the results to ensure that the user aggregate metrics get filtered out in
// case the number of users in that context falls below minimum threshold
func (svc *CollectorService) filterResults(results []*protos.CalculationResult) []*protos.CalculationResult {
	UserStateManager := svc.UserStateManager
	// Nothing to filter
	if UserStateManager == nil {
		return results
	}

	filteredResults := []*protos.CalculationResult{}
	analyticsConfig := svc.analyticsConfig
	for _, result := range results {
		// Enforce the minimum user threshold constraint at network level, gateway level and deployment level
		metricConfig, ok := analyticsConfig.Metrics[result.GetMetricName()]
		if !ok {
			glog.Errorf("Metric Configuration not found for %s", result.GetMetricName())
			continue
		}
		if !metricConfig.Export {
			glog.V(1).Infof("Metric %s export configuration disabled", result.GetMetricName())
			continue
		}
		if !metricConfig.EnforceMinUserThreshold {
			filteredResults = append(filteredResults, result)
			continue
		}
		networkID, networkLabelPresent := result.Labels[metrics.NetworkLabelName]
		gatewayID, gatewayLabelPresent := result.Labels[metrics.GatewayLabelName]

		if networkLabelPresent {
			if gatewayLabelPresent {
				numUsersInGateway := UserStateManager.GetTotalUsersInGateway(networkID, gatewayID)
				if numUsersInGateway > analyticsConfig.MinUserThreshold {
					filteredResults = append(filteredResults, result)
				} else {
					glog.V(1).Infof("Metric %s(gateway label) dropped,active users %d below user threshold %d",
						result.GetMetricName(), numUsersInGateway, analyticsConfig.MinUserThreshold)
				}

			} else {
				numUsersInNetwork := UserStateManager.GetTotalUsersInNetwork(networkID)
				if numUsersInNetwork > analyticsConfig.MinUserThreshold {
					filteredResults = append(filteredResults, result)
				} else {
					glog.V(1).Infof("Metric %s(network label) dropped,active users %d below user threshold %d",
						result.GetMetricName(), numUsersInNetwork, analyticsConfig.MinUserThreshold)
				}
			}
		} else {
			numUsersInDeployment := UserStateManager.GetTotalUsers()
			if numUsersInDeployment > analyticsConfig.MinUserThreshold {
				filteredResults = append(filteredResults, result)
			} else {
				glog.V(1).Infof("Metric %s(deployment) dropped, active users %d below user threshold %d", result.GetMetricName(),
					numUsersInDeployment, analyticsConfig.MinUserThreshold)
			}
		}

	}
	return filteredResults
}

// RegisterResults registers the computed metric with Prometheus based on the metric configuration
func (svc *CollectorService) registerResults(calcParams *calculations.CalculationParams, results []*protos.CalculationResult) {
	filteredResults := []*protos.CalculationResult{}
	analyticsConfig := svc.analyticsConfig
	for _, result := range results {
		if metricConfig, ok := analyticsConfig.Metrics[result.GetMetricName()]; ok {
			if metricConfig.Register {
				filteredResults = append(filteredResults, result)
			}
		} else {
			glog.Errorf("Metric configuration not found for '%s'", result.GetMetricName())
		}
	}
	if len(filteredResults) > 0 {
		glog.V(1).Info("Registering ", filteredResults)
		calculations.RegisterResults(calcParams, filteredResults)
	}
}

// Collect does the operation of running through calculations and returning results
func (svc *CollectorService) Collect(context.Context, *protos.CollectRequest) (*protos.CollectResponse, error) {
	if svc.UserStateManager != nil {
		svc.UserStateManager.Update()
	}

	response := &protos.CollectResponse{}
	for _, calc := range svc.calculations {
		results, err := calc.Calculate(svc.promAPIClient)
		if err != nil {
			glog.Errorf("Error calculating metric: %s", err)
			continue
		}

		svc.registerResults(calc.GetCalculationParams(), results)
		filteredResults := svc.filterResults(results)
		response.Results = append(response.Results, filteredResults...)
	}
	return response, nil
}
