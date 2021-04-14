"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import logging
from typing import List

import grpc
from lte.protos.mobilityd_pb2 import GWInfo, IPAddress
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from magma.common.service_registry import ServiceRegistry
from orc8r.protos.common_pb2 import Void

SERVICE_NAME = "mobilityd"
IPV4_ADDR_KEY = "ipv4_addr"


def get_mobilityd_gw_info() -> List[GWInfo]:
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
        return client.ListGatewayInfo(Void()).gw_list
    except grpc.RpcError as err:
        logging.error(
            "ListGatewayInfo error[%s] %s",
            err.code(),
            err.details())
        return []


def set_mobilityd_gw_info(ip: IPAddress, mac: str, vlan: str):
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
        gwinfo = GWInfo(ip=ip, mac=mac, vlan=vlan)
        client.SetGatewayInfo(gwinfo)
    except grpc.RpcError as err:
        logging.error(
            "SetGatewayInfo error[%s] %s",
            err.code(),
            err.details())


def mobilityd_list_ip_blocks():
    """
    Make RPC call to query all ip-blocks.
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except ValueError:
        logging.error('Cant get RPC channel to %s', SERVICE_NAME)
        return

    client = MobilityServiceStub(chan)
    try:
        resp = client.ListAddedIPv4Blocks(Void())
        return resp
    except grpc.RpcError as err:
        logging.error(
            "List IpBlock error[%s] %s",
            err.code(),
            err.details())
