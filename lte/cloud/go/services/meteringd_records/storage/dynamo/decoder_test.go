/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/storage/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestDecoderImpl_ProtoFromAttributeMap(t *testing.T) {
	in := map[string]*dynamodb.AttributeValue{
		"Id":            {S: aws.String("record1")},
		"Sid":           {S: aws.String("sid1")},
		"GatewayId":     {S: aws.String("gw1")},
		"NetworkId":     {S: aws.String("network")},
		"BytesTx":       {N: aws.String("1")},
		"BytesRx":       {N: aws.String("2")},
		"PktsTx":        {N: aws.String("3")},
		"PktsRx":        {N: aws.String("4")},
		"StartTime":     {N: aws.String("5")},
		"SchemaVersion": {N: aws.String("1")},
	}
	expected := &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: "record1"},
		Sid:       "sid1",
		GatewayId: "gw1",
		BytesTx:   1,
		BytesRx:   2,
		PktsTx:    3,
		PktsRx:    4,
		StartTime: &timestamp.Timestamp{Seconds: 5},
	}

	decoder := dynamo.NewDecoder()
	actual, err := decoder.ProtoFromAttributeMap(in)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
