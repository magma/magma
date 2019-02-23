/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo/mocks"
)

func TestGetGatewayViewsForNetwork(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Set up objects
	gw1ID := "gw1"
	gw2ID := "gw2"

	gw1State := &storage.GatewayState{GatewayID: gw1ID}
	gw2State := &storage.GatewayState{GatewayID: gw2ID}

	gatewayStates := map[string]*storage.GatewayState{
		gw1ID: gw1State,
		gw2ID: gw2State,
	}

	// Create queries
	query1, err := getGatewayViewsForNetworkQueryInput("net1")
	assert.NoError(t, err)
	query1Pages := [][]map[string]*dynamodb.AttributeValue{
		{map[string]*dynamodb.AttributeValue{
			dynamo.GatewayIDName: {S: aws.String(gw1ID)},
		}},
		{map[string]*dynamodb.AttributeValue{
			dynamo.GatewayIDName: {S: aws.String(gw2ID)},
		}},
	}
	mockDB.On("QueryPages", query1, mock.Anything).Return(getMockQueryPagesImpl(query1Pages))
	mockDecoder.On("GetStateFromAttributeMap", mock.MatchedBy(getQueryResultMatcherFn(query1Pages[0]))).Return(gw1State, nil)
	mockDecoder.On("GetStateFromAttributeMap", mock.MatchedBy(getQueryResultMatcherFn(query1Pages[1]))).Return(gw2State, nil)

	query2, err := getGatewayViewsForNetworkQueryInput("badnetwork")
	assert.NoError(t, err)
	mockDB.On("QueryPages", query2, mock.Anything).Return(errors.New("Bad network error"))

	query3, err := getGatewayViewsForNetworkQueryInput("net3")
	assert.NoError(t, err)
	query3Pages := [][]map[string]*dynamodb.AttributeValue{
		{map[string]*dynamodb.AttributeValue{
			dynamo.GatewayIDName: {S: aws.String("badgateway")},
		}},
	}
	mockDB.On("QueryPages", query3, mock.Anything).Return(getMockQueryPagesImpl(query3Pages))
	mockDecoder.On("GetStateFromAttributeMap", mock.MatchedBy(getQueryResultMatcherFn(query3Pages[0]))).Return(nil, errors.New("Bad gateway error"))

	// Run tests
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)

	actual, err := store.GetGatewayViewsForNetwork("net1")
	assert.NoError(t, err)
	assert.Equal(t, gatewayStates, actual)

	_, err = store.GetGatewayViewsForNetwork("badnetwork")
	assert.EqualError(t, err, "Bad network error")

	_, err = store.GetGatewayViewsForNetwork("net3")
	assert.EqualError(t, err, "Bad gateway error")
}

func getGatewayViewsForNetworkQueryInput(networkID string) (*dynamodb.QueryInput, error) {
	keyConditionBuilder := expression.KeyEqual(
		expression.Key(dynamo.NetworkIDName),
		expression.Value(networkID),
	)
	keyCondition, err := expression.NewBuilder().WithKeyCondition(keyConditionBuilder).Build()
	if err != nil {
		return nil, err
	}
	query := &dynamodb.QueryInput{
		TableName:                 aws.String(dynamo.GatewayViewTableName),
		KeyConditionExpression:    keyCondition.KeyCondition(),
		ExpressionAttributeNames:  keyCondition.Names(),
		ExpressionAttributeValues: keyCondition.Values(),
	}
	return query, nil
}

