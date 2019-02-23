/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"magma/orc8r/cloud/go/dynamo"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

type Encoder interface {
	GetUpdateItemInputs(
		networkID string,
		updates map[string]*storage.GatewayUpdateParams) ([]*dynamodb.UpdateItemInput, error)
	GetBatchedWriteRequestsForGatewayViewDeletion(
		networkID string,
		gatewayIDs []string,
		batchSize int) ([][]*dynamodb.WriteRequest, error)
	GetKeysForGetItems(
		networkID string,
		gatewayIDs []string) ([]map[string]*dynamodb.AttributeValue, error)
	GetBatchGetItemInputs(
		keys []map[string]*dynamodb.AttributeValue,
		batchSize int) ([]*dynamodb.BatchGetItemInput, error)
}

type encoderImpl struct {
}

// MaxBatchSize is the maximum batch size for writes to DynamoDB
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_BatchWriteItem.html
const MaxBatchSize = 25

func NewEncoder() Encoder {
	return &encoderImpl{}
}

func (encoder *encoderImpl) GetUpdateItemInputs(
	networkID string,
	updates map[string]*storage.GatewayUpdateParams,
) ([]*dynamodb.UpdateItemInput, error) {
	output := make([]*dynamodb.UpdateItemInput, 0, len(updates))
	for gatewayID, update := range updates {
		dynamoState, err := updateToDynamoType(networkID, gatewayID, update)
		if err != nil {
			return nil, fmt.Errorf("Error converting gateway update to dynamo type: %s", err)
		}
		updateExpression, err := getUpdateExpression(dynamoState, update.NewConfig, update.ConfigsToDelete)
		if err != nil {
			return nil, err
		}
		key := getUpdateKeyFromDynamoType(dynamoState)
		updateInput := &dynamodb.UpdateItemInput{
			TableName:                 aws.String(GatewayViewTableName),
			Key:                       key,
			UpdateExpression:          updateExpression.Update(),
			ExpressionAttributeNames:  updateExpression.Names(),
			ExpressionAttributeValues: updateExpression.Values(),
		}
		output = append(output, updateInput)
	}
	return output, nil
}

func (encoder *encoderImpl) GetBatchedWriteRequestsForGatewayViewDeletion(
	networkID string,
	gatewayIDs []string,
	batchSize int,
) ([][]*dynamodb.WriteRequest, error) {
	if batchSize < 1 {
		return nil, fmt.Errorf("Batch size must be at least 1")
	}
	writeRequests, err := encoder.getWriteRequestsForGatewayViewDeletion(networkID, gatewayIDs)
	if err != nil {
		return nil, err
	}
	return dynamo.BatchWriteRequests(writeRequests, batchSize)
}

func (encoder *encoderImpl) GetKeysForGetItems(
	networkID string,
	gatewayIDs []string,
) ([]map[string]*dynamodb.AttributeValue, error) {
	networkIDAttrValue, err := dynamodbattribute.Marshal(networkID)
	if err != nil {
		return nil, err
	}
	// Create the list of keys to query
	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(gatewayIDs))
	for _, gatewayID := range gatewayIDs {
		gatewayIDAttrValue, err := dynamodbattribute.Marshal(gatewayID)
		if err != nil {
			return nil, err
		}
		key := map[string]*dynamodb.AttributeValue{
			NetworkIDName: networkIDAttrValue,
			GatewayIDName: gatewayIDAttrValue,
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (encoder *encoderImpl) GetBatchGetItemInputs(
	keys []map[string]*dynamodb.AttributeValue,
	batchSize int,
) ([]*dynamodb.BatchGetItemInput, error) {
	if batchSize < 1 {
		return nil, fmt.Errorf("Batch size must be at least 1")
	}
	output := []*dynamodb.BatchGetItemInput{}
	for start := 0; start < len(keys); start += batchSize {
		batchKeys := keys[start:min(start+batchSize, len(keys))]
		getItemInput := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				GatewayViewTableName: &dynamodb.KeysAndAttributes{
					Keys: batchKeys,
				},
			},
		}
		output = append(output, getItemInput)
	}
	return output, nil
}

// Convert the GatewayState type to the simplified type used for DynamoDB marshalling
func updateToDynamoType(networkID string, gatewayID string, src *storage.GatewayUpdateParams) (*gatewayState, error) {
	out := &gatewayState{
		NetworkID: networkID,
		GatewayID: gatewayID,
		Offset:    src.Offset,
	}
	if src.NewStatus != nil {
		statusBytes, err := protos.Marshal(src.NewStatus)
		if err != nil {
			return nil, err
		}
		out.Status = statusBytes
	}
	if src.NewRecord != nil {
		recordBytes, err := protos.Marshal(src.NewRecord)
		if err != nil {
			return nil, err
		}
		out.Record = recordBytes
	}
	return out, nil
}

func getUpdateKeyFromDynamoType(state *gatewayState) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		NetworkIDName: {S: aws.String(state.NetworkID)},
		GatewayIDName: {S: aws.String(state.GatewayID)},
	}
}

func getUpdateExpression(
	dynamoState *gatewayState,
	config map[string]interface{},
	configsToDelete []string,
) (*expression.Expression, error) {
	attributes, err := getAttributeMap(dynamoState, config)
	if err != nil {
		return nil, err
	}
	updateBuilder := expression.UpdateBuilder{}
	for _, name := range sortedKeys(attributes) {
		updateBuilder = updateBuilder.Set(expression.Name(name), expression.Value(attributes[name]))
	}
	for _, configType := range configsToDelete {
		updateBuilder = updateBuilder.Remove(expression.Name(configPrefix + configType))
	}
	out, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	return &out, err
}

func sortedKeys(attributes map[string]interface{}) []string {
	keys := []string{}
	for key := range attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func getAttributeMap(dynamoState *gatewayState, config map[string]interface{}) (map[string]interface{}, error) {
	attributes := map[string]interface{}{
		"Offset": dynamoState.Offset,
	}
	if dynamoState.Status != nil {
		attributes["Status"] = dynamoState.Status
	}
	if dynamoState.Record != nil {
		attributes["Record"] = dynamoState.Record
	}
	for configType, configObj := range config {
		configBytes, err := registry.MarshalConfig(configType, configObj)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling config value: %s", err)
		}
		attributes[configPrefix+configType] = configBytes
	}
	return attributes, nil
}

func (encoder *encoderImpl) getWriteRequestsForGatewayViewDeletion(
	networkID string,
	gatewayIDs []string,
) ([]*dynamodb.WriteRequest, error) {
	output := make([]*dynamodb.WriteRequest, 0, len(gatewayIDs))
	networkIDAttrValue, err := dynamodbattribute.Marshal(networkID)
	if err != nil {
		return nil, fmt.Errorf("error marshaling network ID to dynamoDB AttributeValues: %s", err)
	}
	for _, gatewayID := range gatewayIDs {
		gatewayIDAttrValue, err := dynamodbattribute.Marshal(gatewayID)
		if err != nil {
			return nil, fmt.Errorf("error marshaling gateway ID to dynamoDB AttributeValues: %s", err)
		}
		key := map[string]*dynamodb.AttributeValue{
			NetworkIDName: networkIDAttrValue,
			GatewayIDName: gatewayIDAttrValue,
		}
		request := &dynamodb.WriteRequest{DeleteRequest: &dynamodb.DeleteRequest{Key: key}}
		output = append(output, request)
	}
	return output, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
