"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import warnings
from unittest.mock import MagicMock
from concurrent.futures import Future

from magma.pipelined.app.meter_stats import UsageRecord
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery
from magma.pipelined.app.meter import MeterController
from magma.pipelined.tests.app.start_pipelined import TestSetup,\
    PipelinedController
from magma.pipelined.tests.app.subscriber import SubContextConfig
from magma.pipelined.tests.app.table_isolation import RyuDirectTableIsolator,\
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import start_ryu_app_thread, \
    stop_ryu_app_thread, wait_after_send, create_service_manager, FlowTest, \
    FlowVerifier, wait_for_meter_stats
from lte.protos.mobilityd_pb2 import SubscriberIPTable, SubscriberIPTableEntry
from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.subscriberdb_pb2 import SubscriberID


class SubscriberTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"

    def setUp(self):
        super(SubscriberTest, self).setUp()
        warnings.simplefilter('ignore')
        service_manager = create_service_manager([PipelineD.METERING])
        self._tbl_num = service_manager.get_table_num(MeterController.APP_NAME)

        meter_ref = Future()
        testing_controller_reference = Future()
        meter_stat_ref = Future()
        subscriber_ref = Future()

        def mock_thread_safe(cmd, body):
            cmd(body)

        loop_mock = MagicMock()
        loop_mock.call_soon_threadsafe = mock_thread_safe

        test_setup = TestSetup(
            apps=[
                PipelinedController.Meter,
                PipelinedController.Testing,
                PipelinedController.MeterStats,
                PipelinedController.Subscriber,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.Meter: meter_ref,
                PipelinedController.Testing: testing_controller_reference,
                PipelinedController.MeterStats: meter_stat_ref,
                PipelinedController.Subscriber: subscriber_ref,
                PipelinedController.StartupFlows: Future(),
            },
            config={
                'bridge_name': self.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'meter': {'poll_interval': -1,
                          'enabled': True},
                'subscriber': {'enabled': True, 'poll_interval': -1},
                'clean_restart': True,
            },
            mconfig={},
            loop=loop_mock,
            service_manager=service_manager,
            integ_test=False,
            rpc_stubs={
                'metering_cloud': MagicMock(),
                'mobilityd': MagicMock(),
            }
        )
        BridgeTools.create_bridge(self.BRIDGE, self.BRIDGE)

        self.thread = start_ryu_app_thread(test_setup)

        self.meter_controller = meter_ref.result()
        self.stats_controller = meter_stat_ref.result()
        self.testing_controller = testing_controller_reference.result()
        self.subscriber_controller = subscriber_ref.result()

        self.stats_controller._sync_stats = MagicMock()
        self.subscriber_controller._poll_subscriber_list = MagicMock()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

    def test_process_deleted_subscriber(self):
        """
        With usage polling off, send packets to install metering flows and
        delete one of them by removing the subscriber from the subscriber ip
        table.

        Verifies that the metering flows for the subscriber is deleted when the
        subscriber is deleted.
        """
        # Set up subscribers
        sub1 = SubContextConfig('IMSI001010000000013', '192.168.128.74',
                                self._tbl_num)
        sub2 = SubContextConfig('IMSI001010000000014', '192.168.128.75',
                                self._tbl_num)

        isolator1 = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub1).build_requests(),
            self.testing_controller,
        )
        isolator2 = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub2).build_requests(),
            self.testing_controller,
        )

        # Set up packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        packets = [
            _make_default_pkt(self.MAC_DEST, '45.10.0.1', sub1.ip),
            _make_default_pkt(self.MAC_DEST, '45.10.0.3', sub2.ip),
        ]

        # Initialize subscriber list in subscriber controller.
        subscriber_ip_table = SubscriberIPTable()
        subscriber_ip_table.entries.extend([
            SubscriberIPTableEntry(sid=SubscriberID(id='IMSI001010000000013')),
            SubscriberIPTableEntry(sid=SubscriberID(id='IMSI001010000000014')),
        ])
        fut = Future()
        fut.set_result(subscriber_ip_table)
        self.subscriber_controller._poll_subscriber_list_done(fut)

        # Verify that after the poll, subscriber 2 flows are deleted while
        # subscriber 1 flows remain.
        sub1_query = RyuDirectFlowQuery(
            self._tbl_num, self.testing_controller,
            match=MagmaMatch(imsi=encode_imsi('IMSI001010000000013')))
        sub2_query = RyuDirectFlowQuery(
            self._tbl_num, self.testing_controller,
            match=MagmaMatch(imsi=encode_imsi('IMSI001010000000014')))
        flow_verifier = FlowVerifier([
            FlowTest(sub1_query, 0, 2),
            FlowTest(sub2_query, 0, 0),
        ], lambda: None)

        # Send packets through pipeline and wait.
        with isolator1, isolator2, flow_verifier:
            # Send packets to create the metering flows. Note that these
            # packets will not be matched because the test setup does not
            # support outputting to port.
            for pkt in packets:
                pkt_sender.send(pkt)
            wait_after_send(self.testing_controller)

            # Update the subscriber list to delete subscriber 2.
            subscriber_ip_table = SubscriberIPTable()
            subscriber_ip_table.entries.extend([
                SubscriberIPTableEntry(
                    sid=SubscriberID(id='IMSI001010000000013')),
            ])
            fut = Future()
            fut.set_result(subscriber_ip_table)
            self.subscriber_controller._poll_subscriber_list_done(fut)

        flow_verifier.verify()


class SubscriberWithPollingTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"

    def setUp(self):
        super(SubscriberWithPollingTest, self).setUp()
        warnings.simplefilter('ignore')
        service_manager = create_service_manager([PipelineD.METERING])
        self._tbl_num = service_manager.get_table_num(MeterController.APP_NAME)

        meter_ref = Future()
        testing_controller_reference = Future()
        meter_stat_ref = Future()
        subscriber_ref = Future()

        def mock_thread_safe(cmd, body):
            cmd(body)

        loop_mock = MagicMock()
        loop_mock.call_soon_threadsafe = mock_thread_safe

        test_setup = TestSetup(
            apps=[
                PipelinedController.Meter,
                PipelinedController.Testing,
                PipelinedController.MeterStats,
                PipelinedController.Subscriber,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.Meter: meter_ref,
                PipelinedController.Testing: testing_controller_reference,
                PipelinedController.MeterStats: meter_stat_ref,
                PipelinedController.Subscriber: subscriber_ref,
                PipelinedController.StartupFlows: Future()
            },
            config={
                'bridge_name': self.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'meter': {'poll_interval': 5,
                          'enabled': True},
                'subscriber': {'enabled': True, 'poll_interval': -1},
                'clean_restart': True,
            },
            mconfig={},
            loop=loop_mock,
            service_manager=service_manager,
            integ_test=False,
            rpc_stubs={
                'metering_cloud': MagicMock(),
                'mobilityd': MagicMock(),
            }
        )
        BridgeTools.create_bridge(self.BRIDGE, self.BRIDGE)

        self.thread = start_ryu_app_thread(test_setup)

        self.meter_controller = meter_ref.result()
        self.stats_controller = meter_stat_ref.result()
        self.testing_controller = testing_controller_reference.result()
        self.subscriber_controller = subscriber_ref.result()

        # Mock out poll stats so we have more control over when the poll
        # happens.
        self.real_poll_stats = self.stats_controller._poll_stats
        self.stats_controller._poll_stats = MagicMock()
        self.stats_controller._sync_stats = MagicMock()
        self.subscriber_controller._poll_subscriber_list = MagicMock()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

    def _poll_stats(self):
        for _, datapath in self.stats_controller.dpset.get_all():
            self.real_poll_stats(datapath)

    def test_process_deleted_subscriber(self):
        """
        With usage polling on, send packets to install metering flows and
        delete one of them by removing the subscriber from the subscriber ip
        table.

        Verifies that the metering flows for the subscriber is deleted after
        the correct usage is reported.
        """
        # Set up subscribers
        sub1 = SubContextConfig('IMSI001010000000013', '192.168.128.74',
                                self._tbl_num)
        sub2 = SubContextConfig('IMSI001010000000014', '192.168.128.75',
                                self._tbl_num)

        isolator1 = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub1).build_requests(),
            self.testing_controller,
        )
        isolator2 = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub2).build_requests(),
            self.testing_controller,
        )

        # Set up packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        packets = [
            _make_default_pkt(self.MAC_DEST, '45.10.0.1', sub1.ip),
            _make_default_pkt(self.MAC_DEST, '45.10.0.3', sub2.ip),
        ]

        # Initialize subscriber list in subscriber controller.
        subscriber_ip_table = SubscriberIPTable()
        subscriber_ip_table.entries.extend([
            SubscriberIPTableEntry(sid=SubscriberID(id='IMSI001010000000013')),
            SubscriberIPTableEntry(sid=SubscriberID(id='IMSI001010000000014')),
        ])
        fut = Future()
        fut.set_result(subscriber_ip_table)
        self.subscriber_controller._poll_subscriber_list_done(fut)

        # Verify that after the poll, flows for subscriber 1 and 2 are
        # installed and the second pair of packets sent are matched.
        sub1_query = RyuDirectFlowQuery(
            self._tbl_num, self.testing_controller,
            match=MagmaMatch(imsi=encode_imsi('IMSI001010000000013')))
        sub2_query = RyuDirectFlowQuery(
            self._tbl_num, self.testing_controller,
            match=MagmaMatch(imsi=encode_imsi('IMSI001010000000014')))
        flow_verifier = FlowVerifier([
            FlowTest(sub1_query, 1, 2),
            FlowTest(sub2_query, 1, 2),
        ], lambda: None)

        # Send packets through pipeline and wait.
        with isolator1, isolator2, flow_verifier:
            # Send packets to create the metering flows. Note that these
            # packets will not be matched because the test setup does not
            # support outputting to port.
            for pkt in packets:
                pkt_sender.send(pkt)
            wait_after_send(self.testing_controller)

            # Update the subscriber list to delete subscriber 2.
            subscriber_ip_table = SubscriberIPTable()
            subscriber_ip_table.entries.extend([
                SubscriberIPTableEntry(
                    sid=SubscriberID(id='IMSI001010000000013')),
            ])
            fut = Future()
            fut.set_result(subscriber_ip_table)
            self.subscriber_controller._poll_subscriber_list_done(fut)

            # Send another pair of packets which will be matched.
            for pkt in packets:
                pkt_sender.send(pkt)
            wait_after_send(self.testing_controller)

            # Temporarily mock out _handle_flow_stats because flow_verifier
            # sends a stats request to the meter table, which will trigger
            # the deletion prematurely.
            handle_flow_stats = self.subscriber_controller._handle_flow_stats
            self.subscriber_controller._handle_flow_stats = MagicMock()

        flow_verifier.verify()
        self.subscriber_controller._handle_flow_stats = handle_flow_stats

        # Verify that after the usage is reported, the flows for subscriber 2
        # are deleted.
        sub1_record = UsageRecord()
        sub1_record.bytes_tx = len(packets[0])
        sub2_record = UsageRecord()
        sub2_record.bytes_tx = len(packets[1])
        target_usage = {
            'IMSI001010000000013': sub1_record,
            'IMSI001010000000014': sub2_record,
        }

        flow_verifier = FlowVerifier([
            FlowTest(sub1_query, 0, 2),
            FlowTest(sub2_query, 0, 0),
        ], lambda: wait_for_meter_stats(self.stats_controller, target_usage))

        with flow_verifier:
            self._poll_stats()

        flow_verifier.verify()


def _make_default_pkt(mac_dest, dst, src):
    return IPPacketBuilder() \
        .set_ip_layer(dst, src) \
        .set_ether_layer(mac_dest, "00:00:00:00:00:00") \
        .build()


if __name__ == "__main__":
    unittest.main()
