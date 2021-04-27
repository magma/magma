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
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import (
    convert_ipv4_str_to_ip_proto,
    flow_match_to_magma_match,
)
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery as FlowQuery
from magma.pipelined.tests.app.packet_builder import TCPPacketBuilder
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
    assert_bridge_snapshot_match,
    create_service_manager,
    fake_controller_setup,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)


class RedirectTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP_ADDRESS = '192.168.128.1'
    # TODO test for multiple incoming requests (why we match on tcp ports)

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of enforcement_controller
        """
        super(RedirectTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([PipelineD.ENFORCEMENT])
        cls._tbl_num = cls.service_manager.get_table_num(
            EnforcementController.APP_NAME)

        enforcement_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.Enforcement,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.Enforcement:
                    enforcement_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP_ADDRESS,
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
                'setup_type': 'LTE',
            },
            mconfig=PipelineD(),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)

        cls.enforcement_controller = enforcement_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

        cls.enforcement_controller._redirect_manager._save_redirect_entry =\
            MagicMock()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_url_redirect(self):
        """
        Partial redirection test, checks if flows were added properly for url
        based redirection.

        Assert:
            1 Packet is matched
            Packet bypass flows are added
            Flow learn action is triggered - another flow is added to the table
        """
        fake_controller_setup(self.enforcement_controller)
        redirect_ips = ["185.128.101.5", "185.128.121.4"]
        self.enforcement_controller._redirect_manager._dns_cache.get(
            "about.sha.ddih.org", lambda: redirect_ips, max_age=42
        )
        imsi = 'IMSI010000000088888'
        uplink_tunnel = 0x1234
        sub_ip = '192.168.128.74'
        flow_list = [FlowDescription(match=FlowMatch())]
        policy = VersionedPolicy(
            rule=PolicyRule(
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
            imsi, sub_ip, uplink_tunnel, self.enforcement_controller, self._tbl_num
        ).add_policy(policy)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = TCPPacketBuilder()\
            .set_tcp_layer(42132, 80, 321)\
            .set_tcp_flags("S")\
            .set_ip_layer('151.42.41.122', sub_ip)\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .build()

        # Check if these flows were added (queries should return flows)
        permit_outbound, permit_inbound = [], []
        for ip in redirect_ips:
            permit_outbound.append(FlowQuery(
                self._tbl_num, self.testing_controller,
                match=flow_match_to_magma_match(
                    FlowMatch(ip_dst=convert_ipv4_str_to_ip_proto(ip),
                              direction=FlowMatch.UPLINK))
            ))
            permit_inbound.append(FlowQuery(
                self._tbl_num, self.testing_controller,
                match=flow_match_to_magma_match(
                    FlowMatch(ip_src=convert_ipv4_str_to_ip_proto(ip),
                              direction=FlowMatch.DOWNLINK))
            ))

        learn_action_flow = flow_match_to_magma_match(
            FlowMatch(
                ip_proto=6, direction=FlowMatch.DOWNLINK,
                ip_src=convert_ipv4_str_to_ip_proto(self.BRIDGE_IP_ADDRESS),
                ip_dst=convert_ipv4_str_to_ip_proto(sub_ip))
        )
        learn_action_query = FlowQuery(self._tbl_num, self.testing_controller,
                                       learn_action_flow)

        # =========================== Verification ===========================
        # 1 packet sent, permit rules installed, learn action installed. Since
        # the enforcement table is entered via the DPI table and the scratch
        # enforcement table, the number of packets handled by the table is 2.
        flow_verifier = FlowVerifier(
            [FlowTest(FlowQuery(self._tbl_num, self.testing_controller), 2),
             FlowTest(learn_action_query, 0, flow_count=1)] +
            [FlowTest(query, 0, flow_count=1) for query in permit_outbound] +
            [FlowTest(query, 0, flow_count=1) for query in permit_inbound],
            lambda: wait_after_send(self.testing_controller))

        with isolator, sub_context, flow_verifier:
            pkt_sender.send(packet)
            assert_bridge_snapshot_match(self, self.BRIDGE,
                                         self.service_manager)

        flow_verifier.verify()

    def test_ipv4_redirect(self):
        """
        Partial redirection test, checks if flows were added properly for ipv4
        based redirection.

        Assert:
            1 Packet is matched
            Packet bypass flows are added
            Flow learn action is triggered - another flow is added to the table
        """
        fake_controller_setup(self.enforcement_controller)
        redirect_ip = "54.12.31.42"
        imsi = 'IMSI012000000088888'
        uplink_tunnel = 0x1234
        sub_ip = '192.168.128.74'
        flow_list = [FlowDescription(match=FlowMatch())]
        policy = VersionedPolicy(
            rule=PolicyRule(
                id='redir_ip_test', priority=3, flow_list=flow_list,
                redirect=RedirectInformation(
                    support=1,
                    address_type=0,
                    server_address=redirect_ip
                )
            ),
            version=1,
        )

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel, self.enforcement_controller, self._tbl_num
        ).add_policy(policy)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = TCPPacketBuilder()\
            .set_tcp_layer(42132, 80, 321)\
            .set_tcp_flags("S")\
            .set_ip_layer('151.42.41.122', sub_ip)\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .build()

        # Check if these flows were added (queries should return flows)
        permit_outbound = FlowQuery(
            self._tbl_num, self.testing_controller,
            match=flow_match_to_magma_match(
                FlowMatch(ip_dst=convert_ipv4_str_to_ip_proto(redirect_ip),
                          direction=FlowMatch.UPLINK))
        )
        permit_inbound = FlowQuery(
            self._tbl_num, self.testing_controller,
            match=flow_match_to_magma_match(
                FlowMatch(ip_src=convert_ipv4_str_to_ip_proto(redirect_ip),
                          direction=FlowMatch.DOWNLINK))
        )
        learn_action_flow = flow_match_to_magma_match(
            FlowMatch(
                ip_proto=6, direction=FlowMatch.DOWNLINK,
                ip_src=convert_ipv4_str_to_ip_proto(self.BRIDGE_IP_ADDRESS),
                ip_dst=convert_ipv4_str_to_ip_proto(sub_ip))
        )
        learn_action_query = FlowQuery(self._tbl_num, self.testing_controller,
                                       learn_action_flow)

        # =========================== Verification ===========================
        # 1 packet sent, permit rules installed, learn action installed. Since
        # the enforcement table is entered via the DPI table and the scratch
        # enforcement table, the number of packets handled by the table is 2.
        flow_verifier = FlowVerifier([
            FlowTest(FlowQuery(self._tbl_num, self.testing_controller), 2),
            FlowTest(permit_outbound, 0, flow_count=1),
            FlowTest(permit_inbound, 0, flow_count=1),
            FlowTest(learn_action_query, 0, flow_count=1)
        ], lambda: wait_after_send(self.testing_controller))

        with isolator, sub_context, flow_verifier:
            pkt_sender.send(packet)
            assert_bridge_snapshot_match(self, self.BRIDGE,
                                         self.service_manager)

        flow_verifier.verify()


if __name__ == "__main__":
    unittest.main()
