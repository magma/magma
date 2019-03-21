/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo_test

import (
	"fmt"
	"os"
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records/storage"
	"magma/lte/cloud/go/services/meteringd_records/storage/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

// Note that this connects to a dynamoDB local instance running on cloud VM
func TestDynamoMeteringStorage_Integration(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(getDynamoEndpoint()),
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("id", "secret", "token"),
	})
	assert.NoError(t, err)

	// Run integration tests only if TEST_MODE is enabled and dynamo is running
	if os.Getenv("TEST_MODE") != "1" {
		return
	}

	// Clear existing data
	// This is kind of fragile, but we need an end-to-end test with the local
	// dynamoDB executable and it would suck to make it a manual step to delete the DB
	// file or restart the dynamoDB service every time you want to run go test.
	dynamoClient := dynamodb.New(sess)
	_, err = dynamoClient.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(dynamo.RecordsTableName)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				// No table exists
				break
			default:
				t.Fatalf("Failed to delete dynamoDB data: %s", err)
			}
		} else {
			t.Fatalf("Failed to delete dynamoDB data: %s", err)
		}
	}

	store := dynamo.NewDynamoDBMeteringRecordsStorage(dynamoClient, dynamo.NewEncoder(&dynamo.DefaultTimeProvider{}), dynamo.NewDecoder())
	err = store.InitTables()
	assert.NoError(t, err)

	// flow1.Sid has 1 record, flow2.Sid has 2 records
	flow1 := getFlow("flow1", 1)
	flow2 := getFlow("flow2", 2)
	flow3 := getFlow("flow3", 3)
	flow3.Sid = flow2.Sid
	err = store.UpdateOrCreateRecords(
		"network",
		[]*protos.FlowRecord{flow1, flow2, flow3},
	)
	assert.NoError(t, err)

	// Read back flows
	actual, err := store.GetRecord("network", "flow1")
	assert.NoError(t, err)
	assert.Equal(t, flow1, actual)
	actual, err = store.GetRecord("network", "flow2")
	assert.NoError(t, err)
	assert.Equal(t, flow2, actual)
	actual, err = store.GetRecord("network", "flow3")
	assert.NoError(t, err)
	assert.Equal(t, flow3, actual)

	testSidIndexRead(t, store, flow1.Sid, []*protos.FlowRecord{flow1})
	testSidIndexRead(t, store, flow2.Sid, []*protos.FlowRecord{flow2, flow3})

	// Update and push new flows
	// 3 unique sids: flow1.Sid/flow5.Sid, flow2.Sid/flow3.Sid, flow4.Sid
	flow1.BytesTx = 500
	flow1.BytesRx = 600
	flow3.BytesTx = 700
	flow3.BytesRx = 800
	flow4 := getFlow("flow4", 4)
	flow5 := getFlow("flow5", 5)
	flow5.Sid = flow1.Sid
	err = store.UpdateOrCreateRecords(
		"network",
		[]*protos.FlowRecord{flow1, flow3, flow4, flow5},
	)
	assert.NoError(t, err)

	// Read back flows
	actual, err = store.GetRecord("network", "flow1")
	assert.NoError(t, err)
	assert.Equal(t, flow1, actual)
	actual, err = store.GetRecord("network", "flow2")
	assert.NoError(t, err)
	assert.Equal(t, flow2, actual)
	actual, err = store.GetRecord("network", "flow3")
	assert.NoError(t, err)
	assert.Equal(t, flow3, actual)
	actual, err = store.GetRecord("network", "flow4")
	assert.NoError(t, err)
	assert.Equal(t, flow4, actual)

	testSidIndexRead(t, store, flow1.Sid, []*protos.FlowRecord{flow1, flow5})
	testSidIndexRead(t, store, flow2.Sid, []*protos.FlowRecord{flow2, flow3})
	testSidIndexRead(t, store, flow4.Sid, []*protos.FlowRecord{flow4})

	// Create two flows with the same sid
	flow6 := getFlow("flow6", 6)
	flow7 := getFlow("flow7", 6)

	// Write flow6 and flow 7
	err = store.UpdateOrCreateRecords(
		"network",
		[]*protos.FlowRecord{flow6, flow7},
	)
	assert.NoError(t, err)
	testSidIndexRead(t, store, flow6.Sid, []*protos.FlowRecord{flow6, flow7})

	// Delete flow6 and flow7 using sid
	err = store.DeleteRecordsForSubscriber("network", flow6.Sid)
	assert.NoError(t, err)

	// Check if flow6 and flow7 are deleted
	_, err = store.GetRecordsForSubscriber("network", flow6.Sid)
	assert.NoError(t, err)
	testSidIndexRead(t, store, flow6.Sid, nil)
}

func getDynamoEndpoint() string {
	// Default to localhost:8000, but allow an override via env.
	// This way, you can run this test on a host dev machine by setting
	// endpoint to 192.168.80.10:8000
	envEndpoint := os.Getenv("DYNAMO_ENDPOINT")
	if envEndpoint == "" {
		return "http://127.0.0.1:8000"
	} else {
		return fmt.Sprintf("http://%s", envEndpoint)
	}
}

func testSidIndexRead(t *testing.T, store storage.MeteringRecordsStorage, sid string, expected []*protos.FlowRecord) {
	actualSidRecords, err := store.GetRecordsForSubscriber("network", sid)
	assert.NoError(t, err)
	assert.Equal(t, expected, actualSidRecords)
}
