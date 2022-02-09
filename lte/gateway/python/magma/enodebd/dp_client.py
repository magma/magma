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
import grpc
from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult
from dp.protos.enodebd_dp_pb2_grpc import DPServiceStub
from magma.common.service_registry import ServiceRegistry
from magma.enodebd.logger import EnodebdLogger

logger = EnodebdLogger

DP_SERVICE_NAME = "dp_service"
DEFAULT_GRPC_TIMEOUT = 20


def get_cbsd_state(request: CBSDRequest) -> CBSDStateResult:
    """
    Make RPC call to 'GetCBSDState' method of dp service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(
            DP_SERVICE_NAME,
            ServiceRegistry.CLOUD,
        )
    except ValueError:
        logger.error('Cant get RPC channel to %s', DP_SERVICE_NAME)
        return CBSDStateResult(radio_enabled=False)
    client = DPServiceStub(chan)
    try:
        res = client.GetCBSDState(request, DEFAULT_GRPC_TIMEOUT)
    except grpc.RpcError as err:
        logger.warning(
            "GetCBSDState error: [%s] %s",
            err.code(),
            err.details(),
        )
        return CBSDStateResult(radio_enabled=False)
    return res
