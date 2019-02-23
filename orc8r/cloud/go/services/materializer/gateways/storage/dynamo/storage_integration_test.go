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
	"reflect"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/test_utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestDynamoGatewayViewStorage_Integration(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(getDynamoEndpoint()),
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("id", "secret", "token"),
	})
	assert.NoError(t, err)

	dynamoClient := dynamodb.New(sess)
	_, err = dynamoClient.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(dynamo.GatewayViewTableName)})
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

	// Create table
	store := dynamo.NewDynamoDBGatewayStorage(dynamoClient, dynamo.NewEncoder(), dynamo.NewDecoder())
	err = store.InitTables()
	assert.NoError(t, err)

	// Register configs
	registry.ClearRegistryForTesting()
	config1Manager := test_utils.NewConfig1Manager()
	registry.RegisterConfigManager(config1Manager)
	config2Manager := test_utils.NewConfig2Manager()
	registry.RegisterConfigManager(config2Manager)

	// Create gateway views
	networkID := "net1"
	gatewayStates := map[string]*storage.GatewayState{
		"gw0": getGateway0State(),
		"gw1": getGateway1State(),
		"gw2": getGateway2State(),
	}
	initialUpdates := map[string]*storage.GatewayUpdateParams{
		"gw0": convertStateToUpdateParams(gatewayStates["gw0"]),
		"gw1": convertStateToUpdateParams(gatewayStates["gw1"]),
		"gw2": convertStateToUpdateParams(gatewayStates["gw2"]),
	}
	err = store.UpdateOrCreateGatewayViews(networkID, initialUpdates)
	assert.NoError(t, err)

	// Test GetGatewayViewsForNetwork
	actualStates, err := store.GetGatewayViewsForNetwork(networkID)
	assert.NoError(t, err)
	verifyAllViews(t, gatewayStates, actualStates)

	// Test GetGatewayViewsForNetwork with invalid network
	badIDOutput, err := store.GetGatewayViewsForNetwork("badid")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(badIDOutput))

	// Test GetGatewayViews
	queryGatewayIDs := []string{"gw0", "gw2", "badgw"}
	actualStates, err = store.GetGatewayViews(networkID, queryGatewayIDs)
	assert.NoError(t, err)
	// Assert bad gateway not present
	assert.Equal(t, 2, len(actualStates))
	verifyGatewayView(t, gatewayStates["gw0"], actualStates)
	verifyGatewayView(t, gatewayStates["gw2"], actualStates)
	_, ok := actualStates["gw1"]
	assert.False(t, ok, "gw1 returned in GetGatewayViews query which did not include it")

	// Test update gateway views
	updateParams := map[string]*storage.GatewayUpdateParams{
		"gw0": {
			NewConfig: map[string]interface{}{
				"Config2": &test_utils.Conf2{
					Value1: []string{"abc", "def"},
					Value2: 1234,
				},
			},
			ConfigsToDelete: []string{"Config1"},
			Offset:          2,
		},
		"gw2": {
			NewConfig: map[string]interface{}{
				"Config2": &test_utils.Conf2{
					Value1: []string{"new", "strings"},
					Value2: 5678,
				},
			},
			NewStatus: &protos.GatewayStatus{
				Time:    91011,
				Checkin: nil,
			},
			Offset: 2,
		},
	}
	gatewayStates["gw0"].Config["Config2"] = updateParams["gw0"].NewConfig["Config2"]
	delete(gatewayStates["gw0"].Config, "Config1")
	gatewayStates["gw0"].Offset = updateParams["gw0"].Offset
	gatewayStates["gw2"].Config["Config2"] = updateParams["gw2"].NewConfig["Config2"]
	gatewayStates["gw2"].Status = updateParams["gw2"].NewStatus
	gatewayStates["gw2"].Offset = updateParams["gw2"].Offset

	newGateway := &storage.GatewayState{
		GatewayID: "gw3",
		Config:    make(map[string]interface{}),
		Status: &protos.GatewayStatus{
			Time: 12345,
			Checkin: &protos.CheckinRequest{
				KernelVersionsInstalled: []string{},
			},
		},
		Record: &magmadprotos.AccessGatewayRecord{},
		Offset: 2,
	}
	gatewayStates["gw3"] = newGateway
	updateParams["gw3"] = convertStateToUpdateParams(newGateway)

	err = store.UpdateOrCreateGatewayViews(networkID, updateParams)
	assert.NoError(t, err)
	actualStates, err = store.GetGatewayViewsForNetwork(networkID)
	assert.NoError(t, err)
	verifyAllViews(t, gatewayStates, actualStates)

	idsToDelete := []string{"gw1", "gw2"}
	err = store.DeleteGatewayViews(networkID, idsToDelete)
	assert.NoError(t, err)
	actualStates, err = store.GetGatewayViewsForNetwork(networkID)
	assert.NoError(t, err)
	_, ok = actualStates["gw0"]
	assert.True(t, ok, "Gateway 0 not found after delete of gateways 1 and 2")
	_, ok = actualStates["gw1"]
	assert.False(t, ok, "Gateway 1 view found after it was deleted")
	_, ok = actualStates["gw2"]
	assert.False(t, ok, "Gateway 2 view found after it was deleted")
}

