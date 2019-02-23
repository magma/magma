/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package dynamo contains an implementation for the metering records storage
// interface backed by dynamoDB for persistence.
package dynamo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	RecordsTableName = "meteringd_records"
	SidIndexName     = "sid_idx"
	SchemaVersion    = 1

	SidKeyName = "SubNetworkId"
	IdKeyName  = "Id"
)

type dynamoMeteringStorage struct {
	db      dynamodbiface.DynamoDBAPI
	encoder Encoder
	decoder Decoder
}

func NewDynamoDBMeteringRecordsStorage(db dynamodbiface.DynamoDBAPI, encoder Encoder, decoder Decoder) storage.MeteringRecordsStorage {
	return &dynamoMeteringStorage{db: db, encoder: encoder, decoder: decoder}
}

/*
	Create the metering record tables if they don't already exist.
	The metering records table schema looks like this:

	Field			|	Type	|	Key Type	|	 Notes
	=====================================================
	Id				|	S		|	hash		|
	Sid				|	S		|				|
	GatewayId		|	S		|				|
	NetworkId		|	S		|				|
	SubNetworkId	|	S		|				|	Internal: composite key "Sid|NetworkId" to use as a partition key for sid secondary index
	Match			|	B		|				|
	BytesTx			|	N		|				|
	BytesRx			|	N		|				|
	PktsTx			|	N		|				|
	PktsRx			|	N		|				|
	LastUpdatedTime	|	N		|				|	Timestamp in seconds since epoch
	StartTime		|	N		|				|	Timestamp in seconds since epoch
	SchemaVersion	|	N		|				|	Internal: for future backcompat

	We use the flow record's unique UUID `Id` as the partition key to
	utilize throughput most efficiently. Note that there is no range key in
	this schema as these UUIDs are probabilistically guaranteed to be unique.

	There is 1 global secondary index on this table on the `SidNetworkId` field
	to accommodate queries for a subscriber's usage. The range key on this
	index is the flow UUID to achieve a unique primary key.

	IMPORTANT: DO NOT MODIFY THIS METHOD BEFORE MAKING A CORRESPONDING CHANGE
	IN THE TERRAFORM DYNAMODB RESOURCE. On staging/prod, this method isn't
	called because we rely on Terraform to manage the DynamoDB resources.

	This is only called on local dev VM in order to set up the local
	dynamoDB executable. Changing this method without changing the terraform
	resource will have no effect on staging or prod. If you change the
	terraform resource without modifying this method, your dev VM schema will
	be outdated.
*/
func (ms *dynamoMeteringStorage) InitTables() error {
	sidIndex := &dynamodb.GlobalSecondaryIndex{
		IndexName: aws.String(SidIndexName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(SidKeyName),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String(IdKeyName),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		Projection: &dynamodb.Projection{
			ProjectionType: aws.String(dynamodb.ProjectionTypeAll),
		},
		// Ignored by local dynamoDB
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(100),
			WriteCapacityUnits: aws.Int64(100),
		},
	}

	_, err := ms.db.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(RecordsTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(IdKeyName),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String(SidKeyName),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(IdKeyName),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{sidIndex},
		// Ignored by local dynamoDB
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
				return fmt.Errorf("Error creating metering records tables: %s", err)
			}
		} else {
			return fmt.Errorf("Error creating metering records tables: %s", err)
		}
	}

	return nil
}

// This is implemented here using dynamoDB batched writes.
// IMPORTANT: batches will fail silently, because this should only be called
// from gateways when they sync flow stats, which means the same items will
// get re-synced on the next RPC call from the gateway in 5 seconds.
// If there eventually arises a need for a strongly consistent/reliable
// put call, we can add that to the interface.
func (ms *dynamoMeteringStorage) UpdateOrCreateRecords(networkId string, flows []*protos.FlowRecord) error {
	if flows == nil {
		return nil
	}
	writes, err := ms.encoder.GetBatchedWriteRequestsForFlowTableUpdate(networkId, flows, maxBatchSize)
	if err != nil {
		return err
	}

	for _, writeBatch := range writes {
		// Ignore unprocessed items because they will just get re-written  on
		// the next sync from the gateway
		_, err := ms.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				RecordsTableName: writeBatch,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (ms *dynamoMeteringStorage) GetRecord(networkId string, recordId string) (*protos.FlowRecord, error) {
	keyAttr, err := dynamodbattribute.Marshal(recordId)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal record ID to attribute value: %s", err)
	}

	res, err := ms.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(RecordsTableName),
		Key:       map[string]*dynamodb.AttributeValue{IdKeyName: keyAttr},
	})
	if err != nil {
		return nil, err
	}

	return ms.decoder.ProtoFromAttributeMap(res.Item)
}

