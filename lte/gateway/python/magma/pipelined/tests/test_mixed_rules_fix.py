"""
Test for the _activate_rules_in_enforcement mixed-rules crash fix.

Self-contained: adds the gateway python path first, then stubs only the
things that genuinely aren't available (generated proto pb2 files).
All magma.pipelined.* real code loads normally from disk.

Run with:
    python test_mixed_rules_fix.py -v
"""

import sys
import types
import unittest
from unittest.mock import MagicMock

# Add the gateway python path first so real magma.pipelined code loads
_GATEWAY_PYTHON = r'c:\Users\anurag\.gemini\antigravity\scratch\magma\lte\gateway\python'
if _GATEWAY_PYTHON not in sys.path:
    sys.path.insert(0, _GATEWAY_PYTHON)


# ---------------------------------------------------------------------------
# Fake proto classes — match the real proto interface just enough
# ---------------------------------------------------------------------------

class FakeActivateFlowsResult:
    def __init__(self, policy_results=None):
        self.policy_results = list(policy_results or [])

    def __repr__(self):
        return f"ActivateFlowsResult(results={[r.rule_id for r in self.policy_results]})"


class FakeRuleModResult:
    SUCCESS = 0
    FAILURE = 2

    def __init__(self, rule_id='', result=0, version=0):
        self.rule_id = rule_id
        self.result = result
        self.version = version

    def __repr__(self):
        return f"RuleModResult(rule_id={self.rule_id}, result={self.result})"


class FakeVersionedPolicy:
    def __init__(self, rule=None, version=1):
        self.rule = rule
        self.version = version


class FakeRedirectInformation:
    DISABLED = 0
    ENABLED = 1
    URL = 2

    def __init__(self, support=0, address_type=0, server_address=''):
        self.support = support
        self.address_type = address_type
        self.server_address = server_address


class FakePolicyRule:
    def __init__(self, id='', priority=0, redirect=None):
        self.id = id
        self.priority = priority
        self.redirect = redirect if redirect is not None else FakeRedirectInformation()


class FakeIPAddress:
    IPV4 = 0

    def __init__(self, version=0, address=b''):
        self.version = version
        self.address = address


# ---------------------------------------------------------------------------
# Stub the proto modules so the import chain in rpc_servicer.py works
# ---------------------------------------------------------------------------

def _stub_module(name, **attrs):
    if name not in sys.modules:
        mod = types.ModuleType(name)
        sys.modules[name] = mod
    else:
        mod = sys.modules[name]
    for k, v in attrs.items():
        setattr(mod, k, v)
    return mod


# lte proto stubs
_stub_module('lte')
_stub_module('lte.protos')
_stub_module('lte.protos.pipelined_pb2',
    ActivateFlowsResult=FakeActivateFlowsResult,
    RuleModResult=FakeRuleModResult,
    VersionedPolicy=FakeVersionedPolicy,
    # other stuff rpc_servicer imports at top-level — stub as MagicMock
    ActivateFlowsRequest=MagicMock,
    AllTableAssignments=MagicMock,
    CauseIE=MagicMock,
    DeactivateFlowsRequest=MagicMock,
    DeactivateFlowsResult=MagicMock,
    FlowResponse=MagicMock,
    OffendingIE=MagicMock,
    PdrState=MagicMock,
    RequestOriginType=MagicMock,
    SetupFlowsResult=MagicMock,
    SetupPolicyRequest=MagicMock,
    SetupQuotaRequest=MagicMock,
    SetupUEMacRequest=MagicMock,
    SessionSet=MagicMock,
    TableAssignment=MagicMock,
    UESessionContextResponse=MagicMock,
    UESessionSet=MagicMock,
    UPFSessionContextState=MagicMock,
    VersionedPolicyID=MagicMock,
)

_pb2_grpc = _stub_module('lte.protos.pipelined_pb2_grpc')
_pb2_grpc.PipelinedServicer = object
_pb2_grpc.add_PipelinedServicer_to_server = MagicMock()

_stub_module('lte.protos.apn_pb2', AggregatedMaximumBitrate=MagicMock)
_stub_module('lte.protos.mobilityd_pb2', IPAddress=FakeIPAddress)
_stub_module('lte.protos.session_manager_pb2', RuleRecordTable=MagicMock)
_stub_module('lte.protos.policydb_pb2',
    PolicyRule=FakePolicyRule,
    RedirectInformation=FakeRedirectInformation,
)
_stub_module('lte.protos.subscriberdb_pb2', SubscriberID=MagicMock)

