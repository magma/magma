/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package state

import "magma/orc8r/cloud/go/serde"

func NewStateSerde(stateType string, modelPtr serde.BinaryConvertible) serde.Serde {
	return serde.NewBinarySerde(SerdeDomain, stateType, modelPtr)
}
