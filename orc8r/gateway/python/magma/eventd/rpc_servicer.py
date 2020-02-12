"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.

import logging

from orc8r.protos import eventd_pb2_grpc

from magma.common.rpc_utils import return_void


class EventDRpcServicer(eventd_pb2_grpc.EventServiceServicer):
    """
    gRPC based server for EventD.
    """
    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        eventd_pb2_grpc.add_EventServiceServicer_to_server(self, server)

    @return_void
    def LogEvent(self, request, context):
        """
        Logs an event.
        """
        logging.error("Logging event: %s", request)
