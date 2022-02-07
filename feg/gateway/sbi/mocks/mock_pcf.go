package mocks

import (
	"net/http"

	"magma/feg/gateway/sbi"
	sbi_NpcfSMPolicyControlServer "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControlServer"

	"github.com/labstack/echo/v4"
)

const (
	POLICY_ID1 = "1234"
)

type MockPcf struct {
	*sbi.EchoServer
	policies map[string]sbi_NpcfSMPolicyControlServer.SmPolicyControl
}

func NewMockPcf(localAddr string) (*MockPcf, error) {
	mockPcf := &MockPcf{
		EchoServer: sbi.NewEchoServer(),
		policies:   make(map[string]sbi_NpcfSMPolicyControlServer.SmPolicyControl),
	}
	sbi_NpcfSMPolicyControlServer.RegisterHandlers(mockPcf, mockPcf)
	err := mockPcf.StartWithWait(localAddr)
	if err != nil {
		return nil, err
	}
	return mockPcf, nil
}

func (pcf *MockPcf) PostSmPolicies(ctx echo.Context) error {
	var newPolicy sbi_NpcfSMPolicyControlServer.SmPolicyContextData
	err := ctx.Bind(&newPolicy)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var policy sbi_NpcfSMPolicyControlServer.SmPolicyControl
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
