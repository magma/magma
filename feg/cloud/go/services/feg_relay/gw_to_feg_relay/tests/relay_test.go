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

package tests

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	models2 "magma/feg/cloud/go/services/feg/obsidian/models"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"magma/feg/cloud/go/feg"
	feg_protos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay"
	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay/servicers"
	healthTestUtils "magma/feg/cloud/go/services/health/test_utils"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	service_test_utils "magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const (
	s6aProxyService  = "s6a_proxy"
	testHelloService = "feg_hello"
)

type testS6aProxy struct {
	feg_protos.UnimplementedS6AProxyServer
	resultChan chan string // Calling FeG ID string on success
}

func (tp *testS6aProxy) AuthenticationInformation(
	ctx context.Context,
	req *feg_protos.AuthenticationInformationRequest) (*feg_protos.AuthenticationInformationAnswer, error) {

	if tp == nil {
		return nil, fmt.Errorf("nil test S6a proxy")
	}
	if tp.resultChan == nil {
		return nil, fmt.Errorf("nil test S6a proxy resultChan")
	}
	var targetFegId = "<MISSING METADATA>"
	ctxMetadata, ok := metadata.FromIncomingContext(ctx)
	if ok && ctxMetadata != nil {
		targetFegId = "<MISSING GW ID>"
		values, ok := ctxMetadata[gateway_registry.GatewayIdHeaderKey]
		if !ok {
			values, ok = ctxMetadata[strings.ToLower(gateway_registry.GatewayIdHeaderKey)]
		}
		if ok && len(values) > 0 {
			targetFegId = values[0]
		}
	}
	tp.resultChan <- targetFegId
	return &feg_protos.AuthenticationInformationAnswer{}, nil
}

type testHelloServer struct {
}

func (tp *testHelloServer) SayHello(c context.Context, req *feg_protos.HelloRequest) (*feg_protos.HelloReply, error) {
	return &feg_protos.HelloReply{Greeting: "testHelloService reply to: " + req.GetGreeting()}, nil
}

