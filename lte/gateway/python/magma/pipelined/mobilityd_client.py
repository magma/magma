"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

import grpc
import logging

from magma.common.service_registry import ServiceRegistry
from orc8r.protos.common_pb2 import Void
from lte.protos.mobilityd_pb2 import GWInfo
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub

SERVICE_NAME = "mobilityd"
IPV4_ADDR_KEY = "ipv4_addr"


def get_mobilityd_gw_info() -> GWInfo:
    """
    Make RPC call to 'GetGatewayInfo' method of local mobilityD service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except ValueError:
        logging.error('Cant get RPC channel to %s', SERVICE_NAME)
        return GWInfo()

    client = MobilityServiceStub(chan)
    try:
        return client.GetGatewayInfo(Void())
    except grpc.RpcError as err:
        logging.error(
            "GetGatewayInfoRequest error[%s] %s",
            err.code(),
            err.details())


def set_mobilityd_gw_info(gwinfo: GWInfo):
    """
    Make RPC call to 'SetGatewayInfo' method of local mobilityD service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except ValueError:
        logging.error('Cant get RPC channel to %s', SERVICE_NAME)
        return

    client = MobilityServiceStub(chan)
    try:
        client.SetGatewayInfo(gwinfo)
    except grpc.RpcError as err:
        logging.error(
            "GetGatewayInfoRequest error[%s] %s",
            err.code(),
            err.details())
