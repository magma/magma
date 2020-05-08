"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import warnings
import unittest
from concurrent.futures import Future
from unittest.mock import MagicMock

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule, \
    RedirectInformation
from magma.pipelined.app.gy import GYController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import PipelinedController, \
    TestSetup
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.app.table_isolation import RyuDirectTableIsolator, \
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.pipelined_test_util import SnapshotVerifier, \
    create_service_manager, start_ryu_app_thread, fake_controller_setup, \
    stop_ryu_app_thread, wait_after_send, SnapshotVerifier


class GYTableTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"

    @classmethod
    @unittest.mock.patch('netifaces.ifaddresses',
                return_value=[[{'addr': '00:11:22:33:44:55'}]])
    @unittest.mock.patch('netifaces.AF_LINK', 0)
    def setUpClass(cls, *_):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of gy_controller
        """
        super(GYTableTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls._static_rule_dict = {}
        cls.service_manager = create_service_manager([PipelineD.ENFORCEMENT])
        cls._tbl_num = cls.service_manager.get_table_num(
            GYController.APP_NAME)

        gy_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.GY,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.GY:
                    gy_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'internal_ip_subnet': '192.168.0.0/16',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'enable_queue_pgm': False,
                'clean_restart': True,
            },
            mconfig=PipelineD(
                relay_enabled=True
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)

        cls.gy_controller = gy_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

        cls.gy_controller._policy_dict = cls._static_rule_dict
        cls.gy_controller._redirect_manager._save_redirect_entry = MagicMock()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_subscriber_policy(self):
        """
        Add policy to subscriber, send 4096 packets

        Assert:
            Packets are properly matched with the 'simple_match' policy
            Send /20 (4096) packets, match /16 (256) packets
        """
        fake_controller_setup(self.gy_controller)
        imsi = 'IMSI010000000088888'
        sub_ip = '192.168.128.74'
        redirect_ips = ["185.128.101.5", "185.128.121.4"]
        self.gy_controller._redirect_manager._dns_cache.get(
            "about.sha.ddih.org", lambda: redirect_ips, max_age=42
        )
        flow_list = [FlowDescription(match=FlowMatch())]
        policy = PolicyRule(
            id='redir_test', priority=3, flow_list=flow_list,
            redirect=RedirectInformation(
                support=1,
                address_type=2,
                server_address="http://about.sha.ddih.org/"
            )
        )

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.gy_controller, self._tbl_num
        ).add_dynamic_rule(policy)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder()\
            .set_ip_layer('45.10.0.0/20', sub_ip)\
            .set_ether_layer(self.MAC_DEST, "01:20:10:20:aa:bb")\
            .build()

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             include_stats=False)

        with isolator, sub_context, flow_verifier, snapshot_verifier:
            pkt_sender.send(packet)



if __name__ == "__main__":
    unittest.main()