func TestNHRouting(t *testing.T) {
	testHealthServiser := setupNeutralHostNetworks(t)

	// test # 1: Verify, relay finds the right serving FeG for IMSI's PLMN ID
	foundFegHwId, err := gw_to_feg_relay.FindServingFeGHwId(federatedLteNetworkID, nhImsi)
	assert.NoError(t, err)
	assert.Equal(t, fegHwId, foundFegHwId)

	// test #2: Verify routing of matched user PLMNID
	//
	// Start & register Serving FeG's test S6a Proxy Server
	s6aProxy := &testS6aProxy{resultChan: make(chan string, 3)}
	srv, lis := test_utils.NewTestService(t, "feg", s6aProxyService)
	s6aAddr := lis.Addr().(*net.TCPAddr)
	s6aHost := "localhost"
	t.Logf("Serving FeG S6a Proxy Address: %s", s6aAddr)

	feg_protos.RegisterS6AProxyServer(srv.GrpcServer, s6aProxy)

	// Register FeG's test Hello Server
	feg_protos.RegisterHelloServer(srv.GrpcServer, &testHelloServer{})

	go srv.RunTest(lis)

	// Add Serving FeG Host to directoryd
	directoryd_test_init.StartTestService(t)
	directoryd.MapHWIDToHostname(fegHwId, s6aHost)
	gateway_registry.SetPort(s6aAddr.Port)

	// Start S6a relay Service
	relaySrv, relayLis := test_utils.NewTestService(t, "", s6aProxyService)

	t.Logf("Relay S6a Proxy Address: %s", relayLis.Addr())

	relayRouter := servicers.NewRelayRouter()
	feg_protos.RegisterS6AProxyServer(relaySrv.GrpcServer, relayRouter)
	feg_protos.RegisterHelloServer(relaySrv.GrpcServer, relayRouter)
	go relaySrv.RunTest(relayLis)

	ctx := service_test_utils.GetContextWithCertificate(t, agwHwId)
	connectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := registry.GetClientConnection(connectCtx, relayLis.Addr().String())
	cancel()

	assert.NoError(t, err)
	s6aClient := feg_protos.NewS6AProxyClient(conn)
	aiReq := &feg_protos.AuthenticationInformationRequest{UserName: nhImsi}

	toutctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	_, err = s6aClient.AuthenticationInformation(toutctx, aiReq)
	cancel()
	assert.NoError(t, err)
	select {
	case servingFegHwId := <-s6aProxy.resultChan:
		assert.Equal(t, fegHwId, servingFegHwId)
	case <-time.After(3 * time.Second):
		t.Fatal("Neutral Host Routed S6a Proxy Call timed out")
	}

	// Test SayHello routing & NH argumentation regex
	helloClient := feg_protos.NewHelloClient(conn)
	toutctx, cancel = context.WithTimeout(ctx, 5*time.Second)

	// regex style: @nh-Feg-for Imsi 123456
	helloReq := &feg_protos.HelloRequest{Greeting: "Hello FeG @nh-Feg-for Imsi " + nhImsi}
	helloResp, err := helloClient.SayHello(toutctx, helloReq)
	assert.NoError(t, err)
	assert.Equal(t, "testHelloService reply to: Hello FeG", helloResp.GetGreeting())

	// regex style: @NH-FEG-FOR: IMSI123456
	helloReq = &feg_protos.HelloRequest{Greeting: "Hello FeG test 2 @NH-FEG-FOR: IMSI" + nhImsi}
	helloResp, err = helloClient.SayHello(toutctx, helloReq)
	assert.NoError(t, err)
	assert.Equal(t, "testHelloService reply to: Hello FeG test 2", helloResp.GetGreeting())

	// regex style: @NH-FeG-FOR IMSI123456
	helloReq = &feg_protos.HelloRequest{Greeting: "Hello FeG t3 @NH-FEG-FOR IMSI" + nhImsi}
	helloResp, err = helloClient.SayHello(toutctx, helloReq)
	assert.NoError(t, err)
	assert.Equal(t, "testHelloService reply to: Hello FeG t3", helloResp.GetGreeting())

	// expect failure on absent default local FeG
	helloReq = &feg_protos.HelloRequest{Greeting: "Hi" + nhImsi}
	_, err = helloClient.SayHello(toutctx, helloReq)
	assert.Error(t, err)

	cancel()

	// test #3: Verify failure of routing of unknown PLMN IDs
	aiReq.UserName = nonNhImsi
	toutctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	// The call for non-matching PLMN ID should end up on NH Network's FeG, but none is configured, expect an error
	_, err = s6aClient.AuthenticationInformation(toutctx, aiReq)
	cancel()
	assert.Error(t, err)

	// test #4: Verify serving of non-matching PLMN IDs by default NH FeG (if exist)
	//
	// Add a FeG to serve non-matching PLMN IDs to NH network
	_, err = configurator.CreateEntities(
		nhNetworkID,
		[]configurator.NetworkEntity{
			{
				Type: feg.FegGatewayType, Key: nhFegId,
			},
			{
				Type: orc8r.MagmadGatewayType, Key: nhFegId,
				Name: "nh_feg_gateway", Description: "neutral host federation gateway",
				PhysicalID:   nhFegHwId,
				Config:       &models.MagmadGatewayConfigs{},
				Associations: []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: nhFegId}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: nhFegId},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		nhNetworkID, orc8r.AccessGatewayRecordType, nhFegHwId,
		&models.GatewayDevice{HardwareID: nhFegHwId, Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)

	// Map NH FeG to already running test S6a proxy address
	directoryd.MapHWIDToHostname(nhFegHwId, "localhost")

	// Update Serving FeG Health status
	healthctx := protos.NewGatewayIdentity(nhFegHwId, nhNetworkID, nhFegId).NewContextWithIdentity(context.Background())
	req := healthTestUtils.GetHealthyRequest()
	_, err = testHealthServiser.UpdateHealth(healthctx, req)
	assert.NoError(t, err)

	// Verify that NH FeG will be used as "catch all" for all but nhPlmnId
	foundFegHwId, err = gw_to_feg_relay.FindServingFeGHwId(federatedLteNetworkID, "") // no IMSI
	assert.NoError(t, err)
	assert.Equal(t, nhFegHwId, foundFegHwId)

	foundFegHwId, err = gw_to_feg_relay.FindServingFeGHwId(federatedLteNetworkID, nhImsi) // NH IMSI
	assert.NoError(t, err)
	assert.Equal(t, fegHwId, foundFegHwId)

	toutctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	// Now the call for non-matching PLMN ID should end up on NH Network's FeG
	_, err = s6aClient.AuthenticationInformation(toutctx, aiReq)
	cancel()
	assert.NoError(t, err)
	select {
	case servingFegHwId := <-s6aProxy.resultChan:
		assert.Equal(t, nhFegHwId, servingFegHwId)
	case <-time.After(3 * time.Second):
		t.Fatal("Neutral Host Non Routed S6a Proxy Call timed out")
	}
	// Verify that matching PLMN ID routing still works
	aiReq.UserName = nhImsi
	toutctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	_, err = s6aClient.AuthenticationInformation(toutctx, aiReq)
	cancel()
	assert.NoError(t, err)
	select {
	case servingFegHwId := <-s6aProxy.resultChan:
		assert.Equal(t, fegHwId, servingFegHwId)
	case <-time.After(3 * time.Second):
		t.Fatal("Neutral Host Routed S6a Proxy Call timed out")
	}

	// test #5: Remove Neutral Host settings (making NH FeG network a legacy FeG Network) and verify that
	//			legacy (non NH) relay logic works as expected (GW requests with NH IMSI are routed to NH FeG)
	nhNet, err := configurator.LoadNetwork(nhNetworkID, true, true, serdes.Network)
	assert.NoError(t, err)
	assert.NotNil(t, nhNet)
	cfg, ok := nhNet.Configs[feg.FegNetworkType]
	assert.True(t, ok)
	assert.NotNil(t, cfg)
	fegCfg, ok := cfg.(*models2.NetworkFederationConfigs)
	assert.True(t, ok)
	assert.NotNil(t, fegCfg)
	fegCfg.NhRoutes = nil // delete NH configuration, now FeG Network is just a regular FeG Network
	err = configurator.UpdateNetworkConfig(nhNetworkID, feg.FegNetworkType, fegCfg, serdes.Network)
	assert.NoError(t, err)

	// Verify, relay now finds the NH local FeG for any IMSI
	foundFegHwId, err = gw_to_feg_relay.FindServingFeGHwId(federatedLteNetworkID, nhImsi)
	assert.NoError(t, err)
	assert.Equal(t, nhFegHwId, foundFegHwId)
	foundFegHwId, err = gw_to_feg_relay.FindServingFeGHwId(federatedLteNetworkID, nonNhImsi)
	assert.NoError(t, err)
	assert.Equal(t, nhFegHwId, foundFegHwId)
	foundFegHwId, err = gw_to_feg_relay.FindServingFeGHwId(federatedLteNetworkID, "") // no IMSI
	assert.NoError(t, err)
	assert.Equal(t, nhFegHwId, foundFegHwId)

	aiReq.UserName = nhImsi
	toutctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	// Now the call for non-matching PLMN ID should end up on NH Network's FeG
	_, err = s6aClient.AuthenticationInformation(toutctx, aiReq)
	cancel()
	assert.NoError(t, err)
	select {
	case servingFegHwId := <-s6aProxy.resultChan:
		assert.Equal(t, nhFegHwId, servingFegHwId)
	case <-time.After(3 * time.Second):
		t.Fatal("Non Neutral Host & NH IMSI S6a Proxy Call timed out")
	}

	aiReq.UserName = nonNhImsi
	toutctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	// Now the call for non-matching PLMN ID should end up on NH Network's FeG
	_, err = s6aClient.AuthenticationInformation(toutctx, aiReq)
	cancel()
	assert.NoError(t, err)
	select {
	case servingFegHwId := <-s6aProxy.resultChan:
		assert.Equal(t, nhFegHwId, servingFegHwId)
	case <-time.After(3 * time.Second):
		t.Fatal("Non Neutral Host, Non NH IMSI S6a Proxy Call timed out")
	}
}
