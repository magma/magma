/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package gateway_health provides a client for using the cloud health service from
// federated gateways. This allows the gateway health manager to send health updates
// without knowing about the RPC implementation.
package gateway_health

import (
	"context"
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/orc8r/lib/go/errors"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

// getHealthClient is a utility function to get an RPC connection to the
// cloud Health service from the feg
func getHealthClient(cloudRegistry registry.CloudRegistry) (protos.HealthClient, *grpc.ClientConn, error) {
	if cloudRegistry == nil {
		return nil, nil, fmt.Errorf("Nil cloud registry provided")
	}
	conn, err := cloudRegistry.GetConnection("HEALTH")
	if err != nil {
		initErr := errors.NewInitError(err, "HEALTH")
		glog.Error(initErr)
		return nil, nil, initErr
	}
	return protos.NewHealthClient(conn), conn, nil
}

// UpdateHealth sends a health update using a HealthRequest to the cloud and returns
// back a health response and any potential error that occurred
func UpdateHealth(cloudReg registry.CloudRegistry, req *protos.HealthRequest) (*protos.HealthResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil HealthRequest")
	}
	client, conn, err := getHealthClient(cloudReg)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.UpdateHealth(context.Background(), req)
}