# magma.common
_stub_module('magma.common.sentry', EXCLUDE_FROM_ERROR_MONITORING={})

# magma.pipelined sub-deps that rpc_servicer.py imports at top level
_metrics = _stub_module('magma.pipelined.metrics')
_metrics.ENFORCEMENT_RULE_INSTALL_FAIL = MagicMock()
_metrics.ENFORCEMENT_RULE_INSTALL_FAIL.labels = MagicMock(return_value=MagicMock())
_metrics.ENFORCEMENT_STATS_RULE_INSTALL_FAIL = MagicMock()
_metrics.ENFORCEMENT_STATS_RULE_INSTALL_FAIL.labels = MagicMock(return_value=MagicMock())

_stub_module('magma.pipelined.ng_manager')
_stub_module('magma.pipelined.ng_manager.session_state_manager_util',
    PDRRuleEntry=MagicMock)

_stub_module('magma.pipelined.policy_converters',
    convert_ip_str_to_ip_proto=MagicMock(),
    convert_ipv4_str_to_ip_proto=MagicMock(),
    convert_ipv6_bytes_to_ip_proto=MagicMock(),
)

_stub_module('magma.pipelined.imsi', encode_imsi=lambda x: x)
_stub_module('magma.pipelined.ipv6_prefix_store',
    get_ipv6_interface_id=MagicMock(),
    get_ipv6_prefix=MagicMock(),
)

# App module stubs
_app_map = {
    'check_quota': 'CheckQuotaController',
    'classifier': 'Classifier',
    'dpi': 'DPIController',
    'enforcement': 'EnforcementController',
    'enforcement_stats': 'EnforcementStatsController',
    'ipfix': 'IPFIXController',
    'ng_services': 'NGServiceController',
    'tunnel_learn': 'TunnelLearnController',
    'ue_mac': 'UEMacAddressController',
    'vlan_learn': 'VlanLearnController',
}
for mod_suffix, cls_name in _app_map.items():
    _stub_module(
        f'magma.pipelined.app.{mod_suffix}',
        **{cls_name: type(cls_name, (), {'APP_NAME': cls_name})},
    )

# ---------------------------------------------------------------------------
# Now import the real thing
# ---------------------------------------------------------------------------
from magma.pipelined.rpc_servicer import PipelinedRpcServicer  # noqa: E402


# ---------------------------------------------------------------------------
# Test helpers
# ---------------------------------------------------------------------------

def _make_servicer():
    loop = MagicMock()
    cfg = MagicMock()
    cfg.get.return_value = 5
    svc_mgr = MagicMock()
    svc_mgr.is_app_enabled.return_value = True
    return PipelinedRpcServicer(
        loop,
        MagicMock(), MagicMock(), MagicMock(),  # gy, enforcer, enf_stats
        MagicMock(), MagicMock(), MagicMock(),  # dpi, ue_mac, check_quota
        MagicMock(), MagicMock(), MagicMock(),  # ipfix, vlan_learn, tunnel_learn
        MagicMock(), MagicMock(), MagicMock(), MagicMock(),  # classifier, ingress, middle, egress
        MagicMock(),  # ng_servicer
        cfg, svc_mgr,
    )


def _static(rule_id, version=1):
    return FakeVersionedPolicy(
        rule=FakePolicyRule(
            id=rule_id, priority=100,
            redirect=FakeRedirectInformation(support=FakeRedirectInformation.DISABLED),
        ),
        version=version,
    )


def _redirect(rule_id, version=1):
    return FakeVersionedPolicy(
        rule=FakePolicyRule(
            id=rule_id, priority=100,
            redirect=FakeRedirectInformation(
                support=FakeRedirectInformation.ENABLED,
                address_type=FakeRedirectInformation.URL,
                server_address='http://example.com',
            ),
        ),
        version=version,
    )


def _result(rule_id):
    return FakeActivateFlowsResult(
        policy_results=[FakeRuleModResult(rule_id=rule_id, result=FakeRuleModResult.SUCCESS)],
    )


IP = FakeIPAddress(version=FakeIPAddress.IPV4, address=b'1.2.3.4')


# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------

