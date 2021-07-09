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
from unittest.mock import Mock

from lte.protos import session_manager_pb2
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.pipelined.tests.app.subscriber import (
    SubContextConfig,
    default_ambr_config,
)
from ryu.lib import hub

from .policies import create_uplink_rule, get_packets_for_flows
from .session_manager import (
    create_update_response,
    get_from_queue,
    get_standard_update_response,
)
from .utils import GxGyTestUtil as TestUtil


class FailureScenarioTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        super(FailureScenarioTest, cls).setUpClass()
        # Static policies
        cls.test_util = TestUtil()
        policy = create_uplink_rule("simple_match", 1, '45.10.0.1')
        cls.test_util.static_rules[policy.id] = policy

    @classmethod
    def tearDownClass(cls):
        cls.test_util.cleanup()

    def test_rule_with_no_credit(self):
        """
        Test that when a rule is returned that requires OCS tracking but has
        no credit, data is not allowed to pass
        """
        sub1 = SubContextConfig(
            'IMSI001010000088888',
            '192.168.128.74', default_ambr_config, 4,
        )

        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="simple_match",
                    ),
                ],  # no credit for RG 1
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        self.test_util.sessiond.CreateSession(
            session_manager_pb2.LocalCreateSessionRequest(
                sid=SubscriberID(id=sub1.imsi),
                ue_ipv4=sub1.ip,
            ),
        )
        self.assertEqual(
            self.test_util.controller.mock_create_session.call_count, 1,
        )

        packets = get_packets_for_flows(
            sub1, self.test_util.static_rules["simple_match"].flow_list,
        )

        pkt_diff = self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender([sub1], packets, 1),
        )
        self.assertEqual(pkt_diff, 0)

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )

    def test_rules_with_failed_credit(self):
        """
        Test that when a session is initialized but the OCS either errored out or
        returned 0 GSUs, data is not allowed to flow
        """
        sub1 = SubContextConfig(
            'IMSI001010000088888',
            '192.168.128.74', default_ambr_config, 4,
        )

        rule2 = create_uplink_rule("rule2", 2, '46.10.0.1')
        rule3 = create_uplink_rule("rule3", 3, '47.10.0.1')
        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[
                    # failed update
                    create_update_response(sub1.imsi, 1, 0, success=False),
                    # successful update, no credit
                    create_update_response(sub1.imsi, 1, 0, success=True),
                ],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="simple_match",
                    ),
                ],  # no credit for RG 1
                dynamic_rules=[
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=rule2,
                    ),
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=rule3,
                    ),
                ],
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        self.test_util.sessiond.CreateSession(
            session_manager_pb2.LocalCreateSessionRequest(
                sid=SubscriberID(id=sub1.imsi),
                ue_ipv4=sub1.ip,
            ),
        )
        self.assertEqual(
            self.test_util.controller.mock_create_session.call_count, 1,
        )

        flows = [rule.flow_list[0] for rule in [rule2, rule3]]
        packets = get_packets_for_flows(sub1, flows)
        pkt_diff = self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender([sub1], packets, 1),
        )
        self.assertEqual(pkt_diff, 0)

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )

    def test_ocs_failure(self):
        """
        Test that when the OCS fails to respond to an update request, the service
        is cut off until the update can be completed
        """
        sub1 = SubContextConfig(
            'IMSI001010000088888',
            '192.168.128.74', default_ambr_config, 4,
        )
        quota = 1024

        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[create_update_response(sub1.imsi, 1, quota)],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="simple_match",
                    ),
                ],
            ),
        )

        update_complete = hub.Queue()
        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                update_complete, None, quota, success=False,
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        self.test_util.sessiond.CreateSession(
            session_manager_pb2.LocalCreateSessionRequest(
                sid=SubscriberID(id=sub1.imsi),
                ue_ipv4=sub1.ip,
            ),
        )
        self.assertEqual(
            self.test_util.controller.mock_create_session.call_count, 1,
        )

        packets = get_packets_for_flows(
            sub1, self.test_util.static_rules["simple_match"].flow_list,
        )
        packet_count = int(quota / len(packets[0])) + 1
        sender = self.test_util.get_packet_sender(
            [sub1], packets, packet_count,
        )

        # assert after session init, data can flow
        self.assertGreater(self.test_util.thread.run_in_greenthread(sender), 0)

        # wait for failed update
        self.assertIsNotNone(get_from_queue(update_complete))
        hub.sleep(2)

        # assert that no data can be sent anymore
        self.assertEqual(self.test_util.thread.run_in_greenthread(sender), 0)

        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                update_complete, None, quota, success=True,
            ),
        )
        # wait for second update cycle to reactivate
        hub.sleep(4)
        self.assertGreater(self.test_util.thread.run_in_greenthread(sender), 0)

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )


if __name__ == "__main__":
    unittest.main()
