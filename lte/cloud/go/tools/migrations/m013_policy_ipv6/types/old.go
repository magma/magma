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

import (
	"encoding/json"
)

// OldPolicyRuleConfig policy rule config
// swagger:model policy_rule_config
type OldPolicyRuleConfig struct {
	AppName json.RawMessage `json:"app_name,omitempty"`

	AppServiceType json.RawMessage `json:"app_service_type,omitempty"`

	// flow list
	// Required: true
	FlowList []*OldFlowDescription `json:"flow_list"`

	MonitoringKey json.RawMessage `json:"monitoring_key,omitempty"`

	Priority json.RawMessage `json:"priority"`

	RatingGroup json.RawMessage `json:"rating_group,omitempty"`

	Redirect json.RawMessage `json:"redirect,omitempty"`

	TrackingType json.RawMessage `json:"tracking_type,omitempty"`
}

// OldFlowDescription flow description
// swagger:model flow_description
type OldFlowDescription struct {
	Action json.RawMessage `json:"action"`

	// match
	// Required: true
	Match *OldFlowMatch `json:"match"`
}

// OldFlowMatch flow match
// swagger:model flow_match
type OldFlowMatch struct {

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
