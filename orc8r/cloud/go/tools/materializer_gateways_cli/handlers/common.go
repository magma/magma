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

	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo"
	"magma/orc8r/cloud/go/tools/commands"
)

var Commands = new(commands.Map)

var store storage.GatewayViewStorage
var networkIDArg string
var networkIDsList string

func init() {
	newStore, err := dynamo.GetInitializedDynamoStorage()
	if err != nil {
		log.Fatalf("Could not initialize materialized view storage: %s", err)
	}
	store = newStore
}

func getNetworkIDs() ([]string, error) {
	if networkIDsList == "" {
		return magmad.ListNetworks()
	}
	rawNetworkIDs := strings.Split(networkIDsList, ",")
	networkIDs := make([]string, len(rawNetworkIDs))
	for i, rawNetworkID := range rawNetworkIDs {
		networkIDs[i] = strings.TrimSpace(rawNetworkID)
	}
	return networkIDs, nil
}
