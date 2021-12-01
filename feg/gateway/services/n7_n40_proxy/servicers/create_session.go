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

	"github.com/golang/glog"

	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/lte/cloud/go/protos"
)

// getSmPolicyRules sends the request to PCF and get the response when N7 is enabled,
// otherwise returns an empty SmPolicyDecision
// Returns SmPolicyDecision rules and smPolicyUrl that uniquely identifies the policy session.
func (srv *CentralSessionController) getSmPolicyRules(
	request *protos.CreateSessionRequest,
) (*n7_sbi.SmPolicyDecision, string, error) {
	if srv.cfg.N7Config.DisableN7 {
		// Empty response when disabled
		return &n7_sbi.SmPolicyDecision{
			PccRules: &n7_sbi.SmPolicyDecision_PccRules{},
		}, "", nil
	}
	// Convert and send the request to PCF
	reqBody := n7.GetSmPolicyContextDataN7(request, srv.cfg.N7Config.Client.NotifyApiRoot)
	reqCtx, cancel := context.WithTimeout(context.Background(), srv.cfg.RequestTimeout)
	defer cancel()
	resp, err := srv.policyClient.PostSmPoliciesWithResponse(reqCtx, *reqBody)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("SmPolicyCreate request failure: status-code=%d", resp.StatusCode())
		return nil, "", err
	}
	policyUrl, found := resp.HTTPResponse.Header["Location"]
	if !found {
		err = fmt.Errorf("SmPolicyCreate request failure: Location header not found")
		return nil, "", err
	}
	return resp.JSON201, policyUrl[0], nil
}

func (srv *CentralSessionController) injectOmnipresentRules(policy *n7_sbi.SmPolicyDecision) error {
	if policy == nil || policy.PccRules == nil {
		return fmt.Errorf("policy decision or pcc rules cannot be nil")
	}
	// No Base-Names returned in N7 response, however since these are statically configured,
	// fetch all the omnipresent rules including the ones referred using the base-names. Will
	// achieve the same behavior for both Gx and N7.
	omnipresentRuleIDs, omnipresentBaseNames := srv.dbClient.GetOmnipresentRules()
	if len(omnipresentRuleIDs) == 0 {
		return nil
	}
	glog.V(2).Infof("Adding omnipresent rules %v omnipresent basenames %v", omnipresentRuleIDs, omnipresentBaseNames)
	baseNameRuleIds := srv.dbClient.GetRuleIDsForBaseNames(omnipresentBaseNames)
	omnipresentRuleIDs = append(omnipresentRuleIDs, baseNameRuleIds...)
	for _, ruleId := range omnipresentRuleIDs {
		policy.PccRules.Set(ruleId, n7_sbi.PccRule{
			PccRuleId: ruleId,
		})
	}
	return nil
}

func validateCreateSessionRequest(req *protos.CreateSessionRequest) error {
	subscriber := req.GetCommonContext().GetSid()
	if subscriber == nil || subscriber.GetId() == "" {
		return fmt.Errorf("missing subscriber information on create session request %+v", req)
	}
	if req.SessionId == "" {
		return fmt.Errorf("missing sessionId information on create session request %+v", req)
	}
	return nil
}
