/*
 Copyright 2020 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/golang/glog"
)

const RecordKeySessionID = "session_id"
const RecordKeySpgCTeid = "sgw_c_teid"

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
		return "", fmt.Errorf("directory record's identifiers is nil")
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
	if strippedSid == "" {
		return "", fmt.Errorf("empty session id for record %+v", m)
	}
	return strippedSid, nil
}

// GetCurrentLocation returns LocalHistory if exists, otherwise returns error
func (m *DirectoryRecord) GetCurrentLocation() (string, error) {
	if m == nil {
		return "", fmt.Errorf("directory record is nil. Cant find LocalHistory")
	}
	if len(m.LocationHistory) == 0 {
		return "", fmt.Errorf("empty LocationHistory")
	}
	return m.LocationHistory[0], nil
}

// GetSgwCTeids returns the S8Teids stored in the directory record.
func (m *DirectoryRecord) GetSgwCTeids() ([]string, error) {
	if m.Identifiers == nil {
		return nil, fmt.Errorf("directory record's identifiers is nil")
	}

	storedTeids, ok := m.Identifiers[RecordKeySpgCTeid]
	if !ok {
		return []string{}, nil
	}

	teidsStr, ok := storedTeids.(string)
	if !ok {
		return nil, fmt.Errorf("failed to convert session ID value to string: %v", storedTeids)
	}

	glog.V(2).Infof("S8 Teids: %s", teidsStr)
	return parseSgwCTeids(teidsStr), nil
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

// parseSgwCTeids converts comma separated string into a slice
func parseSgwCTeids(teidsStr string) []string {
	teids := []string{}
	splitTeids := strings.Split(teidsStr, ",")
	for _, teid := range splitTeids {
		teids = append(teids, strings.TrimSpace(teid))
	}
	return teids
}
