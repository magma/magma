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
	"context"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"

	"github.com/golang/glog"
)

//CollectorService implements the operations of collecting the metrics from CWF service
type CollectorService struct {
	calculations  []calculations.Calculation
	promAPIClient query_api.PrometheusAPI
}

//NewCollectorService constructs new collector service
func NewCollectorService(promAPIClient query_api.PrometheusAPI, calculations []calculations.Calculation) *CollectorService {
	return &CollectorService{promAPIClient: promAPIClient, calculations: calculations}
}

//Collect does the operation of running through calculations and returning results
func (svc *CollectorService) Collect(context.Context, *protos.CollectRequest) (*protos.CollectResponse, error) {
	response := &protos.CollectResponse{}
	for _, calc := range svc.calculations {
		results, err := calc.Calculate(svc.promAPIClient)
		if err != nil {
			glog.Errorf("Error calculating metric: %s", err)
			continue
		}
		glog.V(10).Infof("results are %v", results)
		response.Results = append(response.Results, results...)
	}
	glog.V(10).Infof("response is %v", response)
	return response, nil
}
