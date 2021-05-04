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
	"magma/feg/gateway/services/eap/providers/aka"

	cp "magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

type testSwxProxy struct{}

// Test SwxProxyServer implementation
//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s testSwxProxy) Authenticate(
	ctx context.Context,
	req *cp.AuthenticationRequest,
) (*cp.AuthenticationAnswer, error) {
	return &cp.AuthenticationAnswer{
		UserName: req.GetUserName(),
		SipAuthVectors: []*cp.AuthenticationAnswer_SIPAuthVector{
			{
				AuthenticationScheme: req.AuthenticationScheme,
				RandAutn: []byte(
					"\x01\x23\x45\x67\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef" +
						"\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c"),
				Xres:               []byte("\x29\x5c\x00\xea\xe3\x88\x93\x0d"),
				ConfidentialityKey: []byte("\xa8\x35\xcf\x22\xb0\xf4\x3e\x15\x19\xd6\xfd\x23\x4c\x00\xd7\x93"),
				IntegrityKey:       []byte("\xd5\x37\x0f\x13\x79\x6f\x2f\x61\x5c\xbe\x15\xef\x9f\x42\x0a\x98"),
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
	testEapIdentityResp = "\x02\x01\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31" +
		"\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x35\x35\x40\x77\x6c\x61" +
		"\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e" +
		"\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00"

	expectedRand = "\x01\x05\x00\x00\x01\x23\x45\x67\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef"
	expectedAutn = "\x02\x05\x00\x00\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c"
	identityAttr = "\x0e\x0e\x00\x33\x30\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30" +
		"\x30\x30\x35\x35\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31" +
		"\x2e\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77" +
		"\x6f\x72\x6b\x2e\x6f\x72\x67\x00"
)

var (
	expectedTestMac = []byte{11, 5, 0, 0, 187, 28, 77, 175, 111, 216, 83, 74, 247, 124, 169, 254, 40, 141, 169, 189}

	expectedTestEapChallengeResp = []byte{1, 2, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 1, 35, 69, 103, 137, 171, 205, 239, 1,
		35, 69, 103, 137, 171, 205, 239, 2, 5, 0, 0, 84, 171, 100, 74, 144, 81, 185, 185, 94, 133, 193, 34, 62, 14, 241,
		76, 11, 5, 0, 0, 187, 28, 77, 175, 111, 216, 83, 74, 247, 124, 169, 254, 40, 141, 169, 189}
)

func TestChallengeEAPTemplate(t *testing.T) {
	if challengeReqTemplateLen != 68 {
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
	if attr.Type() != aka.AT_RAND || attr.Len() != 20 {
		t.Fatalf("Invalid AT_RAND: %v\n", attr.Marshaled())
	}

	attr, err = scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_AUTN: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_AUTN Attribute")
	}
	if attr.Type() != aka.AT_AUTN || attr.Len() != 20 {
		t.Fatalf("Invalid AT_AUTN: %v\n", attr.Marshaled())
	}

	attr, err = scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_MAC: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_MAC Attribute")
	}
	if attr.Type() != aka.AT_MAC || attr.Len() != 20 {
		t.Fatalf("Invalid AT_MAC: %v\n", attr.Marshaled())
	}

	fullId, imsi, err := getIMSIIdentity(eap.NewRawAttribute([]byte(identityAttr)))
	if err != nil {
		t.Fatalf("getIMSIIdentity Error: %v", err)
	}
	if fullId != "0001010000000055@wlan.mnc001.mcc001.3gppnetwork.org" {
		t.Fatalf("Unexpected full Identity: %s", fullId)
	}
	if imsi != "0001010000000055" {
		t.Fatalf("Unexpected IMSI: %s", imsi)
	}
}

func TestAkaChallenge(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	akaSrv, _ := servicers.NewEapAkaService(nil)
	p, err := identityResponse(akaSrv, &protos.Context{}, eap.Packet(testEapIdentityResp))

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
	if attr.Type() != aka.AT_RAND || !reflect.DeepEqual(attr.Marshaled(), []byte(expectedRand)) {
		t.Fatalf("Invalid AT_RAND:\n\tExpected: %v\n\tReceived: %v\n", []byte(expectedRand), attr.Marshaled())
	}
	attr, err = scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_AUTN: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_AUTN Attribute")
	}
	if attr.Type() != aka.AT_AUTN || !reflect.DeepEqual(attr.Marshaled(), []byte(expectedAutn)) {
		t.Fatalf("Invalid AT_AUTN:\n\tExpected: %v\n\tReceived: %v\n", []byte(expectedAutn), attr.Marshaled())
	}
	attr, err = scanner.Next()
	if err != nil {
		t.Fatalf("Error getting AT_MAC: %v", err)
	}
	if attr == nil {
		t.Fatal("Nil AT_MAC Attribute")
	}
	if attr.Type() != aka.AT_MAC || !reflect.DeepEqual(attr.Marshaled(), expectedTestMac) {
		t.Fatalf("Invalid AT_MAC:\n\tExpected: %v\n\tReceived: %v\n", expectedTestMac, attr.Marshaled())
	}
	if !reflect.DeepEqual([]byte(p), expectedTestEapChallengeResp) {
		t.Fatalf("Unexpected identityResponse EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			p, expectedTestEapChallengeResp)
	}
}
