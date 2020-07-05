/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"fmt"
	"net"
	"testing"

	cloud_service "magma/orc8r/cloud/go/service"
	"magma/orc8r/lib/go/registry"
	platform_service "magma/orc8r/lib/go/service"
)

// NewTestService creates & initializes test magma service on a dynamically
// selected available local port.
// Returns the newly created service and net.Listener, it was registered with.
func NewTestService(t *testing.T, moduleName string, serviceType string) (*platform_service.Service, net.Listener) {
	srvPort, lis, err := getOpenPort(t)
	if err != nil {
		t.Fatal(err)
	}
	registry.AddService(registry.ServiceLocation{Name: serviceType, Host: "localhost", Port: srvPort})

	// Create the service
	srv, err := cloud_service.NewTestService(t, moduleName, serviceType)
	if err != nil {
		t.Fatalf("Error creating service: %s", err)
	}
	return srv, lis
}

// NewOrchestratorTestService creates & initializes a test orchestrator service
// on a dynamically selected available local port for the GRPC server and HTTP
// echo server. Returns the newly created service and net.Listener, it was
// registered with.
func NewOrchestratorTestService(t *testing.T, moduleName string, serviceType string) (*cloud_service.OrchestratorService, net.Listener) {
	srvPort, lis, err := getOpenPort(t)
	if err != nil {
		t.Fatal(err)
	}
	echoPort, echoLis, err := getOpenPort(t)
	if err != nil {
		t.Fatal(err)
	}
	err = echoLis.Close()
	if err != nil {
		t.Fatal(err)
	}
	registry.AddService(registry.ServiceLocation{Name: serviceType, Host: "localhost", EchoPort: echoPort, Port: srvPort})

	// Create the service
	srv, err := cloud_service.NewTestOrchestratorService(t, moduleName, serviceType)
	if err != nil {
		t.Fatalf("Error creating service: %s", err)
	}
	return srv, lis
}

func getOpenPort(t *testing.T) (int, net.Listener, error) {
	lis, err := net.Listen("tcp", "")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create listener: %s", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", lis.Addr().String())
	if err != nil {
		return 0, nil, fmt.Errorf("failed to resolve TCP address: %s", err)
	}
	return addr.Port, lis, err
}
