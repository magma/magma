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
import unittest.mock
from concurrent import futures

import grpc
from lte.protos.policydb_pb2 import (
    ChargingRuleNameSet,
    DisableStaticRuleRequest,
    EnableStaticRuleRequest,
)
from lte.protos.policydb_pb2_grpc import PolicyDBStub
from lte.protos.session_manager_pb2 import (
    PolicyReAuthAnswer,
    PolicyReAuthRequest,
    ReAuthResult,
)
from magma.policydb.reauth_handler import ReAuthHandler
from magma.policydb.servicers.policy_servicer import PolicyRpcServicer
from orc8r.protos.common_pb2 import Void


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


class MockPolicyAssignmentControllerStub:
    """ Always succeeds by not raising an error """

    def __init__(self):
        pass

    def EnableStaticRules(self, _: EnableStaticRuleRequest) -> Void:
        return Void()

    def DisableStaticRules(self, _: DisableStaticRuleRequest) -> Void:
        return Void()


class MockPolicyAssignmentControllerStub2:
    """ Always fails """

    def __init__(self):
        pass

    def EnableStaticRules(self, _: EnableStaticRuleRequest) -> Void:
        raise grpc.RpcError()

    def DisableStaticRules(self, _: DisableStaticRuleRequest) -> Void:
        raise grpc.RpcError()


class PolicyRpcServicerTest(unittest.TestCase):
    def test_EnableStaticRules(self):
        """
        Check the happy path where everything succeeds.
        """
        rules_by_sid = {}
        rules_by_basename = {
            "bn1": ChargingRuleNameSet(
                RuleNames=["p4", "p5"],
            ),
        }
        reauth_handler = ReAuthHandler(
            rules_by_sid,
            MockSessionProxyResponderStub(),
        )

        servicer = PolicyRpcServicer(
            reauth_handler,
            rules_by_basename,
            MockPolicyAssignmentControllerStub(),
        )

        # Bind the rpc server to a free port
        thread_pool = futures.ThreadPoolExecutor(max_workers=10)
        rpc_server = grpc.server(thread_pool)
        port = rpc_server.add_insecure_port('0.0.0.0:0')

        # Create a mock "mconfig" for the servicer to use
        mconfig = unittest.mock.Mock()
        mconfig.ip_block = None

        # Add the servicer
        servicer.add_to_server(rpc_server)
        rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))
        stub = PolicyDBStub(channel)
        rules_by_basename["bn1"] = ChargingRuleNameSet(
            RuleNames=["p4", "p5"],
        )
        req = EnableStaticRuleRequest(
            imsi="s1",
            rule_ids=["p1", "p2", "p3"],
            base_names=["bn1"],
        )
        stub.EnableStaticRules(req)
        self.assertEqual(
            len(rules_by_sid["s1"].installed_policies), 5,
            'After a successful update, Redis should be tracking '
            '5 active rules.',
        )

    def test_FailOrc8r(self):
        """ Check that nothing is updated if orc8r is unreachable """
        rules_by_sid = {}
        rules_by_basename = {
            "bn1": ChargingRuleNameSet(
                RuleNames=["p4", "p5"],
            ),
        }
        reauth_handler = ReAuthHandler(
            rules_by_sid,
            MockSessionProxyResponderStub(),
        )

        servicer = PolicyRpcServicer(
            reauth_handler,
            rules_by_basename,
            MockPolicyAssignmentControllerStub2(),
        )

        # Bind the rpc server to a free port
        thread_pool = futures.ThreadPoolExecutor(max_workers=10)
        rpc_server = grpc.server(thread_pool)
        port = rpc_server.add_insecure_port('0.0.0.0:0')

        # Create a mock "mconfig" for the servicer to use
        mconfig = unittest.mock.Mock()
        mconfig.ip_block = None

        # Add the servicer
        servicer.add_to_server(rpc_server)
        rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))
        stub = PolicyDBStub(channel)
        req = EnableStaticRuleRequest(
            imsi="s1",
            rule_ids=["p1", "p2", "p3"],
            base_names=["bn1"],
        )
        with self.assertRaises(grpc.RpcError):
            stub.EnableStaticRules(req)

        self.assertFalse(
            "s1" in rules_by_sid,
            "There should be no installed policies for s1",
        )
