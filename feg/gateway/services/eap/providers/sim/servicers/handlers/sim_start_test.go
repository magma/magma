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
	"context"
	"os"
	"reflect"
	"testing"

	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim"

	cp "magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

type testSwxProxy struct{}

// Test SwxProxyServer implementation
//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s testSwxProxy) Authenticate(
	_ context.Context,
	req *cp.AuthenticationRequest,
) (*cp.AuthenticationAnswer, error) {
	return &cp.AuthenticationAnswer{
		UserName: req.GetUserName(),
		SipAuthVectors: []*cp.AuthenticationAnswer_SIPAuthVector{
			{
				AuthenticationScheme: req.AuthenticationScheme,
				RandAutn: []byte{57, 22, 40, 33, 82, 189, 193, 89, 219, 31, 18, 64, 95, 197, 50,
					240, 188, 167, 68, 25, 19, 11, 128, 0, 228, 20, 201, 246, 253, 57, 224, 99},
				Xres:               []byte{155, 36, 12, 2, 227, 43, 246, 254},
				ConfidentialityKey: []byte{235, 74, 254, 58, 73, 108, 112, 173, 61, 24, 169, 176, 219, 233, 85, 180},
				IntegrityKey:       []byte{8, 114, 43, 29, 82, 150, 220, 38, 242, 123, 82, 108, 116, 174, 27, 212},
			},
			{
				AuthenticationScheme: req.AuthenticationScheme,
				RandAutn: []byte{127, 70, 44, 220, 221, 96, 68, 186, 152, 38, 223, 29, 92, 21, 1, 60,
					131, 208, 242, 222, 202, 147, 128, 0, 154, 211, 214, 217, 92, 43, 101, 232},
				Xres:               []byte{135, 169, 135, 251, 190, 133, 65, 108},
				ConfidentialityKey: []byte{79, 73, 201, 197, 199, 254, 178, 13, 168, 21, 85, 129, 186, 164, 41, 106},
				IntegrityKey:       []byte{16, 75, 4, 255, 189, 104, 158, 100, 49, 214, 172, 248, 77, 102, 249, 214},
			},
			{
				AuthenticationScheme: req.AuthenticationScheme,
				RandAutn: []byte{209, 140, 161, 3, 11, 175, 197, 177, 85, 62, 129, 182, 231, 17, 8, 115,
					197, 218, 91, 111, 30, 12, 128, 0, 180, 179, 86, 226, 118, 89, 17, 228},
				Xres:               []byte{38, 215, 150, 33, 249, 61, 236, 70},
				ConfidentialityKey: []byte{95, 22, 168, 51, 109, 73, 55, 191, 123, 171, 50, 191, 255, 146, 190, 145},
				IntegrityKey:       []byte{187, 186, 245, 206, 25, 39, 73, 19, 135, 193, 86, 169, 24, 135, 112, 202},
			},
		},
	}, nil
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s testSwxProxy) Register(
	ctx context.Context,
	req *cp.RegistrationRequest,
) (*cp.RegistrationAnswer, error) {
	return &cp.RegistrationAnswer{}, nil
}

// Deregister sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s testSwxProxy) Deregister(
	ctx context.Context,
	req *cp.RegistrationRequest,
) (*cp.RegistrationAnswer, error) {
	return &cp.RegistrationAnswer{}, nil
}

