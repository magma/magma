"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
from lte.protos.policydb_pb2 import AssignedPolicies, ChargingRuleNameSet
from lte.protos.session_manager_pb2 import PolicyReAuthRequest, \
    PolicyReAuthAnswer, ReAuthResult
from magma.policydb.streamer_callback import RuleMappingsStreamerCallback
from magma.policydb.reauth_handler import ReAuthHandler
from orc8r.protos.streamer_pb2 import DataUpdate


class MockSessionProxyResponderStub1:
    """
    This Mock SessionProxyResponderStub will always respond with a success to
    a received RAR
    """
    def __init__(self):
        pass

    def PolicyReAuth(self, _: PolicyReAuthRequest) -> PolicyReAuthAnswer:
        return PolicyReAuthAnswer(
            result=ReAuthResult.Value('UPDATE_INITIATED')
        )


class MockSessionProxyResponderStub2:
    """
    This Mock SessionProxyResponderStub will always respond with a failure to
    a received RAR
    """
    def __init__(self):
        pass

    def PolicyReAuth(self, _: PolicyReAuthRequest) -> PolicyReAuthAnswer:
        return PolicyReAuthAnswer(
            result=ReAuthResult.Value('OTHER_FAILURE')
        )


class MockSessionProxyResponderStub3:
    """
    This Mock SessionProxyResponderStub will always fail to install rule p2
    """
    def __init__(self):
        pass

    def PolicyReAuth(self, _: PolicyReAuthRequest) -> PolicyReAuthAnswer:
        return PolicyReAuthAnswer(
            result=ReAuthResult.Value('UPDATE_INITIATED'),
            failed_rules={
                "p2": PolicyReAuthAnswer.FailureCode.Value("UNKNOWN_RULE_NAME"),
            },
        )


class RuleMappingsStreamerCallbackTest(unittest.TestCase):
    def test_SuccessfulUpdate(self):
        """
        Test the happy path where updates come in for added rules, and sessiond
        accepts the RAR without issue.
        """
        assignments_dict = {}
        basenames_dict = {
            'bn1': ChargingRuleNameSet(RuleNames=['p5']),
            'bn2': ChargingRuleNameSet(RuleNames=['p6']),
        }
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub1()),
            basenames_dict,
            assignments_dict,
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

    def test_FailedUpdate(self):
        """
        Test when sessiond answers to the RAR with a failure for any re-auth.
        """
        assignments_dict = {}
        basenames_dict = {}
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub2()),
            basenames_dict,
            assignments_dict,
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
        basenames_dict = {}
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub3()),
            basenames_dict,
            assignments_dict,
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
        basenames_dict = {}
        callback = RuleMappingsStreamerCallback(
            ReAuthHandler(assignments_dict, MockSessionProxyResponderStub3()),
            basenames_dict,
            assignments_dict,
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
