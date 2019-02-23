/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/tools/commands"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

type StoredMconfigPrettyPrinted struct {
	NetworkId string
	GatewayId string
	Mconfig   map[string]proto.Message
	Offset    int64
}

func init() {
	cmd := Commands.Add(
		"print",
		"Print the computed mconfigs for a network",
		PrettyPrintMconfigs,
	)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage; %s print [OPTIONS]\n", os.Args[0])
		f.PrintDefaults()
	}
	f.StringVar(&networkIdArg, "network", "", "Network to print mconfigs for")
}

func PrettyPrintMconfigs(cmd *commands.Command, args []string) int {
	gatewayIds, err := magmad.ListGateways(networkIdArg)
	if err != nil {
		glog.Errorf("Could not list gateways for network %s: %s", networkIdArg, err)
		return 1
	}
	mconfigs, err := store.GetMconfigs(networkIdArg, gatewayIds)
	if err != nil {
		glog.Errorf("Error loading mconfigs for network %s: %s", networkIdArg, err)
		return 1
	}

	for gwId, storedMconfig := range mconfigs {
		prettyStoredMconfig, err := getPrettyPrintableMconfigsFromStoredMconfigs(storedMconfig)
		if err != nil {
			glog.Errorf("Error converting stored mconfig to pretty print version: %s", err)
			return 1
		}

		marshaledStoredMconfig, err := json.MarshalIndent(prettyStoredMconfig, "", "  ")
		if err != nil {
			glog.Errorf("Error marshaling stored mconfig: %s", err)
			return 1
		}

		fmt.Printf("Mconfig for gateway %s:\n", gwId)
		fmt.Println(string(marshaledStoredMconfig))
		fmt.Println()
	}
	return 0
}

func getPrettyPrintableMconfigsFromStoredMconfigs(stored *storage.StoredMconfig) (*StoredMconfigPrettyPrinted, error) {
	return getPrettyPrintableMconfig(stored.NetworkId, stored.GatewayId, stored.Mconfig, stored.Offset)
}

func getPrettyPrintableMconfig(networkId string, gatewayId string, in *protos.GatewayConfigs, offset int64) (*StoredMconfigPrettyPrinted, error) {
	retMap := map[string]proto.Message{}

	for configKey, cfgAny := range in.ConfigsByKey {
		dAny := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(cfgAny, dAny); err != nil {
			return nil, err
		}
		retMap[configKey] = dAny
	}

	return &StoredMconfigPrettyPrinted{
		NetworkId: networkId,
		GatewayId: gatewayId,
		Mconfig:   retMap,
		Offset:    offset,
	}, nil
}
