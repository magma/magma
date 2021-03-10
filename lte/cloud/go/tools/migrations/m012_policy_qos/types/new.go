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

// PolicyRuleConfig policy rule config
// swagger:model policy_rule_config
type PolicyRuleConfig struct {

	// app name
	// Enum: [NO_APP_NAME FACEBOOK FACEBOOK_MESSENGER INSTAGRAM YOUTUBE GOOGLE GMAIL GOOGLE_DOCS NETFLIX APPLE MICROSOFT REDDIT WHATSAPP GOOGLE_PLAY APPSTORE AMAZON WECHAT TIKTOK TWITTER WIKIPEDIA GOOGLE_MAPS YAHOO IMO]
	AppName string `json:"app_name,omitempty"`

	// app service type
	// Enum: [NO_SERVICE_TYPE CHAT AUDIO VIDEO]
	AppServiceType string `json:"app_service_type,omitempty"`

	// flow list
	// Required: true
	FlowList []*FlowDescription `json:"flow_list"`

	// monitoring key
	MonitoringKey string `json:"monitoring_key,omitempty"`

	// priority
	// Required: true
	Priority *uint32 `json:"priority"`

	// rating group
	RatingGroup uint32 `json:"rating_group,omitempty"`

	// redirect
	Redirect *RedirectInformation `json:"redirect,omitempty"`

	// tracking type
	// Enum: [ONLY_OCS ONLY_PCRF OCS_AND_PCRF NO_TRACKING]
	TrackingType string `json:"tracking_type,omitempty"`
}
