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

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/storage/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

type mockTimeProvider struct {
	fakeTime int64
}

func (m *mockTimeProvider) UnixTimeNow() int64 {
	return m.fakeTime
}

func TestEncoderImpl_GetBatchedWriteRequestsForFlowTableUpdate(t *testing.T) {
	encoder := dynamo.NewEncoder(&mockTimeProvider{fakeTime: 123})

	// 1 batch
	in := []*protos.FlowRecord{
		getFlow("flow1", 1),
		getFlow("flow2", 2),
	}
	expected := [][]*dynamodb.WriteRequest{getExpectedWriteRequests(t, "network", in)}
	actual, err := encoder.GetBatchedWriteRequestsForFlowTableUpdate("network", in, 5)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 2 batch
	in = []*protos.FlowRecord{
		getFlow("flow1", 1),
		getFlow("flow2", 2),
		getFlow("flow2", 3),
		getFlow("flow2", 4),
		getFlow("flow2", 5),
	}
	expectedUnbatched := getExpectedWriteRequests(t, "network", in)
	expected = [][]*dynamodb.WriteRequest{expectedUnbatched[:3], expectedUnbatched[3:]}
	actual, err = encoder.GetBatchedWriteRequestsForFlowTableUpdate("network", in, 3)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Penny and dime (empty input)
	in = []*protos.FlowRecord{}
	expected = [][]*dynamodb.WriteRequest{}
	actual, err = encoder.GetBatchedWriteRequestsForFlowTableUpdate("network", in, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func getFlow(recordId string, seed int) *protos.FlowRecord {
	return &protos.FlowRecord{
		Id:        &protos.FlowRecord_ID{Id: recordId},
		Sid:       fmt.Sprintf("sid%d", seed+1),
		GatewayId: fmt.Sprintf("gw%d", seed+2),
		BytesTx:   uint64(seed + 3),
		BytesRx:   uint64(seed + 4),
		PktsTx:    uint64(seed + 5),
		PktsRx:    uint64(seed + 6),
		StartTime: &timestamp.Timestamp{Seconds: int64(seed + 7)},
	}
}

func getExpectedWriteRequests(t *testing.T, networkId string, flows []*protos.FlowRecord) []*dynamodb.WriteRequest {
	ret := make([]*dynamodb.WriteRequest, 0, len(flows))
	for _, flow := range flows {
		ret = append(ret, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"Id":              {S: aws.String(flow.Id.Id)},
					"Sid":             {S: aws.String(flow.Sid)},
					"NetworkId":       {S: aws.String(networkId)},
					"SubNetworkId":    {S: aws.String(fmt.Sprintf("%s|%s", flow.Sid, networkId))},
					"GatewayId":       {S: aws.String(flow.GatewayId)},
					"BytesTx":         {N: aws.String(fmt.Sprintf("%d", flow.BytesTx))},
					"BytesRx":         {N: aws.String(fmt.Sprintf("%d", flow.BytesRx))},
					"PktsTx":          {N: aws.String(fmt.Sprintf("%d", flow.PktsTx))},
					"PktsRx":          {N: aws.String(fmt.Sprintf("%d", flow.PktsRx))},
					"StartTime":       {N: aws.String(fmt.Sprintf("%d", flow.StartTime.Seconds))},
					"LastUpdatedTime": {N: aws.String("123")},
					"SchemaVersion":   {N: aws.String("1")},
				},
			},
		})
	}
	return ret
}
