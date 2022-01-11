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
from typing import Dict, Optional

from magma.configuration_controller.response_processor.response_db_processor import (
    ResponseDBProcessor,
)
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBGrant,
    DBResponse,
)
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


def process_registration_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process registration response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session
    """

    cbsd_id = response.payload.get("cbsdId", None)
    if response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(response, session)
    elif response.response_code == ResponseCodes.SUCCESS.value and cbsd_id:
        payload = response.request.payload
        cbsd = _find_cbsd_from_registration_request(session, payload)
        cbsd.cbsd_id = cbsd_id
        _change_cbsd_state(cbsd, session, CbsdStates.REGISTERED.value)


def _find_cbsd_from_registration_request(session: Session, payload: Dict) -> DBCbsd:
    return session.query(DBCbsd).filter(
        DBCbsd.cbsd_serial_number == payload["cbsdSerialNumber"],
    ).scalar()


def _change_cbsd_state(cbsd: DBCbsd, session: Session, new_state: str) -> None:
    state = session.query(DBCbsdState).filter(
        DBCbsdState.name == new_state,
    ).scalar()
    cbsd.state = state


def process_spectrum_inquiry_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process spectrum inquiry response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session
    """

    if response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(response, session)
    elif response.response_code == ResponseCodes.SUCCESS.value:
        _create_channels(response, session)


def _create_channels(response: DBResponse, session: Session):
    _terminate_all_grants_from_response(response, session)
    cbsd_id = response.request.payload["cbsdId"]
    cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
    available_channels = response.payload.get("availableChannel")
    if not available_channels:
        logger.warning(
            "Could not create channel from spectrumInquiryResponse. Response missing 'availableChannel' object",
        )
        return
    for ac in available_channels:
        frequency_range = ac["frequencyRange"]
        channel = DBChannel(
            cbsd=cbsd,
            low_frequency=frequency_range["lowFrequency"],
            high_frequency=frequency_range["highFrequency"],
            channel_type=ac["channelType"],
            rule_applied=ac["ruleApplied"],
            max_eirp=ac.get("maxEirp"),
        )
        logger.info(f"Creating channel for {cbsd=}")
        session.add(channel)


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

    grant = _get_or_create_grant_from_response(obj, response, session)
    if not grant:
        return
    _update_grant_from_response(response, grant)
    channel = _get_channel_related_to_grant(response, session)
    if channel:
        channel.last_used_max_eirp = response.request.payload[OPERATION_PARAM]["maxEirp"]
        grant.channel = channel

    # Grant response codes worth considering here also are:
    # 400 - INTERFERENCE
    # 401 - GRANT_CONFLICT
    # Might need better processing, for now we set the state to IDLE in all cases other than SUCCESS
    if response.response_code == ResponseCodes.SUCCESS.value:
        new_state = obj.grant_states_map[GrantStates.GRANTED.value]
    else:
        new_state = obj.grant_states_map[GrantStates.IDLE.value]
    logger.info(
        f'process_grant_responses: Updating grant state from {grant.state} to {new_state}',
    )
    grant.state = new_state


def _get_channel_related_to_grant(response: DBResponse, session: Session) -> DBChannel:
    payload = response.request.payload
    operation_param = payload[OPERATION_PARAM]
    frequency_range = operation_param["operationFrequencyRange"]
    channel = session.query(DBChannel).join(DBCbsd).filter(
        DBCbsd.cbsd_id == payload["cbsdId"],
        DBChannel.low_frequency == frequency_range["lowFrequency"],
        DBChannel.high_frequency == frequency_range["highFrequency"],
    ).scalar()
    return channel


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

    grant = _get_or_create_grant_from_response(obj, response, session)
    if not grant:
        return
    _update_grant_from_response(response, grant)

    if response.response_code == ResponseCodes.SUCCESS.value:
        new_state = obj.grant_states_map[GrantStates.AUTHORIZED.value]
    elif response.response_code == ResponseCodes.SUSPENDED_GRANT.value:
        new_state = obj.grant_states_map[GrantStates.GRANTED.value]
    elif response.response_code in {ResponseCodes.TERMINATED_GRANT.value, ResponseCodes.UNSYNC_OP_PARAM.value}:
        new_state = obj.grant_states_map[GrantStates.IDLE.value]
        _reset_last_used_max_eirp(grant)
    elif response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(response, session)
        return
    else:
        new_state = grant.state
    logger.info(
        f'process_heartbeat_responses: Updating grant state from {grant.state} to {new_state}',
    )
    grant.state = new_state
    grant.last_heartbeat_request_time = datetime.now()


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

    grant = _get_or_create_grant_from_response(obj, response, session)
    if not grant:
        return
    _update_grant_from_response(response, grant)

    if response.response_code == ResponseCodes.SUCCESS.value:
        new_state = obj.grant_states_map[GrantStates.IDLE.value]
        _reset_last_used_max_eirp(grant)
    elif response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(response, session)
        return
    else:
        new_state = grant.state
    logger.info(
        f'process_relinquishment_responses: Updating grant state from {grant.state} to {new_state}',
    )
    grant.state = new_state


def process_deregistration_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    """
    Process deregistration response

    Parameters:
        obj: Response processor object
        response: Database response object
        session: Database session
    """

    _terminate_all_grants_from_response(response, session)
    cbsd_id = response.payload.get("cbsdId", None)
    if cbsd_id:
        cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
        _change_cbsd_state(cbsd, session, CbsdStates.UNREGISTERED.value)


def _reset_last_used_max_eirp(grant: DBGrant) -> None:
    if grant and grant.channel:
        grant.channel.last_used_max_eirp = None


def _get_or_create_grant_from_response(
    obj: ResponseDBProcessor,
    response: DBResponse,
    session: Session,
) -> Optional[DBGrant]:
    cbsd_id = response.payload.get(
        CBSD_ID,
    ) or response.request.payload.get(CBSD_ID)
    grant_id = response.payload.get(
        GRANT_ID,
    ) or response.request.payload.get(GRANT_ID)
    cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
    grant = None
    if grant_id:
        logger.info(f'Getting grant by: {cbsd_id=} {grant_id=}')
        grant = session.query(DBGrant).filter(
            DBGrant.cbsd_id == cbsd.id, DBGrant.grant_id == grant_id,
        ).scalar()

    if grant_id and not grant:
        grant_idle_state = obj.grant_states_map[GrantStates.IDLE.value]
        grant = DBGrant(cbsd=cbsd, grant_id=grant_id, state=grant_idle_state)
        session.add(grant)
        logger.info(f'Created new grant: {grant}')
    return grant


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
    grant.responses.append(response)
    logger.info(f'Updated grant: {grant}')


def _terminate_all_grants_from_response(response: DBResponse, session: Session) -> None:
    cbsd_id = response.payload.get(
        CBSD_ID,
    ) or response.request.payload.get(CBSD_ID)
    if not cbsd_id:
        return
    cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
    grant_ids = [
        grant.id for grant in session.query(
            DBGrant.id,
        ).filter(DBGrant.cbsd == cbsd).all()
    ]
    if grant_ids:
        logger.info(f'Disconnecting responses from grants for {cbsd_id=}')
        session.query(DBResponse).filter(
            DBResponse.grant_id.in_(
                grant_ids,
            ),
        ).update({DBResponse.grant_id: None})
    logger.info(f'Terminating all grants for {cbsd_id=}')
    session.query(DBGrant).filter(DBGrant.cbsd == cbsd).delete()
    logger.info(f"Deleting all channels for {cbsd_id=}")
    session.query(DBChannel).filter(DBChannel.cbsd == cbsd).delete()
