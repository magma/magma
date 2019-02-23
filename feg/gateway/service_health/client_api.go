/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package service_health encapsulates service functionality related to health
// that service303 services can extend themselves with
package service_health

import (
	"context"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/orc8r/cloud/go/errors"
	orcprotos "magma/orc8r/cloud/go/protos"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

// getClient is a utility function to get an RPC connection to
// ServiceHealth
func getClient(service string) (protos.ServiceHealthClient, *grpc.ClientConn, error) {
	conn, err := registry.GetConnection(service)
	if err != nil {
		initErr := errors.NewInitError(err, service)
		glog.Error(initErr)
		return nil, nil, initErr
	}
	return protos.NewServiceHealthClient(conn), conn, nil
}

// Disable disables service functionality for the period of time
// specified in the DisableMessage for the service provided
func Disable(service string, req *protos.DisableMessage) error {
	client, conn, err := getClient(service)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.Disable(context.Background(), req)
	return err
}

// Enable enables service functionality for the service provided
func Enable(service string) error {
	client, conn, err := getClient(service)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.Enable(context.Background(), &orcprotos.Void{})
	return err
}

// GetHealthStatus returns a HealthStatus object that indicates the current health of
// the service provided
func GetHealthStatus(service string) (*protos.HealthStatus, error) {
	client, conn, err := getClient(service)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	healthStatus, err := client.GetHealthStatus(context.Background(), &orcprotos.Void{})
	return healthStatus, err
}
