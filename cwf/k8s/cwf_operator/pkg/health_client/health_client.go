/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package health_client

import (
	"context"
	"fmt"
	"time"

	"magma/feg/cloud/go/protos"
	orc8rprotos "magma/orc8r/lib/go/protos"

	"google.golang.org/grpc"
)

const (
	GrpcMaxDelaySec   = 10
	GrpcMaxTimeoutSec = 10
)

// getClient is a utility function to get an RPC connection to
// the gateway's health service for the provided service address.
func getClient(serviceAddr string) (protos.ServiceHealthClient, error) {
	conn, err := getConnection(serviceAddr)
	if err != nil {
		return nil, err
	}
	client := protos.NewServiceHealthClient(conn)
	return client, nil
}

// GetHealthStatus calls the provided service address to obtain health
// status from a gateway.
func GetHealthStatus(address string) (*protos.HealthStatus, error) {
	client, err := getClient(address)
	if err != nil {
		return nil, err
	}
	return client.GetHealthStatus(context.Background(), &orc8rprotos.Void{})
}

// Enable calls the provided service address to enable gateway functionality
// after a standby gateway is promoted.
func Enable(address string) error {
	client, err := getClient(address)
	if err != nil {
		return err
	}
	_, err = client.Enable(context.Background(), &orc8rprotos.Void{})
	return err
}

// Disable calls the provided service address to disable gateway functionality
// after an active gateway is demoted.
func Disable(address string) error {
	req := &protos.DisableMessage{}
	client, err := getClient(address)
	if err != nil {
		return err
	}
	_, err = client.Disable(context.Background(), req)
	return err
}

// getConnection provides a gRPC connection to a service in the registry.
func getConnection(serviceAddr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GrpcMaxTimeoutSec*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, serviceAddr, grpc.WithBackoffMaxDelay(GrpcMaxDelaySec*time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Address: %s GRPC Dial error: %s", serviceAddr, err)
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return conn, nil
}
