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

package s6a_proxy_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/service_health"
	"magma/feg/gateway/services/s6a_proxy"
	"magma/feg/gateway/services/s6a_proxy/servicers/test"
	"magma/feg/gateway/services/s6a_proxy/test_init"
)

var (
	diamServerAddr = fmt.Sprintf("127.0.0.1:%d", 30000+rand.Intn(1000))
	TCPorSCTP      = systemBasedTCPorSCTP() // sctp if run in linux, tcp if run in MAC
)

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

func TestS6aProxyClient(t *testing.T) {
	err := test_init.StartTestService(t, TCPorSCTP, diamServerAddr)
	if err != nil {
		t.Fatal(err)
		return
	}

	req := &protos.AuthenticationInformationRequest{
		UserName:                   test.TEST_IMSI,
		VisitedPlmn:                []byte(test.TEST_PLMN_ID),
		NumRequestedEutranVectors:  3,
		ImmediateResponsePreferred: true,
	}
	// AIR
	r, err := s6a_proxy.AuthenticationInformation(req)
	if err != nil {
		t.Fatalf("GRPC AIR Error: %v", err)
		return
	}
	t.Logf("GRPC AIA: %#+v", *r)
	if r.ErrorCode != protos.ErrorCode_UNDEFINED {
		t.Errorf("Unexpected AIA Error Code: %d", r.ErrorCode)
	}
	if len(r.EutranVectors) != 3 {
		t.Errorf("Unexpected Number of EutranVectors: %d, Expected: 3", len(r.EutranVectors))
	}

	// AIR with UTRAN Vectors
	req.NumRequestedUtranGeranVectors = 3
	r, err = s6a_proxy.AuthenticationInformation(req)
	if err != nil {
		t.Fatalf("GRPC AIR with UTRAN Error: %v", err)
		return
	}
	t.Logf("GRPC AIA with UTRAN: %#+v", *r)
	if r.ErrorCode != protos.ErrorCode_UNDEFINED {
		t.Errorf("Unexpected AIA Error Code: %d", r.ErrorCode)
	}
	if len(r.EutranVectors) != 3 {
		t.Errorf("Unexpected Number of EutranVectors: %d, Expected: 3", len(r.EutranVectors))
	}
	if len(r.UtranVectors) != 3 {
		t.Errorf("Unexpected Number of UtranVectors: %d, Expected: 3", len(r.UtranVectors))
	}

	// AIR with UTRAN Only Vectors
	req.NumRequestedEutranVectors = 0
	r, err = s6a_proxy.AuthenticationInformation(req)
	if err != nil {
		t.Fatalf("GRPC AIR with UTRAN Error: %v", err)
		return
	}
	t.Logf("GRPC AIA with UTRAN: %#+v", *r)
	if r.ErrorCode != protos.ErrorCode_UNDEFINED {
		t.Errorf("Unexpected AIA Error Code: %d", r.ErrorCode)
	}
	if len(r.EutranVectors) != 0 {
		t.Errorf("Unexpected Number of EutranVectors: %d, Expected: 0", len(r.EutranVectors))
	}
	if len(r.UtranVectors) != 3 {
		t.Errorf("Unexpected Number of UtranVectors: %d, Expected: 3", len(r.UtranVectors))
	}

	ulReq := &protos.UpdateLocationRequest{
		UserName:                     test.TEST_IMSI,
		VisitedPlmn:                  []byte(test.TEST_PLMN_ID),
		SkipSubscriberData:           false,
		InitialAttach:                true,
		DualRegistration_5GIndicator: true,
	}
	// ULR
	ulResp, err := s6a_proxy.UpdateLocation(ulReq)
	if err != nil {
		t.Fatalf("GRPC ULR Error: %v", err)
		return
	}
	t.Logf("GRPC ULA: %#+v", *ulResp)
	if ulResp.ErrorCode != protos.ErrorCode_UNDEFINED {
		t.Errorf("Unexpected ULA Error Code: %d", r.ErrorCode)
	}

	puReq := &protos.PurgeUERequest{
		UserName: test.TEST_IMSI,
	}
	// PUR
	puResp, err := s6a_proxy.PurgeUE(puReq)
	if err != nil {
		t.Fatalf("GRPC PUR Error: %v", err)
	}
	t.Logf("GRPC PUA: %#+v", *puResp)
	if puResp.ErrorCode != protos.ErrorCode_SUCCESS {
		t.Errorf("Unexpected PUA Error Code: %d", r.ErrorCode)
	}

	// Disable connections and ensure subsequent requests fail
	disableReq := &protos.DisableMessage{
		DisablePeriodSecs: 10,
	}
	err = service_health.Disable(registry.S6A_PROXY, disableReq)
	if err != nil {
		t.Fatalf("GRPC ServiceHealth Disable Error: %v", err)
		return
	}

	// AIR should fail
	_, err = s6a_proxy.AuthenticationInformation(req)
	if err == nil {
		t.Errorf("AIR Succeeded, but should have failed due to disabled connections")
	}

	// Enable connections
	err = service_health.Enable(registry.S6A_PROXY)
	if err != nil {
		t.Fatalf("GRPC ServiceHealth Enable Error: %v", err)
		return
	}

	// ULR should pass now
	ulResp, err = s6a_proxy.UpdateLocation(ulReq)
	if err != nil {
		t.Fatalf("GRPC ULR Error: %v", err)
		return
	}
	t.Logf("GRPC ULA: %#+v", *ulResp)
	if ulResp.ErrorCode != protos.ErrorCode_UNDEFINED {
		t.Errorf("Unexpected ULA Error Code: %d", ulResp.ErrorCode)
	}
}
