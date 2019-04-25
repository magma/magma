/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package client_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"magma/feg/cloud/go/protos/mconfig"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/client"
	eap_protos "magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
	"magma/orc8r/cloud/go/test_utils"
)

type testSwxProxy struct{}

type TestData struct {
	EapIdentityResp,
	ExpectedChallengeReq,
	EapChallengeResp,

	RandAutn,
	Xres,
	ConfidentialityKey,
	IntegrityKey []byte

	IMSI,
	MSISDN string
}

const (
	TestIMSI1 = "001010000000055"
	TestIMSI2 = "001010000000043"
)

var (
	tests = map[string]*TestData{
		TestIMSI1: {
			EapIdentityResp: []byte("\x02\x01\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31" +
				"\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x35\x35\x40\x77\x6c\x61" +
				"\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e" +
				"\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00"),
			ExpectedChallengeReq: []byte{
				1, 2, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 1, 35, 69, 103, 137, 171, 205, 239, 1, 35, 69, 103, 137,
				171, 205, 239, 2, 5, 0, 0, 84, 171, 100, 74, 144, 81, 185, 185, 94, 133, 193, 34, 62, 14, 241,
				76, 11, 5, 0, 0, 187, 28, 77, 175, 111, 216, 83, 74, 247, 124, 169, 254, 40, 141, 169, 189,
			},
			EapChallengeResp: []byte("\x02\x02\x00\x28\x17\x01\x00\x00\x0b\x05\x00\x00\xfd\x2b\x50\xbd" +
				"\x32\x24\x7a\xd7\x32\x9d\x9d\x26\x41\x60\x44\x46\x03\x03\x00\x40\x29\x5c\x00\xea\xe3\x88\x93\x0d"),
			RandAutn: []byte("\x01\x23\x45\x67\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef" +
				"\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c"),
			Xres:               []byte("\x29\x5c\x00\xea\xe3\x88\x93\x0d"),
			ConfidentialityKey: []byte("\xa8\x35\xcf\x22\xb0\xf4\x3e\x15\x19\xd6\xfd\x23\x4c\x00\xd7\x93"),
			IntegrityKey:       []byte("\xd5\x37\x0f\x13\x79\x6f\x2f\x61\x5c\xbe\x15\xef\x9f\x42\x0a\x98"),
			IMSI:               TestIMSI1,
			MSISDN:             "123456789",
		},
		TestIMSI2: {
			EapIdentityResp: []byte("\x02\x02\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31\x30" +
				"\x31\x30\x30\x30\x30\x30\x30\x30\x30\x34\x33\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e" +
				"\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00"),
			ExpectedChallengeReq: []byte{
				1, 3, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 254, 220, 186, 152, 118, 84, 50, 16, 254, 220, 186, 152, 118,
				84, 50, 16, 2, 5, 0, 0, 85, 108, 69, 100, 0, 217, 185, 185, 215, 177, 57, 81, 156, 159, 118, 136,
				11, 5, 0, 0, 9, 176, 57, 57, 175, 141, 130, 36, 60, 20, 41, 206, 233, 71, 100, 170,
			},
			EapChallengeResp: []byte("\x02\x03\x00\x28\x17\x01\x00\x00\x0b\x05\x00\x00\x10\xff\x67\x8d\x06" +
				"\xf2\x59\x09\x1b\x6f\x81\x9e\x5a\x62\x7a\x28\x03\x03\x00\x40\xe7\x17\xf3\x2f\x5d\xc8\xa9\x9b"),
			RandAutn: []byte("\xfe\xdc\xba\x98\x76\x54\x32\x10\xfe\xdc\xba\x98\x76\x54\x32" +
				"\x10\x55\x6c\x45\x64\x00\xd9\xb9\xb9\xd7\xb1\x39\x51\x9c\x9f\x76\x88"),
			Xres:               []byte("\xe7\x17\xf3\x2f\x5d\xc8\xa9\x9b"),
			ConfidentialityKey: []byte("\x21\xb1\x64\x48\x9b\xf5\x04\x7e\xae\x88\xc4\xcd\x7c\xcd\xe3\xc2"),
			IntegrityKey:       []byte("\xf4\xcb\x01\x9b\xed\xc8\x4d\x63\xc6\xce\xa7\xe2\xb0\x77\xfd\xb0"),
			IMSI:               TestIMSI2,
			MSISDN:             "456789012",
		},
	}

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

	authenticator = []byte{
		0x9f, 0xe8, 0xff, 0xcb, 0xc9, 0xd4, 0x85, 0x97, 0xb9, 0x5b, 0x79, 0x7c, 0x2d, 0xf5, 0x43, 0x31,
	}
	sharedSecret = []byte("1qaz2wsx")
	msisdn       = "123456789"

	plmnId5      = "00101"
	plmnId6      = "001010"
	wrongPlmnId6 = "001011"
)

