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
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	cp "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	aaa_client "magma/feg/gateway/services/aaa/client"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	eap_client "magma/feg/gateway/services/eap/client"
	eapp "magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
	eap_test "magma/feg/gateway/services/eap/test"
	"magma/orc8r/cloud/go/test_utils"
)

// Test AAA EAP Service
type testAuthenticator struct {
	supportedMethods []byte
}

func (s *testAuthenticator) HandleIdentity(ctx context.Context, in *protos.EapIdentity) (*protos.Eap, error) {
	resp, err := eap_client.HandleIdentityResponse(
		uint8(in.GetMethod()), &protos.Eap{Payload: in.Payload, Ctx: in.Ctx})
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		err = nil
	}
	return resp, err
}
func (s *testAuthenticator) Handle(ctx context.Context, in *protos.Eap) (*protos.Eap, error) {
	resp, err := eap_client.Handle(in)
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		err = nil
	}
	return resp, err

}
func (s *testAuthenticator) SupportedMethods(ctx context.Context, in *protos.Void) (*protos.EapMethodList, error) {
	return &protos.EapMethodList{Methods: s.supportedMethods}, nil
}

var (
	plmnID5      = "00101"
	plmnID6      = "001010"
	wrongPlmnID6 = "001011"
)

type testEapServiceClient struct{}

func (c testEapServiceClient) Handle(in *protos.Eap) (*protos.Eap, error) {
	return aaa_client.Handle(in)
}

func (c testEapServiceClient) HandleIdentity(in *protos.EapIdentity) (*protos.Eap, error) {
	return aaa_client.HandleIdentity(in)
}

// TestEapAkaConcurent tests EAP AKA Provider
func TestEapAkaConcurent(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.SwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{
		Timeout: &mconfig.EapAkaConfig_Timeouts{
			ChallengeMs:            700,
			ErrorNotificationMs:    200,
			SessionMs:              900,
			SessionAuthenticatedMs: 1000,
		}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eapp.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.AAA_SERVER)
	protos.RegisterAuthenticatorServer(rtrSrv.GrpcServer, &testAuthenticator{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	client := &testEapServiceClient{}
	done := make(chan error)
	go eap_test.Auth(t, client, eap_test.IMSI1, 50, done)
	go eap_test.Auth(t, client, eap_test.IMSI2, 47, done)
	eap_test.Auth(t, client, eap_test.IMSI1, 43, nil)
	<-done
	<-done // wait for test 1 & 2 to complete
}

func TestEAPPeerNak(t *testing.T) {
	failureEAP := []byte{4, 237, 0, 4}
	akaPrimeIdentity := eap.NewPacket(
		eap.ResponseCode, 236,
		append([]byte{eap.MethodIdentity}, []byte("6001010000000091@wlan.mnc001.mcc001.3gppnetwork.org")...))
	permIdReq := []byte{0x01, 237, 0x00, 0x0c, 0x17, 0x05, 0x00, 0x00, 0x0a, 0x01, 0x00, 0x00}
	akaPrimeNak := []byte{0x02, 237, 0x00, 0x06, 0x03, 50}
	akaAkaPrimeNak := []byte{0x02, 236, 0x00, 0x07, 0x03, 50, 23}

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(nil)
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eapp.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.AAA_SERVER)
	protos.RegisterAuthenticatorServer(rtrSrv.GrpcServer, &testAuthenticator{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	eapCtx := &protos.Context{SessionId: eap.CreateSessionId()}

	peap, err := aaa_client.HandleIdentity(&protos.EapIdentity{Payload: akaPrimeIdentity, Ctx: eapCtx, Method: 23})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	if !reflect.DeepEqual(peap.GetPayload(), permIdReq) {
		t.Fatalf("Unexpected Identity Responsen\tReceived: %.3v\n\tExpected: %.3v", peap.GetPayload(), permIdReq)
	}
	peap, err = aaa_client.Handle(&protos.Eap{Payload: akaPrimeNak, Ctx: peap.Ctx})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	if !reflect.DeepEqual(peap.GetPayload(), failureEAP) {
		t.Fatalf("Unexpected AKA' Nak Response\n\tReceived: %.3v\n\tExpected: %.3v", peap.GetPayload(), failureEAP)
	}
	peap, err = aaa_client.Handle(&protos.Eap{Payload: akaAkaPrimeNak, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	if !reflect.DeepEqual(peap.GetPayload(), permIdReq) {
		t.Fatalf("Unexpected AKA['] Nak Response\n\tReceived: %.3v\n\tExpected: %.3v", peap.GetPayload(), permIdReq)
	}
}

func TestEAPAkaWrongPlmnId(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.NoUseSwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnID6}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eapp.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.AAA_SERVER)
	protos.RegisterAuthenticatorServer(rtrSrv.GrpcServer, &testAuthenticator{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	tst := eap_test.Units[eap_test.IMSI1]
	eapCtx := &protos.Context{SessionId: eap.CreateSessionId()}
	peap, err := aaa_client.Handle(&protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test EAP: %v", err)
	}
	notifAkaEap := aka.NewAKANotificationReq(eap.Packet(tst.EapIdentityResp).Identifier(), aka.NOTIFICATION_FAILURE)
	if !reflect.DeepEqual(peap.GetPayload(), []byte(notifAkaEap)) {
		t.Fatalf(
			"Unexpected identityResponse Notification\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), notifAkaEap)
	}
}

func TestEAPAkaPlmnId5(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.SwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnID6, plmnID5}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}

	servicer.SetChallengeTimeout(time.Millisecond * 10)
	eapp.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.AAA_SERVER)
	protos.RegisterAuthenticatorServer(rtrSrv.GrpcServer, &testAuthenticator{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	tst := eap_test.Units[eap_test.IMSI1]
	eapCtx := &protos.Context{SessionId: eap.CreateSessionId()}
	peap, err := aaa_client.Handle(&protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test EAP: %v", err)
	}
	if !reflect.DeepEqual(peap.GetPayload(), tst.ExpectedChallengeReq) {
		t.Fatalf(
			"Unexpected identityResponse EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), tst.ExpectedChallengeReq)
	}
}

func TestEAPAkaPlmnId6(t *testing.T) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.SwxProxy
	cp.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnID6, plmnID6}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	servicer.SetChallengeTimeout(time.Millisecond * 10)
	eapp.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.AAA_SERVER)
	protos.RegisterAuthenticatorServer(rtrSrv.GrpcServer, &testAuthenticator{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	tst := eap_test.Units[eap_test.IMSI1]
	eapCtx := &protos.Context{SessionId: eap.CreateSessionId()}
	peap, err := aaa_client.Handle(&protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test EAP: %v", err)
	}
	if !reflect.DeepEqual(peap.GetPayload(), tst.ExpectedChallengeReq) {
		t.Fatalf(
			"Unexpected identityResponse EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), tst.ExpectedChallengeReq)
	}
}
