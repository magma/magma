/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo

import (
	"magma/lte/cloud/go/protos"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// The Decoder interface contains deserialization utils to go from dynamoDB
// types to concrete data types for the dynamoDB implementation of the storage
// interface.
type Decoder interface {
	// Given a dynamoDB attribute map, deserialize into a FlowRecord proto
	// structure
	ProtoFromAttributeMap(map[string]*dynamodb.AttributeValue) (*protos.FlowRecord, error)
}

type decoderImpl struct{}

func NewDecoder() Decoder {
	return &decoderImpl{}
}

func (decoder *decoderImpl) ProtoFromAttributeMap(src map[string]*dynamodb.AttributeValue) (*protos.FlowRecord, error) {
	record := &flowRecord{}
	err := dynamodbattribute.UnmarshalMap(src, record)
	if err != nil {
		return nil, err
	}
	return protoFromFlowRecord(record)
}

func protoFromFlowRecord(src *flowRecord) (*protos.FlowRecord, error) {
	return &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: src.Id},
		Sid:       src.Sid,
		GatewayId: src.GatewayId,
		BytesRx:   src.BytesRx,
		BytesTx:   src.BytesTx,
		PktsRx:    src.PktsRx,
		PktsTx:    src.PktsTx,
		StartTime: &timestamp.Timestamp{Seconds: src.StartTime},
	}, nil
}
