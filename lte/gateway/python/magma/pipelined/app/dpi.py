"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import shlex
import subprocess
import logging

from magma.pipelined.openflow import flows
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.ipfix import IPFIXController
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, DPI_REG
from magma.pipelined.policy_converters import FlowMatchError, \
    flow_match_to_magma_match, flip_flow_match, flow_match_to_actions
from lte.protos.policydb_pb2 import FlowMatch
from lte.protos.pipelined_pb2 import FlowRequest

from ryu.lib.packet import ether_types
from ryu.lib.packet import packet
from ryu.lib.packet import ethernet
from ryu.lib.packet import ipv4
from ryu.lib.packet import tcp
from ryu.lib.packet import udp

# TODO might move to config file
# Current classification will finalize if found in APP_PROTOS, if found in
# PARENT_PROTOS we will also add the SERVICE_IDS id to the final classification
PARENT_PROTOS = {"facebook": 10, "google_gen": 20, "viber": 30, "imo": 40}
APP_PROTOS = {"facebook_messenger": 1, "instagram": 2, "youtube": 3,
              "gmail": 4, "google_docs": 5, "netflix": 6,
              "apple": 7, "microsoft": 8, 'reddit': 9, 'whatsapp': 101,
              "google_play": 102, "appstore": 103, "amazon": 104, "wechat": 105,
              "tiktok": 106, "twitter": 107, "wikipedia": 108, "yahoo": 109}
SERVICE_IDS = {"other": 0, "chat": 1, "audio": 2, "video": 3}
DEFAULT_DPI_ID = 0
# Max register value
UNCLASSIFIED_PROTO_ID = 0xFFFFFFFF

LOG = logging.getLogger('pipelined.app.dpi')