// Test SwxProxyServer implementation
//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s testSwxProxy) Authenticate(
	ctx context.Context,
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {

	time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond * 10))))

	v, ok := tests[req.GetUserName()]
	if !ok {
		return &protos.AuthenticationAnswer{},
			status.Errorf(codes.PermissionDenied, "Unknown User: "+req.GetUserName())
	}
	res := &protos.AuthenticationAnswer{
		UserName: req.GetUserName(),
		SipAuthVectors: []*protos.AuthenticationAnswer_SIPAuthVector{
			&protos.AuthenticationAnswer_SIPAuthVector{
				AuthenticationScheme: req.AuthenticationScheme,
				RandAutn:             v.RandAutn,
				Xres:                 v.Xres,
				ConfidentialityKey:   v.ConfidentialityKey,
				IntegrityKey:         v.IntegrityKey,
			},
		},
	}
	if req.RetrieveUserProfile {
		res.UserProfile = &protos.AuthenticationAnswer_UserProfile{Msisdn: v.MSISDN}
	}
	return res, nil
}

func testAuth(t *testing.T, imsi string, iter int, done chan error) {
	var (
		err  error
		peap *eap_protos.Eap
	)
	defer func() {
		if done != nil {
			done <- err
		}
		if err != nil {
			t.Fatal(err)
		}
	}()

	tst, found := tests[imsi]
	if !found {
		err = fmt.Errorf("Missing Test Data for IMSI: %s", imsi)
		return
	}

	for i := 0; i < iter; i++ {
		eapCtx := &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
		peap, err = client.Handle(&eap_protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
		if err != nil {
			err = fmt.Errorf("Error Handling Test EAP: %v", err)
			return
		}
		if !reflect.DeepEqual([]byte(peap.GetPayload()), tst.ExpectedChallengeReq) {
			err = fmt.Errorf("Unexpected identityResponse EAP\n\tReceived: %s\n\tExpected: %s",
				client.BytesToStr(peap.GetPayload()), client.BytesToStr(tst.ExpectedChallengeReq))
			return
		}
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond * 10))))
		eapCtx = peap.GetCtx()
		peap, err = client.Handle(&eap_protos.Eap{Payload: tst.EapChallengeResp, Ctx: eapCtx})
		if err != nil {
			err = fmt.Errorf("Error Handling Test Challenge EAP: %v", err)
			return
		}
		successp := []byte{eap.SuccessCode, eap.Packet(tst.EapChallengeResp).Identifier(), 0, 4}
		if !reflect.DeepEqual([]byte(peap.GetPayload()), []byte(successp)) {
			err = fmt.Errorf(
				"Unexpected Challenge Response EAP\n\tReceived: %.3v\n\tExpected: %.3v",
				peap.GetPayload(), []byte(successp))
			return
		}
		// Check that we got expected MSISDN with the success EAP
		if peap.GetCtx().Msisdn != tst.MSISDN {
			err = fmt.Errorf("Unexpected MSISDN: %s, expected: %s", eapCtx.Msisdn, tst.MSISDN)
			return
		}
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond * 10))))
	}
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s testSwxProxy) Register(_ context.Context, _ *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{}, nil
}

