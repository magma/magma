/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"os"

	"magma/orc8r/cloud/go/plugin"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gateway_cli",
	Short: "Gateway cli",
}

var networkId string
var gatewayId string

func main() {
	flag.Parse()
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})

	rootCmd.PersistentFlags().StringVar(&networkId, "network", "", "the network id")
	rootCmd.PersistentFlags().StringVar(&gatewayId, "gateway", "", "the gateway id")

	rootCmd.MarkPersistentFlagRequired("network")
	rootCmd.MarkPersistentFlagRequired("gateway")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
