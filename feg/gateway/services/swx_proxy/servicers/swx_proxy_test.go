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
	"context"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/swx_proxy/servicers/test"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

// TestSwxProxyService_VerifyAuthorization tests swx_proxy by sending Authenticate
// gRPC request (with VerifyAuthorization set to true) and then sending Register gRPC
// request and ensuring the request succeeds
func TestSwxProxyService_VerifyAuthorization(t *testing.T) {
	config := getSwxTestConfig(true)
	addr := initSwxTestSetup(t, config)
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()
	client := protos.NewSwxProxyClient(conn)
	swxStandardTest(t, client, TEST_LOOPS)
}

// TestSwxProxyService_VerifyAuthorizationOff tests swx_proxy by sending Authenticate
// gRPC request (with VerifyAuthorization set to false) and then sending Register gRPC
// request and ensuring the request succeeds
func TestSwxProxyService_VerifyAuthorizationOff(t *testing.T) {
	config := getSwxTestConfig(false)
	addr := initSwxTestSetup(t, config)
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()
	client := protos.NewSwxProxyClient(conn)
	swxStandardTest(t, client, TEST_LOOPS)
}

// TestSwxProxyService_ValidationErrors tests the swx proxy service error handling
// by sending a variety of invalid gRPC requests
func TestSwxProxyService_ValidationErrors(t *testing.T) {
	config := getSwxTestConfig(true)
	addr := initSwxTestSetup(t, config)
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()
	client := protos.NewSwxProxyClient(conn)
	swxStandardTest(t, client, TEST_LOOPS)

	// Test Auth Error Handling
	_, err = client.Authenticate(context.Background(), nil)
	assert.EqualError(t, err, "rpc error: code = Internal desc = grpc: error while marshaling: proto: Marshal called with nil")

	emptyAuthReq := &protos.AuthenticationRequest{}
	_, err = client.Authenticate(context.Background(), emptyAuthReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = empty user-name provided in authentication request")

	badNumVectorsReq := &protos.AuthenticationRequest{
		UserName:             "10111011000110",
		AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		SipNumAuthVectors:    0,
	}

	_, err = client.Authenticate(context.Background(), badNumVectorsReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = SIPNumAuthVectors in authentication request must be greater than 0")

	badUserNameReq := &protos.AuthenticationRequest{
		UserName:             "1234567890123456",
		AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		SipNumAuthVectors:    0,
	}
	_, err = client.Authenticate(context.Background(), badUserNameReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = SIPNumAuthVectors in authentication request must be greater than 0")

	// Test Register Error Handling
	_, err = client.Register(context.Background(), nil)
	assert.EqualError(t, err, "rpc error: code = Internal desc = grpc: error while marshaling: proto: Marshal called with nil")

	emptyResReq := &protos.RegistrationRequest{}
	_, err = client.Register(context.Background(), emptyResReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Empty user-name provided in registration request")

	badRegReq := &protos.RegistrationRequest{
		UserName: "1234567890123456",
	}
	_, err = client.Register(context.Background(), badRegReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Provided username 1234567890123456 is greater than 15 digits")
}

func swxStandardTest(t *testing.T, client protos.SwxProxyClient, test_loops int) {
	complChan := make(chan error, test_loops+1)

	// Happy path
	testHappyPath := func(reqId uint32) {
		userName := test.BASE_IMSI + strconv.Itoa(int(reqId))
		authReq := &protos.AuthenticationRequest{
			UserName:             userName,
			SipNumAuthVectors:    1,
			AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		}
		// Authentication Request - MAR
		// proxy will always ask for MinRequestedVectors
		// and will return 1 vector
		for i := uint32(0); i < servicers.MinRequestedVectors; i++ {
			authRes, err := client.Authenticate(context.Background(), authReq)
			if err != nil {
				t.Errorf("GRPC MAR Error: %v", err)
				complChan <- err
				return
			}
			t.Logf("GRPC MAA: %#+v", *authRes)
			assert.Equal(t, userName, authRes.GetUserName())
			if len(authRes.SipAuthVectors) != 1 {
				t.Errorf("Unexpected Number of SIPAuthVectors: %d, Expected: %d", len(authRes.SipAuthVectors), 1)
			}
			v := authRes.SipAuthVectors[0]
			assert.Equal(t, protos.AuthenticationScheme_EAP_AKA, v.GetAuthenticationScheme())
			assert.Equal(t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(i+14))), v.GetRandAutn())
			assert.Equal(t, []byte(test.DefaultSIPAuthorization), v.GetXres())
			assert.Equal(t, []byte(test.DefaultCK), v.GetConfidentialityKey())
			assert.Equal(t, []byte(test.DefaultIK), v.GetIntegrityKey())
		}
		// Registration request - SAR
		regReq := &protos.RegistrationRequest{
			UserName: userName,
		}
		regRes, err := client.Register(context.Background(), regReq)
		// Only must verify that request was successful (no error) to ensure user
		// is registered
		if err != nil {
			t.Errorf("GRPC SAR Register Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC SAA Register: %#+v", *regRes)

		unregRes, err := client.Deregister(context.Background(), regReq)
		// Only must verify that request was successful (no error) to ensure user
		// is de-registered
		if err != nil {
			t.Errorf("GRPC SAR De-register Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC SAA De-register: %#+v", *unregRes)

		// Test EAP-SIM (3 vectors) case
		authReq = &protos.AuthenticationRequest{
			// make sure, IMSI for SIM tests does not overlap with single vector tests
			UserName:             "9" + userName + "9",
			SipNumAuthVectors:    3,
			AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		}
		authRes, err := client.Authenticate(context.Background(), authReq)
		if err != nil {
			t.Errorf("GRPC MAR Error: %v", err)
			complChan <- err
			return
		}
		assert.Equal(t, authReq.UserName, authRes.GetUserName())
		if len(authRes.SipAuthVectors) != 3 {
			t.Errorf("Unexpected Number of SIPAuthVectors: %d, Expected: %d", len(authRes.SipAuthVectors), 3)
		}
		for i := uint32(0); i < 3; i++ {
			v := authRes.SipAuthVectors[i]
			assert.Equal(t, protos.AuthenticationScheme_EAP_AKA, v.GetAuthenticationScheme())
			assert.Equal(t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(i+14))), v.GetRandAutn())
			assert.Equal(t, []byte(test.DefaultSIPAuthorization), v.GetXres())
			assert.Equal(t, []byte(test.DefaultCK), v.GetConfidentialityKey())
			assert.Equal(t, []byte(test.DefaultIK), v.GetIntegrityKey())
		}
		authRes, err = client.Authenticate(context.Background(), authReq)
		if err != nil {
			t.Errorf("GRPC MAR Error: %v", err)
			complChan <- err
			return
		}
		assert.Equal(t, authReq.UserName, authRes.GetUserName())
		if len(authRes.SipAuthVectors) != 3 {
			t.Errorf("Unexpected Number of SIPAuthVectors: %d, Expected: %d", len(authRes.SipAuthVectors), 3)
		}
		assert.Equal(
			t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(17))), authRes.SipAuthVectors[0].GetRandAutn())
		assert.Equal(
			t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(18))), authRes.SipAuthVectors[1].GetRandAutn())
		assert.Equal(
			t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(14))), authRes.SipAuthVectors[2].GetRandAutn())
		authReq.SipNumAuthVectors = 2
		authRes, err = client.Authenticate(context.Background(), authReq)
		if err != nil {
			t.Errorf("GRPC MAR Error: %v", err)
			complChan <- err
			return
		}
		if len(authRes.SipAuthVectors) != 2 {
			t.Errorf("Unexpected Number of SIPAuthVectors: %d, Expected: %d", len(authRes.SipAuthVectors), 3)
		}
		assert.Equal(
			t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(15))), authRes.SipAuthVectors[0].GetRandAutn())
		assert.Equal(
			t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(16))), authRes.SipAuthVectors[1].GetRandAutn())
		complChan <- nil
	}
	go testHappyPath(uint32(test_loops + rand.Intn(100)))
	select {
	case err := <-complChan:
		if err != nil {
			t.Fatal(err)
			return
		}
	case <-time.After(time.Second * 5):
		t.Fatal("Timed out. \n" +
			"!!! If this happened during multiple_swx_porxy_test, that may mean that the " +
			"multiplexor sent the request to the wrong diam server (HSS)")
		return
	}

	// Multi-threaded test ensures session-id logic handling works
	for round := 0; round < test_loops; round++ {
		go testHappyPath(uint32(round))
	}
	for round := 0; round < test_loops; round++ {
		testErr := <-complChan
		if testErr != nil {
			t.Fatal(testErr)
			return
		}
	}
}

