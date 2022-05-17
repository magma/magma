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
    DeactivateFlowsRequest,
    RequestOriginType,
    SetupPolicyRequest,
    VersionedPolicy,
    VersionedPolicyID,
)
from lte.protos.policydb_pb2 import PolicyRule
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


if __name__ == "__main__":
    unittest.main()
