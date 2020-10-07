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
from typing import List, Optional

import grpc

from magma.common.service_registry import ServiceRegistry
from orc8r.protos.common_pb2 import Void
from lte.protos.s1ap_service_pb2_grpc import S1apServiceStub
from magma.enodebd.logger import EnodebdLogger as logger

S1AP_SERVICE_NAME = "s1ap_service"
DEFAULT_GRPC_TIMEOUT = 20


def get_all_enb_connected() -> Optional[List[int]]:
    """
    Make RPC call to 'GetEnbConnected' method of s1ap service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(S1AP_SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except ValueError:
        logger.error('Cant get RPC channel to %s', S1AP_SERVICE_NAME)
        return
    client = S1apServiceStub(chan)
    try:
        res = client.GetEnbConnected(Void(), DEFAULT_GRPC_TIMEOUT)
        return res.enb_ids
    except grpc.RpcError as err:
        logger.warning(
            "GetEnbConnected error: [%s] %s",
            err.code(),
            err.details())
    return []

