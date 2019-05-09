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
	"sync"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

// StateValue includes reported operational states and additional info
type StateValue struct {
	// ID of the entity reporting the state (hwID, cert serial number, etc)
	ReporterID string
	// Checkin Time
	Time uint64
	// Cert expiration Time
	CertExpirationTime int64
	ReportedValue      []byte
}

// StateID contains the identifying information of a state
type StateID struct {
	Type     string
	DeviceID string
}

// Global clientconn that can be reused for this service
var connSingleton = (*grpc.ClientConn)(nil)
var connGuard = sync.Mutex{}

func getStateClient() (protos.StateServiceClient, error) {
	if connSingleton == nil {
		// Reading the conn optimistically to avoid unnecessary overhead
		connGuard.Lock()
		if connSingleton == nil {
			conn, err := registry.GetConnection(ServiceName)
			if err != nil {
				initErr := errors.NewInitError(err, ServiceName)
				glog.Error(initErr)
				connGuard.Unlock()
				return nil, initErr
			}
			connSingleton = conn
		}
		connGuard.Unlock()
	}
	return protos.NewStateServiceClient(connSingleton), nil
}

// GetState returns the state specified by the networkID, typeVal, and hwID
func GetState(networkID string, typeVal string, hwID string) (*StateValue, error) {
	stateValue := &StateValue{}
	client, err := getStateClient()
	if err != nil {
		return nil, err
	}

	stateID := &protos.StateID{
		Type:     typeVal,
		DeviceID: hwID,
	}

	ret, err := client.GetStates(
		context.Background(),
		&protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       []*protos.StateID{stateID},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(ret.States) == 0 {
		return nil, errors.ErrNotFound
	}
	return stateValue, json.Unmarshal(ret.States[0].Value, stateValue)
}

// GetStates returns a map of states specified by the networkID and a list of type and key
func GetStates(networkID string, stateIDs []StateID) (map[StateID]StateValue, error) {
	client, err := getStateClient()
	if err != nil {
		return nil, err
	}

	res, err := client.GetStates(
		context.Background(), &protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       toProtosStateIDs(stateIDs),
		},
	)
	if err != nil {
		return nil, err
	}

	idToValue := map[StateID]StateValue{}
	for _, state := range res.States {
		stateID := StateID{Type: state.Type, DeviceID: state.DeviceID}
		stateValue := StateValue{}
		json.Unmarshal(state.Value, &stateValue)
		idToValue[stateID] = stateValue
	}
	return idToValue, nil
}

// DeleteStates deletes states specified by the networkID and a list of type and key
func DeleteStates(networkID string, stateIDs []StateID) error {
	client, err := getStateClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteStates(
		context.Background(),
		&protos.DeleteStatesRequest{
			NetworkID: networkID,
			Ids:       toProtosStateIDs(stateIDs),
		},
	)
	return err
}

func toProtosStateIDs(stateIDs []StateID) []*protos.StateID {
	ids := []*protos.StateID{}
	for _, state := range stateIDs {
		ids = append(ids, &protos.StateID{Type: state.Type, DeviceID: state.DeviceID})
	}
	return ids
}
