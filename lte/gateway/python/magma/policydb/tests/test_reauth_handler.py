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

import unittest

from lte.protos.session_manager_pb2 import (
    PolicyReAuthAnswer,
    PolicyReAuthRequest,
    ReAuthResult,
    StaticRuleInstall,
)
from magma.policydb.reauth_handler import ReAuthHandler


class MockSessionProxyResponderStub:
    """
    This Mock SessionProxyResponderStub will always respond with a success to
    a received RAR
    """

    def __init__(self):
        pass

    def PolicyReAuth(self, _: PolicyReAuthRequest) -> PolicyReAuthAnswer:
        return PolicyReAuthAnswer(
            result=ReAuthResult.Value('UPDATE_INITIATED'),
        )


class RuleMappingsStreamerCallbackTest(unittest.TestCase):
    def test_SuccessfulUpdate(self):
        """
        Test the happy path where updates come in for added rules, and sessiond
        accepts the RAR without issue.
        """
        install_dict = {}
        handler = ReAuthHandler(
            install_dict,
            MockSessionProxyResponderStub(),
        )

        rar = PolicyReAuthRequest(
            imsi='s1',
            rules_to_install=[
                StaticRuleInstall(rule_id=rule_id) for rule_id in ['p1', 'p2']
            ],
        )
        handler.handle_policy_re_auth(rar)
        s1_policies = install_dict['s1'].installed_policies
        expected = 2
        self.assertEqual(
            len(s1_policies), expected,
            'There should be 2 installed policies for s1',
        )
        self.assertTrue(
            'p1' in s1_policies,
            'Policy p1 should be marked installed for p1',
        )

        rar = PolicyReAuthRequest(
            imsi='s1',
            rules_to_install=[StaticRuleInstall(rule_id='p3')],
        )
        handler.handle_policy_re_auth(rar)
        s1_policies = install_dict['s1'].installed_policies
        expected = 3
        self.assertEqual(
            len(s1_policies), expected,
            'There should be 3 installed policies for s1',
        )
        self.assertTrue(
            'p3' in s1_policies,
            'Policy p1 should be marked installed for p1',
        )
