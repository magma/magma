/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// Package wrappers provides semantic wrappers around the state service's
// client API.
package wrappers

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/errors"

	"github.com/thoas/go-funk"
)

// GetGatewayStatus returns the status for an indicated gateway.
func GetGatewayStatus(networkID string, deviceID string) (*models.GatewayStatus, error) {
	st, err := state.GetState(networkID, orc8r.GatewayStateType, deviceID)
	if err != nil {
		return nil, err
	}
	if st.ReportedState == nil {
		return nil, errors.ErrNotFound
	}
	return fillInGatewayStatusState(st), nil
}

// GetGatewayStatuses returns the status for indicated gateways, keyed by
// device ID.
func GetGatewayStatuses(networkID string, deviceIDs []string) (map[string]*models.GatewayStatus, error) {
	stateIDs := funk.Map(deviceIDs, func(id string) types.ID {
		return types.ID{Type: orc8r.GatewayStateType, DeviceID: id}
	}).([]types.ID)
	res, err := state.GetStates(networkID, stateIDs)
	if err != nil {
		return map[string]*models.GatewayStatus{}, err
	}

	ret := make(map[string]*models.GatewayStatus, len(res))
	for stateID, st := range res {
		ret[stateID.DeviceID] = fillInGatewayStatusState(st)
	}
	return ret, nil
}

func fillInGatewayStatusState(st types.State) *models.GatewayStatus {
	if st.ReportedState == nil {
		return nil
	}
	gwStatus := st.ReportedState.(*models.GatewayStatus)
	gwStatus.CheckinTime = st.TimeMs
	gwStatus.CertExpirationTime = st.CertExpirationTime
	gwStatus.HardwareID = st.ReporterID
	return gwStatus
}
