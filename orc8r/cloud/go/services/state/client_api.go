/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package state

import (
	"context"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	state_types "magma/orc8r/cloud/go/services/state/types"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

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
func GetState(networkID string, typeVal string, hwID string) (state_types.State, error) {
	client, err := GetStateClient()
	if err != nil {
		return state_types.State{}, err
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
		return state_types.State{}, err
	}
	if len(res.States) == 0 {
		return state_types.State{}, merrors.ErrNotFound
	}
	return state_types.MakeState(res.States[0])
}

// GetStates returns a map of states specified by the networkID and a list of type and key
func GetStates(networkID string, stateIDs []state_types.ID) (state_types.StatesByID, error) {
	if len(stateIDs) == 0 {
		return state_types.StatesByID{}, nil
	}

	client, err := GetStateClient()
	if err != nil {
		return nil, err
	}

	res, err := client.GetStates(
		context.Background(), &protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       makeProtoIDs(stateIDs),
		},
	)
	if err != nil {
		return nil, err
	}
	return state_types.MakeStatesByID(res.States)
}

// SearchStates returns all states matching the filter arguments.
// typeFilter and keyFilter are both OR clauses, and the final predicate
// applied to the search will be the AND of both filters.
// e.g.: ["t1", "t2"], ["k1", "k2"] => (t1 OR t2) AND (k1 OR k2)
func SearchStates(networkID string, typeFilter []string, keyFilter []string) (state_types.StatesByID, error) {
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
	return state_types.MakeStatesByID(res.States)
}

// DeleteStates deletes states specified by the networkID and a list of
// type and key.
func DeleteStates(networkID string, stateIDs []state_types.ID) error {
	client, err := GetStateClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteStates(
		context.Background(),
		&protos.DeleteStatesRequest{
			NetworkID: networkID,
			Ids:       makeProtoIDs(stateIDs),
		},
	)
	return err
}

// GetGatewayStatus returns the status for an indicated gateway.
func GetGatewayStatus(networkID string, deviceID string) (*models.GatewayStatus, error) {
	st, err := GetState(networkID, orc8r.GatewayStateType, deviceID)
	if err != nil {
		return nil, err
	}
	if st.ReportedState == nil {
		return nil, merrors.ErrNotFound
	}
	return fillInGatewayStatusState(st), nil
}

// GetGatewayStatuses returns the status for indicated gateways, keyed by
// device ID.
func GetGatewayStatuses(networkID string, deviceIDs []string) (map[string]*models.GatewayStatus, error) {
	stateIDs := funk.Map(deviceIDs, func(id string) state_types.ID {
		return state_types.ID{Type: orc8r.GatewayStateType, DeviceID: id}
	}).([]state_types.ID)
	res, err := GetStates(networkID, stateIDs)
	if err != nil {
		return map[string]*models.GatewayStatus{}, err
	}

	ret := make(map[string]*models.GatewayStatus, len(res))
	for stateID, st := range res {
		ret[stateID.DeviceID] = fillInGatewayStatusState(st)
	}
	return ret, nil
}

func makeProtoIDs(stateIDs []state_types.ID) []*protos.StateID {
	var ids []*protos.StateID
	for _, st := range stateIDs {
		ids = append(ids, &protos.StateID{Type: st.Type, DeviceID: st.DeviceID})
	}
	return ids
}

func fillInGatewayStatusState(st state_types.State) *models.GatewayStatus {
	if st.ReportedState == nil {
		return nil
	}
	gwStatus := st.ReportedState.(*models.GatewayStatus)
	gwStatus.CheckinTime = st.TimeMs
	gwStatus.CertExpirationTime = st.CertExpirationTime
	gwStatus.HardwareID = st.ReporterID
	return gwStatus
}