func TestGetGatewayViews(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Generate input/output objects
	networkID := "net1"
	gatewayIDs := []string{
		"gw0",
		"gw1",
		"gw2",
	}
	keys := []map[string]*dynamodb.AttributeValue{
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String("gw0")},
		},
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String("gw1")},
		},
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String("gw2")},
		},
	}
	batchGetInputs := []*dynamodb.BatchGetItemInput{
		getBatchGetItemInput(keys[0:2]),
		getBatchGetItemInput(keys[2:]),
	}
	responses := keys
	batchGetOutputs := []*dynamodb.BatchGetItemOutput{
		getBatchGetItemOutput(responses[0:2]),
		getBatchGetItemOutput(responses[2:]),
	}
	states := map[string]*storage.GatewayState{
		"gw0": {GatewayID: "gw0"},
		"gw1": {GatewayID: "gw1"},
		"gw2": {GatewayID: "gw2"},
	}

	// Register mock queries
	mockEncoder.On("GetKeysForGetItems", networkID, gatewayIDs).Return(keys, nil)
	mockEncoder.On("GetBatchGetItemInputs", keys, mock.Anything).Return(batchGetInputs, nil)
	for i := 0; i < len(batchGetInputs); i++ {
		mockDB.On("BatchGetItem", batchGetInputs[i]).Return(batchGetOutputs[i], nil)
	}
	for i := 0; i < len(responses); i++ {
		gatewayID := fmt.Sprintf("gw%d", i)
		mockDecoder.On("GetStateFromAttributeMap", responses[i]).Return(states[gatewayID], nil)
	}

	// Run test
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)
	actualStates, err := store.GetGatewayViews(networkID, gatewayIDs)
	assert.NoError(t, err)
	assert.Equal(t, states, actualStates)
}

func TestGetGatewayViewsUnprocessedKeys(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Generate input/output objects
	networkID := "net1"
	gatewayID := "gw1"
	gatewayIDs := []string{gatewayID}
	keys := []map[string]*dynamodb.AttributeValue{
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String(gatewayID)},
		},
	}
	requestItems := map[string]*dynamodb.KeysAndAttributes{dynamo.GatewayViewTableName: {Keys: keys}}
	batchGetInputs := []*dynamodb.BatchGetItemInput{{RequestItems: requestItems}}
	response := &dynamodb.BatchGetItemOutput{
		Responses:       make(map[string][]map[string]*dynamodb.AttributeValue),
		UnprocessedKeys: requestItems,
	}

	// Register mock queries
	mockEncoder.On("GetKeysForGetItems", networkID, gatewayIDs).Return(keys, nil)
	mockEncoder.On("GetBatchGetItemInputs", keys, mock.Anything).Return(batchGetInputs, nil)
	mockDB.On("BatchGetItem", batchGetInputs[0]).Return(response, nil)

	// Run test
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)
	_, err := store.GetGatewayViews(networkID, gatewayIDs)
	assert.True(t, strings.HasPrefix(err.Error(), "Remaining unprocessed keys after three attempts"))
	mockDB.AssertNumberOfCalls(t, "BatchGetItem", 3)
}

func TestGetGatewayViewsBatchGetDBError(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Generate input/output objects
	networkID := "net1"
	gatewayID := "gw1"
	gatewayIDs := []string{gatewayID}
	keys := []map[string]*dynamodb.AttributeValue{
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String(gatewayID)},
		},
	}
	requestItems := map[string]*dynamodb.KeysAndAttributes{
		dynamo.GatewayViewTableName: {
			Keys: keys,
		},
	}
	batchGetInputs := []*dynamodb.BatchGetItemInput{
		{
			RequestItems: requestItems,
		},
	}

	// Register mock queries
	mockEncoder.On("GetKeysForGetItems", networkID, gatewayIDs).Return(keys, nil)
	mockEncoder.On("GetBatchGetItemInputs", keys, mock.Anything).Return(batchGetInputs, nil)
	mockDB.On("BatchGetItem", batchGetInputs[0]).Return(nil, fmt.Errorf("Mock DB error"))

	// Run test
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)
	_, err := store.GetGatewayViews(networkID, gatewayIDs)
	assert.EqualError(t, err, "Error executing batch get: Mock DB error")
}

func getBatchGetItemInput(batch []map[string]*dynamodb.AttributeValue) *dynamodb.BatchGetItemInput {
	return &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			dynamo.GatewayViewTableName: {
				Keys: batch,
			},
		},
	}
}

