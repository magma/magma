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

from netaddr import IPNetwork
import unittest
from unittest.mock import MagicMock
import time
import threading

from nose.plugins.skip import SkipTest
from ryu.app.ofctl.exception import InvalidDatapath
from ryu.base.app_manager import AppManager

from magma.pipelined.openflow.exceptions import MagmaOFError
from magma.pkt_tester.tests.test_topology_builder import check_env
from magma.pipelined.app.base import MagmaController, ControllerType



"""
Writing tests for pipelined
=============================

Most tests for pipelined and associated apps are actually integration tests
-- you want to have OVS running as well as an instance of the controller
process, with each of those elements communicating with each other, while the
test runner exercises each. If you make your test a subclass of
BaseMagmaTest.MagmaControllerTest, this will get taken for you automatically.

Your test case needs to define its own setUp method that defines the apps under
test. Because of how ryu is implemented internally, this is a list of strings
that are import paths to the apps you want included in your test case. See
below for examples.

Next, you'll want to import things that rely on OVS inside your test case,
rather than globally. As part of the base test case setUp, we check to make
sure that the environment supports these controller tests, and skips it if not.
Since the ovs library is built as part of OVS, importing it globally will cause
these tests to raise ImportErrors, rather than skipping.
"""

@unittest.skip("temporarily disabled")
class MagmaControllerTest(unittest.TestCase):
    def setUp(self):
        self.mc = MagmaController()
        self.dp = MagicMock()

    def test_register_table(self):
        self.mc.register_table('foo', 1)
        self.assertEqual(self.mc.TABLES['foo'], 1)

        with self.assertRaises(ValueError):
            self.mc.register_table('foo', 2)

        with self.assertRaises(ValueError):
            self.mc.register_table('bar', 1)

        with self.assertRaises(ValueError):
            self.mc.register_table('bad', 'invalid')

        self.assertEqual(len(self.mc.TABLES), 1)

        self.mc.register_table('bar', 2)

        self.assertEqual(len(self.mc.TABLES), 2)

    def test_sendmsg_fail(self):
        self.dp.send_msg = MagicMock(side_effect=InvalidDatapath("test"))
        with self.assertRaises(MagmaOFError):
            self.mc.send_msg(self.dp, "foo")

    def test_sendmsg_success_after_fail(self):
        returns = (InvalidDatapath("test"), None)
        self.dp.send_msg = MagicMock(side_effect=returns)
        self.mc.send_msg(self.dp, "foo", 2)  # test fails if this raises


class BaseMagmaTest:
    class MagmaControllerTest(unittest.TestCase):
        TEST_BRIDGE = "test_br"
        SRC_PORT = "test_left"
        DST_PORT = "test_right"
        TEST_NETMASK = "255.255.255.0"
        TEST_IPS = {'test_left': "192.168.70.1", "test_right": "192.168.80.1"}
        TEST_NETS = {'test_left': IPNetwork("192.168.70.0/24"),
                     'test_right': IPNetwork("192.168.80.0/24")}

        def setUp(self):
            if not check_env():
                raise SkipTest("Environment doesn't support this test.")

            self._topo_builder = None
            self.apps_under_test = []
            self.mgr = AppManager.get_instance()
            self.mc = None
            self.controller_thread = threading.Thread(target=self._start_controller)

        def tearDown(self):
            if self._topo_builder:
                self._topo_builder.destroy()

            for app in list(self.mgr.applications.values()):
                app.stop()
            self.mgr.close()

            # Ryu doesn't clear out AppManager.applications_cls, which means
            # apps hang around between invocations and breaks isolation between
            # test cases. AppManager uses a singleton pattern, so by setting
            # this to None here we'll ensure the next test case creates a new
            # instance.
            AppManager._instance = None

            self.destroy_bridge()
            self.controller_thread.join()

        def destroy_bridge(self, br_name=None):
            # Destroy the br_name (default: TEST_BRIDGE)
            from ovstest import util

            if not br_name:
                br_name = self.TEST_BRIDGE

            util.start_process(["ovs-vsctl", "del-br", br_name])

            if self.mc:
                for k in list(self.mc.TABLES.keys()): # reset all the tables
                    del self.mc.TABLES[k]

        def _start_controller(self):
            # Actually start running a controller thread
            AppManager.run_apps(self.apps_under_test)

        def _setup_ovs(self):
            """
            This actually causes the switch to come up and connect to the
            controller.
            """
            # Make sure OVS is pointing to the controller.
            from ovstest import util

            # set ovs protocol version
            ret, out, err = util.start_process(["ovs-vsctl", "set", "bridge",
                                                self.TEST_BRIDGE,
                                                "protocols=OpenFlow10,OpenFlow14"])
            # connect to a controller
            ret, out, err = util.start_process(["ovs-vsctl", "set-controller",
                                                self.TEST_BRIDGE,
                                                "tcp:127.0.0.1:6633"])

        def _wait_for_controller(self, app_name="MagmaController"):
            # wait for application
            tries = 0

            while not self.mc and tries < 10:
                try:
                    self.mc = self.mgr.applications[app_name]
                except KeyError:
                    tries += 1
                    time.sleep(0.5)

            if tries > 10:
                raise ValueError("Controller app took too long, failing")

        def _wait_for_datapath(self):
            # wait for connection from switch
            tries = 0

            # we also check to make sure the controller app itself is up
            while not self.mc or (len(self.mc.datapaths) == 0 and tries < 10):
                tries += 1
                time.sleep(0.5)
            if tries > 10:
                raise ValueError("Switch didn't connect in time, failing")

