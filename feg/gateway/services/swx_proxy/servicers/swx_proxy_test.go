/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"context"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/swx_proxy/servicers/test"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	TEST_LOOPS = 33
)

// TestSwxProxyService creates a mock Swx Diameter server, Swx Proxy service
// and runs tests using GRPC client: GRPC Client <--> GRPC Server <--> Swx SCTP Diameter Server
func TestSwxProxyService(t *testing.T) {
	serverAddr, err := test.StartTestSwxServer("sctp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Started Swx Server at %s", serverAddr)

	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
		return
	}
	srvCfg := &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:     serverAddr,
		Protocol: "sctp"}, // tcp/sctp
	}

	clientCfg := &diameter.DiameterClientConfig{
		Host:  "magma-oai.openair4G.eur", // diameter host
		Realm: "openair4G.eur",           // diameter realm
	}
	s := grpc.NewServer()
	service, err := servicers.NewSwxProxy(clientCfg, srvCfg)
	if err != nil {
		t.Fatalf("failed to create SwxProxy: %v", err)
		return

	}
	protos.RegisterSwxProxyServer(s, service)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
			return
		}
	}()
	addr := lis.Addr()
	t.Logf("Started Swx GRPC Proxy on %s", addr.String())

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("GRPC connect error: %v", err)
		return
	}
	defer conn.Close()

	client := protos.NewSwxProxyClient(conn)
	complChan := make(chan error, TEST_LOOPS+1)

	// Happy path
	testHappyPath := func(reqId uint32) {
		userName := test.BASE_IMSI + strconv.Itoa(int(reqId))
		numVectors := (reqId % 5) + 1 // arbitrary number 1-5
		authReq := &protos.AuthenticationRequest{
			UserName:             userName,
			SipNumAuthVectors:    numVectors,
			AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		}
		// Authentication Request - MAR
		authRes, err := client.Authenticate(context.Background(), authReq)
		if err != nil {
			t.Fatalf("GRPC MAR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC MAA: %#+v", *authRes)
		assert.Equal(t, userName, authRes.GetUserName())
		if len(authRes.SipAuthVectors) != int(numVectors) {
			t.Errorf("Unexpected Number of SIPAuthVectors: %d, Expected: %d", len(authRes.SipAuthVectors), numVectors)
		}
		for i, v := range authRes.SipAuthVectors {
			assert.Equal(t, protos.AuthenticationScheme_EAP_AKA, v.GetAuthenticationScheme())
			assert.Equal(t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(i+14))), v.GetRandAutn())
			assert.Equal(t, []byte(test.DefaultSIPAuthorization), v.GetXres())
			assert.Equal(t, []byte(test.DefaultCK), v.GetConfidentialityKey())
			assert.Equal(t, []byte(test.DefaultIK), v.GetIntegrityKey())
		}
		// Registration request - SAR
		regReq := &protos.RegistrationRequest{
			UserName: userName,
		}
		regRes, err := client.Register(context.Background(), regReq)
		// Only must verify that request was successful (no error) to ensure user
		// is registered
		if err != nil {
			t.Fatalf("GRPC SAR Error: %v", err)
			complChan <- err
			return
		}
		t.Logf("GRPC SAA: %#+v", *regRes)
		complChan <- nil
	}
	go testHappyPath(uint32(rand.Intn(100)))
	select {
	case err := <-complChan:
		if err != nil {
			t.Fatal(err)
			return
		}
	case <-time.After(time.Second * 5):
		t.Fatal("Timed out")
		return
	}

	// Multi-threaded test ensures session-id logic handling works
	for round := 0; round < TEST_LOOPS; round++ {
		go testHappyPath(uint32(round))
	}
	for round := 0; round < TEST_LOOPS; round++ {
		testErr := <-complChan
		if testErr != nil {
			t.Fatal(err)
			return
		}
	}

	// Test Auth Error Handling
	_, err = client.Authenticate(context.Background(), nil)
	assert.EqualError(t, err, "rpc error: code = Internal desc = grpc: error while marshaling: proto: Marshal called with nil")

	emptyAuthReq := &protos.AuthenticationRequest{}
	_, err = client.Authenticate(context.Background(), emptyAuthReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Empty user-name provided in authentication request")

	badNumVectorsReq := &protos.AuthenticationRequest{
		UserName:             "10111011000110",
		AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		SipNumAuthVectors:    0,
	}

	_, err = client.Authenticate(context.Background(), badNumVectorsReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = SIPNumAuthVectors in authentication request must be greater than 0")

	badUserNameReq := &protos.AuthenticationRequest{
		UserName:             "1234567890123456",
		AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		SipNumAuthVectors:    0,
	}
	_, err = client.Authenticate(context.Background(), badUserNameReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = SIPNumAuthVectors in authentication request must be greater than 0")

	// Test Register Error Handling
	_, err = client.Register(context.Background(), nil)
	assert.EqualError(t, err, "rpc error: code = Internal desc = grpc: error while marshaling: proto: Marshal called with nil")

	emptyResReq := &protos.RegistrationRequest{}
	_, err = client.Register(context.Background(), emptyResReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Empty user-name provided in registration request")

	badRegReq := &protos.RegistrationRequest{
		UserName: "1234567890123456",
	}
	_, err = client.Register(context.Background(), badRegReq)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Provided username 1234567890123456 is greater than 15 digits")
}
