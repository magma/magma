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
from collections import namedtuple
from threading import Lock
from typing import List

from ryu.lib.packet import ether_types
from ryu.lib.packet.in_proto import IPPROTO_TCP

from .base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import load_direction, Direction, \
    load_passthrough, set_proxy_tag, set_in_port, load_imsi, \
    PROXY_TAG_TO_PROXY, set_tun_id
from magma.pipelined.envoy_client import activate_he_urls_for_ue, \
    deactivate_he_urls_for_ue

from lte.protos.mobilityd_pb2 import IPAddress
from magma.pipelined.bridge_util import BridgeTools, DatapathLookupError
from magma.pipelined.policy_converters import get_ue_ip_match_args, \
    get_eth_type, convert_ipv4_str_to_ip_proto, ipv4_address_to_str

import logging

PROXY_PORT_NAME = 'proxy_port'
HTTP_PORT = 80
PROXY_TABLE = 'proxy'


class UeProxyRuleCounter:
    def __init__(self):
        self._map = {}
        self._lock = Lock()

    def inc(self, ue_ip: str):
        with self._lock:
            cnt = self._map.get(ue_ip, 0)
            cnt = cnt + 1
            self._map[ue_ip] = cnt

    def get(self, ue_ip: str) -> int:
        with self._lock:
            return self._map.get(ue_ip, 0)

    def dec(self, ue_ip: str) -> bool:
        with self._lock:
            cnt = self._map.get(ue_ip, 0)
            if cnt == 0:
                return False
            cnt = cnt - 1
            if cnt == 0:
                self._map.pop(ue_ip)
            else:
                self._map[ue_ip] = cnt
            return True

    def delete(self, ue_ip: str):
        with self._lock:
            self._map.pop(ue_ip, 0)

    def clear(self):
        with self._lock:
            self._map.clear()


