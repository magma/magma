"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder

MAC_DEST = "5e:cc:cc:b1:49:4b"


def create_uplink_rule(
    id, rating_group, ip_dest, m_key=None,
    priority=10, tracking=PolicyRule.ONLY_OCS,
    action=FlowDescription.PERMIT,
):
    """
    Create a rule with a single uplink IP flow, useful for testing
    Args:
        id (string): rule id
        rating_group (int): charging key
        ip_dest (string): IP destination for rule flow
        m_key (optional string): monitoring key, if the rule is tracked by PCRF
        priority (int): priority of flow, the greater the higher the priority
        tracking (PolicyRule.TrackingType): enum to dictate whether OCS or PCRF
            or both is tracking the credit
        action: permit or deny
    Returns:
        PolicyRule with single uplink IP flow
    """
    return PolicyRule(
        id=id,
        priority=priority,
        flow_list=[
            FlowDescription(
                match=FlowMatch(
                    ipv4_dst=ip_dest, direction=FlowMatch.UPLINK,
                ),
                action=action,
            ),
        ],
        tracking_type=tracking,
        rating_group=rating_group,
        monitoring_key=m_key,
    )


def get_packets_for_flows(sub, flows):
    """
    Get packets sent from a subscriber to match a set of flows
    Args:
        sub (SubscriberContext): subscriber to send packets towards
        flows ([FlowDescription]): list of flows to send matching packets to
    Returns:
        list of scapy packets to send
    """
    packets = []
    for flow in flows:
        packet = IPPacketBuilder()\
            .set_ip_layer(flow.match.ipv4_dst, sub.ip)\
            .set_ether_layer(MAC_DEST, "00:00:00:00:00:00")\
            .build()
        packets.append(packet)
    return packets