func initSwxTestSetup(t *testing.T, config *servicers.SwxProxyConfig) string {
	// ---- CORE 3gpp ----
	// create the mockHSS server/servers (depending on the config)
	serverAddr, err := test.StartTestSwxServer(TCPorSCTP, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Started Swx Server at %s", serverAddr)
	// Update config address with address of where test swx server is running
	config.ServerCfg.Addr = serverAddr

	// ---- GRPC ----
	grpcListener, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	service, err := servicers.NewSwxProxy(config)
	if err != nil {
		t.Fatalf("failed to create SwxProxy: %v", err)
	}
	grpcServer := grpc.NewServer()
	protos.RegisterSwxProxyServer(grpcServer, service)
	// start GRPC server
	go func() {
		if err2 := grpcServer.Serve(grpcListener); err2 != nil {
			t.Errorf("failed to serve: %v", err2)
		}
	}()
	grpcAddress := grpcListener.Addr()
	t.Logf("Started Swx GRPC Proxy on %s", grpcAddress.String())
	return grpcAddress.String()
}

func getSwxTestConfig(verify bool) *servicers.SwxProxyConfig {
	return &servicers.SwxProxyConfig{
		ClientCfg: &diameter.DiameterClientConfig{
			Host:  "magma-oai.openair4G.eur", // diameter host
			Realm: "openair4G.eur",           // diameter realm,
		},
		ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:     "",         // to be filled in once server addr is started
			Protocol: TCPorSCTP}, // tcp/sctp
		},
		VerifyAuthorization: verify,
	}
}
