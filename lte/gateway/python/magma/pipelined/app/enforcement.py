"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from typing import List

from lte.protos.pipelined_pb2 import RuleModResult
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.app.policy_mixin import PolicyMixin
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MsgChannel, MessageHub
from magma.pipelined.openflow.registers import Direction, RULE_VERSION_REG
from magma.pipelined.policy_converters import FlowMatchError, \
    flow_match_to_magma_match
from magma.pipelined.redirect import RedirectionManager, RedirectException
from magma.pipelined.qos.qos_rate_limiting import QosQueueMap
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import ether_types
from ryu.ofproto.ofproto_v1_4_parser import OFPFlowStats


class EnforcementController(PolicyMixin, MagmaController):
    """
    EnforcementController

    The enforcement controller installs flows for policy enforcement and
    classification. Each flow installed matches on a rule and an IMSI and then
    classifies the packet with the rule. The flow also redirects and drops
    the packet as specified in the policy.

    NOTE: Enforcement currently relies on the fact that policies do not
    overlap. In this implementation, there is the idea of a 'default rule'
    which is the catch-all. This rule is treated specially and tagged with a
    specific priority.
    """

    APP_NAME = "enforcement"
    APP_TYPE = ControllerType.LOGICAL
    ENFORCE_DROP_PRIORITY = flows.MINIMUM_PRIORITY + 1
    # Should not overlap with the drop flow as drop matches all packets.
    MIN_ENFORCE_PROGRAMMED_FLOW = ENFORCE_DROP_PRIORITY + 1
    MAX_ENFORCE_PRIORITY = flows.MAXIMUM_PRIORITY
    # Effectively range is 2 -> 65535
    ENFORCE_PRIORITY_RANGE = MAX_ENFORCE_PRIORITY - MIN_ENFORCE_PROGRAMMED_FLOW

    def __init__(self, *args, **kwargs):
        super(EnforcementController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_main_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self._enforcement_stats_scratch = self._service_manager.get_table_num(
            EnforcementStatsController.APP_NAME)
        self.loop = kwargs['loop']
        self._relay_enabled = kwargs['mconfig'].relay_enabled
        self._qos_map = QosQueueMap(
            kwargs['config']['nat_iface'],
            kwargs['config']['enodeb_iface'],
            kwargs['config']['enable_queue_pgm'])
        self._msg_hub = MessageHub(self.logger)
        self._redirect_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]
        self._bridge_ip_address = kwargs['config']['bridge_ip_address']
        self._redirect_manager = None
        self._clean_restart = kwargs['config']['clean_restart']
        self._relay_enabled = kwargs['mconfig'].relay_enabled
        if not self._relay_enabled:
            self.logger.info('Relay mode is not enabled, enforcement will not'
                             ' wait for sessiond to push flows.')

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath

        if not self._relay_enabled:
            self._install_default_flows_if_not_installed(datapath, [])

        self._redirect_manager = RedirectionManager(
            self._bridge_ip_address,
            self.logger,
            self.tbl_num,
            self._enforcement_stats_scratch,
            self._redirect_scratch,
            self._session_rule_version_mapper)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        if self._clean_restart:
            self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._redirect_scratch)

    def cleanup_state(self):
        pass

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    def _install_default_flows_if_not_installed(self, datapath,
            existing_flows: List[OFPFlowStats]) -> List[OFPFlowStats]:
        """
        For each direction set the default flows to just forward to next app.
        The enforcement flows for each subscriber would be added when the
        IP session is created, by reaching out to the controller/PCRF.
        If default flows are already installed, do nothing.

        Args:
            datapath: ryu datapath struct
        Returns:
            The list of flows that remain after inserting default flows
        """
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)

        inbound_msg = flows.get_add_resubmit_next_service_flow_msg(
            datapath, self.tbl_num, inbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

        outbound_msg = flows.get_add_resubmit_next_service_flow_msg(
            datapath, self.tbl_num, outbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

        msgs, remaining_flows = self._msg_hub \
            .filter_msgs_if_not_in_flow_list([inbound_msg, outbound_msg],
                                             existing_flows)
        if msgs:
            chan = self._msg_hub.send(msgs, datapath)
            self._wait_for_responses(chan, len(msgs))

        return remaining_flows

    def get_of_priority(self, precedence):
        """
        Lower the precedence higher the importance of the flow in 3GPP.
        Higher the priority higher the importance of the flow in openflow.
        Convert precedence to priority:
        1 - Flows with precedence > 65534 will have min priority which is the
        min priority for a programmed flow = (default drop + 1)
        2 - Flows in the precedence range 0-65534 will have priority 65535 -
        Precedence
        :param precedence:
        :return:
        """
        if precedence >= self.ENFORCE_PRIORITY_RANGE:
            self.logger.warning(
                "Flow precedence is higher than OF range using min priority %d",
                self.MIN_ENFORCE_PROGRAMMED_FLOW)
            return self.MIN_ENFORCE_PROGRAMMED_FLOW
        return self.MAX_ENFORCE_PRIORITY - precedence

    def _get_rule_match_flow_msgs(self, imsi, rule):
        """
        Get a flow msg to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        priority = self.get_of_priority(rule.priority)
        ul_qos = rule.qos.max_req_bw_ul
        dl_qos = rule.qos.max_req_bw_dl

        flow_adds = []
        for flow in rule.flow_list:
            try:
                flow_adds.append(self._get_classify_rule_flow_msg(
                    imsi, flow, rule_num, priority, ul_qos,
                    dl_qos, rule.hard_timeout,
                    rule.id))

            except FlowMatchError as err:  # invalid match
                self.logger.error(
                    "Failed to get flow msg '%s' for subscriber %s: %s",
                    rule.id, imsi, err)
                raise err
        return flow_adds

    def _install_flow_for_rule(self, imsi, ip_addr, rule):
        """
        Install a flow to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """

        if rule.redirect.support == rule.redirect.ENABLED:
            return self._install_redirect_flow(imsi, ip_addr, rule)

        if not rule.flow_list:
            self.logger.error('The flow list for imsi %s, rule.id - %s'
                              'is empty, this shoudn\'t happen', imsi, rule.id)
            return RuleModResult.FAILURE

        flow_adds = []
        try:
            flow_adds = self._get_rule_match_flow_msgs(imsi, rule)
        except FlowMatchError:
            return RuleModResult.FAILURE

        chan = self._msg_hub.send(flow_adds, self._datapath)

        return self._wait_for_rule_responses(imsi, rule, chan)

    def _wait_for_rule_responses(self, imsi, rule, chan):
        def fail(err):
            self.logger.error(
                "Failed to install rule %s for subscriber %s: %s",
                rule.id, imsi, err)
            self._deactivate_flow_for_rule(imsi, rule.id)
            return RuleModResult.FAILURE

        for _ in range(len(rule.flow_list)):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS")
            if not result.ok():
                return fail(result.exception())
        return RuleModResult.SUCCESS

    def _get_classify_rule_flow_msg(self, imsi, flow, rule_num, priority,
                                    ul_qos, dl_qos, hard_timeout, rule_id):
        """
        Install a flow from a rule. If the flow action is DENY, then the flow
        will drop the packet. Otherwise, the flow classifies the packet with
        its matched rule and injects the rule num into the packet's register.
        """
        flow_match = flow_match_to_magma_match(flow.match)
        flow_match.imsi = encode_imsi(imsi)
        flow_match_actions = self._get_classify_rule_of_actions(
            flow, rule_num, imsi, ul_qos, dl_qos, rule_id)
        if flow.action == flow.DENY:
            return flows.get_add_drop_flow_msg(self._datapath,
                                               self.tbl_num,
                                               flow_match,
                                               flow_match_actions,
                                               hard_timeout=hard_timeout,
                                               priority=priority,
                                               cookie=rule_num)

        if self._enforcement_stats_scratch:
            return flows.get_add_resubmit_current_service_flow_msg(
                self._datapath,
                self.tbl_num,
                flow_match,
                flow_match_actions,
                hard_timeout=hard_timeout,
                priority=priority,
                cookie=rule_num,
                resubmit_table=self._enforcement_stats_scratch)

        # If enforcement stats has not claimed a scratch table, resubmit
        # directly to the next app.
        return flows.get_add_resubmit_next_service_flow_msg(
            self._datapath,
            self.tbl_num,
            flow_match,
            flow_match_actions,
            hard_timeout=hard_timeout,
            priority=priority,
            cookie=rule_num,
            resubmit_table=self.next_main_table)

    def _install_redirect_flow(self, imsi, ip_addr, rule):
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        priority = self.get_of_priority(rule.priority)
        redirect_request = RedirectionManager.RedirectRequest(
            imsi=imsi,
            ip_addr=ip_addr,
            rule=rule,
            rule_num=rule_num,
            priority=priority)
        try:
            self._redirect_manager.handle_redirection(
                self._datapath, self.loop, redirect_request)
            return RuleModResult.SUCCESS
        except RedirectException as err:
            self.logger.error(
                'Redirect Exception for imsi %s, rule.id - %s : %s',
                imsi, rule.id, err
            )
            return RuleModResult.FAILURE

    def _get_classify_rule_of_actions(self, flow, rule_num, imsi, ul_qos,
                                      dl_qos, rule_id):
        parser = self._datapath.ofproto_parser
        # encode the rule id in hex
        of_note = parser.NXActionNote(list(rule_id.encode()))
        actions = [of_note]
        if flow.action == flow.DENY:
            return actions

        # QoS Rate-Limiting is currently supported for uplink traffic
        qid = 0
        if ul_qos != 0 and flow.match.direction == flow.match.UPLINK:
            qid = self._qos_map.map_flow_to_queue(imsi, rule_num, ul_qos, True)
        elif dl_qos != 0 and flow.match.direction == flow.match.DOWNLINK:
            qid = self._qos_map.map_flow_to_queue(imsi, rule_num, dl_qos, False)

        if qid != 0:
            actions.append(parser.OFPActionSetField(pkt_mark=qid))

        version = self._session_rule_version_mapper.get_version(imsi, rule_id)
        actions.extend(
            [parser.NXActionRegLoad2(dst='reg2', value=rule_num),
             parser.NXActionRegLoad2(dst=RULE_VERSION_REG, value=version)
             ])

        return actions

    def _get_default_flow_msg_for_subscriber(self, imsi):
        match = MagmaMatch(imsi=encode_imsi(imsi))
        actions = []
        return flows.get_add_drop_flow_msg(self._datapath, self.tbl_num,
            match, actions, priority=self.ENFORCE_DROP_PRIORITY)

    def _install_default_flow_for_subscriber(self, imsi):
        """
        Add a low priority flow to drop a subscriber's traffic in the event
        that all rules have been deactivated.

        Args:
            imsi (string): subscriber id
        """
        match = MagmaMatch(imsi=encode_imsi(imsi))
        actions = []  # empty options == drop
        flows.add_drop_flow(self._datapath, self.tbl_num, match, actions,
                            priority=self.ENFORCE_DROP_PRIORITY)

    def _deactivate_flow_for_rule(self, imsi, rule_id):
        """
        Deactivate a specific rule using the flow cookie for a subscriber
        """
        try:
            num = self._rule_mapper.get_rule_num(rule_id)
        except KeyError:
            self.logger.error('Could not find rule id %s', rule_id)
            return
        cookie, mask = (num, flows.OVS_COOKIE_MATCH_ALL)
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.tbl_num, match,
                          cookie=cookie, cookie_mask=mask)
        self._redirect_manager.deactivate_flow_for_rule(self._datapath, imsi,
                                                        num)
        self._qos_map.del_queue_for_flow(imsi, num)

    def _deactivate_flows_for_subscriber(self, imsi):
        """ Deactivate all rules for a subscriber, ending any enforcement """
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.tbl_num, match)
        self._redirect_manager.deactivate_flows_for_subscriber(self._datapath,
                                                               imsi)
        self._qos_map.del_subscriber_queues(imsi)

    def deactivate_rules(self, imsi, rule_ids):
        """
        Deactivate flows for a subscriber. If only imsi is present, delete all
        rule flows for a subscriber (i.e. end its session). If rule_ids are
        present, delete the rule flows for that subscriber.

        Args:
            imsi (string): subscriber id
            rule_ids (list of strings): policy rule ids
        """
        if not self.init_finished:
            self.logger.error('Pipelined is not initialized')
            return RuleModResult.FAILURE

        if self._datapath is None:
            self.logger.error('Datapath not initialized')
            return

        if not imsi:
            self.logger.error('No subscriber specified')
            return

        if not rule_ids:
            self._deactivate_flows_for_subscriber(imsi)
        else:
            for rule_id in rule_ids:
                self._deactivate_flow_for_rule(imsi, rule_id)
