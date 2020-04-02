/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package s6a_proxy

import (
	"errors"
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/feg/gateway/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func getCloudConn() (*grpc.ClientConn, error) {
	conn, err := registry.Get().GetCloudConnection(feg_relay.ServiceName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to establish connection to cloud FegToGwRelayClient: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return conn, nil
}

// GWS6AProxyCancelLocation forwards CLR to Controller
func GWS6AProxyCancelLocation(in *protos.CancelLocationRequest) (*protos.CancelLocationAnswer, error) {
	conn, err := getCloudConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := protos.NewS6AGatewayServiceClient(conn)
	return client.CancelLocation(context.Background(), in)
}

// GWS6AProxyReset forwards RSR to Controller
func GWS6AProxyReset(in *protos.ResetRequest) (*protos.ResetAnswer, error) {
	conn, err := getCloudConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := protos.NewS6AGatewayServiceClient(conn)
	return client.Reset(context.Background(), in)
}
