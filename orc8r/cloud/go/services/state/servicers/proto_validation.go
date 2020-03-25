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

	"magma/orc8r/cloud/go/serde"
	stateservice "magma/orc8r/cloud/go/services/state"
	"magma/orc8r/lib/go/protos"

	"github.com/thoas/go-funk"
)

// ValidateGetStatesRequest checks that all required fields exist
func ValidateGetStatesRequest(req *protos.GetStatesRequest) error {
	if err := enforceNetworkID(req.NetworkID); err != nil {
		return err
	}
	if funk.IsEmpty(req.Ids) && funk.IsEmpty(req.TypeFilter) && funk.IsEmpty(req.IdFilter) {
		return errors.New("at least one filter criteria must be specified in the request")
	}
	return nil
}

// ValidateDeleteStatesRequest checks that all required fields exist
func ValidateDeleteStatesRequest(req *protos.DeleteStatesRequest) error {
	if err := enforceNetworkID(req.NetworkID); err != nil {
		return err
	}
	if funk.IsEmpty(req.Ids) {
		return errors.New("States value must be specified and non-empty")
	}
	return nil
}

// ValidateSyncStatesRequest checks that all required fields exist
func ValidateSyncStatesRequest(req *protos.SyncStatesRequest) error {
	if req.GetStates() == nil || len(req.GetStates()) == 0 {
		return errors.New("States value must be specified and non-empty")
	}
	return nil
}

// PartitionStatesBySerializability checks that each state is deserializable.
// If a state is not deserializable, we will send back the states type, key, and error.
func PartitionStatesBySerializability(req *protos.ReportStatesRequest) ([]*protos.State, []*protos.IDAndError, error) {
	validatedStates := []*protos.State{}
	invalidStates := []*protos.IDAndError{}

	states := req.GetStates()
	if states == nil || len(states) == 0 {
		return nil, nil, errors.New("States value must be specified and non-empty")
	}
	for _, state := range states {
		model, err := serde.Deserialize(stateservice.SerdeDomain, state.GetType(), state.GetValue())
		if err != nil {
			stateAndError := &protos.IDAndError{
				Type:     state.Type,
				DeviceID: state.DeviceID,
				Error:    err.Error(), // deserialization error
			}
			invalidStates = append(invalidStates, stateAndError)
		} else {
			if err := model.(serde.ValidateableBinaryConvertible).ValidateModel(); err != nil {
				stateAndError := &protos.IDAndError{
					Type:     state.Type,
					DeviceID: state.DeviceID,
					Error:    err.Error(), // validation error
				}
				invalidStates = append(invalidStates, stateAndError)
				continue
			}

			validatedStates = append(validatedStates, state)
		}
	}
	return validatedStates, invalidStates, nil
}

func enforceNetworkID(networkID string) error {
	if len(networkID) == 0 {
		return errors.New("Network ID must be specified")
	}
	return nil
}
