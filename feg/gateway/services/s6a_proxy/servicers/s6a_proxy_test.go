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
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/s6a_proxy/servicers/test"
	orcprotos "magma/orc8r/lib/go/protos"
)

const TEST_LOOPS = 33

var TCPorSCTP = systemBasedTCPorSCTP() // sctp if run in linux, tcp if run in MAC

// systemBasedTCPorSCTP decides to run the test in TCP or SCTP. By default tests should
// be run in SCTP, but if test are run on MacOs, TCP is the only supported protocol
func systemBasedTCPorSCTP() string {
	if runtime.GOOS == "darwin" {
		fmt.Println(
			"Running servers with TCP. MacOS detected, SCTP not supported in this system. " +
				"Use this mode only for debugging!!!")
		return "tcp"
	}
	fmt.Println("Running servers with SCTP")
	return "sctp"
}

// TestS6aProxyService creates a mock S6a Diameter server, S6a S6a Proxy service
// and runs tests using GRPC client: GRPC Client <--> GRPC Server <--> S6a SCTP Diameter Server
func TestS6aProxyService(t *testing.T) {

	config := generateS6aProxyConfig()

	addr := startTestServer(t, config, false)

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()

	c := protos.NewS6AProxyClient(conn)
	hs := protos.NewServiceHealthClient(conn)
	req := &protos.AuthenticationInformationRequest{
		UserName:                   test.TEST_IMSI,
		VisitedPlmn:                []byte(test.TEST_PLMN_ID),
		NumRequestedEutranVectors:  3,
		ImmediateResponsePreferred: true,
	}
	complChan := make(chan error, TEST_LOOPS+1)
	testLoopF := func() {
		// AIR
		r, airErr := c.AuthenticationInformation(context.Background(), req)
		if airErr != nil {
			t.Errorf("GRPC AIR Error: %v", airErr)
			complChan <- airErr
			return
		}
		t.Logf("GRPC AIA: %#+v", *r)
		if r.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected AIA Error Code: %d", r.ErrorCode)
		}
		if len(r.EutranVectors) != 3 {
			t.Errorf("Unexpected Number of EutranVectors: %d, Expected: 3", len(r.EutranVectors))
		}
		ulReq := &protos.UpdateLocationRequest{
			UserName:                     test.TEST_IMSI,
			VisitedPlmn:                  []byte(test.TEST_PLMN_ID),
			SkipSubscriberData:           false,
			InitialAttach:                true,
			DualRegistration_5GIndicator: true,
			FeatureListId_1: &protos.FeatureListId1{
				RegionalSubscription: true,
			},
			FeatureListId_2: &protos.FeatureListId2{
				NrAsSecondaryRat: true,
			},
		}
		// ULR
		ulResp, airErr := c.UpdateLocation(context.Background(), ulReq)
		if airErr != nil {
			t.Errorf("GRPC ULR Error: %v", airErr)
			complChan <- airErr
			return
		}
		t.Logf("GRPC ULA: %#+v", *ulResp)
		if ulResp.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected ULA Error Code: %d", ulResp.ErrorCode)
		}
		assert.NoError(t, airErr)
		if len(ulResp.RegionalSubscriptionZoneCode) != 2 ||
			!bytes.Equal(ulResp.RegionalSubscriptionZoneCode[0], []byte{155, 36, 12, 2, 227, 43, 246, 254}) ||
			!bytes.Equal(ulResp.RegionalSubscriptionZoneCode[1], []byte{1, 1, 0, 1}) {
			t.Errorf("There should be 2 Regional Subscription Zone Codes : %+v", ulResp.RegionalSubscriptionZoneCode)
		}
		assert.NotEmpty(t, ulResp.FeatureListId_1)
		assert.True(t, ulResp.FeatureListId_1.RegionalSubscription)
		assert.NotEmpty(t, ulResp.FeatureListId_2)
		assert.True(t, ulResp.FeatureListId_2.NrAsSecondaryRat)

		assert.NotNil(t, ulResp.TotalAmbr)
		assert.Equal(t, uint32(500), ulResp.TotalAmbr.MaxBandwidthDl)
		assert.Equal(t, uint32(600), ulResp.TotalAmbr.MaxBandwidthUl)
		assert.Equal(t, protos.UpdateLocationAnswer_AggregatedMaximumBitrate_KBPS, ulResp.TotalAmbr.Unit)
		assert.NotEmpty(t, ulResp.Apn)
		assert.Equal(t, uint32(50), ulResp.Apn[0].Ambr.MaxBandwidthDl)
		assert.Equal(t, uint32(60), ulResp.Apn[0].Ambr.MaxBandwidthUl)
		assert.Equal(t, protos.UpdateLocationAnswer_AggregatedMaximumBitrate_BPS, ulResp.Apn[0].Ambr.Unit)
		puReq := &protos.PurgeUERequest{
			UserName: test.TEST_IMSI,
		}
		// PUR
		puResp, airErr := c.PurgeUE(context.Background(), puReq)
		if airErr != nil {
			t.Errorf("GRPC PUR Error: %v", airErr)
			complChan <- airErr
			return
		}
		t.Logf("GRPC PUA: %#+v", *puResp)
		if puResp.ErrorCode != protos.ErrorCode_SUCCESS {
			t.Errorf("Unexpected PUA Error Code: %d", puResp.ErrorCode)
		}
		// End
		complChan <- nil
	}
	go testLoopF()
	select {
	case err := <-complChan:
		if err != nil {
			t.Fatal(err)
			return
		}
	case <-time.After(time.Second):
		t.Fatal("Timed out")
		return
	}

	for round := 0; round < TEST_LOOPS; round++ {
		go testLoopF()
	}
	for round := 0; round < TEST_LOOPS; round++ {
		testErr := <-complChan
		if testErr != nil {
			t.Fatal(err)
			return
		}
	}

	// Test Disabling / Enabling Connections

	disableReq := &protos.DisableMessage{
		DisablePeriodSecs: 10,
	}

	// Disable connections
	_, err = hs.Disable(context.Background(), disableReq)
	if err != nil {
		t.Fatalf("GRPC Disable Error: %v", err)
		return
	}

	// AIR should fail
	_, err = c.AuthenticationInformation(context.Background(), req)
	if err == nil {
		t.Errorf("AIR Succeeded, but should have failed due to disabled connections")
	}

	// Enable connections
	_, err = hs.Enable(context.Background(), &orcprotos.Void{})
	if err != nil {
		t.Fatalf("GRPC Enable Error: %v", err)
		return
	}

	// AIR should pass now
	airResp, err := c.AuthenticationInformation(context.Background(), req)
	if err != nil {
		t.Fatalf("GRPC AIR Error: %v", err)
		return
	}
	t.Logf("GRPC AIA: %#+v", *airResp)
	if airResp.ErrorCode != protos.ErrorCode_UNDEFINED {
		t.Errorf("Unexpected AIA Error Code: %d", airResp.ErrorCode)
	}
}