@unittest.skip
class MagmaControllerPktTest(BaseMagmaTest.MagmaControllerTest):
    def setUp(self):
        super(MagmaControllerPktTest, self).setUp()
        self.apps_under_test = ['pipelined.app.base']

    def _generate_topology(self):
        # import here, after we've checked the environment
        from ovstest import util
        from magma.pkt_tester.topology_builder import TopologyBuilder

        self._topo_builder = TopologyBuilder()

        # set up a simple topology
        bridge = self._topo_builder.create_bridge(self.TEST_BRIDGE)
        self._port_no = {}
        for iface_name, ip_address in self.TEST_IPS.items():
            port = self._topo_builder.bind(iface_name, bridge)
            self._topo_builder.create_interface(iface_name,
                                                ip_address,
                                                self.TEST_NETMASK)
            self._port_no[iface_name] = port.port_no

        self.assertFalse(self._topo_builder.invalid_devices())

    def test_delete_all_flows(self):
        # import here, after we've checked the environment
        from ovstest import util
        from magma.pkt_tester.topology_builder import OvsException

        # basic setup
        self._generate_topology()
        self.controller_thread.start()
        self._wait_for_controller()
        self._setup_ovs()
        self._wait_for_datapath()

        # add flows to the test bridge
        for iface in self.TEST_IPS:
            port = self._port_no[iface]
            flow = "in_port=%d,actions=output:%d" % (port, port)
            ret, out, err = util.start_process(["ovs-ofctl", "add-flow",
                                               self.TEST_BRIDGE, flow])
            ret, out, err = util.start_process(["ovs-ofctl", "dump-flows",
                                               self.TEST_BRIDGE])

        self.mc.reset_all_flows(list(self.mc.datapaths.values())[0])

        time.sleep(1.5) # we gotta wait a while in practice :-(

        ret, out, err = util.start_process(["ovs-ofctl", "dump-flows",
                                           self.TEST_BRIDGE])

        flows = out.decode("utf-8").strip().split('\n')

        # when no flows, you get just one element containing "NXST_FLOW"
        self.assertEqual(len(flows), 1)

    def test_delete_table_flows(self):
        # import here, after we've checked the environment
        from ovstest import util
        from magma.pkt_tester.topology_builder import OvsException

        # basic setup
        self._generate_topology()
        self.controller_thread.start()
        self._wait_for_controller()
        self._setup_ovs()
        self._wait_for_datapath()

        # add flows to the test bridge
        for iface in self.TEST_IPS:
            port = self._port_no[iface]
            flow = "in_port=%d,table=5,actions=output:%d" % (port, port)
            ret, out, err = util.start_process(["ovs-ofctl", "add-flow",
                                               self.TEST_BRIDGE, flow])
            flow = "in_port=%d,table=6,actions=output:%d" % (port, port)
            ret, out, err = util.start_process(["ovs-ofctl", "add-flow",
                                               self.TEST_BRIDGE, flow])
            ret, out, err = util.start_process(["ovs-ofctl", "dump-flows",
                                               self.TEST_BRIDGE])

        dp = list(self.mc.datapaths.values())[0]
        self.mc.delete_all_table_flows(dp, table=5)
        time.sleep(1.5)
        ret, out, err = util.start_process(["ovs-ofctl", "dump-flows",
                                           self.TEST_BRIDGE])

        self.assertTrue("table=6" in str(out))
        self.assertFalse("table=5" in str(out))
