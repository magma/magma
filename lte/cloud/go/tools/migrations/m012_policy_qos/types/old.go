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

// Types copied from lte/cloud/go/services/policydb/obsidian/models

package types

// OldPolicyRuleConfig policy rule config
// swagger:model policy_rule_config
type OldPolicyRuleConfig struct {

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

	// qos
	Qos *FlowQos `json:"qos,omitempty"`

	// rating group
	RatingGroup uint32 `json:"rating_group,omitempty"`

	// redirect
	Redirect *RedirectInformation `json:"redirect,omitempty"`

	// tracking type
	// Enum: [ONLY_OCS ONLY_PCRF OCS_AND_PCRF NO_TRACKING]
	TrackingType string `json:"tracking_type,omitempty"`
}

// FlowDescription flow description
// swagger:model flow_description
type FlowDescription struct {

	// action
	// Required: true
	// Enum: [PERMIT DENY]
	Action *string `json:"action"`

	// match
	// Required: true
	Match *FlowMatch `json:"match"`
}

// FlowMatch flow match
// swagger:model flow_match
type FlowMatch struct {

	// direction
	// Required: true
	// Enum: [UPLINK DOWNLINK]
	Direction *string `json:"direction"`

	// ip proto
	// Required: true
	// Enum: [IPPROTO_IP IPPROTO_TCP IPPROTO_UDP IPPROTO_ICMP]
	IPProto *string `json:"ip_proto"`

	// ipv4 dst
	IPV4Dst string `json:"ipv4_dst,omitempty" magma_alt_name:"Ipv4Dst"`

	// ipv4 src
	IPV4Src string `json:"ipv4_src,omitempty" magma_alt_name:"Ipv4Src"`

	// tcp dst
	TCPDst uint32 `json:"tcp_dst,omitempty" magma_alt_name:"TcpDst"`

	// tcp src
	TCPSrc uint32 `json:"tcp_src,omitempty" magma_alt_name:"TcpSrc"`

	// udp dst
	UDPDst uint32 `json:"udp_dst,omitempty" magma_alt_name:"UdpDst"`

	// udp src
	UDPSrc uint32 `json:"udp_src,omitempty" magma_alt_name:"UdpSrc"`
}

// FlowQos flow qos
// swagger:model flow_qos
type FlowQos struct {

	// max req bw dl
	// Required: true
	MaxReqBwDl *uint32 `json:"max_req_bw_dl"`

	// max req bw ul
	// Required: true
	MaxReqBwUl *uint32 `json:"max_req_bw_ul"`
}

// RedirectInformation redirect information
// swagger:model redirect_information
type RedirectInformation struct {

	// address type
	// Required: true
	// Enum: [IPv4 IPv6 URL SIP_URI]
	AddressType *string `json:"address_type"`

	// server address
	// Required: true
	ServerAddress *string `json:"server_address"`

	// support
	// Required: true
	// Enum: [DISABLED ENABLED]
	Support *string `json:"support"`
}
