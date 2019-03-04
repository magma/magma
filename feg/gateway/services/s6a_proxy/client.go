/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package s6a_proxy provides a thin client for using s6a proxy service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package s6a_proxy

import (
	"errors"
	"fmt"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type s6aProxyClient struct {
	lteprotos.S6AProxyClient
	fegprotos.ServiceHealthClient
}

// getS6aProxyClient is a utility function to get a RPC connection to the
// S6a Proxy service
func getS6aProxyClient() (*s6aProxyClient, *grpc.ClientConn, error) {
	conn, err := registry.GetConnection(registry.S6A_PROXY)
	if err != nil {
		errMsg := fmt.Sprintf("S6a Proxy client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, conn, errors.New(errMsg)
	}
	return &s6aProxyClient{
		lteprotos.NewS6AProxyClient(conn),
		fegprotos.NewServiceHealthClient(conn),
	}, conn, err
}

// AuthenticationInformation sends AIR over diameter connection,
// waits (blocks) for AIA & returns its RPC representation
func AuthenticationInformation(req *lteprotos.AuthenticationInformationRequest) (*lteprotos.AuthenticationInformationAnswer, error) {
	if req == nil {
		return nil, errors.New("Invalid AuthenticationInformationRequest")
	}
	cli, conn, err := getS6aProxyClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return cli.AuthenticationInformation(context.Background(), req)
}

// UpdateLocation sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULA & returns its RPC representation
func UpdateLocation(req *lteprotos.UpdateLocationRequest) (*lteprotos.UpdateLocationAnswer, error) {
	if req == nil {
		return nil, errors.New("Invalid UpdateLocation")
	}
	cli, conn, err := getS6aProxyClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return cli.UpdateLocation(context.Background(), req)
}

// PurgeUE sends PUR (Code 321) over diameter connection,
// waits (blocks) for PUA & returns its RPC representation
func PurgeUE(req *lteprotos.PurgeUERequest) (*lteprotos.PurgeUEAnswer, error) {
	if req == nil {
		return nil, errors.New("Invalid PurgeUE Request")
	}
	cli, conn, err := getS6aProxyClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return cli.PurgeUE(context.Background(), req)
}
