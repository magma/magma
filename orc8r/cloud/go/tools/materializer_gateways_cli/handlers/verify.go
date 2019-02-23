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
		"verify",
		"Verify the materialized views for the specified networks",
		VerifyGatewayViews,
	)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: envdir /var/opt/magma/envdir %s verify [OPTIONS]\n", os.Args[0])
		f.PrintDefaults()
	}
	f.StringVar(&networkIDsList, "networks", "", "Networks to verify - will verify all networks if not specified")
}

func VerifyGatewayViews(cmd *commands.Command, args []string) int {
	networkIDs, err := getNetworkIDs()
	if err != nil {
		glog.Errorf("Error getting network IDs for verification")
		return 1
	}
	for _, networkID := range networkIDs {
		err := verifyGatewayViewsforNetwork(networkID)
		if err != nil {
			glog.Errorf("%s", err)
			return 1
		}
	}
	return 0
}

func verifyGatewayViewsforNetwork(networkID string) error {
	states, err := loadGatewayStatesForNetwork(networkID)
	if err != nil {
		return fmt.Errorf("Error loading true gateway states for network %s: %s", networkID, err)
	}
	views, err := store.GetGatewayViewsForNetwork(networkID)
	if err != nil {
		return err
	}
	for gatewayID, state := range states {
		view, ok := views[gatewayID]
		if !ok {
			return fmt.Errorf("View does not exist for gateway %s in network %s", gatewayID, networkID)
		}
		err = verifyConfigs(state.Config, view.Config)
		if err != nil {
			return fmt.Errorf("Config mismatch in view for gateway %s in network %s: %s", gatewayID, networkID, err)
		}
		err = compareFormatted(state.Status, view.Status)
		if err != nil {
			return fmt.Errorf("Status mismatch in view for gateway %s in network %s: %s", gatewayID, networkID, err)
		}
		err = compareFormatted(state.Record, view.Record)
		if err != nil {
			return fmt.Errorf("Record mismatch in view for gateway %s in network %s: %s", gatewayID, networkID, err)
		}
	}
	return nil
}

func verifyConfigs(expected, actual map[string]interface{}) error {
	for configType, actualObj := range actual {
		expectedObj, ok := expected[configType]
		if !ok {
			return fmt.Errorf("Unexpected config type found in view: %s", configType)
		}
		err := compareFormatted(expectedObj, actualObj)
		if err != nil {
			return fmt.Errorf("Config %s comparison error: %s", configType, err)
		}
	}
	for configType := range expected {
		if _, ok := actual[configType]; !ok {
			return fmt.Errorf("Expected config type %s not found in view", configType)
		}
	}
	return nil
}

func compareFormatted(expected, actual interface{}) error {
	expectedBytes, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		return err
	}
	expectedString := string(expectedBytes)
	actualBytes, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		return err
	}
	actualString := string(actualBytes)
	if expectedString != actualString {
		return fmt.Errorf("Objects do not match.\nExpected:\n%s\nActual:\n%s", expectedString, actualString)
	}
	return nil
}
