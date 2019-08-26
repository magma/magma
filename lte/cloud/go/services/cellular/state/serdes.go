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

// EnodebStateSerde is used to serialize/deserialize enodeb operational states
type EnodebStateSerde struct {
}

func (*EnodebStateSerde) GetDomain() string {
	return state.SerdeDomain
}

func (*EnodebStateSerde) GetType() string {
	return lte.EnodebStateType
}

func (*EnodebStateSerde) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func (*EnodebStateSerde) Deserialize(message []byte) (interface{}, error) {
	res := protos.SingleEnodebStatus{}
	err := json.Unmarshal(message, &res)
	return res, err
}
