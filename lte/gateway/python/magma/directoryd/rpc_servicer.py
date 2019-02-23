"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from orc8r.protos.directoryd_pb2_grpc import DirectoryServiceServicer, \
    DirectoryServiceStub, add_DirectoryServiceServicer_to_server

from magma.common.rpc_utils import return_void
from magma.common.service_registry import ServiceRegistry


class DirectoryServiceRpcServicer(DirectoryServiceServicer):
    """ gRPC based server for the Directoryd. """

    def __init__(self, mconfig, config):
        pass

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_DirectoryServiceServicer_to_server(self, server)

    def GetLocation(self, request_msg, context):
        """ Get the location record of an object

        Args:
            request_msg (GetLocationRequest): get location request

        Returns:
            LocationRecord: location record
        """
        location_record = self._get_grpc_client().GetLocation(request_msg)
        return location_record

    @return_void
    def UpdateLocation(self, request_msg, context):
        """ Update the location record of an object

        Args:
            request_msg (UpdateLocationRequest): update location
            request
        """
        self._get_grpc_client().UpdateLocation(request_msg)

    @return_void
    def DeleteLocation(self, request_msg, context):
        """ Delete the location record of an object

        Args:
            request_msg (DeleteLocationRequest): delete location
            request
        """
        self._get_grpc_client().DeleteLocation(request_msg)

    def _get_grpc_client(self):
        chan = ServiceRegistry.get_rpc_channel(
            'directoryd', ServiceRegistry.CLOUD)
        return DirectoryServiceStub(chan)
