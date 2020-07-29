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
import subprocess
import threading
import unittest
import warnings
from concurrent.futures import Future

from lte.protos.mobilityd_pb2 import IPAddress, GWInfo, IPBlock

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

from magma.pipelined.app import inout

gw_info = GWInfo()


def mocked_get_mobilityd_gw_info() -> GWInfo:
    global gw_info
    return gw_info


def mocked_set_mobilityd_gw_info(updated_gw_info: GWInfo) -> GWInfo:
    global gw_info
    gw_info = updated_gw_info


class InOutNonNatTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    SCRIPT_PATH = "/home/vagrant/magma/lte/gateway/python/magma/mobilityd/"
    UPLINK_BR = "t_up_br0"
    NON_NAT_ARP_EGRESS_PORT = "t1uplink_p0"

    @classmethod
    def setup_uplink_br(cls):
        setup_dhcp_server = cls.SCRIPT_PATH + "scripts/setup-test-dhcp-srv.sh"
        subprocess.check_call([setup_dhcp_server, "t1"])

        BridgeTools.destroy_bridge(cls.UPLINK_BR)
        setup_uplink_br = [cls.SCRIPT_PATH + "scripts/setup-uplink-br.sh",
                           cls.UPLINK_BR,
                           cls.NON_NAT_ARP_EGRESS_PORT]
        subprocess.check_call(setup_uplink_br)
        inout.get_mobilityd_gw_info = mocked_get_mobilityd_gw_info
        inout.set_mobilityd_gw_info = mocked_set_mobilityd_gw_info

    def setUp(self):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """

        cls = self.__class__
        super(InOutNonNatTest, cls).setUpClass()
        warnings.simplefilter('ignore')

        cls.setup_uplink_br()
        cls.service_manager = create_service_manager([])

        inout_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.InOut,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.InOut:
                    inout_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'LTE',
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'enable_nat': False,
                'non_mat_gw_probe_frequency': .2,
                'non_nat_arp_egress_port': cls.UPLINK_BR,
                'ovs_uplink_port_name': 'patch-up'
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        cls.thread = start_ryu_app_thread(test_setup)

        cls.inout_controller = inout_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def testFlowSnapshotMatch(self):
        # wait for atleast one iteration of the ARP probe.
        global gw_info
        ip_addr = ipaddress.ip_address("192.168.128.211")
        gw_info.ip.version = IPBlock.IPV4
        gw_info.ip.address = ip_addr.packed

        while gw_info.mac is None or gw_info.mac == '':
            threading.Event().wait(0.1)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)
        assert gw_info.mac == 'b2:a0:cc:85:80:7a'


if __name__ == "__main__":
    unittest.main()
