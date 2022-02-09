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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/feg/gateway/sbi"
	"magma/feg/gateway/sbi/mocks"
	sbi_NpcfSMPolicyControl "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	sbi_CommonData "magma/feg/gateway/sbi/specs/TS29571CommonData"
)

const (
	ShutdownTimeout = 10 * time.Second
	IMSI1           = "12345678901234"
	IP_ADDR1        = "10.1.2.3"
	LISTEN_ADDR     = "localhost:0"
)

func TestCreateSMPolicy(t *testing.T) {
	// Start the server and wait till the server is started
	mockPcf, err := mocks.NewMockPcf(LISTEN_ADDR)
	require.NoError(t, err, "Error starting mock server")
	defer mockPcf.Shutdown(ShutdownTimeout)

	pcfAddr, err := mockPcf.GetListenerAddr()
	require.NoError(t, err)

	// Create new N7 client object
	client, err := sbi_NpcfSMPolicyControl.NewClientWithResponses(
		fmt.Sprintf("http://%s", pcfAddr.String()),
		sbi_NpcfSMPolicyControl.WithHTTPClient(sbi.NewLoggingHttpClient()))

	require.NoError(t, err, "BaseClientWithNotifier creation failed")

	// Create a new SM Policy
	ipAddr := sbi_CommonData.Ipv4Addr(IP_ADDR1)
	imsi := sbi_CommonData.Supi(IMSI1)
	body := sbi_NpcfSMPolicyControl.PostSmPoliciesJSONRequestBody{
		Ipv4Address:  &ipAddr,
		PduSessionId: 10,
		Supi:         imsi,
	}
	resp, err := client.PostSmPolicies(context.Background(), body)
	require.NoError(t, err, "Failed to create SM Policies")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Get SM Policy given the policyId
	policyResp, err := client.GetSmPoliciesSmPolicyIdWithResponse(context.Background(), mocks.POLICY_ID1)
	require.NoError(t, err, "Failed to get SM Policies")
	assert.Equal(t, http.StatusOK, policyResp.StatusCode())
	assert.Equal(t, sbi_CommonData.Ipv4Addr(IP_ADDR1), *(policyResp.JSON200.Context.Ipv4Address))
	assert.Equal(t, sbi_CommonData.Supi(IMSI1), policyResp.JSON200.Context.Supi)
}
