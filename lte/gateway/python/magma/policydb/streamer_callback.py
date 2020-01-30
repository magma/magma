"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from typing import Any, List, Set
from lte.protos.policydb_pb2 import AssignedPolicies, PolicyRule,\
    ChargingRuleNameSet
from lte.protos.session_manager_pb2 import PolicyReAuthRequest,\
    StaticRuleInstall
from magma.common.streamer import StreamerClient
from orc8r.protos.streamer_pb2 import DataUpdate
from magma.policydb.reauth_handler import ReAuthHandler
from magma.policydb.rule_map_store import RuleAssignmentsDict
from magma.policydb.rule_store import PolicyRuleDict
from magma.policydb.basename_store import BaseNameDict


class PolicyDBStreamerCallback(StreamerClient.Callback):
    """
    Callback implementation for the PolicyDB StreamerClient instance.
    """

    def __init__(self):
        self._policy_dict = PolicyRuleDict()

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name, updates, resync):
        logging.info("Processing %d policy updates (resync=%s)",
                     len(updates), resync)
        if resync:
            policy_ids = set()
            for update in updates:
                policy = PolicyRule()
                policy.ParseFromString(update.value)
                self._store_policy_rule(policy)
                policy_ids.add(policy.id)
            logging.debug("Resync with policies: %s", ','.join(policy_ids))
            self._remove_old_policies(policy_ids)
            self._policy_dict.send_update_notification()
        else:
            pass

    def _store_policy_rule(self, policy):
        self._policy_dict[policy.id] = policy

    def _remove_old_policies(self, id_set):
        """
        Scan the set of ids passes in the streaming update to see which have
        been deleted and delete them in the policy dictionary
        """
        missing_rules = set(self._policy_dict.keys()) - id_set
        for rule in missing_rules:
            del self._policy_dict[rule]


class BaseNamesStreamerCallback(StreamerClient.Callback):
    """
    Callback for the base names streamer policy which persists the basenames
    and rules associated to the basename
    """
    def __init__(
        self,
        basenames_dict: BaseNameDict,
    ):
        self._basenames = basenames_dict

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name: str, updates: List[DataUpdate],
                       resync: bool):
        logging.info('Processing %d basename -> policy updates', len(updates))
        for update in updates:
            basename = ChargingRuleNameSet()
            basename.ParseFromString(update.value)
            self._basenames[update.key] = basename


class RuleMappingsStreamerCallback(StreamerClient.Callback):
    """
    Callback for the rule mapping streamer policy which persists the policies
    and basenames active for a subscriber.
    """
    def __init__(
        self,
        reauth_handler: ReAuthHandler,
        rules_by_basename: BaseNameDict,
        rules_by_sid: RuleAssignmentsDict,
    ):
        self._reauth_handler = reauth_handler
        self._rules_by_basename = rules_by_basename
        self._rules_by_sid = rules_by_sid

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name: str, updates: List[DataUpdate],
                       resync: bool):
        logging.info('Processing %d SID -> policy updates', len(updates))
        for update in updates:
            policies = AssignedPolicies()
            policies.ParseFromString(update.value)
            self._handle_update(update.key, policies)

        # TODO: delta with state in Redis, send RARs, persist new state

    def _handle_update(
        self,
        subscriber_id: str,
        assigned_policies: AssignedPolicies,
    ):
        """
        Based on the streamed updates, find the delta in added and removed
        rules. Then make a RAR to send to sessiond. If all goes successfully,
        update Redis with the currently installed policies for the subscriber.
        """
        prev_rules = self._get_prev_policies(subscriber_id)
        desired_rules = self._get_desired_rules(assigned_policies)

        rar = self._generate_rar(subscriber_id,
                                 list(desired_rules - prev_rules),
                                 list(prev_rules - desired_rules))
        self._reauth_handler.handle_policy_re_auth(rar)

    def _get_desired_rules(
        self,
        assigned_policies: AssignedPolicies,
    ) -> Set[str]:
        """
        Get the desired list of all rules that should be installed for the
        subscriber. This is built with a combination of base names and the
        assigned policies.
        """
        desired_rules = set(assigned_policies.assigned_policies)
        for basename in assigned_policies.assigned_base_names:
            if basename not in self._rules_by_basename:
                # They will be installed when we get the basename definition
                # streamed down from orc8r
                continue
            desired_rules.update(self._rules_by_basename[basename].RuleNames)
        return desired_rules

    def _get_prev_policies(self, subscriber_id: str) -> Set[str]:
        if subscriber_id not in self._rules_by_sid:
            return set()
        return set(self._rules_by_sid[subscriber_id].installed_policies)

    def _generate_rar(
        self,
        subscriber_id: str,
        added_rules: List[str],
        removed_rules: List[str],
    ) -> PolicyReAuthRequest:
        rules_to_install = [
            StaticRuleInstall(rule_id=rule_id) for rule_id in added_rules
        ]
        return PolicyReAuthRequest(
            # Skip the session ID, so apply to all sessions of the subscriber
            imsi=subscriber_id,
            rules_to_install=rules_to_install,
            rules_to_remove=removed_rules,
            # No changes to dynamic rules
            # No event triggers
            # No additional usage monitoring credits
            # No QoS info
        )
