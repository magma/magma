"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This file models common logic of the mandatory services ingress,
middle and egress.
"""

from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import TUN_PORT_REG, Direction


def get_vlan_egress_flow_msgs(
    dp, table_no, eth_type, ip, out_port=None,
    priority=0, direction=Direction.IN, dst_mac=None,
):
    """
    Install egress flows
    Args:
        dp datapath
        table_no table to install flow
        out_port specify egress port, if None reg value is used
        priority flow priority
        direction packet direction.
    """
    msgs = []
    if out_port:
        output_reg = None
    else:
        output_reg = TUN_PORT_REG

    # Pass non vlan packet as it is.
    # TODO: add support to match IPv6 address
    if ip:
        match = MagmaMatch(
            direction=direction,
            eth_type=eth_type,
            vlan_vid=(0x0000, 0x1000),
            ipv4_dst=ip,
        )
    else:
        match = MagmaMatch(
            direction=direction,
            eth_type=eth_type,
            vlan_vid=(0x0000, 0x1000),
        )
    actions = []
    if dst_mac:
        actions.append(dp.ofproto_parser.NXActionRegLoad2(dst='eth_dst', value=dst_mac))

    msgs.append(
        flows.get_add_output_flow_msg(
            dp, table_no, match, actions,
            priority=priority, output_reg=output_reg, output_port=out_port,
        ),
    )

    # remove vlan header for out_port.
    if ip:
        match = MagmaMatch(
            direction=direction,
            eth_type=eth_type,
            vlan_vid=(0x1000, 0x1000),
            ipv4_dst=ip,
        )
    else:
        match = MagmaMatch(
            direction=direction,
            eth_type=eth_type,
            vlan_vid=(0x1000, 0x1000),
        )
    actions = [dp.ofproto_parser.OFPActionPopVlan()]
    if dst_mac:
        actions.append(dp.ofproto_parser.NXActionRegLoad2(dst='eth_dst', value=dst_mac))

    msgs.append(
        flows.get_add_output_flow_msg(
            dp, table_no, match, actions,
            priority=priority, output_reg=output_reg, output_port=out_port,
        ),
    )
    return msgs
