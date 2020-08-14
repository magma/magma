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
from lte.protos.smfupfif_pb2 import (
    SetGroupPDR,
    SetGroupFAR,
    PDI,
    LocalFTeid,
    ApplyAction,
    FwdParam,
    OuterHeaderCreation,
    SessionSet,
)
from magma.pipelined.tests.pipelined_test_util import start_ryu_app_thread, \
     stop_ryu_app_thread, create_service_manager, wait_after_send, \
     SnapshotVerifier
from magma.pipelined.app.gtp_flows import GtpFlows

class GtpFlowsTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    EnodeB_IP = "192.168.60.141"
    MTR_IP = "10.0.2.10"
    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(GtpFlowsTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['gtp_flows'])
        gtp_flows_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.GtpFlows,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.GtpFlows:
                    gtp_flows_reference,
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
                'clean_restart': True,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )
        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        cls.thread = start_ryu_app_thread(test_setup)
        cls.gtp_flows_controller = gtp_flows_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    
    def test_default_rule_gtp(self):
        self.gtp_flows_controller._install_default_tunnel_flows()
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager, include_stats=False)
        with snapshot_verifier:
            pass
    
    def test_install_gtp_tunnel_flows(self):
        seid1 = 5000
        group_pdr1 = SetGroupPDR(precedence=65536, pdi=PDI(local_f_teid=LocalFTeid(teid=1),
                         ue_ip_adr = "192.168.128.30"))
        
        group_pdr2 = SetGroupPDR(precedence =65536, pdi= PDI(local_f_teid=LocalFTeid(teid=2),
                         ue_ip_adr = "192.168.128.31"))
        seid2 = 5001
        group_far1 = SetGroupFAR(fwd_parm=FwdParam(ohdrcr=OuterHeaderCreation(o_teid = 100000, ipv4_adr = self.EnodeB_IP)))
        
        group_far2 = SetGroupFAR(fwd_parm=FwdParam(ohdrcr=OuterHeaderCreation(o_teid = 100001, ipv4_adr = self.EnodeB_IP)))
        
        self.gtp_flows_controller._add_gtp_tunnel_flows(group_pdr1, group_far1, seid1)
        self.gtp_flows_controller._add_gtp_tunnel_flows(group_pdr2, group_far2, seid2)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager, include_stats=False)
        with snapshot_verifier:
            pass

        
    def test_remove_gtp_tunnel_flows(self):
        group_pdr1 = SetGroupPDR(precedence =65536, pdi=PDI(local_f_teid=LocalFTeid(teid = 1),
                         ue_ip_adr = "192.168.128.30"))
        
        group_pdr2 = SetGroupPDR(precedence =65536, pdi=PDI(local_f_teid=LocalFTeid(teid = 2),
                         ue_ip_adr = "192.168.128.31"))
        
        self.gtp_flows_controller._delete_gtp_tunnel_flows(group_pdr1)
        self.gtp_flows_controller._delete_gtp_tunnel_flows(group_pdr2)
        
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager, include_stats=False)
        with snapshot_verifier:
            pass
            
    def test_install_discard_data_gtp_tunnel_flows(self):
        group_pdr1 = SetGroupPDR(precedence =65536, pdi=PDI(local_f_teid=LocalFTeid(teid = 3),
                         ue_ip_adr = "192.168.128.80"))
        
        group_pdr2 = SetGroupPDR(precedence =65536, pdi=PDI(local_f_teid=LocalFTeid(teid = 4),
                         ue_ip_adr = "192.168.128.82"))
                         
        self.gtp_flows_controller._add_discard_data_gtp_tunnel_flows(group_pdr1)
        self.gtp_flows_controller._add_discard_data_gtp_tunnel_flows(group_pdr1)
        
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager, include_stats=False)
        with snapshot_verifier:
            pass
            
    
    def test_install_forward_data_gtp_tunnel_flows(self):
        group_pdr1 = SetGroupPDR(precedence =65536, pdi=PDI(local_f_teid=LocalFTeid(teid = 3),
                         ue_ip_adr = "192.168.128.80"))
        
        group_pdr2 = SetGroupPDR(precedence =65536, pdi=PDI(local_f_teid=LocalFTeid(teid = 4),
                         ue_ip_adr = "192.168.128.82"))
                         
        self.gtp_flows_controller._add_forward_data_gtp_tunnel_flows(group_pdr1)
        self.gtp_flows_controller._add_forward_data_gtp_tunnel_flows(group_pdr1)
        
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager, include_stats=False)
        with snapshot_verifier:
            pass
    
        
if __name__ == "__main__":
    unittest.main()
