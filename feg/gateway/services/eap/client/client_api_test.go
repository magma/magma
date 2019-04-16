/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package client_test

import (
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/eap/client"
	eap_protos "magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
	"magma/orc8r/cloud/go/test_utils"
)

type testSwxProxy struct{}

var rpcResponseDelay time.Duration
// Test SwxProxyServer implementation
//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s testSwxProxy) Authenticate(
	ctx context.Context,
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {

	time.Sleep(rpcResponseDelay)

	res := &protos.AuthenticationAnswer{
		UserName: req.GetUserName(),
		SipAuthVectors: []*protos.AuthenticationAnswer_SIPAuthVector{
			&protos.AuthenticationAnswer_SIPAuthVector{
				AuthenticationScheme: req.AuthenticationScheme,
				RandAutn: []byte(
					"\x01\x23\x45\x67\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef" +
						"\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c"),
				Xres:               []byte("\x29\x5c\x00\xea\xe3\x88\x93\x0d"),
				ConfidentialityKey: []byte("\xa8\x35\xcf\x22\xb0\xf4\x3e\x15\x19\xd6\xfd\x23\x4c\x00\xd7\x93"),
				IntegrityKey:       []byte("\xd5\x37\x0f\x13\x79\x6f\x2f\x61\x5c\xbe\x15\xef\x9f\x42\x0a\x98"),
			},
		},
	}
	if req.RetrieveUserProfile {
		res.UserProfile = &protos.AuthenticationAnswer_UserProfile{Msisdn: msisdn}
	}
	return res, nil
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s testSwxProxy) Register(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{}, nil
}

var (
	testEapIdentityResp = []byte("\x02\x01\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31" +
		"\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x35\x35\x40\x77\x6c\x61" +
		"\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e" +
		"\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00")
	expectedTestEap = []byte{1, 2, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 1, 35, 69, 103, 137, 171, 205, 239, 1, 35, 69,
		103,137, 171, 205, 239, 2, 5, 0, 0, 84, 171, 100, 74, 144, 81, 185, 185, 94, 133, 193, 34, 62, 14, 241,
		76, 11, 5, 0, 0, 187, 28, 77, 175, 111, 216, 83, 74, 247, 124, 169, 254, 40, 141, 169, 189}

	testEapChallengeResp = []byte("\x02\x02\x00\x28\x17\x01\x00\x00\x0b\x05\x00\x00\xfd\x2b\x50\xbd" +
		"\x32\x24\x7a\xd7\x32\x9d\x9d\x26\x41\x60\x44\x46\x03\x03\x00\x40\x29\x5c\x00\xea\xe3\x88\x93\x0d")
	expectedChallengeResp = []byte{1, 2, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 1, 35, 69, 103, 137, 171, 205, 239, 1, 35, 69, 103,
		137, 171, 205, 239, 2, 5, 0, 0, 84, 171, 100, 74, 144, 81, 185, 185, 94, 133, 193, 34, 62, 14, 241, 76, 11, 5,
		0, 0, 180, 191, 23, 199, 219, 210, 244, 54, 3, 41, 254, 37, 158, 216, 47, 19}
	successEAP = []byte{3, 2, 0, 4}

	expectedMppeRecvKey = []byte(
		"\x95\x63\x3c\x3a\xa5\x8b\x48\xbe\xde\x6d\x2c\x1a\x91\x70\x71\xf5" +
		"\x63\xd4\xed\x7f\xba\xb3\xec\x61\xed\x7e\x3a\xf4\x82\x06\x58\x71" +
		"\x8c\xf7\xee\x86\x81\x0d\xf4\xf9\xf4\xb7\xb9\xdd\x14\xca\xc3\xbd\x95\x80")
	expectedMppeRecvKeySalt = []byte("\x95\x63")

	expectedMppeSendKey = []byte(
		"\x9b\x87\x83\x49\x6a\x78\xcc\xaa\x34\x4e\x45\x51\x7f\x15\x37\xf9" +
		"\x30\x94\x26\x07\x60\x68\x97\xf0\xb5\x69\xab\x1d\x61\x9d\x8b\xa9" +
		"\x85\x3c\xc8\xaf\x68\x4b\xaa\x8f\x8f\x77\x5f\x68\x94\xf0\xcd\xc6\xc9\x2f")
	expectedMppeSendKeySalt = []byte("\x9b\x87")

	authenticator         = []byte{
		0x9f, 0xe8, 0xff, 0xcb, 0xc9, 0xd4, 0x85, 0x97, 0xb9, 0x5b, 0x79, 0x7c, 0x2d, 0xf5, 0x43, 0x31}
	sharedSecret = []byte("1qaz2wsx")
	msisdn = "123456789"
)

