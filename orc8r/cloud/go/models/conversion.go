/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

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

func (m *NetworkName) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID:      network.ID,
		NewName: swag.String(string(*m)),
	}, nil
}

func (m *NetworkName) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkName(network.Name)
}

func (m *NetworkType) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID:      network.ID,
		NewType: swag.String(string(*m)),
	}, nil
}

func (m *NetworkType) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkType(network.Type)
}

func (m *NetworkDescription) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID:             network.ID,
		NewDescription: swag.String(string(*m)),
	}, nil
}

func (m *NetworkDescription) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkDescription(network.Description)
}

func (m *GatewayName) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Key:     gatewayID,
			Type:    orc8r.MagmadGatewayType,
			NewName: swag.String(string(*m)),
		},
	}, nil
}

func (m *GatewayName) FromBackendModels(networkID string, gatewayID string) error {
	entity, err := configurator.LoadSerializedEntity(
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{LoadMetadata: true},
	)
	if err != nil {
		return err
	}
	*m = GatewayName(entity.Name)
	return nil
}

func (m *GatewayDescription) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Key:            gatewayID,
			Type:           orc8r.MagmadGatewayType,
			NewDescription: swag.String(string(*m)),
		},
	}, nil
}

func (m *GatewayDescription) FromBackendModels(networkID string, gatewayID string) error {
	entity, err := configurator.LoadSerializedEntity(
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{LoadMetadata: true},
	)
	if err != nil {
		return err
	}
	*m = GatewayDescription(entity.Description)
	return nil
}
