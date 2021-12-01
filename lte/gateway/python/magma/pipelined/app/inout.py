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

from lte.protos.mobilityd_pb2 import IPAddress
from magma.pipelined.app.base import MagmaController
from magma.pipelined.app.li_mirror import LIMirrorController
from magma.pipelined.app.restart_mixin import DefaultMsgsMap, RestartMixin
from magma.pipelined.bridge_util import BridgeTools, DatapathLookupError
from magma.pipelined.mobilityd_client import (
    get_mobilityd_gw_info,
    set_mobilityd_gw_info,
)
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MessageHub, MsgChannel
from magma.pipelined.openflow.registers import (
    PASSTHROUGH_REG_VAL,
    PROXY_TAG_TO_PROXY,
    REG_ZERO_VAL,
    TUN_PORT_REG,
    Direction,
    load_direction,
)
from magma.pipelined.utils import get_virtual_iface_mac
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib import hub
from ryu.lib.packet import ether_types
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL
from scapy.arch import get_if_addr, get_if_hwaddr
from scapy.data import ETH_P_ALL, ETHER_BROADCAST
from scapy.error import Scapy_Exception
from scapy.layers.inet6 import getmacbyip6
from scapy.layers.l2 import ARP, Dot1Q, Ether
from scapy.sendrecv import srp1

# ingress and egress service names -- used by other controllers

INGRESS = "ingress"
EGRESS = "egress"
PHYSICAL_TO_LOGICAL = "middle"
PROXY_PORT_MAC = 'e6:8f:a2:80:80:80'


