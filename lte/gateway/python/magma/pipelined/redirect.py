"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import aiodns
import asyncio
import ipaddress
from collections import namedtuple
from redis import RedisError
from urllib.parse import urlsplit

from magma.configuration.service_configs import get_service_config_value
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import IMSI_REG, DIRECTION_REG, \
    Direction, SCRATCH_REGS, REG_ZERO_VAL, RULE_VERSION_REG
from magma.redirectd.redirect_store import RedirectDict

from ryu.lib.packet import ether_types
from ryu.ofproto.inet import IPPROTO_TCP, IPPROTO_UDP
from memoize import Memoizer


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
        ['imsi', 'ip_addr', 'rule', 'rule_num', 'priority'],
    )

    def __init__(self, bridge_ip, logger, main_tbl_num, next_table,
                 scratch_table_num, session_rule_version_mapper):
        self._bridge_ip = bridge_ip
        self.logger = logger
        self.main_tbl_num = main_tbl_num
        self.next_table = next_table
        self._scratch_tbl_num = scratch_table_num
        self._redirect_dict = RedirectDict()
        self._dns_cache = Memoizer({})
        self._redirect_port = get_service_config_value(
            'redirectd', 'http_port', 8080)
        self._session_rule_version_mapper = session_rule_version_mapper

    def handle_redirection(self, datapath, loop, redirect_request):
        """
        Depending on redirection server address type install redirection rules
        """
        imsi = redirect_request.imsi
        ip_addr = redirect_request.ip_addr
        rule = redirect_request.rule
        rule_num = redirect_request.rule_num
        priority = redirect_request.priority

        # TODO IMPORTANT check that redirectd service is running, as its a
        # dynamic service its not on by default. Will save you some sanity :)

        # TODO figure out what to do with SIP_URI
        if rule.redirect.address_type == rule.redirect.SIP_URI:
            raise RedirectException("SIP_URIs redirection isn't setup")
        if rule.redirect.address_type == rule.redirect.IPv6:
            raise RedirectException("No ipv6 support, so no ipv6 redirect")

        self._save_redirect_entry(ip_addr, rule.redirect)
        self._install_redirect_flows(datapath, loop, imsi, rule, rule_num,
                                     priority)
        return

    def _install_redirect_flows(self, datapath, loop, imsi, rule, rule_num,
                                priority):
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
                                           rule_num, priority)
        elif rule.redirect.address_type == rule.redirect.IPv4:
            self._install_ipv4_bypass_flows(datapath, imsi, rule,
                                            rule_num, priority,
                                            [rule.redirect.server_address])

        self._install_dns_flows(datapath, imsi, rule, rule_num, priority)
        self._install_server_flows(datapath, imsi, rule, rule_num, priority)

    def _install_scratch_table_flows(self, datapath, imsi, rule, rule_num,
                                     priority):
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
        actions += self._load_rule_actions(parser, rule_num, imsi, rule.id)

        flows.add_resubmit_current_service_flow(
            datapath, self._scratch_tbl_num, match_http, actions,
            priority=priority, cookie=rule_num, hard_timeout=rule.hard_timeout,
            resubmit_table=self.main_tbl_num)

        match = MagmaMatch(imsi=encode_imsi(imsi))
        action = []
        flows.add_drop_flow(datapath, self._scratch_tbl_num, match, action,
                            priority=flows.MINIMUM_PRIORITY + 1,
                            cookie=rule_num)

    def _install_not_processed_flows(self, datapath, imsi, rule, rule_num,
                                     priority):
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
                           eth_type=ether_types.ETH_TYPE_IP)
        action = [of_note]
        flows.add_resubmit_current_service_flow(
            datapath, self.main_tbl_num, match, action, priority=priority,
            cookie=rule_num, hard_timeout=rule.hard_timeout,
            resubmit_table=self._scratch_tbl_num)

        match = MagmaMatch(imsi=encode_imsi(imsi),
                           direction=Direction.OUT,
                           reg0=self.REDIRECT_PROCESSED)
        action = [of_note]
        flows.add_resubmit_next_service_flow(
            datapath, self.main_tbl_num, match, action, priority=priority,
            cookie=rule_num, hard_timeout=rule.hard_timeout,
            resubmit_table=self.next_table)

    def _install_server_flows(self, datapath, imsi, rule, rule_num, priority):
        """
        Install the redirect flows to redirect all HTTP traffic to the captive
        portal and to drop all of the rest.
        """
        self._install_scratch_table_flows(datapath, imsi, rule, rule_num,
                                          priority)
        self._install_not_processed_flows(datapath, imsi, rule, rule_num,
                                          priority)

    def _install_url_bypass_flows(self, datapath, loop, imsi, rule, rule_num,
                                  priority):
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
                                            priority, cached_ips)
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
                                            priority, ips)

        asyncio.ensure_future(query, loop=loop).add_done_callback(add_flows)

    def _install_ipv4_bypass_flows(self, datapath, imsi, rule, rule_num,
                                   priority, ips):
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
        actions += self._load_rule_actions(parser, rule_num, imsi, rule.id)

        matches = []
        for ip in ips:
            matches.append(MagmaMatch(
                eth_type=ether_types.ETH_TYPE_IP, direction=Direction.OUT,
                ipv4_dst=ip, imsi=encode_imsi(imsi)
            ))
            matches.append(MagmaMatch(
                eth_type=ether_types.ETH_TYPE_IP, direction=Direction.IN,
                ipv4_src=ip, imsi=encode_imsi(imsi)
            ))
        for match in matches:
            flows.add_resubmit_next_service_flow(
                datapath, self.main_tbl_num, match, actions,
                priority=priority + 1, cookie=rule_num,
                hard_timeout=rule.hard_timeout, resubmit_table=self.next_table)

    def _install_dns_flows(self, datapath, imsi, rule, rule_num, priority):
        """
        Installs flows that allow DNS queries to path through.
        """
        parser = datapath.ofproto_parser
        of_note = parser.NXActionNote(list(rule.id.encode()))
        actions = [
            of_note,
        ]
        actions += self._load_rule_actions(parser, rule_num, imsi, rule.id)
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
                resubmit_table=self.next_table)

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

    def _load_rule_actions(self, parser, rule_num, imsi, rule_id):
        version = self._session_rule_version_mapper.get_version(imsi, rule_id)
        return [
            parser.NXActionRegLoad2(dst='reg2', value=rule_num),
            parser.NXActionRegLoad2(dst=RULE_VERSION_REG, value=version),
        ]
