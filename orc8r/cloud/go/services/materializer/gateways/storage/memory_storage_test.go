/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"

	"github.com/stretchr/testify/assert"
)

func TestGetGatewayViewsForNetwork(t *testing.T) {
	gatewayView := initStorage(t)
	gatewayStates, err := gatewayView.GetGatewayViewsForNetwork("net1")
	assert.NoError(t, err)
	gw1State, ok := gatewayStates["gw1"]
	assert.True(t, ok)
	assert.Equal(t, gw1State.Status, gateway1Status())
	assert.Equal(t, gw1State.Record, gateway1Record())
	assert.Equal(t, gw1State.Config, gateway1Config())

	// Just check that both got inserted
	_, ok = gatewayStates["gw2"]
	assert.True(t, ok)
}

func TestGetGatewayViewsForNetworkBadNetwork(t *testing.T) {
	gatewayView := initStorage(t)
	_, err := gatewayView.GetGatewayViewsForNetwork("badid")
	assert.EqualError(t, err, "Network ID not found: badid")
}

func TestGetGatewayView(t *testing.T) {
	gatewayView := initStorage(t)
	gatewayStates, err := gatewayView.GetGatewayViews("net1", []string{"gw2"})
	assert.NoError(t, err)
	gatewayState := gatewayStates["gw2"]
	assert.Equal(t, gatewayState.Status, gateway2Status())
	assert.Equal(t, gatewayState.Record, gateway2Record())
	assert.Equal(t, gatewayState.Config, gateway2Config())
}

func TestGetGatewayViewBadNetwork(t *testing.T) {
	gatewayView := initStorage(t)
	_, err := gatewayView.GetGatewayViews("badid", []string{"gw1"})
	assert.EqualError(t, err, "Network ID not found: badid")
}

func TestGetGatewayViewBadGateway(t *testing.T) {
	gatewayView := initStorage(t)
	_, err := gatewayView.GetGatewayViews("net1", []string{"badid"})
	assert.EqualError(t, err, "Gateway ID not found: badid")
}

func TestUpdateGateway(t *testing.T) {
	gatewayView := initStorage(t)
	newConfig := map[string]interface{}{
		"conf1": 11,
		"conf3": 103,
	}
	newParams := &storage.GatewayUpdateParams{
		NewConfig: newConfig,
		NewRecord: updatedGatewayRecord(),
		Offset:    2,
	}
	updates := make(map[string]*storage.GatewayUpdateParams)
	updates["gw1"] = newParams
	err := gatewayView.UpdateOrCreateGatewayViews("net1", updates)
	assert.NoError(t, err)
	states, err := gatewayView.GetGatewayViews("net1", []string{"gw1"})
	assert.NoError(t, err)
	state := states["gw1"]
	assert.Equal(t, state.Status, gateway1Status())
	assert.Equal(t, state.Record, updatedGatewayRecord())
	updatedConfig := map[string]interface{}{
		"conf1": 11,
		"conf2": 21,
		"conf3": 103,
	}
	assert.Equal(t, state.Config, updatedConfig)
}

func updatedGatewayRecord() *magmadprotos.AccessGatewayRecord {
	return &magmadprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{Id: "gw1"},
		Name: "Updated Gateway",
		Key: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_SOFTWARE_RSA_SHA256,
			Key:     []byte{4, 5, 6},
		},
		Ip:   "2.2.2.2",
		Port: 10000,
	}
}

func TestUpdateGatewayBadOffset(t *testing.T) {
	gatewayView := initStorage(t)
	newParams := &storage.GatewayUpdateParams{
		NewConfig: map[string]interface{}{
			"conf1": 11,
		},
		Offset: 0,
	}
	updates := make(map[string]*storage.GatewayUpdateParams)
	updates["gw1"] = newParams
	err := gatewayView.UpdateOrCreateGatewayViews("net1", updates)
	assert.EqualError(t, err, "Update offset less than current state")
	states, err := gatewayView.GetGatewayViews("net1", []string{"gw1"})
	assert.NoError(t, err)
	state := states["gw1"]
	assert.Equal(t, state.Config, gateway1Config())
}

