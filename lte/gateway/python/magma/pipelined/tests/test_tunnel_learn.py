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

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from magma.pipelined.app.tunnel_learn import TunnelLearnController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.subscriber import (
    SubContextConfig,
    default_ambr_config,
)
from magma.pipelined.tests.app.table_isolation import (
    RyuDirectTableIsolator,
    RyuForwardFlowArgsBuilder,
)
from magma.pipelined.tests.pipelined_test_util import (
    assert_bridge_snapshot_match,
    create_service_manager,
    start_ryu_app_thread,
    stop_ryu_app_thread,
)


class TunnelLearnTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    INBOUND_TEST_IP = '127.0.0.1'
    OUTBOUND_TEST_IP = '127.1.0.1'
    BOTH_DIR_TEST_IP = '127.2.0.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(TunnelLearnTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([],
            ['ue_mac', 'tunnel_learn'])
        cls._tbl_num = cls.service_manager.get_table_num(
            TunnelLearnController.APP_NAME)

        tunnel_learn_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.TunnelLearnController,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.TunnelLearnController:
                    tunnel_learn_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'CWF',
                'allow_unknown_arps': False,
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'internal_ip_subnet': '192.168.0.0/16',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
                'access_control': {
                    'ip_blocklist': [
                    ]
                }
            },
            mconfig=PipelineD(
                allowed_gre_peers=[],
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.tunnel_learn_controller = \
            tunnel_learn_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_gre_learn(self):
        """
        Test that uplink packet hits the learn rule and a new flow is created
        in the scratch table(snapshot is checked)
        """
        # Set up subscribers
        sub = SubContextConfig('IMSI001010000000013', '192.168.128.74', 0x1234,
                               default_ambr_config, self._tbl_num)

        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub).build_requests(),
            self.testing_controller,
        )

        # Set up packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        pkt = IPPacketBuilder() \
            .set_ip_layer(self.INBOUND_TEST_IP, sub.ip) \
            .set_ether_layer(self.MAC_DEST, "01:02:03:04:05:06") \
            .build()

        with isolator:
            pkt_sender.send(pkt)

        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)


if __name__ == "__main__":
    unittest.main()
