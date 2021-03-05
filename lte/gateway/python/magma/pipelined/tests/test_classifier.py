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
import os
import warnings
from concurrent.futures import Future
from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    assert_bridge_snapshot_match,
)
from magma.pipelined.tests.pipelined_test_util import start_ryu_app_thread, \
     stop_ryu_app_thread, create_service_manager, wait_after_send, \
     SnapshotVerifier
from magma.pipelined.app.classifier import Classifier

class ClassifierTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    EnodeB_IP = "192.168.60.141"
    EnodeB2_IP = "192.168.60.140"
    MTR_IP = "10.0.2.10"
    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(ClassifierTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['classifier'])
        classifier_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.Classifier,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.Classifier:
                    classifier_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'ovs_mtr_port_number': 15577,
                'mtr_ip': cls.MTR_IP,
                'ovs_internal_sampling_port_number': 15578,
                'ovs_internal_sampling_fwd_tbl_number': 201,
                'clean_restart': True,
                'ovs_multi_tunnel': True,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )
        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        cls.thread = start_ryu_app_thread(test_setup)
        cls.classifier_controller = classifier_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_install_internal_pkt_fwd_flow(self):
        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.test_detach_default_tunnel_flows()
        self.classifier_controller._install_internal_pkt_fwd_flow()
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass
    def test_detach_default_tunnel_flows(self):
        self.classifier_controller._delete_all_flows()

    def test_attach_tunnel_flows(self):

        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.test_detach_default_tunnel_flows()

        seid1 = 5000
        self.classifier_controller.add_tunnel_flows(65525, 1, 100000,
                                                     "192.168.128.30",
                                                     self.EnodeB_IP, seid1)

        seid2 = 5001
        self.classifier_controller.add_tunnel_flows(65525, 2,100001,
                                                     "192.168.128.31",
                                                     self.EnodeB_IP, seid2)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_detach_tunnel_flows(self):

        self.classifier_controller.delete_tunnel_flows(1, "192.168.128.30")

        self.classifier_controller.delete_tunnel_flows(2, "192.168.128.31")

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_attach_multi_tunnel_flows(self):

        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.test_detach_default_tunnel_flows()

        seid1 = 5000
        self.classifier_controller.add_tunnel_flows(65525, 1, 100000,
                                                     "192.168.128.30",
                                                     self.EnodeB_IP, seid1)

        seid2 = 5001
        self.classifier_controller.add_tunnel_flows(65525, 2,100001,
                                                     "192.168.128.31",
                                                     self.EnodeB2_IP, seid2)

        self.classifier_controller.add_tunnel_flows(65525, 5,1001,
                                                     "192.168.128.51",
                                                     self.EnodeB2_IP, seid2)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_detach_multi_tunnel_flows(self):

        self.classifier_controller.delete_tunnel_flows(1, "192.168.128.30", self.EnodeB_IP)

        self.classifier_controller.delete_tunnel_flows(2, "192.168.128.31", self.EnodeB2_IP)

        self.classifier_controller.delete_tunnel_flows(5, "192.168.128.51", self.EnodeB2_IP)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_discard_tunnel_flows(self):

        self.classifier_controller._delete_all_flows()
        self.classifier_controller._discard_tunnel_flows(65525, 3,
                                                         "192.168.128.80")

        self.classifier_controller._discard_tunnel_flows(65525, 4,
                                                         "192.168.128.82")
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_resume_tunnel_flows(self):

        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.test_detach_default_tunnel_flows()
        self.classifier_controller._resume_tunnel_flows(65525, 3,
                                                        "192.168.128.80")

        self.classifier_controller._resume_tunnel_flows(65525, 4,
                                                        "192.168.128.82")

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()
