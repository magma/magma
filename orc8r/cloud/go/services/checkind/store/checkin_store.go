/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package store

import (
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind/scribe"
	"magma/orc8r/cloud/go/services/magmad"
)

const GatewaysStatusTableName string = "gwstatus"

var ErrNotFound = errors.New("Status not found")

type CheckinStore struct {
	store datastore.Api
}

// Validate checks if the store is properly initialized
func (s *CheckinStore) Validate() error {
	if s == nil {
		return fmt.Errorf("Nil CheckinStore")
	}
	if s.store == nil {
		return fmt.Errorf("Nil CheckinStore datastore")
	}
	return nil
}

// Create a new Checkin Store
func NewCheckinStore(ds datastore.Api) (*CheckinStore, error) {
	s := &CheckinStore{ds}
	return s, s.Validate()
}

func statusTable(networkId string) string {
	return datastore.GetTableName(networkId, GatewaysStatusTableName)
}

// Updates the given gateway status, the gateway is identified by its hardware
// ID and UpdateGatewayStatus uses Magmad to map the hardware ID to
// the gateway's network and logical IDs
func (s *CheckinStore) UpdateGatewayStatus(status *protos.GatewayStatus) error {
	if status == nil || status.Checkin == nil {
		return fmt.Errorf("Nil Gateway Status/Checkin Request")
	}
	networkId, err := magmad.FindGatewayNetworkId(status.Checkin.GatewayId)
	if err != nil {
		return fmt.Errorf("ID Lookup Error for Gateway '%s': %s",
			status.Checkin.GatewayId, err)
	}
	logicalId, err := magmad.FindGatewayId(networkId, status.Checkin.GatewayId)
	if err != nil {
		return err
	}
	return s.UpdateRegisteredGatewayStatus(networkId, logicalId, status)
}

// UpdateRegisteredGatewayStatus - updates the given registered gateway
// status, the gateway is identified by its network & logical IDs
func (s *CheckinStore) UpdateRegisteredGatewayStatus(networkId, logicalId string, status *protos.GatewayStatus) error {
	if status == nil {
		return fmt.Errorf("Nil Gateway Status")
	}
	if status.Checkin == nil {
		return fmt.Errorf("Nil Gateway Checkin Request")
	}
	marshaledStatus, err := protos.MarshalIntern(status)

	if err != nil {
		return fmt.Errorf(
			"Gateway Status Marshal Error: %s for GW: %s > %s",
			err, status.Checkin.GatewayId, logicalId,
		)
	}
	err = s.store.Put(statusTable(networkId), logicalId, marshaledStatus)
	if err != nil {
		return fmt.Errorf(
			"Gateway Status Write Error: %s for GW: %s > %s",
			err, status.Checkin.GatewayId, logicalId,
		)
	}
	// update checkin status successful, log status to Scribe
	go scribe.LogGatewayStatusToScribe(status, networkId, logicalId)
	return nil
}

// GetGatewayStatus returns the gateway status given its network and logical IDs
// GetGatewayStatus relies only on it's own DB table and does not use any
// external DBs or services
func (s *CheckinStore) GetGatewayStatus(req *protos.GatewayStatusRequest) (*protos.GatewayStatus, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil Gateway Status Request")
	}
	marshaledStatus, _, err := s.store.Get(statusTable(req.NetworkId), req.LogicalId)
	if err == datastore.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf(
			"Gateway Status Read Error: %s for network: %s, Gateway: %s",
			err, req.NetworkId, req.LogicalId,
		)
	}
	status := new(protos.GatewayStatus)
	err = protos.Unmarshal(marshaledStatus, status)
	return status, err
}

// DeleteGatewayStatus deletes the status of given gateway based on its network
// and logical IDs
// DeleteGatewayStatus relies only on it's own DB table and does not use any
// external DBs or services
func (s *CheckinStore) DeleteGatewayStatus(req *protos.GatewayStatusRequest) error {
	if req == nil {
		return fmt.Errorf("Nil Gateway Status Request")
	}
	return s.store.Delete(statusTable(req.NetworkId), req.LogicalId)
}

// DeleteNetworkTable deletes the status table for a given network.
// The table must be empty prior to call to DeleteNetworkTable
// DeleteNetworkTable relies only on it's own DB table and does not use any
// external DBs or services
func (s *CheckinStore) DeleteNetworkTable(networkId string) error {
	allKeys, err := s.store.ListKeys(statusTable(networkId))
	if err != nil {
		return fmt.Errorf("Error while checking if status table is empty: %s", err)
	}
	if len(allKeys) > 0 {
		return fmt.Errorf("Status table for network %s is not empty", networkId)
	}
	return s.store.DeleteTable(statusTable(networkId))
}

// List all logical gateway IDs for a given network.
func (s *CheckinStore) List(networkId string) ([]string, error) {
	return s.store.ListKeys(statusTable(networkId))
}
