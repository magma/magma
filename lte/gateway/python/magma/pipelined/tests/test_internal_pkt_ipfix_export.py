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
from magma.pipelined.app.dpi import DPIController
from magma.pipelined.app.ipfix import IPFIXController
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


class InternalPktIpfixExportTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    DPI_PORT = 'mon1'
    DPI_IP = '1.1.1.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of dpi_controller
        """
        super(InternalPktIpfixExportTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls._static_rule_dict = {}
        cls.service_manager = create_service_manager(
            [PipelineD.DPI], ['ue_mac', 'ipfix'])
        cls._tbl_num = cls.service_manager.get_table_num(
            DPIController.APP_NAME)

        ue_mac_controller_reference = Future()
        dpi_controller_reference = Future()
        ipfix_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.UEMac,
                  PipelinedController.DPI,
                  PipelinedController.IPFIX,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.UEMac:
                    ue_mac_controller_reference,
                PipelinedController.DPI:
                    dpi_controller_reference,
                PipelinedController.Arp:
                    Future(),
                 PipelinedController.IPFIX:
                     ipfix_controller_reference,
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
                'setup_type': 'CWF',
                'dpi': {
                    'enabled': True,
                    'mon_port': 'mon1',
                    'mon_port_number': 32769,
                    'idle_timeout': 42,
                },
                'ipfix': {
                    'enabled': True,
                    'probability': 65,
                    'collector_set_id': 1,
                    'collector_ip': '1.1.1.1',
                    'collector_port': 65010,
                    'cache_timeout': 60,
                    'obs_domain_id': 1,
                    'obs_point_id': 1,
                },
                'conntrackd': {
                    'enabled': True,
                },
                'ovs_gtp_port_number': 32768,
            },
            mconfig=PipelineD(),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        BridgeTools.create_internal_iface(cls.BRIDGE, cls.DPI_PORT,
                                          cls.DPI_IP)

        cls.thread = start_ryu_app_thread(test_setup)

        cls.ue_mac_controller = ue_mac_controller_reference.result()
        cls.dpi_controller = dpi_controller_reference.result()
        cls.ipfix_controller = ipfix_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

        cls.dpi_controller._policy_dict = cls._static_rule_dict

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_subscriber_policy(self):
        """
        Classify DPI flow, verify internal packet is generated

        Assert:
            snapshots math
        """
        imsi = 'IMSI010000000088888'
        ue_mac = '5e:cc:cc:b1:49:4b'

        self.ue_mac_controller.add_ue_mac_flow(imsi, ue_mac)

        flow_match = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP,
            ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.1'),
            ip_src=convert_ipv4_str_to_ip_proto('1.2.3.0'),
            tcp_dst=80, tcp_src=51115, direction=FlowMatch.UPLINK
        )
        self.dpi_controller.add_classify_flow(
            flow_match, FlowRequest.FLOW_FINAL_CLASSIFICATION,
            'base.ip.http.facebook', 'tbd')
        self.ipfix_controller.add_ue_sample_flow(imsi, "magma_is_awesome_msisdn",
            "00:11:22:33:44:55", "apn_name123456789", 145)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             include_stats=False)

        with snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()
