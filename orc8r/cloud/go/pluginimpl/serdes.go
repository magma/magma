/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package pluginimpl

import (
	"encoding/json"

	"magma/orc8r/cloud/go/orc8r"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
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

type GatewayStatusSerde struct{}

func (*GatewayStatusSerde) GetDomain() string {
	return state.SerdeDomain
}

func (s *GatewayStatusSerde) GetType() string {
	return orc8r.GatewayStateType
}

func (s *GatewayStatusSerde) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func (s *GatewayStatusSerde) Deserialize(in []byte) (interface{}, error) {
	response := checkind_models.GatewayStatus{}
	err := json.Unmarshal(in, &response)
	return response, err
}
