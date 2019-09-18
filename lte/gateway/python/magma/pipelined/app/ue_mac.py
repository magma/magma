"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

from collections import namedtuple

from .base import MagmaController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, IMSI_REG, \
    load_direction

from ryu.lib.packet import ether_types
from ryu.ofproto.inet import IPPROTO_TCP, IPPROTO_UDP


class UEMacAddressController(MagmaController):
    """
    UE MAC Address Controller

    This controller controls table 0 which is the first table every packet
    touches. It matches on UE MAC address and sets IMSI metadata
    """

    APP_NAME = "ue_mac"
    UEMacConfig = namedtuple(
        'UEMacConfig',
        ['gre_tunnel_port'],
    )

    def __init__(self, *args, **kwargs):
        super(UEMacAddressController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self._datapath = None

    def _get_config(self, config_dict):
        return self.UEMacConfig(
            # TODO: rename port number to a tunneling protocol agnostic name
            gre_tunnel_port=config_dict['ovs_gtp_port_number'],
        )

    def initialize_on_connect(self, datapath):
        flows.delete_all_flows_from_table(datapath,
                                          self._service_manager.get_table_num(
                                              self.APP_NAME))
        self._datapath = datapath

    def cleanup_on_disconnect(self, datapath):
        flows.delete_all_flows_from_table(datapath,
                                          self._service_manager.get_table_num(
                                              self.APP_NAME))

    def add_ue_mac_flow(self, sid, mac_addr):
        self._add_dhcp_passthrough_flows(sid, mac_addr)
        self._add_dns_passthrough_flows(sid, mac_addr)

        uplink_match = MagmaMatch(in_port=self.config.gre_tunnel_port,
                                  eth_src=mac_addr)
        self._add_resubmit_flow(sid, uplink_match)

        downlink_match = MagmaMatch(in_port=self._datapath.ofproto.OFPP_LOCAL,
                                    eth_dst=mac_addr)
        self._add_resubmit_flow(sid, downlink_match)

    def delete_ue_mac_flow(self, sid, mac_addr):
        self._delete_dhcp_passthrough_flows(sid, mac_addr)
        self._delete_dns_passthrough_flows(sid, mac_addr)

        uplink_match = MagmaMatch(in_port=self.config.gre_tunnel_port,
                                  eth_src=mac_addr)
        self._delete_resubmit_flow(sid, uplink_match)

        downlink_match = MagmaMatch(in_port=self._datapath.ofproto.OFPP_LOCAL,
                                    eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, downlink_match)

    def _add_resubmit_flow(self, sid, match, actions=None,
                           priority=flows.DEFAULT_PRIORITY):
        parser = self._datapath.ofproto_parser
        tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        next_table = self._service_manager.get_next_table_num(self.APP_NAME)

        if actions is None:
            actions = []

        # Add IMSI metadata
        actions.append(
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid)))

        flows.add_resubmit_next_service_flow(self._datapath, tbl_num, match,
                                             actions=actions,
                                             priority=priority,
                                             resubmit_table=next_table)

    def _delete_resubmit_flow(self, sid, match, actions=None):
        parser = self._datapath.ofproto_parser
        tbl_num = self._service_manager.get_table_num(self.APP_NAME)

        if actions is None:
            actions = []

        # Add IMSI metadata
        actions.append(
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid)))

        flows.delete_flow(self._datapath, tbl_num, match, actions=actions)

    def _add_dhcp_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser

        # Set so inout knows to skip tables and send to egress
        actions = [load_direction(parser, Direction.PASSTHROUGH)]

        uplink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_src=68,
                                  udp_dst=67,
                                  eth_src=mac_addr)
        self._add_resubmit_flow(sid, uplink_match, actions,
                                flows.PASSTHROUGH_PRIORITY)

        downlink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    ip_proto=IPPROTO_UDP,
                                    udp_src=67,
                                    udp_dst=68,
                                    eth_dst=mac_addr)
        self._add_resubmit_flow(sid, downlink_match, actions,
                                 flows.PASSTHROUGH_PRIORITY)

    def _delete_dhcp_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser

        # Set so inout knows to skip tables and send to egress
        actions = [load_direction(parser, Direction.PASSTHROUGH)]

        uplink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_src=68,
                                  udp_dst=67,
                                  eth_src=mac_addr)
        self._delete_resubmit_flow(sid, uplink_match, actions)

        downlink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    ip_proto=IPPROTO_UDP,
                                    udp_src=67,
                                    udp_dst=68,
                                    eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, downlink_match, actions)

    def _add_dns_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser
        # Set so inout knows to skip tables and send to egress
        actions = [load_direction(parser, Direction.PASSTHROUGH)]

        # Install UDP flows for DNS
        ulink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_dst=53,
                                     eth_src=mac_addr)
        self._add_resubmit_flow(sid, ulink_match_udp, actions,
                                flows.PASSTHROUGH_PRIORITY)

        dlink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_src=53,
                                     eth_dst=mac_addr)
        self._add_resubmit_flow(sid, dlink_match_udp, actions,
                                flows.PASSTHROUGH_PRIORITY)

        # Install TCP flows for DNS
        ulink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_dst=53,
                                     eth_src=mac_addr)
        self._add_resubmit_flow(sid, ulink_match_tcp, actions,
                                flows.PASSTHROUGH_PRIORITY)

        dlink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_src=53,
                                     eth_dst=mac_addr)
        self._add_resubmit_flow(sid, dlink_match_tcp, actions,
                                flows.PASSTHROUGH_PRIORITY)

    def _delete_dns_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser
        # Set so inout knows to skip tables and send to egress
        actions = [load_direction(parser, Direction.PASSTHROUGH)]

        # Install UDP flows for DNS
        ulink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_dst=53,
                                     eth_src=mac_addr)
        self._delete_resubmit_flow(sid, ulink_match_udp, actions)

        dlink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_src=53,
                                     eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, dlink_match_udp, actions)

        # Install TCP flows for DNS
        ulink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_dst=53,
                                     eth_src=mac_addr)
        self._delete_resubmit_flow(sid, ulink_match_tcp, actions)

        dlink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_src=53,
                                     eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, dlink_match_tcp, actions)
