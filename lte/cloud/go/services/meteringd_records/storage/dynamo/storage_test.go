/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/storage"
	"magma/lte/cloud/go/services/meteringd_records/storage/dynamo"
	"magma/lte/cloud/go/services/meteringd_records/storage/dynamo/mocks"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDynamoMeteringStorage_GetRecord(t *testing.T) {
	// Mock setup - 1 successful return and 1 error return
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// Successful return
	mockAttributes := getMockAttributeMap("attr1", "attr2")
	mockDB.On("GetItem", &dynamodb.GetItemInput{
		TableName: aws.String("meteringd_records"),
		Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String("record1")}},
	}).Return(&dynamodb.GetItemOutput{Item: mockAttributes}, nil)
	// Error return
	mockDB.On("GetItem", &dynamodb.GetItemInput{
		TableName: aws.String("meteringd_records"),
		Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String("record2")}},
	}).Return(nil, errors.New("Mock error"))

	mockFlowRecord := &protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: "record1"}}
	mockDecoder.On("ProtoFromAttributeMap", mock.Anything).Return(mockFlowRecord, nil)

	// Start test
	store := newStorage(mockDB, mockEncoder, mockDecoder)
	actual, err := store.GetRecord("network", "record1")
	assert.NoError(t, err)
	assert.Equal(t, actual, mockFlowRecord)

	_, err = store.GetRecord("network", "record2")
	assert.Error(t, err)
	assert.Equal(t, "Mock error", err.Error())

	mockDB.AssertNumberOfCalls(t, "GetItem", 2)
	mockDB.AssertExpectations(t)
	mockDecoder.AssertNumberOfCalls(t, "ProtoFromAttributeMap", 1)
	mockDecoder.AssertExpectations(t)
	mockEncoder.AssertExpectations(t)
}

func TestDynamoMeteringStorage_UpdateOrCreateRecords(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// 1 happy path, 1 encoder error, 1 dynamo API error
	flowTbl1 := getFlows("record1", "record2")
	writeRequests1 := getMockWriteRequests([][]string{{"attr1", "attr2"}, {"attr3"}})
	flowTbl2 := getFlows("record3")
	flowTbl3 := getFlows("record4", "record5")
	writeRequests3 := getMockWriteRequests([][]string{{"attr4"}, {"attr5"}, {"attr6"}})
	mockEncoder.
		On("GetBatchedWriteRequestsForFlowTableUpdate", "network", flowTbl1, 25).
		Return(writeRequests1, nil)
	mockEncoder.
		On("GetBatchedWriteRequestsForFlowTableUpdate", "network", flowTbl2, 25).
		Return(nil, errors.New("Mock encoder error"))
	mockEncoder.
		On("GetBatchedWriteRequestsForFlowTableUpdate", "network", flowTbl3, 25).
		Return(writeRequests3, nil)
	mockDB.On("BatchWriteItem", getBatchWriteItemInput(writeRequests1[0])).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	mockDB.On("BatchWriteItem", getBatchWriteItemInput(writeRequests1[1])).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	mockDB.On("BatchWriteItem", getBatchWriteItemInput(writeRequests3[0])).Return(&dynamodb.BatchWriteItemOutput{}, errors.New("BatchWriteItem error"))

	// Run test cases
	store := newStorage(mockDB, mockEncoder, mockDecoder)
	err := store.UpdateOrCreateRecords("network", flowTbl1)
	assert.NoError(t, err)
	err = store.UpdateOrCreateRecords("network", flowTbl2)
	assert.Error(t, err)
	assert.Equal(t, "Mock encoder error", err.Error())
	err = store.UpdateOrCreateRecords("network", flowTbl3)
	assert.Error(t, err)
	assert.Equal(t, "BatchWriteItem error", err.Error())

	mockDB.AssertNumberOfCalls(t, "BatchWriteItem", 3)
	mockDB.AssertExpectations(t)
	mockEncoder.AssertExpectations(t)
	mockDecoder.AssertExpectations(t)
}

