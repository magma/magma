/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming_test

import (
	"errors"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/config/streaming"
	"magma/orc8r/cloud/go/services/config/streaming/mocks"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	storage_protos "magma/orc8r/cloud/go/services/config/streaming/storage/protos"
	"magma/orc8r/cloud/go/services/config/streaming/storage/test_protos"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Integration test that populates the registry with mconfig streamers to test
// the system a bit more end-to-end.
func TestStreamProcessor_Run_Integration(t *testing.T) {
	// Setup: updates to config type1 affect gw1 and gw2
	// updates to configs type 2 affect gw3

	// An update to type1 will append an mconfig1 value to the mconfig of all
	// gateways. An update to type2 will append an mconfig2 value to the
	// mconfig of all gateways.
	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(&mockConfigManager{retType: "type1", gatewayIds: []string{"gw1", "gw2"}, retError: nil})
	registry.RegisterConfigManager(&mockConfigManager{retType: "type2", gatewayIds: []string{"gw3"}, retError: nil})

	streaming.ClearRegistryForTesting()
	streaming.RegisterMconfigStreamer(
		&mockConfigStreamer{
			subscribedTypes: []string{"type1", "type2"},
			retKey:          "mconfig1",
			retValue:        getExpectedMconfig1Value(t),
			retErr:          nil,
		},
	)
	streaming.RegisterMconfigStreamer(
		&mockConfigStreamer{
			subscribedTypes: []string{"type2"},
			retKey:          "mconfig2",
			retValue:        getExpectedMconfig2Value(t),
			retErr:          nil,
		},
	)

	ds := test_utils.NewMockDatastore()
	store := storage.NewDatastoreMconfigStorage(ds)

	// Send 2 config update messages and a gateway create, then an error to end the test
	mockConsumer := &mocks.StreamConsumer{}

	message1 := &kafka.Message{Key: []byte("message1"), TopicPartition: kafka.TopicPartition{Offset: 1}}
	message2 := &kafka.Message{Key: []byte("message2"), TopicPartition: kafka.TopicPartition{Offset: 2}}
	message3 := &kafka.Message{Key: []byte("message3"), TopicPartition: kafka.TopicPartition{Offset: 0}}

	mockConsumer.On("SubscribeTopics", mock.Anything, mock.Anything).
		Return(nil)
	mockConsumer.On("Commit").Return(nil, nil)
	mockConsumer.On("Close").Return(nil)

	mockConsumer.On("ReadMessage", mock.Anything).Once().
		Return(message1, nil)
	mockConsumer.On("ReadMessage", mock.Anything).Once().
		Return(message2, nil)
	mockConsumer.On("ReadMessage", mock.Anything).Once().
		Return(message3, nil)
	mockConsumer.On("ReadMessage", mock.Anything).Once().
		Return(nil, errors.New("Ending the test"))

	// Mock the decode so that the first message will be an update to type1
	// affecting gw1 and gw2, and the second will be and update to type2
	// affecting gw3 only.
	// The 3rd message will be a gateway creation
	mockDecoder := &mocks.Decoder{}
	update1 := &streaming.ConfigUpdate{ConfigType: "type1", ConfigKey: "key", NetworkId: "network1"}
	update2 := &streaming.ConfigUpdate{ConfigType: "type2", ConfigKey: "key", NetworkId: "network1"}
	update3 := &streaming.GatewayUpdate{GatewayId: "gw4", NetworkId: "network1", Operation: streaming.CreateOperation}
	mockDecoder.On("GetUpdateFromMessage", message1).
		Return(update1, nil)
	mockDecoder.On("GetUpdateFromMessage", message2).
		Return(update2, nil)
	mockDecoder.On("GetUpdateFromMessage", message3).
		Return(update3, nil)

	processor := streaming.NewStreamProcessor(store, mockDecoder, getMockConsumerFactory(mockConsumer))
	err := processor.Run()
	assert.EqualError(t, err, "Ending the test")

	// Assert mocks first
	mockConsumer.AssertExpectations(t)
	mockConsumer.AssertNumberOfCalls(t, "ReadMessage", 4)
	mockConsumer.AssertNumberOfCalls(t, "Commit", 3)
	mockConsumer.AssertNumberOfCalls(t, "Close", 1)

	mockDecoder.AssertExpectations(t)
	mockDecoder.AssertNumberOfCalls(t, "GetUpdateFromMessage", 3)

	// Now check the datastore values
	expectedGw1AndGw2Value := &storage_protos.StoredMconfig{
		Configs: &protos.GatewayConfigs{
			ConfigsByKey: map[string]*any.Any{
				"mconfig1": getExpectedMconfig1Value(t),
			},
		},
		Offset: 1,
	}
	// Both streamers should have executed and appended their values to gw3
	expectedGw3Value := &storage_protos.StoredMconfig{
		Configs: &protos.GatewayConfigs{
			ConfigsByKey: map[string]*any.Any{
				"mconfig1": getExpectedMconfig1Value(t),
				"mconfig2": getExpectedMconfig2Value(t),
			},
		},
		Offset: 2,
	}
	// And both streamers should have executed on the add gateway event for gw4
	expectedGw4Value := &storage_protos.StoredMconfig{
		Configs: &protos.GatewayConfigs{
			ConfigsByKey: map[string]*any.Any{
				"mconfig1": getExpectedMconfig1Value(t),
				"mconfig2": getExpectedMconfig2Value(t),
			},
		},
		Offset: -1,
	}

	test_utils.AssertDatastoreHasRows(
		t, ds,
		"network1_mconfig_views",
		map[string]interface{}{
			"gw1": expectedGw1AndGw2Value,
			"gw2": expectedGw1AndGw2Value,
			"gw3": expectedGw3Value,
			"gw4": expectedGw4Value,
		},
		deserializeStoredMconfigProto,
	)
}

func getMockConsumerFactory(mockConsumer *mocks.StreamConsumer) streaming.StreamConsumerFactory {
	return func() (streaming.StreamConsumer, error) {
		return mockConsumer, nil
	}
}

type mockConfigManager struct {
	retType    string
	gatewayIds []string
	retError   error
}

func (mcm *mockConfigManager) GetConfigType() string {
	return mcm.retType
}

func (mcm *mockConfigManager) GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error) {
	return mcm.gatewayIds, mcm.retError
}

