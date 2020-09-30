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

from typing import Callable, Dict, List
import unittest
from unittest.mock import Mock
from lte.protos.policydb_pb2 import AssignedPolicies, ChargingRuleNameSet,\
    SubscriberPolicySet, ApnPolicySet, PolicyRule, FlowDescription, FlowMatch
from lte.protos.session_manager_pb2 import SessionRules, RulesPerSubscriber,\
    RuleSet, StaticRuleInstall, DynamicRuleInstall
from magma.policydb.streamer_callback import ApnRuleMappingsStreamerCallback,\
    RuleMappingsStreamerCallback
from magma.policydb.reauth_handler import ReAuthHandler
from magma.policydb.tests.mock_stubs import MockSessionProxyResponderStub1, \
    MockSessionProxyResponderStub2, MockSessionProxyResponderStub3, \
    MockLocalSessionManagerStub
from orc8r.protos.common_pb2 import Void
from orc8r.protos.streamer_pb2 import DataUpdate


class RuleMappingsStreamerCallbackTest(unittest.TestCase):
    def test_SuccessfulUpdate(self):
        """
        Test the happy path where updates come in for added rules, and sessiond
        accepts the RAR without issue.
        """
        assignments_dict = {}
        apn_rules_dict = {} # type: Dict[str, SubscriberPolicySet]
        basenames_dict = {
            'bn1': ChargingRuleNameSet(RuleNames=['p5']),
            'bn2': ChargingRuleNameSet(RuleNames=['p6']),
        }
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub1()),
            basenames_dict,
            assignments_dict,
            apn_rules_dict,
        )

        # Construct a set of updates, keyed by subscriber ID
        updates = [
            DataUpdate(
                key="s1",
                value=AssignedPolicies(
                    assigned_policies=["p1", "p2"],
                    assigned_base_names=["bn1"],
                ).SerializeToString(),
            ),
            DataUpdate(
                key="s2",
                value=AssignedPolicies(
                    assigned_policies=["p2", "p3"],
                ).SerializeToString(),
            ),
        ]

        callback.process_update("stream", updates, False)

        # Since we used a stub which always succeeds when a RAR is made,
        # We should expect the assignments_dict to be updated

        s1_policies = assignments_dict["s1"].installed_policies
        expected = 3
        self.assertEqual(len(s1_policies), expected, 'There should be 3 active '
                         'policies for s1')
        self.assertTrue("p1" in s1_policies, 'Policy p1 should be active for '
                        'subscriber s1')
        self.assertTrue("p5" in s1_policies, 'Policy p5 should be active for '
                                             'subscriber s1')

        s2_policies = assignments_dict["s2"].installed_policies
        expected = 2
        self.assertEqual(len(s2_policies), expected, 'There should be 2 active '
                                                     'policies for s2')
        self.assertTrue("p3" in s2_policies, 'Policy p3 should be active for '
                                             'subscriber s2')

        # Check the ApnRuleAssignmentsDict too
        s1_policies = apn_rules_dict["s1"].global_policies
        s1_base_names = apn_rules_dict["s1"].global_base_names
        self.assertEqual(len(s1_policies), 2,
                         'There should be 2 global policies for s1')
        self.assertTrue("p1" in s1_policies,
                        'Policy p1 should be global for subscriber s1')
        self.assertTrue("bn1" in s1_base_names,
                        'Base name bn1 should be global for subscriber s1')

        s2_policies = apn_rules_dict["s2"].global_policies
        self.assertEqual(len(s2_policies), 2,
                         'There should be 2 global policies for s2')
        self.assertTrue("p3" in s2_policies,
                        'Policy p3 should be global for subscriber s2')


    def test_FailedUpdate(self):
        """
        Test when sessiond answers to the RAR with a failure for any re-auth.
        """
        assignments_dict = {}
        apn_rules_dict = {} # type: Dict[str, SubscriberPolicySet]
        basenames_dict = {}
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub2()),
            basenames_dict,
            assignments_dict,
            apn_rules_dict,
        )

        # Construct a set of updates, keyed by subscriber ID
        updates = [
            DataUpdate(
                key="s1",
                value=AssignedPolicies(
                    assigned_policies=["p1", "p2"],
                ).SerializeToString(),
            ),
            DataUpdate(
                key="s2",
                value=AssignedPolicies(
                    assigned_policies=["p2", "p3"],
                ).SerializeToString(),
            ),
        ]

        callback.process_update("stream", updates, False)

        # Since we used a stub which always succeeds when a RAR is made,
        # We should expect the assignments_dict to be updated

        self.assertFalse("s1" in assignments_dict, 'There should be no entry '
                         'for subscriber s1 since update failed')
        self.assertFalse("s2" in assignments_dict, 'There should be no entry '
                         'for subscriber s2 since update failed')

    def test_FailedOnePolicy(self):
        """
        Test when sessiond answers to the RAR with a failure for installing p2.
        """
        assignments_dict = {}
        apn_rules_dict = {} # type: Dict[str, SubscriberPolicySet]
        basenames_dict = {}
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub3()),
            basenames_dict,
            assignments_dict,
            apn_rules_dict,
        )

        # Construct a set of updates, keyed by subscriber ID
        updates = [
            DataUpdate(
                key="s1",
                value=AssignedPolicies(
                    assigned_policies=["p1", "p2"],
                ).SerializeToString(),
            ),
            DataUpdate(
                key="s2",
                value=AssignedPolicies(
                    assigned_policies=["p2", "p3"],
                ).SerializeToString(),
            ),
        ]

        callback.process_update("stream", updates, False)

        s1_policies = assignments_dict["s1"].installed_policies
        expected = 1
        self.assertEqual(len(s1_policies), expected, 'There should be 1 active '
                                                     'policies for s1')
        self.assertTrue("p1" in s1_policies, 'Policy p1 should be active for '
                                             'subscriber s1')

        s2_policies = assignments_dict["s2"].installed_policies
        expected = 1
        self.assertEqual(len(s2_policies), expected, 'There should be 1 active '
                                                     'policies for s2')
        self.assertTrue("p3" in s2_policies, 'Policy p3 should be active for '
                                             'subscriber s2')

    def test_MultiUpdate(self):
        """
        Test consecutive updates
        """
        assignments_dict = {}
        apn_rules_dict = {} # type: Dict[str, SubscriberPolicySet]
        basenames_dict = {}
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub3()),
            basenames_dict,
            assignments_dict,
            apn_rules_dict,
        )

        # Construct a set of updates, keyed by subscriber ID
        updates = [
            DataUpdate(
                key="s1",
                value=AssignedPolicies(
                    assigned_policies=["p1", "p2"],
                ).SerializeToString(),
            ),
            DataUpdate(
                key="s2",
                value=AssignedPolicies(
                    assigned_policies=["p2", "p3"],
                ).SerializeToString(),
            ),
        ]
        callback.process_update("stream", updates, False)

        updates = [
            DataUpdate(
                key="s1",
                value=AssignedPolicies(
                    assigned_policies=["p4", "p5"],
                ).SerializeToString(),
            ),
            DataUpdate(
                key="s2",
                value=AssignedPolicies(
                    assigned_policies=["p4", "p5"],
                ).SerializeToString(),
            ),
        ]
        callback.process_update("stream", updates, False)

        s1_policies = assignments_dict["s1"].installed_policies
        expected = 2
        self.assertEqual(len(s1_policies), expected,
                         'There should be 2 active policies for s1')
        self.assertTrue("p5" in s1_policies,
                        'Policy p5 should be active for subscriber s1')
        self.assertTrue("p4" in s1_policies,
                        'Policy p4 should be active for subscriber s1')

        s2_policies = assignments_dict["s2"].installed_policies
        expected = 2
        self.assertEqual(len(s2_policies), expected,
                         'There should be 2 active policies for s2')
        self.assertTrue("p5" in s2_policies,
                        'Policy p5 should be active for subscriber s2')
        self.assertTrue("p4" in s2_policies,
                        'Policy p4 should be active for subscriber s2')