func TestDeleteGateways(t *testing.T) {
	gatewayView := initStorage(t)
	toDelete := []string{
		"gw1",
	}
	err := gatewayView.DeleteGatewayViews("net1", toDelete)
	assert.NoError(t, err)
	_, err = gatewayView.GetGatewayViews("net1", []string{"gw1"})
	assert.EqualError(t, err, "Gateway ID not found: gw1")
}

func TestDeleteGatewaysBadNetwork(t *testing.T) {
	gatewayView := initStorage(t)
	err := gatewayView.DeleteGatewayViews("badid", []string{"gw1"})
	assert.EqualError(t, err, "Network ID not found: badid")
}

func TestDeleteGatewaysBadGateway(t *testing.T) {
	gatewayView := initStorage(t)
	err := gatewayView.DeleteGatewayViews("net1", []string{"badid"})
	assert.EqualError(t, err, "Gateway ID not found: badid")
}

func initStorage(t *testing.T) storage.GatewayViewStorage {
	gatewayView := storage.NewMemoryGatewayViewStorage()
	gatewayParams := make(map[string]*storage.GatewayUpdateParams)
	gw1Config := gateway1Config()
	gatewayParams["gw1"] = &storage.GatewayUpdateParams{
		NewConfig: gw1Config,
		NewStatus: gateway1Status(),
		NewRecord: gateway1Record(),
		Offset:    1,
	}
	gw2Config := gateway2Config()
	gatewayParams["gw2"] = &storage.GatewayUpdateParams{
		NewConfig: gw2Config,
		NewStatus: gateway2Status(),
		NewRecord: gateway2Record(),
		Offset:    1,
	}
	err := gatewayView.UpdateOrCreateGatewayViews("net1", gatewayParams)
	assert.NoError(t, err)
	return gatewayView
}

func gateway1Config() map[string]interface{} {
	return map[string]interface{}{
		"conf1": 17,
		"conf2": 21,
	}
}

func gateway1Status() *protos.GatewayStatus {
	return &protos.GatewayStatus{
		Time: 1,
		Checkin: &protos.CheckinRequest{
			GatewayId:       "gw1",
			MagmaPkgVersion: "v1",
			Status: &protos.ServiceStatus{
				Meta: make(map[string]string),
			},
			SystemStatus: nil,
		},
	}
}

func gateway1Record() *magmadprotos.AccessGatewayRecord {
	return &magmadprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{
			Id: "gw1",
		},
		Name: "gateway 1",
		Key: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_ECHO,
			Key:     []byte{1, 2, 3},
		},
		Ip:   "127.0.0.1",
		Port: 8000,
	}
}

func gateway2Config() map[string]interface{} {
	return map[string]interface{}{
		"conf1": 57,
		"conf2": 24,
	}
}

func gateway2Status() *protos.GatewayStatus {
	return &protos.GatewayStatus{
		Time: 2,
		Checkin: &protos.CheckinRequest{
			GatewayId:       "gw2",
			MagmaPkgVersion: "v2",
			Status: &protos.ServiceStatus{
				Meta: make(map[string]string),
			},
			SystemStatus: &protos.SystemStatus{
				Time:         2,
				CpuUser:      10,
				CpuSystem:    15,
				CpuIdle:      3,
				MemTotal:     1000,
				MemAvailable: 800,
				MemUsed:      200,
				MemFree:      0,
			},
		},
	}
}

func gateway2Record() *magmadprotos.AccessGatewayRecord {
	return &magmadprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{
			Id: "gw2",
		},
		Name: "Gateway 2",
		Key: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_SOFTWARE_ECDSA_SHA256,
			Key:     []byte{3, 2, 1},
		},
		Ip:   "1.1.1.1",
		Port: 4000,
	}
}
