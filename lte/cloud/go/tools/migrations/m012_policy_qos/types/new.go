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

// PolicyQosProfile policy qos profile
// swagger:model policy_qos_profile
type PolicyQosProfile struct {

	// arp
	Arp *Arp `json:"arp,omitempty"`

	// class id
	// Required: true
	ClassID QosClassID `json:"class_id"`

	// gbr
	Gbr *Gbr `json:"gbr,omitempty"`

	// id
	// Required: true
	// Min Length: 1
	ID string `json:"id"`

	// max req bw dl
	// Required: true
	MaxReqBwDl *uint32 `json:"max_req_bw_dl"`

	// max req bw ul
	// Required: true
	MaxReqBwUl *uint32 `json:"max_req_bw_ul"`
}

// Arp Allocation and retention priority
// swagger:model arp
type Arp struct {

	// preemption capability
	PreemptionCapability *bool `json:"preemption_capability,omitempty"`

	// preemption vulnerability
	PreemptionVulnerability *bool `json:"preemption_vulnerability,omitempty"`

	// priority level
	// Maximum: 15
	// Minimum: 0
	PriorityLevel *uint32 `json:"priority_level,omitempty"`
}

// QosClassID qos class id
// swagger:model qos_class_id
type QosClassID int32

// Gbr Guaranteed bit rate
// swagger:model gbr
type Gbr struct {

	// downlink
	// Required: true
	Downlink *uint32 `json:"downlink"`

	// uplink
	// Required: true
	Uplink *uint32 `json:"uplink"`
}
