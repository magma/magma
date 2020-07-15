"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
import subprocess
import threading
import unittest
import warnings
from concurrent.futures import Future
from unittest import mock
from magma.common.redis.mocks.mock_redis import MockRedis

from magma.mobilityd.uplink_gw import DHCP_Router_key, DHCP_Router_Mac_key
from magma.mobilityd import mobility_store as store

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


class InOutNonNatTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    script_path = "/home/vagrant/magma/lte/gateway/python/magma/mobilityd/"
    uplink_br = "t_up_br0"
    setup_done = False

    @classmethod
    def setup_uplink_br(cls):
        subprocess.check_call(["redis-cli", "flushall"])

        setup_dhcp_server = cls.script_path + "scripts/setup-test-dhcp-srv.sh"
        subprocess.check_call([setup_dhcp_server, "t1"])

        BridgeTools.destroy_bridge(cls.uplink_br)
        setup_uplink_br = [cls.script_path + "scripts/setup-uplink-br.sh",
                           cls.uplink_br,
                           "t1uplink_p0",
                           "8A:00:00:00:00:01"]
        subprocess.check_call(setup_uplink_br)
        cls.setup_done = True

    @mock.patch("redis.Redis", MockRedis)
    def setUp(self):
        self._dhcp_gw_info_mock = store.GatewayInfoMap()

        self._dhcp_gw_info_mock[DHCP_Router_key] = '192.168.128.211'
        logging.info("set router key [{}]".format(self._dhcp_gw_info_mock[DHCP_Router_key]))

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
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'enable_nat': False,
                'non_mat_gw_probe_frequency': .2,
                'non_nat_arp_egress_port': 't_dhcp0',
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

    @mock.patch("redis.Redis", MockRedis)
    def testFlowSnapshotMatch(self):
        # wait for atleast one iteration of the ARP probe.
        while DHCP_Router_Mac_key not in self._dhcp_gw_info_mock:
            threading.Event().wait(0.01)

        threading.Event().wait(0.5)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)
        print("done")


if __name__ == "__main__":
    unittest.main()
