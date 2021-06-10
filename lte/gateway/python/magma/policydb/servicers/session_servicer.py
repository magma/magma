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
from typing import List, Set

from lte.protos.mconfig import mconfigs_pb2
from lte.protos.policydb_pb2 import (
    ApnPolicySet,
    RatingGroup,
    SubscriberPolicySet,
)
from lte.protos.session_manager_pb2 import (
    CreateSessionRequest,
    CreateSessionResponse,
    CreditLimitType,
    CreditUpdateResponse,
    DynamicRuleInstall,
    SessionTerminateResponse,
    StaticRuleInstall,
    UpdateSessionRequest,
    UpdateSessionResponse,
)
from lte.protos.session_manager_pb2_grpc import (
    CentralSessionControllerServicer,
    add_CentralSessionControllerServicer_to_server,
)
from magma.policydb.apn_rule_map_store import ApnRuleAssignmentsDict
from magma.policydb.basename_store import BaseNameDict
from magma.policydb.default_rules import get_allow_all_policy_rule
from magma.policydb.rating_group_store import RatingGroupsDict
from orc8r.protos.common_pb2 import NetworkID


class SessionRpcServicer(CentralSessionControllerServicer):
    """
    gRPC based server for CentralSessionController service.

    This will act as a bare-bones local PCRF and OCS.

    For all connecting subscribers, an allow-all flow will be installed as
    a dynamic policy rule. In addition, whatever static rules, and static rules
    from base names have been configured via the orc8r will also be installed
    for the subscriber.

    This limited PCRF/OCS is also used for enabling the Captive Portal
    feature.
    """

    def __init__(
        self,
        mconfig: mconfigs_pb2.PolicyDB,
        rating_groups_by_id: RatingGroupsDict,
        rules_by_basename: BaseNameDict,
        apn_rules_by_sid: ApnRuleAssignmentsDict,
    ):
        self._mconfig = mconfig
        self._network_id = NetworkID(id="_")
        self._rating_groups_by_id = rating_groups_by_id
        self._rules_by_basename = rules_by_basename
        self._apn_rules_by_sid = apn_rules_by_sid

    def get_infinite_credit_charging_keys(self) -> List[int]:
        keys = []
        for rating_group in self._rating_groups_by_id.values():
            if rating_group.limit_type == RatingGroup.INFINITE_UNMETERED:
                keys.append(rating_group.id)
        return keys

    def _get_postpay_charging_keys(self) -> List[int]:
        keys = []
        for rating_group in self._rating_groups_by_id.values():
            if rating_group.limit_type == RatingGroup.INFINITE_METERED:
                keys.append(rating_group.id)
        return keys

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_CentralSessionControllerServicer_to_server(
            self, server,
        )

    def CreateSession(
        self,
        request: CreateSessionRequest,
        context,
    ) -> CreateSessionResponse:
        """
        Handles create session request from MME by installing the necessary
        flows in pipelined's enforcement app.

        NOTE: leave the 'IMSI' prefix
        """
        imsi = request.common_context.sid.id
        apn = request.common_context.apn
        logging.info('Creating a session for subscriber ID: %s', imsi)
        return CreateSessionResponse(
            credits=self._get_credits(imsi),
            static_rules=self._get_session_static_rules(imsi, apn),
            dynamic_rules=self._get_default_dynamic_rules(imsi, apn),
            session_id=request.session_id,
        )

    def UpdateSession(
        self,
        request: UpdateSessionRequest,
        context,
    ) -> UpdateSessionResponse:
        """
        On UpdateSession, return an arbitrarily large amount of additional
        credit for the session.

        NOTE: This really shouldn't be called, as no credit should have been
        granted on CreateSession.
        """
        logging.info('UpdateSession called')
        resp = UpdateSessionResponse()
        for credit_usage_update in request.updates:
            resp.responses.extend(
                self._get_credits(credit_usage_update.common_context.sid.id),
            )
        return resp

    def TerminateSession(
        self,
        request: SessionTerminateResponse,
        context,
    ) -> SessionTerminateResponse:
        logging.info('Terminating session: %s', request.session_id)
        return SessionTerminateResponse(
            sid=request.common_context.sid.id,
            session_id=request.session_id,
        )

    def _get_default_dynamic_rules(
        self,
        subscriber_id: str,
        apn: str,
    ) -> List[DynamicRuleInstall]:
        """
        Get a list of dynamic rules to install
        Currently only includes a single rule for allow-all of traffic
        """
        return [
            DynamicRuleInstall(
                policy_rule=get_allow_all_policy_rule(subscriber_id, apn),
            ),
        ]

    def _get_session_static_rules(
        self,
        imsi: str,
        apn: str,
    ) -> List[StaticRuleInstall]:
        """
        Get the list of static rules to be installed for a subscriber
        NOTE: Remove "IMSI" prefix from imsi argument.
        """
        if imsi not in self._apn_rules_by_sid:
            return []

        sub_apn_policies = self._apn_rules_by_sid[imsi]
        assigned_static_rules = []  # type: List[StaticRuleInstall]
        # Add global rules
        global_rules = self._get_global_static_rules(sub_apn_policies)
        assigned_static_rules += \
            list(
                map(
                    lambda id: StaticRuleInstall(rule_id=id),
                    global_rules,
                ),
            )
        # Add APN specific rules
        for apn_policy_set in sub_apn_policies.rules_per_apn:
            if apn_policy_set.apn != apn:
                continue
            # Only add rules if the APN matches
            static_rule_ids = self._get_static_rules(apn_policy_set)
            assigned_static_rules +=\
                list(
                    map(
                        lambda id: StaticRuleInstall(rule_id=id),
                        static_rule_ids,
                    ),
                )

        return assigned_static_rules

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

    def _get_static_rules(
        self,
        policies: ApnPolicySet,
    ) -> Set[str]:
        desired_rules = set(policies.assigned_policies)
        for basename in policies.assigned_base_names:
            if basename not in self._rules_by_basename:
                # Eventually, basename definition will be streamed from orc8r
                continue
            desired_rules.update(
                self._rules_by_basename[basename].RuleNames,
            )
        return desired_rules

    def _get_credits(self, sid: str) -> List[CreditUpdateResponse]:
        infinite_credit_keys = self.get_infinite_credit_charging_keys()
        postpay_keys = self._get_postpay_charging_keys()
        credit_updates = []
        for charging_key in infinite_credit_keys:
            credit_updates.append(
                CreditUpdateResponse(
                    success=True,
                    sid=sid,
                    charging_key=charging_key,
                    limit_type=CreditLimitType.Value("INFINITE_UNMETERED"),
                ),
            )
        for charging_key in postpay_keys:
            credit_updates.append(
                CreditUpdateResponse(
                    success=True,
                    sid=sid,
                    charging_key=charging_key,
                    limit_type=CreditLimitType.Value("INFINITE_METERED"),
                ),
            )
        return credit_updates
