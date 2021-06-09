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

import ipaddress
import os
import socket
import unittest
import warnings
from concurrent.futures import Future
from unittest.mock import MagicMock
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import IPFlowDL, UESessionSet, UESessionState
from magma.pipelined.app.classifier import Classifier
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (PipelinedController,
                                                       TestSetup)
from magma.pipelined.tests.pipelined_test_util import (SnapshotVerifier,
                                                       assert_bridge_snapshot_match,
                                                       create_service_manager,
                                                       start_ryu_app_thread,
                                                       stop_ryu_app_thread,
                                                       wait_after_send)


class ClassifierMmeTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    EnodeB_IP = "192.168.60.178"
    MTR_IP = "10.0.2.10"
    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(ClassifierMmeTest, cls).setUpClass()
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
                'ovs_internal_conntrack_port_number': 15579,
                'ovs_internal_conntrack_fwd_tbl_number': 202,
                'clean_restart': False,
                'ovs_multi_tunnel': True,
                'paging_timeout': 30,
                'classifier_controller_id': 5,
            },
            mconfig=None,
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


    def test_detach_default_tunnel_flows(self):
        self.classifier_controller._delete_all_flows()

    def test_add_tunnel_ip_flow_dl(self):

        # Need to delete all default flows in table 0 before
        # install the specific flows test case.
        self.test_detach_default_tunnel_flows()

        dest_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.128.12"),
        )
        src_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.129.64"),
        ) 
        ip_flow_dl = IPFlowDL(set_params=71, tcp_dst_port=0,
                              tcp_src_port=5002, udp_dst_port=0, udp_src_port=0, ip_proto=6,
                              dest_ip=dest_ip,src_ip=src_ip, precedence=65530)

        seid1 = 5000
        ue_ip_addr = "192.168.128.30"
        self.classifier_controller.gtp_handler(0, 65525, 1, 100000,
                                               IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                               self.EnodeB_IP, seid1, ip_flow_dl=ip_flow_dl)

        dest_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.128.111"),
        )
        src_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.129.4"),
        )

        ip_flow_dl = IPFlowDL(set_params=70, tcp_dst_port=0,
                              tcp_src_port=0, udp_dst_port=80, udp_src_port=5060, ip_proto=17,
                              dest_ip=dest_ip,src_ip=src_ip, precedence=65525)
 
        seid2 = 5001
        ue_ip_addr = "192.168.128.31"
        self.classifier_controller.gtp_handler(0, 65525, 2,100001,
                                               IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                               self.EnodeB_IP, seid2, ip_flow_dl=ip_flow_dl)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_delete_tunnel_ip_flow_dl(self):

        dest_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.128.12"),
        )
        src_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.129.64"),
        )

        ip_flow_dl = IPFlowDL(set_params=71, tcp_dst_port=0,
                              tcp_src_port=5002, udp_dst_port=0, udp_src_port=0, ip_proto=6,
                              dest_ip=dest_ip,src_ip=src_ip, precedence=65530)

        seid1 = 5000
        ue_ip_addr = "192.168.128.30"
        self.classifier_controller.gtp_handler(1, 65525, 1, 100000, 
                                               IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')), 
                                               self.EnodeB_IP, seid1, ip_flow_dl=ip_flow_dl)

        dest_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.128.111"),
        )
        src_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.129.4"),
        )

        ip_flow_dl = IPFlowDL(set_params=70, tcp_dst_port=0,
                              tcp_src_port=0, udp_dst_port=80, udp_src_port=5060, ip_proto=17,
                              dest_ip=dest_ip,src_ip=src_ip, precedence=65525)

        seid2 = 5001
        ue_ip_addr = "192.168.128.31"
        self.classifier_controller.gtp_handler(1, 65525, 2, 100001,
                                               IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                               self.EnodeB_IP, seid2, ip_flow_dl=ip_flow_dl)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass


    def test_discard_tunnel_ip_flow_dl(self):


        dest_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.128.12"),
        )
        src_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.129.64"),
        )

        ip_flow_dl = IPFlowDL(set_params=71, tcp_dst_port=0,
                              tcp_src_port=5002, udp_dst_port=0, udp_src_port=0, ip_proto=6,
                              dest_ip=dest_ip,src_ip=src_ip, precedence=65530)

        seid1 = 5000
        ue_ip_addr = "192.168.128.30"
        self.classifier_controller.gtp_handler(0, 65525, 1, 100000,
                                               IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                               self.EnodeB_IP, seid1, ip_flow_dl=ip_flow_dl)
        self.classifier_controller.gtp_handler(4, 65525, 1, 100000,
                                               IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                               self.EnodeB_IP, seid1, ip_flow_dl=ip_flow_dl)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_resume_tunnel_ip_flow_dl(self):

        dest_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.128.12"),
        )
        src_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.129.64"),
        )
        ip_flow_dl = IPFlowDL(set_params=71, tcp_dst_port=0,
                              tcp_src_port=5002, udp_dst_port=0, udp_src_port=0, ip_proto=6,
                              dest_ip=dest_ip,src_ip=src_ip, precedence=65530)

        ue_ip_addr = "192.168.128.30"
        seid1 = 5000
        self.classifier_controller.gtp_handler(5, 65525, 1, 100000,
                                               IPAddress(version=IPAddress.IPV4,address=ue_ip_addr.encode('utf-8')),
                                               self.EnodeB_IP, seid1, ip_flow_dl=ip_flow_dl)
        
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

if __name__ == "__main__":
    unittest.main()
