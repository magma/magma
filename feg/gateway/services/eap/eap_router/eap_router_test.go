/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// Package eap_router_test implements eap router unit tests
package main_test

import (
	"reflect"
	"testing"
	"time"

	"magma/feg/gateway/services/eap"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	eap_client "magma/feg/gateway/services/eap/client"
	eap_protos "magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
	eap_test "magma/feg/gateway/services/eap/test"
	"magma/orc8r/cloud/go/test_utils"
)

// Test EAP Router Service
type testEapRouter struct {
	supportedMethods []byte
}

func (s *testEapRouter) HandleIdentity(ctx context.Context, in *eap_protos.EapIdentity) (*eap_protos.Eap, error) {
	return eap_client.HandleIdentityResponse(uint8(in.GetMethod()), &eap_protos.Eap{Payload: in.Payload, Ctx: in.Ctx})
}
func (s *testEapRouter) Handle(ctx context.Context, in *eap_protos.Eap) (*eap_protos.Eap, error) {
	return eap_client.Handle(in)
}
func (s *testEapRouter) SupportedMethods(ctx context.Context, in *eap_protos.Void) (*eap_protos.MethodList, error) {
	return &eap_protos.MethodList{Methods: s.supportedMethods}, nil
}

var (
	plmnID5      = "00101"
	plmnID6      = "001010"
	wrongPlmnID6 = "001011"
)

type testEapServiceClient struct {
	eap_protos.EapRouterClient
}

func (c testEapServiceClient) Handle(in *eap_protos.Eap) (*eap_protos.Eap, error) {
	return c.EapRouterClient.Handle(context.Background(), in)
}

func newTestEapClient(t *testing.T, addr string) testEapServiceClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithBackoffMaxDelay(10*time.Second), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Client dial error: %v", err)
	}
	return testEapServiceClient{eap_protos.NewEapRouterClient(conn)}
}

// TestEapAkaConcurent tests EAP AKA Provider
func TestEapAkaConcurent(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.SwxProxy
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

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP)
	eap_protos.RegisterEapRouterServer(rtrSrv.GrpcServer, &testEapRouter{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	client := newTestEapClient(t, rtrLis.Addr().String())
	done := make(chan error)
	go eap_test.Auth(t, client, eap_test.IMSI1, 50, done)
	go eap_test.Auth(t, client, eap_test.IMSI2, 47, done)
	eap_test.Auth(t, client, eap_test.IMSI1, 43, nil)
	<-done
	<-done // wait for test 1 & 2 to complete
}

func TestEAPAkaWrongPlmnId(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.NoUseSwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnID6}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP)
	eap_protos.RegisterEapRouterServer(rtrSrv.GrpcServer, &testEapRouter{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	client := newTestEapClient(t, rtrLis.Addr().String())

	tst := eap_test.Units[eap_test.IMSI1]
	eapCtx := &eap_protos.EapContext{SessionId: eap.CreateSessionId()}
	_, err = client.Handle(&eap_protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
	if err == nil {
		t.Fatalf("Expected Error Handling Filtered PLMN ID")
	}
}

func TestEAPAkaPlmnId5(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	var service eap_test.SwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnID6, plmnID5}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}

	servicer.SetChallengeTimeout(time.Millisecond * 10)
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP)
	eap_protos.RegisterEapRouterServer(rtrSrv.GrpcServer, &testEapRouter{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	client := newTestEapClient(t, rtrLis.Addr().String())

	tst := eap_test.Units[eap_test.IMSI1]
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
	var service eap_test.SwxProxy
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	eapSrv, eapLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP_AKA)
	servicer, err := servicers.NewEapAkaService(&mconfig.EapAkaConfig{PlmnIds: []string{wrongPlmnID6, plmnID6}})
	if err != nil {
		t.Fatalf("failed to create EAP AKA Service: %v", err)
		return
	}
	servicer.SetChallengeTimeout(time.Millisecond * 10)
	eap_protos.RegisterEapServiceServer(eapSrv.GrpcServer, servicer)
	go eapSrv.RunTest(eapLis)

	rtrSrv, rtrLis := test_utils.NewTestService(t, registry.ModuleName, registry.EAP)
	eap_protos.RegisterEapRouterServer(rtrSrv.GrpcServer, &testEapRouter{supportedMethods: eap_client.SupportedTypes()})
	go rtrSrv.RunTest(rtrLis)

	client := newTestEapClient(t, rtrLis.Addr().String())

	tst := eap_test.Units[eap_test.IMSI1]
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
