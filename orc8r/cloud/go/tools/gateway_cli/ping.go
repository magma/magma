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

	"magma/orc8r/cloud/go/services/magmad"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var packets int32

func init() {
	cmdPing := &cobra.Command{
		Use:   "ping <hosts>... [--packets=<packets>]",
		Short: "ping hosts from gateway",
		Args:  cobra.MinimumNArgs(1),
		Run:   pingCmd,
	}

	cmdPing.Flags().Int32Var(&packets, "packets", 4, "number of packets to send for each ping")
	rootCmd.AddCommand(cmdPing)
}

func pingCmd(cmd *cobra.Command, args []string) {
	response, err := magmad.GatewayPing(networkId, gatewayId, packets, args)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", response)
}