func TestDynamoMeteringStorage_GetRecordsForSubscriber(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// 1 happy path, 1 dynamo API error, 1 decoder error
	query1 := getSubscriberIndexQueryInput("network", "sid1")
	query2 := getSubscriberIndexQueryInput("network", "sid2")
	query3 := getSubscriberIndexQueryInput("network", "sid3")
	query1Pages := [][]map[string]*dynamodb.AttributeValue{
		{getMockAttributeMap("attr1", "attr2"), getMockAttributeMap("attr3")},
		{getMockAttributeMap("attr4", "attr5", "attr6")},
	}
	query1Items := flattenNestedQueryPages(query1Pages)
	query1MockResult := &protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: "fake record"}}
	query3Pages := [][]map[string]*dynamodb.AttributeValue{
		{getMockAttributeMap("attr8")},
		{getMockAttributeMap("attr9", "attr10")},
		{getMockAttributeMap("attr11")},
	}
	query3Items := flattenNestedQueryPages(query3Pages)
	query3MockResult := &protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: "fake record 2"}}

	mockDB.On("QueryPages", query1, mock.Anything).Return(getMockQueryPagesImpl(query1Pages))
	mockDecoder.On("ProtoFromAttributeMap", mock.MatchedBy(getQueryResultMatcherFn(query1Items))).Return(query1MockResult, nil)

	mockDB.On("QueryPages", query2, mock.Anything).Return(errors.New("Mock dynamoDB error"))

	mockDB.On("QueryPages", query3, mock.Anything).Return(getMockQueryPagesImpl(query3Pages))
	mockDecoder.On("ProtoFromAttributeMap", query3Items[1]).Return(nil, errors.New("Mock decoder error"))
	mockDecoder.On("ProtoFromAttributeMap", mock.MatchedBy(getQueryResultMatcherFn(query3Items))).Return(query3MockResult, nil)

	// Run test cases
	store := newStorage(mockDB, mockEncoder, mockDecoder)
	actual, err := store.GetRecordsForSubscriber("network", "sid1")
	assert.NoError(t, err)
	// Side note: aren't closures cool?
	assert.Equal(t, []*protos.FlowRecord{query1MockResult, query1MockResult, query1MockResult}, actual)

	_, err = store.GetRecordsForSubscriber("network", "sid2")
	assert.Error(t, err)
	assert.Equal(t, "Mock dynamoDB error", err.Error())

	_, err = store.GetRecordsForSubscriber("network", "sid3")
	assert.Error(t, err)
	assert.Equal(t, "Mock decoder error", err.Error())

	mockDB.AssertExpectations(t)
	mockDB.AssertNumberOfCalls(t, "QueryPages", 3)
	mockDecoder.AssertExpectations(t)
	// For query 3, decode should have been called on the first 2 pages but not the 3rd
	mockDecoder.AssertNumberOfCalls(t, "ProtoFromAttributeMap", 5)
}

