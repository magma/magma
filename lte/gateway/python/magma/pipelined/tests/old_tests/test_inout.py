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
from test_controller import BaseMagmaTest
import unittest


@unittest.skip("temporarily disabled")
class InoutTest(BaseMagmaTest.MagmaControllerTest):
    def setUp(self):
        super(InoutTest, self).setUp()
        self.apps_under_test = ['pipelined.app.base', 'pipelined.app.inout']

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

    def test_add_inout_flows(self):
        from ovstest import util
        from magma.pkt_tester.topology_builder import OvsException

        self._generate_topology()
        self.controller_thread.start()
        self._wait_for_controller()

        # clear out any existing in_blocks and set up for the test
        in_net = self.TEST_NETS[self.SRC_PORT]
        del self.mc.in_blocks[:]
        self.mc.in_blocks.append(in_net)

        # clear out existing net block to of port mappings
        for k in list(self.mc.IPBLOCK_TO_OFPORT.keys()):
            del self.mc.IPBLOCK_TO_OFPORT[k]

        self.mc.IPBLOCK_TO_OFPORT[in_net] = self._port_no[self.SRC_PORT]

        self._setup_ovs()
        self._wait_for_datapath()

        ret, out, err = util.start_process(["ovs-ofctl", "dump-flows",
                                           self.TEST_BRIDGE])
        flow_string = str(out)

        # check if the flows we expect are loaded
        # directions are tagged properly, and resubmit to right table
        expected = "nw_dst=%s actions=set_field:0->metadata,resubmit(,1)" % in_net
        self.assertTrue(expected in flow_string)

        expected = "nw_src=%s actions=set_field:0x10->metadata,resubmit(,1)" % in_net
        self.assertTrue(expected in flow_string)


