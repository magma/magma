"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from lte.protos.pipelined_pb2 import ActivateFlowsResult, RuleModResult
from magma.pipelined.app.base import MagmaController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MessageHub, MsgChannel
from magma.pipelined.openflow.registers import Direction
from magma.pipelined.policy_converters import FlowMatchError, \
    flow_match_to_magma_match
from magma.pipelined.qos.qos_rate_limiting import QosQueueMap
from magma.pipelined.redirect import RedirectException, RedirectionManager
from magma.policydb.rule_store import PolicyRuleDict
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import ether_types


class EnforcementController(MagmaController):
    """
    EnforcementController

    The enforcement controller installs flows for tracking subscriber usage
    per rule and enforcing the usage. Each flow installed matches on a rule
    and an IMSI and then the statistics are sent to sessiond for tracking.

    NOTE: Enforcement currently relies on the fact that policies do not
    overlap. In this implementation, there is the idea of a 'default rule'
    which is the catch-all. This rule is treated specially and tagged with a
    specific priority.
    """

    APP_NAME = "enforcement"
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
        self._datapath = None
        self._rule_mapper = kwargs['rule_id_mapper']
        self.loop = kwargs['loop']
        self._policy_dict = PolicyRuleDict()
        self._qos_map = QosQueueMap(
            kwargs['config']['nat_iface'],
            kwargs['config']['enodeb_iface'],
            kwargs['config']['enable_queue_pgm'])
        self._msg_hub = MessageHub(self.logger)
        self._redirect_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]

        self._redirect_manager = RedirectionManager(
            kwargs['config']['bridge_ip_address'],
            self.logger,
            self.tbl_num,
            self.next_main_table,
            self._redirect_scratch)

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._delete_all_flows(datapath)
        self._install_default_flows(datapath)
        self._datapath = datapath

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self._delete_all_flows(datapath)

    def _delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._redirect_scratch)

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    def _install_default_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next app.
        The enforcement flows for each subscriber would be added when the
        IP session is created, by reaching out to the controller/PCRF.

        Args:
            datapath: ryu datapath struct
        """
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, inbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, outbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

    def _install_flow_for_static_rule(self, imsi, ip_addr, rule_id):
        """
        Install a flow to get stats for a particular static rule id. The rule
        will be loaded from Redis and installed

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule_id (string): policy rule id
        """
        rule = self._policy_dict[rule_id]
        if rule is None:
            self.logger.error("Could not find rule for rule_id: %s", rule_id)
            return RuleModResult.FAILURE
        return self._install_flow_for_rule(imsi, ip_addr, rule)

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
                "Flow precedence is higher than OF range using "
                             "min priority %d",
                self.MIN_ENFORCE_PROGRAMMED_FLOW)
            return self.MIN_ENFORCE_PROGRAMMED_FLOW
        return self.MAX_ENFORCE_PRIORITY - precedence

    def _install_flow_for_rule(self, imsi, ip_addr, rule):
        """
        Install a flow to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        priority = self.get_of_priority(rule.priority)
        ul_qos = rule.qos.max_req_bw_ul

        if rule.redirect.support == rule.redirect.ENABLED:
            # TODO currently if redirection is enabled we ignore other flows
            # from rule.flow_list, confirm that this is the expected behaviour
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

        flow_adds = []
        for flow in rule.flow_list:
            try:
                flow_adds.append(self._get_add_flow_msg(
                    imsi, flow, rule_num, priority, ul_qos, rule.hard_timeout,
                    rule.id))
            except FlowMatchError as err:  # invalid match
                self.logger.error(
                    "Failed to install rule %s for subscriber %s: %s",
                    rule.id, imsi, err)
                return RuleModResult.FAILURE
        chan = self._msg_hub.send(flow_adds, self._datapath)

        return self._wait_for_responses(imsi, rule, chan)

    def _wait_for_responses(self, imsi, rule, chan):
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

    def _get_add_flow_msg(self, imsi, flow, rule_num, priority,
                          ul_qos, hard_timeout, rule_id):
        match = flow.match
        ryu_match = flow_match_to_magma_match(match)
        ryu_match.imsi = encode_imsi(imsi)

        actions = self._get_of_actions_for_flow(flow, rule_num, imsi, ul_qos,
                                                rule_id)

        if flow.action == flow.DENY:
            return flows.get_add_flow_msg(
                self._datapath, self.tbl_num, ryu_match, actions,
                hard_timeout=hard_timeout, priority=priority, cookie=rule_num)

        return flows.get_add_resubmit_next_service_flow_msg(
            self._datapath, self.tbl_num, ryu_match, actions,
            hard_timeout=hard_timeout, priority=priority, cookie=rule_num,
            resubmit_table=self.next_main_table)

    def _get_of_actions_for_flow(self, flow, rule_num, imsi, ul_qos,
                                 rule_id):
        parser = self._datapath.ofproto_parser
        # encode the rule id in hex
        of_note = parser.NXActionNote(list(rule_id.encode()))
        if flow.action == flow.DENY:
            return [of_note]

        # QoS Rate-Limiting is currently supported for uplink traffic
        if ul_qos != 0 and flow.match.direction == flow.match.UPLINK:
            qid = self._qos_map.map_flow_to_queue(imsi, rule_num, ul_qos, True)
            if qid != 0:
                return [parser.OFPActionSetField(pkt_mark=qid),
                        parser.NXActionRegLoad2(dst='reg2', value=rule_num),
                        of_note]

        return [parser.NXActionRegLoad2(dst='reg2', value=rule_num),
                of_note]

    def _install_drop_flow(self, imsi):
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

    def activate_flows(self, imsi, ip_addr, static_rule_ids,
                       dynamic_rules, fut):
        """
        Activate the flows for a subscriber based on the rules stored in Redis.
        During activation, another low priority flow is installed for the
        subscriber in the event that all rules are out of credit.

        Args:
            imsi (string): subscriber id
            ip_addr (string): subscriber session ipv4 address
            static_rule_ids (string []): list of static rules to activate
            dynamic_rules (PolicyRule []): list of dynamic rules to activate
            fut (Future): future to wait on the results of flow activations
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized for adding flows')
            fut.set_result(ActivateFlowsResult(
                static_rule_results=[RuleModResult(
                    rule_id=rule_id,
                    result=RuleModResult.FAILURE,
                ) for rule_id in static_rule_ids],
                dynamic_rule_results=[RuleModResult(
                    rule_id=rule.id,
                    result=RuleModResult.FAILURE,
                ) for rule in dynamic_rules],
            ))
            return
        static_results = []
        for rule_id in static_rule_ids:
            res = self._install_flow_for_static_rule(imsi, ip_addr, rule_id)
            static_results.append(RuleModResult(rule_id=rule_id, result=res))
        dyn_results = []
        for rule in dynamic_rules:
            res = self._install_flow_for_rule(imsi, ip_addr, rule)
            dyn_results.append(RuleModResult(rule_id=rule.id, result=res))

        # No matter what, install base flow to drop packets when all other
        # flows have been deactivated
        self._install_drop_flow(imsi)
        fut.set_result(ActivateFlowsResult(
            static_rule_results=static_results,
            dynamic_rule_results=dyn_results,
        ))

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
        flows.delete_flow(self._datapath, self.tbl_num,
                          MagmaMatch(imsi=encode_imsi(imsi)))
        self._redirect_manager.deactivate_flows_for_subscriber(self._datapath,
                                                               imsi)
        self._qos_map.del_subscriber_queues(imsi)

    def deactivate_flows(self, imsi, rule_ids):
        """
        Deactivate flows for a subscriber. If only imsi is present, delete all
        rule flows for a subscriber (i.e. end its session). If rule_ids are
        present, delete the rule flows for that subscriber.

        Args:
            imsi (string): subscriber id
            rule_ids (list of strings): policy rule ids
        """
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
