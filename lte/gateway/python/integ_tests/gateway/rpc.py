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
import os
from typing import Any, Dict, List

from magma.common.service_registry import create_grpc_channel
from magma.configuration.mconfigs import unpack_mconfig_any
from orc8r.protos.common_pb2 import Void
from orc8r.protos.magmad_pb2_grpc import MagmadStub
from orc8r.protos.mconfig import mconfigs_pb2


def get_rpc_channel(service):
    """
    Returns a RPC channel to the service in the gateway.
    """
    return create_grpc_channel(
        os.environ.get('GATEWAY_IP', '192.168.60.142'),
        os.environ.get('GATEWAY_PORT', '8443'),
        '%s.local' % service,
    )


def get_gateway_hw_id():
    """
    Get the hardware ID of the gateway. Is blocking.

    Returns:
        hw_id (str): hardware ID of the gateway specified by
            env variable GATEWAY_IP
    """
    magmad_stub = MagmadStub(get_rpc_channel('magmad'))
    stub_response = magmad_stub.GetGatewayId(Void())
    gateway_hw_id = stub_response.gateway_id
    return gateway_hw_id


def get_gateway_service_mconfigs(services: List[str]) -> Dict[str, Any]:
    """
    Get the managed configurations of some gateway services.

    Args:
        services: List of service names to fetch configs for

    Returns:
        service mconfigs keyed by name
    """
    ret = {}
    magmad_stub = MagmadStub(get_rpc_channel('magmad'))
    stub_response = magmad_stub.GetConfigs(Void())
    for srv in services:
        ret[srv] = unpack_mconfig_any(stub_response.configs_by_key[srv])
    return ret


def reset_gateway_mconfigs():
    """
    Delete the stored mconfigs from the gateway.
    """
    magmad_stub = MagmadStub(get_rpc_channel('magmad'))
    magmad_stub.SetConfigs(mconfigs_pb2.GatewayConfigs())
