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

import os
import unittest

from nose.plugins.skip import SkipTest


def check_env(need_root=True):
    """
    Checks if the tests can run in the current env.
    True if ovstest is installed and the user is root.
    """
    if need_root and os.geteuid():
        return False

    try:
        import ovstest  # pylint: disable=unused-variable, import-error
    except ImportError:
        return False
    return True


class TestTopologyBuilder(unittest.TestCase):
    """
    Test class to test the topology builder utility.
    """
    TEST_BRIDGE_NAME = "test_br"
    TEST_INT_PREFIX = "test_int"
    TEST_NETMASK = "255.255.255.0"
    TEST_IP_PREFIX = "192.168.70."

    def setUp(self):
        # Sanity check the env.
        if not check_env():
            raise SkipTest("Environment does not support this test")

        self._topology_builder = None

    def tearDown(self):
        self._topology_builder.destroy()

    def test_create_ovs_topology(self):
        """
        Creates a test topology with a test bridge and two interfaces with ip
        addresses assigned.
        """
        from magma.pkt_tester.topology_builder import TopologyBuilder

        self._topology_builder = TopologyBuilder()

        bridge = self._topology_builder.create_bridge(self.TEST_BRIDGE_NAME)
        for i in range(0, 2):
            ip_address = self.TEST_IP_PREFIX + str(i + 2)
            iface_name = self.TEST_INT_PREFIX + str(i)
            self._topology_builder.bind(iface_name, bridge)
            self._topology_builder.create_interface(
                iface_name,
                ip_address,
                self.TEST_NETMASK,
            )

        self.assertFalse(self._topology_builder.invalid_devices())

    def test_ports(self):
        """
        Simple validator for port methods.
        """
        from magma.pkt_tester.topology_builder import (
            TopologyBuilder,
            UseAfterFreeException,
        )
        self._topology_builder = TopologyBuilder()
        bridge = self._topology_builder.create_bridge(self.TEST_BRIDGE_NAME)

        ip_address = self.TEST_IP_PREFIX + "2"
        iface_name = self.TEST_INT_PREFIX + "0"

        port = self._topology_builder.bind(iface_name, bridge)
        self._topology_builder.create_interface(
            iface_name, ip_address,
            self.TEST_NETMASK,
        )

        self.assertEqual(port.bridge_name, self.TEST_BRIDGE_NAME)
        self.assertEqual(port.iface_name, iface_name)
        self.assertFalse(self._topology_builder.invalid_devices())
        # Verify that the port number read back is > 0.
        self.assertTrue(port.port_no >= 0)
        port.destroy(free_resource=False)
        self.assertRaises(UseAfterFreeException, port.destroy)
        self.assertRaises(UseAfterFreeException, port.sanity_check)

        # Cleanup
        self._topology_builder.destroy()

    def test_iface(self):
        """
        Simple validator for iface methods
        """
        from magma.pkt_tester.topology_builder import (
            TopologyBuilder,
            UseAfterFreeException,
        )
        self._topology_builder = TopologyBuilder()
        ip_address = self.TEST_IP_PREFIX + "2"
        iface_name = self.TEST_INT_PREFIX + "0"
        bridge = self._topology_builder.create_bridge(self.TEST_BRIDGE_NAME)
        self._topology_builder.bind(iface_name, bridge)
        iface = self._topology_builder.create_interface(
            iface_name,
            ip_address,
            self.TEST_NETMASK,
        )
        self.assertEqual(iface.name, iface_name)
        self.assertEqual(iface.ip_address, ip_address)
        self.assertEqual(iface.netmask, self.TEST_NETMASK)

        iface.destroy()

        self.assertRaises(UseAfterFreeException, iface.up)
        self.assertRaises(UseAfterFreeException, iface.sanity_check)

    def test_bridge(self):
        """
        Simple validator for bridge methods
        """
        from magma.pkt_tester.topology_builder import (
            TopologyBuilder,
            UseAfterFreeException,
        )
        self._topology_builder = TopologyBuilder()
        bridge = self._topology_builder.create_bridge(self.TEST_BRIDGE_NAME)

        self.assertEqual(bridge.name, self.TEST_BRIDGE_NAME)
        bridge.destroy()
        iface_name = self.TEST_INT_PREFIX + "0"
        self.assertRaises(
            UseAfterFreeException, bridge.add_virtual_port,
            iface_name, "internal",
        )
        self.assertRaises(
            UseAfterFreeException, bridge.add_physical_port,
            "eth0",
        )
        self.assertRaises(UseAfterFreeException, bridge.sanity_check)
