/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/test_utils"
)

func TestGetKeysForGetItems(t *testing.T) {
	encoder := dynamo.NewEncoder()
	networkID := "net1"
	numGateways := 4
	gatewayIDs := make([]string, numGateways)
	for i := 0; i < numGateways; i++ {
		gatewayIDs[i] = fmt.Sprintf("gw%d", i)
	}
	keys, err := encoder.GetKeysForGetItems(networkID, gatewayIDs)
	assert.NoError(t, err)
	assert.Equal(t, numGateways, len(keys))
	for i, key := range keys {
		verifyKey(t, networkID, i, key)
	}
}

func TestGetBatchGetItemInputs(t *testing.T) {
	encoder := dynamo.NewEncoder()
	networkID := "net1"
	numGateways := 3
	batchSize := 2
	gatewayIDs := make([]string, numGateways)
	for i := 0; i < numGateways; i++ {
		gatewayIDs[i] = fmt.Sprintf("gw%d", i)
	}
	keys, err := encoder.GetKeysForGetItems(networkID, gatewayIDs)
	assert.NoError(t, err)
	batchGetItemInputs, err := encoder.GetBatchGetItemInputs(keys, batchSize)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(batchGetItemInputs))
	for batchNum, batchGetItemInput := range batchGetItemInputs {
		requestItems, ok := batchGetItemInput.RequestItems[dynamo.GatewayViewTableName]
		assert.True(t, ok, "Request entries for gateway view table storage not found")
		batchKeys := requestItems.Keys
		assert.True(t, len(batchKeys) <= batchSize, "Batch found with length greater than batch size")
		for i, key := range batchKeys {
			verifyKey(t, networkID, batchNum*batchSize+i, key)
		}
	}
}

func TestGetUpdateItemInputs(t *testing.T) {
	// Set up config
	registry.ClearRegistryForTesting()
	configManager := test_utils.NewConfig1Manager()
	registry.RegisterConfigManager(configManager)
	config1 := &test_utils.Conf1{
		Value1: 4,
		Value2: "Hello, world!",
		Value3: []byte{3, 4, 5},
	}
	config1Bytes, err := registry.MarshalConfig(configManager.GetConfigType(), config1)
	assert.NoError(t, err)
	config := map[string]interface{}{
		configManager.GetConfigType(): config1,
	}

	newStatus := &protos.GatewayStatus{
		Time:    1234,
		Checkin: nil,
	}
	newStatusBytes, err := protos.Marshal(newStatus)
	assert.NoError(t, err)

	encoder := dynamo.NewEncoder()
	networkID := "net1"
	numUpdates := 3
	updates := make(map[string]*storage.GatewayUpdateParams)
	for i := 0; i < numUpdates; i++ {
		gatewayID := fmt.Sprintf("gw%d", i)
		updates[gatewayID] = &storage.GatewayUpdateParams{
			NewConfig:       config,
			NewStatus:       newStatus,
			Offset:          1,
			ConfigsToDelete: []string{"ConfigDeleteType1", "ConfigDeleteType2"},
		}
	}
	updateBuilder := expression.Set(
		expression.Name("Offset"),
		expression.Value(1),
	).Set(
		expression.Name("Status"),
		expression.Value(newStatusBytes),
	).Set(
		expression.Name("cfg_Config1"),
		expression.Value(config1Bytes),
	).Remove(
		expression.Name("cfg_ConfigDeleteType1"),
	).Remove(
		expression.Name("cfg_ConfigDeleteType2"),
	)
	updateExpression, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	assert.NoError(t, err)

	updateInputs, err := encoder.GetUpdateItemInputs(networkID, updates)
	assert.NoError(t, err)

	for _, updateInput := range updateInputs {
		assert.Equal(t, updateExpression.Update(), updateInput.UpdateExpression)
		assert.Equal(t, updateExpression.Names(), updateInput.ExpressionAttributeNames)
		assert.Equal(t, updateExpression.Values(), updateInput.ExpressionAttributeValues)
	}
}

func TestGetBatchedWriteRequestsForGatewayViewDeletion(t *testing.T) {
	networkID := "net1"
	numGateways := 3
	batchSize := 2
	gatewayIDs := make([]string, numGateways)
	for i := 0; i < numGateways; i++ {
		gatewayIDs[i] = fmt.Sprintf("gw%d", i)
	}
	encoder := dynamo.NewEncoder()
	batchedWrites, err := encoder.GetBatchedWriteRequestsForGatewayViewDeletion(networkID, gatewayIDs, batchSize)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(batchedWrites))
	for batchNum, batch := range batchedWrites {
		assert.True(t, len(batch) <= batchSize, "Batch size larger than maximum write batch size")
		for i, write := range batch {
			verifyKey(t, networkID, batchNum*batchSize+i, write.DeleteRequest.Key)
		}
	}
}

func verifyKey(t *testing.T, networkID string, index int, key map[string]*dynamodb.AttributeValue) {
	gatewayID := fmt.Sprintf("gw%d", index)

	networkIDAttrValue := key[dynamo.NetworkIDName]
	var actualNetworkID string
	err := dynamodbattribute.Unmarshal(networkIDAttrValue, &actualNetworkID)
	assert.NoError(t, err)
	assert.Equal(t, networkID, actualNetworkID, "Network IDs do not match")

	gatewayIDAttrValue := key[dynamo.GatewayIDName]
	var actualGatewayID string
	err = dynamodbattribute.Unmarshal(gatewayIDAttrValue, &actualGatewayID)
	assert.NoError(t, err)
	assert.Equal(t, gatewayID, actualGatewayID, "Gateway IDs do not match")
}
