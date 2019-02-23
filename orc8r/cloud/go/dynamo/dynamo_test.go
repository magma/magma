/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/dynamo"
)

func TestBatchWriteRequestsBatchSizeMultiple(t *testing.T) {
	testBatchWriteRequests(t, 6, 2, 3)
}

func TestBatchWriteRequestsWithRemainder(t *testing.T) {
	testBatchWriteRequests(t, 5, 2, 3)
}

func TestBatchWriteRequestsGreaterBatch(t *testing.T) {
	testBatchWriteRequests(t, 2, 3, 1)
}

func TestBatchWriteRequestsBatchSizeOne(t *testing.T) {
	testBatchWriteRequests(t, 3, 1, 3)
}

func TestBatchWriteRequestsBatchSizeZero(t *testing.T) {
	_, err := dynamo.BatchWriteRequests(nil, 0)
	assert.EqualError(t, err, "Batch size must be at least 1")
}

func testBatchWriteRequests(t *testing.T, numRequests int, batchSize int, expectedBatchCount int) {
	requests := make([]*dynamodb.WriteRequest, numRequests)
	for i := 0; i < numRequests; i++ {
		attrValue, err := dynamodbattribute.Marshal(i)
		assert.NoError(t, err)
		requests[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"item": attrValue,
				},
			},
		}
	}
	batchedRequests, err := dynamo.BatchWriteRequests(requests, batchSize)
	assert.NoError(t, err)
	assert.Equal(t, expectedBatchCount, len(batchedRequests))
	count := 0
	for _, batch := range batchedRequests {
		assert.True(t, len(batch) <= batchSize, "Batch size greater than input batch size")
		for _, request := range batch {
			var contents int
			err := dynamodbattribute.Unmarshal(request.PutRequest.Item["item"], &contents)
			assert.NoError(t, err)
			assert.Equal(t, count, contents)
			count++
		}
	}
}
