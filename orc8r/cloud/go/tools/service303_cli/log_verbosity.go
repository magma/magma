/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"magma/orc8r/cloud/go/service/client"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/cloud/go/services/dispatcher/gw_client_apis/service303"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	cmdLogVerbosity := &cobra.Command{
		Use:   "log_verbosity <verbosity> <service> [--gateway-service (--hwid=<hardware-id> | --network=<network-id> --gateway=<gateway-id>)]",
		Short: "Set log verbosity",
		Args:  validateLogVerbosityArgs,
		Run:   logVerbosityCmd,
	}

	rootCmd.AddCommand(cmdLogVerbosity)
}

func validateLogVerbosityArgs(cmd *cobra.Command, args []string) error {
	if err := validateGlobalFlags(); err != nil {
		return err
	}
	if err := setHwIdFlag(); err != nil {
		return err
	}
	if len(args) != 2 {
		return errors.New("requires 2 args")
	}
	if verb, err := strconv.Atoi(args[0]); err != nil || verb < 0 {
		return fmt.Errorf("log_verbosity is not valid")
	}
	if !isGatewayServiceQuery && !isValidService(args[1], services) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[1], services)
	}
	if isGatewayServiceQuery && !isValidGwService(gateway_registry.GwServiceType(args[1]), gwServices) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[1], gwServices)
	}

	return nil
}

func logVerbosityCmd(cmd *cobra.Command, args []string) {
	verb, err := strconv.Atoi(args[0])
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	err = setLogVerbosity(args[1], verb)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}

func setLogVerbosity(service string, verbosity int) error {
	err := setLogVerbosityOrGwLogVerbosity(service, verbosity)
	if err != nil {
		return fmt.Errorf("Failed to SetLogVerbosity for %s: %s", service, err)
	}
	return nil
}

func setLogVerbosityOrGwLogVerbosity(service string, verbosity int) error {
	if isGatewayServiceQuery {
		return service303.GWService303SetLogVerbosity(gateway_registry.GwServiceType(service), hardwareID, &protos.LogVerbosity{Verbosity: int32(verbosity)})
	} else {
		return client.Service303SetLogVerbosity(service, &protos.LogVerbosity{Verbosity: int32(verbosity)})
	}
}
