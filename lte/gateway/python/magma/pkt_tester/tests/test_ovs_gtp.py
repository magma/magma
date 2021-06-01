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

import test_topology_builder
from nose.plugins.skip import SkipTest
from scapy.all import IP, UDP, L2Socket  # pylint: disable=no-name-in-module
from scapy.contrib.gtp import GTPCreatePDPContextRequest, GTPHeader


class TestOvsGtp(unittest.TestCase):
    """
    Simple test class that tests matching on GTP tunnel id.
    """
    TEST_BRIDGE_NAME = "test_br"
    TEST_INT_PREFIX = "test_int"
    TEST_NETMASK = "255.255.255.0"
    SRC_NAME = "test_int0"
    DST_NAME = "test_int1"
    TEST_IPS = {SRC_NAME: "192.168.70.2", DST_NAME: "192.168.70.3"}
    GTP_C_PORT = 2123

    def setUp(self):
        """
        Create a basic ovs topology with two interfaces on a single vswitch.
        """
        # Sanity check the env.

        if not test_topology_builder.check_env():
            raise SkipTest("Environment does not support this test")

        from topology_builder import TopologyBuilder

        self._topology_builder = TopologyBuilder()

        bridge = self._topology_builder.create_bridge(self.TEST_BRIDGE_NAME)
        self._port_no = {}
        for iface_name, ip_address in self.TEST_IPS.items():
            port = self._topology_builder.bind(iface_name, bridge)
            self._topology_builder.create_interface(
                iface_name,
                ip_address,
                self.TEST_NETMASK,
            )
            self._port_no[iface_name] = port.port_no

        self.assertFalse(self._topology_builder.invalid_devices())

    @staticmethod
    def create_match_flow(test_tun_id, src_of_port, dst_of_port):
        """
        Create a match flow that matches on the GTP tunnel id and the in_port
        and outputs to the specified out_port.
        This is better done through an openflow controller.
        Args:
            test_tun_id: The tunnel id to match.
            src_of_port: The openflow port number of the src port
            dst_of_port: The destination port number of the dst port
        """
        from magma.pkt_tester.topology_builder import OvsException
        from ovstest import util  # pylint: disable=import-error

        # Set bridge to secure (to prevent learning)
        ret_val, out, err = util.start_process([
            "ovs-vsctl", "set-fail-mode",
            "test_br", "secure",
        ])
        if ret_val:
            raise OvsException(
                "Failed to set bridge in secure mode %s, %s" % (out, err),
            )

        flow = (
            "in_port=%s,gtp_tun_id=%s,actions=output:%s" % (
                src_of_port,
                test_tun_id,
                dst_of_port,
            )
        )
        ret_val, out, err = util.start_process([
            "ovs-ofctl", "add-flow",
            TestOvsGtp.TEST_BRIDGE_NAME,
            flow,
        ])
        if ret_val:
            raise OvsException(
                "Failed to install gtp match flow %s, "
                "%s" % (out, err),
            )

    def test_gtp_ping(self):
        """
        Send a gtp request packet from a source port to a destination port and
        verify that the pkt recvd on the destination port is the same as the
        source packet.
        """
        test_tun_id = 2087
        self.create_match_flow(
            test_tun_id, self._port_no[self.SRC_NAME],
            self._port_no[self.DST_NAME],
        )
        gtp = IP(
            src=self.TEST_IPS[self.SRC_NAME],
            dst=self.TEST_IPS[self.DST_NAME],
        ) / \
            UDP(dport=self.GTP_C_PORT) / GTPHeader(teid=test_tun_id) / \
            GTPCreatePDPContextRequest()
        src_socket = L2Socket(iface=self.SRC_NAME)
        dst_socket = L2Socket(iface=self.DST_NAME)
        sent_len = src_socket.send(gtp)
        recv = dst_socket.recv()
        self.assertEqual(sent_len, len(recv))
        print(sent_len)
