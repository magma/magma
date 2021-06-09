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
from typing import Set

import grpc
from lte.protos.policydb_pb2 import InstalledPolicies
from lte.protos.session_manager_pb2 import (
    PolicyReAuthAnswer,
    PolicyReAuthRequest,
    ReAuthResult,
)
from lte.protos.session_manager_pb2_grpc import SessionProxyResponderStub
from magma.policydb.rule_map_store import RuleAssignmentsDict


class ReAuthHandler():
    """
    Handles making the RAR to sessiond, and updating Redis with the current
    active policies for subscribers
    """

    def __init__(
        self,
        rules_by_sid: RuleAssignmentsDict,
        sessiond_stub: SessionProxyResponderStub,
    ):
        self._rules_by_sid = rules_by_sid
        self._sessiond_stub = sessiond_stub

    def handle_policy_re_auth(self, rar: PolicyReAuthRequest) -> bool:
        """
        Take the RAR and make the request to sessiond

        Returns whether the call was successful. Partial updates count as
        success.
        """
        if not self._is_valid_rar(rar):
            logging.error(
                'Invalid RAR: Either installing already installed '
                'rules, or uninstalling rules that are not installed',
            )
            return False
        try:
            resp = self._sessiond_stub.PolicyReAuth(rar)
            return self._handle_rar_answer(rar, resp)
        except grpc.RpcError:
            logging.error(
                'Unable to apply policy updates for subscriber %s',
                rar.imsi,
            )
            return False

    def _is_valid_rar(self, rar: PolicyReAuthRequest) -> bool:
        """
        Return false if the RAR is invalid
        RAR is invalid if attempting to remove rules that are not installed,
        or trying to add rules that are already installed.
        """
        prev_rules = self._get_prev_policies(rar.imsi)
        install = {rule.rule_id for rule in rar.rules_to_install}
        if install & prev_rules:
            return False
        if len(set(rar.rules_to_remove) - prev_rules) > 0:
            return False
        return True

    def _handle_rar_answer(
        self,
        rar: PolicyReAuthRequest,
        answer: PolicyReAuthAnswer,
    ) -> bool:
        if answer.result == ReAuthResult.Value('OTHER_FAILURE'):
            logging.error(
                'Failed to apply policy updates for subscriber %s',
                rar.imsi,
            )
            return False
        self._rules_by_sid[rar.imsi] = InstalledPolicies(
            installed_policies=list(self._get_updated_rules(rar, answer)),
        )
        return True

    def _get_updated_rules(
        self,
        rar: PolicyReAuthRequest,
        answer: PolicyReAuthAnswer,
    ) -> Set[str]:
        failed = set(answer.failed_rules)
        installed = {rule.rule_id for rule in rar.rules_to_install} - failed
        uninstalled = set(rar.rules_to_remove) - failed
        return (self._get_prev_policies(rar.imsi) | installed) - uninstalled

    def _get_prev_policies(self, subscriber_id: str) -> Set[str]:
        if subscriber_id not in self._rules_by_sid:
            return set()
        return set(self._rules_by_sid[subscriber_id].installed_policies)
