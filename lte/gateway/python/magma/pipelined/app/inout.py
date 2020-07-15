"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import threading

from collections import namedtuple

from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL
from threading import Thread

from scapy.arch import get_if_hwaddr
from scapy.data import ETHER_BROADCAST, ETH_P_ARP
from scapy.error import Scapy_Exception
from scapy.layers.l2 import ARP, Ether
from scapy.sendrecv import srp1

from .base import MagmaController
from magma.mobilityd import mobility_store as store
from magma.mobilityd.uplink_gw import UplinkGatewayInfo

from magma.pipelined.app.li_mirror import LIMirrorController
from magma.pipelined.openflow import flows
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import load_direction, Direction, \
    PASSTHROUGH_REG_VAL

from ryu.lib.packet import ether_types

# ingress and egress service names -- used by other controllers
INGRESS = "ingress"
EGRESS = "egress"
PHYSICAL_TO_LOGICAL = "middle"


class InOutController(MagmaController):
    """
    A controller that sets up an openflow pipeline for Magma.

    The EPC controls table 0 which is the first table every packet touches.
    This controller owns the ingress and output portions of the pipeline, the
    first table a packet hits after the EPC controller's table 0 and the last
    table a packet hits before exiting the pipeline.
    """

    APP_NAME = "inout"

    InOutConfig = namedtuple(
        'InOutConfig',
        ['gtp_port', 'uplink_port_name', 'mtr_ip', 'mtr_port', 'li_port_name',
         'enable_nat', 'non_mat_gw_probe_frequency', 'non_nat_arp_egress_port'],
    )
    ARP_PROBE_FREQUENCY = 300
    UPLINK_DPCP_PORT_NAME = 'dhcp0'

    def __init__(self, *args, **kwargs):
        super(InOutController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self._uplink_port = OFPP_LOCAL
        self._li_port = None
        # TODO Alex do we want this to be cofigurable from swagger?
        if self.config.mtr_ip:
            self._mtr_service_enabled = True
        else:
            self._mtr_service_enabled = False
        if self.config.uplink_port_name is not None:
            self._uplink_port = BridgeTools.get_ofport(self.config.uplink_port_name)

        if (self._service_manager.is_app_enabled(LIMirrorController.APP_NAME)
                and self.config.li_port_name):
            self._li_port = BridgeTools.get_ofport(self.config.li_port_name)
            self._li_table = self._service_manager.get_table_num(
                LIMirrorController.APP_NAME)
        self._ingress_tbl_num = self._service_manager.get_table_num(INGRESS)
        self._midle_tbl_num = \
            self._service_manager.get_table_num(PHYSICAL_TO_LOGICAL)
        self._egress_tbl_num = self._service_manager.get_table_num(EGRESS)
        # following fields are only used in Non Nat config
        self._dhcp_gw_info = None
        self._gw_mac_monitor = None

    def _get_config(self, config_dict):
        port_name = None
        mtr_ip = None
        mtr_port = None
        li_port_name = None
        if 'ovs_uplink_port_name' in config_dict:
            port_name = config_dict['ovs_uplink_port_name']

        if 'mtr_ip' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']
        if 'li_local_iface' in config_dict:
            li_port_name = config_dict['li_local_iface']

        enable_nat = config_dict.get('enable_nat', True)
        non_mat_gw_probe_freq = config_dict.get('non_mat_gw_probe_frequency',
                                                self.ARP_PROBE_FREQUENCY)
        non_nat_arp_egress_port = config_dict.get('non_nat_arp_egress_port',
                                                  self.UPLINK_DPCP_PORT_NAME)

        return self.InOutConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            uplink_port_name=port_name,
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            li_port_name=li_port_name,
            enable_nat=enable_nat,
            non_mat_gw_probe_frequency=non_mat_gw_probe_freq,
            non_nat_arp_egress_port=non_nat_arp_egress_port,
        )

    def initialize_on_connect(self, datapath):
        self.delete_all_flows(datapath)
        self._install_default_egress_flows(datapath)
        self._install_default_ingress_flows(datapath)
        self._install_default_middle_flows(datapath)
        self._setup_non_nat_monitoring(datapath)

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self._ingress_tbl_num)
        flows.delete_all_flows_from_table(datapath, self._midle_tbl_num)
        flows.delete_all_flows_from_table(datapath, self._egress_tbl_num)

    def _install_default_middle_flows(self, dp):
        """
        Egress table is the last table that a packet touches in the pipeline.
        Output downlink traffic to gtp port, uplink trafic to LOCAL

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        next_tbl = self._service_manager.get_next_table_num(PHYSICAL_TO_LOGICAL)

        # Allow passthrough pkts(skip enforcement and send to egress table)
        ps_match = MagmaMatch(passthrough=PASSTHROUGH_REG_VAL)
        flows.add_resubmit_next_service_flow(dp, self._midle_tbl_num, ps_match,
                                             actions=[], priority=flows.PASSTHROUGH_PRIORITY,
                                             resubmit_table=self._egress_tbl_num)

        match = MagmaMatch()
        flows.add_resubmit_next_service_flow(dp,
                                             self._midle_tbl_num, match,
                                             actions=[], priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_tbl)

        if self._mtr_service_enabled:
            match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                               ipv4_dst=self.config.mtr_ip)
            flows.add_output_flow(dp,
                                  self._midle_tbl_num, match,
                                  [], priority=flows.UE_FLOW_PRIORITY,
                                  output_port=self.config.mtr_port)

    def _install_default_egress_flows(self, dp, mac_addr: str = None):
        """
        Egress table is the last table that a packet touches in the pipeline.
        Output downlink traffic to gtp port, uplink trafic to LOCAL

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        downlink_match = MagmaMatch(direction=Direction.IN)
        flows.add_output_flow(dp, self._egress_tbl_num, downlink_match, [],
                              output_port=self.config.gtp_port)

        uplink_match = MagmaMatch(direction=Direction.OUT)
        actions = []
        if mac_addr is not None:
            parser = dp.ofproto_parser
            actions.append(parser.OFPActionSetField(eth_dst=mac_addr))
            self.logger.info("Using GW: %s actions: %s", mac_addr, str(actions))

        flows.add_output_flow(dp, self._egress_tbl_num, uplink_match, actions,
                              output_port=self._uplink_port)

    def _install_default_ingress_flows(self, dp):
        """
        Sets up the ingress table, the first step in the packet processing
        pipeline.

        This sets up flow rules to annotate packets with a metadata bit
        indicating the direction. Incoming packets are defined as packets
        originating from the LOCAL port, outgoing packets are defined as
        packets originating from the gtp port.

        All other packets bypass the pipeline.

        Note that the ingress rules do *not* install any flows that cause
        PacketIns (i.e., sends packets to the controller).

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        parser = dp.ofproto_parser
        next_table = self._service_manager.get_next_table_num(INGRESS)

        # set traffic direction bits
        # set a direction bit for outgoing (pn -> inet) traffic.
        match = MagmaMatch(in_port=self.config.gtp_port)
        actions = [load_direction(parser, Direction.OUT)]
        flows.add_resubmit_next_service_flow(dp, self._ingress_tbl_num, match,
                                             actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_table)

        # set a direction bit for incoming (internet -> UE) traffic.
        match = MagmaMatch(in_port=OFPP_LOCAL)
        actions = [load_direction(parser, Direction.IN)]
        flows.add_resubmit_next_service_flow(dp, self._ingress_tbl_num, match,
                                             actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_table)

        # set a direction bit for incoming (internet -> UE) traffic.
        match = MagmaMatch(in_port=self._uplink_port)
        actions = [load_direction(parser, Direction.IN)]
        flows.add_resubmit_next_service_flow(dp, self._ingress_tbl_num, match,
                                             actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_table)

        # Send RADIUS requests directly to li table
        if self._li_port:
            match = MagmaMatch(in_port=self._li_port)
            actions = [load_direction(parser, Direction.IN)]
            flows.add_resubmit_next_service_flow(dp, self._ingress_tbl_num,
                                                 match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                                                 resubmit_table=self._li_table)

        # set a direction bit for incoming (mtr -> UE) traffic.
        if self._mtr_service_enabled:
            match = MagmaMatch(in_port=self.config.mtr_port)
            actions = [load_direction(parser, Direction.IN)]
            flows.add_resubmit_next_service_flow(dp, self._ingress_tbl_num,
                                                 match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                                                 resubmit_table=next_table)

    def _get_gw_mac_address(self, gw_ip: str) -> str:
        try:
            self.logger.debug("sending arp for IP: %s ovs egress: %s",
                              gw_ip, self.config.non_nat_arp_egress_port)
            eth_mac_src = get_if_hwaddr(self.config.non_nat_arp_egress_port)

            pkt = Ether(dst=ETHER_BROADCAST, src=eth_mac_src)
            pkt /= ARP(op="who-has", pdst=gw_ip, hwsrc=eth_mac_src, psrc="0.0.0.0")
            self.logger.debug("pkt: %s", pkt.summary())

            res = srp1(pkt,
                       type=ETH_P_ARP,
                       iface=self.config.non_nat_arp_egress_port,
                       timeout=1,
                       verbose=0,
                       nofilter=1,
                       promisc=0)

            if res is not None:
                self.logger.debug("resp: %s ", res.summary())
                mac = res[ARP].hwsrc
                return mac
            else:
                self.logger.debug("Got Null response")

        except Scapy_Exception as ex:
            self.logger.warning("Error in probing Mac address: err %s", ex)
            return None

    def _monitor_and_update(self, datapath):
        current_feq = self.config.non_mat_gw_probe_frequency
        current_mac = None
        if self._dhcp_gw_info.getMac() is not None:
            current_mac = self._dhcp_gw_info.getMac()
            self._install_default_egress_flows(datapath, current_mac)
            flows.set_barrier(datapath)

        while True:
            ip = self._dhcp_gw_info.getIP()
            if ip is not None:
                self.logger.info("GW found: %s", ip)
                latest_mac_addr = self._get_gw_mac_address(ip)
                if latest_mac_addr is not None:
                    # got back to configured frequency.
                    current_feq = self.config.non_mat_gw_probe_frequency
                    if current_mac != latest_mac_addr:
                        self.logger.info("Current mac %s updated gw mac: %s",
                                         current_mac, latest_mac_addr)
                        current_mac = latest_mac_addr

                        self._install_default_egress_flows(datapath, current_mac)
                        flows.set_barrier(datapath)
                        self._dhcp_gw_info.update_mac(current_mac)
            else:
                self.logger.warning("No default GW found.")
                # increase frequency.
                current_feq = 1

            e = threading.Event()
            self.logger.debug("non_mat_gw_probe_frequency: %s ip: %s mac: %s",
                              current_feq,
                              ip, current_mac)
            e.wait(timeout=current_feq)

    def _setup_non_nat_monitoring(self, datapath):
        """
        Setup egress flow to forward traffic to internet GW.
        Start a thread to figure out MAC address of uplink NAT gw.

        :param datapath: datapath to install flows.
        :return: None
        """
        if self.config.enable_nat is True:
            self.logger.info("Nat is on")
            return
        else:
            self.logger.info("Non nat conf: Frequency:%s, egress port: %s, uplink: %s",
                             self.config.non_mat_gw_probe_frequency,
                             self.config.non_nat_arp_egress_port,
                             self._uplink_port)

        self._dhcp_gw_info = UplinkGatewayInfo(store.GatewayInfoMap())
        self._gw_mac_monitor = Thread(target=self._monitor_and_update,
                                      args=(datapath,))
        self._gw_mac_monitor.setDaemon(True)
        self._gw_mac_monitor.start()
        threading.Event().wait(1)
