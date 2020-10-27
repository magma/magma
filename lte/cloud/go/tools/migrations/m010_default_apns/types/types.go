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
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
)

const (
	DefaultAPNName = "oai.ipv4"

	defaultAMBRDownlink = 200000000
	defaultAMBRUplink   = 100000000

	defaultQoSClassID                 = 9
	defaultQoSPriorityLevel           = 15
	defaultQoSPreemptionCapability    = true
	defaultQoSPreemptionVulnerability = false
)

var DefaultAPNVal []byte

var defaultAPN = &ApnConfiguration{
	Ambr: &AggregatedMaximumBitrate{
		MaxBandwidthDl: swag.Uint32(defaultAMBRDownlink),
		MaxBandwidthUl: swag.Uint32(defaultAMBRUplink),
	},
	QosProfile: &QosProfile{
		ClassID:                 swag.Int32(defaultQoSClassID),
		PreemptionCapability:    swag.Bool(defaultQoSPreemptionCapability),
		PreemptionVulnerability: swag.Bool(defaultQoSPreemptionVulnerability),
		PriorityLevel:           swag.Uint32(defaultQoSPriorityLevel),
	},
}

func init() {
	DefaultAPNVal = defaultAPN.MustMarshalBinary()
}

type ApnConfiguration struct {
	// ambr
	// Required: true
	Ambr *AggregatedMaximumBitrate `json:"ambr"`

	// qos profile
	// Required: true
	QosProfile *QosProfile `json:"qos_profile"`
}

func (m *ApnConfiguration) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

func (m *ApnConfiguration) UnmarshalBinary(b []byte) error {
	var res ApnConfiguration
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

func (m *ApnConfiguration) MustMarshalBinary() []byte {
	if m == nil {
		return nil
	}
	bytes, err := swag.WriteJSON(m)
	if err != nil {
		glog.Fatalf("Error marshaling APN configuration: %v", err)
	}
	return bytes
}

func (m *ApnConfiguration) String() string {
	return fmt.Sprintf("{Ambr: %v, QosProfile: %v}", m.Ambr, m.QosProfile)
}

type AggregatedMaximumBitrate struct {
	// max bandwidth dl
	// Required: true
	MaxBandwidthDl *uint32 `json:"max_bandwidth_dl"`

	// max bandwidth ul
	// Required: true
	MaxBandwidthUl *uint32 `json:"max_bandwidth_ul"`
}

func (m *AggregatedMaximumBitrate) String() string {
	return fmt.Sprintf("{MaxBandwidthDl: %v, MaxBandwidthUl: %v}", m.MaxBandwidthDl, m.MaxBandwidthUl)
}

type QosProfile struct {
	// class id
	// Maximum: 255
	// Minimum: 0
	ClassID *int32 `json:"class_id,omitempty" magma_alt_name:"QCI"`

	// preemption capability
	PreemptionCapability *bool `json:"preemption_capability,omitempty"`

	// preemption vulnerability
	PreemptionVulnerability *bool `json:"preemption_vulnerability,omitempty"`

	// priority level
	// Maximum: 15
	// Minimum: 0
	PriorityLevel *uint32 `json:"priority_level,omitempty"`
}

func (m *QosProfile) String() string {
	return fmt.Sprintf(
		"{ClassID: %v, PreemptionCapability: %v, PreemptionVulnerability: %v, PriorityLevel: %v}",
		m.ClassID, m.PreemptionCapability, m.PreemptionVulnerability, m.PriorityLevel,
	)
}
