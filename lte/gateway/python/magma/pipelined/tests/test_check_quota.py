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

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.pipelined_pb2 import SubscriberQuotaUpdate
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    wait_after_send,
    assert_bridge_snapshot_match,
)


class UEMacAddressTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    BRIDGE_IP = '192.168.130.1'

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
        super(UEMacAddressTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([],
            ['arpd', 'check_quota'])
        check_quota_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.CheckQuotaController,
                  PipelinedController.Arp,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.CheckQuotaController:
                    check_quota_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.Arp:
                    Future(),
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'CWF',
                'allow_unknown_arps': False,
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'has_quota_port': 50001,
                'no_quota_port': 50002,
                'quota_check_ip': '1.2.3.4',
                'local_ue_eth_addr': False,
                'clean_restart': True,
            },
            mconfig=PipelineD(
                ue_ip_block='192.168.128.0/24',
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.check_quota_controller = check_quota_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_add_valid_quota_subscriber(self):
        """
        Add flows for two subscribers
        """
        imsi_1 = 'IMSI010000000088888'
        mac_1 = '5e:cc:cc:b1:49:4b'

        # Add subscriber with UE MAC address
        self.check_quota_controller.update_subscriber_quota_state(
            [
                SubscriberQuotaUpdate(
                    sid=SubscriberID(id=imsi_1), mac_addr=mac_1,
                    update_type=SubscriberQuotaUpdate.VALID_QUOTA),
            ]
        )

        wait_after_send(self.testing_controller)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)

    def test_add_three_subscribers(self):
        """
        Add flows for two subscribers
        """
        imsi_1 = 'IMSI010000000088888'
        imsi_2 = 'IMSI010000111111118'
        imsi_3 = 'IMSI010002222222222'
        mac_1 = '5e:cc:cc:b1:49:4b'
        mac_2 = '5e:a:cc:af:aa:fe'
        mac_3 = '5e:bb:cc:aa:aa:fe'

        # Add subscriber with UE MAC address
        self.check_quota_controller.update_subscriber_quota_state(
            [
                SubscriberQuotaUpdate(
                    sid=SubscriberID(id=imsi_1), mac_addr=mac_1,
                    update_type=SubscriberQuotaUpdate.NO_QUOTA),
                SubscriberQuotaUpdate(
                    sid=SubscriberID(id=imsi_2), mac_addr=mac_2,
                    update_type=SubscriberQuotaUpdate.NO_QUOTA),
                SubscriberQuotaUpdate(
                    sid=SubscriberID(id=imsi_3), mac_addr=mac_3,
                    update_type=SubscriberQuotaUpdate.VALID_QUOTA),
            ]
        )
        wait_after_send(self.testing_controller)

        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)
        self.check_quota_controller.update_subscriber_quota_state(
            [
                SubscriberQuotaUpdate(
                    sid=SubscriberID(id=imsi_2), mac_addr=mac_2,
                    update_type=SubscriberQuotaUpdate.TERMINATE),
                SubscriberQuotaUpdate(
                    sid=SubscriberID(id=imsi_3), mac_addr=mac_3,
                    update_type=SubscriberQuotaUpdate.TERMINATE),
            ]
        )


if __name__ == "__main__":
    unittest.main()
