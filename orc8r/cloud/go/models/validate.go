/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"github.com/go-openapi/strfmt"
)

func (m *NetworkName) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkType) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkDescription) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayName) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayDescription) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
