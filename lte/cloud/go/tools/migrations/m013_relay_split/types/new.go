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

// Types copied from lte/cloud/go/services/lte/obsidian/models

package types

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"
)

// OldNetworkCellularConfigs Cellular configuration for a network
// swagger:model network_cellular_configs
type NetworkCellularConfigs struct {

	// epc
	// Required: true
	Epc *NetworkEpcConfigs `json:"epc"`

	// feg network id
	FegNetworkID json.RawMessage `json:"feg_network_id,omitempty"`

	// ran
	// Required: true
	Ran json.RawMessage `json:"ran"`
}

// NetworkEpcConfigs EPC (evolved packet core) cellular configuration for a network
// swagger:model network_epc_configs
type NetworkEpcConfigs struct {

	// cloud subscriberdb enabled
	CloudSubscriberdbEnabled bool `json:"cloud_subscriberdb_enabled,omitempty"`

	// default rule id
	DefaultRuleID string `json:"default_rule_id,omitempty"`

	// gx gy relay enabled
	// Required: true
	GxGyRelayEnabled *bool `json:"gx_gy_relay_enabled"`

	// hss relay enabled
	// Required: true
	HssRelayEnabled *bool `json:"hss_relay_enabled"`

	// lte auth amf
	// Required: true
	// Format: byte
	LteAuthAmf strfmt.Base64 `json:"lte_auth_amf"`

	// lte auth op
	// Required: true
	// Max Length: 16
	// Min Length: 15
	// Format: byte
	LteAuthOp strfmt.Base64 `json:"lte_auth_op"`

	// mcc
	// Required: true
	// Pattern: ^(\d{3})$
	Mcc string `json:"mcc"`

	// mnc
	// Required: true
	// Pattern: ^(\d{2,3})$
	Mnc string `json:"mnc"`

	// mobility
	Mobility json.RawMessage `json:"mobility,omitempty"`

	// Configuration for network services. Services will be instantiated in the listed order.
	NetworkServices []string `json:"network_services,omitempty"`

	// sub profiles
	SubProfiles json.RawMessage `json:"sub_profiles,omitempty"`

	// tac
	// Required: true
	// Maximum: 65535
	// Minimum: 1
	Tac uint32 `json:"tac"`
}
