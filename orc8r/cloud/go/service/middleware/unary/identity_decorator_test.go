/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package unary_test

import (
	"net"
	"os"
	"testing"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/service/middleware/unary/test_utils"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/device"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const testAgHwID = "Test-AGW-Hw-Id"

func TestIdentityInjector(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	_ = serde.RegisterSerdes(serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}))

	// Make sure to "share" in memory magmad DBs with interceptors
	networkID := "identity_decorator_test_network"
	testNetwork := configurator.Network{
		ID:   networkID,
		Name: "Identity Decorator Test",
	}
	err := configurator.CreateNetwork(testNetwork)
	assert.NoError(t, err)

	configurator_test_utils.RegisterGateway(t, networkID, testAgHwID, &models.GatewayDevice{HardwareID: testAgHwID})

	// Create the service
	srv, err := service.NewTestOrchestratorService(t, orc8r.ModuleName, checkind.ServiceName)
	assert.NoError(t, err)

	// Add servicers to the service
	checkindServer, err := NewTestCheckindServer()
	assert.NoError(t, err)
	protos.RegisterCheckindServer(srv.GrpcServer, checkindServer)

	l, err := net.Listen("tcp", "")
	assert.NoError(t, err)
	addr := l.Addr().String()
	// Run the service
	go srv.RunTest(l)

	conn, err := registry.GetClientConnection(context.Background(), addr)
	assert.NoError(t, err)

	// Test GW updating status
	request := protos.CheckinRequest{
		GatewayId:       testAgHwID,
		MagmaPkgVersion: "1.2.3",
		Status: &protos.ServiceStatus{
			Meta: map[string]string{
				"hello": "world",
			},
		},
		SystemStatus: &protos.SystemStatus{
			CpuUser:   31498,
			CpuSystem: 8361,
			CpuIdle:   1869111,
			MemTotal:  1016084,
			MemUsed:   54416,
			MemFree:   412772,
		},
	}

	csn := test_utils.StartMockGwAccessControl(t, []string{testAgHwID})
	magmaCheckindClient := protos.NewCheckindClient(conn)

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", csn[0]))

	_, err = magmaCheckindClient.Checkin(ctx, &request)
	assert.NoError(t, err)
	identity := checkindServer.lastClientIdentity
	assert.NotNil(t, identity)
	assert.True(t, time.Now().Unix() < checkindServer.lastClientCertExpTime)

	cn := identity.ToCommonName()
	assert.NotNil(t, cn)
	assert.Equal(t, *cn, testAgHwID)

	gwid := identity.GetGateway()
	assert.NotNil(t, gwid)
	assert.Equal(t, gwid.HardwareId, testAgHwID)
	assert.Equal(t, gwid.NetworkId, networkID)
	assert.Equal(t, gwid.LogicalId, testAgHwID)

	// Test CTX without any Identification related headers (Identity should
	// not be injected by the middleware)
	_, err = magmaCheckindClient.Checkin(context.Background(), &request)
	assert.NoError(t, err)
	identity = checkindServer.lastClientIdentity
	assert.Nil(t, identity)
	assert.Equal(t, int64(0), checkindServer.lastClientCertExpTime)

	// Test empty x-magma-client-cert-serial header
	// Hack in the identity context
	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", ""))
	_, err = magmaCheckindClient.Checkin(ctx, &request)
	assert.Error(t, err)
	assert.Equal(t, int64(0), checkindServer.lastClientCertExpTime)

	// Test x-magma-client-cert-cn, but not x-magma-client-cert-serial headers
	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-cn", "bla bla bla"))
	_, err = magmaCheckindClient.Checkin(ctx, &request)
	assert.Error(t, err)
	assert.Equal(t, int64(0), checkindServer.lastClientCertExpTime)

	// Unregister GW
	assert.NoError(
		t,
		configurator.DeleteEntity(networkID, orc8r.MagmadGatewayType, gwid.LogicalId))

	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", csn[0]))

	// Expect PermissionDenied error now
	_, err = magmaCheckindClient.Checkin(ctx, &request)
	assert.Error(t, err)
	assert.Equal(t, int64(0), checkindServer.lastClientCertExpTime)
	assert.Equal(
		t,
		"rpc error: code = PermissionDenied desc = Unregistered Gateway Test-AGW-Hw-Id",
		err.Error())
}
