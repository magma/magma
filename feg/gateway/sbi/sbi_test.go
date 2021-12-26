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

package sbi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/feg/gateway/sbi"
	n7_client "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	n7 "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControlServer"
	sbi_common "magma/feg/gateway/sbi/specs/TS29571CommonData"
)

const (
	POLICY_ID1      = "1234"
	IMSI1           = "12345678901234"
	IP_ADDR1        = "10.1.2.3"
	LISTEN_ADDR     = "localhost:0"
	ShutdownTimeout = 10 * time.Second
)

type MockPcf struct {
	sbiServer *sbi.SbiServer
	policies  map[string]n7.SmPolicyControl
}

func NewMockPcf() (*MockPcf, error) {
	mockPcf := &MockPcf{
		sbiServer: sbi.NewSbiServer(LISTEN_ADDR),
		policies:  make(map[string]n7.SmPolicyControl),
	}
	n7.RegisterHandlers(mockPcf.sbiServer.Server, mockPcf)
	err := mockPcf.sbiServer.Start()
	if err != nil {
		return nil, err
	}
	return mockPcf, nil
}

func (pcf *MockPcf) PostSmPolicies(ctx echo.Context) error {
	var newPolicy n7.SmPolicyContextData
	err := ctx.Bind(&newPolicy)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var policy n7.SmPolicyControl
	policy.Context = newPolicy

	pcf.policies[POLICY_ID1] = policy

	return ctx.NoContent(http.StatusOK)
}

// GetSmPoliciesSmPolicyId handles GET /sm-policies/{smPolicyId}
func (pcf *MockPcf) GetSmPoliciesSmPolicyId(ctx echo.Context, smPolicyId string) error {
	policy, found := pcf.policies[smPolicyId]
	if !found {
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, policy)
}

// PostSmPoliciesSmPolicyIdDelete handles POST /sm-policies/{smPolicyId}/delete
func (pcf *MockPcf) PostSmPoliciesSmPolicyIdDelete(ctx echo.Context, smPolicyId string) error {

	return ctx.NoContent(http.StatusOK)
}

// PostSmPoliciesSmPolicyIdUpdate handles POST /sm-policies/{smPolicyId}/update
func (pcf *MockPcf) PostSmPoliciesSmPolicyIdUpdate(ctx echo.Context, smPolicyId string) error {

	return ctx.NoContent(http.StatusOK)
}

func TestCreateSMPolicy(t *testing.T) {
	// Start the server and wait till the server is started
	mockPcf, err := NewMockPcf()
	require.NoError(t, err, "Error starting mock server")
	defer mockPcf.sbiServer.Shutdown(ShutdownTimeout)

	// Create new N7 client object
	client, err := n7_client.NewClientWithResponses(fmt.Sprintf("http://%s", mockPcf.sbiServer.ListenAddr))

	require.NoError(t, err, "Client creation failed")

	// Create a new SM Policy
	ipAddr := sbi_common.Ipv4Addr(IP_ADDR1)
	imsi := sbi_common.Supi(IMSI1)
	body := n7_client.PostSmPoliciesJSONRequestBody{
		Ipv4Address:  &ipAddr,
		PduSessionId: 10,
		Supi:         imsi,
	}
	resp, err := client.PostSmPolicies(context.Background(), body)
	require.NoError(t, err, "Failed to create SM Policies")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Get SM Policy given the policyId
	policyResp, err := client.GetSmPoliciesSmPolicyIdWithResponse(context.Background(), POLICY_ID1)
	require.NoError(t, err, "Failed to get SM Policies")
	assert.Equal(t, http.StatusOK, policyResp.StatusCode())
	assert.Equal(t, sbi_common.Ipv4Addr(IP_ADDR1), *(policyResp.JSON200.Context.Ipv4Address))
	assert.Equal(t, sbi_common.Supi(IMSI1), policyResp.JSON200.Context.Supi)
}
