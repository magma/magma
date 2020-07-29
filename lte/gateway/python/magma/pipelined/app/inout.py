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
import threading

from collections import namedtuple
from ipaddress import IPv4Address

from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

from scapy.arch import get_if_hwaddr
from scapy.data import ETHER_BROADCAST, ETH_P_ARP
from scapy.error import Scapy_Exception
from scapy.layers.l2 import ARP, Ether
from scapy.sendrecv import srp1

from .base import MagmaController
from magma.pipelined.mobilityd_client import get_mobilityd_gw_info, \
    set_mobilityd_gw_info

from magma.pipelined.app.li_mirror import LIMirrorController
from magma.pipelined.openflow import flows
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import load_direction, Direction, \
    PASSTHROUGH_REG_VAL

from ryu.lib import hub
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
         'enable_nat', 'non_mat_gw_probe_frequency', 'non_nat_arp_egress_port',
         'setup_type'],
    )
    ARP_PROBE_FREQUENCY = 300
    NON_NAT_ARP_EGRESS_PORT = 'uplink_br0'

    def __init__(self, *args, **kwargs):
        super(InOutController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self._li_port = None
        # TODO Alex do we want this to be cofigurable from swagger?
        if self.config.mtr_ip:
            self._mtr_service_enabled = True
        else:
            self._mtr_service_enabled = False

        self._uplink_port = OFPP_LOCAL
        if self.config.enable_nat is False and \
                self.config.uplink_port_name is not None:
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
        self._gw_mac_monitor = None
        self._current_upstream_mac = ""
        self._datapath = None

    def _get_config(self, config_dict):
        mtr_ip = None
        mtr_port = None
        li_port_name = None
        port_name = config_dict.get('ovs_uplink_port_name', None)
        setup_type = config_dict.get('setup_type', None)

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
                                                  self.NON_NAT_ARP_EGRESS_PORT)

        return self.InOutConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            uplink_port_name=port_name,
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            li_port_name=li_port_name,
            enable_nat=enable_nat,
            non_mat_gw_probe_frequency=non_mat_gw_probe_freq,
            non_nat_arp_egress_port=non_nat_arp_egress_port,
            setup_type=setup_type,
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

    def _install_default_egress_flows(self, dp, mac_addr: str = ""):
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
        # avoid resetting mac address on switch connect event.
        if mac_addr == "":
            mac_addr = self._current_upstream_mac

        if mac_addr != "":
            parser = dp.ofproto_parser
            actions.append(parser.NXActionRegLoad2(dst='eth_dst', value=mac_addr))

            self.logger.info("Using GW: %s actions: %s", mac_addr, str(actions))
            self._current_upstream_mac = mac_addr

        flows.add_output_flow(dp, self._egress_tbl_num, uplink_match,
                              actions=actions,
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

    def _get_gw_mac_address(self, ip: IPv4Address) -> str:
        try:
            gw_ip = ipaddress.ip_address(ip.address)
            self.logger.debug("sending arp via egress: %s",
                              self.config.non_nat_arp_egress_port)
            eth_mac_src = get_if_hwaddr(self.config.non_nat_arp_egress_port)

            pkt = Ether(dst=ETHER_BROADCAST, src=eth_mac_src)
            pkt /= ARP(op="who-has", pdst=gw_ip, hwsrc=eth_mac_src, psrc="0.0.0.0")
            self.logger.debug("ARP Req pkt %s", pkt.show(dump=True))

            res = srp1(pkt,
                       type=ETH_P_ARP,
                       iface=self.config.non_nat_arp_egress_port,
                       timeout=1,
                       verbose=0,
                       nofilter=1,
                       promisc=0)

            if res is not None:
                self.logger.debug("ARP Res pkt %s", res.show(dump=True))
                mac = res[ARP].hwsrc
                return mac
            else:
                self.logger.debug("Got Null response")
                return ""

        except Scapy_Exception as ex:
            self.logger.warning("Error in probing Mac address: err %s", ex)
            return ""
        except ValueError:
            self.logger.warning("Invalid GW Ip address: [%s]", ip)
            return ""

    def _monitor_and_update(self):
        while True:
            updated_info = get_mobilityd_gw_info()
            if updated_info:
                cached_gw_info = updated_info
            if cached_gw_info and cached_gw_info.ip:
                latest_mac_addr = self._get_gw_mac_address(cached_gw_info.ip)
                if len(latest_mac_addr) == 17 and \
                        self._current_upstream_mac != latest_mac_addr:
                    self._install_default_egress_flows(self._datapath, latest_mac_addr)
                    cached_gw_info.mac = latest_mac_addr
                    set_gw_info = cached_gw_info
                    set_mobilityd_gw_info(set_gw_info)
                elif latest_mac_addr != "":
                    self.logger.warning("invalid mac: %s", latest_mac_addr)
            else:
                self.logger.warning("No default GW found.")

            self.logger.info("non_mat_gw_probe_frequency: %s mac: [%s]",
                             self.config.non_mat_gw_probe_frequency,
                             self._current_upstream_mac)
            hub.sleep(self.config.non_mat_gw_probe_frequency)

    def _setup_non_nat_monitoring(self, datapath):
        """
        Setup egress flow to forward traffic to internet GW.
        Start a thread to figure out MAC address of uplink NAT gw.

        :param datapath: datapath to install flows.
        :return: None
        """
        if self._gw_mac_monitor is not None:
            # No need to multiple probes here.
            return
        if self.config.enable_nat is True:
            self.logger.info("Nat is on")
            return
        elif self.config.setup_type != 'LTE':
            self.logger.info("No GW MAC probe for %s", self.config.setup_type)
            return
        else:
            self.logger.info("Non nat conf: egress port: %s, uplink: %s",
                             self.config.non_nat_arp_egress_port,
                             self._uplink_port)

        self._datapath = datapath
        self._gw_mac_monitor = hub.spawn(self._monitor_and_update)

        threading.Event().wait(1)
