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

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/services/bootstrapper/servicers/registration"
	"magma/orc8r/cloud/go/services/tenants"
	tenantsTestInit "magma/orc8r/cloud/go/services/tenants/test_init"
	"magma/orc8r/lib/go/protos"
)

var (
	registerRequest = &protos.RegisterRequest{
		Token: registration.NonceToToken(registration.GenerateNonce(registration.NonceLength)),
		Hwid: &protos.AccessGatewayID{
			Id: "Id",
		},
		ChallengeKey: &protos.ChallengeKey{
			KeyType: 0,
			Key:     []byte("key"),
		},
	}
	controlProxy       = "controlProxy"
	nextTenantID int64 = 0
)

func TestRegistrationServicer_Register(t *testing.T) {
	registrationServicer := setupMockRegistrationServicer()

	res, err := registrationServicer.Register(context.Background(), registerRequest)
	assert.NoError(t, err)
	expectedRes := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_ControlProxy{ControlProxy: controlProxy},
	}
	assert.Equal(t, expectedRes, res)
}

func TestRegistrationServicer_Register_BadToken(t *testing.T) {
	rpcErr := status.Error(codes.NotFound, "errMessage")

	registrationServicer := setupMockRegistrationServicer()
	registrationServicer.GetGatewayDeviceInfo = func(ctx context.Context, token string) (*protos.GatewayDeviceInfo, error) {
		return nil, rpcErr
	}

	res, err := registrationServicer.Register(context.Background(), registerRequest)
	assert.NoError(t, err)
	expectedRes := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_Error{
			Error: fmt.Sprintf("could not get device info from token %v: %v", registerRequest.Token, rpcErr),
		},
	}
	assert.Equal(t, expectedRes, res)
}

func TestRegistrationServicer_Register_NoControlProxy(t *testing.T) {
	rpcErr := status.Error(codes.NotFound, "errMessage")

	registrationServicer := setupMockRegistrationServicer()
	registrationServicer.GetControlProxy = func(networkID string) (string, error) {
		return "", rpcErr
	}

	res, err := registrationServicer.Register(context.Background(), registerRequest)
	assert.NoError(t, err)
	expectedRes := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_Error{
			Error: fmt.Sprintf("error getting control proxy: %v", rpcErr),
		},
	}
	assert.Equal(t, expectedRes, res)
}

func TestGetControlProxy_NoNetworkID(t *testing.T) {
	setupAddNetworksToTenantsService(t)

	res, err := registration.GetControlProxy(networkID)
	assert.Equal(t, status.Errorf(codes.NotFound, "tenantID for current NetworkID %v not found", networkID), err)
	assert.Equal(t, "", res)
}

func TestGetControlProxy_NoControlProxy(t *testing.T) {
	setupAddNetworksToTenantsService(t)

	networkIDTenant := &protos.Tenant{
		Name:     "tenant",
		Networks: []string{networkID},
	}
	addTenant(t, networkIDTenant)

	res, err := registration.GetControlProxy(networkID)
	assert.Equal(t, "Not found", err.Error())
	assert.Equal(t, "", res)
}

func TestGetControlProxy(t *testing.T) {
	setupAddNetworksToTenantsService(t)

	networkIDTenant := &protos.Tenant{
		Name:     "tenant",
		Networks: []string{networkID},
	}
	id := addTenant(t, networkIDTenant)
	err := tenants.CreateOrUpdateControlProxy(context.Background(), protos.CreateOrUpdateControlProxyRequest{
		Id:           id,
		ControlProxy: controlProxy,
	})
	assert.NoError(t, err)

	res, err := registration.GetControlProxy(networkID)
	assert.NoError(t, err)
	assert.Equal(t, controlProxy, res)
}

func setupMockRegistrationServicer() *registration.RegistrationService {
	registrationService := &registration.RegistrationService{
		GetGatewayDeviceInfo: func(ctx context.Context, token string) (*protos.GatewayDeviceInfo, error) {
			return gatewayDeviceInfo, nil
		},
		RegisterDevice: func(deviceInfo protos.GatewayDeviceInfo, hwid *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) error {
			return nil
		},
		GetControlProxy: func(networkID string) (string, error) {
			return controlProxy, nil
		},
	}

	return registrationService
}

func setupAddNetworksToTenantsService(t *testing.T) {
	var (
		tenant1 = &protos.Tenant{
			Name:     "tenant",
			Networks: []string{"network1", "network2"},
		}
		tenant2 = &protos.Tenant{
			Name:     "tenant",
			Networks: []string{"network3", "network4"},
		}
	)
	tenantsTestInit.StartTestService(t)

	addTenant(t, tenant1)
	addTenant(t, tenant2)
}

func addTenant(t *testing.T, tenant *protos.Tenant) int64 {
	ctx := context.Background()

	tenantRes, err := tenants.CreateTenant(ctx, nextTenantID, tenant)
	assert.NoError(t, err)
	assert.Equal(t, tenant, tenantRes)

	curTenantID := nextTenantID
	nextTenantID = nextTenantID + 1
	return curTenantID
}
