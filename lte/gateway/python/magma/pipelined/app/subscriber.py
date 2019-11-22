"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from collections import namedtuple
import threading

from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib import hub

from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.meter import MeterController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from orc8r.protos.common_pb2 import Void


class SubscriberController(MagmaController):
    """
    This openflow controller manages a cached table of subscriber IDs that
    are currently active, by periodically polling mobilityd. When a subscriber
    leaves the system, the controller deletes the flows from all the
    required tables.
    """

    APP_NAME = 'subscriber'
    APP_TYPE = ControllerType.LOGICAL
    POLL_TIMEOUT = 3

    SubscriberConfig = namedtuple('SubscriberConfig',
                                  ['enabled', 'poll_interval'])

    def __init__(self, *args, **kwargs):
        super(SubscriberController, self).__init__(*args, **kwargs)
        self.dpset = kwargs['dpset']
        self.mobilityd = kwargs['rpc_stubs']['mobilityd']
        self.loop = kwargs['loop']
        self.config = self._get_config(kwargs['config'])
        if not self.config.enabled:
            return

        # List of tables having per subscriber flows, that need to be cleared.
        self.table_nums = [
            self._service_manager.get_table_num(MeterController.APP_NAME)]
        self._subs_list = set()
        self.worker_thread = hub.spawn(self._run)
        # List of subscribers that should have their meter flows deleted
        self._subs_to_delete_for_meter = []
        # Write lock is needed as subscriber list polling and flow deletion
        # happen in a different threads
        self._subs_to_delete_for_meter_lock = threading.Lock()
        self._meter_poll_active = \
            kwargs['config']['meter']['poll_interval'] >= 0

    def _get_config(self, config_dict):
        return self.SubscriberConfig(
            enabled=config_dict['subscriber']['enabled'],
            poll_interval=config_dict['subscriber']['poll_interval'],
        )

    def _run(self):
        while True:
            self._poll_subscriber_list()
            hub.sleep(self.config.poll_interval)

    def _poll_subscriber_list(self):
        """
        Send a local RPC request to mobilityd to get the current subscribers
        """
        future = self.mobilityd.GetSubscriberIPTable.future(
            Void(), self.POLL_TIMEOUT)
        future.add_done_callback(
            lambda future: self.loop.call_soon_threadsafe(
                self._poll_subscriber_list_done, future))

    def _poll_subscriber_list_done(self, future):
        """
        Process response from mobilityd and find deleted subscribers
        """
        err = future.exception()
        if err:
            logging.error('Error polling subscriber list: %s', err)
            return
        new_list = {entry.sid.id for entry in future.result().entries}
        deleted_subs = self._subs_list - new_list
        if len(deleted_subs) > 0:
            self._process_deleted_subscribers(deleted_subs)
        self._subs_list = new_list

    def _process_deleted_subscribers(self, deleted_subs):
        logging.debug('Processing deleted subs: %s', deleted_subs)
        self._process_deleted_subscribers_for_meter(deleted_subs)

    def _process_deleted_subscribers_for_meter(self, deleted_subs):
        if self._meter_poll_active:
            # If polling is active, schedule the deletion for later to ensure
            # the stats are reported before deletion.
            with self._subs_to_delete_for_meter_lock:
                self._subs_to_delete_for_meter.extend(deleted_subs)
        else:
            self._delete_meter_flows(deleted_subs)

    @set_ev_cls(ofp_event.EventOFPFlowStatsReply, MAIN_DISPATCHER)
    def _flow_stats_reply(self, ev):
        """
        Schedule the flow stats handling in the main event loop, so as to
        unblock the ryu event loop
        """
        self.loop.call_soon_threadsafe(self._handle_flow_stats, ev.msg.body)

    def _handle_flow_stats(self, flow_stats):
        """
        If a flow stats reply is received for the meter flow table, then
        meter_stats will have reported stats of all flows, so any scheduled
        flow deletion can be executed.
        """
        # Ignore the stats reply if polling is off, or if it is not for the
        # metering table.
        stats_for_different_table = any(
            flow_stat.table_id != self._service_manager.get_table_num(
                MeterController.APP_NAME) for flow_stat in flow_stats)
        if not self._meter_poll_active or not flow_stats or \
                stats_for_different_table:
            return

        with self._subs_to_delete_for_meter_lock:
            self._delete_meter_flows(self._subs_to_delete_for_meter)
            self._subs_to_delete_for_meter = []

    def _delete_meter_flows(self, deleted_subs):
        for _, datapath in self.dpset.get_all():
            for imsi in deleted_subs:
                match = MagmaMatch(imsi=encode_imsi(imsi))
                for table in self.table_nums:
                    flows.delete_flow(datapath, table, match)