def get_SetSessionRules_side_effect(
    called_with: List[SessionRules],
) -> Callable[[SessionRules], Void]:
    def side_effect(session_rules: SessionRules, timeout: float) -> Void:
        called_with.append(session_rules)
        return Void()
    return side_effect


class ApnRuleMappingsStreamerCallbackTest(unittest.TestCase):
    def test_Update(self):
        """
        Test the happy path where updates come in for rules, and sessiond
        accepts the SessionRules without issue.
        """

        # Expected call arguments to SetSessionRules
        allow_all_flow_list = [
            FlowDescription(
                match=FlowMatch(
                    direction=FlowMatch.Direction.Value(
                        "UPLINK"),
                ),
                action=FlowDescription.Action.Value(
                    "PERMIT"),
            ),
            FlowDescription(
                match=FlowMatch(
                    direction=FlowMatch.Direction.Value(
                        "DOWNLINK"),
                ),
                action=FlowDescription.Action.Value(
                    "PERMIT"),
            ),
        ] # type: List[FlowDescription]
        no_tracking_type = PolicyRule.TrackingType.Value("NO_TRACKING")
        expected = SessionRules(
            rules_per_subscriber=[
                RulesPerSubscriber(
                    imsi='imsi_1',
                    rule_set=[
                        RuleSet(
                            apply_subscriber_wide=False,
                            apn="apn1",
                            static_rules=[
                                StaticRuleInstall(rule_id="p1"),
                            ],
                            dynamic_rules=[
                                DynamicRuleInstall(
                                    policy_rule=PolicyRule(
                                        id="allowlist_sid-imsi_1-apn1",
                                        priority=2,
                                        flow_list=allow_all_flow_list,
                                        tracking_type=no_tracking_type,
                                    )
                                )
                            ],
                        ),
                    ]
                ),
                RulesPerSubscriber(
                    imsi='imsi_2',
                    rule_set=[
                        RuleSet(
                            apply_subscriber_wide=False,
                            apn="apn1",
                            static_rules=[
                                StaticRuleInstall(rule_id="p5"),
                            ],
                            dynamic_rules=[
                                DynamicRuleInstall(
                                    policy_rule=PolicyRule(
                                        id="allowlist_sid-imsi_2-apn1",
                                        priority=2,
                                        flow_list=allow_all_flow_list,
                                        tracking_type=no_tracking_type,
                                    )
                                )
                            ],
                        ),
                    ]
                )
            ]
        )

        # Setup the test
        apn_rules_dict = {}
        basenames_dict = {
            'bn1': ChargingRuleNameSet(RuleNames=['p5']),
            'bn2': ChargingRuleNameSet(RuleNames=['p6']),
        }
        stub = MockLocalSessionManagerStub()

        stub_call_args = [] # type: List[SessionRules]
        side_effect = get_SetSessionRules_side_effect(stub_call_args)
        stub.SetSessionRules = Mock(side_effect=side_effect)


        callback = ApnRuleMappingsStreamerCallback(
            stub,
            basenames_dict,
            apn_rules_dict,
        )

        # Construct a set of updates, keyed by subscriber ID
        updates = [
            DataUpdate(
                key="imsi_1",
                value=SubscriberPolicySet(
                    rules_per_apn=[
                        ApnPolicySet(
                            apn="apn1",
                            assigned_base_names=[],
                            assigned_policies=["p1"],
                        ),
                    ],
                ).SerializeToString(),
            ),
            DataUpdate(
                key="imsi_2",
                value=SubscriberPolicySet(
                    rules_per_apn=[
                        ApnPolicySet(
                            apn="apn1",
                            assigned_base_names=["bn1"],
                            assigned_policies=[],
                        ),
                    ],
                ).SerializeToString(),
            ),
        ]

        callback.process_update("stream", updates, False)

        # Since we used a stub which always succeeds when a RAR is made,
        # We should expect the assignments_dict to be updated
        imsi_1_policies = apn_rules_dict["imsi_1"]
        self.assertEqual(len(imsi_1_policies.rules_per_apn), 1,
                         'There should be 1 active APNs for imsi_1')
        self.assertEqual(len(stub_call_args), 1,
                         'Stub should have been called once')
        called_with = stub_call_args[0].SerializeToString()
        self.assertEqual(called_with, expected.SerializeToString(),
                         'SetSessionRules call has incorrect arguments')

        # Stream down a second update, and now IMSI_1 gets access to a new APN
        updates_2 = [
            DataUpdate(
                key="imsi_1",
                value=SubscriberPolicySet(
                    rules_per_apn=[
                        ApnPolicySet(
                            apn="apn2",
                            assigned_base_names=["bn1"],
                            assigned_policies=[],
                        ),
                    ],
                ).SerializeToString(),
            ),
            DataUpdate(
                key="imsi_2",
                value=SubscriberPolicySet(
                    global_base_names=["bn2"],
                    global_policies=[],
                    rules_per_apn=[
                        ApnPolicySet(
                            apn="apn1",
                            assigned_base_names=[],
                            assigned_policies=[],
                        ),
                    ],
                ).SerializeToString(),
            ),
        ]
        expected_2 = SessionRules(
            rules_per_subscriber=[
                RulesPerSubscriber(
                    imsi='imsi_1',
                    rule_set=[
                        RuleSet(
                            apply_subscriber_wide=False,
                            apn="apn2",
                            static_rules=[
                                StaticRuleInstall(rule_id="p5"),
                            ],
                            dynamic_rules=[
                                DynamicRuleInstall(
                                    policy_rule=PolicyRule(
                                        id="allowlist_sid-imsi_1-apn2",
                                        priority=2,
                                        flow_list=allow_all_flow_list,
                                        tracking_type=no_tracking_type,
                                    )
                                )
                            ],
                        ),
                    ]
                ),
                RulesPerSubscriber(
                    imsi='imsi_2',
                    rule_set=[
                        RuleSet(
                            apply_subscriber_wide=False,
                            apn="apn1",
                            static_rules=[
                                StaticRuleInstall(rule_id="p6"),
                            ],
                            dynamic_rules=[
                                DynamicRuleInstall(
                                    policy_rule=PolicyRule(
                                        id="allowlist_sid-imsi_2-apn1",
                                        priority=2,
                                        flow_list=allow_all_flow_list,
                                        tracking_type=no_tracking_type,
                                    )
                                )
                            ],
                        ),
                    ]
                ),
            ]
        )

        callback.process_update("stream", updates_2, False)

        imsi_1_policies = apn_rules_dict["imsi_1"]
        self.assertEqual(len(imsi_1_policies.rules_per_apn), 1,
                         'There should be 1 active APNs for imsi_1')
        self.assertEqual(len(stub_call_args), 2,
                         'Stub should have been called twice')
        called_with = stub_call_args[1].SerializeToString()
        self.assertEqual(called_with, expected_2.SerializeToString(),
                         'SetSessionRules call has incorrect arguments')
