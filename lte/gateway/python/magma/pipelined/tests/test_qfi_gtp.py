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

from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import IPFlowDL
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    SnapshotVerifier,
    create_service_manager,
    start_ryu_app_thread,
    stop_ryu_app_thread,
)


class QfigtpTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_1 = '5e:cc:cc:b1:49:4b'
    MAC_2 = '0a:00:27:00:00:02'
    BRIDGE_IP = '192.168.128.1'
    EnodeB_IP = "192.168.60.178"
    EnodeB2_IP = "192.168.60.190"
    MTR_IP = "10.0.2.10"
    Dst_nat = '192.168.129.42'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(QfigtpTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['classifier'])
        classifier_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.Classifier,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
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
                'ovs_internal_conntrack_port_number': 15579,
                'ovs_internal_conntrack_fwd_tbl_number': 202,
                'clean_restart': True,
                'ovs_multi_tunnel': False,
                'paging_timeout': 30,
                'classifier_controller_id': 5,
                'enable_nat': True,
                'ovs_uplink_port_name': "patch-up",
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
            rpc_stubs={'sessiond_setinterface': MagicMock()},
        )
        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        cls.thread = start_ryu_app_thread(test_setup)
        cls.classifier_controller = classifier_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_qfi_tunnel_flows(self):

        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.classifier_controller._delete_all_flows()

        seid1 = 5000
        ue_ip_addr = "192.168.128.30"
        ip_flow_dl = IPFlowDL(set_params=0)
        self.classifier_controller.add_tunnel_flows(
            65525, 1, 100000,
            self.EnodeB_IP,
            IPAddress(version=IPAddress.IPV4, address=ue_ip_addr.encode('utf-8')),
            seid1, True, ip_flow_dl=ip_flow_dl, session_qfi=9,
        )

        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            snapshot_name='with_qfi_flows',
            include_stats=False,
        )
        with snapshot_verifier:
            pass

        self.classifier_controller.delete_tunnel_flows(
            1, IPAddress(version=IPAddress.IPV4, address=ue_ip_addr.encode('utf-8')),
            ip_flow_dl=ip_flow_dl, session_qfi=9,
        )

        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            include_stats=False,
            snapshot_name='empty',
        )
        with snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()
