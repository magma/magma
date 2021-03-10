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

// collectorServicer implements the operations of collecting the metrics from CWF service
type collectorServicer struct {
	calculations     []calculations.Calculation
	promAPIClient    query_api.PrometheusAPI
	userStateManager calculations.UserStateManager
	analyticsConfig  *calculations.AnalyticsConfig
}

// NewCollectorServicer constructs new collector service
func NewCollectorServicer(
	analyticsConfig *calculations.AnalyticsConfig,
	promAPIClient query_api.PrometheusAPI,
	calculations []calculations.Calculation,
	userStateManager calculations.UserStateManager,
) protos.AnalyticsCollectorServer {
	return &collectorServicer{
		promAPIClient:    promAPIClient,
		calculations:     calculations,
		analyticsConfig:  analyticsConfig,
		userStateManager: userStateManager,
	}
}

// Collect does the operation of running through calculations and returning results
func (c *collectorServicer) Collect(context.Context, *protos.CollectRequest) (*protos.CollectResponse, error) {
	if c.userStateManager != nil {
		c.userStateManager.Update()
	}

	response := &protos.CollectResponse{}
	for _, calc := range c.calculations {
		results, err := calc.Calculate(c.promAPIClient)
		if err != nil {
			glog.Errorf("Error %v calculating metric for %v", err, calc.GetCalculationParams())
			continue
		}

		c.registerResults(calc.GetCalculationParams(), results)
		filteredResults := c.filterResults(results)
		response.Results = append(response.Results, filteredResults...)
	}
	return response, nil
}

// FilterResults filters the results to ensure that the user aggregate metrics
// get filtered out in case the number of users in that context falls below
// minimum threshold
func (c *collectorServicer) filterResults(results []*protos.CalculationResult) []*protos.CalculationResult {
	userStateManager := c.userStateManager
	// Nothing to filter
	if userStateManager == nil {
		return results
	}

	filteredResults := []*protos.CalculationResult{}
	for _, result := range results {
		if c.filterResult(result) {
			filteredResults = append(filteredResults, result)
		}
	}
	return filteredResults
}

func (c *collectorServicer) filterResult(result *protos.CalculationResult) bool {
	analyticsConfig := c.analyticsConfig
	metricConfig, ok := analyticsConfig.Metrics[result.GetMetricName()]
	if !ok {
		glog.Errorf("Metric Configuration not found for %s", result.GetMetricName())
		return false
	}
	if !metricConfig.Export {
		glog.V(1).Infof("Metric %s export configuration disabled", result.GetMetricName())
		return false
	}
	if !metricConfig.EnforceMinUserThreshold {
		glog.V(2).Infof("Metric(%s) not a user metric, skipping filtering", result.GetMetricName())
		return true
	}

	networkID, networkLabelPresent := result.Labels[metrics.NetworkLabelName]
	gatewayID, gatewayLabelPresent := result.Labels[metrics.GatewayLabelName]

	numUsers := 0
	minUserVerificationScope := "Deployment Wide"
	if networkLabelPresent && gatewayLabelPresent {
		numUsers = c.userStateManager.GetTotalUsersInGateway(networkID, gatewayID)
		minUserVerificationScope = "Site Wide"
	} else if networkLabelPresent {
		numUsers = c.userStateManager.GetTotalUsersInNetwork(networkID)
		minUserVerificationScope = "Network Wide"
	} else {
		numUsers = c.userStateManager.GetTotalUsers()
	}

	if numUsers > analyticsConfig.MinUserThreshold {
		glog.V(2).Infof("Metric(%s) network label(%s) gateway label(%s) can "+
			"be exported numUsers(%d) in the %s greater than minUserThreshold(%d)",
			result.GetMetricName(),
			networkID,
			gatewayID,
			numUsers,
			minUserVerificationScope,
			analyticsConfig.MinUserThreshold)
		return true
	}
	glog.V(1).Infof(
		"Metric(%s) network label(%s) gateway label(%s) dropped, active "+
			"users(%d) in the %s below user threshold (%d)\n",
		result.GetMetricName(),
		networkID,
		gatewayID,
		numUsers,
		minUserVerificationScope,
		analyticsConfig.MinUserThreshold)
	return false
}

// RegisterResults registers the computed metric with Prometheus based on the metric configuration
func (c *collectorServicer) registerResults(calcParams calculations.CalculationParams, results []*protos.CalculationResult) {
	filteredResults := []*protos.CalculationResult{}
	analyticsConfig := c.analyticsConfig
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
