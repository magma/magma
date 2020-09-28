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
import time
import warnings
from concurrent.futures import Future

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from magma.pipelined.openflow.registers import DIRECTION_REG
from magma.pipelined.tests.app.flow_query import RyuRestFlowQuery
from magma.pipelined.app.ipv6_router_solicitation import \
    IPV6RouterSolicitationController
from magma.pipelined.tests.app.table_isolation import RyuRestTableIsolator,\
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder,\
    ARPPacketBuilder
from magma.pipelined.tests.app.start_pipelined import TestSetup, \
    PipelinedController
from magma.pipelined.openflow.registers import DIRECTION_REG, Direction
from magma.pipelined.tests.app.table_isolation import RyuDirectTableIsolator, \
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.pipelined_test_util import start_ryu_app_thread, \
    stop_ryu_app_thread, create_service_manager, wait_after_send, \
    SnapshotVerifier

from scapy.arch import get_if_hwaddr, get_if_addr
from scapy.data import ETHER_BROADCAST, ETH_P_ALL
from scapy.error import Scapy_Exception
from scapy.layers.l2 import ARP, Ether, Dot1Q
from scapy.layers.inet6 import IPv6, ICMPv6ND_RS, ICMPv6NDOptSrcLLAddr, \
    ICMPv6NDOptPrefixInfo, ICMPv6ND_NS
from scapy.sendrecv import srp1, sendp

from ryu.lib import hub

def _pkt_total(stats):
    return sum(n.packets for n in stats)


class IPV6RouterSolicitationTableTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    UE_BLOCK = '192.168.128.0/24'
    UE_MAC = '5e:cc:cc:b1:49:4b'
    UE_IP = '192.168.128.22'
    OTHER_MAC = '0a:00:27:00:00:02'
    OTHER_IP = '1.2.3.4'

    @classmethod
    @unittest.mock.patch('netifaces.ifaddresses',
                return_value=[[{'addr': '00:11:22:33:44:55'}]])
    @unittest.mock.patch('netifaces.AF_LINK', 0)
    def setUpClass(cls, *_):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(IPV6RouterSolicitationTableTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([],
            ['ipv6_router_solicitation'])
        cls._tbl_num = cls.service_manager.get_table_num(IPV6RouterSolicitationController.APP_NAME)

        ipv6_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.IPV6RouterSolicitation,
                PipelinedController.Testing,
                PipelinedController.StartupFlows
            ],
            references={
                PipelinedController.IPV6RouterSolicitation:
                    ipv6_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'LTE',
                'allow_unknown_arps': False,
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'virtual_interface': cls.BRIDGE,
                'local_ue_eth_addr': True,
                'quota_check_ip': '1.2.3.4',
                'clean_restart': True,
                'enable_nat': True,
            },
            mconfig=PipelineD(
                ue_ip_block=cls.UE_BLOCK,
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.solicitation_controller = ipv6_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_ipv6_flows(self):
        """
        Verify that a UPLINK->UE arp request is properly matched
        """
        hub.sleep(5)
        ll_addr = get_if_hwaddr('testing_br')
        print(ll_addr)
        
        pkt_sender = ScapyPacketInjector(self.IFACE)

        pkt_rs = Ether(dst=self.UE_MAC, src=self.OTHER_MAC)
        pkt_rs /= IPv6(src='fe80::24c3:d0ff:fef3:dd82',
                       dst='ff02::2')
        pkt_rs /= ICMPv6ND_RS()
        pkt_rs /= ICMPv6NDOptSrcLLAddr(lladdr=ll_addr)

        pkt_ns = Ether(dst=self.UE_MAC, src=self.OTHER_MAC)
        pkt_ns /= IPv6(src='fe80::24c3:d0ff:fef3:dd82',
                       dst='ff02::2')
        pkt_ns /= ICMPv6ND_NS()
        pkt_ns /= ICMPv6NDOptSrcLLAddr(lladdr=ll_addr)

        print(pkt_rs.show())

        dlink_args = RyuForwardFlowArgsBuilder(self._tbl_num) \
            .set_eth_match(eth_dst=self.UE_MAC, eth_src=self.OTHER_MAC) \
            .set_reg_value(DIRECTION_REG, Direction.IN) \
            .build_requests()
        isolator = RyuDirectTableIsolator(dlink_args, self.testing_controller)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, snapshot_verifier:
            pkt_sender.send(pkt_rs)
            pkt_sender.send(pkt_ns)
            wait_after_send(self.testing_controller)


        hub.sleep(1)



if __name__ == "__main__":
    unittest.main()