class TestActivateRulesInEnforcementMixedFix(unittest.TestCase):

    def setUp(self):
        self.srv = _make_servicer()
        self.enforcer = self.srv._enforcer_app

    def _call(self, policies):
        return self.srv._activate_rules_in_enforcement(
            'imsi01', b'msisdn', 0, IP, None, policies, 0,
        )

    # --- edge cases ---

    def test_empty_policies_no_call(self):
        """Empty list -> activate_rules never called, empty result."""
        result = self._call([])
        self.enforcer.activate_rules.assert_not_called()
        self.assertEqual(result.policy_results, [])

    # --- single-type inputs (regression guard) ---

    def test_only_static_rules_one_call(self):
        """All static -> single activate_rules call, both policies in it."""
        p1, p2 = _static('s1'), _static('s2')
        self.enforcer.activate_rules.return_value = FakeActivateFlowsResult(
            policy_results=[
                FakeRuleModResult(rule_id='s1'),
                FakeRuleModResult(rule_id='s2'),
            ],
        )

        result = self._call([p1, p2])

        self.assertEqual(self.enforcer.activate_rules.call_count, 1)
        sent = self.enforcer.activate_rules.call_args.args[5]
        self.assertEqual([p.rule.id for p in sent], ['s1', 's2'])
        self.assertEqual(len(result.policy_results), 2)

    def test_only_redirect_rules_one_call(self):
        """All redirect -> single activate_rules call."""
        p1 = _redirect('r1')
        self.enforcer.activate_rules.return_value = _result('r1')

        result = self._call([p1])

        self.assertEqual(self.enforcer.activate_rules.call_count, 1)
        self.assertEqual(len(result.policy_results), 1)

    # --- the actual bug scenario ---

    def test_mixed_rules_two_calls_no_crash(self):
        """Mixed static + redirect -> two separate activate_rules calls."""
        sp, rp = _static('static_rule'), _redirect('redir_rule')
        self.enforcer.activate_rules.side_effect = [
            _result('static_rule'),
            _result('redir_rule'),
        ]

        result = self._call([sp, rp])

        # two calls — not one mixed call
        self.assertEqual(self.enforcer.activate_rules.call_count, 2)

        first_batch = self.enforcer.activate_rules.call_args_list[0].args[5]
        second_batch = self.enforcer.activate_rules.call_args_list[1].args[5]

        self.assertEqual(first_batch[0].rule.id, 'static_rule')
        self.assertEqual(second_batch[0].rule.id, 'redir_rule')

        result_ids = {r.rule_id for r in result.policy_results}
        self.assertEqual(result_ids, {'static_rule', 'redir_rule'})

    def test_static_always_processed_before_redirect(self):
        """Static batch goes first regardless of input order."""
        sp, rp = _static('s'), _redirect('r')
        self.enforcer.activate_rules.side_effect = [_result('s'), _result('r')]

        # give redirect first in the list
        self._call([rp, sp])

        first_sent = self.enforcer.activate_rules.call_args_list[0].args[5]
        self.assertEqual(first_sent[0].rule.id, 's',
                         "static should always go first")

    def test_multiple_static_multiple_redirect(self):
        """3 static + 2 redirect -> two calls, 5 results total."""
        statics = [_static(f's{i}') for i in range(3)]
        redirs = [_redirect(f'r{i}') for i in range(2)]

        self.enforcer.activate_rules.side_effect = [
            FakeActivateFlowsResult(
                policy_results=[FakeRuleModResult(rule_id=f's{i}') for i in range(3)],
            ),
            FakeActivateFlowsResult(
                policy_results=[FakeRuleModResult(rule_id=f'r{i}') for i in range(2)],
            ),
        ]

        result = self._call(statics + redirs)

        self.assertEqual(self.enforcer.activate_rules.call_count, 2)
        self.assertEqual(len(result.policy_results), 5)

    def test_no_phantom_entries_in_merged_result(self):
        """Every entry in merged result must have a real rule_id."""
        sp, rp = _static('s'), _redirect('r')
        self.enforcer.activate_rules.side_effect = [_result('s'), _result('r')]

        result = self._call([sp, rp])

        for r in result.policy_results:
            self.assertTrue(r.rule_id, f"phantom entry with empty rule_id: {r}")

    def test_failure_in_static_still_returns_redirect_results(self):
        """Even if static call succeeds partially, redirect results still merged."""
        sp, rp = _static('s'), _redirect('r')
        self.enforcer.activate_rules.side_effect = [
            FakeActivateFlowsResult(
                policy_results=[FakeRuleModResult(rule_id='s', result=FakeRuleModResult.FAILURE)],
            ),
            _result('r'),
        ]

        result = self._call([sp, rp])

        self.assertEqual(len(result.policy_results), 2)
        result_ids = {r.rule_id for r in result.policy_results}
        self.assertIn('s', result_ids)
        self.assertIn('r', result_ids)


if __name__ == '__main__':
    unittest.main(verbosity=2)
