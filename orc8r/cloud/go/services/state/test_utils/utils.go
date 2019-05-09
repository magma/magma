package test_utils

import (
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	middlewareTests "magma/orc8r/cloud/go/service/middleware/unary/interceptors/tests"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	stateh "magma/orc8r/cloud/go/services/state/obsidian/handlers"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func ReportGatewayStatus(t *testing.T, ctx context.Context, req *checkind_models.GatewayStatus) {
	conn, err := registry.GetConnection(state.ServiceName)
	assert.NoError(t, err)
	client := protos.NewStateServiceClient(conn)

	serializedGWStatus, err := serde.Serialize(state.SerdeDomain, "gw_state", req)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     "gw_state",
			DeviceID: req.HardwareID,
			Value:    serializedGWStatus,
		},
	}
	_, err = client.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
}

func GetContextWithCertificate(t *testing.T, hwID string) context.Context {
	csn := middlewareTests.StartMockGwAccessControl(t, []string{hwID})
	return metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(identity.CLIENT_CERT_SN_KEY, csn[0]))
}

func GetGWStatusViaHTTPNoError(t *testing.T, url string, networkID string, key string) {
	status, response, err := tests.SendHttpRequest("GET", url, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	expected, err := stateh.GetGWStatus(networkID, key)
	assert.NoError(t, err)
	expectedJSON, err := json.Marshal(*expected)
	assert.NoError(t, err)
	assert.Equal(t, string(expectedJSON), response)
}

func GetGWStatusExpectNotFound(t *testing.T, url string) {
	status, _, err := tests.SendHttpRequest("GET", url, "")
	assert.NoError(t, err)
	assert.Equal(t, 404, status)
}
