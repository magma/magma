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
	lte_calculations "magma/lte/cloud/go/services/lte/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics"
	"magma/orc8r/cloud/go/services/analytics/calculations"
)

// GetAnalyticsCalculations gets all the calculations provided by lte analytics
func GetAnalyticsCalculations(config *calculations.AnalyticsConfig) []calculations.Calculation {
	if config == nil {
		return nil
	}

	calcs := make([]calculations.Calculation, 0)

	calcs = append(calcs, &lte_calculations.UserMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{AnalyticsConfig: config},
		},
	})
	calcs = append(calcs, &lte_calculations.SiteMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{AnalyticsConfig: config},
		},
	})
	for _, d := range []calculations.ConsumptionDirection{calculations.ConsumptionDown, calculations.ConsumptionUp} {
		calcs = append(calcs, &lte_calculations.UserThroughputCalculation{
			BaseCalculation: calculations.BaseCalculation{
				CalculationParams: calculations.CalculationParams{AnalyticsConfig: config},
			},
			Direction: d,
		})
	}
	calcs = append(calcs, analytics.GetRawMetricsCalculations(config)...)
	return calcs
}
