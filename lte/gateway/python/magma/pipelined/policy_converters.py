"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import ipaddress

from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction

from ryu.lib.packet import ether_types


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


def flow_match_to_magma_match(match):
    '''
    Convert a FlowMatch to a MagmaMatch object

    Args:
        match: FlowMatch
    '''
    _check_pkt_protocol(match)
    match_kwargs = {'eth_type': ether_types.ETH_TYPE_IP}
    attributes = ['ipv4_dst', 'ipv4_src',
                  'ip_proto', 'tcp_src', 'tcp_dst',
                  'udp_src', 'udp_dst', 'app_name']
    for attrib in attributes:
        value = getattr(match, attrib, None)
        if not value:
            continue
        if attrib in {'ipv4_dst', 'ipv4_src'}:
            value = _get_ip_tuple(value)
            if value is None:
                return
        if attrib == 'app_name':
            attrib = 'reg3'

        match_kwargs[attrib] = value
    return MagmaMatch(direction=_get_direction_for_match(match),
                      **match_kwargs)


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
    block_tuple = '{}'.format(ip_block.network_address),\
                  '{}'.format(ip_block.netmask)
    return block_tuple


def _get_direction_for_match(flow_match):
    if flow_match.direction == flow_match.UPLINK:
        return Direction.OUT
    return Direction.IN
