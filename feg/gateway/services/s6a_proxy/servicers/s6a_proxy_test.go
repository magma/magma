/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/s6a_proxy/servicers/test"
	orcprotos "magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const TEST_LOOPS = 33

// TestS6aProxyService creates a mock S6a Diameter server, S6a S6a Proxy service
// and runs tests using GRPC client: GRPC Client <--> GRPC Server <--> S6a SCTP Diameter Server
func TestS6aProxyService(t *testing.T) {

	diamAddr := fmt.Sprintf("127.0.0.1:%d", 29000+rand.Intn(1900))
	err := test.StartTestS6aServer("sctp", diamAddr)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Started S6a Server")

	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
		return
	}
	srvCfg := &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:     diamAddr, // use "192.168.60.145:3868" to send diam messages to OAI HSS VM
		Protocol: "sctp"},  // tcp/sctp
	}

	clientCfg := &diameter.DiameterClientConfig{
		Host:  "magma-oai.openair4G.eur", // diameter host
		Realm: "openair4G.eur",           // diameter realm
	}
	s := grpc.NewServer()
	service, err := servicers.NewS6aProxy(clientCfg, srvCfg)
	if err != nil {
		t.Fatalf("failed to create S6aProxy: %v", err)
		return

	}
	protos.RegisterS6AProxyServer(s, service)
	protos.RegisterServiceHealthServer(s, service)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
			return
		}
	}()
	addr := lis.Addr()
	t.Logf("Started S6a GRPC Proxy on %s", addr.String())

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()

	c := protos.NewS6AProxyClient(conn)
	hs := protos.NewServiceHealthClient(conn)
	req := &protos.AuthenticationInformationRequest{
		UserName:                   test.TEST_IMSI,
		VisitedPlmn:                []byte(test.TEST_PLMN_ID),
		NumRequestedEutranVectors:  3,
		ImmediateResponsePreferred: true,
	}
	complChan := make(chan error, TEST_LOOPS+1)
	testLoopF := func() {
		// AIR
		r, err := c.AuthenticationInformation(context.Background(), req)
		if err != nil {
			t.Fatalf("GRPC AIR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC AIA: %#+v", *r)
		if r.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected AIA Error Code: %d", r.ErrorCode)
		}
		if len(r.EutranVectors) != 3 {
			t.Errorf("Unexpected Number of EutranVectors: %d, Expected: 3", len(r.EutranVectors))
		}
		ulReq := &protos.UpdateLocationRequest{
			UserName:           test.TEST_IMSI,
			VisitedPlmn:        []byte(test.TEST_PLMN_ID),
			SkipSubscriberData: false,
			InitialAttach:      true,
		}
		// ULR
		ulResp, err := c.UpdateLocation(context.Background(), ulReq)
		if err != nil {
			t.Fatalf("GRPC ULR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC ULA: %#+v", *ulResp)
		if ulResp.ErrorCode != protos.ErrorCode_UNDEFINED {
			t.Errorf("Unexpected ULA Error Code: %d", ulResp.ErrorCode)
		}

		puReq := &protos.PurgeUERequest{
			UserName: test.TEST_IMSI,
		}
		// PUR
		puResp, err := c.PurgeUE(context.Background(), puReq)
		if err != nil {
			t.Fatalf("GRPC PUR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC PUA: %#+v", *puResp)
		if puResp.ErrorCode != protos.ErrorCode_SUCCESS {
			t.Errorf("Unexpected PUA Error Code: %d", puResp.ErrorCode)
		}
		complChan <- nil
	}
	go testLoopF()
	select {
	case err := <-complChan:
		if err != nil {
			t.Fatal(err)
			return
		}
	case <-time.After(time.Second):
		t.Fatal("Timed out")
		return
	}

	for round := 0; round < TEST_LOOPS; round++ {
		go testLoopF()
	}
	for round := 0; round < TEST_LOOPS; round++ {
		testErr := <-complChan
		if testErr != nil {
			t.Fatal(err)
			return
		}
	}

	// Test Disabling / Enabling Connections

	disableReq := &protos.DisableMessage{
		DisablePeriodSecs: 10,
	}

	// Disable connections
	_, err = hs.Disable(context.Background(), disableReq)
	if err != nil {
		t.Fatalf("GRPC Disable Error: %v", err)
		return
	}

	// AIR should fail
	_, err = c.AuthenticationInformation(context.Background(), req)
	if err == nil {
		t.Errorf("AIR Succeeded, but should have failed due to disabled connections")
	}

	// Enable connections
	_, err = hs.Enable(context.Background(), &orcprotos.Void{})
	if err != nil {
		t.Fatalf("GRPC Enable Error: %v", err)
		return
	}

	// AIR should pass now
	airResp, err := c.AuthenticationInformation(context.Background(), req)
	if err != nil {
		t.Fatalf("GRPC AIR Error: %v", err)
		return
	}
	t.Logf("GRPC AIA: %#+v", *airResp)
	if airResp.ErrorCode != protos.ErrorCode_UNDEFINED {
		t.Errorf("Unexpected AIA Error Code: %d", airResp.ErrorCode)
	}
}
