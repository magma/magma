"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import threading

from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import packet
from ryu.lib.packet import ether_types, dhcp
from ryu.ofproto.inet import IPPROTO_TCP, IPPROTO_UDP

from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.inout import INGRESS
from magma.pipelined.directoryd_client import update_record
from magma.pipelined.imsi import encode_imsi, decode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.exceptions import MagmaOFError
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import IMSI_REG, load_passthrough


class UEMacAddressController(MagmaController):
    """
    UE MAC Address Controller

    This controller controls table 0 which is the first table every packet
    touches. It matches on UE MAC address and sets IMSI metadata
    """

    APP_NAME = "ue_mac"
    APP_TYPE = ControllerType.SPECIAL

    def __init__(self, *args, **kwargs):
        super(UEMacAddressController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_table_num(INGRESS)
        self.arpd_controller_fut = kwargs['app_futures']['arpd']
        self.arp_contoller = None
        self._datapath = None
        self._dhcp_learn_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]

    def initialize_on_connect(self, datapath):
        self.delete_all_flows(datapath)
        self._datapath = datapath
        self._install_default_flows()

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._dhcp_learn_scratch)

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

    def add_arp_response_flow(self, imsi, yiaddr, chaddr):
        if self.arp_contoller or self.arpd_controller_fut.done():
            if not self.arp_contoller:
                self.arp_contoller = self.arpd_controller_fut.result()
            self.arp_contoller.add_ue_arp_flows(self._datapath,
                                                yiaddr, chaddr)
            self.logger.debug("Learned arp for imsi %s, ip %s", imsi, yiaddr)

            # Associate IMSI to IPv4 addr in directory service
            threading.Thread(target=update_record, args=(str(imsi),
                                                         yiaddr)).start()
        else:
            self.logger.error("ARPD controller not ready, ARP learn FAILED")

    def _add_resubmit_flow(self, sid, match, action=None,
                           priority=flows.DEFAULT_PRIORITY,
                           next_table=None):
        parser = self._datapath.ofproto_parser

        if action is None:
            actions = []
        else:
            actions = [action]
        if next_table is None:
            next_table = self.next_table

        # Add IMSI metadata
        actions.append(
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid)))

        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             match, actions=actions,
                                             priority=priority,
                                             resubmit_table=next_table)

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
        # Set so packet skips enforcement and send to egress
        action = load_passthrough(parser)

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
        # Set so packet skips enforcement controller
        action = load_passthrough(parser)

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
        ofproto, parser = self._datapath.ofproto, self._datapath.ofproto_parser

        # Set so packet skips enforcement controller
        action = load_passthrough(parser)
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
        # Set so triggers packetin and we can learn the ip to do arp response
        self._add_resubmit_flow(sid, downlink_match, action,
              flows.PASSTHROUGH_PRIORITY, next_table=self._dhcp_learn_scratch)

        # Install default flow for dhcp learn scratch
        imsi_match = MagmaMatch(imsi=encode_imsi(sid))
        flows.add_output_flow(self._datapath, self._dhcp_learn_scratch,
                              match=imsi_match, actions=[],
                              priority=flows.PASSTHROUGH_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              copy_table=self.next_table,
                              max_len=ofproto.OFPCML_NO_BUFFER)

    def _delete_dhcp_passthrough_flows(self, sid, mac_addr):
        parser = self._datapath.ofproto_parser

        # Set so packet skips enforcement controller
        action = load_passthrough(parser)
        uplink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_src=68,
                                  udp_dst=67,
                                  eth_src=mac_addr)
        self._delete_resubmit_flow(sid, uplink_match, action)

        downlink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    ip_proto=IPPROTO_UDP,
                                    udp_src=67,
                                    udp_dst=68,
                                    eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, downlink_match, action)
        imsi_match = MagmaMatch(imsi=encode_imsi(sid))
        flows.delete_flow(self._datapath, self._dhcp_learn_scratch, imsi_match)

    def _add_uplink_arp_allow_flow(self):
        arp_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP)
        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             arp_match, actions=[],
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=self.next_table)

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _learn_arp_entry(self, ev):
        """
        Learn action to process PacketIn DHCP packets, dhcp ack packets will
        be used to learn the ARP entry for the UE to install rules in the arp
        table. The DHCP packets will then be sent thorugh the pipeline.
        """
        msg = ev.msg

        if self._dhcp_learn_scratch != msg.table_id:
            # Intended for other application
            return

        try:
            encoded_imsi = _get_encoded_imsi_from_packetin(msg)
            # Decode the imsi to properly save in directoryd
            imsi = decode_imsi(encoded_imsi)
        except MagmaOFError as e:
            # No packet direction, but intended for this table
            self.logger.error("Error obtaining IMSI from pkt-in: %s", e)
            return

        pkt = packet.Packet(msg.data)
        dhcp_header = pkt.get_protocols(dhcp.dhcp)[0]
        # DHCP yiaddr is the client(UE) ip addr
        #      chaddr is the client mac address
        self.add_arp_response_flow(imsi, dhcp_header.yiaddr, dhcp_header.chaddr)

    def _install_default_flows(self):
        """
        Install default flows
        """
        # Allows arp packets from uplink(no eth dst set) to go to the arp table
        self._add_uplink_arp_allow_flow()

        # TODO We might want a default drop all rule with min priority, but
        # adding it breakes all unit tests for this controller(needs work)


def _get_encoded_imsi_from_packetin(msg):
    """
    Retrieve encoded imsi from the Packet-In message, or raise an exception if
    it doesn't exist.
    """
    imsi = msg.match.get(IMSI_REG)
    if imsi is None:
        raise MagmaOFError('IMSI not found in OFPMatch')
    return imsi
