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
from dp.protos.cbsd_pb2 import (
    CBSDStateResult,
    EnodebdUpdateCbsdRequest,
    InstallationParam,
)
from dp.protos.cbsd_pb2_grpc import CbsdManagementStub
from google.protobuf.json_format import MessageToJson
from google.protobuf.wrappers_pb2 import (  # pylint: disable=no-name-in-module
    BoolValue,
    DoubleValue,
    StringValue,
)
from magma.common.service_registry import ServiceRegistry
from magma.enodebd.logger import EnodebdLogger

logger = EnodebdLogger

DP_SERVICE_NAME = "dp_service"
DP_ORC8R_SERVICE_NAME = "dp"
DEFAULT_GRPC_TIMEOUT = 20


def _indoortobool(s: str):
    return s.lower() in ['true', 't', '1', 'yes', 'indoor']


def build_enodebd_update_cbsd_request(
    serial_number: str,
    latitude_deg: str,
    longitude_deg: str,
    indoor_deployment: str,
    antenna_height: str,
    antenna_height_type: str,
    cbsd_category: str,
) -> EnodebdUpdateCbsdRequest:
    # cbsd category and antenna height type should be converted to lowercase
    # for the gRPC call
    antenna_height_type = antenna_height_type.lower()
    cbsd_category = cbsd_category.lower()
    # lat and long values are part of tr181 specification, but they are kept in device config
    # transformed and eventually kept within the device config as strings representing degrees
    latitude_deg_float = float(latitude_deg)
    longitude_deg_float = float(longitude_deg)

    indoor_deployment_bool = _indoortobool(indoor_deployment)
    antenna_height_float = float(antenna_height)

    installation_param = InstallationParam(
        latitude_deg=DoubleValue(value=latitude_deg_float),
        longitude_deg=DoubleValue(value=longitude_deg_float),
        indoor_deployment=BoolValue(value=indoor_deployment_bool),
        height_m=DoubleValue(value=antenna_height_float),
        height_type=StringValue(value=antenna_height_type),
    )

    return EnodebdUpdateCbsdRequest(
        serial_number=serial_number,
        installation_param=installation_param,
        cbsd_category=cbsd_category,
    )


def enodebd_update_cbsd(request: EnodebdUpdateCbsdRequest) -> CBSDStateResult:
    """
    Make RPC call to 'EnodebdUpdateCbsd' method of dp orc8r service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(
            DP_ORC8R_SERVICE_NAME,
            ServiceRegistry.CLOUD,
        )
    except ValueError:
        logger.error("Can't get RPC channel to %s", DP_ORC8R_SERVICE_NAME)
        return CBSDStateResult(radio_enabled=False)
    client = CbsdManagementStub(chan)
    try:
        msg_json = MessageToJson(request, including_default_value_fields=True, preserving_proto_field_name=True)
        logger.debug(f"Sending EnodebdUpdateCbsd request: {msg_json}")
        res = client.EnodebdUpdateCbsd(request, DEFAULT_GRPC_TIMEOUT)
        msg_json = MessageToJson(res, including_default_value_fields=True, preserving_proto_field_name=True)
        logger.debug(f"Received EnodebdUpdateCbsd reply: {msg_json}")
    except grpc.RpcError as err:
        logger.warning(
            "EnodebdUpdateCbsd error: [%s] %s",
            err.code(),
            err.details(),
        )
        return CBSDStateResult(radio_enabled=False)
    return res
