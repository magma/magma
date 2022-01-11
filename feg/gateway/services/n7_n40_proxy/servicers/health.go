/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"context"
	"fmt"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/n7_n40_proxy/metrics"
	orcprotos "magma/orc8r/lib/go/protos"
)

// Disable closes all existing pcf connections and disables
// connection creation for the time specified in the request
func (srv *CentralSessionController) Disable(ctx context.Context, req *fegprotos.DisableMessage) (*orcprotos.Void, error) {
	// PCF Connections are stateless HTTP connections, nothing to do here
	return &orcprotos.Void{}, nil
}

// Enable enables pcf connection creation
func (srv *CentralSessionController) Enable(ctx context.Context, void *orcprotos.Void) (*orcprotos.Void, error) {
	// PCF Connections are stateless HTTP connections, nothing to do here
	return &orcprotos.Void{}, nil
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
func (srv *CentralSessionController) GetHealthStatus(ctx context.Context, void *orcprotos.Void) (*fegprotos.HealthStatus, error) {
	currentMetrics, err := metrics.GetCurrentHealthMetrics()
	if err != nil {
		return &fegprotos.HealthStatus{
			Health:        fegprotos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("Error occurred while retrieving health metrics on N7N40Proxy: %s", err),
		}, err
	}
	deltaMetrics, err := srv.healthTracker.Metrics.GetDelta(currentMetrics)
	if err != nil {
		return &fegprotos.HealthStatus{
			Health:        fegprotos.HealthStatus_UNHEALTHY,
			HealthMessage: err.Error(),
		}, err
	}
	n7ReqTotal := deltaMetrics.SmPolicyCreateTotal + deltaMetrics.SmpolicyCreateFailures +
		deltaMetrics.SmPolicyUpdateTotal + deltaMetrics.SmPolicyUpdateFailures +
		deltaMetrics.SmPolicyDeleteTotal + deltaMetrics.SmPolicyDeleteFailures
	n7FailureTotal := deltaMetrics.SmpolicyCreateFailures + deltaMetrics.SmPolicyUpdateFailures +
		deltaMetrics.SmPolicyDeleteFailures + deltaMetrics.N7Timeouts

	n7Status := srv.getHealthStatusForN7Requests(n7FailureTotal, n7ReqTotal)
	if n7Status.Health == fegprotos.HealthStatus_UNHEALTHY {
		return n7Status, nil
	}

	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "All metrics appear healthy",
	}, nil
}

func (srv *CentralSessionController) getHealthStatusForN7Requests(failures, total int64) *fegprotos.HealthStatus {
	if !srv.config.DisableN7 {
		n7ExceedsThreshold := total >= int64(srv.healthTracker.MinimumRequestThreshold) &&
			float64(failures)/float64(total) >= float64(srv.healthTracker.RequestFailureThreshold)
		if n7ExceedsThreshold {
			unhealthyMsg := fmt.Sprintf("Metric N7 Request Failure Ratio >= threshold %f; %d / %d",
				srv.healthTracker.RequestFailureThreshold,
				failures,
				total,
			)
			return &fegprotos.HealthStatus{
				Health:        fegprotos.HealthStatus_UNHEALTHY,
				HealthMessage: unhealthyMsg,
			}
		}
	}
	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "N7 metrics appear healthy",
	}
}
