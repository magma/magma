/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package state

import (
	"context"
	"encoding/json"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
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

// GetStateClient returns a client to the state service.
func GetStateClient() (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewStateServiceClient(conn), nil
}

// GetState returns the state specified by the networkID, typeVal, and hwID
func GetState(networkID string, typeVal string, hwID string) (State, error) {
	client, err := GetStateClient()
	if err != nil {
		return State{}, err
	}

	stateID := &protos.StateID{
		Type:     typeVal,
		DeviceID: hwID,
	}

	res, err := client.GetStates(
		context.Background(),
		&protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       []*protos.StateID{stateID},
		},
	)
	if err != nil {
		return State{}, err
	}
	if len(res.States) == 0 {
		return State{}, merrors.ErrNotFound
	}
	return makeState(res.States[0])
}

// GetStates returns a map of states specified by the networkID and a list of type and key
func GetStates(networkID string, stateIDs []ID) (StatesByID, error) {
	if len(stateIDs) == 0 {
		return StatesByID{}, nil
	}

	client, err := GetStateClient()
	if err != nil {
		return nil, err
	}

	res, err := client.GetStates(
		context.Background(), &protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       toProtoIDs(stateIDs),
		},
	)
	if err != nil {
		return nil, err
	}
	return MakeStatesByID(res.States)
}

// SearchStates returns all states matching the filter arguments.
// typeFilter and keyFilter are both OR clauses, and the final predicate
// applied to the search will be the AND of both filters.
// e.g.: ["t1", "t2"], ["k1", "k2"] => (t1 OR t2) AND (k1 OR k2)
func SearchStates(networkID string, typeFilter []string, keyFilter []string) (StatesByID, error) {
	client, err := GetStateClient()
	if err != nil {
		return nil, err
	}

	res, err := client.GetStates(context.Background(), &protos.GetStatesRequest{
		NetworkID:  networkID,
		TypeFilter: typeFilter,
		IdFilter:   keyFilter,
		LoadValues: true,
	})
	if err != nil {
		return nil, err
	}
	return MakeStatesByID(res.States)
}

// DeleteStates deletes states specified by the networkID and a list of
// type and key.
func DeleteStates(networkID string, stateIDs []ID) error {
	client, err := GetStateClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteStates(
		context.Background(),
		&protos.DeleteStatesRequest{
			NetworkID: networkID,
			Ids:       toProtoIDs(stateIDs),
		},
	)
	return err
}

// GetGatewayStatus returns the status for an indicated gateway.
func GetGatewayStatus(networkID string, deviceID string) (*models.GatewayStatus, error) {
	state, err := GetState(networkID, orc8r.GatewayStateType, deviceID)
	if err != nil {
		return nil, err
	}
	if state.ReportedState == nil {
		return nil, merrors.ErrNotFound
	}
	return fillInGatewayStatusState(state), nil
}

// GetGatewayStatuses returns the status for indicated gateways, keyed by
// device ID.
func GetGatewayStatuses(networkID string, deviceIDs []string) (map[string]*models.GatewayStatus, error) {
	stateIDs := funk.Map(deviceIDs, func(id string) ID {
		return ID{Type: orc8r.GatewayStateType, DeviceID: id}
	}).([]ID)
	res, err := GetStates(networkID, stateIDs)
	if err != nil {
		return map[string]*models.GatewayStatus{}, err
	}

	ret := make(map[string]*models.GatewayStatus, len(res))
	for stateID, state := range res {
		ret[stateID.DeviceID] = fillInGatewayStatusState(state)
	}
	return ret, nil
}

// MakeStatesByID converts state protos to state structs, keyed by state ID.
func MakeStatesByID(states []*protos.State) (StatesByID, error) {
	byID := StatesByID{}
	for _, p := range states {
		id := ID{Type: p.Type, DeviceID: p.DeviceID}
		state, err := makeState(p)
		if err != nil {
			return nil, err
		}
		byID[id] = state
	}
	return byID, nil
}

// makeState converts a state proto to a state structs.
func makeState(p *protos.State) (State, error) {
	// Recover state struct
	serialized := &SerializedStateWithMeta{}
	err := json.Unmarshal(p.Value, serialized)
	if err != nil {
		return State{}, errors.Wrap(err, "failed to unmarshal json-encoded state proto value")
	}

	// Recover reported state within state struct
	iReportedState, err := serde.Deserialize(SerdeDomain, p.Type, serialized.SerializedReportedState)
	state := State{
		ReporterID:         serialized.ReporterID,
		TimeMs:             serialized.TimeMs,
		CertExpirationTime: serialized.CertExpirationTime,
		ReportedState:      iReportedState,
		Type:               p.Type,
		Version:            p.Version,
	}

	return state, err
}

func fillInGatewayStatusState(state State) *models.GatewayStatus {
	if state.ReportedState == nil {
		return nil
	}
	gwStatus := state.ReportedState.(*models.GatewayStatus)
	gwStatus.CheckinTime = state.TimeMs
	gwStatus.CertExpirationTime = state.CertExpirationTime
	gwStatus.HardwareID = state.ReporterID
	return gwStatus
}

func toProtoIDs(stateIDs []ID) []*protos.StateID {
	var ids []*protos.StateID
	for _, state := range stateIDs {
		ids = append(ids, &protos.StateID{Type: state.Type, DeviceID: state.DeviceID})
	}
	return ids
}
