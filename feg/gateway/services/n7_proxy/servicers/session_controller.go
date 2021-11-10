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
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/policydb"
	n7_client "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/n7_proxy/metrics"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"
)

type CentralSessionController struct {
	policyClient  *n7_client.ClientWithResponses
	dbClient      policydb.PolicyDBClient
	cfg           *SessionControllerConfig
	healthTracker *metrics.SessionHealthTracker
}

type SessionControllerConfig struct {
	PCFConfig      *PCFConfig
	RequestTimeout time.Duration
	DisableN7      bool
}

func NewCentralSessionController(
	policyClient *n7_client.ClientWithResponses,
	dbClient policydb.PolicyDBClient,
	cfg *SessionControllerConfig,
) *CentralSessionController {
	return &CentralSessionController{
		policyClient:  policyClient,
		dbClient:      dbClient,
		cfg:           cfg,
		healthTracker: metrics.NewSessionHealthTracker(),
	}
}

// CreateSession begins a UE session by requesting rules from PCF and returning them.
func (srv *CentralSessionController) CreateSession(
	_ context.Context,
	request *protos.CreateSessionRequest,
) (*protos.CreateSessionResponse, error) {
	// TODO convert and make a policyClient call to PCF

	return &protos.CreateSessionResponse{}, nil
}

// UpdateSession handles periodic updates from gateways that include quota
// exhaustion and terminations.
func (srv *CentralSessionController) UpdateSession(
	ctx context.Context,
	request *protos.UpdateSessionRequest,
) (*protos.UpdateSessionResponse, error) {
	// TODO convert and make a policyClient call to PCF

	return &protos.UpdateSessionResponse{}, nil
}

// TerminateSession handles a session termination
func (srv *CentralSessionController) TerminateSession(
	ctx context.Context,
	request *protos.SessionTerminateRequest,
) (*protos.SessionTerminateResponse, error) {
	// TODO convert and make a policyClient call to PCF

	return &protos.SessionTerminateResponse{
		Sid:       request.GetCommonContext().GetSid().GetId(),
		SessionId: request.SessionId,
	}, nil
}

// Disable closes all existing pcf connections and disables
// connection creation for the time specified in the request
func (srv *CentralSessionController) Disable(
	ctx context.Context,
	req *fegprotos.DisableMessage,
) (*orcprotos.Void, error) {
	if req == nil {
		return nil, fmt.Errorf("nil disable request")
	}
	// PCF Connections are stateless HTTP connections. Don't have to disable them.
	if !srv.cfg.DisableN7 {
		// No new requestes are made
		srv.cfg.DisableN7 = true
		disablePeriod := time.Duration(req.DisablePeriodSecs) * time.Second
		time.AfterFunc(disablePeriod, func() { srv.cfg.DisableN7 = false })
	}
	return &orcprotos.Void{}, nil
}

// Enable enables pcf connection creation
func (srv *CentralSessionController) Enable(
	ctx context.Context,
	void *orcprotos.Void,
) (*orcprotos.Void, error) {
	// PCF Connections are stateless HTTP connections.
	srv.cfg.DisableN7 = false
	return &orcprotos.Void{}, nil
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
func (srv *CentralSessionController) GetHealthStatus(
	ctx context.Context,
	void *orcprotos.Void,
) (*fegprotos.HealthStatus, error) {
	currentMetrics, err := metrics.GetCurrentHealthMetrics()
	if err != nil {
		return &fegprotos.HealthStatus{
			Health:        fegprotos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("Error occurred while retrieving health metrics: %s", err),
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
	if !srv.cfg.DisableN7 {
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
