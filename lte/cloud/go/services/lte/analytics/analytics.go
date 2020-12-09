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
	lte_service "magma/lte/cloud/go/services/lte"
	lte_calculations "magma/lte/cloud/go/services/lte/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/calculations"
)

//GetAnalyticsCalculations gets all the calculations provided by lte analytics
func GetAnalyticsCalculations(config *lte_service.Config) []calculations.Calculation {
	if config == nil {
		return nil
	}

	allCalculations := make([]calculations.Calculation, 0)
	allCalculations = append(allCalculations, &lte_calculations.GeneralMetricsCalculation{
		CalculationParams: calculations.CalculationParams{
			MetricConfig: config.Analytics.Metrics,
		},
	})
	allCalculations = append(allCalculations, &lte_calculations.UserMetricsCalculation{
		CalculationParams: calculations.CalculationParams{
			MetricConfig: config.Analytics.Metrics,
		},
	})
	allCalculations = append(allCalculations, &lte_calculations.SiteMetricsCalculation{
		CalculationParams: calculations.CalculationParams{
			MetricConfig: config.Analytics.Metrics,
		},
	})
	allCalculations = append(allCalculations, calculations.GetRawMetricsCalculations(config.Analytics.Metrics)...)
	return allCalculations
}
