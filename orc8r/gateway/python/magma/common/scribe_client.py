"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import logging

import random
from orc8r.protos.logging_service_pb2 import LogRequest, LoggerDestination
from orc8r.protos.logging_service_pb2_grpc import LoggingServiceStub

from magma.common.service_registry import ServiceRegistry


class ScribeClient(object):
    """
    ScribeClient is the client_api to call Log() into logging_service on cloud.
    User is responsible of formatting a list of LogEntry.
    """
    def __init__(self, loop=None):
        self._loop = loop if loop else asyncio.get_event_loop()

    @staticmethod
    def should_log(sampling_rate):
        return random.random() < sampling_rate

    def log_to_scribe_with_sampling_rate(self, entries, sampling_rate=1):
        """
        Client API to log entries to scribe.
        Args:
            entries: a list of LogEntry, where each contains a
                              category(str),
                              a required timestamp(int),
                              an optional int_map(map<str, int>),
                              an optional normal_map(map<str, str>,
                              an optional tag_set(arr of str),
                              an optional normvector(arr of str),
                              and an optional hw_id(str).
                              to be sent.
            sampling_rate: defaults to 1, and will be logged always if
                              it's 1. Otherwise, all entries of this
                              specific call will be logged or dropped
                              based on a coin flip.

        Returns:
            n/a. If an exception has occurred, the error will be logged.

        """
        self.log_entries_to_dest(entries,
                                 LoggerDestination.Value("SCRIBE"),
                                 sampling_rate)

    def log_entries_to_dest(self, entries, destination, sampling_rate=1):
        """
        Client API to log entries to destination.
        Args:
            entries: a list of LogEntry, where each contains a
                              category(str),
                              a required timestamp(int),
                              an optional int_map(map<str, int>),
                              an optional normal_map(map<str, str>,
                              an optional tag_set(arr of str),
                              an optional normvector(arr of str),
                              and an optional hw_id(str).
                              to be sent.
            destination: the LoggerDestination to log to. Has to be
                            defined as a enum in LoggerDestination
                            in the proto file.
            sampling_rate: defaults to 1, and will be logged always if
                              it's 1. Otherwise, all entries of this
                              specific call will be logged or dropped
                              based on a coin flip.

        Returns:
             n/a. If an exception has occurred, the error will be logged.
        """
        if not self.should_log(sampling_rate):
            return
        chan = ServiceRegistry.get_rpc_channel('logger',
                                               ServiceRegistry.CLOUD)
        client = LoggingServiceStub(chan)
        log_request = LogRequest(
            Entries=entries,
            Destination=destination
        )
        future = client.Log.future(log_request)
        future.add_done_callback(lambda future:
                                 self._loop.call_soon_threadsafe(
                                     self.log_done, future))

    def log_done(self, log_future):
        """
        Log callback to handle exceptions
        """
        err = log_future.exception()
        if err:
            logging.error("Log Error! [%s] %s",
                          err.code(), err.details())
