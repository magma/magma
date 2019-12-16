"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


import grpc
from typing import Dict, List

from orc8r.protos.directoryd_pb2 import DirectoryField, AllDirectoryRecords
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceServicer, \
    DirectoryServiceServicer, DirectoryServiceStub, \
    add_DirectoryServiceServicer_to_server, \
    add_GatewayDirectoryServiceServicer_to_server
from magma.common.misc_utils import get_gateway_hwid
from magma.common.rpc_utils import return_void
from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisFlatDict
from magma.common.redis.serializers import RedisSerde, get_json_serializer, \
    get_json_deserializer
from magma.common.service_registry import ServiceRegistry

DIRECTORYD_REDIS_TYPE = "directory_record"
LOCATION_MAX_LEN = 5


class DirectoryRecord:
    """
    DirectoryRecord holds a location history of up to five gateway hwids
    and a dict of identifiers for the record.
    """
    __slots__ = ['location_history', 'identifiers']

    def __init__(self, location_history: List[str], identifiers: Dict[str, str]):
        self.location_history = location_history
        self.identifiers = identifiers

class DirectoryServiceRpcServicer(DirectoryServiceServicer):
    """ gRPC based server for the Directoryd. """

    def __init__(self, mconfig, config):
        pass

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_DirectoryServiceServicer_to_server(self, server)

    def GetLocation(self, request, context):
        """ Get the location record of an object

        Args:
            request (GetLocationRequest): get location request

        Returns:
            LocationRecord: location record
        """
        location_record = self._get_grpc_client().GetLocation(request)
        return location_record

    @return_void
    def UpdateLocation(self, request, context):
        """ Update the location record of an object

        Args:
            request (UpdateLocationRequest): update location
            request
        """
        self._get_grpc_client().UpdateLocation(request)

    @return_void
    def DeleteLocation(self, request, context):
        """ Delete the location record of an object

        Args:
            request (DeleteLocationRequest): delete location
            request
        """
        self._get_grpc_client().DeleteLocation(request)

    def _get_grpc_client(self):
        chan = ServiceRegistry.get_rpc_channel(
            'directoryd', ServiceRegistry.CLOUD)
        return DirectoryServiceStub(chan)


class GatewayDirectoryServiceRpcServicer(GatewayDirectoryServiceServicer):
    """ gRPC based server for the Directoryd Gateway service. """

    def __init__(self):
        serde = RedisSerde(DIRECTORYD_REDIS_TYPE,
                           get_json_serializer(),
                           get_json_deserializer())
        self._redis_dict = RedisFlatDict(get_default_client(), serde)

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_GatewayDirectoryServiceServicer_to_server(self, server)

    @return_void
    def UpdateRecord(self, request, context):
        """ Update the directory record of an object

        Args:
            request (UpdateRecordRequest): update record request
        """
        if len(request.id) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("ID argument cannot be empty in "
                                "UpdateRecordRequest")
            return

        # Lock Redis for requested key until update is complete
        with self._redis_dict.lock(request.id):
            hwid = get_gateway_hwid()
            record = self._redis_dict.get(request.id) or \
                     DirectoryRecord(location_history=[hwid], identifiers={})

            if record.location_history[0] != hwid:
                record.location_history = [hwid] + record.location_history

            for field_key in request.fields:
                record.identifiers[field_key] = request.fields[field_key]

            # Truncate location history to the five most recent hwid's
            record.location_history = \
                record.location_history[:LOCATION_MAX_LEN]
            self._redis_dict[request.id] = record

    @return_void
    def DeleteRecord(self, request, context):
        """ Delete the directory record for an ID

        Args:
             request (DeleteRecordRequest): delete record request
         """
        if len(request.id) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("ID argument cannot be empty in "
                            "DeleteRecordRequest")
            return

        # Lock Redis for requested key until delete is complete
        with self._redis_dict.lock(request.id):
            if request.id not in self._redis_dict:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Record for ID %s was not found." %
                                    request.id)
                return
            self._redis_dict.mark_as_garbage(request.id)

    def GetDirectoryField(self, request, context):
        """ Get the directory record field for an ID and key

        Args:
             request (GetDirectoryFieldRequest): get directory field request
         """
        if len(request.id) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("ID argument cannot be empty in "
                               "GetDirectoryFieldRequest")
            return
        if len(request.field_key) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("Field key argument cannot be empty in "
                                "GetDirectoryFieldRequest")
            return

        # Lock Redis for requested key until get is complete
        with self._redis_dict.lock(request.id):
            if request.id not in self._redis_dict:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Record for ID %s was not found." %
                                    request.id)
                return
            record = self._redis_dict[request.id]
            if request.field_key not in record.identifiers:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Field %s was not found in record for "
                                    "ID %s" % (request.field_key, request.id))
                return
            return DirectoryField(key=request.field_key,
                                  value=record.identifiers[request.field_key])

    def GetAllDirectoryRecords(self, request, context):
        """ Get all directory records

        Args:
             request (Void): void
        """
        response = AllDirectoryRecords()
        for key in self._redis_dict.keys():
            with self._redis_dict.lock(key):
                # Lookup may produce an exception if the key has been deleted
                # between the call to __iter__ and lock
                try:
                    stored_record = self._redis_dict[key]
                except KeyError:
                    continue
                directory_record = response.records.add()
                directory_record.id = key
                directory_record.location_history[:] = \
                    stored_record.location_history
                for identifier_key in stored_record.identifiers:
                    directory_record.fields[identifier_key] = \
                        stored_record.identifiers[identifier_key]

        return response
