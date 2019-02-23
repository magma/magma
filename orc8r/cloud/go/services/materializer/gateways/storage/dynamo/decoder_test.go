/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo_test

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/test_utils"
)

func TestDecoderStateFromAttrMap(t *testing.T) {
	// Set up registry
	registry.ClearRegistryForTesting()
	configManager := test_utils.NewConfig1Manager()
	registry.RegisterConfigManager(configManager)
	config1 := &test_utils.Conf1{
		Value1: 4,
		Value2: "Hello, world!",
		Value3: []byte{3, 4, 5},
	}
	config1Bytes, err := registry.MarshalConfig(configManager.GetConfigType(), config1)
	assert.NoError(t, err)
	config := map[string]interface{}{
		configManager.GetConfigType(): config1,
	}
	status := test_utils.GetDefaultStatus(t)
	statusBytes, err := protos.Marshal(status)
	assert.NoError(t, err)
	record := test_utils.GetDefaultRecord(t)
	recordBytes, err := protos.Marshal(record)
	assert.NoError(t, err)

	// Create attribute value map
	in := map[string]*dynamodb.AttributeValue{
		"NetworkID":                            {S: aws.String("net1")},
		"GatewayID":                            {S: aws.String("gw1")},
		"Status":                               {B: statusBytes},
		"Record":                               {B: recordBytes},
		"Offset":                               {N: aws.String("1")},
		"cfg_" + configManager.GetConfigType(): {B: config1Bytes},
	}

	// Run test
	decoder := dynamo.NewDecoder()
	state, err := decoder.GetStateFromAttributeMap(in)
	assert.NoError(t, err)

	// Verify result
	expectedState := &storage.GatewayState{
		GatewayID: "gw1",
		Status:    status,
		Record:    record,
		Offset:    1,
		Config:    config,
	}
	assert.Equal(t, expectedState, state)
}

func TestInvalidConfigTypeFromAttrMap(t *testing.T) {
	registry.ClearRegistryForTesting()
	config1Manager := test_utils.NewConfig1Manager()
	registry.RegisterConfigManager(config1Manager)
	in := map[string]*dynamodb.AttributeValue{
		"NetworkID":   {S: aws.String("net1")},
		"GatewayID":   {S: aws.String("gw1")},
		"cfg_Config2": {B: []byte{1, 2, 3}},
	}
	decoder := dynamo.NewDecoder()
	_, err := decoder.GetStateFromAttributeMap(in)
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Error unmarshaling config object of type Config2 from bytes"))
}

func TestDecoderStateFromAttrMap_NilFields(t *testing.T) {
	record := test_utils.GetDefaultRecord(t)
	recordBytes, err := protos.Marshal(record)
	assert.NoError(t, err)

	in := map[string]*dynamodb.AttributeValue{
		"NetworkID": {S: aws.String("net1")},
		"GatewayID": {S: aws.String("gw1")},
		"Record":    {B: recordBytes},
		"Offset":    {N: aws.String("1")},
	}

	decoder := dynamo.NewDecoder()
	state, err := decoder.GetStateFromAttributeMap(in)
	assert.NoError(t, err)

	expectedState := &storage.GatewayState{
		Config:    make(map[string]interface{}),
		GatewayID: "gw1",
		Record:    record,
		Offset:    1,
	}
	assert.Equal(t, expectedState, state)
}
