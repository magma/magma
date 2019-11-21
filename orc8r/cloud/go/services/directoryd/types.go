/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */
package directoryd

import (
	"encoding/json"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

type DirectoryRecord struct {
	LocationHistory []string `json:"location_history"`

	Identifiers map[string]interface{} `json:"identifiers"`
}

// ValidateModel is a wrapper to validate this directory record
func (m *DirectoryRecord) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

// Validate validates this directory record
func (m *DirectoryRecord) Validate(formats strfmt.Registry) error {
	if err := validate.Required("location_history", "body", m.LocationHistory); err != nil {
		return err
	}
	return nil
}

// MarshalBinary interface implementation
func (m *DirectoryRecord) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary interface implementation
func (m *DirectoryRecord) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, m)
}
