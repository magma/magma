/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package service303

import (
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func getGWClient(service gateway_registry.GwServiceType, hwId string) (protos.Service303Client, *grpc.ClientConn, context.Context, error) {
	conn, ctx, err := gateway_registry.GetGatewayConnection(service, hwId)
	if err != nil {
		errMsg := fmt.Sprintf("service303 gwClient initialization error: %s", err)
		glog.Error(errMsg)
		return nil, nil, nil, errors.New(errMsg)
	}
	return protos.NewService303Client(conn), conn, ctx, nil

}

func GWService303GetServiceInfo(service gateway_registry.GwServiceType, hwId string) (*protos.ServiceInfo, error) {
	client, conn, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return client.GetServiceInfo(ctx, new(protos.Void))
}

func GWService303GetMetrics(service gateway_registry.GwServiceType, hwId string) (*protos.MetricsContainer, error) {
	client, conn, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return client.GetMetrics(ctx, new(protos.Void))
}

func GWService303StopService(service gateway_registry.GwServiceType, hwId string) error {
	client, conn, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.StopService(ctx, new(protos.Void))
	return err
}

func GWService303SetLogLevel(service gateway_registry.GwServiceType, hwId string, in *protos.LogLevelMessage) error {
	client, conn, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.SetLogLevel(ctx, in)
	return err
}

func GWService303SetLogVerbosity(service gateway_registry.GwServiceType, hwId string, in *protos.LogVerbosity) error {
	client, conn, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.SetLogVerbosity(ctx, in)
	return err
}
