/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package utils includes common helper functions used by FeG Rely components/services
package utils

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"magma/feg/cloud/go/services/controller/config"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
)

// GetAllGatewaysIDs returns all Gateways served by calling FeG,
// Calling FeG ID is set by Identity framework & retrieved from ctx
func GetAllGatewayIDs(ctx context.Context) ([]string, error) {
	var err error
	res := []string{}
	fegId := protos.GetClientGateway(ctx)
	if fegId == nil {
		ctxMetadata, _ := metadata.FromIncomingContext(ctx)
		return res, fmt.Errorf(
			"Failed to get Identity of calling Federated Gateway from CTX Metadata: %+v", ctxMetadata)
	}
	networkId := fegId.GetNetworkId()
	logicalId := fegId.GetLogicalId()
	if len(networkId) == 0 || len(logicalId) == 0 {
		return res, fmt.Errorf("Unregistered Federated Gateway: %s", fegId.String())
	}
	cfg, err := config.GetGatewayConfig(networkId, logicalId)
	if err != nil {
		return res, fmt.Errorf("Error getting Federated Gateway %s:%s configs: %v", networkId, logicalId, err)
	}
	// Find as many gateways as possible, don't exit on error, just return last error to the caller along with
	// the list of GWs found
	for _, network := range cfg.GetServedNetworkIds() {
		gateways, err := magmad.ListGateways(network)
		if err != nil {
			err = fmt.Errorf("List Network '%s' Gateways error: %v", network, err)
			continue
		}
		for _, gw := range gateways {
			record, err := magmad.FindGatewayRecord(network, gw)
			if err != nil {
				err = fmt.Errorf("Find Gateway Record Error: %v for Gateway %s:%s", err, network, gw)
				continue
			}
			hwId := record.GetHwId().GetId()
			if len(hwId) > 0 {
				res = append(res, hwId)
			} else {
				err = fmt.Errorf("Empty Harware ID for Gateway %s:%s", network, gw)
			}
		}
	}
	return res, err
}
