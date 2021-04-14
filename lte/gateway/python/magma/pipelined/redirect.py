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

import asyncio
import ipaddress
from collections import namedtuple
from urllib.parse import urlsplit

import aiodns
import netifaces
from magma.configuration.service_configs import get_service_config_value
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import (
    DIRECTION_REG,
    IMSI_REG,
    REG_ZERO_VAL,
    RULE_NUM_REG,
    RULE_VERSION_REG,
    SCRATCH_REGS,
    Direction,
)
from magma.redirectd.redirect_store import RedirectDict
from memoize import Memoizer
from redis import RedisError
from ryu.lib.packet import ether_types
from ryu.ofproto.inet import IPPROTO_TCP, IPPROTO_UDP
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL


class RedirectException(Exception):
    pass


class RedirectionManager:
    """
    RedirectionManager

    The redirection manager handles subscribers who have redirection enabled,
    it adds the flows into ovs for redirecting user to the redirection server
    """
    DNS_TIMEOUT_SECS = 15

    REDIRECT_NOT_PROCESSED = REG_ZERO_VAL
    REDIRECT_PROCESSED = 0x1

    RedirectRequest = namedtuple(
        'RedirectRequest',
        ['imsi', 'ip_addr', 'rule', 'rule_num', 'rule_version', 'priority'],
    )

    def __init__(self, bridge_ip, logger, main_tbl_num, stats_table, next_table,
                 scratch_table_num, session_rule_version_mapper):
        self._bridge_ip = bridge_ip
        self.logger = logger
        self.main_tbl_num = main_tbl_num
        self.stats_table = stats_table
        self.next_table = next_table
        self._scratch_tbl_num = scratch_table_num
        self._redirect_dict = RedirectDict()
        self._dns_cache = Memoizer({})
        self._redirect_port = get_service_config_value(
            'redirectd', 'http_port', 8080)
        self._session_rule_version_mapper = session_rule_version_mapper

        self._cwf_args_set = False
        self._mac_rewrite_scratch = None
        self._internal_ip_allocator = None
        self._arpd_controller_fut = None
        self._arp_contoller = None
        self._egress_table = None
        self._bridge_mac = None

    def set_cwf_args(self, internal_ip_allocator, arp, mac_rewrite,
                     bridge_name, egress_table):
        self._mac_rewrite_scratch = mac_rewrite
        self._internal_ip_allocator = internal_ip_allocator
        self._arpd_controller_fut = arp
        self._arp_contoller = None
        self._egress_table = egress_table

        def get_virtual_iface_mac(iface):
            virt_ifaddresses = netifaces.ifaddresses(iface)
            return virt_ifaddresses[netifaces.AF_LINK][0]['addr']
        self._bridge_mac = get_virtual_iface_mac(bridge_name)
        self._cwf_args_set = True
        return self

    def setup_lte_redirect(self, datapath, loop, redirect_request):
        """
        Depending on redirection server address type install redirection rules
        """
        imsi = redirect_request.imsi
        ip_addr = redirect_request.ip_addr
        rule = redirect_request.rule
        rule_num = redirect_request.rule_num
        rule_version = redirect_request.rule_version
        priority = redirect_request.priority

        # TODO IMPORTANT check that redirectd service is running, as its a
        # dynamic service its not on by default. Will save you some sanity :)

        # TODO figure out what to do with SIP_URI
        if rule.redirect.address_type == rule.redirect.SIP_URI:
            raise RedirectException("SIP_URIs redirection isn't setup")
        if rule.redirect.address_type == rule.redirect.IPv6:
            raise RedirectException("No ipv6 support, so no ipv6 redirect")

        self._save_redirect_entry(ip_addr, rule.redirect)
        self._install_redirect_flows(datapath, loop, imsi, ip_addr, rule,
                                     rule_num, rule_version, priority)
        return

    def _install_redirect_flows(self, datapath, loop, imsi, ip_addr, rule,
                                rule_num, rule_version, priority):
        """
        Add flows to forward traffic to the redirection server.

        1) Intercept tcp traffic to the web to the redirection server, which
            completes the tcp handshake. This is done by adding an OVS flow
            with a learn action (flow catches inbound tcp packets, while learn
            action creates another flow that sends packets back from server)

        2) Add flows to allow UDP traffic so DNS queries can go through.
            Finally add flows with a higher priority that allow traffic to and
            from the address provided in redirect rule.
        """

        if rule.redirect.address_type == rule.redirect.URL:
            self._install_url_bypass_flows(datapath, loop, imsi, rule,
                                           rule_num, rule_version, priority,
                                           ue_ip=ip_addr)
        elif rule.redirect.address_type == rule.redirect.IPv4:
            self._install_ipv4_bypass_flows(datapath, imsi, rule,
                                            rule_num, rule_version, priority,
                                            [rule.redirect.server_address],
                                            ue_ip=ip_addr)

        self._install_dns_flows(datapath, imsi, rule, rule_num, rule_version,
                                priority)
        self._install_server_flows(datapath, imsi, ip_addr, rule, rule_num,
                                   rule_version, priority)

    def _install_scratch_table_flows(self, datapath, imsi, rule, rule_num,
                                     rule_version, priority):
        """
        The flow action for subscribers that need to be redirected does 2 things
            * Forward requests from subscriber to the internal http server
            * Instantiate a flow that matches response packets from the server
              and sends them back to subscriber
        Match: incoming tcp traffic with port 80, direction out
        Action:
            1) Set reg2 to rule_num
            2) Set ip dst to server ip
            3) Output to table 20
            4) Apply LearnAction:
            LearnAction(adds new flow for every pkt flow that hits this rule)
                1) Match ip packets
                2) Match tcp protocol
                3) Match packets from LOCAL port
                4) Match ip src = server ip
                5) Match ip dst = current flow ip src
                6) Match tcp src = current flow tcp dst
                7) Match tcp dst = current flow tcp src
                8) Load ip src = current flow ip dst
                9) Output through gtp0
        """
        parser = datapath.ofproto_parser
        match_http = MagmaMatch(
            eth_type=ether_types.ETH_TYPE_IP, ip_proto=IPPROTO_TCP,
            tcp_dst=80, imsi=encode_imsi(imsi), direction=Direction.OUT)
        of_note = parser.NXActionNote(list(rule.id.encode()))

        actions = [
            parser.NXActionLearn(
                table_id=self.main_tbl_num,
                priority=priority,
                cookie=rule_num,
                specs=[
                    parser.NXFlowSpecMatch(
                        src=ether_types.ETH_TYPE_IP, dst=('eth_type_nxm', 0),
                        n_bits=16
                    ),
                    parser.NXFlowSpecMatch(
                        src=IPPROTO_TCP, dst=('ip_proto_nxm', 0), n_bits=8
                    ),
                    parser.NXFlowSpecMatch(
                        src=Direction.IN,
                        dst=(DIRECTION_REG, 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecMatch(
                        src=int(ipaddress.IPv4Address(self._bridge_ip)),
                        dst=('ipv4_src_nxm', 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecMatch(
                        src=('ipv4_src_nxm', 0),
                        dst=('ipv4_dst_nxm', 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecMatch(
                        src=('tcp_src_nxm', 0),
                        dst=('tcp_dst_nxm', 0),
                        n_bits=16
                    ),
                    parser.NXFlowSpecMatch(
                        src=self._redirect_port,
                        dst=('tcp_src_nxm', 0),
                        n_bits=16
                    ),
                    parser.NXFlowSpecMatch(
                        src=encode_imsi(imsi),
                        dst=(IMSI_REG, 0),
                        n_bits=64
                    ),
                    parser.NXFlowSpecLoad(
                        src=('ipv4_dst_nxm', 0),
                        dst=('ipv4_src_nxm', 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecLoad(
                        src=80,
                        dst=('tcp_src_nxm', 0),
                        n_bits=16
                    ),
                    # Learn doesn't support resubmit to table, so send directly
                    parser.NXFlowSpecOutput(
                        src=('in_port', 0), dst="", n_bits=16
                    ),
                ]
            ),
            parser.NXActionRegLoad2(dst=SCRATCH_REGS[0],
                                    value=self.REDIRECT_PROCESSED),
            parser.OFPActionSetField(ipv4_dst=self._bridge_ip),
            parser.OFPActionSetField(tcp_dst=self._redirect_port),
            of_note,
        ]
        actions += self._load_rule_actions(parser, rule_num, rule_version)

        flows.add_resubmit_current_service_flow(
            datapath, self._scratch_tbl_num, match_http, actions,
            priority=priority, cookie=rule_num, hard_timeout=rule.hard_timeout,
            resubmit_table=self.main_tbl_num)

        match = MagmaMatch(imsi=encode_imsi(imsi))
        action = []
        flows.add_drop_flow(datapath, self._scratch_tbl_num, match, action,
                            priority=flows.MINIMUM_PRIORITY + 1,
                            cookie=rule_num)

    def _install_not_processed_flows(self, datapath, imsi, ip_addr, rule,
                                     rule_num, priority):
        """
        Redirect all traffic to the scratch table to only allow redirected
        http traffic to go through, the rest will be dropped. reg0 is used as
        a boolean to know whether the drop rule was processed.
        """
        parser = datapath.ofproto_parser
        of_note = parser.NXActionNote(list(rule.id.encode()))

        match = MagmaMatch(imsi=encode_imsi(imsi),
                           direction=Direction.OUT,
                           reg0=self.REDIRECT_NOT_PROCESSED,
                           eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ip_addr)
        action = [of_note]
        flows.add_resubmit_current_service_flow(
            datapath, self.main_tbl_num, match, action, priority=priority,
            cookie=rule_num, hard_timeout=rule.hard_timeout,
            resubmit_table=self._scratch_tbl_num)

        match = MagmaMatch(imsi=encode_imsi(imsi),
                           direction=Direction.OUT,
                           reg0=self.REDIRECT_PROCESSED,
                           eth_type=ether_types.ETH_TYPE_IP,
                           ipv4_src=ip_addr)
        action = [of_note]
        flows.add_resubmit_next_service_flow(
            datapath, self.main_tbl_num, match, action, priority=priority,
            cookie=rule_num, hard_timeout=rule.hard_timeout,
            copy_table=self.stats_table, resubmit_table=self.next_table)

    def setup_cwf_redirect(self, datapath, loop, redirect_request):
        """
        Add flows to forward traffic to the redirection server for cwf networks

        1) Intercept tcp traffic to the web to the redirection server, which
            completes the tcp handshake. Also overwrite UE src ip to match the
            subnet of the redirection server.
            This is done by assigning an internal IP per each subscriber.
            Add an OVS flow with a learn action (flow catches inbound tcp http
            packets, while learn action creates another flow that rewrites
            packet back to send to ue)

        2) Add flows to allow UDP traffic so DNS queries can go through.
            Add flows with a higher priority that allow traffic to and
            from the address provided in redirect rule.

        TODO we might want to track stats for these rules and report to sessiond
        """
        if not self._cwf_args_set:
            raise RedirectException("Can't install cwf redirection, missing"
                                    "cwf specific args, call set_cwf_args()")
        imsi = redirect_request.imsi
        rule = redirect_request.rule
        rule_num = redirect_request.rule_num
        rule_version = redirect_request.rule_version
        priority = redirect_request.priority
        if rule.redirect.address_type == rule.redirect.URL:
            self._install_url_bypass_flows(datapath, loop, imsi, rule,
                                           rule_num, rule_version, priority)
        elif rule.redirect.address_type == rule.redirect.IPv4:
            self._install_ipv4_bypass_flows(datapath, imsi, rule,
                                            rule_num, rule_version, priority,
                                            [rule.redirect.server_address])

        parser = datapath.ofproto_parser
        # TODO use subscriber ip_addr to generate internal IP and release
        # internal IP when subscriber disconnects or redirection flow is removed
        internal_ip = self._internal_ip_allocator.next_ip()

        self._save_redirect_entry(internal_ip, rule.redirect)
        #TODO check if we actually need this, dns might already be allowed
        self._install_dns_flows(datapath, imsi, rule, rule_num, rule_version,
                                priority)

        match_tcp_80 = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.OUT, tcp_dst=80
        )
        match_tcp_8008 = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.OUT, tcp_dst=8080
        )
        match_tcp_8080 = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.OUT, tcp_dst=8008
        )
        actions = [
            parser.NXActionLearn(
                table_id=self._mac_rewrite_scratch,
                priority=flows.UE_FLOW_PRIORITY,
                cookie=rule_num,
                specs=[
                    parser.NXFlowSpecMatch(
                        src=ether_types.ETH_TYPE_IP, dst=('eth_type_nxm', 0),
                        n_bits=16
                    ),
                    parser.NXFlowSpecMatch(
                        src=IPPROTO_TCP,
                        dst=('ip_proto_nxm', 0), n_bits=8
                    ),
                    parser.NXFlowSpecMatch(
                        src=int(ipaddress.IPv4Address(self._bridge_ip)),
                        dst=('ipv4_src_nxm', 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecMatch(
                        src=int(internal_ip),
                        dst=('ipv4_dst_nxm', 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecMatch(
                        src=('tcp_src_nxm', 0),
                        dst=('tcp_dst_nxm', 0),
                        n_bits=16
                    ),
                    parser.NXFlowSpecMatch(
                        src=self._redirect_port,
                        dst=('tcp_src_nxm', 0),
                        n_bits=16
                    ),
                    parser.NXFlowSpecLoad(
                        src=('eth_src_nxm', 0),
                        dst=('eth_dst_nxm', 0),
                        n_bits=48
                    ),
                    parser.NXFlowSpecLoad(
                        src=encode_imsi(imsi),
                        dst=(IMSI_REG, 0),
                        n_bits=64
                    ),
                    parser.NXFlowSpecLoad(
                        src=('ipv4_src_nxm', 0),
                        dst=('ipv4_dst_nxm', 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecLoad(
                        src=('ipv4_dst_nxm', 0),
                        dst=('ipv4_src_nxm', 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecLoad(
                        src=('tcp_dst_nxm', 0),
                        dst=('tcp_src_nxm', 0),
                        n_bits=16
                    ),
                ]
            ),
            parser.OFPActionSetField(ipv4_src=str(internal_ip)),
            parser.OFPActionSetField(ipv4_dst=self._bridge_ip),
            parser.OFPActionSetField(eth_dst=self._bridge_mac),
            parser.OFPActionSetField(tcp_dst=self._redirect_port),
        ]
        flows.add_output_flow(
            datapath, self.main_tbl_num, match_tcp_80, actions,
            priority=flows.UE_FLOW_PRIORITY, cookie=rule_num,
            output_port=OFPP_LOCAL)
        flows.add_output_flow(
            datapath, self.main_tbl_num, match_tcp_8008, actions,
            priority=flows.UE_FLOW_PRIORITY, cookie=rule_num,
            output_port=OFPP_LOCAL)
        flows.add_output_flow(
            datapath, self.main_tbl_num, match_tcp_8080, actions,
            priority=flows.UE_FLOW_PRIORITY, cookie=rule_num,
            output_port=OFPP_LOCAL)

        # Add flows for vlan traffic too (we need to pop vlan for flask server)
        # In ryu vlan_vid=(0x1000, 0x1000) matches all vlans
        match_tcp_80_vlan = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.OUT, tcp_dst=80,
            vlan_vid=(0x1000, 0x1000)
        )
        match_tcp_8008_vlan = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.OUT, tcp_dst=8080,
            vlan_vid=(0x1000, 0x1000)
        )
        match_tcp_8080_vlan = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.OUT, tcp_dst=8008,
            vlan_vid=(0x1000, 0x1000)
        )
        actions.append(parser.OFPActionPopVlan())
        flows.add_output_flow(
            datapath, self.main_tbl_num, match_tcp_80_vlan, actions,
            priority=flows.UE_FLOW_PRIORITY + 1, cookie=rule_num,
            output_port=OFPP_LOCAL)
        flows.add_output_flow(
            datapath, self.main_tbl_num, match_tcp_8008_vlan, actions,
            priority=flows.UE_FLOW_PRIORITY + 1, cookie=rule_num,
            output_port=OFPP_LOCAL)
        flows.add_output_flow(
            datapath, self.main_tbl_num, match_tcp_8080_vlan, actions,
            priority=flows.UE_FLOW_PRIORITY + 1, cookie=rule_num,
            output_port=OFPP_LOCAL)

        # TODO cleanup, make this a default rule in the ue_mac table
        ue_tbl = 0
        ue_next_tbl = 1
        # Allows traffic back from the flask server
        match = MagmaMatch(in_port=OFPP_LOCAL)
        actions = [
            parser.NXActionResubmitTable(table_id=self._mac_rewrite_scratch)]
        flows.add_resubmit_next_service_flow(datapath, ue_tbl,
                                             match, actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=ue_next_tbl)
        match = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.IN, in_port=OFPP_LOCAL)
        flows.add_resubmit_next_service_flow(
            datapath, self.main_tbl_num, match, [],
            priority=flows.DEFAULT_PRIORITY, cookie=rule_num,
            resubmit_table=self._egress_table
        )

        # Mac doesn't matter as we rewrite it anwyays
        mac_addr = '01:02:03:04:05:06'
        if self._arp_contoller or self._arpd_controller_fut.done():
            if not self._arp_contoller:
                self._arp_contoller = self._arpd_controller_fut.result()
            self._arp_contoller.set_incoming_arp_flows(datapath, internal_ip,
                                                       mac_addr)

        # Drop all other traffic that doesn't match
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.add_drop_flow(datapath, self.main_tbl_num, match, [],
                            priority=flows.MINIMUM_PRIORITY + 1,
                            cookie=rule_num)

    def _install_server_flows(self, datapath, imsi, ip_addr, rule, rule_num,
                              rule_version, priority):
        """
        Install the redirect flows to redirect all HTTP traffic to the captive
        portal and to drop all of the rest.
        """
        self._install_scratch_table_flows(datapath, imsi, rule, rule_num,
                                          rule_version, priority)
        self._install_not_processed_flows(datapath, imsi, ip_addr, rule,
                                          rule_num, priority)

    def _install_url_bypass_flows(self, datapath, loop, imsi, rule, rule_num,
                                  rule_version, priority, ue_ip=None):
        """
        Resolve DNS queries to get the ip address of redirect url, this is done
        to allow traffic to safely pass through as we want subscribers to have
        full access to the url they are redirected to.

        First check cache for redirect url, if not in cache submit a DNS query
        """
        redirect_addr_host = urlsplit(rule.redirect.server_address).netloc
        cached_ips = self._dns_cache.get(redirect_addr_host)
        if cached_ips is not None:
            self.logger.debug(
                "DNS cache hit for {}, entry expires in {} sec".format(
                    redirect_addr_host, self._dns_cache.ttl(redirect_addr_host)
                )
            )
            self._install_ipv4_bypass_flows(datapath, imsi, rule, rule_num,
                                            rule_version, priority, cached_ips,
                                            ue_ip)
            return

        resolver = aiodns.DNSResolver(timeout=self.DNS_TIMEOUT_SECS, loop=loop)
        query = resolver.query(redirect_addr_host, 'A')

        def add_flows(dns_resolve_future):
            """
            Callback for when DNS query is resolved, adds the bypass flows
            """
            try:
                ips = [entry.host for entry in dns_resolve_future.result()]
                ttl = min(entry.ttl for entry in dns_resolve_future.result())
            except aiodns.error.DNSError as err:
                self.logger.error("Error: ip lookup for {}: {}".format(
                                  redirect_addr_host, err))
                return
            self._dns_cache.get(redirect_addr_host, lambda: ips, max_age=ttl)
            self._install_ipv4_bypass_flows(datapath, imsi, rule, rule_num,
                                            rule_version, priority, ips, ue_ip)

        asyncio.ensure_future(query, loop=loop).add_done_callback(add_flows)

    def _install_ipv4_bypass_flows(self, datapath, imsi, rule, rule_num,
                                   rule_version, priority, ips, ue_ip=None):
        """
        Installs flows for traffic that is allowed to pass through for
        subscriber who has redirection enabled. Allow access to all passed ips.

        Allow UDP traffic(for DNS queries), traffic to/from redirection address
        """
        parser = datapath.ofproto_parser
        of_note = parser.NXActionNote(list(rule.id.encode()))
        actions = [
            of_note,
        ]
        actions += self._load_rule_actions(parser, rule_num, rule_version)

        matches = []
        uplink_ip_match = {}
        downlink_ip_match = {}
        if ue_ip != None:
            uplink_ip_match['ipv4_src'] = ue_ip
            downlink_ip_match['ipv4_dst'] = ue_ip
        for ip in ips:
            matches.append(MagmaMatch(
                eth_type=ether_types.ETH_TYPE_IP, direction=Direction.OUT,
                ipv4_dst=ip, imsi=encode_imsi(imsi), **uplink_ip_match
            ))
            matches.append(MagmaMatch(
                eth_type=ether_types.ETH_TYPE_IP, direction=Direction.IN,
                ipv4_src=ip, imsi=encode_imsi(imsi), **downlink_ip_match
            ))
        for match in matches:
            flows.add_resubmit_next_service_flow(
                datapath, self.main_tbl_num, match, actions,
                priority=priority + 1, cookie=rule_num,
                hard_timeout=rule.hard_timeout,
                copy_table=self.stats_table, resubmit_table=self.next_table)

    def _install_dns_flows(self, datapath, imsi, rule, rule_num, rule_version,
                           priority):
        """
        Installs flows that allow DNS queries to path through.
        """
        parser = datapath.ofproto_parser
        of_note = parser.NXActionNote(list(rule.id.encode()))
        actions = [
            of_note,
        ]
        actions += self._load_rule_actions(parser, rule_num, rule_version)
        matches = []
        # Install UDP flows for DNS
        matches.append(MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_src=53,
                                  direction=Direction.IN,
                                  imsi=encode_imsi(imsi)))
        matches.append(MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_UDP,
                                  udp_dst=53,
                                  direction=Direction.OUT,
                                  imsi=encode_imsi(imsi)))
        # Install TCP flows for DNS
        matches.append(MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_TCP,
                                  tcp_src=53,
                                  direction=Direction.IN,
                                  imsi=encode_imsi(imsi)))
        matches.append(MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                  ip_proto=IPPROTO_TCP,
                                  tcp_dst=53,
                                  direction=Direction.OUT,
                                  imsi=encode_imsi(imsi)))
        for match in matches:
            flows.add_resubmit_next_service_flow(
                datapath, self.main_tbl_num, match, actions, priority=priority,
                cookie=rule_num, hard_timeout=rule.hard_timeout,
                copy_table=self.stats_table, resubmit_table=self.next_table)

    def _save_redirect_entry(self, ip_addr, redirect_info):
        """
        Saves the redirect entry in Redis.

        Throws:
            RedirectException: on error
        """
        try:
            # Verify if ip_addr is in the correct format, and also
            # normalize the variants of the address into a single format
            ip_str = str(ipaddress.ip_address(ip_addr))
        except ValueError as exp:
            raise RedirectException(exp)
        try:
            self._redirect_dict[ip_str] = redirect_info
        except RedisError as exp:
            raise RedirectException(exp)
        self.logger.info("Saved redirect rule for %s in Redis" % ip_str)

    def deactivate_flow_for_rule(self, datapath, imsi, rule_num):
        """
        Deactivate a specific rule using the flow cookie for a subscriber
        """
        cookie, mask = (rule_num, flows.OVS_COOKIE_MATCH_ALL)
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(datapath, self._scratch_tbl_num, match,
                          cookie=cookie, cookie_mask=mask)

    def deactivate_flows_for_subscriber(self, datapath, imsi):
        """
        Deactivate all rules for a subscriber
        """
        flows.delete_flow(datapath, self._scratch_tbl_num,
                          MagmaMatch(imsi=encode_imsi(imsi)))

    def _load_rule_actions(self, parser, rule_num, rule_version):
        return [
            parser.NXActionRegLoad2(dst=RULE_NUM_REG, value=rule_num),
            parser.NXActionRegLoad2(dst=RULE_VERSION_REG, value=rule_version),
        ]
