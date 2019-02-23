/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package xservice

import (
	"errors"
	"fmt"

	magmaerrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/config"
	storage2 "magma/orc8r/cloud/go/services/config/storage"
	"magma/orc8r/cloud/go/services/magmad"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

type xserviceStorage struct{}

func NewCrossServiceGatewayViewsStorage() storage.GatewayViewStorage {
	return &xserviceStorage{}
}

type gatewayConfigs map[storage2.TypeAndKey]interface{}

func (*xserviceStorage) InitTables() error {
	return nil
}

func (ss *xserviceStorage) GetGatewayViewsForNetwork(networkID string) (map[string]*storage.GatewayState, error) {
	gatewayIDs, err := magmad.ListGateways(networkID)
	if err != nil {
		return map[string]*storage.GatewayState{}, fmt.Errorf("Error loading gateway IDs for network view: %s", err)
	}
	return ss.GetGatewayViews(networkID, gatewayIDs)
}

func (*xserviceStorage) GetGatewayViews(networkID string, gatewayIDs []string) (map[string]*storage.GatewayState, error) {
	ret := make(map[string]*storage.GatewayState, len(gatewayIDs))
	for _, gatewayID := range gatewayIDs {
		state, err := loadGatewayView(networkID, gatewayID)
		if err != nil {
			return map[string]*storage.GatewayState{}, fmt.Errorf("Error loading gateway %s view: %s", gatewayID, err)
		}
		ret[gatewayID] = state
	}
	return ret, nil
}

func (*xserviceStorage) UpdateOrCreateGatewayViews(networkID string, updates map[string]*storage.GatewayUpdateParams) error {
	return errors.New("xservice gateway views storage is read only!")
}

func (*xserviceStorage) DeleteGatewayViews(networkID string, gatewayIDs []string) error {
	return errors.New("xservice gateway views storage is read only!")
}

func loadGatewayView(networkID string, gatewayID string) (*storage.GatewayState, error) {
	record, err := magmad.FindGatewayRecord(networkID, gatewayID)
	if err != nil {
		return nil, fmt.Errorf("Error loading record: %s", err)
	}
	configs, err := config.GetConfigsByKey(networkID, gatewayID)
	if err != nil {
		return nil, fmt.Errorf("Error loading configs: %s", err)
	}
	status, err := loadStatus(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	return createGatewayView(gatewayID, record, configs, status), nil
}

func loadStatus(networkID string, gatewayID string) (*protos.GatewayStatus, error) {
	status, err := checkind.GetStatus(networkID, gatewayID)
	if err == magmaerrors.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error loading status: %s", err)
	}
	return status, nil
}

func createGatewayView(gatewayID string, record *magmadprotos.AccessGatewayRecord, configs gatewayConfigs, status *protos.GatewayStatus) *storage.GatewayState {
	ret := &storage.GatewayState{
		GatewayID: gatewayID,
		Status:    status,
		Record:    record,
		Config:    make(map[string]interface{}, len(configs)),
	}
	for typeAndKey, configVal := range configs {
		ret.Config[typeAndKey.Type] = configVal
	}
	return ret
}
