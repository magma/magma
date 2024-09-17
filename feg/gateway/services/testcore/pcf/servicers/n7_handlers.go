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
	"net/http"

	"github.com/golang/glog"
	"github.com/labstack/echo/v4"

	n7_server "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControlServer"
)

// PostSmPolicies handles POST /sm-policies
func (srv *MockPCFServer) PostSmPolicies(ctx echo.Context) error {
	var context n7_server.SmPolicyContextData
	err := ctx.Bind(&context)
	if err != nil {
		glog.Errorf("PostSmPolicies error binding the request body: %s", err)
		return ctx.NoContent(http.StatusBadRequest)
	}
	if srv.serviceConfig.UseMockDriver {
		return srv.getAnswerFromExpectations(ctx, &context, http.StatusCreated)
	}

	decision, err := srv.fetchPolicyDecision(string(context.Supi), uint32(context.PduSessionId))
	if err != nil {
		glog.Errorf("PostSmPolicies unable to fetch account: %s", err)
		return ctx.NoContent(http.StatusNotFound)
	}
	policy := n7_server.SmPolicyControl{
		Context: context,
		Policy:  *decision,
	}
	policyId := srv.createPolicySession(&context, &policy)
	ctx.Response().Header()["Location"] = []string{srv.getSmPolicyUrl(policyId)}

	return ctx.JSON(http.StatusCreated, policy)
}

// GetSmPoliciesSmPolicyId handles GET /sm-policies/{smPolicyId}
func (srv *MockPCFServer) GetSmPoliciesSmPolicyId(ctx echo.Context, smPolicyId string) error {
	sess, found := srv.policySessions[smPolicyId]
	if !found {
		glog.Errorf("GetSmPoliciesSmPolicyId failed to fetch policy session for id %s", smPolicyId)
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, sess.policy)
}

// PostSmPoliciesSmPolicyIdDelete handles POST /sm-policies/{smPolicyId}/delete
func (srv *MockPCFServer) PostSmPoliciesSmPolicyIdDelete(ctx echo.Context, smPolicyId string) error {
	var context n7_server.SmPolicyDeleteData
	err := ctx.Bind(&context)
	if err != nil {
		glog.Errorf("PostSmPoliciesSmPolicyIdDelete error binding the request body: %s", err)
		return ctx.NoContent(http.StatusBadRequest)
	}
	if srv.serviceConfig.UseMockDriver {
		return srv.getAnswerFromExpectations(ctx, &context, http.StatusNoContent)
	}
	err = srv.deletePolicySession(smPolicyId)
	if err != nil {
		glog.Errorf("PostSmPoliciesSmPolicyIdDelete for policy %s failed: %s", smPolicyId, err)
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.NoContent(http.StatusNoContent)
}

// PostSmPoliciesSmPolicyIdUpdate handles POST /sm-policies/{smPolicyId}/update
func (srv *MockPCFServer) PostSmPoliciesSmPolicyIdUpdate(ctx echo.Context, smPolicyId string) error {
	var context n7_server.SmPolicyUpdateContextData
	err := ctx.Bind(&context)
	if err != nil {
		glog.Errorf("PostSmPoliciesSmPolicyIdUpdate error binding the request body: %s", err)
		return ctx.NoContent(http.StatusBadRequest)
	}
	if srv.serviceConfig.UseMockDriver {
		return srv.getAnswerFromExpectations(ctx, &context, http.StatusOK)
	}
	sess, found := srv.policySessions[smPolicyId]
	if !found {
		glog.Errorf("PostSmPoliciesSmPolicyIdUpdate failed to fetch policy session for id %s", smPolicyId)
		return ctx.NoContent(http.StatusNotFound)
	}
	decision, err := srv.fetchPolicyDecision(sess.imsi, sess.pduSessionId)
	if err != nil {
		glog.Errorf("PostSmPolicies unable to fetch account: %s", err)
		return ctx.NoContent(http.StatusNotFound)
	}
	sess.policy.Policy = *decision
	return ctx.JSON(http.StatusOK, decision)
}

func (srv *MockPCFServer) getAnswerFromExpectations(ctx echo.Context, request interface{}, statusCode int) error {
	srv.mockDriver.Lock()
	ans := srv.mockDriver.GetAnswerFromExpectations(request)
	srv.mockDriver.Unlock()
	if ans == nil {
		glog.Errorf("Error getting answer from mock driver")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if len(ans.(string)) == 0 {
		// empty response
		return ctx.NoContent(http.StatusNoContent)
	}
	// ctx.Response().Header()["Location"] = []string{srv.getSmPolicyUrl(policyId)}
	return ctx.JSONBlob(statusCode, []byte(ans.(string)))
}
