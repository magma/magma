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

	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
)

// UserConsumptionCalculation input params, direction can be user consumption volume
// during upload or download
type UserConsumptionCalculation struct {
	calculations.BaseCalculation
	Direction calculations.ConsumptionDirection
}

// Calculate computes the total volume of data consumed by subscribers over a required timeperiod
func (x *UserConsumptionCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.Infof("Calculating User Consumption. Days: %d, Hours: %d, Direction: %s", x.Days, x.Hours, x.Direction)

	var consumptionQuery string
	// Measure consumption over x.Hours if exists
	if x.Hours > 0 {
		consumptionQuery = fmt.Sprintf(`sum(increase(octets_%s[%dh])) by (%s)`, x.Direction, x.Hours, metrics.NetworkLabelName)
	} else {
		consumptionQuery = fmt.Sprintf(`sum(increase(octets_%s[%dd])) by (%s)`, x.Direction, x.Days, metrics.NetworkLabelName)
	}

	vec, err := query_api.QueryPrometheusVector(prometheusClient, consumptionQuery)
	if err != nil {
		return nil, fmt.Errorf("user Consumption query error: %s", err)
	}

	baseLabels := calculations.CombineLabels(x.Labels, map[string]string{calculations.DirectionLabel: string(x.Direction)})
	results := calculations.MakeVectorResults(vec, baseLabels, x.Name)
	return results, nil
}
