/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package servicesrs implements various relay RPCs to relay messages from FeG to Gateways via Controller
package servicers

import (
	"fmt"
	"strings"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/plugin/models"
	"magma/feg/cloud/go/services/feg_relay/utils"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

// FegToGwRelayServer is a server serving requests from FeG to Access Gateway
type FegToGwRelayServer struct {
}

// NewFegToGwRelayServer creates a new FegToGwRelayServer
func NewFegToGwRelayServer() (*FegToGwRelayServer, error) {
	return &FegToGwRelayServer{}, nil
}

func getHwIDFromIMSI(ctx context.Context, imsi string) (string, error) {
	gw := protos.GetClientGateway(ctx)
	// directoryd prefixes imsi with "IMSI" when updating the location
	if !strings.HasPrefix(imsi, "IMSI") {
		imsi = fmt.Sprintf("IMSI%s", imsi)
	}
	servedIds, err := getFegServedIds(gw.GetNetworkId())
	if err != nil {
		return "", err
	}
	for _, nid := range servedIds {
		hwId, err := directoryd.GetHardwareIdByIMSI(imsi, nid)
		if err == nil && len(hwId) != 0 {
			glog.V(2).Infof("IMSI to send is %v\n", imsi)
			return hwId, nil
		}
	}
	return "", fmt.Errorf("could not find gateway location for IMSI: %s", imsi)
}

func validateFegContext(ctx context.Context) error {
	fegId := protos.GetClientGateway(ctx)
	if fegId == nil {
		ctxMetadata, _ := metadata.FromIncomingContext(ctx)
		errorStr := fmt.Sprintf(
			"Failed to get Identity of calling Federated Gateway from CTX Metadata: %+v",
			ctxMetadata,
		)
		glog.Error(errorStr)
		return fmt.Errorf(errorStr)
	}
	if !fegId.Registered() {
		return fmt.Errorf("federated gateway not registered")
	}
	return nil
}

func getGWSGSServiceConnCtx(ctx context.Context, imsi string) (*grpc.ClientConn, context.Context, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, nil, err
	}
	hwId, err := getHwIDFromIMSI(ctx, imsi)
	if err != nil {
		errorStr := fmt.Sprintf(
			"unable to get HwID from IMSI %v. err: %v\n",
			imsi,
			err,
		)
		glog.Error(errorStr)
		return nil, nil, fmt.Errorf(errorStr)
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(
		gateway_registry.GwSgsService, hwId)
	if err != nil {
		errorStr := fmt.Sprintf(
			"unable to get connection to the gateway: %v",
			err,
		)
		return nil, nil, fmt.Errorf(errorStr)
	}
	return conn, ctx, nil
}

func getAllGWSGSServiceConnCtx(ctx context.Context) ([]*grpc.ClientConn, []context.Context, error) {
	var connList []*grpc.ClientConn
	var ctxList []context.Context

	hwIds, err := utils.GetAllGatewayIDs(ctx)
	if err != nil {
		return connList, ctxList, err
	}
	for _, hwId := range hwIds {
		conn, ctx, err := gateway_registry.GetGatewayConnection(
			gateway_registry.GwSgsService,
			hwId,
		)
		if err != nil {
			return connList, ctxList, err
		}
		connList = append(connList, conn)
		ctxList = append(ctxList, ctx)
	}

	return connList, ctxList, nil
}

func getFegServedIds(networkId string) ([]string, error) {
	if len(networkId) == 0 {
		return []string{}, fmt.Errorf("Empty networkID provided.")
	}
	fegCfg, err := configurator.LoadNetworkConfig(networkId, feg.FegNetworkType)
	if err != nil || fegCfg == nil {
		return []string{}, fmt.Errorf("unable to retrieve config for federation network: %s", networkId)
	}
	networkFegConfigs, ok := fegCfg.(*models.NetworkFederationConfigs)
	if !ok || networkFegConfigs == nil {
		return []string{}, fmt.Errorf("invalid federation network config found for network: %s", networkId)
	}
	return networkFegConfigs.ServedNetworkIds, nil
}