class HeaderEnrichmentController(MagmaController):
    """
    A controller that tags related HTTP proxy flows.

        1. From UE to Proxy
        2. From Proxy to UE
        3. From Proxy to Upstream server
        4. From upstream server to Proxy

    This controller is also responsible for setting direction for traffic
    egressing proxy_port.
    details are in: docs/readmes/proposals/p006_header_enrichment.md

    """

    APP_NAME = "proxy"
    APP_TYPE = ControllerType.PHYSICAL

    UplinkConfig = namedtuple(
        'heConfig',
        ['he_proxy_port',
         'he_enabled',
         'uplink_port',
         'gtp_port'],
    )

    def __init__(self, *args, **kwargs):
        super(HeaderEnrichmentController, self).__init__(*args, **kwargs)
        self._datapath = None

        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        self.config = self._get_config(kwargs['config'])
        self._ue_rule_counter = UeProxyRuleCounter()
        self.logger.info("Header Enrichment app config: %s", self.config)

    def _get_config(self, config_dict) -> namedtuple:
        try:
            he_proxy_port = BridgeTools.get_ofport(config_dict.get('proxy_port_name'))

            he_enabled = config_dict.get('he_enabled', True)
            uplink_port = config_dict.get('uplink_port', None)
        except DatapathLookupError:
            he_enabled = False
            uplink_port = 0
            he_proxy_port = 0

        return self.UplinkConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            he_proxy_port=he_proxy_port,
            he_enabled=he_enabled,
            uplink_port=uplink_port)

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self._install_default_flows(self._datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        self._ue_rule_counter.clear()

    def _install_default_flows(self, dp):
        match = MagmaMatch(in_port=self.config.he_proxy_port)
        flows.add_drop_flow(dp, self.tbl_num, match,
                            priority=flows.MINIMUM_PRIORITY + 1)
        match = MagmaMatch()
        flows.add_resubmit_next_service_flow(dp, self.tbl_num, match,
                                             [],
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def _set_he_target_urls(self, ue_addr: str, rule_id: str, urls: List[str], imsi: str, msisdn: bytes) -> bool:
        if msisdn:
            msisdn_str = msisdn.decode("utf-8")
        else:
            msisdn_str = None
        ip_addr = convert_ipv4_str_to_ip_proto(ue_addr)
        return activate_he_urls_for_ue(ip_addr, rule_id, urls, imsi, msisdn_str)

    def get_subscriber_he_flows(self, rule_id: str, direction: Direction,
                                ue_addr: str, uplink_tunnel: int, ip_dst: str,
                                rule_num: int, urls: List[str], imsi: str,
                                msisdn: bytes):
        """
        Add flow to steer traffic to and from proxy port.
        Args:
            rule_id(str) Rule id
            direction(Direction): HE rules are only used for upstream traffic.
            ue_addr(str): IP address of UE
            uplink_tunnel(int) Tunnel ID of the session
            ip_dst(str): HTTP server dst IP (CIDR)
            rule_num(int): rule num of the policy rule
            urls(List[str]): list of HTTP server URLs
            imsi (string): subscriber to install rule for
            msisdn (bytes): subscriber MSISDN
        """
        if not self.config.he_enabled:
            return []

        if direction != Direction.OUT:
            return []

        dp = self._datapath
        parser = dp.ofproto_parser
        tunnel_id = 0
        try:
            if uplink_tunnel:
                tunnel_id = int(uplink_tunnel)
        except ValueError:
            self.logger.error("parsing tunnel id: [%s], HE might not work in every case", uplink_tunnel)

        if urls is None or len(urls) == 0:
            return []

        if ip_dst is None or ip_dst == '':
            logging.error("Missing dst ip, ignoring HE rule.")
            return []

        logging.info("Add HE: ue_addr %s, rule_id: %s, uplink_tunnel %s, ip_dst %s, rule_num %s "
                     "urls %s, imsi %s, msisdn %s", ue_addr, rule_id, uplink_tunnel, ip_dst,
                     str(rule_num), str(urls), imsi, str(msisdn))

        success = self._set_he_target_urls(ue_addr, rule_id, urls, imsi, msisdn)
        if not success:
            return []
        msgs = []
        # 1.a. Going to UE: from uplink send to proxy
        match = MagmaMatch(in_port=self.config.uplink_port,
                           eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ip_dst,
                           ipv4_dst=ue_addr,
                           ip_proto=IPPROTO_TCP,
                           tcp_src=HTTP_PORT,
                           proxy_tag=0)
        actions = [load_direction(parser, Direction.IN),
                   load_passthrough(parser),
                   set_proxy_tag(parser)]
        msgs.append(
            flows.get_add_resubmit_current_service_flow_msg(dp, self.tbl_num,
                                                            match, cookie=rule_num,
                                                            actions=actions,
                                                            priority=flows.DEFAULT_PRIORITY,
                                                            resubmit_table=self.next_table))
        # 1.b. Going to UE: from proxy send to UE
        match = MagmaMatch(in_port=self.config.he_proxy_port,
                           eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ip_dst,
                           ipv4_dst=ue_addr,
                           ip_proto=IPPROTO_TCP,
                           tcp_src=HTTP_PORT)
        actions = [set_in_port(parser, self.config.uplink_port),
                   set_proxy_tag(parser)]
        msgs.append(
            flows.get_add_resubmit_current_service_flow_msg(dp, self.tbl_num,
                                                            match, cookie=rule_num,
                                                            actions=actions,
                                                            priority=flows.DEFAULT_PRIORITY,
                                                            resubmit_table=0))

        # 1.c. continue (1.b) Going to UE: from proxy send to UE
        match = MagmaMatch(in_port=self.config.uplink_port,
                           eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ip_dst,
                           ipv4_dst=ue_addr,
                           ip_proto=IPPROTO_TCP,
                           tcp_src=HTTP_PORT,
                           proxy_tag=PROXY_TAG_TO_PROXY)
        actions = [set_proxy_tag(parser, 0)]
        msgs.append(
            flows.get_add_resubmit_current_service_flow_msg(dp, self.tbl_num,
                                                            match, cookie=rule_num,
                                                            actions=actions,
                                                            priority=flows.DEFAULT_PRIORITY,
                                                            resubmit_table=self.next_table))

        # 2.a. To internet from proxy port, send to uplink
        match = MagmaMatch(in_port=self.config.he_proxy_port,
                           eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ue_addr,
                           ipv4_dst=ip_dst,
                           ip_proto=IPPROTO_TCP,
                           tcp_dst=HTTP_PORT,
                           proxy_tag=0)
        actions = [set_in_port(parser, self.config.gtp_port),
                   set_tun_id(parser, tunnel_id),
                   set_proxy_tag(parser),
                   load_imsi(parser, imsi)]
        msgs.append(
            flows.get_add_resubmit_current_service_flow_msg(dp, self.tbl_num,
                                                            match,
                                                            cookie=rule_num,
                                                            actions=actions,
                                                            priority=flows.MEDIUM_PRIORITY,
                                                            resubmit_table=0))

        # 2.b. Continue from 2.a -> To internet from proxy port, send to uplink
        match = MagmaMatch(in_port=self.config.gtp_port,
                           eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ue_addr,
                           ipv4_dst=ip_dst,
                           ip_proto=IPPROTO_TCP,
                           tcp_dst=HTTP_PORT,
                           proxy_tag=PROXY_TAG_TO_PROXY)
        actions = [set_proxy_tag(parser, 0)]
        msgs.append(
            flows.get_add_resubmit_current_service_flow_msg(dp, self.tbl_num,
                                                            match,
                                                            cookie=rule_num,
                                                            actions=actions,
                                                            priority=flows.DEFAULT_PRIORITY,
                                                            resubmit_table=self.next_table))

        # 2.c. To internet from ue send to proxy, this is coming from HE port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ue_addr,
                           ipv4_dst=ip_dst,
                           ip_proto=IPPROTO_TCP,
                           tcp_dst=HTTP_PORT,
                           proxy_tag=0)
        actions = [load_direction(parser, Direction.OUT),
                   load_passthrough(parser),
                   set_proxy_tag(parser)]
        msgs.append(
            flows.get_add_resubmit_current_service_flow_msg(dp, self.tbl_num,
                                                            match, cookie=rule_num,
                                                            actions=actions,
                                                            priority=flows.DEFAULT_PRIORITY,
                                                            resubmit_table=self.next_table))
        self._ue_rule_counter.inc(ue_addr)
        return msgs

    def remove_subscriber_he_flows(self, ue_addr: IPAddress, rule_id: str = "",
                                   rule_num: int = -1):
        """
        Remove proxy flows of give policy rule of the subscriber.
        Args:
            ue_addr(str): IP address of UE
            rule_id(str) Rule id
            rule_num(int): rule num of the policy rule
        """
        ue_ip_str = ipv4_address_to_str(ue_addr)

        if self._ue_rule_counter.get(ue_ip_str) == 0:
            return
        logging.info("Del HE rule: ue-ip: %s rule_id: %s rule %d",
                     ue_addr, rule_id, rule_num)

        if rule_num == -1:
            ip_match_in = get_ue_ip_match_args(ue_addr, Direction.IN)
            match_in = MagmaMatch(eth_type=get_eth_type(ue_addr),
                                  **ip_match_in)
            flows.delete_flow(self._datapath, self.tbl_num, match_in)

            ip_match_out = get_ue_ip_match_args(ue_addr, Direction.OUT)
            match_out = MagmaMatch(eth_type=get_eth_type(ue_addr),
                                   **ip_match_out)
            flows.delete_flow(self._datapath, self.tbl_num, match_out)
        else:
            match = MagmaMatch()
            flows.delete_flow(self._datapath, self.tbl_num, match,
                              cookie=rule_num, cookie_mask=flows.OVS_COOKIE_MATCH_ALL)

        success = deactivate_he_urls_for_ue(ue_addr, rule_id)
        if success:
            if rule_num == -1:
                self._ue_rule_counter.delete(ue_ip_str)
            else:
                self._ue_rule_counter.dec(ue_ip_str)
