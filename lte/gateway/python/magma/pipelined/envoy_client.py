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
from typing import List

import grpc
import logging

from magma.common.service_registry import ServiceRegistry
from feg.protos.envoy_controller_pb2_grpc import EnvoyControllerStub
from feg.protos.envoy_controller_pb2 import AddUEHeaderEnrichmentRequest

SERVICE_NAME = "envoyd"
IMSI_HDR = 'imsi'
MSISDN_HDR = 'msisdn'


def set_he_urls_for_ue(ip: str, urls: List[str], imsi: str, msisdn: str):
    """
    Make RPC call to 'SetGatewayInfo' method of local mobilityD service
    """

    try:
        chan = ServiceRegistry.get_rpc_channel(SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except grpc.RpcError:
        logging.error('Cant get RPC channel to %s', SERVICE_NAME)
        return

    client = EnvoyControllerStub(chan)
    try:
        h1 = {IMSI_HDR: imsi}
        h2 = {MSISDN_HDR: msisdn}

        he_info = AddUEHeaderEnrichmentRequest(ue_ip=ip,
                                               websites=urls,
                                               headers=[h1, h2])
        client.AddUEHeaderEnrichment(he_info)
    except grpc.RpcError as err:
        logging.error(
            "SetGatewayInfo error[%s] %s",
            err.code(),
            err.details())
