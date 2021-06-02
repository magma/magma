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

import logging
from typing import Any, List, Set

import grpc
from lte.protos.policydb_pb2 import (
    ApnPolicySet,
    ChargingRuleNameSet,
    PolicyRule,
    RatingGroup,
    SubscriberPolicySet,
)
from lte.protos.session_manager_pb2 import (
    DynamicRuleInstall,
    RuleSet,
    RulesPerSubscriber,
    SessionRules,
    StaticRuleInstall,
)
from lte.protos.session_manager_pb2_grpc import LocalSessionManagerStub
from magma.common.streamer import StreamerClient
from magma.policydb.apn_rule_map_store import ApnRuleAssignmentsDict
from magma.policydb.basename_store import BaseNameDict
from magma.policydb.default_rules import get_allow_all_policy_rule
from magma.policydb.rating_group_store import RatingGroupsDict
from magma.policydb.rule_store import PolicyRuleDict
from orc8r.protos.streamer_pb2 import DataUpdate


class PolicyDBStreamerCallback(StreamerClient.Callback):
    """
    Callback implementation for the PolicyDB StreamerClient instance.
    """

    def __init__(self):
        self._policy_dict = PolicyRuleDict()

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name, updates, resync):
        logging.info(
            "Processing %d policy updates (resync=%s)",
            len(updates), resync,
        )
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

    def process_update(
        self, stream_name: str, updates: List[DataUpdate],
        resync: bool,
    ):
        logging.info('Processing %d basename -> policy updates', len(updates))
        for update in updates:
            basename = ChargingRuleNameSet()
            basename.ParseFromString(update.value)
            self._basenames[update.key] = basename


class ApnRuleMappingsStreamerCallback(StreamerClient.Callback):
    """
    Callback for the apn rule mappings streamer policy which persists
    the mapping of (imsi, subscriber) tuples -> rules
    """

    def __init__(
        self,
        session_mgr_stub: LocalSessionManagerStub,
        rules_by_basename: BaseNameDict,
        apn_rules_by_sid: ApnRuleAssignmentsDict,
    ):
        self._session_mgr_stub = session_mgr_stub
        self._rules_by_basename = rules_by_basename
        self._apn_rules_by_sid = apn_rules_by_sid

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(
        self,
        stream_name: str,
        updates: List[DataUpdate],
        resync: bool,
    ):
        logging.info(
            'Processing %d SID -> apn -> policy updates',
            len(updates),
        )
        all_subscriber_rules = []  # type: List[RulesPerSubscriber]
        found_update = False
        for update in updates:
            imsi = update.key
            subApnPolicies = SubscriberPolicySet()
            subApnPolicies.ParseFromString(update.value)
            is_updated = self._are_sub_policies_updated(imsi, subApnPolicies)
            if is_updated:
                found_update = True
                all_subscriber_rules.append(
                    self._build_sub_rule_set(imsi, subApnPolicies),
                )
                self._apn_rules_by_sid[imsi] = subApnPolicies
        if not found_update:
            logging.debug(
                "No IMSIs with APN->Policy assignments found. "
                "Not sending an update to SessionD",
            )
            return
        logging.info(
            'Updating %d IMSIs with new APN->policy assignments',
            len(all_subscriber_rules),
        )
        update = SessionRules(rules_per_subscriber=all_subscriber_rules)
        try:
            self._session_mgr_stub.SetSessionRules(update, timeout=5)
        except grpc.RpcError as e:
            logging.error('Unable to apply apn->policy updates %s', str(e))

    def _are_sub_policies_updated(
        self,
        subscriber_id: str,
        subApnPolicies: SubscriberPolicySet,
    ) -> bool:
        if subscriber_id not in self._apn_rules_by_sid:
            return True
        prev = self._apn_rules_by_sid[subscriber_id]
        # TODO: (8/21/2020) repeated fields may not be ordered the same, use a
        #       different method to compare later
        return subApnPolicies.SerializeToString() != prev.SerializeToString()

    def _build_sub_rule_set(
        self,
        subscriber_id: str,
        sub_apn_policies: SubscriberPolicySet,
    ) -> RulesPerSubscriber:
        apn_rule_sets = []  # type: List[RuleSet]
        global_rules = self._get_global_static_rules(sub_apn_policies)

        for apn_policy_set in sub_apn_policies.rules_per_apn:
            # Static rule installs
            static_rule_ids = self._get_desired_static_rules(apn_policy_set)
            static_rules = []  # type: List[StaticRuleInstall]
            for rule_id in static_rule_ids:
                static_rules.append(StaticRuleInstall(rule_id=rule_id))
            # Add global rules
            for rule_id in global_rules:
                static_rules.append(StaticRuleInstall(rule_id=rule_id))

            # Dynamic rule installs
            dynamic_rules = []  # type: List[DynamicRuleInstall]
            # Build the rule id to be globally unique
            rule = DynamicRuleInstall(
                policy_rule=get_allow_all_policy_rule(
                    subscriber_id,
                    apn_policy_set.apn,
                ),
            )
            dynamic_rules.append(rule)

            # Build the APN rule set
            apn_rule_sets.append(
                RuleSet(
                    apply_subscriber_wide=False,
                    apn=apn_policy_set.apn,
                    static_rules=static_rules,
                    dynamic_rules=dynamic_rules,
                ),
            )

        return RulesPerSubscriber(
            imsi=subscriber_id,
            rule_set=apn_rule_sets,
        )

    def _get_global_static_rules(
        self,
        sub_apn_policies: SubscriberPolicySet,
    ) -> Set[str]:
        global_rules = set(sub_apn_policies.global_policies)
        for basename in sub_apn_policies.global_base_names:
            if basename not in self._rules_by_basename:
                # Eventually, basename definition will be streamed from orc8r
                continue
            global_rules.update(
                self._rules_by_basename[basename].RuleNames,
            )
        return global_rules

    def _get_desired_static_rules(
        self,
        policies: ApnPolicySet,
    ) -> Set[str]:
        desired_rules = set(policies.assigned_policies)
        for basename in policies.assigned_base_names:
            if basename not in self._rules_by_basename:
                # Eventually, basename definition will be streamed from orc8r
                continue
            desired_rules.update(self._rules_by_basename[basename].RuleNames)
        return desired_rules


class RatingGroupsStreamerCallback(StreamerClient.Callback):
    """
    Callback for the rating groups streamer which persists the rating groups
    """

    def __init__(
            self,
            rating_groups_dict: RatingGroupsDict,
    ):
        self._rating_groups = rating_groups_dict

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(
        self, stream_name: str, updates: List[DataUpdate],
        resync: bool,
    ):
        logging.info('Processing %d rating group updates', len(updates))
        for update in updates:
            rg = RatingGroup()
            rg.ParseFromString(update.value)
            self._rating_groups[update.key] = rg
