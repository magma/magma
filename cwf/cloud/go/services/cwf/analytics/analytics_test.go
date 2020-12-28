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
	"testing"

	cwf_calculations "magma/cwf/cloud/go/services/cwf/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/calculations"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestGetXAPCalculations(t *testing.T) {
	xapGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: activeUsersMetricName}, xapLabels)
	calcs := getXAPCalculations([]uint{1, 7, 30}, xapGauge, "metricName")
	assert.Len(t, calcs, 3)
	for _, calc := range calcs {
		c := calc.(*cwf_calculations.XAPCalculation)
		assert.Equal(t, fmt.Sprint(c.Days), c.Labels[calculations.DaysLabel])
	}
}

func TestGetUserThroughputCalculations(t *testing.T) {
	userThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userThroughputMetricName}, userThroughputLabels)
	calcs := getUserThroughputCalculations([]uint{1, 7, 30}, userThroughputGauge, "metricName")
	assert.Len(t, calcs, 6)
	for _, calc := range calcs {
		c := calc.(*cwf_calculations.UserThroughputCalculation)
		assert.Equal(t, fmt.Sprint(c.Days), c.Labels[calculations.DaysLabel])
	}
}

func TestGetUserConsumptionCalculations(t *testing.T) {
	userConsumptionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userConsumptionMetricName}, userConsumptionLabels)
	calcs := getUserConsumptionCalculations([]uint{1, 7, 30}, userConsumptionGauge, "metricName")
	assert.Len(t, calcs, 6)
	for _, calc := range calcs {
		c := calc.(*cwf_calculations.UserConsumptionCalculation)
		assert.Equal(t, fmt.Sprint(c.Days), c.Labels[calculations.DaysLabel])
	}
}

func TestGet1hourUserConsumptionCalculations(t *testing.T) {
	userConsumptionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: userConsumptionMetricName}, hourlyUserConsumptionLabels)
	calcs := get1hourConsumptionCalculation(userConsumptionGauge, "metricName")
	assert.Len(t, calcs, 2)
	for _, calc := range calcs {
		c := calc.(*cwf_calculations.UserConsumptionCalculation)
		assert.Equal(t, fmt.Sprint(c.Hours), c.Labels["hours"])
	}
}

func TestGetAPThroughputCalculations(t *testing.T) {
	apThroughputGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: apThroughputMetricName}, apThroughputLabels)
	calcs := getAPNThroughputCalculations([]uint{1, 7, 30}, apThroughputGauge, "metricName")
	assert.Len(t, calcs, 6)
	for _, calc := range calcs {
		c := calc.(*cwf_calculations.APNThroughputCalculation)
		assert.Equal(t, fmt.Sprint(c.Days), c.Labels[calculations.DaysLabel])
	}
}

func TestGetAuthenticationCalculations(t *testing.T) {
	authenticationsGague := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: authenticationsMetricName}, authenticationsLabels)
	calcs := getAuthenticationCalculations(daysToCalculate, authenticationsGague, "metricName")
	assert.Len(t, calcs, 3)
	for _, calc := range calcs {
		c := calc.(*cwf_calculations.AuthenticationsCalculation)
		assert.Equal(t, fmt.Sprint(c.Days), c.Labels[calculations.DaysLabel])
	}
}
