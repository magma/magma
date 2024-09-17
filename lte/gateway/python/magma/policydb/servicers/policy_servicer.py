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
from typing import List

import grpc
from lte.protos.policydb_pb2 import (
    DisableStaticRuleRequest,
    EnableStaticRuleRequest,
)
from lte.protos.policydb_pb2_grpc import (
    PolicyAssignmentControllerStub,
    PolicyDBServicer,
    add_PolicyDBServicer_to_server,
)
from lte.protos.session_manager_pb2 import (
    PolicyReAuthRequest,
    StaticRuleInstall,
)
from magma.policydb.basename_store import BaseNameDict
from magma.policydb.reauth_handler import ReAuthHandler
from orc8r.protos.common_pb2 import Void


class PolicyRpcServicer(PolicyDBServicer):
    """
    gRPC based server for PolicyDB service.

    This will act as a bare-bones local PCRF and OCS.
    In current implementation, it is only used for enabling the Captive Portal
    feature.
    """

    def __init__(
        self,
        reauth_handler: ReAuthHandler,
        rules_by_basename: BaseNameDict,
        subscriberdb_stub: PolicyAssignmentControllerStub,
    ):
        self._reauth_handler = reauth_handler
        self._rules_by_basename = rules_by_basename
        self._subscriberdb_stub = subscriberdb_stub

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_PolicyDBServicer_to_server(
            self, server,
        )

    def EnableStaticRules(
        self,
        request: EnableStaticRuleRequest,
        context,
    ) -> Void:
        """
        Associate the static rules with the specified subscriber.
        Also send a RAR to sessiond to install the specified rules for the
        subscriber.
        """
        try:
            self._subscriberdb_stub.EnableStaticRules(request)
        except grpc.RpcError:
            logging.error(
                'Unable to enable rules for subscriber %s. ',
                request.imsi,
            )
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('Failed to update rule assignments in orc8r')
            return Void()

        rules_to_install = self._get_rules(
            request.rule_ids, request.base_names,
        )
        rar = PolicyReAuthRequest(
            # Leave session id empty, re-auth for all sessions
            imsi=request.imsi,
            rules_to_install=[
                StaticRuleInstall(rule_id=rule_id)
                for rule_id in rules_to_install
            ],
        )
        success = self._reauth_handler.handle_policy_re_auth(rar)
        if not success:
            context.set_code(grpc.StatusCode.UNKNOWN)
            context.set_details(
                'Failed to enable all static rules for '
                'subscriber. Partial update may have succeeded',
            )
        return Void()

    def DisableStaticRules(
        self,
        request: DisableStaticRuleRequest,
        context,
    ) -> Void:
        """
        Unassociate the static rules with the specified subscriber.
        Also send a RAR to sessiond to uninstall the specified rules for the
        subscriber.
        """
        try:
            self._subscriberdb_stub.DisableStaticRules(request)
        except grpc.RpcError:
            logging.error(
                'Unable to disable rules for subscriber %s. ',
                request.imsi,
            )
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('Failed to update rule assignments in orc8r')
            return Void()

        rar = PolicyReAuthRequest(
            # Leave session id empty, re-auth for all sessions
            imsi=request.imsi,
            rules_to_remove=self._get_rules(
                request.rule_ids,
                request.base_names,
            ),
        )

        success = self._reauth_handler.handle_policy_re_auth(rar)
        if not success:
            context.set_code(grpc.StatusCode.UNKNOWN)
            context.set_details(
                'Failed to enable all static rules for '
                'subscriber. Partial update may have succeeded',
            )
        return Void()

    def _get_rules(
        self,
        rule_ids: List[str],
        basenames: List[str],
    ) -> List[str]:
        rules = set(rule_ids)
        for basename in basenames:
            if basename in self._rules_by_basename:
                rules.update(self._rules_by_basename[basename].RuleNames)
        return list(rules)
