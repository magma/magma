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

// XAPCalculation holds the parameters needed to run a XAP query and the registered
// prometheus gauge that the resulting value should be stored in
type XAPCalculation struct {
	calculations.BaseCalculation
	ThresholdBytes int
}

// Calculate returns the number of unique users who have had a session in the
// past X days and have used over `thresholdBytes` data in that time
func (x *XAPCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.Infof("Calculating XAP. Days: %d", x.Days)

	// List the users who have had an active session over the last X days
	uniqueUsersQuery := fmt.Sprintf(`count(max_over_time(active_sessions[%dd]) >= 1) by (imsi,networkID)`, x.Days)
	// List the users who have used at least x.ThresholdBytes of data in the last X days
	usersOverThresholdQuery := fmt.Sprintf(`count(sum(increase(octets_in[%dd])) by (imsi,networkID) > %d) by (imsi,networkID)`, x.Days, x.ThresholdBytes)
	// Count the users who match both conditions
	intersectionQuery := fmt.Sprintf(`count(%s and %s) by (%s)`, uniqueUsersQuery, usersOverThresholdQuery, metrics.NetworkLabelName)

	vec, err := query_api.QueryPrometheusVector(prometheusClient, intersectionQuery)
	if err != nil {
		return nil, fmt.Errorf("user Consumption query error: %s", err)
	}

	results := calculations.MakeVectorResults(vec, x.Labels, x.Name)
	return results, nil
}
