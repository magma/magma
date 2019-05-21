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

	"magma/feg/cloud/go/services/feg_relay/utils"
	"magma/orc8r/cloud/go/protos"

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

func getHwIDFromIMSI(imsi string) (string, error) {
	// directoryd prefixes imsi with "IMSI" when updating the location
	if !strings.HasPrefix(imsi, "IMSI") {
		imsi = fmt.Sprintf("IMSI%s", imsi)
	}
	hwId, err := directoryd.GetHardwareIdByIMSI(imsi)
	if err != nil {
		return hwId, err
	}
	glog.V(2).Infof("IMSI to send is %v\n", imsi)
	return hwId, nil
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
	hwId, err := getHwIDFromIMSI(imsi)
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
