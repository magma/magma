/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package magmad provides functions for taking actions at connected gateways.
package magmad

import (
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

// getGWMagmadClient gets a GRPC client to the magmad service running on the gateway specified by (network ID, gateway ID).
// If gateway not found by configurator, returns ErrNotFound from magma/orc8r/lib/go/errors.
func getGWMagmadClient(networkID string, gatewayID string) (protos.MagmadClient, context.Context, error) {
	hwID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return nil, nil, err
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwMagmad, hwID)
	if err != nil {
		errMsg := fmt.Sprintf("gateway magmad client initialization error: %s", err)
		glog.Errorf(errMsg, err)
		return nil, nil, errors.New(errMsg)
	}
	return protos.NewMagmadClient(conn), ctx, nil
}

// GatewayReboot reboots a gateway.
// If gateway not registered, returns ErrNotFound from magma/orc8r/lib/go/errors.
func GatewayReboot(networkId string, gatewayId string) error {
	client, ctx, err := getGWMagmadClient(networkId, gatewayId)
	if err != nil {
		return err
	}
	_, err = client.Reboot(ctx, new(protos.Void))
	return err
}

// GatewayRestartServices restarts services at a gateway.
// If gateway not registered, returns ErrNotFound from magma/orc8r/lib/go/errors.
func GatewayRestartServices(networkId string, gatewayId string, services []string) error {
	client, ctx, err := getGWMagmadClient(networkId, gatewayId)
	if err != nil {
		return err
	}
	_, err = client.RestartServices(ctx, &protos.RestartServicesRequest{Services: services})
	return err
}

// GatewayPing sends pings from a gateway to a set of hosts.
// If gateway not registered, returns ErrNotFound from magma/orc8r/lib/go/errors.
func GatewayPing(networkId string, gatewayId string, packets int32, hosts []string) (*protos.NetworkTestResponse, error) {
	client, ctx, err := getGWMagmadClient(networkId, gatewayId)
	if err != nil {
		return nil, err
	}

	var pingParams []*protos.PingParams
	for _, host := range hosts {
		pingParams = append(pingParams, &protos.PingParams{HostOrIp: host, NumPackets: packets})
	}
	return client.RunNetworkTests(ctx, &protos.NetworkTestRequest{Pings: pingParams})
}

// GatewayGenericCommand runs a generic command at a gateway.
// If gateway not registered, returns ErrNotFound from magma/orc8r/lib/go/errors.
func GatewayGenericCommand(networkId string, gatewayId string, params *protos.GenericCommandParams) (*protos.GenericCommandResponse, error) {
	client, ctx, err := getGWMagmadClient(networkId, gatewayId)
	if err != nil {
		return nil, err
	}

	return client.GenericCommand(ctx, params)
}

// TailGatewayLogs
// If gateway not registered, returns ErrNotFound from magma/orc8r/lib/go/errors.
func TailGatewayLogs(networkId string, gatewayId string, service string) (protos.Magmad_TailLogsClient, error) {
	client, ctx, err := getGWMagmadClient(networkId, gatewayId)
	if err != nil {
		return nil, err
	}

	stream, err := client.TailLogs(ctx, &protos.TailLogsRequest{Service: service})
	if err != nil {
		return nil, err
	}

	return stream, nil
}
