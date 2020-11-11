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
	"strconv"
	"time"

	cwf_service "magma/cwf/cloud/go/services/cwf"
	cwf_calculations "magma/cwf/cloud/go/services/cwf/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/lib/go/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	activeUsersMetricName           = "active_users_over_time"
	userThroughputMetricName        = "user_throughput"
	userConsumptionMetricName       = "user_consumption"
	userConsumptionHourlyMetricName = "user_consumption_hourly"
	apThroughputMetricName          = "throughput_per_ap"
	authenticationsMetricName       = "authentications_over_time"
)

var (
	// Map from number of days to query to size the step should be to get best granularity
	// without causes prometheus to reject the query for having too many datapoints
	daysToQueryStepSize = map[int]time.Duration{1: 15 * time.Second, 7: time.Minute, 30: 5 * time.Minute}

	daysToCalculate = []int{1, 7, 30}
)

var (
	xapLabels                   = []string{calculations.DaysLabel, metrics.NetworkLabelName}
	userThroughputLabels        = []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.DirectionLabel}
	userConsumptionLabels       = []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.DirectionLabel}
	hourlyUserConsumptionLabels = []string{"hours", metrics.NetworkLabelName, calculations.DirectionLabel}
	apThroughputLabels          = []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.DirectionLabel, calculations.APNLabel}
	authenticationsLabels       = []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.AuthCodeLabel}
)

// GetAnalyticsCalculations ..
func GetAnalyticsCalculations(config *cwf_service.Config) []calculations.Calculation {
	if config == nil {
		return nil
	}

	xapGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: activeUsersMetricName}, xapLabels)
	userThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userThroughputMetricName}, userThroughputLabels)
	userConsumptionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userConsumptionMetricName}, userConsumptionLabels)
	hourlyUserConsumptionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userConsumptionHourlyMetricName}, hourlyUserConsumptionLabels)
	apThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: apThroughputMetricName}, apThroughputLabels)
	authenticationsGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: authenticationsMetricName}, authenticationsLabels)

	prometheus.MustRegister(xapGauge, userThroughputGauge, userConsumptionGauge,
		hourlyUserConsumptionGauge, apThroughputGauge, authenticationsGauge)

	allCalculations := make([]calculations.Calculation, 0)

	metricConfig := config.Analytics.Metrics
	// MAP, WAP, DAP Calculations
	allCalculations = append(allCalculations, getXAPCalculations(daysToCalculate, xapGauge, activeUsersMetricName, metricConfig)...)

	// User Throughput Calculations
	allCalculations = append(allCalculations, getUserThroughputCalculations(daysToCalculate, userThroughputGauge, userThroughputMetricName, metricConfig)...)

	// AP Throughput Calculations
	allCalculations = append(allCalculations, getAPNThroughputCalculations(daysToCalculate, apThroughputGauge, apThroughputMetricName, metricConfig)...)

	// User Consumption Calculations
	allCalculations = append(allCalculations, getUserConsumptionCalculations(daysToCalculate, userConsumptionGauge, userConsumptionMetricName, metricConfig)...)
	allCalculations = append(allCalculations, get1hourConsumptionCalculation(hourlyUserConsumptionGauge, userConsumptionHourlyMetricName, metricConfig)...)

	// Authentication Calculations
	allCalculations = append(allCalculations, getAuthenticationCalculations(daysToCalculate, authenticationsGauge, authenticationsMetricName, metricConfig)...)

	// Raw Metrics
	allCalculations = append(allCalculations, calculations.GetRawMetricsCalculations(metricConfig)...)

	return allCalculations
}

func getXAPCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string, metricConfig map[string]calculations.MetricConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		calcs = append(calcs, &cwf_calculations.XAPCalculation{
			CalculationParams: calculations.CalculationParams{
				Days:                dayParam,
				RegisteredGauge:     gauge,
				Labels:              prometheus.Labels{calculations.DaysLabel: strconv.Itoa(dayParam)},
				Name:                metricName,
				ExpectedGaugeLabels: xapLabels,
				MetricConfig:        metricConfig,
			},
		})
	}
	return calcs
}

func getUserThroughputCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string, metricConfig map[string]calculations.MetricConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
			calcs = append(calcs, &cwf_calculations.UserThroughputCalculation{
				CalculationParams: calculations.CalculationParams{
					Days:                dayParam,
					RegisteredGauge:     gauge,
					Labels:              prometheus.Labels{calculations.DaysLabel: strconv.Itoa(dayParam)},
					Name:                metricName,
					ExpectedGaugeLabels: userThroughputLabels,
					MetricConfig:        metricConfig,
				},
				Direction:     dir,
				QueryStepSize: daysToQueryStepSize[dayParam],
			})
		}
	}
	return calcs
}

func getAPNThroughputCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string, metricConfig map[string]calculations.MetricConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
			calcs = append(calcs, &cwf_calculations.APNThroughputCalculation{
				CalculationParams: calculations.CalculationParams{
					Days:                dayParam,
					RegisteredGauge:     gauge,
					Labels:              prometheus.Labels{calculations.DaysLabel: strconv.Itoa(dayParam)},
					Name:                metricName,
					ExpectedGaugeLabels: apThroughputLabels,
					MetricConfig:        metricConfig,
				},
				Direction:     dir,
				QueryStepSize: daysToQueryStepSize[dayParam],
			})
		}
	}
	return calcs
}

func getUserConsumptionCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string, metricConfig map[string]calculations.MetricConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
			calcs = append(calcs, &cwf_calculations.UserConsumptionCalculation{
				CalculationParams: calculations.CalculationParams{
					Days:                dayParam,
					RegisteredGauge:     gauge,
					Labels:              prometheus.Labels{calculations.DaysLabel: strconv.Itoa(dayParam)},
					Name:                metricName,
					ExpectedGaugeLabels: userConsumptionLabels,
					MetricConfig:        metricConfig,
				},
				Direction: dir,
			})
		}
	}
	return calcs
}

func get1hourConsumptionCalculation(gauge *prometheus.GaugeVec, metricName string, metricConfig map[string]calculations.MetricConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
		calcs = append(calcs, &cwf_calculations.UserConsumptionCalculation{
			CalculationParams: calculations.CalculationParams{
				Hours:               1,
				RegisteredGauge:     gauge,
				Labels:              prometheus.Labels{"hours": "1"},
				Name:                metricName,
				ExpectedGaugeLabels: hourlyUserConsumptionLabels,
				MetricConfig:        metricConfig,
			},
			Direction: dir,
		})
	}
	return calcs
}

func getAuthenticationCalculations(daysList []int, gauge *prometheus.GaugeVec, metricName string, metricConfig map[string]calculations.MetricConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		calcs = append(calcs, &cwf_calculations.AuthenticationsCalculation{
			CalculationParams: calculations.CalculationParams{
				Days:                dayParam,
				RegisteredGauge:     gauge,
				Labels:              prometheus.Labels{calculations.DaysLabel: strconv.Itoa(dayParam)},
				Name:                metricName,
				ExpectedGaugeLabels: authenticationsLabels,
				MetricConfig:        metricConfig,
			},
		})
	}
	return calcs
}