const (
	testEapStartResp = "\x02\x77\x00\x58\x12\x0a\x00\x00\x10\x01\x00\x01\x07\x05\x00\x00" +
		"\x89\x8f\x4d\xe9\x40\x1f\x13\xcc\xe4\xb7\x8b\xd4\xa6\x6e\xd3\x4b" +
		"\x0e\x0e\x00\x33\x31\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30" +
		"\x30\x31\x31\x39\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31" +
		"\x2e\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77" +
		"\x6f\x72\x6b\x2e\x6f\x72\x67\x00"

	expectedChellengeRand = "\x01\x0d\x00\x00\x39\x16\x28\x21\x52\xbd\xc1\x59\xdb\x1f\x12\x40" +
		"\x5f\xc5\x32\xf0\x7f\x46\x2c\xdc\xdd\x60\x44\xba\x98\x26\xdf\x1d" +
		"\x5c\x15\x01\x3c\xd1\x8c\xa1\x03\x0b\xaf\xc5\xb1\x55\x3e\x81\xb6" +
		"\xe7\x11\x08\x73"

	expectedChellengeMac = "\x0b\x05\x00\x00\xf4\x2e\x55\xe4\x97\x3c\xfb\x11\xcb\xc9\xc6\xa5" +
		"\x2a\xde\x5c\x94"

	identityAttr = "\x0e\x0e\x00\x33\x31\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30" +
		"\x30\x31\x31\x39\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31" +
		"\x2e\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77" +
		"\x6f\x72\x6b\x2e\x6f\x72\x67\x00"

	expectedTestEapChallengeReq = "\x01\x78\x00\x50\x12\x0b\x00\x00\x01\x0d\x00\x00\x39\x16\x28\x21" +
		"\x52\xbd\xc1\x59\xdb\x1f\x12\x40\x5f\xc5\x32\xf0\x7f\x46\x2c\xdc" +
		"\xdd\x60\x44\xba\x98\x26\xdf\x1d\x5c\x15\x01\x3c\xd1\x8c\xa1\x03" +
		"\x0b\xaf\xc5\xb1\x55\x3e\x81\xb6\xe7\x11\x08\x73\x0b\x05\x00\x00" +
		"\xf4\x2e\x55\xe4\x97\x3c\xfb\x11\xcb\xc9\xc6\xa5\x2a\xde\x5c\x94"
)

func TestChallengeEAPTemplate(t *testing.T) {
	if challengeReqTemplateLen != 80 {
		t.Fatalf("Invalid challengeReqTemplateLen: %d", challengeReqTemplateLen)
	}
	scanner, _ := eap.NewAttributeScanner(challengeReqTemplate)
	if scanner == nil {
		t.Fatal("Nil Attribute Scanner")
	}
	attr, err := scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_RAND: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_RAND Attribute")
	}
	if attr.Type() != sim.AT_RAND || attr.Len() != 52 {
		t.Fatalf("Invalid AT_RAND: %v\n", attr.Marshaled())
	}
	attr, err = scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_MAC: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_MAC Attribute")
	}
	if attr.Type() != sim.AT_MAC || attr.Len() != 20 {
		t.Fatalf("Invalid AT_MAC: %v\n", attr.Marshaled())
	}
	fullId, imsi, err := getIMSIIdentity(eap.NewRawAttribute([]byte(identityAttr)))
	if err != nil {
		t.Fatalf("getIMSIIdentity Error: %v", err)
	}
	if fullId != "1001010000000119@wlan.mnc001.mcc001.3gppnetwork.org" {
		t.Fatalf("Unexpected full Identity: %s", fullId)
	}
	if imsi != "1001010000000119" {
		t.Fatalf("Unexpected IMSI: %s", imsi)
	}
}

func TestSimChallenge(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	simSrv, _ := servicers.NewEapSimService(nil)
	p, err := startResponse(simSrv, &protos.Context{}, eap.Packet(testEapStartResp))
	if err != nil {
		t.Fatalf("Unexpected identityResponse error: %v", err)
	}
	scanner, _ := eap.NewAttributeScanner(p)
	if scanner == nil {
		t.Fatal("Nil Attribute Scanner")
	}
	attr, err := scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_RAND: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_RAND Attribute")
	}
	if attr.Type() != sim.AT_RAND || !reflect.DeepEqual(attr.Marshaled(), []byte(expectedChellengeRand)) {
		t.Fatalf("Invalid AT_RAND:\n\tExpected: %v\n\tReceived: %v\n", []byte(expectedChellengeRand), attr.Marshaled())
	}
	attr, err = scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_MAC: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_MAC Attribute")
	}
	if attr.Type() != sim.AT_MAC || !reflect.DeepEqual(attr.Marshaled(), []byte(expectedChellengeMac)) {
		t.Fatalf("Invalid AT_AUTN:\n\tExpected: %v\n\tReceived: %v\n", []byte(expectedChellengeMac), attr.Marshaled())
	}
	if !reflect.DeepEqual([]byte(p), []byte(expectedTestEapChallengeReq)) {
		t.Fatalf("Unexpected SIM Start Challenge EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			p, []byte(expectedTestEapChallengeReq))
	}
}
