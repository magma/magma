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

package servicers_test

import (
	"net"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/multiplex"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/swx_proxy/servicers/test"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// How multiple_swx_proxy_test works:
// In order to test MultiSwxProxy we will create N diam server. Only one of those servers will
// respond properly (the others are created with StartEmptyDiameterServer). If multiplex feature
// works properly MockMultiplexor will send the request to the active server. Otherwise the
// request will time out.
//
// To extend the test cases you can follow TestMultiSwxProxyService_VerifyAuthorization. The loop
// in that tests just tests all the possible positions of the active server. activeServerIndex in
// each iterations represent the server that should be responding

// TestMultiSwxProxyService_VerifyAuthorization does the same as TestSwxProxyService_VerifyAuthorizationtests
// but in this case it uses multiple server
func TestMultiSwxProxyService_VerifyAuthorization(t *testing.T) {
	configs := getSwxTestConfigs(true)
	for activeServerIndex := range configs {
		t.Logf("Starting tests against SWx at at position %d", activeServerIndex)
		addr := initMultiSwxTestSetup(t, configs, activeServerIndex)
		// Set up a connection to the server.
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("GRPC connect error: %v", err)
			return
		}
		defer conn.Close()
		client := protos.NewSwxProxyClient(conn)
		swxStandardTest(t, client, 5)
	}
}

// initMultiSwxTestSetup starts one valid HSS and n-1 empty diameter servers
func initMultiSwxTestSetup(t *testing.T, configs []*servicers.SwxProxyConfig, activeServerIndex int) string {
	// ---- CORE 3gpp ----
	// create the mockHSS server/servers (depending on the config)
	for i, config := range configs {
		var err error
		var serverAddr string
		if activeServerIndex == i {
			serverAddr, err = test.StartTestSwxServer(TCPorSCTP, "127.0.0.1:0")
			t.Logf("Started Swx Server at %s", serverAddr)
		} else {
			serverAddr, err = test.StartEmptyDiameterServer(TCPorSCTP, "127.0.0.1:0")
			t.Logf("Started Diam Server at %s", serverAddr)
		}
		if err != nil {
			t.Fatal(err)
		}
		// Update config address with address of where test swx server is running
		config.ServerCfg.Addr = serverAddr
	}
	// create mux. Here is were we tell the multiplexer where to send the request for this setup
	mux := &MockMultiplexor{t: t, fixedServer: activeServerIndex}

	// create a SWx Service with default mockMultiplexer
	swxService, err := servicers.NewSwxProxies(configs, mux)
	if err != nil {
		t.Fatalf("failed to create SwxProxy: %v", err)
	}
	// ---- GRPC ----
	// create the gRPC server
	grpcListener, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	protos.RegisterSwxProxyServer(grpcServer, swxService)

	// start gRPC service
	go func() {
		if err2 := grpcServer.Serve(grpcListener); err2 != nil {
			t.Errorf("failed to serve: %v", err2)
		}
	}()
	addr := grpcListener.Addr()
	t.Logf("Started Swx GRPC Proxy on %s", addr.String())

	// returns the address of the client at SWx side.
	return addr.String()
}

// Produces a slice of two HSS configurations
func getSwxTestConfigs(verify bool) []*servicers.SwxProxyConfig {
	return []*servicers.SwxProxyConfig{
		{
			ClientCfg: &diameter.DiameterClientConfig{
				Host:  "magma-oai.openair4G.eur", // diameter host
				Realm: "openair4G.eur",           // diameter realm,
			},
			ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:     "",         // to be filled in once server addr is started
				Protocol: TCPorSCTP}, // tcp/sctp
			},
			VerifyAuthorization: verify,
		},
		{
			ClientCfg: &diameter.DiameterClientConfig{
				Host:  "magma2-oai.openair4G.eur", // diameter host
				Realm: "openair4G.eur",            // diameter realm,
			},
			ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:     "",         // to be filled in once server addr is started
				Protocol: TCPorSCTP}, // tcp/sctp
			},
			VerifyAuthorization: verify,
		},
	}
}

//  ---- MockMultiplexor ----
type MockMultiplexor struct {
	mock.Mock
	t           *testing.T
	fixedServer int
}

func (mp *MockMultiplexor) GetIndex(muxCtx *multiplex.Context) (int, error) {
	imsi, err := muxCtx.GetIMSI()
	if err != nil {
		mp.t.Fatal(err)
	}
	mp.t.Logf("Multiplexor selector sent %d to proxy #%d", imsi, mp.fixedServer)
	return mp.fixedServer, nil
}
