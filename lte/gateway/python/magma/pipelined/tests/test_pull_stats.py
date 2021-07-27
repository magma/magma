import unittest
import warnings
from concurrent.futures import Future
from unittest.mock import MagicMock

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.pipelined_pb2 import VersionedPolicy
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow import flows
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.pipelined_test_util import (
    SnapshotVerifier,
    create_service_manager,
    fake_controller_setup,
    start_ryu_app_thread,
    stop_ryu_app_thread,
)
from scapy.all import IP


class PullStatsTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    DEFAULT_DROP_FLOW_NAME = '(ノಠ益ಠ)ノ彡┻━┻'

    def setUp(self):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.

        Mocks the redis policy_dictionary of enforcement_controller.
        Mocks the loop for testing EnforcementStatsController
        """
        super(PullStatsTest, self).setUpClass()
        warnings.simplefilter('ignore')
        self.service_manager = create_service_manager([PipelineD.ENFORCEMENT])
        self._main_tbl_num = self.service_manager.get_table_num(
            EnforcementController.APP_NAME,
        )

        enforcement_controller_reference = Future()
        testing_controller_reference = Future()
        enf_stat_ref = Future()

        """
        Enforcement_stats reports data by using loop.call_soon_threadsafe, but
        as we don't have an eventloop in testing, just directly call the stats
        handling function

        Here is how the mocked function is used in EnforcementStatsController:
        self.loop.call_soon_threadsafe(self._handle_flow_stats, ev.msg.body)
        """
        def mock_thread_safe(cmd, body):
            cmd(body)
        loop_mock = MagicMock()
        loop_mock.call_soon_threadsafe = mock_thread_safe

        test_setup = TestSetup(
            apps=[
                PipelinedController.Enforcement,
                PipelinedController.Enforcement_stats,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.Enforcement:
                    enforcement_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.Enforcement_stats:
                    enf_stat_ref,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': self.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'enforcement': {
                    'poll_interval': 2,
                    'default_drop_flow_name': self.DEFAULT_DROP_FLOW_NAME,
                    'periodic_stats_reporting': False,
                },
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
                'redis_enabled': False,
                'setup_type': 'LTE',
            },
            mconfig=PipelineD(),
            loop=loop_mock,
            service_manager=self.service_manager,
            integ_test=False,
            rpc_stubs={'sessiond': MagicMock()},
        )

        BridgeTools.create_bridge(self.BRIDGE, self.IFACE)

        self.thread = start_ryu_app_thread(test_setup)

        self.enforcement_stats_controller = enf_stat_ref.result()
        self._scratch_tbl_num = self.enforcement_stats_controller.tbl_num
        self.enforcement_controller = enforcement_controller_reference.result()
        self.testing_controller = testing_controller_reference.result()

        self.enforcement_stats_controller._report_usage = MagicMock()
        self.enforcement_controller._redirect_manager._save_redirect_entry = \
            MagicMock()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

    def test_poll(self):
        """
        Unit test to help verify stats polling using cookie and cookie_mask
        """
        fake_controller_setup(
            self.enforcement_controller,
            self.enforcement_stats_controller,
        )
        imsi = 'IMSI001010000000013'
        sub_ip = '192.168.128.74'

        flow_list = [
            FlowDescription(
                match=FlowMatch(
                    ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.0/25'),
                    direction=FlowMatch.UPLINK,
                ),
                action=FlowDescription.PERMIT,
            ),
        ]
        policy = VersionedPolicy(
            rule=PolicyRule(id='rule1', priority=3, flow_list=flow_list),
            version=1,
        )
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, convert_ipv4_str_to_ip_proto(sub_ip), 'rule1', 1,
        )

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller,
        ).add_policy(policy)

        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
        )
        with sub_context, snapshot_verifier:
            rule_map = self.enforcement_stats_controller.get_stats()
            if (rule_map.records[0].rule_id == self.DEFAULT_DROP_FLOW_NAME):
                rule_record = rule_map.records[1]
            else:
                rule_record = rule_map.records[0]
            self.assertEqual(rule_record.sid, imsi)
            self.assertEqual(rule_record.rule_id, "rule1")
            self.assertEqual(rule_record.bytes_tx, 0)
            self.assertEqual(rule_record.bytes_rx, 0)
            rule_map_cookie = self.enforcement_stats_controller.get_stats(1, 0)
            if (rule_map_cookie.records[0].rule_id == self.DEFAULT_DROP_FLOW_NAME):
                rule_record_cookie = rule_map_cookie.records[1]
            else:
                rule_record_cookie = rule_map_cookie.records[0]
            self.assertEqual(rule_record_cookie.sid, imsi)
            self.assertEqual(rule_record_cookie.rule_id, "rule1")
            self.assertEqual(rule_record_cookie.bytes_tx, 0)
            self.assertEqual(rule_record_cookie.bytes_rx, 0)
            rule_map_cookie_and_mask = self.enforcement_stats_controller.get_stats(1, 1)
            rule_record_cookie_and_mask = rule_map_cookie_and_mask.records[0]
            self.assertEqual(rule_record_cookie_and_mask.sid, imsi)
            self.assertEqual(rule_record_cookie_and_mask.rule_id, "rule1")
            self.assertEqual(rule_record_cookie_and_mask.bytes_tx, 0)
            self.assertEqual(rule_record_cookie_and_mask.bytes_rx, 0)


if __name__ == "__main__":
    unittest.main()
