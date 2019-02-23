/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	storage_mocks "magma/orc8r/cloud/go/services/materializer/gateways/storage/mocks"
	"magma/orc8r/cloud/go/services/materializer/gateways/streaming"
	"magma/orc8r/cloud/go/services/materializer/gateways/streaming/mocks"

	"github.com/stretchr/testify/assert"
)

func TestApply_UpdateOrCreateStatus(t *testing.T) {
	status := &protos.GatewayStatus{
		Time: 12345,
		Checkin: &protos.CheckinRequest{
			GatewayId:       "gw1",
			MagmaPkgVersion: "v1",
		},
		CertExpirationTime: 12345,
	}
	statusBytes, err := json.Marshal(status)
	assert.NoError(t, err)
	base64Bytes := make([]byte, base64.StdEncoding.EncodedLen(len(statusBytes)))
	base64.StdEncoding.Encode(base64Bytes, statusBytes)
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Statuses,
		Operation:  "c",
		NetworkID:  "net1",
		Payload: &streaming.GatewayStatusUpdate{
			GatewayID:   "gw1",
			StatusBytes: string(base64Bytes),
		},
	}
	updateParams := map[string]*storage.GatewayUpdateParams{
		"gw1": {
			NewStatus: status,
			Offset:    12345,
		},
	}

	mockStore := &storage_mocks.GatewayViewStorage{}
	mockStore.On("UpdateOrCreateGatewayViews", "net1", updateParams).Return(nil)

	err = update.ApplyUpdate(mockStore, 12345)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockStore.AssertCalled(t, "UpdateOrCreateGatewayViews", "net1", updateParams)
}

func TestApply_DeleteStatus(t *testing.T) {
	// Tests that store/payload items never touched on a status delete
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Statuses,
		Operation:  "d",
		NetworkID:  "net1",
		Payload:    &streaming.GatewayStatusUpdate{},
	}
	err := update.ApplyUpdate(nil, 0)
	assert.NoError(t, err)
}

func TestApply_UpdateOrCreateRecord(t *testing.T) {
	record := &magmadprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{
			Id: "hwid1",
		},
		Name: "gateway 1",
		Key: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_ECHO,
		},
	}
	recordBytes, err := json.Marshal(record)
	assert.NoError(t, err)
	base64Bytes := make([]byte, base64.StdEncoding.EncodedLen(len(recordBytes)))
	base64.StdEncoding.Encode(base64Bytes, recordBytes)
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Records,
		Operation:  "c",
		NetworkID:  "net1",
		Payload: &streaming.GatewayRecordUpdate{
			GatewayID:   "gw1",
			RecordBytes: string(base64Bytes),
		},
	}
	updateParams := map[string]*storage.GatewayUpdateParams{
		"gw1": &storage.GatewayUpdateParams{
			NewRecord: record,
			Offset:    12345,
		},
	}

	mockStore := &storage_mocks.GatewayViewStorage{}
	mockStore.On("UpdateOrCreateGatewayViews", "net1", updateParams).Return(nil)

	err = update.ApplyUpdate(mockStore, 12345)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockStore.AssertCalled(t, "UpdateOrCreateGatewayViews", "net1", updateParams)
}

func TestApply_DeleteRecord(t *testing.T) {
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Records,
		Operation:  "d",
		NetworkID:  "net1",
		Payload: &streaming.GatewayRecordUpdate{
			GatewayID: "gw1",
		},
	}
	mockStore := &storage_mocks.GatewayViewStorage{}
	mockStore.On("DeleteGatewayViews", "net1", []string{"gw1"}).Return(nil)

	err := update.ApplyUpdate(mockStore, 12345)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockStore.AssertCalled(t, "DeleteGatewayViews", "net1", []string{"gw1"})
}

