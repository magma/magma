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

package libgraphql

import (
	"encoding/json"
)

/*
 * This demo user of the graphql.Client type was taken from FBC XWF meter service (see: fbc/xwf/sync/graphql/ for a full blown client sample).
 */

// AppCustomer is the struct used by the application logic - incoming objects should be converted to the types used by the App logic
type AppCustomer struct {
	// mobile number
	// Min Length: 10
	MobileNumber string `json:"mobile_number,omitempty" gorm:"not null;unique"`
}

// Customer is the GraphQL response for CreateCustomer mutation.
type Customer struct {
	FBID         uint64 `json:"id,string,omitempty"`
	FullName     string `json:"full_name,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
}

// UpsertCustomer holds the sync_upsert_customer mutation logic.
type UpsertCustomer struct {
	MobileNumber string `json:"mobile_number,omitempty"`
	// private state. used to store the GraphQL response.
	resp *Customer
}

// NewCreateCustomer creates a new CreateCustomer mutation from meter's model.
func NewUpsertCustomer(c *AppCustomer) *UpsertCustomer {
	return &UpsertCustomer{
		MobileNumber: c.MobileNumber,
	}
}

// Doc returns the doc string of the sync_create_customer mutation.
func (*UpsertCustomer) Doc() string {
	return `
mutation sync_upsert_customer($data: SyncUpsertCustomerData!) {
  sync_upsert_customer(data: $data) {
    client_mutation_id
    customer {
      id
      mobile_number
    }
  }
}`
}

// Vars returns the variables for the sync_create_customer mutation.
func (c *UpsertCustomer) Vars() (string, error) {
	return Vars{"mobile_number": c.MobileNumber}.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface. Used by the libgraphql.Client.
func (c *UpsertCustomer) UnmarshalJSON(b []byte) error {
	var resp struct {
		Errors Errors `json:"errors,omitempty"`
		Data   struct {
			SyncUpsertCustomer struct {
				Customer         *Customer `json:"customer,omitempty"`
				ClientMutationID string    `json:"client_mutation_id,omitempty"`
			} `json:"sync_upsert_customer,omitempty"`
		} `json:"data,omitempty"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}
	if len(resp.Errors) > 0 {
		return resp.Errors
	}
	c.resp = resp.Data.SyncUpsertCustomer.Customer
	return nil
}

// Response returns the GraphQL customer.
func (c *UpsertCustomer) Response() *Customer { return c.resp }
