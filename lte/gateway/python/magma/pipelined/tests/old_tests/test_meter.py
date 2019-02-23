"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import time

from test_controller import BaseMagmaTest


class MeterTest(BaseMagmaTest.MagmaControllerTest):
    def setUp(self):
        super(MeterTest, self).setUp()
        self.apps_under_test = ['pipelined.app.meter']

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

    def test_add_meter_flows(self):
        from ovstest import util
        from magma.pkt_tester.topology_builder import OvsException

        self._generate_topology()
        self.controller_thread.start()
        self._wait_for_controller("MeterController")

        # clear out any existing in_blocks and set up for the test
        in_net = self.TEST_NETS[self.SRC_PORT]

        # clear out existing net block to of port mappings
        for k in list(self.mc.IPBLOCK_TO_OFPORT.keys()):
            del self.mc.IPBLOCK_TO_OFPORT[k]

        self.mc.IPBLOCK_TO_OFPORT[in_net] = self._port_no[self.SRC_PORT]

        self._setup_ovs()
        self._wait_for_datapath()

        ret, out, err = util.start_process(["ovs-ofctl", "dump-flows",
                                           self.TEST_BRIDGE])

        dpid = list(self.mc.datapaths.keys())[0]
        self.mc._poll_stats(self.mc.datapaths[dpid])

        time.sleep(0.5) # give the vswitch some time to respond

        # check if we're tracking usage for each user
        # it should be zero since there's no traffic
        for sid in self.mc.ip_to_sid.values():
            self.assertTrue(sid in self.mc.usage)
            ur = self.mc.usage[sid]
            self.assertTrue(ur.bytes_tx == 0)
            self.assertTrue(ur.bytes_rx == 0)
            self.assertTrue(ur.pkts_tx == 0)
            self.assertTrue(ur.pkts_rx == 0)
