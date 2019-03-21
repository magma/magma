/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models_test

import (
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/protos"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
)

func TestGatewayStateToModel(t *testing.T) {
	state := &view_factory.GatewayState{
		GatewayID: "gw0",
		Config: map[string]interface{}{
			"Hello": "World!",
		},
		Status: &protos.GatewayStatus{
			Time: 12345,
			Checkin: &protos.CheckinRequest{
				GatewayId:       "gw0",
				MagmaPkgVersion: "v1",
				Status: &protos.ServiceStatus{
					Meta: map[string]string{
						"Test": "Map",
					},
				},
				SystemStatus: &protos.SystemStatus{
					Time:    12345,
					CpuUser: 213,
					MemUsed: 142,
				},
			},
		},
		Record: &magmadprotos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: "gw0"},
			Name: "Gateway 0",
			Key: &protos.ChallengeKey{
				KeyType: protos.ChallengeKey_ECHO,
			},
		},
	}
	expectedModel := &magmad_models.GatewayStateType{
		Config: map[string]interface{}{
			"Hello": "World!",
		},
		GatewayID: "gw0",
		Status: &checkind_models.GatewayStatus{
			CheckinTime: 12345,
			HardwareID:  "gw0",
			Meta: map[string]string{
				"Test": "Map",
			},
			SystemStatus: &checkind_models.SystemStatus{
				Time:    12345,
				CPUUser: 213,
				MemUsed: 142,
			},
			Version: "v1",
		},
		Record: &magmad_models.AccessGatewayRecord{
			HwID: &magmad_models.HwGatewayID{ID: "gw0"},
			Key: &magmad_models.ChallengeKey{
				KeyType: magmad_models.ChallengeKeyKeyTypeECHO,
			},
			Name: "Gateway 0",
		},
	}
	actualModel, err := magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)
}

func TestGatewayStateToModelNilFields(t *testing.T) {
	state := &view_factory.GatewayState{
		GatewayID: "gw0",
		Config:    make(map[string]interface{}),
		Status:    nil,
		Record:    nil,
	}
	expectedModel := &magmad_models.GatewayStateType{
		Config:    make(map[string]interface{}),
		GatewayID: "gw0",
		Status:    nil,
		Record:    nil,
	}
	actualModel, err := magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)

	state.Status = &protos.GatewayStatus{
		Time:    12345,
		Checkin: nil,
	}
	expectedModel.Status = &checkind_models.GatewayStatus{
		CheckinTime:  12345,
		HardwareID:   "",
		Meta:         map[string]string{},
		SystemStatus: nil,
		Version:      "",
	}
	actualModel, err = magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)

	state.Record = &magmadprotos.AccessGatewayRecord{
		HwId: nil,
		Name: "gw0",
		Key:  nil,
	}
	expectedModel.Record = &magmad_models.AccessGatewayRecord{
		HwID: nil,
		Key:  nil,
		Name: "gw0",
	}
	actualModel, err = magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)
}

func TestJSONMapToProtobufStruct(t *testing.T) {
	jsonMap := map[string]interface{}{
		"nil":    nil,
		"number": 1.0,
		"string": "str",
		"struct": map[string]interface{}{
			"a": 2.0,
		},
		"list": []interface{}{1.0, "foo"},
	}
	marshaled, err := json.Marshal(jsonMap)
	assert.NoError(t, err)
	expectedProtobufStruct := &structpb.Struct{}
	err = jsonpb.UnmarshalString(string(marshaled), expectedProtobufStruct)
	assert.NoError(t, err)

	actualProtobufStruct, err := magmad_models.JSONMapToProtobufStruct(jsonMap)

	assert.NoError(t, err)
	assert.Equal(t, expectedProtobufStruct, actualProtobufStruct)
}

func TestProtobufStructToJSONMap(t *testing.T) {
	expectedJsonMap := map[string]interface{}{
		"nil":    nil,
		"number": 1.0,
		"string": "str",
		"struct": map[string]interface{}{
			"a": 2.0,
		},
		"list": []interface{}{1.0, "foo"},
	}
	marshaled, err := json.Marshal(expectedJsonMap)
	assert.NoError(t, err)
	protobufStruct := &structpb.Struct{}
	err = jsonpb.UnmarshalString(string(marshaled), protobufStruct)
	assert.NoError(t, err)

	actualJsonMap, err := magmad_models.ProtobufStructToJSONMap(protobufStruct)

	assert.NoError(t, err)
	assert.Equal(t, expectedJsonMap, actualJsonMap)
}