class InOutController(RestartMixin, MagmaController):
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
        [
            'gtp_port', 'uplink_port', 'mtr_ip', 'mtr_port', 'li_port_name',
            'enable_nat', 'non_nat_gw_probe_frequency', 'non_nat_arp_egress_port',
            'setup_type', 'uplink_gw_mac', 'he_proxy_port', 'he_proxy_eth_mac',
            'mtr_mac', 'virtual_mac',
        ],
    )
    ARP_PROBE_FREQUENCY = 300
    NON_NAT_ARP_EGRESS_PORT = 'dhcp0'
    UPLINK_OVS_BRIDGE_NAME = 'uplink_br0'

    def __init__(self, *args, **kwargs):
        super(InOutController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.logger.info("inout config: %s", self.config)

        self._li_port = None
        # TODO Alex do we want this to be cofigurable from swagger?
        if self.config.mtr_ip:
            self._mtr_service_enabled = True
        else:
            self._mtr_service_enabled = False

        if (
            self._service_manager.is_app_enabled(LIMirrorController.APP_NAME)
            and self.config.li_port_name
        ):
            self._li_port = BridgeTools.get_ofport(self.config.li_port_name)
            self._li_table = self._service_manager.get_table_num(
                LIMirrorController.APP_NAME,
            )
        self._ingress_tbl_num = self._service_manager.get_table_num(INGRESS)
        self._midle_tbl_num = \
            self._service_manager.get_table_num(PHYSICAL_TO_LOGICAL)
        self._egress_tbl_num = self._service_manager.get_table_num(EGRESS)
        # following fields are only used in Non Nat config
        self._tbls = [
            self._ingress_tbl_num, self._midle_tbl_num,
            self._egress_tbl_num,
        ]
        self._gw_mac_monitor = None
        self._current_upstream_mac_map = {}  # maps vlan to upstream gw mac
        self._clean_restart = kwargs['config']['clean_restart']
        self._msg_hub = MessageHub(self.logger)
        self._datapath = None
        self._gw_mac_monitor_on = False

    def _get_config(self, config_dict):
        mtr_ip = None
        mtr_port = None
        li_port_name = None
        port_no = config_dict.get('uplink_port', None)
        setup_type = config_dict.get('setup_type', None)

        he_proxy_port = 0
        he_proxy_eth_mac = ''
        try:
            if 'proxy_port_name' in config_dict:
                he_proxy_port = BridgeTools.get_ofport(config_dict.get('proxy_port_name'))
                he_proxy_eth_mac = config_dict.get('he_proxy_eth_mac', PROXY_PORT_MAC)
        except DatapathLookupError:
            # ignore it
            self.logger.debug("could not parse proxy port config")

        if 'mtr_ip' in config_dict and 'mtr_interface' in config_dict and 'ovs_mtr_port_number' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']
            mtr_mac = get_virtual_iface_mac(config_dict['mtr_interface'])
        else:
            mtr_ip = None
            mtr_mac = None
            mtr_port = None

        if 'li_local_iface' in config_dict:
            li_port_name = config_dict['li_local_iface']

        enable_nat = config_dict.get('enable_nat', True)
        non_nat_gw_probe_freq = config_dict.get(
            'non_nat_gw_probe_frequency',
            self.ARP_PROBE_FREQUENCY,
        )
        # In case of vlan tag on uplink_bridge, use separate port.
        sgi_vlan = config_dict.get('sgi_management_iface_vlan', "")
        if not sgi_vlan:
            non_nat_arp_egress_port = config_dict.get(
                'non_nat_arp_egress_port',
                self.UPLINK_OVS_BRIDGE_NAME,
            )
        else:
            non_nat_arp_egress_port = config_dict.get(
                'non_nat_arp_egress_port',
                self.NON_NAT_ARP_EGRESS_PORT,
            )
        virtual_iface = config_dict.get('virtual_interface', None)
        if enable_nat is True or setup_type != 'LTE':
            if virtual_iface is not None:
                virtual_mac = get_virtual_iface_mac(virtual_iface)
            else:
                virtual_mac = ""
        else:
            # override virtual mac from config file.
            virtual_mac = config_dict.get('virtual_mac', "")

        uplink_gw_mac = config_dict.get(
            'uplink_gw_mac',
            "ff:ff:ff:ff:ff:ff",
        )
        return self.InOutConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            uplink_port=port_no,
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            li_port_name=li_port_name,
            enable_nat=enable_nat,
            non_nat_gw_probe_frequency=non_nat_gw_probe_freq,
            non_nat_arp_egress_port=non_nat_arp_egress_port,
            setup_type=setup_type,
            uplink_gw_mac=uplink_gw_mac,
            he_proxy_port=he_proxy_port,
            he_proxy_eth_mac=he_proxy_eth_mac,
            mtr_mac=mtr_mac,
            virtual_mac=virtual_mac,
        )

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self._setup_non_nat_monitoring()
        # TODO possibly investigate stateless XWF(no sessiond)
        if self.config.setup_type == 'XWF':
            self.delete_all_flows(datapath)
            self._install_default_flows(datapath)

    def _get_default_flow_msgs(self, datapath) -> DefaultMsgsMap:
        """
        Gets the default flow msgs for pkt routing

        Args:
            datapath: ryu datapath struct
        Returns:
            The list of default msgs to add
        """
        return {
            self._ingress_tbl_num: self._get_default_ingress_flow_msgs(datapath),
            self._midle_tbl_num: self._get_default_middle_flow_msgs(datapath),
            self._egress_tbl_num: self._get_default_egress_flow_msgs(datapath, mac_addr=self.config.virtual_mac),
        }

    def _install_default_flows(self, datapath):
        default_msg_map = self._get_default_flow_msgs(datapath)
        default_msgs = []

        for _, msgs in default_msg_map.items():
            default_msgs.extend(msgs)
        chan = self._msg_hub.send(default_msgs, datapath)
        self._wait_for_responses(chan, len(default_msgs))

    def cleanup_on_disconnect(self, datapath):
        if self._clean_restart:
            self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self._ingress_tbl_num)
        flows.delete_all_flows_from_table(datapath, self._midle_tbl_num)
        flows.delete_all_flows_from_table(datapath, self._egress_tbl_num)

    def _get_default_middle_flow_msgs(self, dp):
        """
        Egress table is the last table that a packet touches in the pipeline.
        Output downlink traffic to gtp port, uplink trafic to LOCAL

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        msgs = []
        next_tbl = self._service_manager.get_next_table_num(PHYSICAL_TO_LOGICAL)

        # Allow passthrough pkts(skip enforcement and send to egress table)
        ps_match = MagmaMatch(passthrough=PASSTHROUGH_REG_VAL)
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp,
                self._midle_tbl_num, ps_match, actions=[],
                priority=flows.PASSTHROUGH_PRIORITY,
                resubmit_table=self._egress_tbl_num,
            ),
        )

        match = MagmaMatch()
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp,
                self._midle_tbl_num, match, actions=[],
                priority=flows.DEFAULT_PRIORITY, resubmit_table=next_tbl,
            ),
        )

        if self._mtr_service_enabled:
            msgs.extend(
                _get_vlan_egress_flow_msgs(
                    dp,
                    self._midle_tbl_num,
                    ether_types.ETH_TYPE_IP,
                    self.config.mtr_ip,
                    self.config.mtr_port,
                    priority=flows.UE_FLOW_PRIORITY,
                    direction=Direction.OUT,
                    dst_mac=self.config.mtr_mac,
                ),
            )
        return msgs

    def _get_default_egress_flow_msgs(
        self, dp, mac_addr: str = "", vlan: str = "",
        ipv6: bool = False,
    ):
        """
        Egress table is the last table that a packet touches in the pipeline.
        Output downlink traffic to gtp port, uplink trafic to LOCAL
        Args:
            mac_addr: In Non NAT mode, this is upstream internet GW mac address
            vlan: in multi APN this is vlan_id of the upstream network.

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        msgs = []
        if self.config.setup_type == 'LTE':
            msgs.extend(
                _get_vlan_egress_flow_msgs(
                    dp,
                    self._egress_tbl_num,
                    ether_types.ETH_TYPE_IP,
                    None,
                ),
            )
            msgs.extend(
                _get_vlan_egress_flow_msgs(
                    dp,
                    self._egress_tbl_num,
                    ether_types.ETH_TYPE_IPV6,
                    None,
                ),
            )
            msgs.extend(self._get_proxy_flow_msgs(dp))
        else:
            # Use regular match for Non LTE setup.
            downlink_match = MagmaMatch(direction=Direction.IN)
            msgs.append(
                flows.get_add_output_flow_msg(
                    dp, self._egress_tbl_num, downlink_match, [],
                    output_port=self.config.gtp_port,
                ),
            )

        if ipv6:
            uplink_match = MagmaMatch(
                eth_type=ether_types.ETH_TYPE_IPV6,
                direction=Direction.OUT,
            )
        elif vlan.isdigit():
            vid = 0x1000 | int(vlan)
            uplink_match = MagmaMatch(
                direction=Direction.OUT,
                vlan_vid=(vid, 0x1fff),
            )
        else:
            uplink_match = MagmaMatch(direction=Direction.OUT)
        actions = []
        # avoid resetting mac address on switch connect event.
        if mac_addr == "":
            mac_addr = self._current_upstream_mac_map.get(vlan, "")
        if mac_addr == "" and self.config.enable_nat is False and \
            self.config.setup_type == 'LTE':
            mac_addr = self.config.uplink_gw_mac

        if mac_addr != "":
            parser = dp.ofproto_parser
            actions.append(
                parser.NXActionRegLoad2(
                    dst='eth_dst',
                    value=mac_addr,
                ),
            )
            upstream_mac_key = vlan + '_' + str(ipv6)
            if self._current_upstream_mac_map.get(upstream_mac_key, "") != mac_addr:
                self.logger.info(
                    "Using GW: mac: %s match %s actions: %s",
                    mac_addr,
                    str(uplink_match.ryu_match),
                    str(actions),
                )

                self._current_upstream_mac_map[upstream_mac_key] = mac_addr

        if vlan.isdigit():
            priority = flows.UE_FLOW_PRIORITY
        elif mac_addr != "":
            priority = flows.DEFAULT_PRIORITY
        else:
            priority = flows.MINIMUM_PRIORITY

        if ipv6:
            # IPV6 flows would have higher priority than all IPv4
            priority += flows.UE_FLOW_PRIORITY

        msgs.append(
            flows.get_add_output_flow_msg(
                dp, self._egress_tbl_num, uplink_match, priority=priority,
                actions=actions, output_port=self.config.uplink_port,
            ),
        )

        return msgs

    def _get_default_ingress_flow_msgs(self, dp):
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
        msgs = []

        # set traffic direction bits

        # set a direction bit for incoming (internet -> UE) traffic.
        match = MagmaMatch(in_port=OFPP_LOCAL)
        actions = [load_direction(parser, Direction.IN)]
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp,
                self._ingress_tbl_num, match, actions=actions,
                priority=flows.DEFAULT_PRIORITY, resubmit_table=next_table,
            ),
        )

        # set a direction bit for incoming (internet -> UE) traffic.
        match = MagmaMatch(in_port=self.config.uplink_port)
        actions = [load_direction(parser, Direction.IN)]
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp, self._ingress_tbl_num, match,
                actions=actions,
                priority=flows.DEFAULT_PRIORITY,
                resubmit_table=next_table,
            ),
        )

        # Send RADIUS requests directly to li table
        if self._li_port:
            match = MagmaMatch(in_port=self._li_port)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num,
                    match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                    resubmit_table=self._li_table,
                ),
            )

        # set a direction bit for incoming (mtr -> UE) traffic.
        if self._mtr_service_enabled:
            match = MagmaMatch(in_port=self.config.mtr_port)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num,
                    match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                    resubmit_table=next_table,
                ),
            )

        if self.config.he_proxy_port != 0:
            match = MagmaMatch(in_port=self.config.he_proxy_port)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num,
                    match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                    resubmit_table=next_table,
                ),
            )

        if self.config.setup_type == 'CWF':
            # set a direction bit for outgoing (pn -> inet) traffic for remaining traffic
            ps_match_out = MagmaMatch()
            actions = [load_direction(parser, Direction.OUT)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num, ps_match_out,
                    actions=actions,
                    priority=flows.MINIMUM_PRIORITY,
                    resubmit_table=next_table,
                ),
            )
        else:
            # set a direction bit for outgoing (pn -> inet) traffic for remaining traffic
            # Passthrough is zero for packets from eNodeB GTP tunnels
            ps_match_out = MagmaMatch(passthrough=REG_ZERO_VAL)
            actions = [load_direction(parser, Direction.OUT)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num, ps_match_out,
                    actions=actions,
                    priority=flows.MINIMUM_PRIORITY,
                    resubmit_table=next_table,
                ),
            )

            # Passthrough is one for packets from remote PGW GTP tunnels, set direction
            # flag to IN for such packets.
            ps_match_in = MagmaMatch(passthrough=PASSTHROUGH_REG_VAL)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num, ps_match_in,
                    actions=actions,
                    priority=flows.MINIMUM_PRIORITY,
                    resubmit_table=next_table,
                ),
            )

        return msgs

    def _get_gw_mac_address_v4(self, ip: IPAddress, vlan: str = "") -> str:
        try:
            gw_ip = ipaddress.ip_address(ip.address)
            self.logger.debug(
                "sending arp via egress: %s",
                self.config.non_nat_arp_egress_port,
            )
            eth_mac_src = get_if_hwaddr(self.config.non_nat_arp_egress_port)
            psrc = "0.0.0.0"
            egress_port_ip = get_if_addr(self.config.non_nat_arp_egress_port)
            if egress_port_ip:
                psrc = egress_port_ip

            pkt = Ether(dst=ETHER_BROADCAST, src=eth_mac_src)
            if vlan.isdigit():
                pkt /= Dot1Q(vlan=int(vlan))
            pkt /= ARP(op="who-has", pdst=gw_ip, hwsrc=eth_mac_src, psrc=psrc)
            self.logger.debug("ARP Req pkt %s", pkt.show(dump=True))

            res = srp1(
                pkt,
                type=ETH_P_ALL,
                iface=self.config.non_nat_arp_egress_port,
                timeout=1,
                verbose=0,
                nofilter=1,
                promisc=0,
            )

            if res is not None:
                self.logger.debug("ARP Res pkt %s", res.show(dump=True))
                if str(res[ARP].psrc) != str(gw_ip):
                    self.logger.warning(
                        "Unexpected IP in ARP response. expected: %s pkt: %s",
                        str(gw_ip),
                        res.show(dump=True),
                    )
                    return ""
                if vlan.isdigit():
                    if Dot1Q in res and str(res[Dot1Q].vlan) == vlan:
                        mac = res[ARP].hwsrc
                    else:
                        self.logger.warning(
                            "Unexpected vlan in ARP response. expected: %s pkt: %s",
                            vlan,
                            res.show(dump=True),
                        )
                        return ""
                else:
                    mac = res[ARP].hwsrc
                return mac
            else:
                self.logger.debug("Got Null response")
                return ""

        except Scapy_Exception as ex:
            self.logger.warning("Error in probing Mac address: err %s", ex)
            return ""
        except ValueError:
            self.logger.warning(
                "Invalid GW Ip address: [%s] or vlan %s",
                str(ip), vlan,
            )
            return ""

    def _get_gw_mac_address_v6(self, ip: IPAddress) -> str:
        try:
            gw_ip = ipaddress.ip_address(ip.address)
            mac = getmacbyip6(str(gw_ip))
            self.logger.debug("Got mac %s for IP: %s", mac, gw_ip)
            return mac

        except Scapy_Exception as ex:
            self.logger.warning("Error in probing Mac address: err %s", ex)
            return ""
        except ValueError:
            self.logger.warning(
                "Invalid GW Ip address: [%s]",
                str(ip),
            )
            return ""

    def _get_gw_mac_address(self, ip: IPAddress, vlan: str = "") -> str:
        if ip.version == IPAddress.IPV4:
            return self._get_gw_mac_address_v4(ip, vlan)
        if ip.version == IPAddress.IPV6:
            if vlan == "NO_VLAN":
                return self._get_gw_mac_address_v6(ip)
            else:
                gw_ip = ipaddress.ip_address(ip.address)
                self.logger.error("Not supported: GW IPv6: %s over vlan %d", str(gw_ip), vlan)
                return None

    def _monitor_and_update(self):
        while self._gw_mac_monitor_on:
            gw_info_list = get_mobilityd_gw_info()
            for gw_info in gw_info_list:
                if gw_info and gw_info.ip:
                    latest_mac_addr = self._get_gw_mac_address(gw_info.ip, gw_info.vlan)
                    if latest_mac_addr is None or latest_mac_addr == "":
                        latest_mac_addr = gw_info.mac
                    self.logger.debug("mac [%s] for vlan %s", latest_mac_addr, gw_info.vlan)
                    msgs = self._get_default_egress_flow_msgs(
                        self._datapath,
                        latest_mac_addr,
                        gw_info.vlan,
                        ipv6=(gw_info.ip.version == IPAddress.IPV6),
                    )

                    chan = self._msg_hub.send(msgs, self._datapath)
                    self._wait_for_responses(chan, len(msgs))

                    if latest_mac_addr and latest_mac_addr != "":
                        set_mobilityd_gw_info(
                            gw_info.ip,
                            latest_mac_addr,
                            gw_info.vlan,
                        )
                else:
                    self.logger.warning("No default GW found.")

            hub.sleep(self.config.non_nat_gw_probe_frequency)

    def _setup_non_nat_monitoring(self):
        """
        Setup egress flow to forward traffic to internet GW.
        Start a thread to figure out MAC address of uplink NAT gw.

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
            self.logger.info(
                "Non nat conf: egress port: %s, uplink: %s",
                self.config.non_nat_arp_egress_port,
                self.config.uplink_port,
            )

        self._gw_mac_monitor_on = True
        self._gw_mac_monitor = hub.spawn(self._monitor_and_update)

        threading.Event().wait(1)

    def _stop_gw_mac_monitor(self):
        if self._gw_mac_monitor:
            self._gw_mac_monitor_on = False
            self._gw_mac_monitor.wait()

    def _get_proxy_flow_msgs(self, dp):
        """
        Install egress flows
        Args:
            dp datapath
            table_no table to install flow
            out_port specify egress port, if None reg value is used
            priority flow priority
            direction packet direction.
        """
        if self.config.he_proxy_port <= 0:
            return []

        parser = dp.ofproto_parser
        match = MagmaMatch(proxy_tag=PROXY_TAG_TO_PROXY)
        actions = [
            parser.NXActionRegLoad2(
                dst='eth_dst',
                value=self.config.he_proxy_eth_mac,
            ),
        ]
        return [
            flows.get_add_output_flow_msg(
                dp, self._egress_tbl_num, match,
                priority=flows.UE_FLOW_PRIORITY, actions=actions,
                output_port=self.config.he_proxy_port,
            ),
        ]

    def _wait_for_responses(self, chan, response_count):
        def fail(err):
            self.logger.error("Failed to install rule with error: %s", err)

        for _ in range(response_count):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS msg channel")
            if not result.ok():
                return fail(result.exception())

    def _get_ue_specific_flow_msgs(self, _):
        return {}

    def finish_init(self, _):
        pass

    def cleanup_state(self):
        pass

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)


def _get_vlan_egress_flow_msgs(
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
