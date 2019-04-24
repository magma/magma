/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"errors"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/serde"
	stateservice "magma/orc8r/cloud/go/services/state"
)

// ValidateGetStatesRequest checks that all required fields exist
func ValidateGetStatesRequest(req *protos.GetStatesRequest) error {
	if err := checkNonEmptyInput(req.GetNetworkID(), req.GetIds()); err != nil {
		return err
	}
	return nil
}

// ValidateReportStatesRequest checks that all required fields exist
func ValidateReportStatesRequest(req *protos.ReportStatesRequest) error {
	return validateStates(req)
}

// ValidateDeleteStatesRequest checks that all required fields exist
func ValidateDeleteStatesRequest(req *protos.DeleteStatesRequest) error {
	if err := checkNonEmptyInput(req.GetNetworkID(), req.GetIds()); err != nil {
		return err
	}
	return nil
}

func validateStates(req *protos.ReportStatesRequest) error {
	states := req.GetStates()
	if states == nil || len(states) == 0 {
		return errors.New("States value must be specified and non-empty")
	}
	for _, state := range states {
		_, err := serde.Deserialize(stateservice.SerdeDomain, state.GetType(), state.GetValue())
		if err != nil {
			return err
		}
	}
	return nil
}

func checkNonEmptyInput(networkID string, ids []*protos.StateID) error {
	if len(networkID) == 0 {
		return errors.New("Network ID must be specified")
	}
	if ids == nil || len(ids) == 0 {
		return errors.New("States value must be specified and non-empty")
	}
	return nil
}
