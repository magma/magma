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

package analytics

import (
	"encoding/json"
	"fbc/cwf/radius/modules/analytics/graphql"
	"time"
)

// RadiusSession this is the radius session we create in WWW - information returned as mutation result
type RadiusSession struct {
	// f b ID
	FBID uint64 `json:"FBID,omitempty" gorm:"index"`

	// created at
	CreatedAt time.Time `json:"created_at,omitempty"`

	// id
	ID uint64 `json:"id,omitempty" gorm:"primary_key"`

	// updated at
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// n a s IP address
	NASIPAddress string `json:"NASIPAddress,omitempty"`

	// n a s identifier
	NASIdentifier string `json:"NASIdentifier,omitempty"`

	// acct session id
	AcctSessionID string `json:"acct_session_id,omitempty"`

	// called station id
	CalledStationID string `json:"called_station_id,omitempty" gorm:"index:radius_session_index"`

	// calling station id
	CallingStationID string `json:"calling_station_id,omitempty" gorm:"index:radius_session_index"`

	// framed ip address
	FramedIPAddress string `json:"framed_ip_address,omitempty"`

	// normalized mac address
	NormalizedMacAddress string `json:"normalized_mac_address,omitempty"`

	// upload bytes
	UploadBytes int64 `json:"upload_bytes,omitempty"`

	// download bytes
	DownloadBytes int64 `json:"download_bytes,omitempty"`

	//radius_server_id
	RADIUSServerID int64 `json:"radius_server_id,omitempty"`

	// Vendor enum
	Vendor int64 `json:"vendor,omitempty"`
}

// Session is a GraphQL response for create and update radius sessions.
type Session struct {
	FBID             uint64 `json:"id,string,omitempty"`
	AcctSessionID    string `json:"acct_session_id,omitempty"`
	CallingStationID string `json:"calling_station_id,omitempty"`
	CalledStationID  string `json:"called_station_id,omitempty"`
	NASIdentifier    string `json:"nas_identifier,omitempty"`
	NASIPAddress     string `json:"nas_ip_address,omitempty"`
	Vendor           string `json:"vendor_name,omitempty"`
	RADIUSServerID   int64  `json:"radius_server_id,omitempty"`
}

/*
 * CreateSessionOp holds the graphql create_cwfradius_session mutation.
 * Sample usage using GraphiQL tool:
Query/Mutation:
mutation create_cwfradius_session($data: CreateCwfradiusSessionData!) {
mutation create_cwfradius_session($data: CreateCwfradiusSessionData!) {
	create_cwfradius_session(data: $data) {
		client_mutation_id
    radius_session {
      id
    }
  }
}
Query Variables:
{
  "data": {
    "client_mutation_id": 7,
    "acct_session_id": "6",
    "calling_station_id": "5",
    "called_station_id": "4",
    "normalized_mac_address": "aa:bb:cc:dd:ee:ff",
    "nas_identifier": "nas",
    "nas_ip_address": "nas ip",
    "upload_bytes": 100,
    "download_bytes": 100,
    "radius_server_id": 0,
    "vendor_name": "CAMBIUM",
    "framed_ip_address": "1.1.1.1"
  }
}
*/
type CreateSessionOp struct {
	Session *RadiusSession
	// the GraphQL response.
	resp *Session
}

// NewCreateSessionOp creates a new CreateSession mutation from network model.
func NewCreateSessionOp(s *RadiusSession) *CreateSessionOp {
	return &CreateSessionOp{Session: s}
}

// Doc returns the doc string of the update_cwfradius_session mutation.
func (*CreateSessionOp) Doc() string {
	return `
mutation create_cwfradius_session($data: CreateCwfradiusSessionData!) {
  create_cwfradius_session(data: $data) {
    client_mutation_id
    radius_session {
      id
      acct_session_id
    }
  }
}`
}

