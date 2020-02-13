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
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/spf13/cobra"
)

func init() {
	cmdGenericCommand := &cobra.Command{
		Use:   "generic_command <command> <params>",
		Short: "Execute generic command on gateway",
		Args:  cobra.ExactArgs(2),
		Run:   genericCommandCmd,
	}

	rootCmd.AddCommand(cmdGenericCommand)
}

func genericCommandCmd(cmd *cobra.Command, args []string) {
	paramsStruct := structpb.Struct{}
	err := jsonpb.UnmarshalString(args[1], &paramsStruct)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	genericCommandParams := protos.GenericCommandParams{
		Command: args[0],
		Params:  &paramsStruct,
	}

	response, err := magmad.GatewayGenericCommand(networkId, gatewayId, &genericCommandParams)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", response)
}