func TestEAPClientApi(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(nil)
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	go testAuth(t, TestIMSI2, 10, nil) // start IMSI2 tests in parallel

	tst := tests[TestIMSI1]
	eapCtx := &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	peap, err := client.Handle(&eap_protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test EAP: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), tst.ExpectedChallengeReq) {
		t.Fatalf(
			"Unexpected identityResponse EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), tst.ExpectedChallengeReq)
	}

	servicer.SetSessionAuthenticatedTimeout(time.Millisecond * 200)

	eapCtx = peap.GetCtx()
	peap, err = client.Handle(&eap_protos.Eap{Payload: tst.EapChallengeResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test Challenge EAP: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), []byte(successEAP)) {
		t.Fatalf(
			"Unexpected Challenge Response EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), []byte(successEAP))
	}
	// Check that we got expected MSISDN with the success EAP
	if peap.GetCtx().Msisdn != tst.MSISDN {
		t.Fatalf("Unexpected MSISDN: %s, expected: %s", eapCtx.Msisdn, tst.MSISDN)
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
	peap, err = client.Handle(&eap_protos.Eap{Payload: tst.EapChallengeResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Second Test Challenge EAP within Auth timeout window: %v", err)
	}

	time.Sleep(servicer.SessionAuthenticatedTimeout() + time.Millisecond*10)

	eapCtx = peap.GetCtx()
	peap, err = client.Handle(&eap_protos.Eap{Payload: tst.EapChallengeResp, Ctx: eapCtx})
	if err == nil {
		t.Fatalf("Expected Error for removed Session ID: %s", eapCtx.SessionId)
	}
	grpcCode := status.Convert(err).Code()
	if grpcCode != codes.FailedPrecondition {
		t.Fatalf("Unexpected Error Copde: %d", grpcCode)
	}

	// Test timeout
	servicer.SetChallengeTimeout(time.Millisecond * 100)
	eapCtx = &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	peap, err = client.Handle(&eap_protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling second Test EAP: %v", err)
	}
	time.Sleep(servicer.ChallengeTimeout() + time.Millisecond*20)

	eapCtx = peap.GetCtx()
	peap, err = client.Handle(&eap_protos.Eap{Payload: tst.EapChallengeResp, Ctx: eapCtx})
	if err == nil {
		t.Fatalf("Expected Error for timed out Session ID: %s", eapCtx.SessionId)
	}
}

func TestEAPClientApiConcurent(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{
		Timeout: &mconfig.EapAkaConfig_Timeouts{
			ChallengeMs:            300,
			ErrorNotificationMs:    200,
			SessionMs:              500,
			SessionAuthenticatedMs: 1000,
		}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	done := make(chan error)
	go testAuth(t, TestIMSI1, 99, done)
	go testAuth(t, TestIMSI2, 88, done)
	testAuth(t, TestIMSI1, 77, nil)
	<-done
	<-done // wait for test 1 & 2 to complete
}

// Not used (panic) SwxProxyServer implementation
type noUseSwxProxy struct{}

//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s noUseSwxProxy) Authenticate(
	ctx context.Context,
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {

	return nil, fmt.Errorf("Authenticate is NOT IMPLEMENTED")
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s noUseSwxProxy) Register(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{}, fmt.Errorf("Register is NOT IMPLEMENTED")
}

func TestEAPAkaWrongPlmnId(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service noUseSwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnId6}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	tst := tests[TestIMSI1]
	eapCtx := &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	_, err = client.Handle(&eap_protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err == nil {
		t.Fatalf("Expected Error Handling Filtered PLMN ID")
	}
}

func TestEAPAkaPlmnId5(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnId6, plmnId5}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	servicer.SetChallengeTimeout(time.Millisecond * 10)
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	tst := tests[TestIMSI1]
	eapCtx := &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	peap, err := client.Handle(&eap_protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test EAP: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), tst.ExpectedChallengeReq) {
		t.Fatalf(
			"Unexpected identityResponse EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), tst.ExpectedChallengeReq)
	}
}

func TestEAPAkaPlmnId6(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service testSwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnId6, plmnId6}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	servicer.SetChallengeTimeout(time.Millisecond * 10)
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	tst := tests[TestIMSI1]
	eapCtx := &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	peap, err := client.Handle(&eap_protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err != nil {
		t.Fatalf("Error Handling Test EAP: %v", err)
	}
	if !reflect.DeepEqual([]byte(peap.GetPayload()), tst.ExpectedChallengeReq) {
		t.Fatalf(
			"Unexpected identityResponse EAP\n\tReceived: %.3v\n\tExpected: %.3v",
			peap.GetPayload(), tst.ExpectedChallengeReq)
	}
}
