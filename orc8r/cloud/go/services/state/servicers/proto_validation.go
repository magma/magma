/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// ValidateGetStatesRequest checks that all required fields exist.
func ValidateGetStatesRequest(req *protos.GetStatesRequest) error {
	if !funk.IsEmpty(req.Ids) && funk.IsEmpty(req.NetworkID) {
		return errors.New("network ID must be non-empty for non-empty state IDs")
	}
	return nil
}

// ValidateDeleteStatesRequest checks that all required fields exist.
func ValidateDeleteStatesRequest(req *protos.DeleteStatesRequest) error {
	if err := enforceNetworkID(req.NetworkID); err != nil {
		return err
	}
	if funk.IsEmpty(req.Ids) {
		return errors.New("states value must be specified and non-empty")
	}
	return nil
}

// ValidateSyncStatesRequest checks that all required fields exist.
func ValidateSyncStatesRequest(req *protos.SyncStatesRequest) error {
	if req.GetStates() == nil || len(req.GetStates()) == 0 {
		return errors.New("states value must be specified and non-empty")
	}
	return nil
}

// PartitionStatesBySerializability checks that each state is deserializable.
// If a state is not deserializable, we will send back the states type, key, and error.
func PartitionStatesBySerializability(req *protos.ReportStatesRequest) ([]*protos.State, []*protos.IDAndError, error) {
	var validatedStates []*protos.State
	var invalidStates []*protos.IDAndError

	states := req.GetStates()
	if states == nil || len(states) == 0 {
		return nil, nil, errors.New("states value must be specified and non-empty")
	}
	for _, st := range states {
		model, err := serde.Deserialize(state.SerdeDomain, st.GetType(), st.GetValue())
		if err != nil {
			stateAndError := &protos.IDAndError{
				Type:     st.Type,
				DeviceID: st.DeviceID,
				Error:    err.Error(), // deserialization error
			}
			invalidStates = append(invalidStates, stateAndError)
		} else {
			if err := model.(serde.ValidateableBinaryConvertible).ValidateModel(); err != nil {
				stateAndError := &protos.IDAndError{
					Type:     st.Type,
					DeviceID: st.DeviceID,
					Error:    err.Error(), // validation error
				}
				invalidStates = append(invalidStates, stateAndError)
				continue
			}

			validatedStates = append(validatedStates, st)
		}
	}
	return validatedStates, invalidStates, nil
}

func enforceNetworkID(networkID string) error {
	if len(networkID) == 0 {
		return errors.New("network ID must be specified")
	}
	return nil
}
