/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"os"

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/tools/commands"
)

func init() {
	cmd := Commands.Add(
		"delete",
		"Delete the materialized views for the specified networks",
		DeleteNetworks,
	)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: envdir /var/opt/magma/envdir %s delete [OPTIONS]\n", os.Args[0])
		f.PrintDefaults()
	}
	f.StringVar(&networkIDsList, "networks", "", "Networks to delete - will delete all networks if not specified")
}

func DeleteNetworks(cmd *commands.Command, args []string) int {
	networkIDs, err := getNetworkIDs()
	if err != nil {
		glog.Errorf("Error getting network IDs to delete: %s", err)
		return 1
	}
	for _, networkID := range networkIDs {
		err = deleteNetworkViews(networkID)
		if err != nil {
			glog.Errorf("Error deleting views for network %s: %s", networkID, err)
			return 1
		}
	}
	return 0
}

func deleteNetworkViews(networkID string) error {
	views, err := store.GetGatewayViewsForNetwork(networkID)
	if err != nil {
		return fmt.Errorf("Error getting list of views for network: %s", err)
	}
	gatewayIDs := make([]string, 0, len(views))
	for gatewayID := range views {
		gatewayIDs = append(gatewayIDs, gatewayID)
	}
	return store.DeleteGatewayViews(networkID, gatewayIDs)
}
