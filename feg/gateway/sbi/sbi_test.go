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
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	n7_client "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	n7 "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControlServer"
	sbi "magma/feg/gateway/sbi/specs/TS29571CommonData"
)

const (
	POLICY_ID1 = "1234"
	IMSI1      = "12345678901234"
	IP_ADDR1   = "10.1.2.3"
)

type MockPcf struct {
	ListenAddr string
	Policies   map[string]n7.SmPolicyControl
}

func NewMockPcf() *MockPcf {
	return &MockPcf{
		Policies: make(map[string]n7.SmPolicyControl),
	}
}

func (pcf *MockPcf) PostSmPolicies(ctx echo.Context) error {
	var newPolicy n7.SmPolicyContextData
	err := ctx.Bind(&newPolicy)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var policy n7.SmPolicyControl
	policy.Context = newPolicy

	pcf.Policies[POLICY_ID1] = policy

	return ctx.NoContent(http.StatusOK)
}

// GetSmPoliciesSmPolicyId handles GET /sm-policies/{smPolicyId}
func (pcf *MockPcf) GetSmPoliciesSmPolicyId(ctx echo.Context, smPolicyId string) error {
	policy, found := pcf.Policies[smPolicyId]
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
	mockPcf := startServer(t)

	// Create new N7 client object
	client, err := n7_client.NewClientWithResponses(fmt.Sprintf("http://%s", mockPcf.ListenAddr))

	require.NoError(t, err, "Client creation failed")

	// Create a new SM Policy
	ipAddr := sbi.Ipv4Addr(IP_ADDR1)
	imsi := sbi.Supi(IMSI1)
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
	assert.Equal(t, sbi.Ipv4Addr(IP_ADDR1), *(policyResp.JSON200.Context.Ipv4Address))
	assert.Equal(t, sbi.Supi(IMSI1), policyResp.JSON200.Context.Supi)
}

func startServer(t *testing.T) *MockPcf {
	errChan := make(chan error)
	eserver := echo.New()

	mockPcf := NewMockPcf()
	go func() {
		n7.RegisterHandlers(eserver, mockPcf)

		err := eserver.Start("localhost:0")
		if err != nil {
			errChan <- err
		}
	}()

	var err error
	mockPcf.ListenAddr, err = waitForTestServerAndGetAddr(eserver, errChan)
	require.NoError(t, err)
	return mockPcf
}

// waitForTestServer waits for the Echo server to be launched and returns the Listener Address.
func waitForTestServerAndGetAddr(e *echo.Echo, errChan <-chan error) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			addr := e.ListenerAddr()
			if addr != nil && strings.Contains(addr.String(), ":") {
				// server started
				return addr.String(), nil
			}
		case err := <-errChan:
			if err == http.ErrServerClosed {
				return "", nil
			}
			return "", err
		}
	}
}
