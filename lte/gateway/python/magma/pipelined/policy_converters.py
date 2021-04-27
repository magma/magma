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
import ipaddress

from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.policydb_pb2 import FlowMatch
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import (
    DPI_REG,
    INGRESS_TUN_ID_REG,
    Direction,
    load_direction,
)
from ryu.lib.packet import ether_types

MATCH_ATTRIBUTES = ['metadata', 'reg0', 'reg1', 'reg2', 'reg3', 'reg4', 'reg5',
                    'reg6', 'reg8', 'reg9', 'reg10',
                    'in_port', 'dl_vlan', 'vlan_tci',
                    'eth_type', 'dl_dst', 'dl_src',
                    'arp_tpa', 'arp_spa', 'arp_op',
                    'ipv4_dst', 'ipv4_src', 'ipv6_src', 'ipv6_dst',
                    'ip_proto', 'tcp_src', 'tcp_dst', 'udp_src', 'udp_dst']


class FlowMatchError(Exception):
    pass


def _check_pkt_protocol(match):
    '''
    Verify that the match flags are set properly

    Args:
        match: FlowMatch
    '''
    if (match.tcp_dst or match.tcp_src) and (match.ip_proto !=
                                             match.IPPROTO_TCP):
        raise FlowMatchError("To use tcp rules set ip_proto to IPPROTO_TCP")
    if (match.udp_dst or match.udp_src) and (match.ip_proto !=
                                             match.IPPROTO_UDP):
        raise FlowMatchError("To use udp rules set ip_proto to IPPROTO_UDP")
    return True


def flow_match_to_magma_match(match, ip_addr=None, uplink_teid: int = None):
    '''
    Convert a FlowMatch to a MagmaMatch object

    Args:
        match: FlowMatch
    '''
    _check_pkt_protocol(match)
    match_kwargs = {'eth_type': ether_types.ETH_TYPE_IP}
    attributes = ['ip_dst', 'ip_src',
                  'ip_proto', 'tcp_src', 'tcp_dst',
                  'udp_src', 'udp_dst', 'app_name']
    for attrib in attributes:
        value = getattr(match, attrib, None)
        if not value:
            continue
        if attrib in {'ip_dst', 'ip_src'}:
            if not value.address:
                continue
            decoded_ip = _get_ip_tuple(value.address.decode('utf-8'))
            if value is None:
                return

            if value.version == IPAddress.IPV4:
                if attrib == 'ip_src':
                    match_kwargs['ipv4_src'] = decoded_ip
                elif attrib == 'ip_dst':
                    match_kwargs['ipv4_dst'] = decoded_ip
            else:
                match_kwargs['eth_type'] = ether_types.ETH_TYPE_IPV6
                if attrib == 'ip_src':
                    match_kwargs['ipv6_src'] = decoded_ip
                elif attrib == 'ip_dst':
                    match_kwargs['ipv6_dst'] = decoded_ip
            continue
        elif attrib == 'app_name':
            attrib = DPI_REG

        match_kwargs[attrib] = value

    # Specific UE IP match
    if ip_addr:
        if ip_addr.version == IPAddress.IPV4:
            ip_src_reg = 'ipv4_src'
            ip_dst_reg = 'ipv4_dst'
        else:
            match_kwargs['eth_type'] = ether_types.ETH_TYPE_IPV6
            ip_src_reg = 'ipv6_src'
            ip_dst_reg = 'ipv6_dst'

        if ip_addr.address.decode('utf-8'):
            if get_direction_for_match(match) == Direction.OUT:
                match_kwargs[ip_src_reg] = ip_addr.address.decode('utf-8')
            else:
                match_kwargs[ip_dst_reg] = ip_addr.address.decode('utf-8')

    if uplink_teid is not None and uplink_teid != 0:
        match_kwargs[INGRESS_TUN_ID_REG] = uplink_teid

    return MagmaMatch(direction=get_direction_for_match(match),
                      **match_kwargs)


