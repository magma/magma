/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"log"
	"strings"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/tools/commands"
)

var Commands = new(commands.Map)

var db datastore.Api
var store storage.MconfigStorage

// Destination variable for networks flag
var networkIdsList string

// Destination variable for single network flag
var networkIdArg string

func init() {
	newDB, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		log.Fatalf("Could not initialize datastore: %s", err)
	}
	db = newDB

	store = storage.NewDatastoreMconfigStorage(db)
}

// Get network IDs from the value of the networks flag
func getNetworkIds() ([]string, error) {
	if len(networkIdsList) > 0 {
		return strings.Split(networkIdsList, ","), nil
	}
	// If no network IDs specified, load all of them
	return magmad.ListNetworks()
}

// Get gateway IDs for networks specified by the networks flag
func loadGatewayIdsForNetworks() (map[string][]string, error) {
	nwIds, err := getNetworkIds()
	if err != nil {
		return map[string][]string{}, err
	}

	ret := make(map[string][]string, len(nwIds))
	for _, networkId := range nwIds {
		gwIds, err := magmad.ListGateways(networkId)
		if err != nil {
			return map[string][]string{}, err
		}
		ret[networkId] = gwIds
	}
	return ret, nil
}
