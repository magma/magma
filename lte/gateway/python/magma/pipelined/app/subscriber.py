"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from collections import namedtuple

from magma.pipelined.app.base import MagmaController
from magma.pipelined.app.meter import MeterController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from orc8r.protos.common_pb2 import Void
from ryu.lib import hub


class SubscriberController(MagmaController):
    """
    This openflow controller manages a cached table of subscriber IDs that
    are currently active, by periodically polling mobilityd. When a subscriber
    leaves the system, the controller deletes the flows from all the
    required tables.
    TODO: T34106498: We use this controller to guarantee the lifetime of
    the flow to be the same as lifetime of the subscriber. Use enforcement
    stats instead for this usecase.
    """

    APP_NAME = 'subscriber'
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
        for _, datapath in self.dpset.get_all():
            for imsi in deleted_subs:
                match = MagmaMatch(imsi=encode_imsi(imsi))
                for table in self.table_nums:
                    flows.delete_flow(datapath, table, match)
