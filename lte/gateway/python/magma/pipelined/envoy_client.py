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
from feg.protos.envoy_controller_pb2 import (
    AddUEHeaderEnrichmentRequest,
    AddUEHeaderEnrichmentResult,
    DeactivateUEHeaderEnrichmentRequest,
    DeactivateUEHeaderEnrichmentResult,
    Header,
)
from feg.protos.envoy_controller_pb2_grpc import EnvoyControllerStub
from lte.protos.mobilityd_pb2 import IPAddress
from magma.common.service_registry import ServiceRegistry

SERVICE_NAME = "envoy_controller"
IMSI_HDR = 'imsi'
MSISDN_HDR = 'msisdn'
TIMEOUT_SEC = 30


def activate_he_urls_for_ue(ip: IPAddress, rule_id: str, urls: List[str],
                            imsi: str, msisdn: str) -> bool:
    """
    Make RPC call to 'Envoy Controller' to add target URLs to envoy datapath.
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except grpc.RpcError:
        logging.error('Cant get RPC channel to %s', SERVICE_NAME)
        return False

    client = EnvoyControllerStub(chan)
    try:
        headers = [Header(name=IMSI_HDR, value=imsi)]
        if msisdn:
            headers.append(Header(name=MSISDN_HDR, value=msisdn))
        he_info = AddUEHeaderEnrichmentRequest(ue_ip=ip,
                                               rule_id=rule_id,
                                               websites=urls,
                                               headers=headers)
        ret = client.AddUEHeaderEnrichment(he_info, timeout=TIMEOUT_SEC)
        return ret.result == AddUEHeaderEnrichmentResult.SUCCESS
    except grpc.RpcError as err:
        logging.error(
            "Activate HE proxy error[%s] %s",
            err.code(),
            err.details())

    return False


def deactivate_he_urls_for_ue(ip: IPAddress, rule_id: str) -> bool:
    """
    Make RPC call to 'Envoy Controller' to remove the proxy rule for the UE.
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(SERVICE_NAME,
                                               ServiceRegistry.LOCAL)
    except grpc.RpcError:
        logging.error('Cant get RPC channel to %s', SERVICE_NAME)
        return False

    client = EnvoyControllerStub(chan)

    try:
        he_info = DeactivateUEHeaderEnrichmentRequest(ue_ip=ip, rule_id=rule_id)
        ret = client.DeactivateUEHeaderEnrichment(he_info, timeout=TIMEOUT_SEC)
        return ret.result == DeactivateUEHeaderEnrichmentResult.SUCCESS
    except grpc.RpcError as err:
        logging.error(
            "Deactivate HE proxy error[%s] %s",
            err.code(),
            err.details())

    return False

