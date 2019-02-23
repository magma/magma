/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/tools/commands"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

var whitelist string
var networkBlacklist string

func init() {
	verifyCmd := Commands.Add(
		"verify",
		"Verify stored mconfigs against built mconfigs",
		VerifyStoredMconfigsForNetworks,
	)

	verifyFlags := verifyCmd.Flags()
	verifyFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: %s verify [OPTIONS]\n", os.Args[0])
		verifyFlags.PrintDefaults()
	}
	verifyFlags.StringVar(
		&networkIdsList, "networks", "",
		"Comma-separated list of specific networks to verify mconfigs for. If unspecified or blank, all networks will be selected.",
	)
	verifyFlags.StringVar(
		&networkBlacklist, "blacklistedNetworks", "",
		"Comma-separated list of networks to skip verification for.",
	)
	verifyFlags.StringVar(
		&whitelist, "whitelist", "",
		"Comma-separated list of mconfig keys to whitelist when verifying. If unspecified or blank, all fields will be compared.",
	)
}

func VerifyStoredMconfigsForNetworks(cmd *commands.Command, args []string) int {
	gatewayIdsByNetworkId, err := loadGatewayIdsForNetworks()
	if err != nil {
		glog.Errorf("Could not load gateway IDs: %s", err)
		return 1
	}
	gatewayIdsByNetworkId = applyNetworkBlacklist(gatewayIdsByNetworkId)

	storedMconfigsByNetwork, err := loadStoredMconfigsForNetworks(gatewayIdsByNetworkId)
	if err != nil {
		glog.Error(err)
		return 1
	}
	builtMconfigsByNetwork, err := buildMconfigsForNetworks(gatewayIdsByNetworkId)
	if err != nil {
		glog.Error(err)
		return 1
	}

	err = compareStoredAndBuiltMconfigs(storedMconfigsByNetwork, builtMconfigsByNetwork, getWhitelist())
	if err != nil {
		glog.Errorf("Error occurred while comparing mconfigs: %s", err)
		return 1
	}

	fmt.Println("All mconfigs are consistent")
	return 0
}

func applyNetworkBlacklist(gatewayIdsByNetworkId map[string][]string) map[string][]string {
	if len(networkBlacklist) == 0 {
		return gatewayIdsByNetworkId
	}
	blacklist := strings.Split(networkBlacklist, ",")
	for _, blacklistedNetwork := range blacklist {
		delete(gatewayIdsByNetworkId, blacklistedNetwork)
	}
	return gatewayIdsByNetworkId
}

func compareStoredAndBuiltMconfigs(
	storedMconfigs map[string]NetworkStoredMconfigs,
	builtMconfigs map[string]NetworkBuiltMconfigs,
	whitelistedFields []string,
) error {
	for networkId, networkMconfigs := range storedMconfigs {
		for gatewayId, cfg := range networkMconfigs {
			builtMconfig, ok := builtMconfigs[networkId][gatewayId]
			if !ok {
				return fmt.Errorf("Built mconfigs do not contain an entry for gateway %s in network %s", gatewayId, networkId)
			}
			if !areMconfigsEqual(cfg.Mconfig, builtMconfig, whitelistedFields) {
				fmt.Printf("Built and stored mconfigs not equal for gateway %s in network %s\n", gatewayId, networkId)

				prettyStoredMconfig, err := getPrettyPrintableMconfig(networkId, gatewayId, cfg.Mconfig, cfg.Offset)
				if err != nil {
					return fmt.Errorf("Could not pretty print stored mconfig: %s", err)
				}
				marshaledStoredMconfig, err := json.MarshalIndent(prettyStoredMconfig, "", "  ")
				if err != nil {
					return fmt.Errorf("Could not marshal stored mconfig: %s", err)
				}

				prettyBuiltMconfig, err := getPrettyPrintableMconfig(networkId, gatewayId, builtMconfig, -1)
				if err != nil {
					return fmt.Errorf("Could not pretty print built mconfig: %s", err)
				}
				marshaledBuiltMconfig, err := json.MarshalIndent(prettyBuiltMconfig, "", "  ")
				if err != nil {
					return fmt.Errorf("Could not marshal built mconfig: %s", err)
				}

				fmt.Println("Stored Mconfig:")
				fmt.Println(string(marshaledStoredMconfig))
				fmt.Println("Built Mconfig:")
				fmt.Println(string(marshaledBuiltMconfig))

				fmt.Println()
				return errors.New("Inconsistency detected, see printout above")
			}
		}
	}

	return nil
}

func getWhitelist() []string {
	if len(whitelist) == 0 {
		return []string{}
	}
	return strings.Split(whitelist, ",")
}

func areMconfigsEqual(
	storedMconfig *protos.GatewayConfigs,
	builtMconfig *protos.GatewayConfigs,
	whitelistedFields []string,
) bool {
	// Delete any non-whitelisted fields so we can just use an equals compare
	if len(whitelistedFields) > 0 {
		stripNonWhitelistedFields(storedMconfig, whitelistedFields)
		stripNonWhitelistedFields(builtMconfig, whitelistedFields)
	}
	return proto.Equal(storedMconfig, builtMconfig)
}

func stripNonWhitelistedFields(out *protos.GatewayConfigs, whitelistedFields []string) {
	whitelistSet := stringArrayToSet(whitelistedFields)
	for key := range out.ConfigsByKey {
		if _, ok := whitelistSet[key]; !ok {
			delete(out.ConfigsByKey, key)
		}
	}
}

func stringArrayToSet(in []string) map[string]struct{} {
	ret := make(map[string]struct{}, len(in))
	for _, elt := range in {
		ret[elt] = struct{}{}
	}
	return ret
}
