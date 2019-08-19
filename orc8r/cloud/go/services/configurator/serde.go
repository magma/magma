/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package configurator

import (
	"magma/orc8r/cloud/go/serde"
)

// NewNetworkConfigSerde returns a network config domain Serde implementation
// for a pointer to a structure implementing both BinaryMarshaler and
// BinaryUnmarshaler.
// If the modelPtr argument is not a pointer to a struct matching those
// requirements, this function will panic.
func NewNetworkConfigSerde(configType string, modelPtr serde.BinaryConvertible) serde.Serde {
	return serde.NewBinarySerde(NetworkConfigSerdeDomain, configType, modelPtr)
}

// NewNetworkEntityConfigSerde returns a network entity config domain Serde
// implementation/ for a pointer to a structure implementing both
// BinaryMarshaler and BinaryUnmarshaler.
// If the modelPtr argument is not a pointer to a struct matching those
// requirements, this function will panic.
func NewNetworkEntityConfigSerde(configType string, modelPtr serde.BinaryConvertible) serde.Serde {
	return serde.NewBinarySerde(NetworkEntitySerdeDomain, configType, modelPtr)
}
