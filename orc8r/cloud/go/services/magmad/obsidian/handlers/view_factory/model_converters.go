/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package view_factory

import (
	"fmt"

	"magma/orc8r/cloud/go/protos"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/protobuf/ptypes/struct"
)

// GatewayStateType is the manually defined model type for Gateway State
type GatewayStateType struct {
	Config    map[string]interface{}         `json:"config"`
	GatewayID string                         `json:"gateway_id"`
	Record    *models.AccessGatewayRecord    `json:"record"`
	Status    *checkind_models.GatewayStatus `json:"status"`
}

// GatewayStateToModel converts a storage.GatewayState object to the equivalent
// model.GatewayStateType
func GatewayStateToModel(state *GatewayState) (*GatewayStateType, error) {
	modelState := &GatewayStateType{
		GatewayID: state.GatewayID,
		Config:    state.Config,
	}
	modelStatus, err := gatewayStatusToModel(state.LegacyStatus)
	if err != nil {
		return nil, err
	}
	modelRecord, err := gatewayRecordToModel(state.LegacyRecord)
	if err != nil {
		return nil, err
	}
	modelState.Status = modelStatus
	modelState.Record = modelRecord
	return modelState, nil
}

// GatewayStateMapToModelList converts a map of storage.GatewayState objects
// to an equivalent list of model.GatewayStateType objects
func GatewayStateMapToModelList(states map[string]*GatewayState) ([]*GatewayStateType, error) {
	models := make([]*GatewayStateType, 0, len(states))
	for _, state := range states {
		gatewayState, err := GatewayStateToModel(state)
		if err != nil {
			return nil, err
		}
		models = append(models, gatewayState)
	}
	return models, nil
}

// JSONMapToProtobufStruct converts a map[string]interface{} JSON object to
// the equivalent protobuf Struct
func JSONMapToProtobufStruct(m map[string]interface{}) (*structpb.Struct, error) {
	pbStruct := &structpb.Struct{}
	pbStruct.Fields = make(map[string]*structpb.Value)
	for key, value := range m {
		val, err := jsonValueToProtobufValue(value)
		if err != nil {
			return nil, err
		}
		pbStruct.Fields[key] = val
	}
	return pbStruct, nil
}

// ProtobufStructToJSONMap converts a protobuf Struct to the equivalent
// map[string]interface{} JSON object
func ProtobufStructToJSONMap(s *structpb.Struct) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for key, value := range s.Fields {
		val, err := protobufValueToJSONValue(value)
		if err != nil {
			return nil, err
		}
		m[key] = val
	}
	return m, nil
}

func gatewayStatusToModel(status *protos.GatewayStatus) (*checkind_models.GatewayStatus, error) {
	if status == nil {
		return nil, nil
	}
	modelStatus := &checkind_models.GatewayStatus{}
	err := modelStatus.FromMconfig(status)
	return modelStatus, err
}

func gatewayRecordToModel(record *magmadprotos.AccessGatewayRecord) (*models.AccessGatewayRecord, error) {
	if record == nil {
		return nil, nil
	}
	modelRecord := &models.AccessGatewayRecord{}
	err := modelRecord.FromMconfig(record)
	return modelRecord, err
}

func jsonValueToProtobufValue(jsonValue interface{}) (*structpb.Value, error) {
	switch t := jsonValue.(type) {
	case nil:
		return &structpb.Value{Kind: &structpb.Value_NullValue{NullValue: structpb.NullValue_NULL_VALUE}}, nil
	case float64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: t}}, nil
	case string:
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: t}}, nil
	case map[string]interface{}:
		pbStruct, err := JSONMapToProtobufStruct(t)
		if err != nil {
			return nil, err
		}
		return &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: pbStruct}}, nil
	case []interface{}:
		pbListValue, err := jsonValueSliceToProtobufListValue(t)
		if err != nil {
			return nil, err
		}
		return &structpb.Value{Kind: &structpb.Value_ListValue{ListValue: pbListValue}}, nil
	}
	return nil, fmt.Errorf("Could not convert map value of type %T to protobuf value", jsonValue)
}

func jsonValueSliceToProtobufListValue(jsonValueList []interface{}) (*structpb.ListValue, error) {
	pbList := make([]*structpb.Value, 0)
	for _, listVal := range jsonValueList {
		val, err := jsonValueToProtobufValue(listVal)
		if err != nil {
			return nil, err
		}
		pbList = append(pbList, val)
	}
	return &structpb.ListValue{Values: pbList}, nil
}

func protobufValueToJSONValue(protobufValue *structpb.Value) (interface{}, error) {
	switch t := protobufValue.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil, nil
	case *structpb.Value_NumberValue:
		return t.NumberValue, nil
	case *structpb.Value_StringValue:
		return t.StringValue, nil
	case *structpb.Value_StructValue:
		jsonStruct, err := ProtobufStructToJSONMap(t.StructValue)
		if err != nil {
			return nil, err
		}
		return jsonStruct, nil
	case *structpb.Value_ListValue:
		jsonList, err := protobufListValueToJSONValueSlice(t.ListValue)
		if err != nil {
			return nil, err
		}
		return jsonList, nil
	}
	return nil, fmt.Errorf("Could not convert protobuf value of type %T to JSON value", protobufValue.Kind)
}

func protobufListValueToJSONValueSlice(protobufListValue *structpb.ListValue) ([]interface{}, error) {
	jsonList := make([]interface{}, 0)
	for _, listVal := range protobufListValue.Values {
		val, err := protobufValueToJSONValue(listVal)
		if err != nil {
			return nil, err
		}
		jsonList = append(jsonList, val)
	}
	return jsonList, nil
}
