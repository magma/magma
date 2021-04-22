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
import threading
from typing import List

from lte.protos.pipelined_pb2 import (
    FlowResponse,
    SetupFlowsResult,
    UEMacFlowRequest,
)
from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.app.inout import INGRESS
from magma.pipelined.app.ipfix import IPFIXController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.directoryd_client import update_record
from magma.pipelined.imsi import decode_imsi, encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.exceptions import MagmaOFError
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import IMSI_REG, load_passthrough
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import dhcp, ether_types, packet
from ryu.ofproto.inet import IPPROTO_TCP, IPPROTO_UDP


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
        self._loop = kwargs['loop']
        self._datapath = None
        tbls = self._service_manager.allocate_scratch_tables(self.APP_NAME, 2)
        self._passthrough_set_tbl = tbls[0]
        self._dhcp_learn_scratch = tbls[1]
        self._li_port = None
        self._imsi_set_tbl_num = \
            self._service_manager.INTERNAL_IMSI_SET_TABLE_NUM
        self._ipfix_sample_tbl_num = \
            self._service_manager.INTERNAL_IPFIX_SAMPLE_TABLE_NUM
        self._app_set_tbl_num = self._service_manager.INTERNAL_APP_SET_TABLE_NUM
        if 'li_local_iface' in kwargs['config']:
            self._li_port = \
                BridgeTools.get_ofport(kwargs['config']['li_local_iface'])
        self._dpi_port = \
                BridgeTools.get_ofport(kwargs['config']['dpi']['mon_port'])

    def initialize_on_connect(self, datapath):
        self.delete_all_flows(datapath)
        self._datapath = datapath
        self._install_default_flows()

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def handle_restart(self, ue_requests: List[UEMacFlowRequest]
                       ) -> SetupFlowsResult:
        """
        Setup current check quota flows.
        """
        # TODO Potentially we can run a diff logic but I don't think there is
        # benefit(we don't need stats here)
        self.delete_all_flows(self._datapath)
        self._install_default_flows()

        for ue_req in ue_requests:
            self.add_ue_mac_flow(ue_req.sid.id, ue_req.mac_addr)

        self._loop.call_soon_threadsafe(self._setup_arp, ue_requests)

        self.init_finished = True
        return SetupFlowsResult(result=SetupFlowsResult.SUCCESS)

    def _setup_arp(self, ue_requests: List[UEMacFlowRequest]):
        if self.arp_contoller or self.arpd_controller_fut.done():
            if not self.arp_contoller:
                self.arp_contoller = self.arpd_controller_fut.result()
            self.arp_contoller.handle_restart(ue_requests)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._passthrough_set_tbl)
        flows.delete_all_flows_from_table(datapath, self._dhcp_learn_scratch)
        flows.delete_all_flows_from_table(datapath, self._imsi_set_tbl_num,
                                          cookie=self.tbl_num)

    def add_ue_mac_flow(self, sid, mac_addr):
        # TODO report add flow result back to sessiond
        if self._datapath is None:
            return FlowResponse(result=FlowResponse.FAILURE)

        uplink_match = MagmaMatch(eth_src=mac_addr)
        self._add_resubmit_flow(sid, uplink_match,
                                priority=flows.UE_FLOW_PRIORITY,
                                next_table=self._passthrough_set_tbl)

        downlink_match = MagmaMatch(eth_dst=mac_addr)
        self._add_resubmit_flow(sid, downlink_match,
                                priority=flows.UE_FLOW_PRIORITY,
                                next_table=self._passthrough_set_tbl)

        # For handling internal ipfix pkt sampling
        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            self._add_resubmit_flow(sid, uplink_match,
                                    priority=flows.UE_FLOW_PRIORITY,
                                    tbl_num=self._imsi_set_tbl_num,
                                    cookie=self.tbl_num,
                                    next_table=self._ipfix_sample_tbl_num)
            self._add_resubmit_flow(sid, downlink_match,
                                    priority=flows.UE_FLOW_PRIORITY,
                                    tbl_num=self._imsi_set_tbl_num,
                                    cookie=self.tbl_num,
                                    next_table=self._ipfix_sample_tbl_num)

        return FlowResponse(result=FlowResponse.SUCCESS)

    def delete_ue_mac_flow(self, sid, mac_addr):
        # TODO report add flow result back to sessiond
        if self._datapath is None:
            return

        uplink_match = MagmaMatch(eth_src=mac_addr)
        self._delete_resubmit_flow(sid, uplink_match)

        downlink_match = MagmaMatch(eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, downlink_match)

        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            self._delete_resubmit_flow(sid, uplink_match,
                                       tbl_num=self._imsi_set_tbl_num)
            self._delete_resubmit_flow(sid, downlink_match,
                                       tbl_num=self._imsi_set_tbl_num)

    def add_arp_response_flow(self, imsi, yiaddr, chaddr):
        if self.arp_contoller or self.arpd_controller_fut.done():
            if not self.arp_contoller:
                self.arp_contoller = self.arpd_controller_fut.result()
            self.arp_contoller.add_ue_arp_flows(self._datapath,
                                                yiaddr, chaddr)
            self.logger.debug("From DHCP learn: IMSI %s, has ip %s and mac %s",
                              imsi, yiaddr, chaddr)

            # Associate IMSI to IPv4 addr in directory service
            threading.Thread(target=update_record, args=(str(imsi),
                                                         yiaddr)).start()
        else:
            self.logger.error("ARPD controller not ready, ARP learn FAILED")

    def _add_resubmit_flow(self, sid, match, action=None,
                           priority=flows.DEFAULT_PRIORITY,
                           next_table=None, tbl_num=None, cookie=0):
        parser = self._datapath.ofproto_parser

        if action is None:
            actions = []
        else:
            actions = [action]
        if next_table is None:
            next_table = self.next_table
        if tbl_num is None:
            tbl_num = self.tbl_num

        # Add IMSI metadata
        if sid:
            actions.append(parser.NXActionRegLoad2(dst=IMSI_REG,
                                                   value=encode_imsi(sid)))

        flows.add_resubmit_next_service_flow(self._datapath, tbl_num,
                                             match, actions=actions,
                                             priority=priority, cookie=cookie,
                                             resubmit_table=next_table)

    def _delete_resubmit_flow(self, sid, match, action=None, tbl_num=None):
        parser = self._datapath.ofproto_parser

        if action is None:
            actions = []
        else:
            actions = [action]
        if tbl_num is None:
            tbl_num = self.tbl_num

        # Add IMSI metadata
        actions.append(
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid)))

        flows.delete_flow(self._datapath, tbl_num, match, actions=actions)

    def _add_dns_passthrough_flows(self):
        parser = self._datapath.ofproto_parser
        # Set so packet skips enforcement and send to egress
        action = load_passthrough(parser)

        # Install UDP flows for DNS
        ulink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_dst=53)
        self._add_resubmit_flow(None, ulink_match_udp, action,
                                flows.PASSTHROUGH_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

        dlink_match_udp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_UDP,
                                     udp_src=53)
        self._add_resubmit_flow(None, dlink_match_udp, action,
                                flows.PASSTHROUGH_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

        # Install TCP flows for DNS
        ulink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_dst=53)
        self._add_resubmit_flow(None, ulink_match_tcp, action,
                                flows.PASSTHROUGH_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

        dlink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_src=53)
        self._add_resubmit_flow(None, dlink_match_tcp, action,
                                flows.PASSTHROUGH_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

        # Install TCP flows for DNS over tls
        ulink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_dst=853)
        self._add_resubmit_flow(None, ulink_match_tcp, action,
                                flows.PASSTHROUGH_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

        dlink_match_tcp = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                     ip_proto=IPPROTO_TCP,
                                     tcp_src=853)
        self._add_resubmit_flow(None, dlink_match_tcp, action,
                                flows.PASSTHROUGH_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

    def _add_dhcp_passthrough_flows(self):
        ofproto, parser = self._datapath.ofproto, self._datapath.ofproto_parser

        # Set so packet skips enforcement controller
        action = load_passthrough(parser)
        uplink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_src=68,
                                  udp_dst=67)
        self._add_resubmit_flow(None, uplink_match, action,
                                flows.PASSTHROUGH_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

        downlink_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    ip_proto=IPPROTO_UDP,
                                    udp_src=67,
                                    udp_dst=68)
        # Set so triggers packetin and we can learn the ip to do arp response
        self._add_resubmit_flow(None, downlink_match, action,
              flows.PASSTHROUGH_PRIORITY, next_table=self._dhcp_learn_scratch,
              tbl_num=self._passthrough_set_tbl)

        # Install default flow for dhcp learn scratch
        flows.add_output_flow(self._datapath, self._dhcp_learn_scratch,
                              match=MagmaMatch(), actions=[],
                              priority=flows.PASSTHROUGH_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              copy_table=self.next_table,
                              max_len=ofproto.OFPCML_NO_BUFFER)

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

        self._add_dhcp_passthrough_flows()
        self._add_dns_passthrough_flows()

        self._add_resubmit_flow(None, MagmaMatch(),
                                priority=flows.MINIMUM_PRIORITY,
                                tbl_num=self._passthrough_set_tbl)

        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            self._add_resubmit_flow(None, MagmaMatch(in_port=self._dpi_port),
                                    priority=flows.PASSTHROUGH_PRIORITY,
                                    next_table=self._app_set_tbl_num)

        if self._li_port:
            match = MagmaMatch(in_port=self._li_port)
            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                match, actions=[], priority=flows.DEFAULT_PRIORITY,
                resubmit_table=self.next_table)

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