func getBatchGetItemOutput(batch []map[string]*dynamodb.AttributeValue) *dynamodb.BatchGetItemOutput {
	return &dynamodb.BatchGetItemOutput{
		Responses: map[string][]map[string]*dynamodb.AttributeValue{
			dynamo.GatewayViewTableName: batch,
		},
	}
}

func TestGetGatewayViewsNetworkIDCaseInsensitivity(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Generate input/output objects
	networkID := "net1"
	capsNetworkID := "Net1"
	gatewayIDs := []string{
		"gw0",
		"gw1",
		"gw2",
	}
	keys := []map[string]*dynamodb.AttributeValue{
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String("gw0")},
		},
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String("gw1")},
		},
		{
			dynamo.NetworkIDName: {S: aws.String(networkID)},
			dynamo.GatewayIDName: {S: aws.String("gw2")},
		},
	}
	batchGetInputs := []*dynamodb.BatchGetItemInput{
		getBatchGetItemInput(keys[0:2]),
		getBatchGetItemInput(keys[2:]),
	}
	responses := keys
	batchGetOutputs := []*dynamodb.BatchGetItemOutput{
		getBatchGetItemOutput(responses[0:2]),
		getBatchGetItemOutput(responses[2:]),
	}
	states := map[string]*storage.GatewayState{
		"gw0": {GatewayID: "gw0"},
		"gw1": {GatewayID: "gw1"},
		"gw2": {GatewayID: "gw2"},
	}

	// Register mock queries
	mockEncoder.On("GetKeysForGetItems", networkID, gatewayIDs).Return(keys, nil)
	mockEncoder.On("GetBatchGetItemInputs", keys, mock.Anything).Return(batchGetInputs, nil)
	for i := 0; i < len(batchGetInputs); i++ {
		mockDB.On("BatchGetItem", batchGetInputs[i]).Return(batchGetOutputs[i], nil)
	}
	for i := 0; i < len(responses); i++ {
		gatewayID := fmt.Sprintf("gw%d", i)
		mockDecoder.On("GetStateFromAttributeMap", responses[i]).Return(states[gatewayID], nil)
	}

	// Run test
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)
	actualStates, err := store.GetGatewayViews(capsNetworkID, gatewayIDs)
	assert.NoError(t, err)
	assert.Equal(t, states, actualStates)
}

func TestUpdateOrCreateGatewayViews(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Generate input/output objects
	networkID := "net1"
	updates := map[string]*storage.GatewayUpdateParams{
		"gw1": {
			NewStatus: &protos.GatewayStatus{},
			Offset:    1,
		},
		"gw2": {
			Offset: 5,
		},
	}
	inputs := []*dynamodb.UpdateItemInput{
		{
			TableName: aws.String("Fake Table 1"),
		},
		{
			TableName: aws.String("Fake Table 2"),
		},
	}

	// Register mock queries
	mockEncoder.On("GetUpdateItemInputs", networkID, updates).Return(inputs, nil)
	mockDB.On("UpdateItem", inputs[0]).Return(nil, nil)
	mockDB.On("UpdateItem", inputs[1]).Return(nil, nil)

	// Run test
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)
	err := store.UpdateOrCreateGatewayViews(networkID, updates)
	assert.NoError(t, err)
	mockDB.AssertNumberOfCalls(t, "UpdateItem", 2)
}

