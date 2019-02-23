"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from typing import Any

from lte.protos.s6a_service_pb2 import DeleteSubscriberRequest
from lte.protos.s6a_service_pb2_grpc import S6aServiceStub
from lte.protos.subscriberdb_pb2 import SubscriberData

from magma.common.service_registry import ServiceRegistry
from magma.common.streamer import StreamerClient


class SubscriberDBStreamerCallback(StreamerClient.Callback):
    """
    Callback implementation for the SubscriberDB StreamerClient instance.
    """

    def __init__(self, store, loop):
        self._store = store
        self._loop = loop

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name, updates, resync):

        logging.info("Processing %d subscriber updates (resync=%s)",
                     len(updates), resync)

        if resync:
            # TODO:
            # - handle database exceptions
            keys = []
            subscribers = []
            for update in updates:
                sub = SubscriberData()
                sub.ParseFromString(update.value)
                subscribers.append(sub)
                keys.append(update.key)
            old_sub_ids = self._store.list_subscribers()
            self.detach_deleted_subscribers(old_sub_ids, keys)
            logging.debug("Resync with subscribers: %s", ','.join(keys))
            self._store.resync(subscribers)
        else:
            # TODO: implement updates
            pass

    def detach_deleted_subscribers(self, old_sub_ids, new_sub_ids):
        """
        Compares current subscriber ids and new subscriber ids list
        just streamed from the cloud to figure out the deleted subscribers.
        Then send grpc DeleteSubscriber request to mme to detach all the
        deleted subscribers.
        :param old_sub_ids: a list of old subscriber ids in the store.
        :param new_sub_ids: a list of new subscriber ids
                just streamed from the cloud
        :return: n/a
        """
        deleted_sub_ids = [sub_id for sub_id in old_sub_ids
                           if sub_id not in set(new_sub_ids)]
        if len(deleted_sub_ids) == 0:
            return
        # send detach request to mme for all deleted subscribers.
        chan = ServiceRegistry.get_rpc_channel('s6a_service',
                                               ServiceRegistry.LOCAL)
        client = S6aServiceStub(chan)
        req = DeleteSubscriberRequest()
        req.imsi_list.extend(deleted_sub_ids)
        future = client.DeleteSubscriber.future(req)
        future.add_done_callback(lambda future:
                                 self._loop.call_soon_threadsafe(
                                     self.detach_deleted_subscribers_done,
                                     future)
                                 )

    def detach_deleted_subscribers_done(self, delete_future):
        """
        Detach deleted subscribers callback to handle exceptions
        """
        err = delete_future.exception()
        if err:
            logging.error("Detach Deleted Subscribers Error! [%s] %s",
                          err.code(), err.details())
