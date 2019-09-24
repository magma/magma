/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package state

import (
	"encoding/json"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/state"
)

// EnodebStateProtosSerde is used to serialize/deserialize enodeb operational states
// todo deprecate after switching to swagger model serde below
type EnodebStateProtosSerde struct {
}

func (*EnodebStateProtosSerde) GetDomain() string {
	return state.SerdeDomain
}

func (*EnodebStateProtosSerde) GetType() string {
	return lte.LegacyEnodebStateType
}

func (*EnodebStateProtosSerde) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func (*EnodebStateProtosSerde) Deserialize(message []byte) (interface{}, error) {
	res := protos.SingleEnodebStatus{}
	err := json.Unmarshal(message, &res)
	return res, err
}
