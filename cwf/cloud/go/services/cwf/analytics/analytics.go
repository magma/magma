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
	"fmt"
	"time"

	cwf_calculations "magma/cwf/cloud/go/services/cwf/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics"
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
	daysToQueryStepSize = map[uint]time.Duration{1: 15 * time.Second, 7: time.Minute, 30: 5 * time.Minute}

	daysToCalculate = []uint{1, 7, 30}
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
func GetAnalyticsCalculations(config *calculations.AnalyticsConfig) []calculations.Calculation {
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

	calcs := make([]calculations.Calculation, 0)

	// MAP, WAP, DAP Calculations
	calcs = append(calcs, getXAPCalculations(daysToCalculate, xapGauge, activeUsersMetricName)...)

	// User Throughput Calculations
	calcs = append(calcs, getUserThroughputCalculations(daysToCalculate, userThroughputGauge, userThroughputMetricName)...)

	// AP Throughput Calculations
	calcs = append(calcs, getAPNThroughputCalculations(daysToCalculate, apThroughputGauge, apThroughputMetricName)...)

	// User Consumption Calculations
	calcs = append(calcs, getUserConsumptionCalculations(daysToCalculate, userConsumptionGauge, userConsumptionMetricName)...)
	calcs = append(calcs, get1hourConsumptionCalculation(hourlyUserConsumptionGauge, userConsumptionHourlyMetricName)...)

	// Authentication Calculations
	calcs = append(calcs, getAuthenticationCalculations(daysToCalculate, authenticationsGauge, authenticationsMetricName)...)

	// Raw Metrics
	calcs = append(calcs, analytics.GetRawMetricsCalculations(config)...)

	return calcs
}

func getXAPCalculations(daysList []uint, gauge *prometheus.GaugeVec, metricName string) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		calcs = append(calcs, &cwf_calculations.XAPCalculation{
			BaseCalculation: calculations.BaseCalculation{
				CalculationParams: calculations.CalculationParams{
					Days:                dayParam,
					RegisteredGauge:     gauge,
					Labels:              prometheus.Labels{calculations.DaysLabel: fmt.Sprint(dayParam)},
					Name:                metricName,
					ExpectedGaugeLabels: xapLabels,
				}},
		})
	}
	return calcs
}

func getUserThroughputCalculations(daysList []uint, gauge *prometheus.GaugeVec, metricName string) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)

	for _, dayParam := range daysList {
		for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
			calcs = append(calcs, &cwf_calculations.UserThroughputCalculation{
				BaseCalculation: calculations.BaseCalculation{
					CalculationParams: calculations.CalculationParams{
						Days:                dayParam,
						RegisteredGauge:     gauge,
						Labels:              prometheus.Labels{calculations.DaysLabel: fmt.Sprint(dayParam)},
						Name:                metricName,
						ExpectedGaugeLabels: userThroughputLabels,
					},
				},
				Direction:     dir,
				QueryStepSize: daysToQueryStepSize[dayParam],
			})
		}
	}
	return calcs
}

func getAPNThroughputCalculations(daysList []uint, gauge *prometheus.GaugeVec, metricName string) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
			calcs = append(calcs, &cwf_calculations.APNThroughputCalculation{
				BaseCalculation: calculations.BaseCalculation{
					CalculationParams: calculations.CalculationParams{
						Days:                dayParam,
						RegisteredGauge:     gauge,
						Labels:              prometheus.Labels{calculations.DaysLabel: fmt.Sprint(dayParam)},
						Name:                metricName,
						ExpectedGaugeLabels: apThroughputLabels,
					}},
				Direction:     dir,
				QueryStepSize: daysToQueryStepSize[dayParam],
			})
		}
	}
	return calcs
}

func getUserConsumptionCalculations(daysList []uint, gauge *prometheus.GaugeVec, metricName string) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
			calcs = append(calcs, &cwf_calculations.UserConsumptionCalculation{
				BaseCalculation: calculations.BaseCalculation{
					CalculationParams: calculations.CalculationParams{
						Days:                dayParam,
						RegisteredGauge:     gauge,
						Labels:              prometheus.Labels{calculations.DaysLabel: fmt.Sprint(dayParam)},
						Name:                metricName,
						ExpectedGaugeLabels: userConsumptionLabels,
					}},
				Direction: dir,
			})
		}
	}
	return calcs
}

func get1hourConsumptionCalculation(gauge *prometheus.GaugeVec, metricName string) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dir := range []calculations.ConsumptionDirection{calculations.ConsumptionIn, calculations.ConsumptionOut} {
		calcs = append(calcs, &cwf_calculations.UserConsumptionCalculation{
			BaseCalculation: calculations.BaseCalculation{
				CalculationParams: calculations.CalculationParams{
					Hours:               1,
					RegisteredGauge:     gauge,
					Labels:              prometheus.Labels{"hours": "1"},
					Name:                metricName,
					ExpectedGaugeLabels: hourlyUserConsumptionLabels,
				}},
			Direction: dir,
		})
	}
	return calcs
}

func getAuthenticationCalculations(daysList []uint, gauge *prometheus.GaugeVec, metricName string) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for _, dayParam := range daysList {
		calcs = append(calcs, &cwf_calculations.AuthenticationsCalculation{
			BaseCalculation: calculations.BaseCalculation{
				CalculationParams: calculations.CalculationParams{
					Days:                dayParam,
					RegisteredGauge:     gauge,
					Labels:              prometheus.Labels{calculations.DaysLabel: fmt.Sprint(dayParam)},
					Name:                metricName,
					ExpectedGaugeLabels: authenticationsLabels,
				},
			},
		})
	}
	return calcs
}