func TestDynamoMeteringStorage_DeleteRecordsForSubscriber(t *testing.T) {
	// Mock setup
	mockDB := &mocks.DynamoDBAPI{}
	mockEncoder := &mocks.Encoder{}
	mockDecoder := &mocks.Decoder{}

	// 1 happy path, 1 path with BatchWriteItem error, 1 path with 2 unprocessed item returned
	deleteRequests1 := getMockDeleteRequests([][]string{{"fake record"}})
	deleteRequests2 := getMockDeleteRequests([][]string{{"fake record 2"}})
	deleteRequests3 := getMockDeleteRequests([][]string{{"fake record 3"}, {"fake record 4"}})

	mockEncoder.
		On("GetBatchedWriteRequestsForFlowDeletion", "network", []string{"fake record"}, 25).
		Return(deleteRequests1, nil)
	mockEncoder.
		On("GetBatchedWriteRequestsForFlowDeletion", "network", []string{"fake record 2"}, 25).
		Return(deleteRequests2, nil)
	mockEncoder.
		On("GetBatchedWriteRequestsForFlowDeletion", "network", []string{"fake record 3", "fake record 4"}, 25).
		Return(deleteRequests3, nil)

	mockDB.On("BatchWriteItem", getBatchWriteItemInput(deleteRequests1[0])).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	mockDB.On("BatchWriteItem", getBatchWriteItemInput(deleteRequests2[0])).Return(&dynamodb.BatchWriteItemOutput{}, errors.New("BatchWriteItem error"))
	mockDB.
		On("BatchWriteItem", getBatchWriteItemInput(deleteRequests3[0])).
		Return(&dynamodb.BatchWriteItemOutput{UnprocessedItems: map[string][]*dynamodb.WriteRequest{"temp": deleteRequests3[0]}}, nil)
	mockDB.
		On("BatchWriteItem", getBatchWriteItemInput(deleteRequests3[1])).
		Return(&dynamodb.BatchWriteItemOutput{UnprocessedItems: map[string][]*dynamodb.WriteRequest{"temp": deleteRequests3[1]}}, nil)

	query1 := getSubscriberIndexQueryInput("network", "sid1")
	query1Pages := [][]map[string]*dynamodb.AttributeValue{
		{getMockAttributeMap("attr1")},
	}
	query1Items := flattenNestedQueryPages(query1Pages)
	query1MockResult := &protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: "fake record"}}
	mockDB.On("QueryPages", query1, mock.Anything).Return(getMockQueryPagesImpl(query1Pages))
	mockDecoder.On("ProtoFromAttributeMap", mock.MatchedBy(getQueryResultMatcherFn(query1Items))).Return(query1MockResult, nil)

	query2 := getSubscriberIndexQueryInput("network", "sid2")
	query2Pages := [][]map[string]*dynamodb.AttributeValue{
		{getMockAttributeMap("attr2")},
	}
	query2Items := flattenNestedQueryPages(query2Pages)
	query2MockResult := &protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: "fake record 2"}}
	mockDB.On("QueryPages", query2, mock.Anything).Return(getMockQueryPagesImpl(query2Pages))
	mockDecoder.On("ProtoFromAttributeMap", mock.MatchedBy(getQueryResultMatcherFn(query2Items))).Return(query2MockResult, nil)

	query3 := getSubscriberIndexQueryInput("network", "sid3")
	query3Pages := [][]map[string]*dynamodb.AttributeValue{
		{getMockAttributeMap("attr3")},
		{getMockAttributeMap("attr4")},
	}
	query3Items := flattenNestedQueryPages(query3Pages)
	query3MockResult := []interface{}{
		&protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: "fake record 3"}},
		&protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: "fake record 4"}},
	}
	mockDB.On("QueryPages", query3, mock.Anything).Return(getMockQueryPagesImpl(query3Pages))
	mockDecoder.
		On("ProtoFromAttributeMap", mock.MatchedBy(getQueryResultExactMatcherFn([]map[string]*dynamodb.AttributeValue{query3Items[0]}))).
		Return(query3MockResult[0], nil)
	mockDecoder.
		On("ProtoFromAttributeMap", mock.MatchedBy(getQueryResultExactMatcherFn([]map[string]*dynamodb.AttributeValue{query3Items[1]}))).
		Return(query3MockResult[1], nil)

	// Run test cases
	store := newStorage(mockDB, mockEncoder, mockDecoder)
	err := store.DeleteRecordsForSubscriber("network", "sid1")
	assert.NoError(t, err)

	err = store.DeleteRecordsForSubscriber("network", "sid2")
	assert.Error(t, err)
	assert.Equal(t, "BatchWriteItem error", err.Error())

	err = store.DeleteRecordsForSubscriber("network", "sid3")
	marshaledUnprocessedItem1, _ := json.Marshal(map[string][]*dynamodb.WriteRequest{"temp": deleteRequests3[0]})
	marshaledUnprocessedItem2, _ := json.Marshal(map[string][]*dynamodb.WriteRequest{"temp": deleteRequests3[1]})
	unprocessedItemsErrMsg := "Unprocessed Items:\n" + string(marshaledUnprocessedItem1) + "\n" + string(marshaledUnprocessedItem2)
	assert.Error(t, err)
	assert.Equal(t, errors.New(unprocessedItemsErrMsg), err)

	mockDB.AssertNumberOfCalls(t, "BatchWriteItem", 4)
	mockDB.AssertExpectations(t)
	mockEncoder.AssertExpectations(t)
	mockEncoder.AssertNumberOfCalls(t, "GetBatchedWriteRequestsForFlowDeletion", 3)
	mockDecoder.AssertExpectations(t)
}

