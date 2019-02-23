/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"os"

	"magma/orc8r/cloud/go/services/magmad"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	cmdRestartServices := &cobra.Command{
		Use:   "restart_services [<services>...] --network=<network-id> --gateway=<gateway-id>",
		Short: "restart gateway services. If no services specified, restart all services",
		Run:   restartServicesCmd,
	}

	rootCmd.AddCommand(cmdRestartServices)
}

func restartServicesCmd(cmd *cobra.Command, args []string) {
	err := magmad.GatewayRestartServices(networkId, gatewayId, args)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}
