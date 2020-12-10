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
	"fmt"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"sort"
	"strings"

	"github.com/golang/glog"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

const (
	//APNLabel defines label string literal for APN
	APNLabel = "apn"

	//AuthCodeLabel defines label string literal for authentication code
	AuthCodeLabel = "code"

	//DaysLabel defines label string literal for days
	DaysLabel = "days"

	//DirectionLabel defines label string literal for direction
	DirectionLabel = "direction"
)

//Calculation interface
type Calculation interface {
	Calculate(query_api.PrometheusAPI) ([]*protos.CalculationResult, error)
}

//MetricConfig is the expected configuration for a specific metric in the config file
type MetricConfig struct {
	Register             bool     `yaml:"register"`
	Export               bool     `yaml:"export"`
	Expr                 string   `yaml:"expr"`
	Labels               []string `yaml:"labels"`
	EnforceAggrThreshold bool     `yaml:"enforceMinThreshold"`
}

//AnalyticsConfig represents the configuration provided to analytics components
type AnalyticsConfig struct {
	//AggrMinThreshold minimum threshold for exporting aggregate metrics
	AggrMinThreshold int                     `yaml:"aggrMetricThreshold"`
	Metrics          map[string]MetricConfig `yaml:"metrics"`
}

//CalculationParams calculations paramters
type CalculationParams struct {
	Hours               int
	Days                int
	RegisteredGauge     *prometheus.GaugeVec
	Labels              prometheus.Labels
	Name                string
	ExpectedGaugeLabels []string
	MetricConfig        map[string]MetricConfig
}

//ConsumptionDirection defines the direction type
type ConsumptionDirection string

const (
	//ConsumptionIn string literal for identfying incoming data
	ConsumptionIn ConsumptionDirection = "in"
	//ConsumptionOut string literal for identifying outgoign data
	ConsumptionOut ConsumptionDirection = "out"
)

//AverageDatapoints method to compute average of datapoints
func AverageDatapoints(samples []model.SamplePair) float64 {
	sum := float64(0)
	for _, val := range samples {
		sum += float64(val.Value)
	}
	return sum / float64(len(samples))
}

//MakeVectorResults build results from vector
func MakeVectorResults(vec model.Vector, baseLabels prometheus.Labels, metricName string) []*protos.CalculationResult {
	var results []*protos.CalculationResult
	for _, v := range vec {
		// Get labels from query result
		queryLabels := make(map[string]string, 0)
		for label, value := range v.Metric {
			queryLabels[string(label)] = string(value)
		}
		combinedLabels := CombineLabels(baseLabels, queryLabels)
		results = append(results, &protos.CalculationResult{
			MetricName: metricName,
			Labels:     combinedLabels,
			Value:      float64(v.Value),
		})
	}
	return results
}

//NewResult builds a new protos.CalculationResult
func NewResult(value float64, metricName string, labels prometheus.Labels) *protos.CalculationResult {
	return &protos.CalculationResult{
		Value:      value,
		MetricName: metricName,
		Labels:     labels,
	}
}

//CombineLabels combine all the label
func CombineLabels(l1, l2 map[string]string) map[string]string {
	retLabels := make(map[string]string)
	for l, v := range l1 {
		retLabels[l] = v
	}
	for l, v := range l2 {
		retLabels[l] = v
	}
	return retLabels
}

//RegisterResults exports the metrics to prometheus
func RegisterResults(calc CalculationParams, results []*protos.CalculationResult) {
	for _, res := range results {
		if !CheckLabelsMatch(calc.ExpectedGaugeLabels, res.Labels) {
			glog.Errorf("Unmatched labels in Calculation. Expected: %s, Received: %s", calc.ExpectedGaugeLabels, printLabels(res.Labels))
			continue
		}
		if calc.MetricConfig != nil {
			metricConfig, ok := calc.MetricConfig[res.MetricName]
			if ok && metricConfig.Register == false {
				glog.V(2).Infof("%s result registration skipped due to configuration", calc.Name)
				return
			}
		}
		calc.RegisteredGauge.With(res.Labels).Set(res.Value)
		glog.V(10).Infof("Set metric %s{%s} value: %f\n", res.MetricName, printLabels(res.Labels), res.Value)
	}
}

//CheckLabelsMatch check if labels match
func CheckLabelsMatch(expectedLabels []string, labels prometheus.Labels) bool {
	givenLabels := []string{}
	for l := range labels {
		givenLabels = append(givenLabels, l)
	}
	sort.Strings(givenLabels)
	sort.Strings(expectedLabels)
	return cmp.Equal(givenLabels, expectedLabels)
}

func printLabels(labels prometheus.Labels) string {
	str := strings.Builder{}
	str.WriteString("{")
	for key, val := range labels {
		str.WriteString(fmt.Sprintf("%s=\"%s\"", key, val))
	}
	str.WriteString("}")
	return str.String()
}

//GetRawMetricsCalculations ...
func GetRawMetricsCalculations(metrics map[string]MetricConfig) []Calculation {
	allCalculations := make([]Calculation, 0)
	for metricName, metricConfig := range metrics {
		if metricConfig.Expr != "" {
			glog.V(10).Infof("Adding RawMetrics Calculation for %s", metricName)
			allCalculations = append(allCalculations, &RawMetricsCalculation{
				CalculationParams: CalculationParams{
					Name:         metricName,
					MetricConfig: metrics,
				},
				MetricExpr: metricConfig.Expr,
			})
		}
	}
	return allCalculations
}