func (ms *dynamoMeteringStorage) GetRecordsForSubscriber(networkId string, sid string) ([]*protos.FlowRecord, error) {
	// Build the key condition and do a paginated query
	keyConditionBuilder := expression.KeyEqual(
		expression.Key(SidKeyName),
		expression.Value(fmt.Sprintf("%s%s%s", sid, CompositeKeyDelimiter, networkId)),
	)
	keyCondition, err := expression.NewBuilder().WithKeyCondition(keyConditionBuilder).Build()
	if err != nil {
		return nil, err
	}
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(RecordsTableName),
		IndexName:                 aws.String(SidIndexName),
		KeyConditionExpression:    keyCondition.KeyCondition(),
		ExpressionAttributeNames:  keyCondition.Names(),
		ExpressionAttributeValues: keyCondition.Values(),
	}

	var ret []*protos.FlowRecord
	var callbackErr error
	err = ms.db.QueryPages(queryInput, func(result *dynamodb.QueryOutput, lastPage bool) bool {
		decodedPageItems, err := ms.getProtosFromPageItems(result.Items)
		if err != nil {
			callbackErr = err
			return false
		}
		ret = append(ret, decodedPageItems...)
		return !lastPage
	})

	if callbackErr != nil {
		return nil, callbackErr
	}
	return ret, err
}

func (ms *dynamoMeteringStorage) DeleteRecordsForSubscriber(networkId string, sid string) error {
	flowIds, err := ms.getIdsForFlows(networkId, sid)
	if err != nil {
		return err
	}

	// Delete flow records using flow ids
	writes, err := ms.encoder.GetBatchedWriteRequestsForFlowDeletion(networkId, flowIds, maxBatchSize)
	if err != nil {
		return err
	}

	// Execute all batches
	unprocessedItems := []map[string][]*dynamodb.WriteRequest{}
	for _, writeBatch := range writes {
		output, err := ms.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				RecordsTableName: writeBatch,
			},
		})
		if err != nil {
			return err
		}
		if len(output.UnprocessedItems) > 0 {
			unprocessedItems = append(unprocessedItems, output.UnprocessedItems)
		}
	}

	return getUnprocessedItemsError(unprocessedItems)
}

// Get all flow ids of the subscriber
func (ms *dynamoMeteringStorage) getIdsForFlows(networkId string, sid string) ([]string, error) {
	flowRecords, err := ms.GetRecordsForSubscriber(networkId, sid)
	if err != nil {
		return nil, err
	}
	var flowIds []string
	for _, flowRecord := range flowRecords {
		flowIds = append(flowIds, flowRecord.GetId().GetId())
	}
	return flowIds, nil
}

func (ms *dynamoMeteringStorage) getProtosFromPageItems(items []map[string]*dynamodb.AttributeValue) ([]*protos.FlowRecord, error) {
	ret := make([]*protos.FlowRecord, 0, len(items))
	for _, attrValueMap := range items {
		val, err := ms.decoder.ProtoFromAttributeMap(attrValueMap)
		if err != nil {
			return nil, err
		}
		ret = append(ret, val)
	}
	return ret, nil
}

// Aggregate all unprocessed items as an error message
func getUnprocessedItemsError(unprocessedItems []map[string][]*dynamodb.WriteRequest) error {
	if len(unprocessedItems) == 0 {
		return nil
	}

	ret := bytes.NewBufferString("Unprocessed Items:")
	for _, unprocessedItem := range unprocessedItems {
		marshaledUnprocessedItem, err := json.Marshal(unprocessedItem)
		if err != nil {
			return err
		}
		ret.WriteString("\n")
		ret.Write(marshaledUnprocessedItem)
	}

	return errors.New(ret.String())
}
