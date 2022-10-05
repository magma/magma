/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package unary_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/service/middleware/unary"
	"magma/orc8r/cloud/go/service/middleware/unary/test"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	configuratorTestUtils "magma/orc8r/cloud/go/services/configurator/test_utils"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const testAgHwID = "Test-AGW-Hw-Id"

type testStateServer struct {
	lastClientIdentity    *protos.Identity
	lastClientCertExpTime int64
}

func NewTestStateServer() (*testStateServer, error) {
	return &testStateServer{}, nil
}

func (srv *testStateServer) GetStates(ctx context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	srv.lastClientIdentity =
		proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	return nil, nil
}

func (srv *testStateServer) ReportStates(ctx context.Context, req *protos.ReportStatesRequest) (*protos.ReportStatesResponse, error) {
	srv.lastClientIdentity = proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	srv.lastClientCertExpTime = protos.GetClientCertExpiration(ctx)
	return &protos.ReportStatesResponse{UnreportedStates: []*protos.IDAndError{}}, nil
}

func (srv *testStateServer) DeleteStates(ctx context.Context, req *protos.DeleteStatesRequest) (*protos.Void, error) {
	srv.lastClientIdentity =
		proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	return &protos.Void{}, nil
}

func (srv *testStateServer) SyncStates(ctx context.Context, req *protos.SyncStatesRequest) (*protos.SyncStatesResponse, error) {
	srv.lastClientIdentity = proto.Clone(protos.GetClientIdentity(ctx)).(*protos.Identity)
	srv.lastClientCertExpTime = protos.GetClientCertExpiration(ctx)
	return &protos.SyncStatesResponse{UnsyncedStates: []*protos.IDAndVersion{}}, nil
}

func TestIdentityInjector(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)

	networkID := "identity_decorator_test_network"
	configuratorTestUtils.RegisterNetwork(t, networkID, "Identity Decorator Test")
	configuratorTestUtils.RegisterGateway(t, networkID, testAgHwID, &models.GatewayDevice{HardwareID: testAgHwID})

	// Create the service
	srv, err := service.NewTestOrchestratorService(t, orc8r.ModuleName, state.ServiceName)
	assert.NoError(t, err)

	// Add servicers to the service
	stateServer, err := NewTestStateServer()
	assert.NoError(t, err)
	protos.RegisterStateServiceServer(srv.GrpcServer, stateServer)

	l, err := net.Listen("tcp", "")
	assert.NoError(t, err)
	addr := l.Addr().String()
	// Run the service
	go srv.RunTest(l, nil)

	conn, err := registry.GetClientConnection(context.Background(), addr)
	assert.NoError(t, err)
	stateClient := protos.NewStateServiceClient(conn)
	csn := test.StartMockGwAccessControl(t, []string{testAgHwID})
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", csn[0]))

	gwState := &models.GatewayStatus{
		Meta: map[string]string{
			"foo": "bar",
		},
	}
	serializedGWStatus, err := serde.Serialize(gwState, orc8r.GatewayStateType, serdes.State)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     orc8r.GatewayStateType,
			DeviceID: testAgHwID,
			Value:    serializedGWStatus,
		},
	}
	_, err = stateClient.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
	identity := stateServer.lastClientIdentity
	assert.NotNil(t, identity)
	assert.True(t, time.Now().Unix() < stateServer.lastClientCertExpTime)

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
	_, err = stateClient.ReportStates(
		context.Background(),
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
	identity = stateServer.lastClientIdentity
	assert.Nil(t, identity)
	assert.Equal(t, int64(0), stateServer.lastClientCertExpTime)

	// Test empty x-magma-client-cert-serial header
	// Hack in the identity context
	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", ""))
	_, err = stateClient.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.Error(t, err)
	assert.Equal(t, int64(0), stateServer.lastClientCertExpTime)

	// Test x-magma-client-cert-cn, but not x-magma-client-cert-serial headers
	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-cn", "bla bla bla"))
	_, err = stateClient.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.Error(t, err)
	assert.Equal(t, int64(0), stateServer.lastClientCertExpTime)

	// Unregister GW
	assert.NoError(
		t,
		configurator.DeleteEntity(context.Background(), networkID, orc8r.MagmadGatewayType, gwid.LogicalId))

	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", csn[0]))

	// Expect PermissionDenied error now
	_, err = stateClient.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.Error(t, err)
	assert.Equal(t, int64(0), stateServer.lastClientCertExpTime)
	assert.Equal(
		t,
		"rpc error: code = PermissionDenied desc = Unregistered Gateway Test-AGW-Hw-Id",
		err.Error())
}

type testAddr string

func (a testAddr) String() string {
	return string(a)
}

func (a testAddr) Network() string {
	return string(a)
}

type testCase struct {
	ctx                context.Context
	serverInfo         *grpc.UnaryServerInfo
	expectedError      error
	expectedContextNil bool
	reachesAllowCheck  bool
}

