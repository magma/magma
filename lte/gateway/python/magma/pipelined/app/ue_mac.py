"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

from collections import namedtuple

from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import packet
from ryu.lib.packet import ether_types, dhcp
from ryu.ofproto.inet import IPPROTO_TCP, IPPROTO_UDP

from .base import MagmaController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, IMSI_REG, \
    load_direction


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
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        self._datapath = None
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.arpd_controller_fut = kwargs['app_futures']['arpd']
        self.arp_contoller = None

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
        self._install_default_flows()

    def cleanup_on_disconnect(self, datapath):
        flows.delete_all_flows_from_table(datapath,
                                          self._service_manager.get_table_num(
                                              self.APP_NAME))

    def add_ue_mac_flow(self, sid, mac_addr):
        self._add_dhcp_passthrough_flows(sid, mac_addr)
        self._add_dns_passthrough_flows(sid, mac_addr)

        uplink_match = MagmaMatch(eth_src=mac_addr)
        self._add_resubmit_flow(sid, uplink_match,
                                priority=flows.UE_FLOW_PRIORITY)

        downlink_match = MagmaMatch(eth_dst=mac_addr)
        self._add_resubmit_flow(sid, downlink_match,
                                priority=flows.UE_FLOW_PRIORITY)

    def delete_ue_mac_flow(self, sid, mac_addr):
        self._delete_dhcp_passthrough_flows(sid, mac_addr)
        self._delete_dns_passthrough_flows(sid, mac_addr)

        uplink_match = MagmaMatch(eth_src=mac_addr)
        self._delete_resubmit_flow(sid, uplink_match)

        downlink_match = MagmaMatch(eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, downlink_match)

    def add_arp_response_flow(self, yiaddr, chaddr):
        if self.arp_contoller or self.arpd_controller_fut.done():
            if not self.arp_contoller:
                self.arp_contoller = self.arpd_controller_fut.result()
            self.arp_contoller.add_ue_arp_flows(self._datapath,
                                                yiaddr, chaddr)
        else:
            self.logger.error("ARPD controller not ready, ARP learn FAILED")

    def _add_resubmit_flow(self, sid, match, action=None,
                           priority=flows.DEFAULT_PRIORITY):
        parser = self._datapath.ofproto_parser

        if action is None:
            actions = []
        else:
            actions = [action]

        # Add IMSI metadata
        actions.append(
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid)))

        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             match, actions=actions,
                                             priority=priority,
                                             resubmit_table=self.next_table)

    def _delete_resubmit_flow(self, sid, match, action=None):
        parser = self._datapath.ofproto_parser

        if action is None:
            actions = []
        else:
            actions = [action]

        # Add IMSI metadata
        actions.append(
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid)))

        flows.delete_flow(self._datapath, self.tbl_num, match, actions=actions)

    def _add_dns_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser
        # Set so inout knows to skip tables and send to egress
        action = load_direction(parser, Direction.PASSTHROUGH)

        # Install UDP flows for DNS
        ulink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_dst=53,
                                     eth_src=mac_addr)
        self._add_resubmit_flow(sid, ulink_match_udp, action,
                                flows.PASSTHROUGH_PRIORITY)

        dlink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_src=53,
                                     eth_dst=mac_addr)
        self._add_resubmit_flow(sid, dlink_match_udp, action,
                                flows.PASSTHROUGH_PRIORITY)

        # Install TCP flows for DNS
        ulink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_dst=53,
                                     eth_src=mac_addr)
        self._add_resubmit_flow(sid, ulink_match_tcp, action,
                                flows.PASSTHROUGH_PRIORITY)

        dlink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_src=53,
                                     eth_dst=mac_addr)
        self._add_resubmit_flow(sid, dlink_match_tcp, action,
                                flows.PASSTHROUGH_PRIORITY)

    def _delete_dns_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser
        # Set so inout knows to skip tables and send to egress
        action = load_direction(parser, Direction.PASSTHROUGH)

        # Install UDP flows for DNS
        ulink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_dst=53,
                                     eth_src=mac_addr)
        self._delete_resubmit_flow(sid, ulink_match_udp, action)

        dlink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_src=53,
                                     eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, dlink_match_udp, action)

        # Install TCP flows for DNS
        ulink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_dst=53,
                                     eth_src=mac_addr)
        self._delete_resubmit_flow(sid, ulink_match_tcp, action)

        dlink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_src=53,
                                     eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, dlink_match_tcp, action)

    def _add_dhcp_passthrough_flows(self, sid, mac_addr):
        ofproto = self._datapath.ofproto
        parser = self._datapath.ofproto_parser

        # Set so inout knows to skip tables and send to egress
        action = load_direction(parser, Direction.PASSTHROUGH)
        uplink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_src=68,
                                  udp_dst=67,
                                  eth_src=mac_addr)
        self._add_resubmit_flow(sid, uplink_match, action,
                                flows.PASSTHROUGH_PRIORITY)

        downlink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    ip_proto=IPPROTO_UDP,
                                    udp_src=67,
                                    udp_dst=68,
                                    eth_dst=mac_addr)
        next_table = self._service_manager.get_next_table_num(self.APP_NAME)
        actions = [load_direction(parser, Direction.PASSTHROUGH)]
        # Set so triggers packetin and we can learn the ip to do arp response
        flows.add_output_flow(self._datapath, self.tbl_num, downlink_match,
                              actions, priority=flows.PASSTHROUGH_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              copy_table=next_table,
                              max_len=ofproto.OFPCML_NO_BUFFER)

    def _delete_dhcp_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser

        # Set so inout knows to skip tables and send to egress
        action = load_direction(parser, Direction.PASSTHROUGH)
        uplink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_src=68,
                                  udp_dst=67,
                                  eth_src=mac_addr)
        self._delete_resubmit_flow(sid, uplink_match, action)

        actions = [load_direction(parser, Direction.PASSTHROUGH)]
        downlink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    ip_proto=IPPROTO_UDP,
                                    udp_src=67,
                                    udp_dst=68,
                                    eth_dst=mac_addr)
        # Set so triggers packetin and we can learn the ip to do arp response
        flows.delete_flow(self._datapath, self.tbl_num, downlink_match,
                          actions=actions,
                          priority=flows.PASSTHROUGH_PRIORITY + 1)

    def _add_uplink_arp_allow_flow(self):
        next_table = self._service_manager.get_next_table_num(self.APP_NAME)

        actions = []
        arp_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP)
        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             arp_match, actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_table)

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _learn_arp_entry(self, ev):
        """
        Learn action to process PacketIn DHCP packets, dhcp ack packets will
        be used to learn the ARP entry for the UE to install rules in the arp
        table. The DHCP packets will then be sent thorugh the pipeline.
        """
        msg = ev.msg

        if self.tbl_num != msg.table_id:
            # Intended for other application
            return

        pkt = packet.Packet(msg.data)
        dhcp_header = pkt.get_protocols(dhcp.dhcp)[0]
        # DHCP yiaddr is the client(UE) ip addr
        #      chaddr is the client mac address
        self.add_arp_response_flow(dhcp_header.yiaddr, dhcp_header.chaddr)

    def _install_default_flows(self):
        """
        Install default flows
        """
        # Allows arp packets from uplink(no eth dst set) to go to the arp table
        self._add_uplink_arp_allow_flow()

        # TODO We might want a default drop all rule with min priority, but
        # adding it breakes all unit tests for this controller(needs work)
