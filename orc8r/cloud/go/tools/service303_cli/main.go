/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"fmt"
	"os"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/cloud/go/services/magmad"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "service303_cli",
	Short: "Management CLI for Service303",
}

var services []string
var gwServices []gateway_registry.GwServiceType

var isGatewayServiceQuery bool
var hwId string
var networkId string
var gatewayId string

func main() {
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})

	services = registry.ListAllServices()
	gwServices = gateway_registry.ListAllGwServices()

	rootCmd.PersistentFlags().BoolVar(&isGatewayServiceQuery, "gateway-service", false, "query a gateway service")
	rootCmd.PersistentFlags().StringVar(&hwId, "hwid", "", "the hardware id of the gateway to send command to")
	rootCmd.PersistentFlags().StringVar(&networkId, "network", "", "the network id")
	rootCmd.PersistentFlags().StringVar(&gatewayId, "gateway", "", "the gateway id")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func validateGlobalFlags() error {
	if !isGatewayServiceQuery && hwId == "" && networkId == "" && gatewayId == "" {
		return nil
	}
	if isGatewayServiceQuery {
		if hwId != "" && networkId == "" && gatewayId == "" {
			return nil
		}
		if hwId == "" && networkId != "" && gatewayId != "" {
			return nil
		}
	}
	return fmt.Errorf("invalid flag combination")
}

func setHwIdFlag() error {
	if networkId == "" || gatewayId == "" {
		return nil
	}
	var err error
	hwId, err = getHwId(networkId, gatewayId)
	if err != nil {
		return err
	}
	return nil
}

func getHwId(networkId string, logicalId string) (string, error) {
	gwRecord, err := magmad.FindGatewayRecord(networkId, logicalId)
	if err != nil {
		return "", err
	}
	return gwRecord.HwId.Id, nil
}

func isValidService(service string, services []string) bool {
	for _, serv := range services {
		if serv == service {
			return true
		}
	}
	return false
}

func isValidGwService(service gateway_registry.GwServiceType, services []gateway_registry.GwServiceType) bool {
	for _, serv := range services {
		if serv == service {
			return true
		}
	}
	return false
}
