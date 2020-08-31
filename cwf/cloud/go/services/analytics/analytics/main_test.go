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

package main

import (
	"magma/cwf/cloud/go/services/analytics/calculations"
	"strconv"
	"testing"

	"magma/orc8r/lib/go/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestGetXAPCalculations(t *testing.T) {
	xapGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: activeUsersMetricName}, []string{calculations.DaysLabel, metrics.NetworkLabelName})
	calcs := getXAPCalculations([]int{1, 7, 30}, xapGauge, "metricName")
	for _, calc := range calcs {
		xapCalc := calc.(*calculations.XAPCalculation)
		assert.Equal(t, strconv.Itoa(xapCalc.Days), xapCalc.Labels[calculations.DaysLabel])
	}
}

func TestGetUserThroughputCalculations(t *testing.T) {
	userThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userThroughputMetricName}, []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.DirectionLabel, "hours"})
	calcs := getUserThroughputCalculations([]int{1, 7, 30}, userThroughputGauge, "metricName")
	for _, calc := range calcs {
		c := calc.(*calculations.UserThroughputCalculation)
		assert.Equal(t, strconv.Itoa(c.Days), c.Labels[calculations.DaysLabel])
	}
}

func TestGetUserConsumptionCalculations(t *testing.T) {
	userConsumptionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userConsumptionMetricName}, []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.DirectionLabel, "hours"})
	calcs := getUserConsumptionCalculations([]int{1, 7, 30}, userConsumptionGauge, "metricName")
	for _, calc := range calcs {
		c := calc.(*calculations.UserConsumptionCalculation)
		assert.Equal(t, strconv.Itoa(c.Days), c.Labels[calculations.DaysLabel])
	}
}

func TestGet1hourUserConsumptionCalculations(t *testing.T) {
	userConsumptionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userConsumptionMetricName}, []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.DirectionLabel})
	calcs := get1hourConsumptionCalculation(userConsumptionGauge, "metricName")
	for _, calc := range calcs {
		c := calc.(*calculations.UserConsumptionCalculation)
		assert.Equal(t, strconv.Itoa(c.Hours), c.Labels["hours"])
	}
}

func TestGetAPThroughputCalculations(t *testing.T) {
	apThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: apThroughputMetricName}, []string{calculations.DaysLabel, metrics.NetworkLabelName, calculations.DirectionLabel, calculations.APNLabel})
	calcs := getAPThroughputCalculations([]int{1, 7, 30}, apThroughputGauge, "metricName")
	for _, calc := range calcs {
		c := calc.(*calculations.APThroughputCalculation)
		assert.Equal(t, strconv.Itoa(c.Days), c.Labels[calculations.DaysLabel])
	}
}
