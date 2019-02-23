/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

type Decoder interface {
	GetStateFromAttributeMap(map[string]*dynamodb.AttributeValue) (*storage.GatewayState, error)
}

type decoderImpl struct{}

func NewDecoder() Decoder {
	return &decoderImpl{}
}

func (decoder *decoderImpl) GetStateFromAttributeMap(src map[string]*dynamodb.AttributeValue) (*storage.GatewayState, error) {
	state := &gatewayState{}
	err := dynamodbattribute.UnmarshalMap(src, state)
	if err != nil {
		return nil, err
	}
	config, err := parseConfig(src)
	if err != nil {
		return nil, err
	}
	storageState, err := getStateFromDynamoType(state, config)
	if err != nil {
		return nil, err
	}
	return storageState, nil
}

func parseConfig(src map[string]*dynamodb.AttributeValue) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	for key, attrValue := range src {
		if !strings.HasPrefix(key, configPrefix) {
			continue
		}
		configType := key[len(configPrefix):]
		configObj, err := unmarshalConfigValue(configType, attrValue)
		if err != nil {
			return nil, err
		}
		config[configType] = configObj
	}
	return config, nil
}

func unmarshalConfigValue(configType string, attrValue *dynamodb.AttributeValue) (interface{}, error) {
	var configBytes []byte
	err := dynamodbattribute.Unmarshal(attrValue, &configBytes)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling config value to bytes from AttributeValue: %s", err)
	}
	configObj, err := registry.UnmarshalConfig(configType, configBytes)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling config object of type %s from bytes: %s", configType, err)
	}
	return configObj, nil
}

func getStateFromDynamoType(src *gatewayState, config map[string]interface{}) (*storage.GatewayState, error) {
	status, err := getStatus(src)
	if err != nil {
		return nil, err
	}
	record, err := getRecord(src)
	if err != nil {
		return nil, err
	}
	state := &storage.GatewayState{
		GatewayID: src.GatewayID,
		Config:    config,
		Status:    status,
		Record:    record,
		Offset:    src.Offset,
	}
	return state, nil
}

func getStatus(src *gatewayState) (*protos.GatewayStatus, error) {
	if len(src.Status) > 0 {
		status := &protos.GatewayStatus{}
		err := protos.Unmarshal(src.Status, status)
		return status, err
	}
	return nil, nil
}

func getRecord(src *gatewayState) (*magmadprotos.AccessGatewayRecord, error) {
	if len(src.Record) > 0 {
		record := &magmadprotos.AccessGatewayRecord{}
		err := protos.Unmarshal(src.Record, record)
		return record, err
	}
	return nil, nil
}