func TestApply_UpdateOrCreateConfig(t *testing.T) {
	mockConfigManager := &mocks.ConfigManager{}
	configType := "test_config"
	networkID := "net1"
	config := "Hello, World!"

	mockConfigManager.On("GetConfigType").Return(configType)
	mockConfigManager.On("GetGatewayIdsForConfig", networkID, networkID).Return([]string{"gw1"}, nil)
	mockConfigManager.On("UnmarshalConfig", []byte(config)).Return(config, nil)

	registry.ClearRegistryForTesting()
	err := registry.RegisterConfigManager(mockConfigManager)
	assert.NoError(t, err)

	configBytes := []byte(config)
	base64Bytes := make([]byte, base64.StdEncoding.EncodedLen(len(configBytes)))
	base64.StdEncoding.Encode(base64Bytes, configBytes)

	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Configurations,
		Operation:  "c",
		NetworkID:  "net1",
		Payload: &streaming.GatewayConfigUpdate{
			ConfigKey:   "net1",
			ConfigType:  configType,
			ConfigBytes: string(base64Bytes),
		},
	}
	updateParam := &storage.GatewayUpdateParams{
		NewConfig: map[string]interface{}{
			configType: config,
		},
		Offset: 12345,
	}
	updateParams := map[string]*storage.GatewayUpdateParams{
		"gw1": updateParam,
	}

	mockStore := &storage_mocks.GatewayViewStorage{}
	mockStore.On("UpdateOrCreateGatewayViews", "net1", updateParams).Return(nil)

	err = update.ApplyUpdate(mockStore, 12345)
	assert.NoError(t, err)

	// Verify that we skip non-gateway level configs - have the config manager
	// return 2 gateway IDs
	update.Payload = &streaming.GatewayConfigUpdate{
		ConfigKey:   "net2",
		ConfigType:  configType,
		ConfigBytes: string(base64Bytes),
	}
	mockConfigManager.On("GetGatewayIdsForConfig", networkID, "net2").Return([]string{"gw1", "gw2"}, nil)
	err = update.ApplyUpdate(mockStore, 12345)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockStore.AssertCalled(t, "UpdateOrCreateGatewayViews", "net1", updateParams)
	mockStore.AssertNumberOfCalls(t, "UpdateOrCreateGatewayViews", 1)
}

func TestApply_DeleteConfig(t *testing.T) {
	mockConfigManager := &mocks.ConfigManager{}
	configType := "test_config"
	networkID := "net1"
	mockConfigManager.On("GetConfigType").Return(configType)
	mockConfigManager.On("GetGatewayIdsForConfig", networkID, networkID).Return([]string{"gw1", "gw2"}, nil)

	registry.ClearRegistryForTesting()
	err := registry.RegisterConfigManager(mockConfigManager)
	assert.NoError(t, err)

	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Configurations,
		Operation:  "d",
		NetworkID:  "net1",
		Payload: &streaming.GatewayConfigUpdate{
			ConfigKey:  networkID,
			ConfigType: configType,
		},
	}
	updateParam := &storage.GatewayUpdateParams{
		ConfigsToDelete: []string{configType},
		Offset:          12345,
	}
	updateParams := map[string]*storage.GatewayUpdateParams{
		"gw1": updateParam,
		"gw2": updateParam,
	}
	mockStore := &storage_mocks.GatewayViewStorage{}
	mockStore.On("UpdateOrCreateGatewayViews", networkID, updateParams).Return(nil)

	err = update.ApplyUpdate(mockStore, 12345)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockStore.AssertCalled(t, "UpdateOrCreateGatewayViews", networkID, updateParams)
}

func TestGetUpdateParams_RealRecord(t *testing.T) {
	// Actual string seen in Kafka during integration testing
	recordBytes := "ewogImh3SWQiOiB7CiAgImlkIjogIjIyZmZlYTEwLTdmYzQtNDQyNy05NzVhLWI5ZTRjZThmNmY0ZCIKIH0sCiAibmFtZSI6ICJOZXcgbmFtZSAyIiwKICJrZXkiOiB7CiAgImtleVR5cGUiOiAiRUNITyIsCiAgImtleSI6IG51bGwKIH0sCiAiaXAiOiAiIiwKICJwb3J0IjogMAp9"
	aggregatedUpdate := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Records,
		Operation:  "c",
		NetworkID:  "net1",
		Payload: &streaming.GatewayRecordUpdate{
			GatewayID:   "gw1",
			RecordBytes: recordBytes,
		},
	}
	expectedParams := map[string]*storage.GatewayUpdateParams{
		"gw1": {
			NewRecord: &magmadprotos.AccessGatewayRecord{
				HwId: &protos.AccessGatewayID{Id: "22ffea10-7fc4-4427-975a-b9e4ce8f6f4d"},
				Name: "New name 2",
				Key: &protos.ChallengeKey{
					KeyType: protos.ChallengeKey_ECHO,
				},
				Ip:   "",
				Port: 0,
			},
			Offset: 12345,
		},
	}

	mockStore := &storage_mocks.GatewayViewStorage{}
	mockStore.On("UpdateOrCreateGatewayViews", "net1", expectedParams).Return(nil)

	err := aggregatedUpdate.ApplyUpdate(mockStore, 12345)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockStore.AssertCalled(t, "UpdateOrCreateGatewayViews", "net1", expectedParams)
}
