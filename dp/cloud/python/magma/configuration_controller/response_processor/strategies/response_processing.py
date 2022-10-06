"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import logging
from datetime import datetime
from typing import Callable, Optional

from magma.configuration_controller.custom_types.custom_types import DBResponse
from magma.configuration_controller.response_processor.response_db_processor import (
    ResponseDBProcessor,
)
from magma.db_service.models import DBCbsd, DBCbsdState, DBGrant, DBGrantState
from magma.db_service.session_manager import Session
from magma.mappings.types import CbsdStates, GrantStates, ResponseCodes

logger = logging.getLogger(__name__)


CBSD_ID = "cbsdId"
GRANT_ID = "grantId"
GRANT_EXPIRE_TIME = "grantExpireTime"
HEARTBEAT_INTERVAL = "heartbeatInterval"
TRANSMIT_EXPIRE_TIME = "transmitExpireTime"
CHANNEL_TYPE = "channelType"
OPERATION_PARAM = "operationParam"


def unregister_cbsd_on_response_condition(process_response_func) -> Callable:
    """
    Unregister a CBSD on specific SAS response code conditions.

    This decorator is applied to any process response functions
    which should react to response codes that require Domain Proxy
    to internally unregister the CBSD.

    Currently a CBSD should be marked as unregistered on Domain Proxy if:
    * SAS returns a response with responseCode 105 (ResponseCodes.DEREGISTER)
    * SAS returns a response with responseCode 103 (ResponseCodes.INVALID_VALUE)

    Parameters:
        process_response_func: Response processing function

    Returns:
        response processing function wrapper
    """
    def process_response_wrapper(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
        if response.response_code in {ResponseCodes.DEREGISTER.value, ResponseCodes.INVALID_VALUE.value}:
            logger.info(f'SAS {response.payload} implies CBSD immedaite unregistration')
            _unregister_cbsd(response, session)
            return
        process_response_func(obj, response, session)

    return process_response_wrapper


@unregister_cbsd_on_response_condition
def process_registration_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process registration response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session
    """

    cbsd_id = response.payload.get("cbsdId", None)
    if response.response_code == ResponseCodes.SUCCESS.value and cbsd_id:
        _process_registration_response(cbsd_id, response, session)


def _process_registration_response(cbsd_id: str, response: DBResponse, session: Session):
    payload = response.request.payload

    where = _get_cbsd_filter(payload)
    if not where:
        return

    state_id = session.query(DBCbsdState.id). \
        filter(DBCbsdState.name == CbsdStates.REGISTERED.value)
    _update_cbsd(session, where, {"state_id": state_id.subquery(), "cbsd_id": cbsd_id})


def _get_cbsd_filter(payload):
    if "cbsdSerialNumber" in payload:
        return {"cbsd_serial_number": payload.get("cbsdSerialNumber")}
    elif "cbsdId" in payload:
        return {"cbsd_id": payload.get("cbsdId")}
    else:
        logger.warning(f'Could not find a CBSD identifier in {payload=}')
        return


def _update_cbsd(session, where_clause, update_clause):
    session.query(DBCbsd). \
        filter_by(**where_clause). \
        update(update_clause)


@unregister_cbsd_on_response_condition
def process_spectrum_inquiry_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process spectrum inquiry response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session
    """
    if response.response_code == ResponseCodes.SUCCESS.value:
        _create_channels(response, session)
        return
    logger.warning(f'process_spectrum_inquiry_response: Received an unsuccessful SAS response, {response.payload}=')


def _create_channels(response: DBResponse, session: Session):
    _terminate_all_grants_from_response(response, session)
    cbsd_id = response.request.payload["cbsdId"]
    cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
    available_channels = response.payload.get("availableChannel")
    cbsd.available_frequencies = None
    if not available_channels:
        logger.warning(
            "Could not create channel from spectrumInquiryResponse. Response missing 'availableChannel' object",
        )
        return

    channels = []
    for ac in available_channels:
        frequency_range = ac["frequencyRange"]
        channels.append({
            "low_frequency": frequency_range["lowFrequency"],
            "high_frequency": frequency_range["highFrequency"],
            "max_eirp": ac.get("maxEirp", 37),
        })
    cbsd.channels = channels


@unregister_cbsd_on_response_condition
def process_grant_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process grant response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session

    Returns:
        None
    """
    # Grant response codes worth considering here also are:
    # 400 - INTERFERENCE
    if response.response_code == ResponseCodes.GRANT_CONFLICT.value:
        _unsync_conflict_from_response(obj, response, session)
        return
    if response.response_code != ResponseCodes.SUCCESS.value:
        _remove_grant_from_response(response, session, unset_freq=True)
        return

    new_state = obj.grant_states_map[GrantStates.GRANTED.value]
    grant = _get_grant_from_response(response, session)
    if not grant:
        grant = _create_grant_from_response(response, new_state, session)
    else:
        logger.info(
            f'process_grant_responses: Updating grant state from {grant.state} to {new_state}',
        )
        grant.state = new_state
    _update_grant_from_response(response, grant)


@unregister_cbsd_on_response_condition
def process_heartbeat_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process heartbeat response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session

    Returns:
        None
    """
    if response.response_code == ResponseCodes.TERMINATED_GRANT.value:
        _remove_grant_from_response(response, session, unset_freq=True)
        return

    if response.response_code == ResponseCodes.SUCCESS.value:
        new_state = obj.grant_states_map[GrantStates.AUTHORIZED.value]
    elif response.response_code == ResponseCodes.SUSPENDED_GRANT.value:
        new_state = obj.grant_states_map[GrantStates.GRANTED.value]
    elif response.response_code == ResponseCodes.UNSYNC_OP_PARAM.value:
        new_state = obj.grant_states_map[GrantStates.UNSYNC.value]
    else:
        return

    grant = _get_grant_from_response(response, session)
    if not grant:
        grant = _create_grant_from_response(response, new_state, session)
    else:
        logger.info(
            f'process_heartbeat_responses: Updating grant state from {grant.state} to {new_state}',
        )
        grant.state = new_state

    _update_grant_from_response(response, grant)
    grant.last_heartbeat_request_time = datetime.now()


@unregister_cbsd_on_response_condition
def process_relinquishment_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process relinquishment response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session

    Returns:
        None
    """
    if response.response_code == ResponseCodes.SUCCESS.value:
        _remove_grant_from_response(response, session)
        return

    # FIXME This code is not unit-tested
    grant = _get_grant_from_response(response, session)
    _update_grant_from_response(response, grant)


def process_deregistration_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process deregistration response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session
    """

    logger.info(
        f'process_deregistration_response: Unregistering {response.payload}',
    )
    _unregister_cbsd(response, session)


def unset_frequency(grant: DBGrant):
    """
    Unset available frequency on the nth position of available frequencies for the given frequency

    Args:
        grant (DBGrant): Grant whose low and high frequencies are the base for the calculation

    Returns:
        None
    """
    frequencies = grant.cbsd.available_frequencies
    low = grant.low_frequency
    high = grant.high_frequency

    if not all([frequencies, low, high]):
        return

    bw_hz = high - low
    mid = (low + high) // 2
    bit_to_unset = (mid - int(3550 * 1e6)) // int(5 * 1e6)
    bw_index = bw_hz // int(5 * 1e6) - 1

    frequencies[bw_index] = frequencies[bw_index] & ~(1 << int(bit_to_unset))  # noqa: WPS465


def _get_grant_from_response(
        response: DBResponse,
        session: Session,
) -> Optional[DBGrant]:
    cbsd_id = response.cbsd_id
    grant_id = response.grant_id
    if not grant_id:
        return None

    grant = session.query(DBGrant). \
        join(DBCbsd). \
        filter(DBCbsd.cbsd_id == cbsd_id, DBGrant.grant_id == grant_id). \
        scalar()
    return grant


def _create_grant_from_response(
        response: DBResponse,
        state: DBGrantState,
        session: Session,
        grant_id: str = None,
) -> Optional[DBGrant]:
    grant_id = grant_id or response.grant_id
    if not grant_id:
        return None
    cbsd_id = session.query(DBCbsd.id).filter(DBCbsd.cbsd_id == response.cbsd_id)
    grant = DBGrant(cbsd_id=cbsd_id.subquery(), grant_id=grant_id, state=state)
    _update_grant_from_request(response, grant)
    session.add(grant)

    logger.info(f'Created new grant {grant}')
    return grant


def _remove_grant_from_response(
        response: DBResponse, session: Session, unset_freq: bool = False,
) -> None:
    grant = _get_grant_from_response(response, session)
    if not grant:
        return

    logger.info(f'Terminating grant {grant.grant_id}')

    if unset_freq:
        unset_frequency(grant)
    session.delete(grant)


def _update_grant_from_request(response: DBResponse, grant: DBGrant) -> None:
    payload = response.request.payload
    operation_param = payload.get(OPERATION_PARAM, {})
    frequency_range = operation_param.get("operationFrequencyRange", {})
    grant.max_eirp = operation_param.get("maxEirp", 0)
    grant.low_frequency = frequency_range.get("lowFrequency", 0)
    grant.high_frequency = frequency_range.get("highFrequency", 0)


def _update_grant_from_response(response: DBResponse, grant: DBGrant) -> None:
    if not grant:
        return
    grant_expire_time = response.payload.get(GRANT_EXPIRE_TIME, None)
    heartbeat_interval = response.payload.get(HEARTBEAT_INTERVAL, None)
    transmit_expire_time = response.payload.get(TRANSMIT_EXPIRE_TIME, None)
    channel_type = response.payload.get(CHANNEL_TYPE, None)
    if grant_expire_time:
        grant.grant_expire_time = grant_expire_time
    if heartbeat_interval:
        grant.heartbeat_interval = int(heartbeat_interval)
    if transmit_expire_time:
        grant.transmit_expire_time = transmit_expire_time
    if channel_type:
        grant.channel_type = channel_type
    logger.info(f'Updated grant: {grant}')


def _terminate_all_grants_from_response(response: DBResponse, session: Session) -> None:
    cbsd_id = response.cbsd_id
    if not cbsd_id:
        return

    with session.no_autoflush:
        logger.info(f'Terminating all grants for {cbsd_id=}')
        session.query(DBGrant).filter(
            DBGrant.cbsd_id == DBCbsd.id, DBCbsd.cbsd_id == cbsd_id,
        ).delete(synchronize_session=False)

        logger.info(f"Deleting all channels for {cbsd_id=}")
        session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).update({DBCbsd.channels: []})


def _unsync_conflict_from_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    state = obj.grant_states_map[GrantStates.UNSYNC.value]

    conflicts_ids = response.payload.get("response", {}).get("responseData", [])
    existing_grants = session.query(DBGrant.grant_id).filter(DBGrant.grant_id.in_(conflicts_ids)).all()
    existing_grants_ids = {g.grant_id for g in existing_grants}

    for grant_id in conflicts_ids:
        if grant_id in existing_grants_ids:
            continue

        _create_grant_from_response(response, state, session, grant_id=grant_id)
        return


def _unregister_cbsd(response: DBResponse, session: Session) -> None:
    payload = response.request.payload
    where = _get_cbsd_filter(payload)
    if not where:
        return
    state_id = session.query(DBCbsdState.id). \
        filter(DBCbsdState.name == CbsdStates.UNREGISTERED.value)
    _terminate_all_grants_from_response(response, session)
    _update_cbsd(session, where, {"state_id": state_id.subquery()})