func TestEAPClientApi(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService()
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	eapCtx := &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	peap, err := client.Handle(&eap_protos.Eap{Payload: testEapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test EAP: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), []byte(expectedTestEap)) {
		t.Fatalf(
			"Unexpected identityResponse EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), expectedTestEap)
	}

	aka.SetSessionAuthenticatedTimeout(time.Millisecond * 200)

	eapCtx = peap.GetCtx()
	peap, err = client.Handle(&eap_protos.Eap{Payload: testEapChallengeResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test Challenge EAP: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), []byte(successEAP)) {
		t.Fatalf(
			"Unexpected Challenge Response EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), expectedChallengeResp)
	}

	// We should get a valid MSR within the auth success EAP Ctx, verify that we generated valid
	// MS-MPPE-Recv-Key & MS-MPPE-Send-Key according to https://tools.ietf.org/html/rfc2548
	genMS_MPPE_Recv_Key := append(
		expectedMppeRecvKeySalt,
		eap.EncodeMsMppeKey(expectedMppeRecvKeySalt, peap.GetCtx().Msk[0:32], authenticator, sharedSecret)...)

	genMS_MPPE_Send_Key := append(
		expectedMppeSendKeySalt,
		eap.EncodeMsMppeKey(expectedMppeSendKeySalt, peap.GetCtx().Msk[32:], authenticator, sharedSecret)...)

	if !reflect.DeepEqual(genMS_MPPE_Recv_Key, expectedMppeRecvKey) {
		t.Fatalf(
			"MS_MPPE_Recv_Keys mismatch.\n\tGenerated MS_MPPE_Recv_Key(%d): %v\n\tExpected  MS_MPPE_Recv_Key(%d): %v",
			len(genMS_MPPE_Recv_Key), genMS_MPPE_Recv_Key, len(expectedMppeRecvKey), expectedMppeRecvKey)
	}
	if !reflect.DeepEqual(genMS_MPPE_Send_Key, expectedMppeSendKey) {
		t.Fatalf(
			"MS_MPPE_Send_Keys mismatch.\n\tGenerated MS_MPPE_Send_Key(%d): %v\n\tExpected  MS_MPPE_Send_Key(%d): %v",
			len(genMS_MPPE_Send_Key), genMS_MPPE_Send_Key, len(expectedMppeSendKey), expectedMppeSendKey)
	}

	time.Sleep(time.Millisecond * 10)

	eapCtx = peap.GetCtx()
	peap, err = client.Handle(&eap_protos.Eap{Payload: testEapChallengeResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Second Test Challenge EAP within Auth timeout window: %v", err)
	}

	time.Sleep(aka.SessionAuthenticatedTimeout() + time.Millisecond*10)

	eapCtx = peap.GetCtx()
	peap, err = client.Handle(&eap_protos.Eap{Payload: testEapChallengeResp, Ctx: eapCtx})
	if err == nil {
		t.Fatalf("Expected Error for removed Session ID: %s", eapCtx.SessionId)
	}
	grpcCode := status.Convert(err).Code()
	if grpcCode != codes.FailedPrecondition {
		t.Fatalf("Unexpected Error Copde: %d", grpcCode)
	}

	// Test timeout
	currentTout := aka.ChallengeTimeout()
	aka.SetChallengeTimeout(time.Millisecond * 100)
	eapCtx = &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	peap, err = client.Handle(&eap_protos.Eap{Payload: testEapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling second Test EAP: %v", err)
	}
	if eapCtx.Msisdn != msisdn {
		t.Fatalf("Unexpected MSISDN: %s, expected: %s", eapCtx.Msisdn, msisdn)
	}
	time.Sleep(aka.ChallengeTimeout() + time.Millisecond*20)

	eapCtx = peap.GetCtx()
	peap, err = client.Handle(&eap_protos.Eap{Payload: testEapChallengeResp, Ctx: eapCtx})
	if err == nil {
		t.Fatalf("Expected Error for timed out Session ID: %s", eapCtx.SessionId)
	}
	aka.SetChallengeTimeout(currentTout)
}
