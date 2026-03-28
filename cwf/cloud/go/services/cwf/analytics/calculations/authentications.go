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

	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"
)

// AuthenticationsCalculation holds the parameters needed to run an authentication
// query and the registered prometheus gauge that the resulting value should be stored in
type AuthenticationsCalculation struct {
	calculations.BaseCalculation
}

// Calculate returns the number of authentications over the past X days segmented
// by result code and networkID
func (x *AuthenticationsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.Infof("Calculating Authentications. Days: %d", x.Days)

	query := fmt.Sprintf(`sum(increase(eap_auth[%dd])) by (code, %s)`, x.Days, metrics.NetworkLabelName)

	vec, err := query_api.QueryPrometheusVector(prometheusClient, query)
	if err != nil {
		return nil, fmt.Errorf("user Consumption query error: %s", err)
	}

	results := calculations.MakeVectorResults(vec, x.Labels, x.Name)
	return results, nil
}