func TestDeleteGatewayViews(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Generate input/output objects
	networkID := "net1"
	gatewayIDs := []string{"gw1", "gw2"}
	batchedWrites := [][]*dynamodb.WriteRequest{
		{
			{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]*dynamodb.AttributeValue{
						dynamo.NetworkIDName: {S: aws.String(networkID)},
						dynamo.GatewayIDName: {S: aws.String(gatewayIDs[0])},
					},
				},
			},
		},
		{
			{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]*dynamodb.AttributeValue{
						dynamo.NetworkIDName: {S: aws.String(networkID)},
						dynamo.GatewayIDName: {S: aws.String(gatewayIDs[1])},
					},
				},
			},
		},
	}
	batchInputs := []*dynamodb.BatchWriteItemInput{
		{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				dynamo.GatewayViewTableName: batchedWrites[0],
			},
		},
		{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				dynamo.GatewayViewTableName: batchedWrites[1],
			},
		},
	}
	response := &dynamodb.BatchWriteItemOutput{
		UnprocessedItems: make(map[string][]*dynamodb.WriteRequest),
	}

	// Register mock queries
	mockEncoder.On("GetBatchedWriteRequestsForGatewayViewDeletion", networkID, gatewayIDs, mock.Anything).Return(batchedWrites, nil)
	mockDB.On("BatchWriteItem", batchInputs[0]).Return(response, nil)
	mockDB.On("BatchWriteItem", batchInputs[1]).Return(response, nil)

	// Run test
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)
	err := store.DeleteGatewayViews(networkID, gatewayIDs)
	assert.NoError(t, err)
	mockDB.AssertNumberOfCalls(t, "BatchWriteItem", 2)
}

func TestDeleteGatewayViewsUnprocessedWrites(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Generate input/output objects
	// Testing with 2 batches. One needs to be full since the function rebatches unprocessed writes.
	batch1 := make([]*dynamodb.WriteRequest, dynamo.MaxBatchSize)
	for i := 0; i < len(batch1); i++ {
		batch1[i] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					dynamo.NetworkIDName: {S: aws.String("net1")},
					dynamo.GatewayIDName: {S: aws.String(fmt.Sprintf("gw%d", i))},
				},
			},
		}
	}
	batch2 := []*dynamodb.WriteRequest{
		{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					dynamo.NetworkIDName: {S: aws.String("net1")},
					dynamo.GatewayIDName: {S: aws.String("extraGW")},
				},
			},
		},
	}
	batchedWrites := [][]*dynamodb.WriteRequest{batch1, batch2}
	batchWriteInput1 := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			dynamo.GatewayViewTableName: batch1,
		},
	}
	response1 := &dynamodb.BatchWriteItemOutput{
		UnprocessedItems: map[string][]*dynamodb.WriteRequest{
			dynamo.GatewayViewTableName: batch1,
		},
	}
	batchWriteInput2 := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			dynamo.GatewayViewTableName: batch2,
		},
	}
	response2 := &dynamodb.BatchWriteItemOutput{
		UnprocessedItems: map[string][]*dynamodb.WriteRequest{
			dynamo.GatewayViewTableName: batch2,
		},
	}

	// Register mock queries
	mockEncoder.On("GetBatchedWriteRequestsForGatewayViewDeletion", mock.Anything, mock.Anything, mock.Anything).Return(batchedWrites, nil)
	mockDB.On("BatchWriteItem", batchWriteInput1).Return(response1, nil)
	mockDB.On("BatchWriteItem", batchWriteInput2).Return(response2, nil)

	// Run test
	store := dynamo.NewDynamoDBGatewayStorage(mockDB, mockEncoder, mockDecoder)
	err := store.DeleteGatewayViews("net1", []string{})
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Unwritten writes"))
	mockDB.AssertNumberOfCalls(t, "BatchWriteItem", 6)
}

func getMockQueryPagesImpl(queryPages [][]map[string]*dynamodb.AttributeValue) func(*dynamodb.QueryInput, func(*dynamodb.QueryOutput, bool) bool) error {
	return func(queryInput *dynamodb.QueryInput, pageHandler func(result *dynamodb.QueryOutput, lastPage bool) bool) error {
		for i, itemList := range queryPages {
			lastPage := i == len(queryPages)-1
			if !pageHandler(&dynamodb.QueryOutput{Items: itemList}, lastPage) {
				return nil
			}
		}
		return nil
	}
}

func getQueryResultMatcherFn(mockQueryResults []map[string]*dynamodb.AttributeValue) func(map[string]*dynamodb.AttributeValue) bool {
	return func(in map[string]*dynamodb.AttributeValue) bool {
		for _, queryResult := range mockQueryResults {
			if reflect.DeepEqual(queryResult, in) {
				return true
			}
		}
		return false
	}
}
