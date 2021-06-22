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
import itertools
import unittest
from unittest.mock import Mock

from integ_tests.gxgy_tests.policies import (
    create_uplink_rule,
    get_packets_for_flows,
)
from integ_tests.gxgy_tests.session_manager import create_update_response
from lte.protos import session_manager_pb2
from lte.protos.policydb_pb2 import PolicyRule
from lte.protos.session_manager_pb2 import (
    CreateSessionResponse,
    LocalCreateSessionRequest,
    PolicyReAuthRequest,
    SessionTerminateResponse,
)
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.pipelined.tests.app.subscriber import (
    SubContextConfig,
    default_ambr_config,
)
from ryu.lib import hub

from .utils import GxGyTestUtil as TestUtil


class GxReauthTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        super(GxReauthTest, cls).setUpClass()

        # Static policies
        cls.test_util = TestUtil()
        policy1 = create_uplink_rule(
            'policy1', 1, '45.10.0.1',
            tracking=PolicyRule.NO_TRACKING,
        )
        cls.test_util.static_rules[policy1.id] = policy1
        policy2 = create_uplink_rule(
            'policy2', 1, '45.10.10.2',
            tracking=PolicyRule.NO_TRACKING,
        )
        cls.test_util.static_rules[policy2.id] = policy2
        hub.sleep(2)  # wait for static rule to sync

    @classmethod
    def tearDownClass(cls):
        cls.test_util.cleanup()

    def test_reauth(self):
        """
        Send a Gx reauth request which installs one new static rule, one new
        dynamic rule, and removes one static and one dynamic rule.
        """
        dynamic_rule1 = create_uplink_rule(
            'dynamic1', 1, '46.10.10.1',
            tracking=PolicyRule.NO_TRACKING,
        )
        dynamic_rule2 = create_uplink_rule(
            'dynamic2', 1, '46.10.10.2',
            tracking=PolicyRule.NO_TRACKING,
        )

        # Initialize sub with 1 static and 1 dynamic rule
        sub = SubContextConfig(
            'IMSI001010000088888',
            '192.168.128.74', default_ambr_config, 4,
        )
        self.test_util.controller.mock_create_session = Mock(
            return_value=CreateSessionResponse(
                credits=[create_update_response(sub.imsi, 1, 1024)],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id='policy1',
                    ),
                ],
                dynamic_rules=[
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=dynamic_rule1,
                    ),
                ],
                usage_monitors=[],
            ),
        )
        self.test_util.controller.mock_terminate_session = Mock(
            return_value=SessionTerminateResponse(),
        )
        self.test_util.sessiond.CreateSession(
            LocalCreateSessionRequest(
                sid=SubscriberID(id=sub.imsi),
                ue_ipv4=sub.ip,
            ),
        )
        self.assertEqual(
            self.test_util.controller.mock_create_session.call_count,
            1,
        )

        # first, send some packets so we know that the uplink rules are
        # accepting traffic
        self._assert_rules(
            sub,
            [
                session_manager_pb2.DynamicRuleInstall(
                    policy_rule=self.test_util.static_rules['policy1'],
                ),
                session_manager_pb2.DynamicRuleInstall(
                    policy_rule=dynamic_rule1,
                ),
            ],
        )

        # Now via reauth, remove the old rules and install new uplink rules
        # Verify the new uplink rules allow traffic
        reauth_result = self.test_util.proxy_responder.PolicyReAuth(
            PolicyReAuthRequest(
                imsi=sub.imsi,
                rules_to_remove=['dynamic1', 'policy1'],
                rules_to_install=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id='policy2',
                    ),
                ],
                dynamic_rules_to_install=[
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=dynamic_rule2,
                    ),
                ],
            ),
        )
        self.assertEqual(
            reauth_result.result,
            session_manager_pb2.UPDATE_INITIATED,
        )
        self.assertEqual(len(reauth_result.failed_rules), 0)
        self._assert_rules(
            sub,
            [
                session_manager_pb2.DynamicRuleInstall(
                    policy_rule=self.test_util.static_rules['policy2'],
                ),
                session_manager_pb2.DynamicRuleInstall(
                    policy_rule=dynamic_rule2,
                ),
            ],
        )

        # Verify the old rules no longer allow traffic (uninstalled)
        self._assert_rules(
            sub,
            [
                session_manager_pb2.DynamicRuleInstall(
                    policy_rule=self.test_util.static_rules['policy1'],
                ),
                session_manager_pb2.DynamicRuleInstall(
                    policy_rule=dynamic_rule1,
                ),
            ],
            expected=0,
        )

    def _assert_rules(self, sub, rules, expected=-1):
        flows = list(
            itertools.chain(*[rule.policy_rule.flow_list for rule in rules]),
        )
        packets = get_packets_for_flows(sub, flows)
        packet_sender = self.test_util.get_packet_sender([sub], packets, 1)

        num_packets = self.test_util.thread.run_in_greenthread(packet_sender)
        if expected == -1:
            self.assertEqual(num_packets, len(flows))
        else:
            self.assertEqual(num_packets, expected)
