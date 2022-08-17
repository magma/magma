"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from datetime import datetime, timezone
from typing import Optional, Union

from magma.db_service.utils import get_cbsd_basic_params
from magma.fluentd_client.client import DPLog
from magma.mappings.request_response_mapping import request_response

Message = Union['DBRequest', 'DBResponse', 'CBSDRequest', 'CBSDStateResult']


def make_dp_log(
        message: Message,
        method_name: Optional[str] = None,
        cbsd: Optional['DBCbsd'] = None,  # noqa: F821
        serial_number: Optional[str] = None,
):
    """
    Create DPLog from a Message

    Args:
        message (Message): Message to be converted to a DPLog
        method_name (Optional[str]): Message to be converted to a DPLog
        cbsd (Optional['DBCbsd']): Message to be converted to a DPLog
        serial_number (Optional[str]): Message to be converted to a DPLog

    Returns:
        DPLog
    """

    mapping = {
        'DBRequest': (_make_log_from_db_request, [message]),
        'DBResponse': (_make_log_from_db_response, [message]),
        'CBSDRequest': (_make_dp_log_from_grpc_request, [method_name, message, cbsd]),
        'CBSDStateResult': (_make_dp_log_from_grpc_response, [method_name, message, cbsd, serial_number]),
    }

    msg_type_name = message.__class__.__name__
    func, args = mapping.get(msg_type_name)
    if not func:
        raise TypeError(
            f"{msg_type_name} is not a valid message type. "
            f"Choose one of: DBRequest, DBResponse, CBSDRequest, CBSDStateResult",
        )
    return func(*args)


def now() -> int:
    return int(datetime.now(timezone.utc).timestamp())


def _make_log_from_db_request(request: 'DBRequest') -> DPLog:  # noqa: F821
    fcc_id, network_id, serial_number = get_cbsd_basic_params(request.cbsd)
    return DPLog(
        event_timestamp=now(),
        log_from='DP',
        log_to='SAS',
        log_name=str(request.type.name),
        log_message=str(request.payload),
        cbsd_serial_number=str(serial_number),
        network_id=str(network_id),
        fcc_id=str(fcc_id),
    )


def _make_log_from_db_response(response: 'DBResponse') -> DPLog:  # noqa: F821
    cbsd = response.request.cbsd
    fcc_id, network_id, serial_number = get_cbsd_basic_params(cbsd)
    log_name = request_response[response.request.type.name]
    response_code = response.payload.get(
        'response', {},
    ).get('responseCode', None)
    return DPLog(
        event_timestamp=now(),
        log_from='SAS',
        log_to='DP',
        log_name=str(log_name),
        log_message=str(response.payload),
        cbsd_serial_number=str(serial_number),
        network_id=str(network_id),
        fcc_id=str(fcc_id),
        response_code=response_code,
    )


def _make_dp_log_from_grpc_request(
        method_name: str,
        request: 'CBSDRequest',  # noqa: F821
        cbsd: 'DBCbsd',  # noqa: F821
):
    fcc_id, network_id, _ = get_cbsd_basic_params(cbsd)

    return DPLog(
        event_timestamp=now(),
        log_from='CBSD',
        log_to='DP',
        log_name=method_name + 'Request',
        log_message=str(request),
        cbsd_serial_number=str(request.serial_number),
        network_id=str(network_id),
        fcc_id=str(fcc_id),
    )


def _make_dp_log_from_grpc_response(
        method_name: str,
        result: 'CBSDStateResult',  # noqa: F821
        cbsd: 'DBCbsd',  # noqa: F821
        serial_number: str,
):
    fcc_id, network_id, _ = get_cbsd_basic_params(cbsd)

    return DPLog(
        event_timestamp=now(),
        log_from='DP',
        log_to='CBSD',
        log_name=method_name + 'Response',
        log_message=str(result),
        cbsd_serial_number=str(serial_number),
        network_id=str(network_id),
        fcc_id=str(fcc_id),
    )
