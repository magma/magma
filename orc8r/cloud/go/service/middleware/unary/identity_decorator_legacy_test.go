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
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/service/middleware/unary/test_utils"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type testCheckindServer struct {
	lastClientIdentity    *protos.Identity
	lastClientCertExpTime int64
}

func NewTestCheckindServer() (*testCheckindServer, error) {
	return &testCheckindServer{}, nil
}

// Gateway periodic checkin
func (srv *testCheckindServer) Checkin(
	ctx context.Context,
	req *protos.CheckinRequest) (*protos.CheckinResponse, error) {

	srv.lastClientIdentity = proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	srv.lastClientCertExpTime = protos.GetClientCertExpiration(ctx)
	return &protos.CheckinResponse{Action: protos.CheckinResponse_NONE,
			Time: uint64(time.Now().UnixNano()) / uint64(time.Millisecond)},
		nil
}

// Gateway real time status retrieval
func (srv *testCheckindServer) GetStatus(
	ctx context.Context,
	req *protos.GatewayStatusRequest) (*protos.GatewayStatus, error) {

	srv.lastClientIdentity =
		proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	return new(protos.GatewayStatus), nil
}

// Removes Gateway status record from the Gateway's network table
func (srv *testCheckindServer) DeleteGatewayStatus(
	ctx context.Context,
	req *protos.GatewayStatusRequest) (*protos.Void, error) {

	srv.lastClientIdentity =
		proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	return &protos.Void{}, nil
}

// Deletes the network's status table
func (srv *testCheckindServer) DeleteNetwork(
	ctx context.Context, networkId *protos.NetworkID) (*protos.Void, error) {

	srv.lastClientIdentity =
		proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	return &protos.Void{}, nil
}

// Returns a list of all logical gateway IDs
func (srv *testCheckindServer) List(
	ctx context.Context, networkId *protos.NetworkID) (*protos.IDList, error) {

	srv.lastClientIdentity =
		proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	return new(protos.IDList), nil
}

func TestIdentityInjectorLegacy(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "0")
	magmad_test_init.StartTestService(t)
	// Make sure to "share" in memory magmad DBs with interceptors

	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Identity Decorator Test"},
		"identity_decorator_test_network")
	assert.NoError(t, err)

	t.Logf("New Registered Network: %s", testNetworkId)

	hwId := protos.AccessGatewayID{Id: testAgHwID}
	logicalId, err := magmad.RegisterGateway(testNetworkId, &magmad_protos.AccessGatewayRecord{HwId: &hwId, Name: "Test GW Name"})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalId, "")

	// Create the service
	srv, err := service.NewTestOrchestratorService(t, "", checkind.ServiceName)
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
	assert.Equal(t, gwid.NetworkId, testNetworkId)
	assert.Equal(t, gwid.LogicalId, logicalId)

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
		magmad.RemoveGateway(testNetworkId, request.GatewayId))

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
