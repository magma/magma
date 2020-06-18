/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package state

import (
	"encoding/json"

	"magma/orc8r/cloud/go/serde"
)

const StringMapSerdeType = "string_map"

func NewStateSerde(stateType string, modelPtr serde.ValidateableBinaryConvertible) serde.Serde {
	return serde.NewBinarySerde(SerdeDomain, stateType, modelPtr)
}

// StringToStringMap is a generic map that holds key value pair both of type
// string. This is used on the gateway side in checkin_cli.py to simply test
// the connection between the cloud and the gateway.
type StringToStringMap map[string]string

func (m *StringToStringMap) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *StringToStringMap) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *StringToStringMap) ValidateModel() error {
	return nil
}
