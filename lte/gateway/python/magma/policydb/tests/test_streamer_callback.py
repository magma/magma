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
from typing import Callable, List
from unittest.mock import Mock

from lte.protos.policydb_pb2 import (
    ApnPolicySet,
    ChargingRuleNameSet,
    FlowDescription,
    FlowMatch,
    PolicyRule,
    SubscriberPolicySet,
)
from lte.protos.session_manager_pb2 import (
    DynamicRuleInstall,
    RuleSet,
    RulesPerSubscriber,
    SessionRules,
    StaticRuleInstall,
)
from magma.policydb.streamer_callback import ApnRuleMappingsStreamerCallback
from magma.policydb.tests.mock_stubs import MockLocalSessionManagerStub
from orc8r.protos.common_pb2 import Void
from orc8r.protos.streamer_pb2 import DataUpdate


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
                        "UPLINK",
                    ),
                ),
                action=FlowDescription.Action.Value(
                    "PERMIT",
                ),
            ),
            FlowDescription(
                match=FlowMatch(
                    direction=FlowMatch.Direction.Value(
                        "DOWNLINK",
                    ),
                ),
                action=FlowDescription.Action.Value(
                    "PERMIT",
                ),
            ),
        ]  # type: List[FlowDescription]
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
                                    ),
                                ),
                            ],
                        ),
                    ],
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
                                    ),
                                ),
                            ],
                        ),
                    ],
                ),
            ],
        )

        # Setup the test
        apn_rules_dict = {}
        basenames_dict = {
            'bn1': ChargingRuleNameSet(RuleNames=['p5']),
            'bn2': ChargingRuleNameSet(RuleNames=['p6']),
        }
        stub = MockLocalSessionManagerStub()

        stub_call_args = []  # type: List[SessionRules]
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
        self.assertEqual(
            len(imsi_1_policies.rules_per_apn), 1,
            'There should be 1 active APNs for imsi_1',
        )
        self.assertEqual(
            len(stub_call_args), 1,
            'Stub should have been called once',
        )
        called_with = stub_call_args[0].SerializeToString()
        self.assertEqual(
            called_with, expected.SerializeToString(),
            'SetSessionRules call has incorrect arguments',
        )

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
                                    ),
                                ),
                            ],
                        ),
                    ],
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
                                    ),
                                ),
                            ],
                        ),
                    ],
                ),
            ],
        )

        callback.process_update("stream", updates_2, False)

        imsi_1_policies = apn_rules_dict["imsi_1"]
        self.assertEqual(
            len(imsi_1_policies.rules_per_apn), 1,
            'There should be 1 active APNs for imsi_1',
        )
        self.assertEqual(
            len(stub_call_args), 2,
            'Stub should have been called twice',
        )
        called_with = stub_call_args[1].SerializeToString()
        self.assertEqual(
            called_with, expected_2.SerializeToString(),
            'SetSessionRules call has incorrect arguments',
        )
