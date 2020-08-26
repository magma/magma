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
	"magma/cwf/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"
)

type UserConsumptionCalculation struct {
	CalculationParams
	Direction ConsumptionDirection
	Hours     int
}

func (x *UserConsumptionCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]Result, error) {

	var consumptionQuery string
	// Measure consumption over x.Hours if exists
	if x.Hours > 0{
		consumptionQuery = fmt.Sprintf(`sum(increase(octets_%s[%dh])) by (%s)`, x.Direction, x.Hours, metrics.NetworkLabelName)
	} else {
		consumptionQuery = fmt.Sprintf(`sum(increase(octets_%s[%dd])) by (%s)`, x.Direction, x.Days, metrics.NetworkLabelName)
	}

	vec, err := query_api.QueryPrometheusVector(prometheusClient, consumptionQuery)
	if err != nil {
		return nil, fmt.Errorf("user Consumption query error: %s", err)
	}

	baseLabels := combineLabels(x.Labels, map[string]string{DirectionLabel: string(x.Direction)})
	results := makeVectorResults(vec, baseLabels, x.Name)
	for _, res := range results {
		x.RegisteredGauge.With(res.labels).Set(res.value)
	}
	return results, nil
}
