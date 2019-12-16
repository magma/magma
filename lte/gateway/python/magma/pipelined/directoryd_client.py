"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import grpc
import logging

from magma.common.service_registry import ServiceRegistry
from orc8r.protos.directoryd_pb2 import UpdateRecordRequest
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub

DIRECTORYD_SERVICE_NAME = "directoryd"
DEFAULT_GRPC_TIMEOUT = 10
IPV4_ADDR_KEY = "ipv4_addr"


def update_record(imsi: str, ip_addr: str) -> None:
    """
    Make RPC call to 'UpdateRecord' method of local directoryD service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(DIRECTORYD_SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except ValueError:
        logging.error('Cant get RPC channel to %s', DIRECTORYD_SERVICE_NAME)
        return
    client = GatewayDirectoryServiceStub(chan)
    if not imsi.startswith("IMSI"):
        imsi = "IMSI" + imsi
    try:
        # Location will be filled in by directory service
        req = UpdateRecordRequest(id=imsi, location="hwid")
        req.fields[IPV4_ADDR_KEY] = ip_addr
        client.UpdateRecord(req, DEFAULT_GRPC_TIMEOUT)
    except grpc.RpcError as err:
        logging.error(
            "UpdateRecordRequest error for id: %s, ipv4_addr: %s! [%s] %s",
            imsi,
            ip_addr,
            err.code(),
            err.details())
