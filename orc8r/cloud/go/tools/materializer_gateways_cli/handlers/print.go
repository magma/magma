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

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/tools/commands"
)

func init() {
	cmd := Commands.Add(
		"print",
		"Print the materialized gateway views for a network",
		PrintGatewayViews,
	)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: envdir /var/opt/magma/envdir %s print [OPTIONS]\n", os.Args[0])
		f.PrintDefaults()
	}
	f.StringVar(&networkIDArg, "network", "", "Network to print materialized gateway views for")
}

func PrintGatewayViews(cmd *commands.Command, args []string) int {
	views, err := store.GetGatewayViewsForNetwork(networkIDArg)
	if err != nil {
		glog.Errorf("Error getting gateway views for network %s: %s", networkIDArg, err)
		return 1
	}
	marshalled, err := json.MarshalIndent(views, "", "  ")
	if err != nil {
		glog.Errorf("Error marshaling gateway views to JSON: %s", err)
		return 1
	}
	fmt.Println(string(marshalled))
	fmt.Println()
	return 0
}
