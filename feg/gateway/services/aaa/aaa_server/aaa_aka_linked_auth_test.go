// +build link_local_service

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

// Package eap_router_test implements eap router unit tests
package main_test

import (
	"os"
	"reflect"
	"testing"

	cp "magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	aaa_client "magma/feg/gateway/services/aaa/client"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	eap_client "magma/feg/gateway/services/eap/client"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
	eap_test "magma/feg/gateway/services/eap/test"
	"magma/orc8r/cloud/go/test_utils"
)

// TestEapAkaConcurent tests EAP AKA Provider
func TestLinkedEapAkaConcurent(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.SwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.AAA_SERVER)
	protos.RegisterAuthenticatorServer(rtrSrv.GrpcServer, &testAuthenticator{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	client := &testEapServiceClient{}
	done := make(chan error)
	go eap_test.Auth(t, client, eap_test.IMSI1, 50, done)
	go eap_test.Auth(t, client, eap_test.IMSI2, 47, done)
	eap_test.Auth(t, client, eap_test.IMSI1, 43, nil) // no need for done chan, running in main thread
	<-done
	<-done // wait for test 1 & 2 to complete
}

// TestLinkedEAPPeerNak verifies EAP legacy NAK mechanism triggered by unsupported EAP Method.
// See: https://tools.ietf.org/html/rfc3748#section-5.3 for more details
//
// "...Where a peer receives a Request for an unacceptable authentication
//      Type (4-253,255), or a peer lacking support for Expanded Types
//      receives a Request for Type 254, a Nak Response (Type 3) MUST be
//      sent.  The Type-Data field of the Nak Response (Type 3) MUST
//      contain one or more octets indicating the desired authentication
//      Type(s), one octet per Type, or the value zero (0) to indicate no
//      proposed alternative."
//
func TestLinkedEAPPeerNak(t *testing.T) {
	failureEAP := []byte{4, 237, 0, 4}
	akaPrimeIdentity := eap.NewPacket(
		eap.ResponseCode, 236,
		append([]byte{eap.MethodIdentity}, []byte("6001010000000091@wlan.mnc001.mcc001.3gppnetwork.org")...))
	permIdReq := []byte{0x01, 237, 0x00, 0x0c, 0x17, 0x05, 0x00, 0x00, 0x0a, 0x01, 0x00, 0x00}
	akaPrimeNak := []byte{0x02, 237, 0x00, 0x06, 0x03, 50}
	akaAkaPrimeNak := []byte{0x02, 236, 0x00, 0x07, 0x03, 50, 23}

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.AAA_SERVER)
	protos.RegisterAuthenticatorServer(rtrSrv.GrpcServer, &testAuthenticator{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	eapCtx := &protos.Context{SessionId: eap.CreateSessionId()}

	peap, err := aaa_client.HandleIdentity(&protos.EapIdentity{Payload: akaPrimeIdentity, Ctx: eapCtx, Method: 23})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), permIdReq) {
		t.Fatalf("Unexpected Identity Responsen\tReceived: %.3v\n\tExpected: %.3v", peap.GetPayload(), permIdReq)
	}
	peap, err = aaa_client.Handle(&protos.Eap{Payload: akaPrimeNak, Ctx: peap.Ctx})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), failureEAP) {
		t.Fatalf("Unexpected AKA' Nak Response\n\tReceived: %.3v\n\tExpected: %.3v", peap.GetPayload(), failureEAP)
	}
	peap, err = aaa_client.Handle(&protos.Eap{Payload: akaAkaPrimeNak, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), permIdReq) {
		t.Fatalf("Unexpected AKA['] Nak Response\n\tReceived: %.3v\n\tExpected: %.3v", peap.GetPayload(), permIdReq)
	}
}
