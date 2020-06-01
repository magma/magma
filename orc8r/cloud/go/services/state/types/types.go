/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// Package types contains the types and associated methods for the state service.
package types

import (
	"encoding/json"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
)

const (
	// SerdeDomain is a copy of the state service's serde domain, from
	// orc8r/cloud/go/services/state/doc.go.
	SerdeDomain = "state"
)

// StatesByID maps state IDs to state.
// A state and its ID collectively contains all information for a piece of state.
type StatesByID map[ID]State

// ID identifies a piece of state.
// A piece of state is uniquely identified by the triplet {network ID, device ID, type}.
type ID struct {
	// Type determines how the value is deserialized and validated on the cloud service side.
	Type string
	// DeviceID is the ID of the entity with which the state is associated (IMSI, serial number, etc).
	DeviceID string
}

// State includes reported operational state and additional info about the reporter.
type State struct {
	// ReportedState is the actual state reported by the device.
	ReportedState interface{}

	// Type determines how the reported state value is deserialized and validated on the cloud service side.
	Type string
	// Version is the reported version of the state.
	Version uint64
	// ReporterID is the hardware ID of the gateway which reported the state.
	ReporterID string
	// TimeMs is the time the state was received in milliseconds.
	TimeMs uint64
	// CertExpirationTime is the expiration time in milliseconds.
	CertExpirationTime int64
}

// IDsByNetwork are a set of state IDs, keyed by network ID.
type IDsByNetwork map[string][]ID

// SerializedStateWithMeta includes reported operational states and additional info
type SerializedStateWithMeta struct {
	ReporterID              string
	TimeMs                  uint64
	CertExpirationTime      int64
	SerializedReportedState []byte
	Version                 uint64
}

// MakeStatesByID converts state protos to state structs, keyed by state ID.
func MakeStatesByID(states []*protos.State) (StatesByID, error) {
	byID := StatesByID{}
	for _, p := range states {
		id := ID{Type: p.Type, DeviceID: p.DeviceID}
		st, err := MakeState(p)
		if err != nil {
			return nil, err
		}
		byID[id] = st
	}
	return byID, nil
}

// MakeState converts a state proto to a state structs.
func MakeState(p *protos.State) (State, error) {
	// Recover state struct
	serialized := &SerializedStateWithMeta{}
	err := json.Unmarshal(p.Value, serialized)
	if err != nil {
		return State{}, errors.Wrap(err, "failed to unmarshal json-encoded state proto value")
	}

	// Recover reported state within state struct
	iReportedState, err := serde.Deserialize(SerdeDomain, p.Type, serialized.SerializedReportedState)
	st := State{
		ReporterID:         serialized.ReporterID,
		TimeMs:             serialized.TimeMs,
		CertExpirationTime: serialized.CertExpirationTime,
		ReportedState:      iReportedState,
		Type:               p.Type,
		Version:            p.Version,
	}

	return st, err
}
