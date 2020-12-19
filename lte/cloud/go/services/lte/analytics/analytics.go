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
	"magma/orc8r/cloud/go/services/analytics/calculations"
)

// GetAnalyticsCalculations gets all the calculations provided by lte analytics
func GetAnalyticsCalculations(analyticsConfig *calculations.AnalyticsConfig) []calculations.Calculation {
	if analyticsConfig == nil {
		return nil
	}

	calcs := make([]calculations.Calculation, 0)
	calcs = append(calcs, &lte_calculations.GeneralMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: &calculations.CalculationParams{AnalyticsConfig: analyticsConfig},
		},
	})
	calcs = append(calcs, &lte_calculations.UserMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: &calculations.CalculationParams{AnalyticsConfig: analyticsConfig},
		},
	})
	calcs = append(calcs, &lte_calculations.SiteMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: &calculations.CalculationParams{AnalyticsConfig: analyticsConfig},
		},
	})
	for _, d := range []calculations.ConsumptionDirection{calculations.ConsumptionDown, calculations.ConsumptionDown} {
		calcs = append(calcs, &lte_calculations.UserThroughputCalculation{
			BaseCalculation: calculations.BaseCalculation{
				CalculationParams: &calculations.CalculationParams{AnalyticsConfig: analyticsConfig},
			},
			Direction: d,
		})
	}
	calcs = append(calcs, calculations.GetRawMetricsCalculations(analyticsConfig)...)
	return calcs
}
