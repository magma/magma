/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package dynamo contains an implementation of GatewayViewStorage which is
// backed by DynamoDB
package dynamo

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"magma/orc8r/cloud/go/dynamo"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

const (
	GatewayViewTableName = "gateway_storage_view"
	SchemaVersion        = 1
	NetworkIDName        = "NetworkID"
	GatewayIDName        = "GatewayID"
	// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_BatchGetItem.html
	maxGetBatchSize = 100
)

type dynamoGatewayStorage struct {
	db      dynamodbiface.DynamoDBAPI
	encoder Encoder
	decoder Decoder
}

/*
	Create the gateway view tables if they don't already exist.
	The gateway view schema looks like this

	Field			|	Type	|	Key Type	|	 Notes
	=====================================================
	NetworkID       |   S       |   hash        |
	GatewayID       |   S       |   range       |
	Status          |   B       |               |    Stored as JSON-encoded byte array
	Record          |   B       |               |    Stored as JSON-encoded byte array
	Offset          |   N       |               |
	cfg_*           |   B       |               |    JSON-encoded byte arrays for each config value

	IMPORTANT: DO NOT MODIFY THIS METHOD BEFORE MAKING A CORRESPONDING CHANGE
	IN THE TERRAFORM DYNAMODB RESOURCE. On staging/prod, this method isn't
	called because we rely on Terraform to manage the DynamoDB resources.

	This is only called on local dev VM in order to set up the local
	dynamoDB executable. Changing this method without changing the terraform
	resource will have no effect on staging or prod. If you change the
	terraform resource without modifying this method, your dev VM schema will
	be outdated.
*/
func (gs *dynamoGatewayStorage) InitTables() error {
	_, err := gs.db.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(GatewayViewTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(NetworkIDName),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String(GatewayIDName),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(NetworkIDName),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String(GatewayIDName),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(100),
			WriteCapacityUnits: aws.Int64(100),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// Table already exists, we're good to go
			case dynamodb.ErrCodeResourceInUseException:
				return nil
			default:
				return fmt.Errorf("Error creating gateway view tables: %s", err)
			}
		} else {
			return fmt.Errorf("Error creating gateway view tables: %s", err)
		}
	}
	return nil
}

func NewDynamoDBGatewayStorage(db dynamodbiface.DynamoDBAPI, encoder Encoder, decoder Decoder) storage.GatewayViewStorage {
	return &dynamoGatewayStorage{db: db, encoder: encoder, decoder: decoder}
}

func (gs *dynamoGatewayStorage) GetGatewayViewsForNetwork(networkID string) (map[string]*storage.GatewayState, error) {
	keyConditionBuilder := expression.KeyEqual(
		expression.Key(NetworkIDName),
		expression.Value(networkID),
	)
	keyCondition, err := expression.NewBuilder().WithKeyCondition(keyConditionBuilder).Build()
	if err != nil {
		return nil, err
	}
	query := &dynamodb.QueryInput{
		TableName:                 aws.String(GatewayViewTableName),
		KeyConditionExpression:    keyCondition.KeyCondition(),
		ExpressionAttributeNames:  keyCondition.Names(),
		ExpressionAttributeValues: keyCondition.Values(),
	}
	output := make(map[string]*storage.GatewayState)
	var callbackErr error
	err = gs.db.QueryPages(query, func(result *dynamodb.QueryOutput, lastPage bool) bool {
		pageStateMap, err := gs.getStatesFromPageItems(result.Items)
		if err != nil {
			callbackErr = err
			return false
		}
		for gatewayID, state := range pageStateMap {
			output[gatewayID] = state
		}
		return !lastPage
	})
	if callbackErr != nil {
		return nil, callbackErr
	}
	return output, err
}

func (gs *dynamoGatewayStorage) getStatesFromPageItems(
	items []map[string]*dynamodb.AttributeValue,
) (map[string]*storage.GatewayState, error) {
	stateMap := make(map[string]*storage.GatewayState)
	for _, attributeValues := range items {
		state, err := gs.decoder.GetStateFromAttributeMap(attributeValues)
		if err != nil {
			return nil, err
		}
		stateMap[state.GatewayID] = state
	}
	return stateMap, nil
}