func (mcm *mockConfigManager) MarshalConfig(config interface{}) ([]byte, error) {
	return []byte(fmt.Sprintf("%s mock value", mcm.retType)), nil
}

func (mcm *mockConfigManager) UnmarshalConfig(message []byte) (interface{}, error) {
	return fmt.Sprintf("%s mock value", mcm.retType), nil
}

type mockConfigStreamer struct {
	subscribedTypes []string
	retKey          string
	retValue        *any.Any
	retErr          error
}

func (mcs *mockConfigStreamer) GetSubscribedConfigTypes() []string {
	return mcs.subscribedTypes
}

func (mcs *mockConfigStreamer) SeedNewGatewayMconfig(
	networkId string,
	gatewayId string,
	mconfigOut *protos.GatewayConfigs, // output parameter
) error {
	mconfigOut.ConfigsByKey[mcs.retKey] = mcs.retValue
	return mcs.retErr
}

func (mcs *mockConfigStreamer) ApplyMconfigUpdate(
	update *streaming.ConfigUpdate,
	oldMconfigsByGatewayId map[string]*protos.GatewayConfigs,
) (map[string]*protos.GatewayConfigs, error) {
	for _, cfg := range oldMconfigsByGatewayId {
		cfg.ConfigsByKey[mcs.retKey] = mcs.retValue
	}
	return oldMconfigsByGatewayId, mcs.retErr
}

func getExpectedMconfig1Value(t *testing.T) *any.Any {
	val := &test_protos.Config1{Field: "mconfig1"}
	ret, err := ptypes.MarshalAny(val)
	assert.NoError(t, err)
	return ret
}

func getExpectedMconfig2Value(t *testing.T) *any.Any {
	val := &test_protos.Config2{Field1: "mconfig2Field1", Field2: "mconfig2Field2"}
	ret, err := ptypes.MarshalAny(val)
	assert.NoError(t, err)
	return ret
}

func deserializeStoredMconfigProto(msg []byte) (interface{}, error) {
	ret := &storage_protos.StoredMconfig{}
	err := protos.Unmarshal(msg, ret)
	return ret, err
}
