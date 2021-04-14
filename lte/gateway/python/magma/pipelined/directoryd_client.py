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

import grpc
from magma.common.service_registry import ServiceRegistry
from orc8r.protos.common_pb2 import Void
from orc8r.protos.directoryd_pb2 import (
    GetDirectoryFieldRequest,
    UpdateRecordRequest,
)
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub
from ryu.lib import hub

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


def get_record(imsi: str, field: str) -> str:
    """
    Make RPC call to 'GetDirectoryField' method of local directoryD service
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
        req = GetDirectoryFieldRequest(id=imsi, field_key=field)
        res = client.GetDirectoryField(req, DEFAULT_GRPC_TIMEOUT)
        if res.value is not None:
            return res.value
    except grpc.RpcError as err:
        logging.error(
            "GetDirectoryFieldRequest error for id: %s! [%s] %s",
            imsi,
            err.code(),
            err.details())
    return None


def get_all_records(retries: int = 3, sleep_time: float = 0.1) -> [dict]:
    """
    Make RPC call to 'GetAllDirectoryRecords' method of local directoryD service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(DIRECTORYD_SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except ValueError:
        logging.error('Cant get RPC channel to %s', DIRECTORYD_SERVICE_NAME)
        return
    client = GatewayDirectoryServiceStub(chan)
    for _ in range(0, retries):
        try:
            res = client.GetAllDirectoryRecords(Void(), DEFAULT_GRPC_TIMEOUT)
            if res.records is not None:
                return res.records
            hub.sleep(sleep_time)
        except grpc.RpcError as err:
            logging.error(
                "GetAllDirectoryRecords error! [%s] %s",
                err.code(),
                err.details())
    return []