func (gs *dynamoGatewayStorage) GetGatewayViews(
	networkID string,
	gatewayIDs []string,
) (map[string]*storage.GatewayState, error) {
	keys, err := gs.encoder.GetKeysForGetItems(strings.ToLower(networkID), gatewayIDs)
	if err != nil {
		return nil, fmt.Errorf("Error getting DynamoDB keys for gateways: %s", err)
	}
	// Get views in batches
	output := make(map[string]*storage.GatewayState)
	for attempts := 0; len(keys) > 0 && attempts < 3; attempts++ {
		batchGetInputs, err := gs.encoder.GetBatchGetItemInputs(keys, maxGetBatchSize)
		if err != nil {
			return nil, fmt.Errorf("Error getting batch get inputs: %s", err)
		}
		states, unprocessedKeys, err := gs.batchGetItems(batchGetInputs)
		if err != nil {
			return nil, fmt.Errorf("Error executing batch get: %s", err)
		}
		keys = unprocessedKeys
		for gatewayID, state := range states {
			output[gatewayID] = state
		}
	}
	if len(keys) > 0 {
		return output, fmt.Errorf("Remaining unprocessed keys after three attempts: %v", keys)
	}
	return output, nil
}

func (gs *dynamoGatewayStorage) batchGetItems(
	inputs []*dynamodb.BatchGetItemInput,
) (map[string]*storage.GatewayState, []map[string]*dynamodb.AttributeValue, error) {
	states := make(map[string]*storage.GatewayState)
	unprocessedKeys := []map[string]*dynamodb.AttributeValue{}
	for _, input := range inputs {
		response, err := gs.db.BatchGetItem(input)
		if err != nil {
			return nil, nil, err
		}
		for _, attributeMap := range response.Responses[GatewayViewTableName] {
			state, err := gs.decoder.GetStateFromAttributeMap(attributeMap)
			if err != nil {
				return nil, nil, err
			}
			states[state.GatewayID] = state
		}
		if batchUnprocessed, ok := response.UnprocessedKeys[GatewayViewTableName]; ok {
			unprocessedKeys = append(unprocessedKeys, batchUnprocessed.Keys...)
		}
	}
	return states, unprocessedKeys, nil
}

func (gs *dynamoGatewayStorage) UpdateOrCreateGatewayViews(networkID string, updates map[string]*storage.GatewayUpdateParams) error {
	updateInputs, err := gs.encoder.GetUpdateItemInputs(networkID, updates)
	if err != nil {
		return err
	}
	for _, updateInput := range updateInputs {
		_, err := gs.db.UpdateItem(updateInput)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *dynamoGatewayStorage) DeleteGatewayViews(networkID string, gatewayIDs []string) error {
	writes, err := gs.encoder.GetBatchedWriteRequestsForGatewayViewDeletion(networkID, gatewayIDs, MaxBatchSize)
	if err != nil {
		return err
	}
	return gs.executeBatchedWrites(writes)
}

func (gs *dynamoGatewayStorage) executeBatchedWrites(writes [][]*dynamodb.WriteRequest) error {
	for attempts := 0; len(writes) > 0 && attempts < 3; attempts++ {
		unprocessedItems := []*dynamodb.WriteRequest{}
		for _, writeBatch := range writes {
			output, err := gs.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{
					GatewayViewTableName: writeBatch,
				},
			})
			if err != nil {
				return err
			}
			unprocessedItems = append(unprocessedItems, output.UnprocessedItems[GatewayViewTableName]...)
		}
		var err error
		writes, err = dynamo.BatchWriteRequests(unprocessedItems, MaxBatchSize)
		if err != nil {
			return err
		}
	}
	if len(writes) > 0 {
		return fmt.Errorf("Unwritten writes: %v", writes)
	}
	return nil
}

func GetInitializedDynamoStorage() (storage.GatewayViewStorage, error) {
	sess, err := dynamo.GetAWSSession()
	if err != nil {
		return nil, fmt.Errorf("Error creating AWS session: %s", err)
	}
	store := NewDynamoDBGatewayStorage(dynamodb.New(sess), NewEncoder(), NewDecoder())
	if dynamo.ShouldInitTables() {
		err = store.InitTables()
		if err != nil {
			return nil, fmt.Errorf("Error initializing dynamoDB tables: %s", err)
		}
	}
	return store, nil
}