func TestS6aProxyServiceWitPLMNlist(t *testing.T) {
	config := generateS6aProxyConfigWithPLMNs()
	addr := startTestServer(t, config, false)

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()

	c := protos.NewS6AProxyClient(conn)

	complChan := make(chan error, 1)
	go func() {
		req := &protos.AuthenticationInformationRequest{
			UserName:                   test.TEST_IMSI,
			VisitedPlmn:                []byte(test.TEST_PLMN_ID),
			NumRequestedEutranVectors:  3,
			ImmediateResponsePreferred: true,
		}

		// AIR
		r, airErr := c.AuthenticationInformation(context.Background(), req)
		if airErr != nil {
			t.Errorf("GRPC AIR with PLMN IMSI1 Error: %v", airErr)
			complChan <- airErr
			return
		}
		t.Logf("GRPC AIA: %#+v", *r)
		if r.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected AIA with PLMN IMSI1 Error Code: %d", r.ErrorCode)
		}
		if len(r.EutranVectors) != 3 {
			t.Errorf("Unexpected Number of EutranVectors with PLMN IMSI1: %d, Expected: 3", len(r.EutranVectors))
		}

		// Use an IMSI that is not on the PLMN list
		req.UserName = test.TEST_IMSI_2
		r, airErr = c.AuthenticationInformation(context.Background(), req)
		if airErr != nil {
			t.Errorf("GRPC AIR with PLMN IMSI2 Error: %v", airErr)
			complChan <- airErr
			return
		}
		t.Logf("GRPC AIA: %#+v", *r)
		if r.ErrorCode != protos.ErrorCode_AUTHENTICATION_REJECTED {
			t.Errorf("Authentication Rejected was expected but AIA with PLMN IMSI2 got Error Code: %d", r.ErrorCode)
		}

		// End
		complChan <- nil
	}()

	select {
	case err := <-complChan:
		if err != nil {
			t.Fatal(err)
			return
		}
	case <-time.After(time.Second):
		t.Fatal("Timed out")
		return
	}
}

