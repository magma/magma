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
from concurrent.futures import Future
import asyncio
import warnings
from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import flow_match_to_magma_match
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery \
    as FlowQuery
from magma.pipelined.tests.app.start_pipelined import PipelinedController, \
    TestSetup
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.app.table_isolation import RyuDirectTableIsolator, \
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.pipelined_test_util import FlowTest, FlowVerifier, \
    PktsToSend, SubTest, create_service_manager, start_ryu_app_thread, \
    stop_ryu_app_thread, wait_after_send, SnapshotVerifier, \
    fake_controller_setup
from unittest.mock import MagicMock
from magma.pipelined.qos.common import QosImplType, QosManager
from magma.pipelined.qos.types import QosInfo
from collections import defaultdict
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from lte.protos.pipelined_pb2 import VersionedPolicy

class EnforcementTableTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    MAC_2 = "0a:00:27:00:00:02"
    
    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps
        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of enforcement_controller
        """
        super(EnforcementTableTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls._static_rule_dict = {}
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
                'bridge_ip_address': '192.168.128.1',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
                'uplink_port': 20,
                'enable_nat': True,
                'ovs_gtp_port_number': 10,
                'setup_type': 'LTE',
            },
            mconfig=PipelineD(),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)

        cls.enforcement_controller = enforcement_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

        cls.enforcement_controller._policy_dict = cls._static_rule_dict

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_enforcemnet_rules(self):
        """
        Add QOS policy to enforcement table into OVS. 
        """
        fake_controller_setup(self.enforcement_controller)
        imsi = 'IMSI001010000000013'
        sub_ip = '192.168.128.30'
        flow_list1 = [FlowDescription(match=FlowMatch(
                            direction=FlowMatch.UPLINK),
                            action=FlowDescription.PERMIT),
                            FlowDescription(match=FlowMatch(
                            ip_dst=convert_ipv4_str_to_ip_proto("192.168.0.0/24"),
                            direction=FlowMatch.DOWNLINK),
                            action=FlowDescription.PERMIT)]
 
        self.enforcement_controller.activate_rules(imsi, None, 0, convert_ipv4_str_to_ip_proto(sub_ip),
                                       		   None,  policies = [VersionedPolicy(
                                                   rule=PolicyRule(id='rule1', priority=65530,flow_list=flow_list1),
                                                   version=1,),],local_f_teid_ng=100)
 
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with snapshot_verifier:
            pass

    def test_enforcemnet_rule_multiple_sessions(self):
        """
        Add QOS policy to enforcement table into OVS.
        """
        fake_controller_setup(self.enforcement_controller)
        imsi = 'IMSI001000000000088'
        sub_ip = '192.168.128.40'
        flow_list = [FlowDescription(match=FlowMatch(
                            direction=FlowMatch.UPLINK),
                            action=FlowDescription.PERMIT),
                            FlowDescription(match=FlowMatch(
                            ip_dst=convert_ipv4_str_to_ip_proto("192.168.0.0/24"),
                            direction=FlowMatch.DOWNLINK),
                            action=FlowDescription.PERMIT)]

        self.enforcement_controller.activate_rules(imsi, None, 0, convert_ipv4_str_to_ip_proto(sub_ip),
                                       		   None, policies = [VersionedPolicy(
                                                   rule=PolicyRule(id='rule1', priority=65530,flow_list=flow_list),
                                                   version=2,),],local_f_teid_ng=555)
        imsi = 'IMSI001000000100088'
        sub_ip = '192.168.128.150'
        flow_list = [FlowDescription(match=FlowMatch(
                            direction=FlowMatch.UPLINK),
                            action=FlowDescription.PERMIT),
                            FlowDescription(match=FlowMatch(
                            ip_dst=convert_ipv4_str_to_ip_proto("192.168.0.0/24"),
                            direction=FlowMatch.DOWNLINK),
                            action=FlowDescription.PERMIT)]

        self.enforcement_controller.activate_rules(imsi, None, 0, convert_ipv4_str_to_ip_proto(sub_ip),
                                       		   None, policies = [VersionedPolicy(
                                                   rule=PolicyRule(id='rule2', priority=65536,flow_list=flow_list),
                                                   version=3,),],local_f_teid_ng=5000)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with snapshot_verifier:
            pass
 
if __name__ == "__main__":
    unittest.main()
