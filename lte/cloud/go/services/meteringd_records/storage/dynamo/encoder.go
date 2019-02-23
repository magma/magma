/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo

import (
	"errors"
	"fmt"
	"time"

	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/dynamo"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Wrap time in an interface for testing
type TimeProvider interface {
	// Returns seconds since unix epoch
	UnixTimeNow() int64
}

// TimeProvider implementation which wraps golang's time.Now()
type DefaultTimeProvider struct{}

func (*DefaultTimeProvider) UnixTimeNow() int64 {
	return time.Now().Unix()
}

// The Encoder interface contains dynamoDB serialization utils for dynamoDB
// storage implementation.
type Encoder interface {
	// Transform a set of flow records into a batched array of write requests
	// for dynamoDB
	GetBatchedWriteRequestsForFlowTableUpdate(networkId string, flows []*protos.FlowRecord, batchSize int) ([][]*dynamodb.WriteRequest, error)
	GetBatchedWriteRequestsForFlowDeletion(networkId string, flowIds []string, batchSize int) ([][]*dynamodb.WriteRequest, error)
}

type encoderImpl struct {
	timeProvider TimeProvider
}

func NewEncoder(timeProvider TimeProvider) Encoder {
	return &encoderImpl{timeProvider: timeProvider}
}

// Maximum batch size for BatchWrite according to dynamoDB
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_BatchWriteItem.html
const maxBatchSize = 25

func (encoder *encoderImpl) GetBatchedWriteRequestsForFlowTableUpdate(
	networkId string,
	flows []*protos.FlowRecord,
	batchSize int,
) ([][]*dynamodb.WriteRequest, error) {
	if batchSize < 1 {
		return nil, errors.New("Batch size must be at least 1")
	}

	writeRequests, err := encoder.getWriteRequestsForFlowTable(networkId, flows)
	if err != nil {
		return nil, err
	}
	return dynamo.BatchWriteRequests(writeRequests, batchSize)
}

func (encoder *encoderImpl) getWriteRequestsForFlowTable(networkId string, flows []*protos.FlowRecord) ([]*dynamodb.WriteRequest, error) {
	ret := make([]*dynamodb.WriteRequest, 0, len(flows))
	for _, flow := range flows {
		record, err := encoder.flowRecordFromProto(networkId, flow)
		if err != nil {
			return nil, err
		}

		marshaledItem, err := dynamodbattribute.MarshalMap(record)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling flow record to dynamoDB AttributeValues: %s", err)
		}

		request := &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: marshaledItem}}
		ret = append(ret, request)
	}
	return ret, nil
}

func (encoder *encoderImpl) GetBatchedWriteRequestsForFlowDeletion(
	networkId string,
	flowIds []string,
	batchSize int,
) ([][]*dynamodb.WriteRequest, error) {
	if batchSize < 1 {
		return nil, errors.New("Batch size must be at least 1")
	}

	writeRequests, err := encoder.getWriteRequestsForFlowDeletion(networkId, flowIds)
	if err != nil {
		return nil, err
	}
	return dynamo.BatchWriteRequests(writeRequests, batchSize)
}

func (encoder *encoderImpl) getWriteRequestsForFlowDeletion(networkId string, flowIds []string) ([]*dynamodb.WriteRequest, error) {
	ret := make([]*dynamodb.WriteRequest, 0, len(flowIds))

	for _, flowId := range flowIds {
		marshaledFlowId, err := dynamodbattribute.Marshal(flowId)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling flow id to dynamoDB AttributeValues: %s", err)
		}
		marshaledKey := map[string]*dynamodb.AttributeValue{IdKeyName: marshaledFlowId}

		request := &dynamodb.WriteRequest{DeleteRequest: &dynamodb.DeleteRequest{Key: marshaledKey}}
		ret = append(ret, request)
	}
	return ret, nil
}

func (encoder *encoderImpl) flowRecordFromProto(networkId string, src *protos.FlowRecord) (*flowRecord, error) {
	dest := &flowRecord{}
	dest.Id = src.GetId().GetId()
	dest.Sid = src.GetSid()
	dest.NetworkId = networkId
	dest.SubNetworkId = fmt.Sprintf("%s%s%s", src.GetSid(), CompositeKeyDelimiter, networkId)
	dest.GatewayId = src.GetGatewayId()
	dest.BytesTx = src.GetBytesTx()
	dest.BytesRx = src.GetBytesRx()
	dest.PktsTx = src.GetPktsTx()
	dest.PktsRx = src.GetPktsRx()
	dest.SchemaVersion = SchemaVersion
	if src.GetStartTime() != nil {
		dest.StartTime = src.GetStartTime().Seconds
	}
	dest.LastUpdatedTime = encoder.timeProvider.UnixTimeNow()
	return dest, nil
}
