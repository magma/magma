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
	cmdReboot := &cobra.Command{
		Use:   "reboot --network=<network-id> --gateway=<gateway-id>",
		Short: "reboot gateway device",
		Run:   rebootCmd,
	}

	rootCmd.AddCommand(cmdReboot)
}

func rebootCmd(cmd *cobra.Command, args []string) {
	err := magmad.GatewayReboot(networkId, gatewayId)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}
