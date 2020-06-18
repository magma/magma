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
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/golang/glog"
)

const RecordKeySessionID = "session_id"

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

// GetSessionID returns the session ID stored in the directory record.
// If no session ID is found, returns empty string.
func (m *DirectoryRecord) GetSessionID() (string, error) {
	if m.Identifiers == nil {
		return "", errors.New("directory record's identifiers is nil")
	}

	sid, ok := m.Identifiers[RecordKeySessionID]
	if !ok {
		return "", nil
	}

	sidStr, ok := sid.(string)
	if !ok {
		return "", fmt.Errorf("failed to convert session ID value to string: %v", sid)
	}

	glog.V(2).Infof("Full session ID: %s", sid)
	strippedSid := stripIMSIFromSessionID(sidStr)
	return strippedSid, nil
}

// stripIMSIFromSessionID removes an IMSI prefix from the session ID.
// This exists for backwards compatibility -- in some cases the session ID
// is passed as a dash-separated concatenation of the IMSI and session ID,
// e.g. "IMSI156304337849371-155129".
func stripIMSIFromSessionID(sessionID string) string {
	if strings.HasPrefix(sessionID, "IMSI") && strings.Contains(sessionID, "-") {
		return strings.Split(sessionID, "-")[1]
	}
	return sessionID
}