func TestS6aProxyWithHSS_AIA(t *testing.T) {
	config := generateS6aProxyConfig()
	addr := startTestServer(t, config, true)
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()

	c := protos.NewS6AProxyClient(conn)
	req := &protos.AuthenticationInformationRequest{
		UserName:                   test.TEST_IMSI,
		VisitedPlmn:                []byte(test.TEST_PLMN_ID),
		NumRequestedEutranVectors:  3,
		ImmediateResponsePreferred: true,
	}
	complChan := make(chan error, 1)
	go func() {
		t.Logf("TestS6aProxyWithHSS_AIA - AIA RPC Req: %s", req.String())
		r, err := c.AuthenticationInformation(context.Background(), req)
		if err != nil {
			t.Errorf("TestS6aProxyWithHSS_AIA - GRPC AIR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC AIA Resp: %#+v", *r)
		if r.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected AIA with PLMN IMSI1 Error Code: %d", r.ErrorCode)
		}
		assert.Len(t, r.EutranVectors, 3)
		assert.Equal(t,
			[]byte("\x15\x9a\xbf\x21\xca\xe2\xbf\x0a\xdb\xcb\xf1\x47\xef\x87\x74\x9d"),
			r.EutranVectors[0].Rand)
		assert.Equal(t,
			[]byte("\x63\x82\xb8\x54\x48\x59\x80\x00\xf5\xaf\x37\xa5\xe9\x6d\x76\x58"),
			r.EutranVectors[1].Autn)
		assert.Equal(t,
			[]byte("\x74\x60\x79\x2b\x8d\x5e\xb1\x62\xfd\x88\x28\xc2\x1a\x3b\xa0\xc5"+
				"\x6e\x06\xed\xbf\x5b\x20\x54\x72\x50\x06\x36\xc5\xfa\xd9\x0b\x84"),
			r.EutranVectors[2].Kasme)
		// End
		complChan <- nil
	}()

	select {
	case err := <-complChan:
		if err != nil {
			t.Fatal(err)
			return
		}
	case <-time.After(time.Second * 2):
		t.Fatal("TestS6aProxyWithHSS_AIA Timed out")
		return
	}
}

func generateS6aProxyConfig() *servicers.S6aProxyConfig {

	diamAddr := fmt.Sprintf("127.0.0.1:%d", 29000+rand.Intn(1900))

	return &servicers.S6aProxyConfig{
		ClientCfg: &diameter.DiameterClientConfig{
			Host:  "magma-oai.openair4G.eur", // diameter host
			Realm: "openair4G.eur",           // diameter realm,
		},
		ServerCfg: &diameter.DiameterServerConfig{
			DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:     diamAddr,   // use "192.168.60.145:3868" to send diam messages to OAI HSS VM
				Protocol: TCPorSCTP}, // tcp/sctp
		},
		PlmnIds: plmn_filter.PlmnIdVals{},
	}
}

func generateS6aProxyConfigWithPLMNs() *servicers.S6aProxyConfig {

	diamAddr := fmt.Sprintf("127.0.0.1:%d", 29000+rand.Intn(1900))

	return &servicers.S6aProxyConfig{
		ClientCfg: &diameter.DiameterClientConfig{
			Host:  "magma-oai.openair4G.eur", // diameter host
			Realm: "openair4G.eur",           // diameter realm,
		},
		ServerCfg: &diameter.DiameterServerConfig{
			DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:     diamAddr,   // use "192.168.60.145:3868" to send diam messages to OAI HSS VM
				Protocol: TCPorSCTP}, // tcp/sctp
		},
		PlmnIds: plmn_filter.GetPlmnVals([]string{"00101", "00102"}),
	}
}

func startTestServer(t *testing.T, config *servicers.S6aProxyConfig, useStaticResp bool) string {
	// ---- CORE 3gpp ----
	// create the mockHSS server/servers (depending on the config)
	err := test.StartTestS6aServer(TCPorSCTP, config.ServerCfg.Addr, useStaticResp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Started S6a Server")

	// ---- GRPC ----
	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service, err := servicers.NewS6aProxy(config)
	if err != nil {
		t.Fatalf("failed to create S6aProxy: %v", err)
	}
	protos.RegisterS6AProxyServer(s, service)
	protos.RegisterServiceHealthServer(s, service)
	go func() {
		if errSrv := s.Serve(lis); errSrv != nil {
			t.Errorf("test server failed to serve: %v", errSrv)
			return
		}
	}()
	addr := lis.Addr().String()
	t.Logf("Started S6a GRPC Proxy on %s", addr)
	return addr
}
