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
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

const (
	testEapChallengeResp = "\x02\x02\x00\x28\x17\x01\x00\x00\x0b\x05\x00\x00\xfd\x2b\x50\xbd" +
		"\x32\x24\x7a\xd7\x32\x9d\x9d\x26\x41\x60\x44\x46\x03\x03\x00\x40\x29\x5c\x00\xea\xe3\x88\x93\x0d"
)

var (
	expectedChallengeResp = []byte{1, 2, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 1, 35, 69, 103, 137, 171, 205, 239, 1, 35, 69, 103,
		137, 171, 205, 239, 2, 5, 0, 0, 84, 171, 100, 74, 144, 81, 185, 185, 94, 133, 193, 34, 62, 14, 241, 76, 11, 5,
		0, 0, 180, 191, 23, 199, 219, 210, 244, 54, 3, 41, 254, 37, 158, 216, 47, 19}
	successEAP = []byte{3, 2, 0, 4}
)

func TestAkaChallengeResp(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	akaSrv, _ := servicers.NewEapAkaService(nil)
	eapCtx := &protos.Context{}
	// Initialize CTX
	p, err := identityResponse(akaSrv, eapCtx, eap.Packet(testEapIdentityResp))
	if err != nil {
		t.Fatalf("Unexpected identityResponse error: %v", err)
	}
	if len(eapCtx.SessionId) == 0 {
		t.Fatal("Empty Session ID")
	}
	if !reflect.DeepEqual([]byte(p), expectedTestEapChallengeResp) {
		t.Fatalf("Unexpected identityResponse EAP\n\tReceived: %v\n\tExpected: %v", p, expectedTestEapChallengeResp)
	}

	p, err = challengeResponse(akaSrv, eapCtx, eap.Packet(testEapChallengeResp))
	if err != nil {
		t.Fatalf("Unexpected challengeResponse error: %v", err)
	}
	if !reflect.DeepEqual([]byte(p), successEAP) {
		t.Fatalf("Unexpected challengeResponse EAP\n\tReceived: %v\n\tExpected: %v", p, expectedTestEapChallengeResp)
	}

}
