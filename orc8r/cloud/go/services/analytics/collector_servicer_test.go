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

package analytics_test

import (
	"context"
	"magma/orc8r/cloud/go/services/analytics"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"
	"reflect"
	"sort"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUserThresholdEnforcement(t *testing.T) {
	analyticsConfig := calculations.AnalyticsConfig{
		MinUserThreshold: 100,
		Metrics: map[string]calculations.MetricConfig{
			"test_user_network_threshold_metric": {
				EnforceMinUserThreshold: true,
				Export:                  true,
			},
			"test_user_gateway_threshold_metric": {
				EnforceMinUserThreshold: true,
				Export:                  true,
			},
			"test_user_deployment_threshold_metric": {
				EnforceMinUserThreshold: true,
				Export:                  true,
			},
		},
	}

	userStateMgr := mockUserStateManager{
		totalUsers: 100,
		usersNetworkTable: map[string]*mockNetworkState{
			"mpk_network": {
				usersGatewayTable: map[string]int{
					"mpk_gateway_1": 10,
					"mpk_gateway_2": 5,
				},
				totalUsersPerNetwork: 15,
			},
		},
	}

	calcs := []calculations.Calculation{
		&TestUserCalculations{calculations.BaseCalculation{}},
	}
	collectorServicer := analytics.NewCollectorServicer(&analyticsConfig, nil, calcs, &userStateMgr)
	resp, err := collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(resp.GetResults()), 0)

	// reduce the min user threshold for deployment and check userDeployment is passed
	analyticsConfig.MinUserThreshold = 50
	resp, err = collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(resp.GetResults()), 1)
	result := resp.GetResults()[0]
	assert.Equal(t, result.GetMetricName(), "test_user_deployment_threshold_metric")

	// reduce the min user threshold and check if network metric is also passed
	analyticsConfig.MinUserThreshold = 12
	resp, err = collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(resp.GetResults()), 2)
	resultMetrics := []string{}
	for _, r := range resp.GetResults() {
		resultMetrics = append(resultMetrics, r.GetMetricName())
	}
	expMetrics := []string{"test_user_deployment_threshold_metric", "test_user_network_threshold_metric"}
	sort.Strings(resultMetrics)
	sort.Strings(expMetrics)
	assert.Equal(t, reflect.DeepEqual(resultMetrics, expMetrics), true)

	analyticsConfig.MinUserThreshold = 5
	resp, err = collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(resp.GetResults()), 3)
	resultMetrics = []string{}
	for _, r := range resp.GetResults() {
		resultMetrics = append(resultMetrics, r.GetMetricName())
	}
	expMetrics = []string{"test_user_deployment_threshold_metric", "test_user_network_threshold_metric", "test_user_gateway_threshold_metric"}
	sort.Strings(resultMetrics)
	sort.Strings(expMetrics)
	assert.Equal(t, reflect.DeepEqual(resultMetrics, expMetrics), true)
}

func TestExportEnforcement(t *testing.T) {
	metricName := "test_reliability_metric"
	analyticsConfig := calculations.AnalyticsConfig{
		MinUserThreshold: 100,
		Metrics: map[string]calculations.MetricConfig{
			metricName: {
				Export: false,
			},
		},
	}

	userStateMgr := mockUserStateManager{}
	calcs := []calculations.Calculation{
		&TestNetworkCalculations{calculations.BaseCalculation{}},
	}

	collectorServicer := analytics.NewCollectorServicer(&analyticsConfig, nil, calcs, &userStateMgr)
	resp, err := collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(resp.GetResults()), 0)

	// enable metric export
	analyticsConfig.Metrics[metricName] = calculations.MetricConfig{
		Export: true,
	}
	resp, err = collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(resp.GetResults()), 1)
}

func TestRegisterEnforcement(t *testing.T) {
	metricName := "test_reliability_metric"
	analyticsConfig := calculations.AnalyticsConfig{
		MinUserThreshold: 100,
		Metrics: map[string]calculations.MetricConfig{
			metricName: {
				Register: false,
				Export:   true,
			},
		},
	}

	userStateMgr := mockUserStateManager{}
	reliabilityGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricName}, []string{metrics.NetworkLabelName})

	calcs := []calculations.Calculation{
		&TestNetworkCalculations{calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{
				RegisteredGauge:     reliabilityGauge,
				Name:                metricName,
				ExpectedGaugeLabels: []string{metrics.NetworkLabelName},
			},
		}},
	}

	collectorServicer := analytics.NewCollectorServicer(&analyticsConfig, nil, calcs, &userStateMgr)
	resp, err := collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	results := resp.GetResults()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 0, testutil.CollectAndCount(reliabilityGauge))

	analyticsConfig.Metrics[metricName] = calculations.MetricConfig{
		Register: true,
		Export:   true,
	}
	resp, err = collectorServicer.Collect(context.Background(), &protos.CollectRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(resp.GetResults()), 1)

	// Check if reliability gauge value is set to the calculation returned by TestNetworkCalculations
	assert.Equal(t, 1, testutil.CollectAndCount(reliabilityGauge))
	v := testutil.ToFloat64(reliabilityGauge)
	assert.Equal(t, v, 5.0)
}

type TestUserCalculations struct {
	calculations.BaseCalculation
}

func (testUserCalc *TestUserCalculations) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	results := []*protos.CalculationResult{}

	results = append(results, calculations.NewResult(5, "test_user_network_threshold_metric", map[string]string{
		metrics.NetworkLabelName: "mpk_network",
	}))
	results = append(results, calculations.NewResult(5, "test_user_gateway_threshold_metric", map[string]string{
		metrics.NetworkLabelName: "mpk_network",
		metrics.GatewayLabelName: "mpk_gateway_1",
	}))
	results = append(results, calculations.NewResult(5, "test_user_deployment_threshold_metric", map[string]string{}))
	return results, nil
}

type TestNetworkCalculations struct {
	calculations.BaseCalculation
}

func (testNetworkCalc *TestNetworkCalculations) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	results := []*protos.CalculationResult{}
	results = append(results, calculations.NewResult(5, "test_reliability_metric", map[string]string{metrics.NetworkLabelName: "mpk_network"}))
	return results, nil
}

type mockNetworkState struct {
	totalUsersPerNetwork int
	usersGatewayTable    map[string]int
}

type mockUserStateManager struct {
	totalUsers        int
	usersNetworkTable map[string]*mockNetworkState
}

func (u *mockUserStateManager) Update() {}

func (u *mockUserStateManager) GetTotalUsers() int {
	return u.totalUsers
}

func (u *mockUserStateManager) GetTotalUsersInNetwork(networkID string) int {
	v, ok := u.usersNetworkTable[networkID]
	if !ok {
		return 0
	}
	return v.totalUsersPerNetwork
}

func (u *mockUserStateManager) GetTotalUsersInGateway(networkID string, gatewayID string) int {
	v, ok := u.usersNetworkTable[networkID]
	if !ok {
		return 0
	}
	g, ok := v.usersGatewayTable[gatewayID]
	if !ok {
		return 0
	}
	return g
}
