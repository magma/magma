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
from lte.protos.pipelined_pb2 import SetupUEMacRequest, UEMacFlowRequest
from magma.pipelined.app.base import global_epoch
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    FlowTest,
    FlowVerifier,
    SnapshotVerifier,
    create_service_manager,
    fake_cwf_setup,
    get_enforcement_stats,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
    wait_for_enforcement_stats,
)
from magma.subscriberdb.sid import SIDUtils
from orc8r.protos.directoryd_pb2 import DirectoryRecord


class CWFRestartResilienceTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP_ADDRESS = '192.168.128.1'
    UE_BLOCK = '192.168.128.0/24'
    DPI_PORT = 'mon1'
    DPI_IP = '1.1.1.1'

    @classmethod
    @unittest.mock.patch('netifaces.ifaddresses',
                         return_value=[[{'addr': '00:11:22:33:44:55'}]])
    @unittest.mock.patch('netifaces.AF_LINK', 0)
    def setUpClass(cls, *_):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of ue_mac_controller
        """
        super(CWFRestartResilienceTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['ue_mac', 'arpd'])

        ue_mac_controller_reference = Future()
        arp_controller_reference = Future()
        testing_controller_reference = Future()

        def mock_thread_safe(cmd, body):
            cmd(body)
        loop_mock = MagicMock()
        loop_mock.call_soon_threadsafe = mock_thread_safe

        test_setup = TestSetup(
            apps=[PipelinedController.UEMac,
                  PipelinedController.Arp,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows,
                  ],
            references={
                PipelinedController.UEMac:
                    ue_mac_controller_reference,
                PipelinedController.Arp:
                    arp_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'CWF',
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP_ADDRESS,
                'enforcement': {'poll_interval': 5},
                'internal_ip_subnet': '192.168.0.0/16',
                'nat_iface': 'eth2',
                'local_ue_eth_addr': False,
                'allow_unknown_arps': False,
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': False,
                'quota_check_ip': '1.2.3.4',
                'enable_nat': False,
                'dpi': {
                    'enabled': False,
                    'mon_port': 'mon1',
                    'mon_port_number': 32769,
                    'idle_timeout': 42,
                },
            },
            mconfig=PipelineD(
                ue_ip_block=cls.UE_BLOCK,
            ),
            loop=loop_mock,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        BridgeTools.create_internal_iface(cls.BRIDGE, cls.DPI_PORT,
                                          cls.DPI_IP)

        cls.thread = start_ryu_app_thread(test_setup)

        cls.ue_mac_controller = ue_mac_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()
        cls.arp_controller = arp_controller_reference.result()
        cls.arp_controller.add_arp_response_flow = MagicMock()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    @unittest.mock.patch('magma.pipelined.app.arp.get_all_records')
    def test_ue_mac_restart(self, directoryd_mock):
        """
        Verify that default flows are properly installed with empty setup
        Verify that ue mac flows are properly restored, with arp recovery from
        directoryd
        """

        imsi1 = 'IMSI111111111111111'
        imsi2 = 'IMSI222222222222222'
        ip1 = '152.81.12.41'
        mac1 = '5e:cc:cc:b1:aa:aa'
        mac2 = 'b2:6a:f3:b3:2f:4c'
        ap_mac_addr1 = '11:22:33:44:55:66'
        ap_mac_addr2 = '12:12:13:24:25:26'

        directoryd_mock.return_value = [
            DirectoryRecord(id=imsi1, fields={'ipv4_addr': ip1,
                                              'mac_addr': mac1})
        ]

        fake_cwf_setup(
            ue_mac_controller=self.ue_mac_controller)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             'default_flows',
                                             include_stats=False)
        with snapshot_verifier:
            pass

        setup_ue_mac_request = SetupUEMacRequest(
            requests=[
                UEMacFlowRequest(
                    sid=SIDUtils.to_pb(imsi1),
                    mac_addr=mac1,
                    msisdn='123456',
                    ap_mac_addr=ap_mac_addr1,
                    ap_name='magma',
                    pdp_start_time=1,
                ),
                UEMacFlowRequest(
                    sid=SIDUtils.to_pb(imsi2),
                    mac_addr=mac2,
                    msisdn='654321',
                    ap_mac_addr=ap_mac_addr2,
                    ap_name='amgam',
                    pdp_start_time=9,
                ),
            ],
            epoch=global_epoch,
        )

        fake_cwf_setup(
            ue_mac_controller=self.ue_mac_controller,
            setup_ue_mac_request=setup_ue_mac_request)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             'recovery_flows',
                                             include_stats=False)
        with snapshot_verifier:
            pass
