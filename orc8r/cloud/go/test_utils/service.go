/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"net"
	"testing"

	"magma/orc8r/cloud/go/service"
	"magma/orc8r/lib/go/registry"
)

// Creates & Initializes test magma service on a dynamically selected available local port.
// Returns the newly created service and net.Listener, it was registered with.
func NewTestService(t *testing.T, moduleName string, serviceType string) (*service.Service, net.Listener) {
	// Create the server socket for gRPC
	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("failed to create listener: %s", err)
	}

	addr, err := net.ResolveTCPAddr("tcp", lis.Addr().String())
	if err != nil {
		t.Fatalf("failed to resolve TCP address: %s", err)
	}
	registry.AddService(registry.ServiceLocation{Name: serviceType, Host: "localhost", Port: addr.Port})

	// Create the service
	srv, err := service.NewTestOrchestratorService(t, moduleName, serviceType)
	if err != nil {
		t.Fatalf("Error creating service: %s", err)
	}

	return srv, lis
}