func newStorage(mockDB *mocks.DynamoDBAPI, mockEncoder *mocks.Encoder, mockDecoder *mocks.Decoder) storage.MeteringRecordsStorage {
	return dynamo.NewDynamoDBMeteringRecordsStorage(mockDB, mockEncoder, mockDecoder)
}

func getMockAttributeMap(attrNames ...string) map[string]*dynamodb.AttributeValue {
	ret := make(map[string]*dynamodb.AttributeValue, len(attrNames))
	for _, attr := range attrNames {
		ret[attr] = &dynamodb.AttributeValue{S: aws.String(attr + "_val")}
	}
	return ret
}

func getFlows(flowIds ...string) []*protos.FlowRecord {
	ret := make([]*protos.FlowRecord, 0, len(flowIds))
	for _, flowId := range flowIds {
		ret = append(ret, &protos.FlowRecord{Id: &protos.FlowRecord_ID{Id: flowId}})
	}
	return ret
}

func getMockWriteRequests(attrNames [][]string) [][]*dynamodb.WriteRequest {
	ret := make([][]*dynamodb.WriteRequest, 0, len(attrNames))
	for _, attrList := range attrNames {
		req := &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: getMockAttributeMap(attrList...)}}
		ret = append(ret, []*dynamodb.WriteRequest{req})
	}
	return ret
}

func getMockDeleteRequests(attrNames [][]string) [][]*dynamodb.WriteRequest {
	ret := make([][]*dynamodb.WriteRequest, 0, len(attrNames))
	for _, attrList := range attrNames {
		req := &dynamodb.WriteRequest{DeleteRequest: &dynamodb.DeleteRequest{Key: getMockAttributeMap(attrList...)}}
		ret = append(ret, []*dynamodb.WriteRequest{req})
	}
	return ret
}

func getBatchWriteItemInput(reqs []*dynamodb.WriteRequest) *dynamodb.BatchWriteItemInput {
	return &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"meteringd_records": reqs,
		},
	}
}

func getSubscriberIndexQueryInput(networkId string, sid string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:                 aws.String("meteringd_records"),
		IndexName:                 aws.String("sid_idx"),
		KeyConditionExpression:    aws.String("#0 = :0"),
		ExpressionAttributeNames:  map[string]*string{"#0": aws.String("SubNetworkId")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":0": {S: aws.String(sid + "|" + networkId)}},
	}
}

func flattenNestedQueryPages(in [][]map[string]*dynamodb.AttributeValue) []map[string]*dynamodb.AttributeValue {
	var ret []map[string]*dynamodb.AttributeValue
	for _, itemList := range in {
		ret = append(ret, itemList...)
	}
	return ret
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

// return true only when the result is matched exactly
func getQueryResultExactMatcherFn(mockQueryResults []map[string]*dynamodb.AttributeValue) func(map[string]*dynamodb.AttributeValue) bool {
	return func(in map[string]*dynamodb.AttributeValue) bool {
		for _, queryResult := range mockQueryResults {
			if !reflect.DeepEqual(queryResult, in) {
				return false
			}
		}
		return true
	}
}
