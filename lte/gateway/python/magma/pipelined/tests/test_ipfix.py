"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import warnings
from concurrent.futures import Future

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
from ryu.lib import hub


@unittest.skip("Disabled")
class IPFIXTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(IPFIXTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], include_ipfix=True)

        ipfix_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.IPFIX,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.IPFIX:
                    ipfix_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'ipfix': {
                    'enabled': 'true',
                    'probability': 30000,
                    'collector_ip': '10.22.20.116',
                    'collector_port': 4740,
                    'collector_set_id': 1,
                    'obs_domain_id': 1,
                    'obs_point_id': 1,
                },
                'clean_restart': True,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.ipfix_controller = ipfix_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def testFlowSnapshotMatch(self):
        """
        Ensure 2 IMSIs sample flows are installed

        Assert:
            Snapshots match
        """
        imsi1 = 'IMSI010000000088888'
        imsi2 = 'IMSI010000000011111'
        msidn = 'BigTower'
        apn_mac = '08-00-27-cd-32-07'
        apn_name = 'MagmaBox'
        self.ipfix_controller.add_ue_sample_flow(imsi1, msidn, apn_mac,
                                                 apn_name)
        self.ipfix_controller.add_ue_sample_flow(imsi2, msidn, apn_mac,
                                                 apn_name)

        # Big rule wait a bit for it to appear
        hub.sleep(2)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)


if __name__ == "__main__":
    unittest.main()
