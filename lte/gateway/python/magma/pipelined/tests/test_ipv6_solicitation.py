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
from magma.pipelined.app.ipv6_solicitation import IPV6SolicitationController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.ipv6_prefix_store import (
    get_ipv6_interface_id,
    get_ipv6_prefix,
)
from magma.pipelined.openflow.registers import DIRECTION_REG, Direction
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.table_isolation import (
    RyuDirectTableIsolator,
    RyuForwardFlowArgsBuilder,
)
from magma.pipelined.tests.pipelined_test_util import (
    SnapshotVerifier,
    create_service_manager,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)
from scapy.arch import get_if_hwaddr
from scapy.layers.inet6 import (
    ICMPv6ND_NS,
    ICMPv6ND_RS,
    ICMPv6NDOptSrcLLAddr,
    IPv6,
)
from scapy.layers.l2 import Ether


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
    def setUpClass(cls, *_):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(IPV6RouterSolicitationTableTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager(
            [],
            ['ipv6_solicitation'],
        )
        cls._tbl_num = cls.service_manager.get_table_num(IPV6SolicitationController.APP_NAME)

        ipv6_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.IPV6RouterSolicitation,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
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
                'ipv6_router_addr': 'd88d:aba4:472f:fc95:7e7d:8457:5301:ebce',
                'clean_restart': True,
                'virtual_mac': 'd6:34:bc:81:5d:40',
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

        cls._prefix_dict = {}
        cls.solicitation_controller._prefix_mapper._prefix_by_interface = \
            cls._prefix_dict

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_ipv6_flows(self):
        """
        Verify that a UPLINK->UE arp request is properly matched
        """
        ll_addr = get_if_hwaddr('testing_br')

        pkt_sender = ScapyPacketInjector(self.IFACE)

        pkt_rs = Ether(dst=self.OTHER_MAC, src=self.UE_MAC)
        pkt_rs /= IPv6(
            src='fe80:24c3:d0ff:fef3:9d21:4407:d337:1928',
            dst='ff02::2',
        )
        pkt_rs /= ICMPv6ND_RS()
        pkt_rs /= ICMPv6NDOptSrcLLAddr(lladdr=ll_addr)

        pkt_ns = Ether(dst=self.OTHER_MAC, src=self.UE_MAC)
        pkt_ns /= IPv6(
            src='fe80::9d21:4407:d337:1928',
            dst='ff02::2',
        )
        pkt_ns /= ICMPv6ND_NS(tgt='abcd:87:3::')
        pkt_ns /= ICMPv6NDOptSrcLLAddr(lladdr=ll_addr)

        ipv6_addr = 'ab22:5:6c:9:9d21:4407:d337:1928'
        interface = get_ipv6_interface_id(ipv6_addr)
        prefix = get_ipv6_prefix(ipv6_addr)
        self.service_manager.interface_to_prefix_mapper.save_prefix(
            interface, prefix,
        )

        ulink_args = RyuForwardFlowArgsBuilder(self._tbl_num) \
            .set_eth_match(eth_dst=self.OTHER_MAC, eth_src=self.UE_MAC) \
            .set_reg_value(DIRECTION_REG, Direction.OUT) \
            .build_requests()
        isolator = RyuDirectTableIsolator(ulink_args, self.testing_controller)

        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            include_stats=False,
        )

        with isolator, snapshot_verifier:
            pkt_sender.send(pkt_rs)
            pkt_sender.send(pkt_ns)
            wait_after_send(self.testing_controller, wait_time=5)


if __name__ == "__main__":
    unittest.main()
