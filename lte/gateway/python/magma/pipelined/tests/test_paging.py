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
from lte.protos.mobilityd_pb2 import IPAddress
from magma.pipelined.app.classifier import Classifier
from magma.pipelined.app.inout import EGRESS, INGRESS
from magma.pipelined.app.ue_mac import UEMacAddressController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.tests.app.flow_query import \
    RyuDirectFlowQuery as FlowQuery
from magma.pipelined.tests.app.packet_builder import (ARPPacketBuilder,
                                                      DHCPPacketBuilder,
                                                      EtherPacketBuilder,
                                                      UDPPacketBuilder)
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (PipelinedController,
                                                       TestSetup)
from magma.pipelined.tests.pipelined_test_util import (FlowTest, FlowVerifier,
                                                       SnapshotVerifier,
                                                       create_service_manager,
                                                       start_ryu_app_thread,
                                                       stop_ryu_app_thread,
                                                       wait_after_send)
from ryu.lib import hub
from scapy.all import *
from scapy.contrib.gtp import GTP_U_Header


class PagingTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_1 = '5e:cc:cc:b1:49:4b'
    MAC_2 = '0a:00:27:00:00:02'
    BRIDGE_IP = '192.168.128.1'
    EnodeB_IP = '192.168.60.141'
    Dst_nat = '192.168.129.42'
    CLASSIFIER_CONTROLLER_ID = 5

    @classmethod
    @unittest.mock.patch('netifaces.ifaddresses',
                return_value=[[{'addr': '00:aa:bb:cc:dd:ee'}]])
    @unittest.mock.patch('netifaces.AF_LINK', 0)
    def setUpClass(cls, *_):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(PagingTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['classifier'])
        cls._tbl_num = cls.service_manager.get_table_num(Classifier.APP_NAME)
        
        testing_controller_reference = Future()
        classifier_reference = Future()
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
                'internal_ip_subnet': '192.168.0.0/16',
                'ovs_gtp_port_number': 32768,
                'ovs_mtr_port_number': 15577,
                'ovs_internal_sampling_port_number': 15578,
                'ovs_internal_sampling_fwd_tbl_number': 201,
                'ovs_internal_conntrack_port_number': 15579,
                'ovs_internal_conntrack_fwd_tbl_number': 202,
                'clean_restart': True,
                'paging_timeout': 30,
                'classifier_controller_id': 5,
            },
            mconfig=PipelineD(),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
            rpc_stubs={'sessiond_setinterface': MagicMock()}
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        cls.thread = start_ryu_app_thread(test_setup)
        cls.classifier_controller = classifier_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_install_paging_flow(self):
        """
           Add paging flow in table 0
        """
        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.classifier_controller._delete_all_flows()

        ue_ip_addr = "192.168.128.30"
        self.classifier_controller.install_paging_flow(IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                                        200, True) 

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_remove_paging_flow(self):
        """
           Delete the paging flow from table 0
        """
        ue_ip_addr = "192.168.128.30"
        self.classifier_controller.remove_paging_flow(IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')))

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_traffic_paging_flow(self):
        """
           Add paging flow in table 0
        """
        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.classifier_controller._delete_all_flows() 

        ue_ip_addr = "192.168.128.30"
        self.classifier_controller.install_paging_flow(IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                                        200, True)
        # Create a set of packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        eth = Ether(dst=self.MAC_1, src=self.MAC_2)
        ip = IP(src=self.Dst_nat, dst='192.168.128.30')
        o_udp = UDP(sport=2152, dport=2152)
        i_udp = UDP(sport=1111, dport=2222)
        i_tcp = TCP(seq=1, sport=1111, dport=2222)
        i_ip = IP(src='192.168.60.142', dst=self.EnodeB_IP)

        gtp_packet_udp = eth / ip / o_udp / GTP_U_Header(teid=0x1, length=28,gtp_type=255) / i_ip / i_udp
        gtp_packet_tcp = eth / ip / o_udp / GTP_U_Header(teid=0x1, length=68, gtp_type=255) / i_ip / i_tcp

        # Check if these flows were added (queries should return flows)
        flow_queries = [
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(tunnel_id=1, in_port=32768)),
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(ipv4_dst='192.168.128.30'))
        ]
        # =========================== Verification ===========================
        # Verify 2 flows installed for classifier table (2 pkts matched)

        flow_verifier = FlowVerifier(
            [
                FlowTest(FlowQuery(self._tbl_num,
                                   self.testing_controller), 2, 2),
            ], lambda: wait_after_send(self.testing_controller))

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with flow_verifier, snapshot_verifier:
            pkt_sender.send(gtp_packet_udp)
            pkt_sender.send(gtp_packet_tcp)

        flow_verifier.verify()


if __name__ == "__main__":
    unittest.main()
