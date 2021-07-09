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
from lte.protos.policydb_pb2 import PolicyRule
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.pipelined.tests.app.subscriber import (
    SubContextConfig,
    default_ambr_config,
)
from ryu.lib import hub

from .policies import create_uplink_rule, get_packets_for_flows
from .session_manager import (
    create_monitor_response,
    create_update_response,
    get_from_queue,
    get_standard_update_response,
)
from .utils import GxGyTestUtil as TestUtil


class UsageMonitorTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        super(UsageMonitorTest, cls).setUpClass()

        cls.test_util = TestUtil()
        # default rule
        policy = create_uplink_rule(
            "monitor_rule", 0, '45.10.0.1',
            m_key="mkey1",
            tracking=PolicyRule.ONLY_PCRF,
        )
        cls.test_util.static_rules[policy.id] = policy
        hub.sleep(2)  # wait for static rule to sync

    @classmethod
    def tearDownClass(cls):
        cls.test_util.cleanup()

    def test_basic_init(self):
        """
        Initiate subscriber, return 1 static policy with monitoring key, send
        traffic to match the policy, verify monitoring update is sent, terminate
        subscriber
        """
        sub1 = SubContextConfig(
            'IMSI001010000088888',
            '192.168.128.74', default_ambr_config, 4,
        )
        quota = 1024  # bytes

        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="monitor_rule",
                    ),
                ],
                dynamic_rules=[],
                usage_monitors=[
                    create_monitor_response(
                        sub1.imsi, "mkey1", quota, session_manager_pb2.PCC_RULE_LEVEL,
                    ),
                ],
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        monitor_complete = hub.Queue()
        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                None, monitor_complete, quota,
            ),
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
            sub1, self.test_util.static_rules["monitor_rule"].flow_list,
        )
        packet_count = int(quota / len(packets[0])) + 1

        self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender([sub1], packets, packet_count),
        )
        self.assertIsNotNone(get_from_queue(monitor_complete))
        self.assertEqual(
            self.test_util.controller.mock_update_session.call_count, 1,
        )

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )

    def test_mixed_monitors_and_updates(self):
        """
        Test a mix of usage monitors, session monitors, and charging credits to
        PCRF and OCS.
        """
        sub1 = SubContextConfig(
            'IMSI001010000088888',
            '192.168.128.74', default_ambr_config, 4,
        )
        quota = 1024  # bytes

        pcrf_rule = create_uplink_rule(
            "pcrf_rule", 0, '46.10.0.1',
            m_key="key1",
            tracking=PolicyRule.ONLY_PCRF,
        )
        ocs_rule = create_uplink_rule(
            "ocs_rule", 1, '47.10.0.1',
            tracking=PolicyRule.ONLY_OCS,
        )
        both_rule = create_uplink_rule(
            "both_rule", 2, '48.10.0.1',
            m_key="key2",
            tracking=PolicyRule.OCS_AND_PCRF,
        )

        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[
                    create_update_response("", 1, quota),
                    create_update_response("", 2, quota),
                ],
                dynamic_rules=[
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=pcrf_rule,
                    ),
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=ocs_rule,
                    ),
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=both_rule,
                    ),
                ],
                usage_monitors=[
                    create_monitor_response(
                        sub1.imsi,
                        "key1",
                        quota,
                        session_manager_pb2.PCC_RULE_LEVEL,
                    ),
                    create_monitor_response(
                        sub1.imsi,
                        "key2",
                        quota,
                        session_manager_pb2.PCC_RULE_LEVEL,
                    ),
                    create_monitor_response(
                        sub1.imsi,
                        "key3",
                        quota,
                        session_manager_pb2.SESSION_LEVEL,
                    ),
                ],
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        charging_complete = hub.Queue()
        monitor_complete = hub.Queue()
        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                charging_complete, monitor_complete, quota,
            ),
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
        flows = [
            rule.flow_list[0]
            for rule in [pcrf_rule, ocs_rule, both_rule]
        ]
        packets = get_packets_for_flows(sub1, flows)
        packet_count = int(quota / len(packets[0])) + 1
        self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender([sub1], packets, packet_count),
        )

        # Wait for responses for keys 1 and 2 (ocs_rule and both_rule)
        charging_keys = {1, 2}
        for _ in range(len(charging_keys)):
            update = get_from_queue(charging_complete)
            self.assertTrue(update.usage.charging_key in charging_keys)
            charging_keys.remove(update.usage.charging_key)

        # Wait for responses for mkeys key1 (pcrf_rule), key2 (both_rule),
        # key3 (session rule)
        monitoring_keys = ["key1", "key2", "key3"]
        for _ in range(len(monitoring_keys)):
            monitor = get_from_queue(monitor_complete)
            self.assertTrue(monitor.update.monitoring_key in monitoring_keys)
            monitoring_keys.remove(monitor.update.monitoring_key)

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )


if __name__ == "__main__":
    unittest.main()