class DPIController(MagmaController):
    """
    DPI controller.

    The DPI controller is responsible for marking a flow with an App ID derived
    from DPI. The APP ID should be stored in register 3
    """

    APP_NAME = "dpi"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(DPIController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.setup_type = kwargs['config']['setup_type']
        self._datapath = None
        self._dpi_enabled = kwargs['config']['dpi']['enabled']
        self._mon_port = kwargs['config']['dpi']['mon_port']
        self._mon_port_number = kwargs['config']['dpi']['mon_port_number']
        self._idle_timeout = kwargs['config']['dpi']['idle_timeout']
        self._bridge_name = kwargs['config']['bridge_name']
        self._app_set_tbl_num = self._service_manager.INTERNAL_APP_SET_TABLE_NUM
        self._imsi_set_tbl_num = \
            self._service_manager.INTERNAL_IMSI_SET_TABLE_NUM
        if self._dpi_enabled:
            self._create_monitor_port()

        tcp_pkt = packet.Packet()
        tcp_pkt.add_protocol(ethernet.ethernet(ethertype=ether_types.ETH_TYPE_IP))
        tcp_pkt.add_protocol(ipv4.ipv4(proto=6))
        tcp_pkt.add_protocol(tcp.tcp())
        tcp_pkt.serialize()
        self.tcp_pkt_data = tcp_pkt.data
        udp_pkt = packet.Packet()
        udp_pkt.add_protocol(ethernet.ethernet(ethertype=ether_types.ETH_TYPE_IP))
        udp_pkt.add_protocol(ipv4.ipv4(proto=17))
        udp_pkt.add_protocol(udp.udp())
        udp_pkt.serialize()
        self.udp_pkt_data = udp_pkt.data

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)
        self._install_default_flows(datapath)
        self._datapath = datapath

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._app_set_tbl_num)

    def add_classify_flow(self, flow_match, flow_state, app: str,
                          service_type: str, src_mac: str, dst_mac: str):
        """
        Parse DPI output and set the register for future packets matching this
        flow. APP is split into tokens as the top level app is not supported,
        but the parent protocol might be.
        Example we care about google traffic, but don't neccessarily want to
        classify every specific google service.
        """
        parser = self._datapath.ofproto_parser

        app_id = get_app_id(app, service_type)

        try:
            ul_match = flow_match_to_magma_match(flow_match)
            ul_match.direction = None
            dl_match = flow_match_to_magma_match(flip_flow_match(flow_match))
            dl_match.direction = None
        except FlowMatchError as e:
            self.logger.error(e)
            return

        actions = [parser.NXActionRegLoad2(dst=DPI_REG, value=app_id)]
        actions_w_mirror = \
            [parser.OFPActionOutput(self._mon_port_number)] + actions
        # No reason to create a flow here
        if flow_state != FlowRequest.FLOW_CREATED:
            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                ul_match, actions_w_mirror, priority=flows.DEFAULT_PRIORITY,
                resubmit_table=self.next_table, idle_timeout=self._idle_timeout)
            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                dl_match, actions_w_mirror, priority=flows.DEFAULT_PRIORITY,
                resubmit_table=self.next_table, idle_timeout=self._idle_timeout)

        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            if (
                flow_state == FlowRequest.FLOW_PARTIAL_CLASSIFICATION
                and app_id == DEFAULT_DPI_ID
            ):
                return
            self._generate_ipfix_sampling_pkt(flow_match, src_mac, dst_mac)
            flows.add_resubmit_next_service_flow(
                self._datapath, self._app_set_tbl_num, ul_match, actions,
                priority=flows.DEFAULT_PRIORITY,
                resubmit_table=self._imsi_set_tbl_num,
                idle_timeout=self._idle_timeout)
            flows.add_resubmit_next_service_flow(
                self._datapath, self._app_set_tbl_num,
                dl_match, actions, priority=flows.DEFAULT_PRIORITY,
                resubmit_table=self._imsi_set_tbl_num,
                idle_timeout=self._idle_timeout)

    def remove_classify_flow(self, flow_match, src_mac: str, dst_mac: str):
        try:
            ul_match = flow_match_to_magma_match(flow_match)
            ul_match.direction = None
            dl_match = flow_match_to_magma_match(flip_flow_match(flow_match))
            dl_match.direction = None
        except FlowMatchError as e:
            self.logger.error(e)
            return False

        flows.delete_flow(self._datapath, self.tbl_num, ul_match)
        flows.delete_flow(self._datapath, self.tbl_num, dl_match)

        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            self._generate_ipfix_sampling_pkt(flow_match, src_mac, dst_mac)
        return True

    def _generate_ipfix_sampling_pkt(self, flow_match, src_mac: str,
                                     dst_mac: str):
        """
        By generating a fake packet trigger the ipfix sampling OVS flow.
        """
        parser = self._datapath.ofproto_parser
        ofproto = self._datapath.ofproto

        if flow_match.ip_proto not in [FlowMatch.IPPROTO_TCP,
                                       FlowMatch.IPPROTO_UDP]:
            self.logger.warning("Ignoring non tcp/udp dpi classification")
            self.logger.warning(flow_match)
            return
        actions = \
            flow_match_to_actions(self._datapath, flow_match)
        actions.extend([
            parser.OFPActionSetField(eth_src=src_mac),
            parser.OFPActionSetField(eth_dst=dst_mac),
            parser.NXActionResubmitTable(table_id=self._app_set_tbl_num)
        ])

        if flow_match.ip_proto == FlowMatch.IPPROTO_TCP:
            bin_packet = self.tcp_pkt_data
        elif flow_match.ip_proto == FlowMatch.IPPROTO_UDP:
            bin_packet = self.udp_pkt_data

        out = parser.OFPPacketOut(datapath=self._datapath,
                                  buffer_id=ofproto.OFP_NO_BUFFER,
                                  in_port=ofproto.OFPP_CONTROLLER,
                                  actions=actions,
                                  data=bin_packet)
        self._datapath.send_msg(out)

    def _install_default_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next table.
        The policies for each subscriber would be added when the IP session is
        created, by reaching out to the controller/PCRF.

        Args:
            datapath: ryu datapath struct
        """
        parser = datapath.ofproto_parser
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)

        actions = [parser.NXActionRegLoad2(dst=DPI_REG,
                                           value=UNCLASSIFIED_PROTO_ID)]
        if self._dpi_enabled:
            actions.append(parser.OFPActionOutput(self._mon_port_number))

        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             inbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)
        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             outbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def _create_monitor_port(self):
        """
        For cwf we set this up when running docker compose as we can't modify
        interfaces from inside the container

        For lte just add the port.
        """
        if self.setup_type == 'CWF':
            self._mon_port_number = BridgeTools.get_ofport(self._mon_port)
            return

        add_cmd = "sudo ovs-vsctl add-port {} mon1 -- set interface {} \
            ofport_request={} type=internal" \
            .format(self._bridge_name, self._mon_port, self._mon_port_number)

        args = shlex.split(add_cmd)
        ret = subprocess.call(args)
        self.logger.debug("Created monitor port ret %d", ret)

        enable_cmd = "sudo ifconfig {} up".format(self._mon_port)
        args = shlex.split(enable_cmd)
        ret = subprocess.call(args)
        self.logger.debug("Enabled monitor port ret %d", ret)


def get_app_id(app: str, service_type: str) -> int:
    """
    Classify the app/service_type to a numeric identifier to export
    """
    if not app or not service_type:
        return DEFAULT_DPI_ID

    app = app.lower()
    service_type = service_type.lower()
    tokens = app.split('.')
    app_match = [app for app in tokens if app in APP_PROTOS]
    if len(app_match) > 1:
        LOG.warning("Found more than 1 app match in %s", app)
        return DEFAULT_DPI_ID

    if (len(app_match) == 1):
        app_id = APP_PROTOS[app_match[0]]
        LOG.debug("Classified %s-%s as %d", app, service_type,
                            app_id)
        return app_id
    parent_match = [app for app in tokens if app in PARENT_PROTOS]

    # This shoudn't happen as we confirmed the match exists
    if len(parent_match) == 0:
        LOG.debug("Didn't find a match for app name %s", app)
        return DEFAULT_DPI_ID
    if len(parent_match) > 1:
        LOG.debug("Found more than 1 parent app match in %s", app)
        return DEFAULT_DPI_ID
    app_id = PARENT_PROTOS[parent_match[0]]

    service_id = SERVICE_IDS['other']
    for serv in SERVICE_IDS:
        if serv in service_type:
            service_id = SERVICE_IDS[serv]
            break
    app_id += service_id
    LOG.debug("Classified %s-%s as %d", app, service_type, app_id)
    return app_id
