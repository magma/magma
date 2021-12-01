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
	"time"

	"github.com/golang/glog"

	"magma/feg/gateway/policydb"
	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/n7_n40_proxy/metrics"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/lte/cloud/go/protos"
)

const (
	DefaultN7Timeout = 3 * time.Second
)

type CentralSessionController struct {
	policyClient  n7_sbi.ClientWithResponsesInterface
	dbClient      policydb.PolicyDBClient
	cfg           *SessionControllerConfig
	healthTracker *metrics.SessionHealthTracker
}

type SessionControllerConfig struct {
	N7Config       *n7.N7Config
	RequestTimeout time.Duration
}

func NewCentralSessionController(
	n7config *n7.N7Config,
	dbClient policydb.PolicyDBClient,
	policyClient n7_sbi.ClientWithResponsesInterface,
) (*CentralSessionController, error) {
	return &CentralSessionController{
		policyClient: policyClient,
		dbClient:     dbClient,
		cfg: &SessionControllerConfig{
			N7Config:       n7config,
			RequestTimeout: DefaultN7Timeout,
		},
		healthTracker: metrics.NewSessionHealthTracker(),
	}, nil
}

// CreateSession begins a UE session by requesting rules from PCF and returning them.
func (srv *CentralSessionController) CreateSession(
	ctx context.Context,
	request *protos.CreateSessionRequest,
) (*protos.CreateSessionResponse, error) {
	if err := validateCreateSessionRequest(request); err != nil {
		return nil, err
	}

	policy, policyId, err := srv.getSmPolicyRules(request)
	metrics.ReportCreateSmPolicy(err)
	if err != nil {
		return nil, err
	}
	err = srv.injectOmnipresentRules(policy)
	if err != nil {
		glog.Errorf("Failed to inject omnipresent rules %s", err)
	}
	return n7.GetCreateSessionResponseProto(request, policy, policyId), nil
}

// UpdateSession handles periodic updates from gateways that include quota
// exhaustion and terminations.
func (srv *CentralSessionController) UpdateSession(
	ctx context.Context,
	request *protos.UpdateSessionRequest,
) (*protos.UpdateSessionResponse, error) {

	return (&protos.UnimplementedCentralSessionControllerServer{}).UpdateSession(ctx, request)
}

// TerminateSession handles a session termination
func (srv *CentralSessionController) TerminateSession(
	ctx context.Context,
	request *protos.SessionTerminateRequest,
) (*protos.SessionTerminateResponse, error) {

	return (&protos.UnimplementedCentralSessionControllerServer{}).TerminateSession(ctx, request)
}
