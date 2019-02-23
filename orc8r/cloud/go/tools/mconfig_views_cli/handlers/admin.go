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

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/tools/commands"

	"github.com/golang/glog"
)

var reseedKeepOffsets bool

func init() {
	delCmd := Commands.Add(
		"delete",
		"Delete stored mconfigs for networks",
		DeleteStoredMconfigsForNetworks,
	)

	delFlags := delCmd.Flags()
	delFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: %s delete [OPTIONS]\n", os.Args[0])
		delFlags.PrintDefaults()
	}
	delFlags.StringVar(
		&networkIdsList, "networks", "",
		"Comma-separated list of specific networks to delete mconfigs for. If unspecified or blank, all networks will be modified",
	)

	reseedCmd := Commands.Add(
		"reseed",
		"Reseed stored mconfigs for networks from mconfig builders",
		ReSeedStoredMconfigsWithBuiltMconfigs,
	)

	reseedFlags := reseedCmd.Flags()
	reseedFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: %s reseed [OPTIONS]\n", os.Args[0])
		reseedFlags.PrintDefaults()
	}

	reseedFlags.StringVar(
		&networkIdsList, "networks", "",
		"Comma-separated list of specific networks to reseed. If unspecified or blank, all networks will be modified,",
	)
	reseedFlags.BoolVar(
		&reseedKeepOffsets, "keepOffsets", true,
		"Maintain offsets for existing views when reseeding.",
	)
}

// Delete command will delete all stored mconfigs for networks
func DeleteStoredMconfigsForNetworks(cmd *commands.Command, args []string) int {
	networkIds, err := getNetworkIds()
	if err != nil {
		glog.Errorf("Could not load network IDs: %s", err)
		return 1
	}

	for _, networkId := range networkIds {
		err := db.DeleteTable(datastore.GetTableName(networkId, "mconfig_views"))
		if err != nil {
			glog.Errorf("Error deleting stored mconfigs for network %s: %s", networkId, err)
			return 1
		}
	}

	glog.V(2).Info("Successfully deleted stored mconfigs")
	return 0
}

// Reseed command will seed network mconfig views with values calculated from
// mconfig builders
func ReSeedStoredMconfigsWithBuiltMconfigs(cmd *commands.Command, args []string) int {
	gwIdsByNwId, err := loadGatewayIdsForNetworks()
	if err != nil {
		glog.Errorf("Could not load gateway IDs for networks: %s", err)
		return 1
	}
	mconfigsByNetwork, err := buildMconfigsForNetworks(gwIdsByNwId)
	if err != nil {
		glog.Errorf("Could not build mconfigs: %s", err)
		return 1
	}
	storedMconfigsByNetwork, err := loadStoredMconfigsForNetworks(gwIdsByNwId)
	if err != nil {
		glog.Errorf("Could not load existing stored mconfigs: %s", err)
		return 1
	}

	for nwId, mconfigsByGwId := range mconfigsByNetwork {
		err := store.CreateOrUpdateMconfigs(nwId, builtMconfigsToCreateCriteria(mconfigsByGwId, storedMconfigsByNetwork[nwId]))
		if err != nil {
			glog.Errorf("Error while storing built mconfigs for network %s: %s", nwId, err)
			return 1
		}
	}
	return 0
}

func builtMconfigsToCreateCriteria(mconfigsByGwId NetworkBuiltMconfigs, storedMconfigs NetworkStoredMconfigs) []*storage.MconfigUpdateCriteria {
	ret := make([]*storage.MconfigUpdateCriteria, 0, len(mconfigsByGwId))
	for gwId, mcfg := range mconfigsByGwId {
		ret = append(ret, &storage.MconfigUpdateCriteria{
			GatewayId:  gwId,
			NewMconfig: mcfg,
			Offset:     getOffset(gwId, storedMconfigs),
		})
	}
	return ret
}

func getOffset(gatewayId string, storedMconfigs NetworkStoredMconfigs) int64 {
	if !reseedKeepOffsets {
		return -1
	}

	storedMconfig, ok := storedMconfigs[gatewayId]
	if !ok {
		return -1
	}
	return storedMconfig.Offset
}