// Vars returns the variables for the create_cwfradius_session mutation.
func (c *CreateSessionOp) Vars() (string, error) {
	return graphql.Vars{
		"acct_session_id":        c.Session.AcctSessionID,
		"called_station_id":      c.Session.CalledStationID,
		"calling_station_id":     c.Session.CallingStationID,
		"nas_identifier":         c.Session.NASIdentifier,
		"nas_ip_address":         c.Session.NASIPAddress,
		"framed_ip_address":      c.Session.FramedIPAddress,
		"normalized_mac_address": c.Session.NormalizedMacAddress,
		"upload_bytes":           c.Session.UploadBytes,
		"download_bytes":         c.Session.DownloadBytes,
		"radius_server_id":       c.Session.RADIUSServerID,
		"vendor_name":            Vendor(c.Session.Vendor).String(),
	}.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface. Used by the graphql.Client.
func (c *CreateSessionOp) UnmarshalJSON(b []byte) error {
	var resp struct {
		Errors graphql.Errors `json:"errors,omitempty"`
		Data   struct {
			CreateCWFSession struct {
				Session          *Session `json:"radius_session,omitempty"`
				ClientMutationID string   `json:"client_mutation_id,omitempty"`
			} `json:"create_cwfradius_session,omitempty"`
		} `json:"data,omitempty"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}
	if len(resp.Errors) > 0 {
		return resp.Errors
	}
	c.resp = resp.Data.CreateCWFSession.Session
	return nil
}

// Response returns the GraphQL session response.
func (c *CreateSessionOp) Response() *Session { return c.resp }

/*
 * UpdateSessionOp holds the graphql update_cwfradius_session mutation.
 * Sample usage using GraphiQL tool:
Query/Mutation:
mutation update_cwfradius_session($data: UpdateCwfradiusSessionData!) {
	update_cwfradius_session(data: $data) {
		client_mutation_id
    radius_session {
      id
    }
  }
}
Query Variables:
{
  "data": {
    "client_mutation_id": "3932561590",
    "radius_session_id": "2530984173795952",
    "acct_session_id": "6",
    "calling_station_id": "5",
    "called_station_id": "4",
    "normalized_mac_address": "aa:bb:cc:dd:ee:ff",
    "nas_identifier": "nas",
    "nas_ip_address": "nas ip",
    "upload_bytes": 100,
    "download_bytes": 100,
    "radius_server_id": 0,
    "vendor_name": "CAMBIUM",
    "framed_ip_address": "1.1.1.1"
  }
}*/

// UpdateSession holds the GraphQL update_cwfradius_session mutation.
type UpdateSessionOp struct {
	Session *RadiusSession
	// private state. used to store the GraphQL response.
	resp *Session
}

// NewUpdateSessionOp creates a new UpdateSession mutation from the network model.
func NewUpdateSessionOp(s *RadiusSession) *UpdateSessionOp {
	return &UpdateSessionOp{Session: s}
}

// Doc returns the doc string of the update_cwfradius_session mutation.
func (*UpdateSessionOp) Doc() string {
	return `
mutation update_cwfradius_session($data: UpdateCwfradiusSessionData!) {
  update_cwfradius_session(data: $data) {
    client_mutation_id
    radius_session {
      id
      acct_session_id
    }
  }
}`
}

// Vars returns the variables for the update_cwfradius_session mutation.
func (u *UpdateSessionOp) Vars() (string, error) {
	v := graphql.Vars{
		"radius_session_id":      u.Session.FBID,
		"acct_session_id":        u.Session.AcctSessionID,
		"called_station_id":      u.Session.CalledStationID,
		"calling_station_id":     u.Session.CallingStationID,
		"nas_identifier":         u.Session.NASIdentifier,
		"nas_ip_address":         u.Session.NASIPAddress,
		"framed_ip_address":      u.Session.FramedIPAddress,
		"normalized_mac_address": u.Session.NormalizedMacAddress,
		"upload_bytes":           u.Session.UploadBytes,
		"download_bytes":         u.Session.DownloadBytes,
		"radius_server_id":       u.Session.RADIUSServerID,
		"vendor_name":            Vendor(u.Session.Vendor).String(),
	}
	return v.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface. Used by the graphql.Client.
func (u *UpdateSessionOp) UnmarshalJSON(b []byte) error {
	var resp struct {
		Errors graphql.Errors `json:"errors,omitempty"`
		Data   struct {
			UpdateCWFSession struct {
				Session          *Session `json:"radius_session,omitempty"`
				ClientMutationID string   `json:"client_mutation_id,omitempty"`
			} `json:"update_cwfradius_session,omitempty"`
		} `json:"data,omitempty"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}
	if len(resp.Errors) > 0 {
		return resp.Errors
	}
	u.resp = resp.Data.UpdateCWFSession.Session
	return nil
}

// Response returns the GraphQL session response.
func (u *UpdateSessionOp) Response() *Session { return u.resp }
