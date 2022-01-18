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

	"github.com/golang/glog"

	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/n7_n40_proxy/metrics"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/lte/cloud/go/protos"
)

const (
	DefaultN7Timeout = 3 * time.Second
)

type CentralSessionController struct {
	policyClient  *n7.N7Client
	dbClient      policydb.PolicyDBClient
	config        *SessionControllerConfig
	healthTracker *metrics.SessionHealthTracker
}

type SessionControllerConfig struct {
	DisableN7      bool
	RequestTimeout time.Duration
}

func NewCentralSessionController(
	n7config *n7.N7Config,
	dbClient policydb.PolicyDBClient,
	policyClient *n7.N7Client,
) (*CentralSessionController, error) {

	cfg := &SessionControllerConfig{
		DisableN7:      n7config.DisableN7,
		RequestTimeout: DefaultN7Timeout,
	}
	return &CentralSessionController{
		policyClient:  policyClient,
		dbClient:      dbClient,
		config:        cfg,
		healthTracker: metrics.NewSessionHealthTracker(),
	}, nil
}

// CreateSession begins a UE session by requesting rules from PCF and returning them.
func (srv *CentralSessionController) CreateSession(
	ctx context.Context,
	request *protos.CreateSessionRequest,
) (*protos.CreateSessionResponse, error) {
	if err := validateCreateSessionRequest(request); err != nil {
		err = fmt.Errorf("CreateSessionRequest failed to validate: %s", err)
		glog.Error(err)
		return nil, err
	}

	policy, policyId, err := srv.getSmPolicyRules(request)
	metrics.ReportCreateSmPolicy(err)
	if err != nil {
		err = fmt.Errorf("CreateSessionRequest failed to get SMPolicyRules: %s", err)
		glog.Error(err)
		return nil, err
	}
	err = srv.injectOmnipresentRules(policy)
	if err != nil {
		glog.Errorf("CreateSessionRequest Failed to inject omnipresent rules %s", err)
	}
	return n7.GetCreateSessionResponseProto(request, policy, policyId), nil
}

// UpdateSession handles periodic updates from gateways that include quota
// exhaustion and terminations.
func (srv *CentralSessionController) UpdateSession(
	ctx context.Context,
	request *protos.UpdateSessionRequest,
) (*protos.UpdateSessionResponse, error) {
	reqCtxts := n7.GetSmPolicyUpdateRequestsN7(request.UsageMonitors)
	responses := srv.sendMutlipleSmPolicyUpdateRequests(reqCtxts)
	return &protos.UpdateSessionResponse{
		UsageMonitorResponses: responses,
	}, nil
}

// TerminateSession handles a session termination
func (srv *CentralSessionController) TerminateSession(
	ctx context.Context,
	request *protos.SessionTerminateRequest,
) (*protos.SessionTerminateResponse, error) {
	if err := validateSessionTerminateRequest(request); err != nil {
		err = fmt.Errorf("SessionTerminateRequest failed to validate: %s", err)
		glog.Error(err)
		return nil, err
	}
	smPolicyId, err := n7.GetSmPolicyId(request.GetTgppCtx())
	if err != nil {
		err = fmt.Errorf("TerminateSession failed to get policyId: %s", err)
		glog.Error(err)
		return nil, err
	}
	reqBody := n7.GetSmPolicyDeleteReqBody(request)
	err = srv.sendSmPolicyDelete(smPolicyId, reqBody)
	metrics.ReportDeleteSmPolicy(err)
	if err != nil {
		err = fmt.Errorf("SessionTerminateRequest failed to send SM Policy Delete: %s", err)
		glog.Error(err)
		return nil, err
	}
	return &protos.SessionTerminateResponse{
		Sid:       request.GetCommonContext().GetSid().GetId(),
		SessionId: request.SessionId,
	}, nil
}

// Close gracefully shuts down the CentralSessionController
func (srv *CentralSessionController) Close() {
	srv.policyClient.NotifyServer.Server.Close()
}
