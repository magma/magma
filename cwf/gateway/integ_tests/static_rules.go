/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package integration

import (
	"magma/feg/cloud/go/protos"
	lteProtos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"github.com/go-openapi/swag"
)

func getUsageInformation(monitorKey string, quota uint64) *protos.UsageMonitoringInformation {
	return &protos.UsageMonitoringInformation{
		MonitoringLevel: protos.MonitoringLevel_RuleLevel,
		MonitoringKey:   []byte(monitorKey),
		Octets:          &protos.Octets{TotalOctets: quota},
	}
}

const MATCH_ALL = "0.0.0.0/0"

func getStaticPassAll(
	ruleID string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32, qos *lteProtos.FlowQos,
) *lteProtos.PolicyRule {
	return getStaticPassTraffic(ruleID, MATCH_ALL, MATCH_ALL, monitoringKey, ratingGroup, trackingType, priority, qos)
}

func getStaticDenyAll(ruleID string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32) *lteProtos.PolicyRule {
	rule := &models.PolicyRuleConfig{
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("DENY"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPDst: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: "0.0.0.0/0",
					},
					IPSrc: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: "0.0.0.0/0",
					},
				},
			},
			{
				Action: swag.String("DENY"),
				Match: &models.FlowMatch{
					Direction: swag.String("DOWNLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPDst: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: "0.0.0.0/0",
					},
					IPSrc: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: "0.0.0.0/0",
					},
				},
			},
		},
		MonitoringKey: monitoringKey,
		Priority:      swag.Uint32(priority),
		TrackingType:  trackingType,
		RatingGroup:   ratingGroup,
	}

	return rule.ToProto(ruleID, nil)
}

func getStaticPassTraffic(
	ruleID string, srcIP string, dstIP string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32, qos *lteProtos.FlowQos,
) *lteProtos.PolicyRule {
	rule := &models.PolicyRuleConfig{
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPDst: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: dstIP,
					},
					IPSrc: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: srcIP,
					},
					IPV4Dst: dstIP,
					IPV4Src: srcIP,
				},
			},
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("DOWNLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPDst: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: srcIP,
					},
					IPSrc: &models.IPAddress{
						Version: models.IPAddressVersionIPV4,
						Address: dstIP,
					},
				},
			},
		},
		MonitoringKey: monitoringKey,
		Priority:      swag.Uint32(priority),
		TrackingType:  trackingType,
		RatingGroup:   ratingGroup,
	}
	return rule.ToProto(ruleID, qos)
}
