"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging

import grpc
from lte.protos import subscriberdb_pb2, subscriberdb_pb2_grpc

from magma.common.rpc_utils import return_void
from magma.subscriberdb.sid import SIDUtils
from .store.base import (
    DuplicateSubscriberError,
    SubscriberNotFoundError,
)


class SubscriberDBRpcServicer(subscriberdb_pb2_grpc.SubscriberDBServicer):
    """
    gRPC based server for the SubscriberDB.
    """

    def __init__(self, store):
        """
        Store should be thread-safe since we use a thread pool for requests.
        """
        self._store = store

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        subscriberdb_pb2_grpc.add_SubscriberDBServicer_to_server(self, server)

    @return_void
    def AddSubscriber(self, request, context):
        """
        Adds a subscriber to the store
        """
        sid = SIDUtils.to_str(request.sid)
        logging.debug("Add subscriber rpc for sid: %s", sid)
        try:
            self._store.add_subscriber(request)
        except DuplicateSubscriberError:
            context.set_details('Duplicate subscriber: %s' % sid)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)

    @return_void
    def DeleteSubscriber(self, request, context):
        """
        Deletes a subscriber from the store
        """
        sid = SIDUtils.to_str(request)
        logging.debug("Delete subscriber rpc for sid: %s", sid)
        try:
            self._store.delete_subscriber(sid)
        except SubscriberNotFoundError:
            context.set_details("Subscriber not found: %s" % sid)
            context.set_code(grpc.StatusCode.NOT_FOUND)

    @return_void
    def UpdateSubscriber(self, request, context):
        """
        Updates the subscription data
        """
        sid = SIDUtils.to_str(request.data.sid)
        try:
            with self._store.edit_subscriber(sid, request) as subs:
                request.mask.MergeMessage(request.data, subs)
        except SubscriberNotFoundError:
            context.set_details('Subscriber not found: %s' % sid)
            context.set_code(grpc.StatusCode.NOT_FOUND)

    def GetSubscriberData(self, request, context):
        """
        Returns the subscription data for the subscriber
        """
        sid = SIDUtils.to_str(request)
        try:
            return self._store.get_subscriber_data(sid)
        except SubscriberNotFoundError:
            context.set_details('Subscriber not found: %s' % sid)
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return subscriberdb_pb2.SubscriberData()

    def ListSubscribers(self, request, context):  # pylint:disable=unused-argument
        """
        Returns a list of subscribers from the store
        """
        sids = self._store.list_subscribers()
        sid_msgs = [SIDUtils.to_pb(sid) for sid in sids]
        return subscriberdb_pb2.SubscriberIDSet(sids=sid_msgs)
