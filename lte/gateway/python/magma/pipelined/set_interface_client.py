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
import logging

from magma.common.service_registry import ServiceRegistry
from lte.protos.session_manager_pb2 import UPFNodeState
from lte.protos.session_manager_pb2_grpc import SetInterfaceForUserPlaneStub

SERVICE_NAME = "sessiond"


def set_interface_node_state_association (node_state_info: UPFNodeState):
    """
    Make RPC call to 'SetGatewayInfo' method of local mobilityD service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except ValueError:
        logging.error('Cant get RPC channel to %s', SERVICE_NAME)
        return False

    client = SetInterfaceForUserPlaneStub(chan)
    try:
        client.SetUPFNodeState(node_state_info)
        return True
    except grpc.RpcError as err:
        logging.error(
            "set_upf_node_state_association error[%s] %s",
            err.code(),
            err.details())
