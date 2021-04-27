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
import warnings
from concurrent.futures import Future
from unittest.mock import MagicMock

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.pipelined_pb2 import VersionedPolicy
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    PolicyRule,
    RedirectInformation,
)
from magma.pipelined.app.gy import GYController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import (
    convert_ipv4_str_to_ip_proto,
    flow_match_to_magma_match,
)
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery as FlowQuery
from magma.pipelined.tests.app.packet_builder import (
    IPPacketBuilder,
    TCPPacketBuilder,
)
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.app.table_isolation import (
    RyuDirectTableIsolator,
    RyuForwardFlowArgsBuilder,
)
from magma.pipelined.tests.pipelined_test_util import (
    FlowTest,
    FlowVerifier,
    PktsToSend,
    SnapshotVerifier,
    SubTest,
    create_service_manager,
    fake_controller_setup,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)


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
        cls.service_manager = create_service_manager(
            [PipelineD.ENFORCEMENT], ['arpd'])
        cls._tbl_num = cls.service_manager.get_table_num(
            GYController.APP_NAME)

        gy_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.GY,
                  PipelinedController.Arp,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.GY:
                    gy_controller_reference,
                PipelinedController.Arp:
                    Future(),
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'CWF',
                'allow_unknown_arps': False,
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'internal_ip_subnet': '192.168.0.0/16',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'enable_queue_pgm': False,
                'local_ue_eth_addr': False,
                'qos': {'enable': False},
                'dpi': {'enable': False},
                'clean_restart': True,
            },
            mconfig=PipelineD(
                ue_ip_block='192.168.128.0/24'
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)

        cls.gy_controller = gy_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

        cls.gy_controller._redirect_manager._save_redirect_entry = MagicMock()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_subscriber_redirect_policy(self):
        """
        Add redirect policy to subscriber, send 4096 packets

        Assert:
            Packets are properly matched with the 'simple_match' policy
            Send /20 (4096) packets, match /16 (256) packets
        """
        fake_controller_setup(self.gy_controller)
        imsi = 'IMSI010000000088888'
        sub_ip = '192.168.128.74'
        uplink_tunnel = 0x1234
        redirect_ips = ["185.128.101.5", "185.128.121.4"]
        self.gy_controller._redirect_manager._dns_cache.get(
            "about.sha.ddih.org", lambda: redirect_ips, max_age=42
        )
        flow_list = [FlowDescription(match=FlowMatch())]
        policy = VersionedPolicy(
            rule= PolicyRule(
                id='redir_test', priority=3, flow_list=flow_list,
                redirect=RedirectInformation(
                    support=1,
                    address_type=2,
                    server_address="http://about.sha.ddih.org/"
                )
            ),
            version=1,
        )

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel, self.gy_controller, self._tbl_num
        ).add_policy(policy)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = TCPPacketBuilder()\
            .set_tcp_layer(42132, 80, 2)\
            .set_tcp_flags("S")\
            .set_ip_layer('151.42.41.122', sub_ip)\
            .set_ether_layer(self.MAC_DEST, "01:20:10:20:aa:bb")\
            .build()

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             include_stats=False)

        with isolator, sub_context, snapshot_verifier:
            pkt_sender.send(packet)

    def test_subscriber_restrict_policy(self):
        """
        Add restrict policy to subscriber, send 4096 packets

        Assert:
            Packets are properly matched with the 'restrict_match' policy
            Send /20 (4096) packets, match /16 (256) packets
        """
        fake_controller_setup(self.gy_controller)
        imsi = 'IMSI010000000088888'
        uplink_tunnel = 0x1234
        sub_ip = '192.168.128.74'
        flow_list1 = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('8.8.8.0/24'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policies = [
            VersionedPolicy(
                rule=PolicyRule(id='restrict_match', priority=2, flow_list=flow_list1),
                version=1,
            )
        ]
        pkts_matched = 256
        pkts_sent = 4096

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel, self.gy_controller, self._tbl_num
        ).add_policy(policies[0])
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder()\
            .set_ip_layer('8.8.8.8', sub_ip)\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .build()
        flow_query = FlowQuery(
            self._tbl_num, self.testing_controller,
            match=flow_match_to_magma_match(flow_list1[0].match)
        )

        # =========================== Verification ===========================
        # Verify aggregate table stats, subscriber 1 'simple_match' pkt count
        flow_verifier = FlowVerifier([
            FlowTest(FlowQuery(self._tbl_num, self.testing_controller),
                     pkts_sent),
            FlowTest(flow_query, pkts_matched)
        ], lambda: wait_after_send(self.testing_controller))
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             include_stats=False)

        with isolator, sub_context, flow_verifier, snapshot_verifier:
            pkt_sender.send(packet)

if __name__ == "__main__":
    unittest.main()
