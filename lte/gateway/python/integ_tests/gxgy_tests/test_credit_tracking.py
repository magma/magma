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
from magma.pipelined.tests.app.subscriber import SubContextConfig
from ryu.lib import hub

from .policies import create_uplink_rule, get_packets_for_flows
from .session_manager import (
    create_update_response,
    get_from_queue,
    get_standard_update_response,
)
from .utils import GxGyTestUtil as TestUtil


class CreditTrackingTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        super(CreditTrackingTest, cls).setUpClass()
        # Static policies
        cls.test_util = TestUtil()
        policy = create_uplink_rule("simple_match", 1, '45.10.0.1')
        cls.test_util.static_rules[policy.id] = policy
        hub.sleep(2)  # wait for static rule to sync

    @classmethod
    def tearDownClass(cls):
        cls.test_util.cleanup()

    def test_basic_init(self):
        """
        Initiate subscriber, return 1 static policy, send traffic to match the
        policy, verify update is sent, terminate subscriber
        """
        sub1 = SubContextConfig('IMSI001010000088888', '192.168.128.74', 4)
        quota = 1024  # bytes

        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[create_update_response(sub1.imsi, 1, quota)],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="simple_match",
                    ),
                ],
                dynamic_rules=[],
                usage_monitors=[],
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        update_complete = hub.Queue()
        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                update_complete, None, quota, is_final=False,
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
            sub1, self.test_util.static_rules["simple_match"].flow_list,
        )
        packet_count = int(quota / len(packets[0])) + 1

        self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender([sub1], packets, packet_count),
        )
        self.assertIsNotNone(get_from_queue(update_complete))
        self.assertEqual(
            self.test_util.controller.mock_update_session.call_count, 1,
        )

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )

    def test_input_output(self):
        """
        """
        sub1 = SubContextConfig('IMSI001010000088888', '192.168.128.74', 4)
        quota = 1024  # bytes

        # return only rx (downlink) packets
        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[
                    session_manager_pb2.CreditUpdateResponse(
                        success=True,
                        sid=sub1.imsi,
                        charging_key=1,
                        credit=session_manager_pb2.ChargingCredit(
                            granted_units=session_manager_pb2.GrantedUnits(
                                rx=session_manager_pb2.CreditUnit(
                                    is_valid=True,
                                    volume=quota,
                                ),
                            ),
                        ),
                    ),
                ],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="simple_match",
                    ),
                ],
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        update_complete = hub.Queue()
        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                update_complete, None, quota, is_final=False,
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
            sub1, self.test_util.static_rules["simple_match"].flow_list,
        )
        packet_count = int(quota / len(packets[0])) + 1

        self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender([sub1], packets, packet_count),
        )
        self.assertIsNone(get_from_queue(update_complete))
        self.assertEqual(
            self.test_util.controller.mock_update_session.call_count, 0,
        )

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )

        # now attach with tx (uplink packets)
        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[
                    session_manager_pb2.CreditUpdateResponse(
                        success=True,
                        sid=sub1.imsi,
                        charging_key=1,
                        credit=session_manager_pb2.ChargingCredit(
                            granted_units=session_manager_pb2.GrantedUnits(
                                tx=session_manager_pb2.CreditUnit(
                                    is_valid=True,
                                    volume=quota,
                                ),
                            ),
                        ),
                    ),
                ],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="simple_match",
                    ),
                ],
            ),
        )
        self.test_util.sessiond.CreateSession(
            session_manager_pb2.LocalCreateSessionRequest(
                sid=SubscriberID(id=sub1.imsi),
                ue_ipv4=sub1.ip,
            ),
        )
        self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender([sub1], packets, packet_count),
        )
        self.assertIsNotNone(get_from_queue(update_complete))
        self.assertEqual(
            self.test_util.controller.mock_update_session.call_count, 1,
        )
        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 2,
        )

    def test_out_of_credit(self):
        """
        Initiate subscriber, return 1 static policy, send traffic to match the
        policy, verify update is sent, return final credits, use up final
        credits, ensure that no traffic can be sent
        """
        sub1 = SubContextConfig('IMSI001010000088888', '192.168.128.74', 4)
        quota = 1024  # bytes

        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[
                    session_manager_pb2.CreditUpdateResponse(
                        success=True,
                        sid=sub1.imsi,
                        charging_key=1,
                        credit=session_manager_pb2.ChargingCredit(
                            granted_units=session_manager_pb2.GrantedUnits(
                                total=session_manager_pb2.CreditUnit(
                                    is_valid=True,
                                    volume=quota,
                                ),
                            ),
                        ),
                    ),
                ],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="simple_match",
                    ),
                ],
                dynamic_rules=[],
                usage_monitors=[],
            ),
        )

        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )

        update_complete = hub.Queue()
        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                update_complete, None, quota, is_final=True,
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
            sub1, self.test_util.static_rules["simple_match"].flow_list,
        )
        packet_count = int(quota / len(packets[0])) + 1
        send_packets = self.test_util.get_packet_sender(
            [sub1], packets, packet_count,
        )

        self.test_util.thread.run_in_greenthread(send_packets)
        self.assertIsNotNone(get_from_queue(update_complete))
        self.assertEqual(
            self.test_util.controller.mock_update_session.call_count, 1,
        )

        # use up last credits
        self.test_util.thread.run_in_greenthread(send_packets)
        hub.sleep(3)  # wait for sessiond to terminate rule after update

        pkt_diff = self.test_util.thread.run_in_greenthread(send_packets)
        self.assertEqual(pkt_diff, 0)

        self.test_util.proxy_responder.ChargingReAuth(
            session_manager_pb2.ChargingReAuthRequest(
                charging_key=1,
                sid=sub1.imsi,
            ),
        )
        get_from_queue(update_complete)
        self.assertEqual(
            self.test_util.controller.mock_update_session.call_count, 2,
        )
        # wait for 1 update to trigger credit request, another to trigger
        # rule activation
        # TODO Add future to track when flows are added/deleted
        hub.sleep(5)
        pkt_diff = self.test_util.thread.run_in_greenthread(send_packets)
        self.assertGreater(pkt_diff, 0)

        self.test_util.sessiond.EndSession(SubscriberID(id=sub1.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, 1,
        )

    def test_multiple_subscribers(self):
        """
        Test credit tracking with multiple rules and 32 subscribers, each using
        up their quota and reporting to the OCS
        """
        subs = [
            SubContextConfig(
                'IMSI0010100000888{}'.format(i),
                '192.168.128.{}'.format(i),
                4,
            ) for i in range(32)
        ]
        quota = 1024  # bytes

        # create some rules
        rule1 = create_uplink_rule("rule1", 2, '46.10.0.1')
        rule2 = create_uplink_rule(
            "rule2", 0, '47.10.0.1',
            tracking=PolicyRule.NO_TRACKING,
        )
        rule3 = create_uplink_rule("rule3", 3, '49.10.0.1')
        self.test_util.static_rules["rule1"] = rule1
        self.test_util.static_rules["rule2"] = rule2
        hub.sleep(2)  # wait for policies

        # set up mocks
        self.test_util.controller.mock_create_session = Mock(
            return_value=session_manager_pb2.CreateSessionResponse(
                credits=[
                    create_update_response("", 2, quota),
                    create_update_response("", 3, quota),
                ],
                static_rules=[
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="rule1",
                    ),
                    session_manager_pb2.StaticRuleInstall(
                        rule_id="rule2",
                    ),
                ],
                dynamic_rules=[
                    session_manager_pb2.DynamicRuleInstall(
                        policy_rule=rule3,
                    ),
                ],
            ),
        )
        self.test_util.controller.mock_terminate_session = Mock(
            return_value=session_manager_pb2.SessionTerminateResponse(),
        )
        update_complete = hub.Queue()
        self.test_util.controller.mock_update_session = Mock(
            side_effect=get_standard_update_response(
                update_complete, None, quota, is_final=True,
            ),
        )

        # initiate sessions
        for sub in subs:
            self.test_util.sessiond.CreateSession(
                session_manager_pb2.LocalCreateSessionRequest(
                    sid=SubscriberID(id=sub.imsi),
                    ue_ipv4=sub.ip,
                ),
            )
        self.assertEqual(
            self.test_util.controller.mock_create_session.call_count, len(
                subs,
            ),
        )

        # send packets towards all 3 rules
        flows = [rule.flow_list[0] for rule in [rule1, rule2, rule3]]
        packets = []
        for sub in subs:
            packets.extend(get_packets_for_flows(sub, flows))
        packet_count = int(quota / len(packets[0])) + 1
        self.test_util.thread.run_in_greenthread(
            self.test_util.get_packet_sender(subs, packets, packet_count),
        )

        # wait for responses for keys 2 and 3 (key 1 is not tracked)
        expected_keys = {(sub.imsi, key) for sub in subs for key in [2, 3]}
        for _ in range(len(expected_keys)):
            update = get_from_queue(update_complete)
            self.assertIsNotNone(update)
            imsiKey = (update.sid, update.usage.charging_key)
            self.assertTrue(imsiKey in expected_keys)
            expected_keys.remove(imsiKey)

        for sub in subs:
            self.test_util.sessiond.EndSession(SubscriberID(id=sub.imsi))
        self.assertEqual(
            self.test_util.controller.mock_terminate_session.call_count, len(
                subs,
            ),
        )


if __name__ == "__main__":
    unittest.main()
