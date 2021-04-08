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
from lte.protos.pipelined_pb2 import FlowRequest
from lte.protos.policydb_pb2 import FlowMatch
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
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


class DPITest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    BRIDGE_IP = '192.168.128.1'
    DPI_PORT = 'mon1'
    DPI_IP = '1.1.1.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(DPITest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([PipelineD.DPI], [])

        dpi_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.DPI,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.DPI:
                    dpi_controller_reference,
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
                'setup_type': 'LTE',
                'dpi': {
                    'enabled': True,
                    'mon_port': cls.DPI_PORT,
                    'mon_port_number': 32769,
                    'idle_timeout': 42,
                },
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        BridgeTools.create_internal_iface(cls.BRIDGE, cls.DPI_PORT,
                                          cls.DPI_IP)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.dpi_controller = dpi_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_add_app_rules(self):
        """
        Test DPI classifier flows are properly added

        Assert:
            1 FLOW_CREATED -> no rule added as its not classified yet
            1 App not tracked -> no rule installed(`notanAPP`)
            3 App types are matched on:
                facebook other
                google_docs other
                viber audio
        """
        flow_match1 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP,
            ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.8'),
            ip_src=convert_ipv4_str_to_ip_proto('1.2.3.4'),
            tcp_dst=80, tcp_src=51115, direction=FlowMatch.UPLINK
        )
        flow_match2 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP,
            ip_dst=convert_ipv4_str_to_ip_proto('1.10.0.1'),
            ip_src=convert_ipv4_str_to_ip_proto('6.2.3.1'),
            tcp_dst=111, tcp_src=222, direction=FlowMatch.UPLINK
        )
        flow_match3 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_UDP,
            ip_dst=convert_ipv4_str_to_ip_proto('22.2.2.24'),
            ip_src=convert_ipv4_str_to_ip_proto('15.22.32.2'),
            udp_src=111, udp_dst=222, direction=FlowMatch.UPLINK
        )
        flow_match_for_no_proto = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_UDP,
            ip_dst=convert_ipv4_str_to_ip_proto('1.1.1.1')
        )
        flow_match_not_added = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_UDP,
            ip_src=convert_ipv4_str_to_ip_proto('22.22.22.22')
        )
        self.dpi_controller.add_classify_flow(
            flow_match_not_added, FlowRequest.FLOW_CREATED,
            'nickproto', 'bestproto')
        self.dpi_controller.add_classify_flow(
            flow_match_for_no_proto, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'notanAPP', 'null')
        self.dpi_controller.add_classify_flow(
            flow_match1, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'base.ip.http.facebook', 'NotReal')
        self.dpi_controller.add_classify_flow(
            flow_match2, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'base.ip.https.google_gen.google_docs', 'MAGMA',)
        self.dpi_controller.add_classify_flow(
            flow_match3, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'base.ip.udp.viber', 'AudioTransfer Receiving',)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_remove_app_rules(self):
        """
        Test DPI classifier flows are properly removed

        Assert:
            Remove the facebook match flow
        """
        flow_match1 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP,
            ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.8'),
            ip_src=convert_ipv4_str_to_ip_proto('1.2.3.4'),
            tcp_dst=80, tcp_src=51115, direction=FlowMatch.UPLINK
        )
        self.dpi_controller.remove_classify_flow(flow_match1)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()
