/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/tools/commands"
)

var keepOffsets bool

func init() {
	cmd := Commands.Add(
		"reseed",
		"Reseed the views for the specified networks",
		ReseedNetworks,
	)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: envdir /var/opt/magma/envdir %s reseed [OPTIONS]\n", os.Args[0])
		f.PrintDefaults()
	}
	f.StringVar(&networkIDsList, "networks", "", "Networks to reseed - will reseed all networks if not specified")
	f.BoolVar(&keepOffsets, "keepOffsets", true, "Specifies if reseed should keep the stored Kafka offsets")
}

func ReseedNetworks(cmd *commands.Command, args []string) int {
	networkIDs, err := getNetworkIDs()
	if err != nil {
		glog.Errorf("Error getting network IDs: %s", err)
		return 1
	}
	for _, networkID := range networkIDs {
		err := reseedNetwork(networkID)
		if err != nil {
			glog.Errorf("Error reseeding network %s: %s", networkID, err)
			return 1
		}
	}
	return 0
}

func reseedNetwork(networkID string) error {
	reseedStates, err := loadGatewayStatesForNetwork(networkID)
	if err != nil {
		return fmt.Errorf("Error loading true gateway states: %s", err)
	}

	// List the current views for this network
	existingViews, err := store.GetGatewayViewsForNetwork(networkID)
	if err != nil {
		return fmt.Errorf("Error loading old gateway views")
	}
	existingViewGatewayIds := make([]string, 0, len(existingViews))

	// Keep offsets if necessary
	for gatewayID, view := range existingViews {
		existingViewGatewayIds = append(existingViewGatewayIds, gatewayID)
		if state, ok := reseedStates[gatewayID]; keepOffsets && ok {
			state.Offset = view.Offset
		}
	}

	// Delete all views in the network
	err = store.DeleteGatewayViews(networkID, existingViewGatewayIds)
	if err != nil {
		return fmt.Errorf("Error deleting old gateway views: %s", err)
	}

	// When writing the reseed, explicitly lowercase the network ID
	return store.UpdateOrCreateGatewayViews(strings.ToLower(networkID), convertStateToUpdateParams(reseedStates))
}

func convertStateToUpdateParams(states map[string]*storage.GatewayState) map[string]*storage.GatewayUpdateParams {
	updateParams := make(map[string]*storage.GatewayUpdateParams)
	for gatewayID, state := range states {
		updateParams[gatewayID] = &storage.GatewayUpdateParams{
			NewConfig: state.Config,
			NewStatus: state.Status,
			NewRecord: state.Record,
			Offset:    state.Offset,
		}
	}
	return updateParams
}
