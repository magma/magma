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
	"github.com/golang/glog"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"magma/cwf/cloud/go/services/analytics/query_api"
	"sort"
	"strings"
)

const (
	APNLabel       = "apn"
	AuthCodeLabel  = "code"
	DaysLabel      = "days"
	DirectionLabel = "direction"
)

type Calculation interface {
	Calculate(query_api.PrometheusAPI) ([]Result, error)
}

type CalculationParams struct {
	Days                int
	RegisteredGauge     *prometheus.GaugeVec
	Labels              prometheus.Labels
	Name                string
	ExpectedGaugeLabels []string
}

type Result struct {
	value      float64
	metricName string
	labels     prometheus.Labels
}

func NewResult(value float64, metricName string, labels prometheus.Labels) Result {
	return Result{
		value:      value,
		metricName: metricName,
		labels:     labels,
	}
}

func (r *Result) Value() float64 {
	return r.value
}

func (r *Result) MetricName() string {
	return r.metricName
}

func (r *Result) Labels() prometheus.Labels {
	return r.labels
}

type ConsumptionDirection string

const (
	ConsumptionIn  ConsumptionDirection = "in"
	ConsumptionOut ConsumptionDirection = "out"
)

func averageDatapoints(samples []model.SamplePair) float64 {
	sum := float64(0)
	for _, val := range samples {
		sum += float64(val.Value)
	}
	return sum / float64(len(samples))
}

func makeVectorResults(vec model.Vector, baseLabels prometheus.Labels, metricName string) []Result {
	var results []Result
	for _, v := range vec {
		// Get labels from query result
		queryLabels := make(map[string]string, 0)
		for label, value := range v.Metric {
			queryLabels[string(label)] = string(value)
		}
		combinedLabels := combineLabels(baseLabels, queryLabels)
		results = append(results, Result{
			metricName: metricName,
			labels:     combinedLabels,
			value:      float64(v.Value),
		})
	}
	return results
}

func combineLabels(l1, l2 map[string]string) map[string]string {
	retLabels := make(map[string]string)
	for l, v := range l1 {
		retLabels[l] = v
	}
	for l, v := range l2 {
		retLabels[l] = v
	}
	return retLabels
}

func registerResults(calc CalculationParams, results []Result) {
	for _, res := range results {
		if !checkLabelsMatch(calc.ExpectedGaugeLabels, res.labels) {
			glog.Errorf("Unmatched labels in APThroughput Calculation. Expected: %s, Received: %s", calc.ExpectedGaugeLabels, printLabels(res.labels))
			continue
		}
		calc.RegisteredGauge.With(res.labels).Set(res.value)
		glog.Infof("Set metric %s{%s} value: %f\n", res.metricName, printLabels(res.labels), res.value)
	}
}

func checkLabelsMatch(expectedLabels []string, labels prometheus.Labels) bool {
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
