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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers/registration"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/tenants"
	tenant_protos "magma/orc8r/cloud/go/services/tenants/protos"
	tenant_TestInit "magma/orc8r/cloud/go/services/tenants/test_init"
	"magma/orc8r/lib/go/protos"
)

var (
	registerRequest = &protos.RegisterRequest{
		Token: registration.NonceToToken(registration.GenerateNonce(registration.NonceLength)),
		Hwid: &protos.AccessGatewayID{
			Id: hardwareID,
		},
		ChallengeKey: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_ECHO,
			Key:     challengeKey,
		},
	}
	challengeKey       = []byte("MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEQrZVdmuZpvciEXdznTErWUelOcgdBwPKQfOZDL7Wkl8ALSBtKvJWDPyhS6rkW9/xJdgPD4QK3Jqc4Eox5NT6SVYYuHWLv7b28493rwFvuC2+YurmfYj+LZh9VBVTvlwk")
	controlProxy       = "controlProxy"
	nextTenantID int64 = 0
	hardwareID         = "foo-bar-hardware-id"
)

func TestRegistrationServicer_Register(t *testing.T) {
	registrationServicer := setupMockRegistrationServicer(t)

	res, err := registrationServicer.Register(context.Background(), registerRequest)
	assert.NoError(t, err)

	expectedRes := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_ControlProxy{ControlProxy: controlProxy},
	}
	assert.Equal(t, expectedRes, res)

	checkRegisteredGateway(t)
}

func TestRegistrationServicer_Register_BadToken(t *testing.T) {
	rpcErr := status.Error(codes.NotFound, "errMessage")

	registrationServicer := setupMockRegistrationServicer(t)
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

	registrationServicer := setupMockRegistrationServicer(t)
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
	assert.Equal(t, fmt.Errorf("could not get control-proxy from tenant with network ID %s: Not found", networkID).Error(), err.Error())
	assert.Equal(t, "", res)
}

func TestGetControlProxy_NoControlProxy(t *testing.T) {
	setupAddNetworksToTenantsService(t)

	networkIDTenant := &tenant_protos.Tenant{
		Name:     "tenant",
		Networks: []string{networkID},
	}
	addTenant(t, networkIDTenant)

	res, err := registration.GetControlProxy(networkID)
	assert.Equal(t, fmt.Errorf("could not get control-proxy from tenant with network ID %s: Not found", networkID).Error(), err.Error())
	assert.Equal(t, "", res)
}

func TestGetControlProxy(t *testing.T) {
	setupAddNetworksToTenantsService(t)

	networkIDTenant := &tenant_protos.Tenant{
		Name:     "tenant",
		Networks: []string{networkID},
	}
	id := addTenant(t, networkIDTenant)
	err := tenants.CreateOrUpdateControlProxy(context.Background(), &tenant_protos.CreateOrUpdateControlProxyRequest{
		Id:           id,
		ControlProxy: controlProxy,
	})
	assert.NoError(t, err)

	res, err := registration.GetControlProxy(networkID)
	assert.NoError(t, err)
	assert.Equal(t, controlProxy, res)
}

func setupMockRegistrationServicer(t *testing.T) *registration.RegistrationService {
	registrationService := &registration.RegistrationService{
		GetGatewayDeviceInfo: func(ctx context.Context, token string) (*protos.GatewayDeviceInfo, error) {
			return gatewayDeviceInfo, nil
		},
		RegisterDevice: func(deviceInfo *protos.GatewayDeviceInfo, hwid *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) error {
			return nil
		},
		GetControlProxy: func(networkID string) (string, error) {
			return controlProxy, nil
		},
	}

	stateTestInit.StartTestService(t)
	configuratorTestInit.StartTestService(t)

	createUnregisteredGateway(t)

	return registrationService
}

// createUnregisteredGateway creates an unregistered gateway, i.e. a gateway without its device field
func createUnregisteredGateway(t *testing.T) {
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: networkID}, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(context.Background(), networkID, []configurator.NetworkEntity{
		{
			Type:   orc8r.MagmadGatewayType,
			Key:    logicalID,
			Config: &models.MagmadGatewayConfigs{},
		},
	}, serdes.Entity)
	assert.NoError(t, err)
}

func checkRegisteredGateway(t *testing.T) {
	ent, err := configurator.LoadEntity(
		context.Background(),
		networkID, orc8r.MagmadGatewayType, logicalID,
		configurator.EntityLoadCriteria{},
		serdes.Entity,
	)
	assert.Equal(t, ent.PhysicalID, hardwareID)
	assert.NoError(t, err)
}

func setupAddNetworksToTenantsService(t *testing.T) {
	var (
		tenant1 = &tenant_protos.Tenant{
			Name:     "tenant",
			Networks: []string{"network1", "network2"},
		}
		tenant2 = &tenant_protos.Tenant{
			Name:     "tenant",
			Networks: []string{"network3", "network4"},
		}
	)
	tenant_TestInit.StartTestService(t)

	addTenant(t, tenant1)
	addTenant(t, tenant2)
}

func addTenant(t *testing.T, tenant *tenant_protos.Tenant) int64 {
	ctx := context.Background()

	tenantRes, err := tenants.CreateTenant(ctx, nextTenantID, tenant)
	assert.NoError(t, err)
	assert.Equal(t, tenant, tenantRes)

	curTenantID := nextTenantID
	nextTenantID = nextTenantID + 1
	return curTenantID
}