def flow_match_to_actions(datapath, match):
    '''
    Convert a FlowMatch to list of actions to get the same packet

    Args:
        match: FlowMatch
    '''
    parser = datapath.ofproto_parser
    _check_pkt_protocol(match)
    # Eth type and ip proto are read only, can't set them here (set on pkt init)
    actions = [
        parser.OFPActionSetField(ipv4_src=getattr(match, 'ipv4_src', '1.1.1.1')),
        parser.OFPActionSetField(ipv4_dst=getattr(match, 'ipv4_dst', '1.2.3.4')),
        load_direction(parser, get_direction_for_match(match)),
        parser.NXActionRegLoad2(dst=DPI_REG, value=getattr(match, 'app_id', 0)),
    ]
    if match.ip_proto == FlowMatch.IPPROTO_TCP:
        actions.extend([
            parser.OFPActionSetField(tcp_src=getattr(match, 'tcp_src', 0)),
            parser.OFPActionSetField(tcp_dst=getattr(match, 'tcp_dst', 0))
        ])
    elif match.ip_proto == FlowMatch.IPPROTO_UDP:
        actions.extend([
            parser.OFPActionSetField(udp_src=getattr(match, 'udp_src', 0)),
            parser.OFPActionSetField(udp_dst=getattr(match, 'udp_dst', 0))
        ])
    return actions


def flip_flow_match(match):
    '''
    Flips FlowMatch(ip/ports/direction)

    Args:
        match: FlowMatch
    '''
    if getattr(match, 'direction', None) == match.DOWNLINK:
        direction = match.UPLINK
    else:
        direction = match.DOWNLINK

    return FlowMatch(
        ip_src=getattr(match, 'ip_dst', None),
        ip_dst=getattr(match, 'ip_src', None),
        tcp_src=getattr(match, 'tcp_dst', None),
        tcp_dst=getattr(match, 'tcp_src', None),
        udp_src=getattr(match, 'udp_dst', None),
        udp_dst=getattr(match, 'udp_src', None),
        ip_proto=getattr(match, 'ip_proto', None),
        direction=direction,
        app_name=getattr(match, 'app_name', None)
    )


def get_flow_ip_dst(match):
    ip_dst = getattr(match, 'ip_dst', None)
    if ip_dst is None:
        return
    decoded_ip = ip_dst.address.decode('utf-8')

    if ip_dst.version == IPAddress.IPV4:
        return decoded_ip
    else:
        return None


def ipv4_address_to_str(ipaddr: IPAddress):

    decoded_ip = ipaddr.address.decode('utf-8')

    if ipaddr.version == IPAddress.IPV4:
        return decoded_ip
    else:
        return None


def get_ue_ip_match_args(ip_addr: IPAddress, direction: Direction):
    ip_match = {}

    if ip_addr:
        if ip_addr.version == ip_addr.IPV4:
            ip_src_reg = 'ipv4_src'
            ip_dst_reg = 'ipv4_dst'
        else:
            ip_src_reg = 'ipv6_src'
            ip_dst_reg = 'ipv6_dst'

        if not ip_addr.address.decode('utf-8'):
            return ip_match

        if direction == Direction.OUT:
            ip_match = {ip_src_reg: ip_addr.address.decode('utf-8')}
        else:
            ip_match = {ip_dst_reg: ip_addr.address.decode('utf-8')}
    return ip_match


def get_eth_type(ip_addr: IPAddress):
    if not ip_addr:
        return ether_types.ETH_TYPE_IP
    if ip_addr.version == IPAddress.IPV4:
        return ether_types.ETH_TYPE_IP
    else:
        return ether_types.ETH_TYPE_IPV6


def _get_ip_tuple(ip_str):
    '''
    Convert an ip string to a formatted block tuple

    Args:
        ip_str (string): ip string to parse
    '''
    try:
        ip_block = ipaddress.ip_network(ip_str)
    except ValueError as err:
        raise FlowMatchError("Invalid Ip block: %s" % err)
    block_tuple = '{}'.format(ip_block.network_address), \
                  '{}'.format(ip_block.netmask)
    return block_tuple


def get_direction_for_match(flow_match):
    if flow_match.direction == flow_match.UPLINK:
        return Direction.OUT
    return Direction.IN


def convert_ipv4_str_to_ip_proto(ipv4_str):
    return IPAddress(version=IPAddress.IPV4,
                     address=ipv4_str.encode('utf-8'))


def convert_ipv6_bytes_to_ip_proto(ipv6_bytes):
    return IPAddress(version=IPAddress.IPV6,
                     address=ipv6_bytes)


def convert_ip_str_to_ip_proto(ip_str: str):
    if ip_str.count(":") >= 2:
        ip_addr = \
            convert_ipv6_bytes_to_ip_proto(ip_str.encode('utf-8'))
    else:
        ip_addr = convert_ipv4_str_to_ip_proto(ip_str)
    return ip_addr


def ovs_flow_match_to_magma_match(flow):
    attribute_dict = {}
    for a in MATCH_ATTRIBUTES:
        val = flow.match.get(a, None)
        if val:
            attribute_dict[a] = val
    return MagmaMatch(**attribute_dict)
