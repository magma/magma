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
from unittest.mock import MagicMock

from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import (
    ActivateFlowsRequest,
    ActivateFlowsResult,
    DeactivateFlowsRequest,
    RequestOriginType,
    RuleModResult,
    SetupPolicyRequest,
    VersionedPolicy,
    VersionedPolicyID,
)
from lte.protos.policydb_pb2 import PolicyRule, RedirectInformation
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.pipelined.rpc_servicer import PipelinedRpcServicer


class RPCServicerTest(unittest.TestCase):
    def setUp(self):
        def call_soon_threadsafe(func, arg, fut=None):
            if fut:
                res = func(arg, fut)
            else:
                res = func(arg)
            return res
        self._loop = MagicMock()
        self._loop.call_soon_threadsafe = MagicMock(side_effect=call_soon_threadsafe)
        self._gy_app = MagicMock()
        self._enforcer_app = MagicMock()
        self._enforcement_stats = MagicMock()
        self._dpi_app = MagicMock()
        self._ue_mac_app = MagicMock()
        self._check_quota_app = MagicMock()
        self._ipfix_app = MagicMock()
        self._vlan_learn_app = MagicMock()
        self._tunnel_learn_app = MagicMock()
        self._classifier_app = MagicMock()
        self._ingress_app = MagicMock()
        self._middle_app = MagicMock()
        self._egress_app = MagicMock()
        self._ng_servicer_app = MagicMock()
        self._service_config = MagicMock()
        self._service_manager = MagicMock()
        self._service_manager.is_app_enabled.side_effect = lambda x: True

        for controller in [
            self._gy_app, self._enforcer_app,
            self._enforcement_stats,
        ]:
            controller.check_setup_request_epoch.side_effect = lambda x: None
            controller.is_controller_ready = lambda: True

        self.pipelined_srv = PipelinedRpcServicer(
            self._loop,
            self._gy_app,
            self._enforcer_app,
            self._enforcement_stats,
            self._dpi_app,
            self._ue_mac_app,
            self._check_quota_app,
            self._ipfix_app,
            self._vlan_learn_app,
            self._tunnel_learn_app,
            self._classifier_app,
            self._ingress_app,
            self._middle_app,
            self._egress_app,
            self._ng_servicer_app,
            self._service_config,
            self._service_manager,
        )

    def test_setup_flows_req(self):
        gx_req1 = ActivateFlowsRequest()
        gx_req2 = ActivateFlowsRequest()
        gy_req = ActivateFlowsRequest(
            request_origin=RequestOriginType(type=RequestOriginType.GY),
        )
        setup_req = SetupPolicyRequest(requests=[gx_req1, gx_req2, gy_req])

        self.pipelined_srv.SetupPolicyFlows(setup_req, MagicMock())
        self._enforcer_app.handle_restart.assert_called_with([gx_req1, gx_req2])
        self._enforcement_stats.handle_restart.assert_called_with([gx_req1, gx_req2])
        self._gy_app.handle_restart.assert_called_with([gy_req])

    def test_activate_flows_req(self):
        rule = PolicyRule(id="rule1", priority=100, flow_list=[])
        policies = [VersionedPolicy(rule=rule, version=1)]
        req = ActivateFlowsRequest(
            sid=SubscriberID(id="imsi12345"),
            ip_addr="1.2.3.4",
            msisdn=b'magma',
            uplink_tunnel=0x1,
            downlink_tunnel=0x2,
            policies=policies,
        )
        ip_addr = IPAddress(
            version=IPAddress.IPV4,
            address=req.ip_addr.encode('utf-8'),
        )

        self.pipelined_srv.ActivateFlows(req, MagicMock())
        # Not using assert_called_with because protos comparison

        assert self._enforcement_stats.activate_rules.call_args.args[0] == req.sid.id
        assert self._enforcement_stats.activate_rules.call_args.args[1] == req.msisdn
        assert self._enforcement_stats.activate_rules.call_args.args[2] == req.uplink_tunnel
        assert self._enforcement_stats.activate_rules.call_args.args[3].version == ip_addr.version
        assert self._enforcement_stats.activate_rules.call_args.args[3].address == ip_addr.address
        assert self._enforcement_stats.activate_rules.call_args.args[4] == req.apn_ambr
        assert self._enforcement_stats.activate_rules.call_args.args[5][0].version == policies[0].version
        assert self._enforcement_stats.activate_rules.call_args.args[6] == req.shard_id
        assert self._enforcement_stats.activate_rules.call_args.args[7] == 0

        assert self._enforcer_app.activate_rules.call_args.args[0] == req.sid.id
        assert self._enforcer_app.activate_rules.call_args.args[1] == req.msisdn
        assert self._enforcer_app.activate_rules.call_args.args[2] == req.uplink_tunnel
        assert self._enforcer_app.activate_rules.call_args.args[3].version == ip_addr.version
        assert self._enforcer_app.activate_rules.call_args.args[3].address == ip_addr.address
        assert self._enforcer_app.activate_rules.call_args.args[4] == req.apn_ambr
        assert self._enforcer_app.activate_rules.call_args.args[5][0].version == policies[0].version
        assert self._enforcer_app.activate_rules.call_args.args[6] == req.shard_id
        assert self._enforcer_app.activate_rules.call_args.args[7] == 0

    def test_deactivate_flows_req(self):
        policies = [VersionedPolicyID(rule_id="rule1", version=1)]
        req = DeactivateFlowsRequest(
            sid=SubscriberID(id="imsi12345"),
            ip_addr="1.2.3.4",
            uplink_tunnel=0x1,
            downlink_tunnel=0x2,
            policies=policies,
        )
        ip_addr = IPAddress(
            version=IPAddress.IPV4,
            address=req.ip_addr.encode('utf-8'),
        )

        self.pipelined_srv.DeactivateFlows(req, MagicMock())
        assert self._enforcer_app.deactivate_rules.call_args.args[0] == req.sid.id
        assert self._enforcer_app.deactivate_rules.call_args.args[1].version == ip_addr.version
        assert self._enforcer_app.deactivate_rules.call_args.args[1].address == ip_addr.address
        assert self._enforcer_app.deactivate_rules.call_args.args[2] == ["rule1"]

    # -----------------------------------------------------------------------
    # Tests for _activate_rules_in_enforcement mixed-rule fix
    # -----------------------------------------------------------------------

    def _make_static_policy(self, rule_id, version=1):
        rule = PolicyRule(
            id=rule_id, priority=100,
            redirect=RedirectInformation(support=RedirectInformation.DISABLED),
        )
        return VersionedPolicy(rule=rule, version=version)

    def _make_redirect_policy(self, rule_id, version=1):
        rule = PolicyRule(
            id=rule_id, priority=100,
            redirect=RedirectInformation(
                support=RedirectInformation.ENABLED,
                address_type=RedirectInformation.URL,
                server_address="https://example.com",
            ),
        )
        return VersionedPolicy(rule=rule, version=version)

    def _fake_activate_rules_result(self, rule_id):
        return ActivateFlowsResult(
            policy_results=[
                RuleModResult(rule_id=rule_id, result=RuleModResult.SUCCESS),
            ],
        )

    def test_activate_rules_in_enforcement_only_static(self):
        # Only static rules — should call activate_rules exactly once
        ip = IPAddress(version=IPAddress.IPV4, address=b'1.2.3.4')
        p1 = self._make_static_policy('static1')
        p2 = self._make_static_policy('static2')

        self._enforcer_app.activate_rules.side_effect = [
            self._fake_activate_rules_result('static1'),
            self._fake_activate_rules_result('static2'),
        ]

        result = self.pipelined_srv._activate_rules_in_enforcement(
            'imsi01', b'msisdn', 0, ip, None, [p1, p2], 0,
        )

        # called once — both static go in one shot
        self.assertEqual(self._enforcer_app.activate_rules.call_count, 1)
        passed_policies = self._enforcer_app.activate_rules.call_args.args[5]
        self.assertEqual(len(passed_policies), 2)
        self.assertEqual(len(result.policy_results), 1)

    def test_activate_rules_in_enforcement_only_redirect(self):
        # Only redirect rules — should call activate_rules exactly once
        ip = IPAddress(version=IPAddress.IPV4, address=b'1.2.3.4')
        p1 = self._make_redirect_policy('redir1')

        self._enforcer_app.activate_rules.return_value = \
            self._fake_activate_rules_result('redir1')

        result = self.pipelined_srv._activate_rules_in_enforcement(
            'imsi01', b'msisdn', 0, ip, None, [p1], 0,
        )

        self.assertEqual(self._enforcer_app.activate_rules.call_count, 1)
        passed_policies = self._enforcer_app.activate_rules.call_args.args[5]
        self.assertEqual(len(passed_policies), 1)
        self.assertEqual(passed_policies[0].rule.id, 'redir1')
        self.assertEqual(len(result.policy_results), 1)

    def test_activate_rules_in_enforcement_mixed_no_crash(self):
        # The actual bug: mixed static + redirect — should call activate_rules
        # twice (once per group) and merge results cleanly
        ip = IPAddress(version=IPAddress.IPV4, address=b'1.2.3.4')
        static_p = self._make_static_policy('static_rule')
        redir_p = self._make_redirect_policy('redir_rule')

        self._enforcer_app.activate_rules.side_effect = [
            self._fake_activate_rules_result('static_rule'),
            self._fake_activate_rules_result('redir_rule'),
        ]

        result = self.pipelined_srv._activate_rules_in_enforcement(
            'imsi01', b'msisdn', 0, ip, None, [static_p, redir_p], 0,
        )

        # two calls — one for each group
        self.assertEqual(self._enforcer_app.activate_rules.call_count, 2)

        first_call_policies = self._enforcer_app.activate_rules.call_args_list[0].args[5]
        second_call_policies = self._enforcer_app.activate_rules.call_args_list[1].args[5]

        # first batch is static, second is redirect
        self.assertEqual(len(first_call_policies), 1)
        self.assertEqual(first_call_policies[0].rule.id, 'static_rule')
        self.assertEqual(len(second_call_policies), 1)
        self.assertEqual(second_call_policies[0].rule.id, 'redir_rule')

        # both results merged
        self.assertEqual(len(result.policy_results), 2)
        result_ids = {r.rule_id for r in result.policy_results}
        self.assertIn('static_rule', result_ids)
        self.assertIn('redir_rule', result_ids)

    def test_activate_rules_in_enforcement_empty_policies(self):
        # edge case: empty list — should not call activate_rules at all
        ip = IPAddress(version=IPAddress.IPV4, address=b'1.2.3.4')

        result = self.pipelined_srv._activate_rules_in_enforcement(
            'imsi01', b'msisdn', 0, ip, None, [], 0,
        )

        self._enforcer_app.activate_rules.assert_not_called()
        self.assertEqual(len(result.policy_results), 0)


if __name__ == "__main__":
    unittest.main()
