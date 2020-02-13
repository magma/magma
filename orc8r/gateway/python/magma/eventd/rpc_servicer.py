"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import json
import logging
import socket
from contextlib import closing
from magma.common.rpc_utils import return_void
from orc8r.protos import eventd_pb2_grpc, eventd_pb2


class EventDRpcServicer(eventd_pb2_grpc.EventServiceServicer):
    """
    gRPC based server for EventD.
    """
    def __init__(self, fluent_bit_port: int, tcp_timeout: int):
        self.fluent_bit_port = fluent_bit_port
        self.tcp_timeout = tcp_timeout

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        eventd_pb2_grpc.add_EventServiceServicer_to_server(self, server)

    @return_void
    def LogEvent(self, request: eventd_pb2.Event, context):
        """
        Logs an event.
        """
        logging.error("Logging event: %s", request)
        try:
            with closing(socket.create_connection(
                    ('localhost', self.fluent_bit_port),
                    timeout=self.tcp_timeout)) as sock:
                sock.sendall(json.dumps({
                    'stream_name': request.stream_name,
                    'event_type': request.event_type,
                    'tag': request.tag,
                    'value': request.value.hex()
                }).encode('utf-8'))
        except socket.error as e:
            raise e
