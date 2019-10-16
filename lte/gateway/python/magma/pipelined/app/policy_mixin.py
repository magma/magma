"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from abc import ABCMeta, abstractmethod

from lte.protos.pipelined_pb2 import RuleModResult, ActivateFlowsResult
from magma.policydb.rule_store import PolicyRuleDict


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
            self.logger.info(
                "In PolicyMixin result of trying to install rule for imsi:%s, rule_id:%s res:%s",
                imsi, rule.id, res)
            dyn_results.append(RuleModResult(rule_id=rule.id, result=res))

        # Install a base flow for when no rule is matched.
        self._install_default_flow_for_subscriber(imsi)
        res = ActivateFlowsResult(
            static_rule_results=static_results,
            dynamic_rule_results=dyn_results,
        )
        self.logger.info("ActivateFlowsResult : %s", res)
        return res

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
