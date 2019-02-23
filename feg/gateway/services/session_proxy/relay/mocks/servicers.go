/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mocks

import (
	"log"
	"net"
	"testing"
	"time"

	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// Create and start a mock SessionProxyResponder as a service. This returns the
// mock object that is registered as the servicer to configure in test code and
// an instance of a `MockCloudRegistry` bound to the server address allocated
// to the mock local session manager.
func StartMockSessionProxyResponder(t *testing.T) (*SessionProxyResponderServer, *MockCloudRegistry) {
	lis, err := net.Listen("tcp", "")
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	sm := &SessionProxyResponderServer{}
	protos.RegisterSessionProxyResponderServer(grpcServer, sm)
	serverStarted := make(chan struct{})
	go func() {
		log.Printf("Starting server")
		serverStarted <- struct{}{}
		grpcServer.Serve(lis)
	}()
	<-serverStarted
	time.Sleep(time.Millisecond)

	return sm, &MockCloudRegistry{ServerAddr: lis.Addr().String()}
}
