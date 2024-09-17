"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import threading
from collections import namedtuple

from lte.protos.mobilityd_pb2 import IPAddress
from magma.pipelined.app.base import MagmaController
from magma.pipelined.app.restart_mixin import DefaultMsgsMap, RestartMixin
from magma.pipelined.bridge_util import BridgeTools, DatapathLookupError
from magma.pipelined.gw_mac_address import get_gw_mac_address
from magma.pipelined.ifaces import get_mac_address_from_iface
from magma.pipelined.mobilityd_client import (
    get_mobilityd_gw_info,
    set_mobilityd_gw_info,
)
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MessageHub, MsgChannel
from magma.pipelined.openflow.registers import PROXY_TAG_TO_PROXY, Direction
from magma.pipelined.vlan_utils import get_vlan_egress_flow_msgs
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib import hub
from ryu.lib.packet import ether_types

EGRESS = "egress"


class EgressController(RestartMixin, MagmaController):
    APP_NAME = "egress"

    EgressConfig = namedtuple(
        'EgressConfig',
        [
            'gtp_port', 'uplink_port', 'li_port_name',
            'enable_nat', 'non_nat_gw_probe_frequency', 'non_nat_arp_egress_port',
            'setup_type', 'uplink_gw_mac', 'he_proxy_port', 'he_proxy_eth_mac',
            'virtual_mac',
        ],
    )
    ARP_PROBE_FREQUENCY = 300
    NON_NAT_ARP_EGRESS_PORT = 'dhcp0'
    UPLINK_OVS_BRIDGE_NAME = 'uplink_br0'
    PROXY_PORT_MAC = 'e6:8f:a2:80:80:80'

    def __init__(self, *args, **kwargs):
        super(EgressController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.logger.info("egress config: %s", self.config)
        self._egress_tbl_num = self._service_manager.get_table_num(EGRESS)
        # following fields are only used in Non Nat config
        self._clean_restart = kwargs['config']['clean_restart']
        self._msg_hub = MessageHub(self.logger)
        self._datapath = None
        self.tbl_num = self._egress_tbl_num
        self._gw_mac_monitor = None
        self._current_upstream_mac_map = {}  # maps vlan to upstream gw mac
        self._gw_mac_monitor_on = False

    def _get_config(self, config_dict):
        port_no = config_dict.get('uplink_port', None)

        he_proxy_port = 0
        he_proxy_eth_mac = ''
        try:
            if 'proxy_port_name' in config_dict:
                he_proxy_port = BridgeTools.get_ofport(config_dict.get('proxy_port_name'))
                he_proxy_eth_mac = config_dict.get('he_proxy_eth_mac', self.PROXY_PORT_MAC)
        except DatapathLookupError:
            # ignore it
            self.logger.debug("could not parse proxy port config")

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
        setup_type = config_dict.get('setup_type', None)
        if enable_nat is True or setup_type != 'LTE':
            if virtual_iface is not None:
                virtual_mac = get_mac_address_from_iface(virtual_iface)
            else:
                virtual_mac = ""
        else:
            # override virtual mac from config file.
            virtual_mac = config_dict.get('virtual_mac', "")

        uplink_gw_mac = config_dict.get(
            'uplink_gw_mac',
            "ff:ff:ff:ff:ff:ff",
        )

        li_port_name = None
        if 'li_local_iface' in config_dict:
            li_port_name = config_dict['li_local_iface']

        return self.EgressConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            uplink_port=port_no,
            li_port_name=li_port_name,
            enable_nat=enable_nat,
            non_nat_gw_probe_frequency=non_nat_gw_probe_freq,
            non_nat_arp_egress_port=non_nat_arp_egress_port,
            setup_type=setup_type,
            uplink_gw_mac=uplink_gw_mac,
            he_proxy_port=he_proxy_port,
            he_proxy_eth_mac=he_proxy_eth_mac,
            virtual_mac=virtual_mac,
        )

    def _get_default_flow_msgs(self, datapath) -> DefaultMsgsMap:
        return {
            self._egress_tbl_num: self._get_default_egress_flow_msgs(datapath, mac_addr=self.config.virtual_mac),
        }

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self._setup_non_nat_monitoring()

    def _get_default_egress_flow_msgs(
        self, dp, mac_addr: str = "", vlan: str = "",
        ipv6: bool = False,
    ):
        """
        Egress table is the last table that a packet touches in the pipeline.
        Output downlink traffic to gtp port, uplink traffic to LOCAL
        Args:
            mac_addr: In Non NAT mode, this is upstream internet GW mac address
            vlan: in multi APN this is vlan_id of the upstream network.

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        msgs = []
        if self.config.setup_type == 'LTE':
            msgs.extend(
                get_vlan_egress_flow_msgs(
                    dp,
                    self._egress_tbl_num,
                    ether_types.ETH_TYPE_IP,
                    None,
                ),
            )
            msgs.extend(
                get_vlan_egress_flow_msgs(
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

    def _monitor_and_update(self):
        while self._gw_mac_monitor_on:
            gw_info_list = get_mobilityd_gw_info()
            for gw_info in gw_info_list:
                if gw_info and gw_info.ip:
                    latest_mac_addr = get_gw_mac_address(gw_info.ip, gw_info.vlan, self.config.non_nat_arp_egress_port)
                    if latest_mac_addr == "":
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

                    if latest_mac_addr:
                        set_mobilityd_gw_info(
                            gw_info.ip,
                            latest_mac_addr,
                            gw_info.vlan,
                        )
                else:
                    self.logger.warning("No default GW found.")

            hub.sleep(self.config.non_nat_gw_probe_frequency)

    def _stop_gw_mac_monitor(self):
        if self._gw_mac_monitor:
            self._gw_mac_monitor_on = False
            self._gw_mac_monitor.wait()

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

    def cleanup_on_disconnect(self, datapath):
        if self._clean_restart:
            self.delete_all_flows(datapath)

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    def cleanup_state(self):
        pass

    def _get_ue_specific_flow_msgs(self, _):
        return {}

    def finish_init(self, _):
        pass

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
