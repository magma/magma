"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from typing import List
from abc import ABCMeta, abstractmethod

from ryu.ofproto.ofproto_v1_4_parser import OFPFlowStats

from lte.protos.pipelined_pb2 import RuleModResult, SetupFlowsResult, \
    ActivateFlowsResult, ActivateFlowsRequest, SubscriberQuotaUpdate
from magma.pipelined.openflow import flows
from magma.policydb.rule_store import PolicyRuleDict
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, IMSI_REG, \
    DIRECTION_REG
from magma.pipelined.openflow.messages import MsgChannel
from magma.pipelined.policy_converters import FlowMatchError


class PolicyMixin(metaclass=ABCMeta):
    """
    PolicyMixin

    Mixin class for policy enforcement apps that includes common methods
    used for rule activation/deactivation.
    """

    def __init__(self, *args, **kwargs):
        super(PolicyMixin, self).__init__(*args, **kwargs)
        self._datapath = None
        self._policy_dict = PolicyRuleDict()
        self._rule_mapper = kwargs['rule_id_mapper']
        self._session_rule_version_mapper = kwargs[
            'session_rule_version_mapper']
        self._relay_enabled = kwargs['mconfig'].relay_enabled
        if not self._relay_enabled:
            self.logger.info('Relay mode is not enabled, init finished')
            self.init_finished = True

    # pylint:disable=unused-argument
    def setup(self, requests: List[ActivateFlowsRequest],
              quota_updates: List[SubscriberQuotaUpdate],
              startup_flows: List[OFPFlowStats]) -> SetupFlowsResult:
        """
        Setup flows for subscribers, used on restart.

        Args:
            rules (List[OFPFlowStats]): list of subcriber policyrules
        """
        self.logger.debug('Setting up enforcer default rules')
        remaining_flows = self._install_default_flows_if_not_installed(
            self._datapath, startup_flows)

        self.logger.debug('Startup flows before filtering -> %s',
            [flow.match for flow in startup_flows])
        extra_flows = self._add_missing_flows(requests, remaining_flows)

        self.logger.debug('Startup flows after filtering will be deleted -> %s',
            [flow.match for flow in startup_flows])
        self._remove_extra_flows(extra_flows)

        # For now just reinsert redirection rules, this is a bit of a hack but
        # redirection relies on async dns request to be setup and we can't
        # currently do this from out synchronous setup request. So just reinsert
        self._process_redirection_rules(requests)

        self.init_finished = True
        return SetupFlowsResult.SUCCESS

    def _remove_extra_flows(self, extra_flows):
        msg_list = []
        for flow in extra_flows:
            if DIRECTION_REG in flow.match:
                direction = Direction(flow.match.get(DIRECTION_REG, None))
            else:
                direction = None
            match = MagmaMatch(imsi=flow.match.get(IMSI_REG, None),
                direction=direction)
            self.logger.debug('Sending msg for deletion -> %s',
                flow.match.get('reg1', None))
            msg_list.append(flows.get_delete_flow_msg(
                self._datapath, self.tbl_num, match, cookie=flow.cookie,
                cookie_mask=flows.OVS_COOKIE_MATCH_ALL))
        if msg_list:
            chan = self._msg_hub.send(msg_list, self._datapath)
            self._wait_for_responses(chan, len(msg_list))

    def _add_missing_flows(self, requests, current_flows):
        msg_list = []
        for add_flow_req in requests:
            imsi = add_flow_req.sid.id
            static_rule_ids = add_flow_req.rule_ids
            dynamic_rules = add_flow_req.dynamic_rules

            for rule_id in static_rule_ids:
                rule = self._policy_dict[rule_id]
                if rule is None:
                    self.logger.error("Could not find rule for rule_id: %s",
                        rule_id)
                    continue
                try:
                    if rule.redirect.support == rule.redirect.ENABLED:
                        continue
                    flow_adds = self._get_rule_match_flow_msgs(imsi, rule)
                    msg_list.extend(flow_adds)
                except FlowMatchError:
                    self.logger.error("Failed to verify rule_id: %s", rule_id)

            for rule in dynamic_rules:
                try:
                    if rule.redirect.support == rule.redirect.ENABLED:
                        continue
                    flow_adds = self._get_rule_match_flow_msgs(imsi, rule)
                    msg_list.extend(flow_adds)
                except FlowMatchError:
                    self.logger.error("Failed to verify rule_id: %s", rule.id)

            flow_add = self._get_default_flow_msg_for_subscriber(imsi)
            if flow_add:
                msg_list.append(flow_add)

        msgs_to_send, remaining_flows = \
            self._msg_hub.filter_msgs_if_not_in_flow_list(msg_list,
                                                          current_flows)
        if msgs_to_send:
            chan = self._msg_hub.send(msgs_to_send, self._datapath)
            self._wait_for_responses(chan, len(msgs_to_send))

        return remaining_flows

    def _process_redirection_rules(self, requests):
        for add_flow_req in requests:
            imsi = add_flow_req.sid.id
            ip_addr = add_flow_req.ip_addr
            static_rule_ids = add_flow_req.rule_ids
            dynamic_rules = add_flow_req.dynamic_rules
            for rule_id in static_rule_ids:
                rule = self._policy_dict[rule_id]
                if rule is None:
                    self.logger.error("Could not find rule for rule_id: %s",
                        rule_id)
                    continue
                if rule.redirect.support == rule.redirect.ENABLED:
                    self._install_redirect_flow(imsi, ip_addr, rule)

            for rule in dynamic_rules:
                if rule.redirect.support == rule.redirect.ENABLED:
                    self._install_redirect_flow(imsi, ip_addr, rule)

    def activate_rules(self, imsi, ip_addr, static_rule_ids, dynamic_rules):
        """
        Activate the flows for a subscriber based on the rules stored in Redis.
        During activation, a default flow may be installed for the subscriber.

        Args:
            imsi (string): subscriber id
            ip_addr (string): subscriber session ipv4 address
            static_rule_ids (string []): list of static rules to activate
            dynamic_rules (PolicyRule []): list of dynamic rules to activate
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized for adding flows')
            return ActivateFlowsResult(
                static_rule_results=[RuleModResult(
                    rule_id=rule_id,
                    result=RuleModResult.FAILURE,
                ) for rule_id in static_rule_ids],
                dynamic_rule_results=[RuleModResult(
                    rule_id=rule.id,
                    result=RuleModResult.FAILURE,
                ) for rule in dynamic_rules],
            )
        static_results = []
        for rule_id in static_rule_ids:
            res = self._install_flow_for_static_rule(imsi, ip_addr, rule_id)
            static_results.append(RuleModResult(rule_id=rule_id, result=res))
        dyn_results = []
        for rule in dynamic_rules:
            res = self._install_flow_for_rule(imsi, ip_addr, rule)
            dyn_results.append(RuleModResult(rule_id=rule.id, result=res))

        # Install a base flow for when no rule is matched.
        self._install_default_flow_for_subscriber(imsi)
        return ActivateFlowsResult(
            static_rule_results=static_results,
            dynamic_rule_results=dyn_results,
        )

    def _install_flow_for_static_rule(self, imsi, ip_addr, rule_id):
        """
        Install a flow to get stats for a particular static rule id. The rule
        will be loaded from Redis and installed.

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

    def _wait_for_responses(self, chan, response_count):
        def fail(err):
            #TODO need to rework setup to return all rule specific success/fails
            self.logger.error("Failed to install rule for subscriber: %s", err)

        for _ in range(response_count):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS")
            if not result.ok():
                return fail(result.exception())

    @abstractmethod
    def _install_flow_for_rule(self, imsi, ip_addr, rule):
        """
        Install a flow given a rule. Subclass should implement this.

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        raise NotImplementedError

    @abstractmethod
    def _install_default_flow_for_subscriber(self, imsi):
        """
        Install a flow for the subscriber in the event no rule is matched.
        Subclass should implement this.

        Args:
            imsi (string): subscriber id
        """
        raise NotImplementedError

    @abstractmethod
    def _install_redirect_flow(self, imsi, ip_addr, rule):
        """
        Install a redirection flow for the subscriber.
        Subclass should implement this.

        Args:
            imsi (string): subscriber id
            ip_addr (string): subscriber ip
            rule (PolicyRule): policyrule to install
        """
        raise NotImplementedError

    @abstractmethod
    def _install_default_flows_if_not_installed(self, datapath,
            existing_flows: List[OFPFlowStats]) -> List[OFPFlowStats]:
        """
        Install default flows(if not already installed), if no other flows are
        matched.

        Returns:
            The list of flows that remain after inserting default flows
        """
        raise NotImplementedError