func verifyAllViews(
	t *testing.T,
	expectedStates map[string]*storage.GatewayState,
	actualStates map[string]*storage.GatewayState,
) {
	for _, expectedState := range expectedStates {
		verifyGatewayView(t, expectedState, actualStates)
	}
}

func verifyGatewayView(
	t *testing.T,
	expectedState *storage.GatewayState,
	actualStates map[string]*storage.GatewayState,
) {
	actualState, ok := actualStates[expectedState.GatewayID]
	assert.True(t, ok, fmt.Sprintf("Did not find %s in map from GetGatewayViewsForNetwork", expectedState.GatewayID))
	assert.True(t, reflect.DeepEqual(expectedState, actualState), "Gateway states not equal")
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

func getGateway0State() *storage.GatewayState {
	return &storage.GatewayState{
		GatewayID: "gw0",
		Config: map[string]interface{}{
			"Config1": &test_utils.Conf1{
				Value1: 13,
				Value2: "TEST",
				Value3: []byte{1, 3, 5, 7},
			},
		},
		Status: &protos.GatewayStatus{
			Time: 1,
			Checkin: &protos.CheckinRequest{
				GatewayId:               "gw0",
				MagmaPkgVersion:         "v1",
				Status:                  nil,
				SystemStatus:            nil,
				KernelVersionsInstalled: []string{},
			},
		},
		Record: &magmadprotos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: "hwid0"},
			Name: "Gateway 0",
			Key:  &protos.ChallengeKey{},
			Ip:   "127.0.0.1",
			Port: 1000,
		},
		Offset: 1,
	}
}

func getGateway1State() *storage.GatewayState {
	return &storage.GatewayState{
		GatewayID: "gw1",
		Config: map[string]interface{}{
			"Config1": &test_utils.Conf1{
				Value1: 11,
				Value2: "TESTSTR",
				Value3: []byte{3, 2, 1},
			},
			"Config2": &test_utils.Conf2{
				Value1: []string{"Hello", "World", "!"},
				Value2: 29,
			},
		},
		Status: &protos.GatewayStatus{
			Time: 100,
			Checkin: &protos.CheckinRequest{
				GatewayId:       "gw1",
				MagmaPkgVersion: "v2",
				Status: &protos.ServiceStatus{
					Meta: map[string]string{
						"Hello": "World",
					},
				},
				SystemStatus: &protos.SystemStatus{
					Time:      100,
					CpuUser:   1,
					CpuSystem: 3,
					MemTotal:  1000,
				},
				KernelVersionsInstalled: []string{},
			},
		},
		Record: &magmadprotos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: "hwid1"},
			Name: "Gateway 1",
			Key: &protos.ChallengeKey{
				KeyType: protos.ChallengeKey_SOFTWARE_ECDSA_SHA256,
				Key:     []byte{1, 2, 3, 4, 5},
			},
			Ip:   "127.0.0.2",
			Port: 2000,
		},
		Offset: 2,
	}
}

func getGateway2State() *storage.GatewayState {
	return &storage.GatewayState{
		GatewayID: "gw2",
		Config: map[string]interface{}{
			"Config2": &test_utils.Conf2{
				Value1: []string{"TEST", "STRING", "ARRAY"},
				Value2: 103,
			},
		},
		Status: &protos.GatewayStatus{},
		Record: &magmadprotos.AccessGatewayRecord{},
		Offset: 3,
	}
}

func convertStateToUpdateParams(state *storage.GatewayState) *storage.GatewayUpdateParams {
	return &storage.GatewayUpdateParams{
		NewConfig: state.Config,
		NewStatus: state.Status,
		NewRecord: state.Record,
		Offset:    state.Offset,
	}
}
