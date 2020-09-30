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
package handlers

import (
	"os"
	"reflect"
	"testing"

	cp "magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

var (
	testEapChallengeResp = "\x02\x78\x00\x1c\x12\x0b\x00\x00\x0b\x05\x00\x00\x16\x96\xb1\x64" +
		"\x14\x9f\xd6\x39\x70\xb0\xe3\x36\xc8\x5d\x00\x61"
	successEAP = []byte{3, 120, 0, 4}
)

func TestSimChallengeResp(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	simSrv, _ := servicers.NewEapSimService(nil)
	eapCtx := &protos.Context{}
	// Initialize CTX
	p, err := startResponse(simSrv, eapCtx, eap.Packet(testEapStartResp))
	if err != nil {
		t.Fatalf("Unexpected identityResponse error: %v", err)
	}
	if len(eapCtx.SessionId) == 0 {
		t.Fatal("Empty Session ID")
	}
	if !reflect.DeepEqual([]byte(p), []byte(expectedTestEapChallengeReq)) {
		t.Fatalf("Unexpected identityResponse EAP\n\tReceived: %v\n\tExpected: %v", p, []byte(expectedTestEapChallengeReq))
	}
	p, err = challengeResponse(simSrv, eapCtx, eap.Packet(testEapChallengeResp))
	if err != nil {
		t.Fatalf("Unexpected challengeResponse error: %v", err)
	}
	if !reflect.DeepEqual([]byte(p), successEAP) {
		t.Fatalf("Unexpected challengeResponse EAP\n\tReceived: %v\n\tExpected: %v", p, successEAP)
	}
}
