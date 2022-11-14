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

from magma.pipelined.app import egress, ingress, middle
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    assert_bridge_snapshot_match,
    create_service_manager,
    fake_mandatory_controller_setup,
    start_ryu_app_thread,
    stop_ryu_app_thread,
)
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL


def mocked_get_mac_address_from_iface(interface: str) -> str:
    if interface == 'test_mtr1':
        return 'ae:fa:b2:76:37:5d'
    if interface == 'testing_br':
        return 'bb:fa:b2:76:37:5d'
    raise ValueError(f"No mac address found for interface {interface}")


class InOutTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    MTR_PORT = "test_mtr1"

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(InOutTest, cls).setUpClass()
        middle.get_mac_address_from_iface = mocked_get_mac_address_from_iface
        egress.get_mac_address_from_iface = mocked_get_mac_address_from_iface
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([])

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        BridgeTools.create_internal_iface(
            cls.BRIDGE,
            cls.MTR_PORT, None,
        )
        mtr_port_no = BridgeTools.get_ofport(cls.MTR_PORT)

        ingress_controller_reference = Future()
        middle_controller_reference = Future()
        egress_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.Ingress,
                PipelinedController.Middle,
                PipelinedController.Egress,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.Ingress:
                    ingress_controller_reference,
                PipelinedController.Middle:
                    middle_controller_reference,
                PipelinedController.Egress:
                    egress_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'enable_nat': True,
                'uplink_gw_mac': '11:22:33:44:55:66',
                'uplink_port': OFPP_LOCAL,
                'virtual_interface': cls.BRIDGE,
                'mtr_ip': '5.6.7.8',
                'mtr_interface': cls.MTR_PORT,
                'ovs_mtr_port_number': mtr_port_no,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        cls.thread = start_ryu_app_thread(test_setup)
        cls.ingress_controller = ingress_controller_reference.result()
        cls.middle_controller = middle_controller_reference.result()
        cls.egress_controller = egress_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def testFlowSnapshotMatch(self):
        fake_mandatory_controller_setup(self.ingress_controller)
        fake_mandatory_controller_setup(self.middle_controller)
        fake_mandatory_controller_setup(self.egress_controller)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)


# LTE with incomplete MTR config
class InOutTestLTE(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(InOutTestLTE, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([])

        ingress_controller_reference = Future()
        middle_controller_reference = Future()
        egress_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.Ingress,
                PipelinedController.Middle,
                PipelinedController.Egress,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.Ingress:
                    ingress_controller_reference,
                PipelinedController.Middle:
                    middle_controller_reference,
                PipelinedController.Egress:
                    egress_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'LTE',
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'enable_nat': True,
                'uplink_gw_mac': '11:22:33:44:55:66',
                'mtr_ip': '1.2.3.4',
                'ovs_mtr_port_number': 211,
                'uplink_port': OFPP_LOCAL,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.ingress_controller = ingress_controller_reference.result()
        cls.middle_controller = middle_controller_reference.result()
        cls.egress_controller = egress_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def testFlowSnapshotMatch(self):
        fake_mandatory_controller_setup(self.ingress_controller)
        fake_mandatory_controller_setup(self.middle_controller)
        fake_mandatory_controller_setup(self.egress_controller)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)
