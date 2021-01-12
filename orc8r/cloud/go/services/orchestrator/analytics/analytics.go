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
	"magma/orc8r/cloud/go/services/analytics"
	"magma/orc8r/cloud/go/services/analytics/calculations"
)

// GetAnalyticsCalculations returns all calculations computed by the component
func GetAnalyticsCalculations(config *calculations.AnalyticsConfig) []calculations.Calculation {
	if config == nil {
		return nil
	}

	calcs := analytics.GetLogMetricsCalculations(config)
	calcs = append(calcs, analytics.GetRawMetricsCalculations(config)...)
	return calcs
}
