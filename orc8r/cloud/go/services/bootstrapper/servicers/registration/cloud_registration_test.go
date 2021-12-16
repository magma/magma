/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers/registration"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

var (
	rootCA          = "rootCA"
	timeoutDuration = 30 * time.Minute

	networkID = "networkID"
	logicalID = "logicalID"

	gatewayDeviceInfo = &protos.GatewayDeviceInfo{
		NetworkId: networkID,
		LogicalId: logicalID,
	}
)

func TestCloudRegistrationServicer_GetGatewayRegistrationInfo(t *testing.T) {
	ctx, cloudRegistration := cloudRegistrationTestSetup(t)

	getGatewayRegistrationInfoRes, err := cloudRegistration.GetGatewayRegistrationInfo(ctx, &protos.GetGatewayRegistrationInfoRequest{})
	expectedRes := &protos.GetGatewayRegistrationInfoResponse{
		RootCa:     rootCA,
		DomainName: registration.NotImplementedWarning,
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedRes, getGatewayRegistrationInfoRes)
}

// TestCloudRegistrationServicer_Registration tests GetGatewayDeviceInfo and GetToken
// Tests that the functions interplay together with expected behavior
func TestCloudRegistrationServicer_Registration(t *testing.T) {
	ctx, cloudRegistration := cloudRegistrationTestSetup(t)

	nonce := registration.GenerateNonce(registration.NonceLength)
	assert.Equal(t, registration.NonceLength, len(nonce))
	token := registration.NonceToToken(nonce)

	// Try getting device info when token is invalid
	getGatewayDeviceInfoRes, err := cloudRegistration.GetGatewayDeviceInfo(ctx, &protos.GetGatewayDeviceInfoRequest{Token: token})
	expectedErrorRes := &protos.GetGatewayDeviceInfoResponse_Error{
		Error: fmt.Sprintf("could not get token info from token %v: %v", token, "Not found"),
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedErrorRes, getGatewayDeviceInfoRes.Response)

	// Register some device info
	getTokenRes, err := cloudRegistration.GetToken(ctx, &protos.GetTokenRequest{
		GatewayDeviceInfo: gatewayDeviceInfo,
		Refresh:           false,
	})
	assert.NoError(t, err)
	assert.True(t, clock.Now().Before(registration.GetTime(getTokenRes.Timeout)))
	token = getTokenRes.Token
	timeout := getTokenRes.Timeout

	// Get device info
	getGatewayDeviceInfoRes, err = cloudRegistration.GetGatewayDeviceInfo(ctx, &protos.GetGatewayDeviceInfoRequest{Token: token})
	assert.NoError(t, err)
	assert.Equal(t, &protos.GetGatewayDeviceInfoResponse_GatewayDeviceInfo{GatewayDeviceInfo: gatewayDeviceInfo}, getGatewayDeviceInfoRes.Response)

	// Refresh device info
	getTokenRes, err = cloudRegistration.GetToken(ctx, &protos.GetTokenRequest{
		GatewayDeviceInfo: gatewayDeviceInfo,
		Refresh:           true,
	})
	newToken := getTokenRes.Token
	assert.NoError(t, err)
	assert.True(t, registration.GetTime(timeout).Before(registration.GetTime(getTokenRes.Timeout)))
	assert.NotEqual(t, token, newToken)

	// Get device info with new token
	getGatewayDeviceInfoRes, err = cloudRegistration.GetGatewayDeviceInfo(ctx, &protos.GetGatewayDeviceInfoRequest{Token: newToken})
	assert.NoError(t, err)
	assert.Equal(t, &protos.GetGatewayDeviceInfoResponse_GatewayDeviceInfo{GatewayDeviceInfo: gatewayDeviceInfo}, getGatewayDeviceInfoRes.Response)

	// Old token should still work
	getGatewayDeviceInfoRes, err = cloudRegistration.GetGatewayDeviceInfo(ctx, &protos.GetGatewayDeviceInfoRequest{Token: token})
	assert.NoError(t, err)
	assert.Equal(t, &protos.GetGatewayDeviceInfoResponse_GatewayDeviceInfo{GatewayDeviceInfo: gatewayDeviceInfo}, getGatewayDeviceInfoRes.Response)

	// Test that token expires properly
	clock.SetAndFreezeClock(t, time.Now().Add(timeoutDuration))
	getGatewayDeviceInfoRes, err = cloudRegistration.GetGatewayDeviceInfo(ctx, &protos.GetGatewayDeviceInfoRequest{Token: token})
	expectedErrorRes = &protos.GetGatewayDeviceInfoResponse_Error{
		Error: fmt.Sprintf("token %v has expired", token),
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedErrorRes, getGatewayDeviceInfoRes.Response)
}

func cloudRegistrationTestSetup(t *testing.T) (context.Context, protos.CloudRegistrationServer) {
	factory := test_utils.NewSQLBlobstore(t, bootstrapper.BlobstoreTableName)
	store := registration.NewBlobstoreStore(factory)

	cloudRegistration, err := registration.NewCloudRegistrationServicer(store, rootCA, timeoutDuration, true)
	assert.NoError(t, err)

	return context.Background(), cloudRegistration
}