func TestSetIdentityFromContext(t *testing.T) {
	csn := test.StartMockGwAccessControl(t, []string{testAgHwID})

	testCases := []testCase{
		{
			// Return an error when context metadata is missing.
			// This means the call is not authenticated. Return an error.
			ctx:                context.Background(),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      status.Error(codes.Unauthenticated, unary.ERROR_MSG_NO_METADATA),
			expectedContextNil: true,
			reachesAllowCheck:  false,
		}, {
			// Context metadata exists but empty. No CSN is given and the call
			// cannot be authenticated. Return error since this
			// is not a local client call.
			ctx:                metadata.NewIncomingContext(context.Background(), nil),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      status.Error(codes.PermissionDenied, unary.ERROR_MSG_UNKNOWN_CLIENT),
			expectedContextNil: true,
			reachesAllowCheck:  true,
		}, {
			// Context metadata exists but empty. No CSN is given and the call
			// cannot be authenticated. Return an error since this
			// is not a local client call. `serverInfo` is nil, make no
			// allowList check. Logs "Undefined" instead of
			// `serverInfo.FullMethod`.
			ctx:                metadata.NewIncomingContext(context.Background(), nil),
			serverInfo:         nil,
			expectedError:      status.Error(codes.PermissionDenied, unary.ERROR_MSG_UNKNOWN_CLIENT),
			expectedContextNil: true,
			reachesAllowCheck:  false,
		}, {
			// Context metadata exists but empty. No CSN is given and the call
			// cannot be authenticated. Return an error since this is a
			// local client call.
			ctx:                metadata.NewIncomingContext(peer.NewContext(context.Background(), &peer.Peer{Addr: testAddr("127.168.0.1:4567"), AuthInfo: nil}), nil),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      nil,
			expectedContextNil: true,
			reachesAllowCheck:  true,
		}, {
			// If the context contains the CN key but not the right value,
			// return an error.
			ctx:                metadata.NewIncomingContext(context.Background(), metadata.Pairs(unary.CLIENT_CERT_CN_KEY, "val")),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      status.Error(codes.Unauthenticated, "Inconsistent Request Signature"),
			expectedContextNil: true,
			reachesAllowCheck:  true,
		}, {
			// If the context contains keys different from CN and SN, return
			// an error
			ctx:                metadata.NewIncomingContext(context.Background(), metadata.Pairs("key", "val")),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      status.Error(codes.PermissionDenied, unary.ERROR_MSG_UNKNOWN_CLIENT),
			expectedContextNil: true,
			reachesAllowCheck:  true,
		}, {
			// If the SN key is present with the right value, return no error.
			ctx:                metadata.NewIncomingContext(context.Background(), metadata.Pairs(unary.CLIENT_CERT_SN_KEY, registry.ORC8R_CLIENT_CERT_VALUE)),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      nil,
			expectedContextNil: true,
			reachesAllowCheck:  false,
		}, {
			// If multiple SN keys are present, return an error.
			ctx:                metadata.NewIncomingContext(context.Background(), metadata.Pairs(unary.CLIENT_CERT_SN_KEY, registry.ORC8R_CLIENT_CERT_VALUE, unary.CLIENT_CERT_SN_KEY, "other value")),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      status.Error(codes.Unauthenticated, "Multiple CSNs present"),
			expectedContextNil: true,
			reachesAllowCheck:  true,
		}, {
			// If CN key is present with the value, return no error.
			ctx:                metadata.NewIncomingContext(context.Background(), metadata.Pairs(unary.CLIENT_CERT_SN_KEY, csn[0])),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      nil,
			expectedContextNil: false,
			reachesAllowCheck:  false,
		}, {
			// If CN key is present with the wrong value, return an error.
			ctx:                metadata.NewIncomingContext(context.Background(), metadata.Pairs(unary.CLIENT_CERT_SN_KEY, "wrong CSN")),
			serverInfo:         &grpc.UnaryServerInfo{},
			expectedError:      status.Error(codes.PermissionDenied, "Unknown Client Certificate"),
			expectedContextNil: true,
			reachesAllowCheck:  true,
		},
	}

	for _, testCase := range testCases {
		testCaseSetIdentityFromContext(t, testCase)
	}
}

func testCaseSetIdentityFromContext(t *testing.T, tc testCase) {
	newCtx, newReq, resp, err := unary.SetIdentityFromContext(tc.ctx, nil, tc.serverInfo)

	if tc.expectedContextNil {
		assert.Nil(t, newCtx)
	} else {
		// Don't explicitly check what a non-nil context might look like.
		// This depends on things such as system time.
		assert.NotNil(t, newCtx)
	}
	assert.Nil(t, newReq) // newReq is always nil
	assert.Nil(t, resp)   // resp is always nil
	if tc.expectedError == nil {
		assert.NoError(t, err)
	} else {
		assert.EqualError(t, err, tc.expectedError.Error())
	}

	if tc.reachesAllowCheck {
		// Check that the behavior is correct if the FullMethod is on the
		// Allow list.
		serverInfo := grpc.UnaryServerInfo{FullMethod: "/magma.orc8r.Bootstrapper/GetChallenge"}
		newCtx, newReq, resp, err = unary.SetIdentityFromContext(tc.ctx, nil, &serverInfo)
		assert.Nil(t, newCtx)
		assert.Nil(t, newReq)
		assert.Nil(t, resp)
		assert.NoError(t, err)
	}

}
