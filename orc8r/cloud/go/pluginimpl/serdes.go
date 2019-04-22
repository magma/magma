/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package pluginimpl

import (
	"fmt"
	"reflect"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/inventory"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
)

// Serdes for user-facing types delegate to swagger's MarshalBinary, which
// writes the bytearray representation of a JSON swagger struct.
// The advantage for us here is that compatibility for changes in the model
// is easier (we can serialize all fields and ignore unknown fields on
// deserialization), and we don't have to define a secondary type for
// serialization to the datastore.

// The context behind this decision is that we used to convert swagger to
// protobuf and save the converted protobuf to the datastore. The idea here
// being that protobuf was a more mature library for this use case. But
// swagger now seems stable enough for storage, and going straight from
// swagger struct to the datastore means that we no longer have to
// define the same struct twice (swagger <-> proto).

// GatewayRecordSerde is an inventory serde for the AccessGatewayRecord type
type GatewayRecordSerde struct{}

func (*GatewayRecordSerde) GetDomain() string {
	return inventory.SerdeDomain
}

func (*GatewayRecordSerde) GetType() string {
	return "access_gateway_record"
}

func (*GatewayRecordSerde) Serialize(in interface{}) ([]byte, error) {
	agr, ok := in.(*models.AccessGatewayRecord)
	if !ok {
		return []byte{}, fmt.Errorf("Could not serialize gateway record. Expected *models.AccessGatewayRecord, got %T", in)
	}
	return agr.MarshalBinary()
}

func (*GatewayRecordSerde) Deserialize(in []byte) (interface{}, error) {
	ret := &models.AccessGatewayRecord{}
	err := ret.UnmarshalBinary(in)
	return ret, err
}

// CheckinRequestSerde is a state serde for the CheckinRequest type
type CheckinRequestSerde struct{}

func (*CheckinRequestSerde) GetDomain() string {
	return state.SerdeDomain
}

func (s *CheckinRequestSerde) GetType() string {
	return "checkin_request"
}

func (s *CheckinRequestSerde) Serialize(in interface{}) ([]byte, error) {
	castedState, ok := in.(*protos.CheckinRequest)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid gateway state type. Expected *CheckinRequest, received %s",
			reflect.TypeOf(in),
		)
	}
	return protos.MarshalIntern(castedState)
}

func (s *CheckinRequestSerde) Deserialize(in []byte) (interface{}, error) {
	response := &protos.CheckinRequest{}
	err := protos.Unmarshal(in, response)
	return response, err
}
