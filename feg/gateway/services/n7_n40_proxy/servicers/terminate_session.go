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
	"net/http"

	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/lte/cloud/go/protos"
)

func validateSessionTerminateRequest(req *protos.SessionTerminateRequest) error {
	subscriber := req.GetCommonContext().GetSid()
	if subscriber == nil || subscriber.GetId() == "" {
		return fmt.Errorf("missing subscriber information on create session request %+v", req)
	}
	if req.SessionId == "" {
		return fmt.Errorf("missing magma sessionId information on create session request %+v", req)
	}
	return nil
}

func (srv *CentralSessionController) sendSmPolicyDelete(
	smPolicyId string,
	reqBody *n7_sbi.PostSmPoliciesSmPolicyIdDeleteJSONRequestBody,
) error {
	reqCtx, cancel := context.WithTimeout(context.Background(), srv.cfg.RequestTimeout)
	defer cancel()
	resp, err := srv.policyClient.PostSmPoliciesSmPolicyIdDeleteWithResponse(reqCtx, smPolicyId, *reqBody)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("SmPolicyDelete request failure: status-code=%d policy-id=%s", resp.StatusCode(), smPolicyId)
	}
	return nil
}
