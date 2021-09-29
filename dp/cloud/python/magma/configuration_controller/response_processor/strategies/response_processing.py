import logging
from datetime import datetime, timezone
from typing import Optional, Dict

from magma.configuration_controller.response_processor.response_db_processor import (
    ResponseDBProcessor,
)
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBGrant,
    DBRequest,
    DBResponse,
)
from magma.db_service.session_manager import Session
from magma.mappings.types import (
    CbsdStates,
    GrantStates,
    ResponseCodes,
)

logger = logging.getLogger(__name__)


CBSD_ID = "cbsdId"
GRANT_ID = "grantId"
GRANT_EXPIRE_TIME = "grantExpireTime"
HEARTBEAT_INTERVAL = "heartbeatInterval"
TRANSMIT_EXPIRE_TIME = "transmitExpireTime"
CHANNEL_TYPE = "channelType"
OPERATION_PARAM = "operationParam"


def process_registration_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    cbsd_id = response.payload.get("cbsdId", None)
    if response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(obj, response, session)
    elif response.response_code == ResponseCodes.SUCCESS.value and cbsd_id:
        payload = response.request.payload
        cbsd = _find_cbsd_from_registration_request(session, payload)
        cbsd.cbsd_id = cbsd_id
        _change_cbsd_state(cbsd, session, CbsdStates.REGISTERED.value)


def _find_cbsd_from_registration_request(session: Session, payload: Dict) -> DBCbsd:
    return session.query(DBCbsd).filter(
        DBCbsd.fcc_id == payload["fccId"],
        DBCbsd.user_id == payload["userId"],
        DBCbsd.cbsd_serial_number == payload["cbsdSerialNumber"],
    ).scalar()


def _change_cbsd_state(cbsd: DBCbsd, session: Session, new_state: str) -> None:
    state = session.query(DBCbsdState).filter(DBCbsdState.name == new_state).scalar()
    cbsd.state = state


def process_spectrum_inquiry_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    if response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(obj, response, session)
    elif response.response_code == ResponseCodes.SUCCESS.value:
        _create_channels(response, session)


def _create_channels(response: DBResponse, session: Session):
    cbsd_id = response.request.payload["cbsdId"]
    cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
    logger.info(f"Deleting all channels for {cbsd}")
    session.query(DBChannel).filter(DBChannel.cbsd == cbsd).delete()
    available_channels = response.payload.get("availableChannel")
    if not available_channels:
        logger.warning(
            "Could not create channel from spectrumInquiryResponse. Response missing 'availableChannel' object")
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
        logger.info(f"Creating channel for {cbsd}")
        session.add(channel)


def process_grant_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    grant = _get_or_create_grant_from_response(obj, response, session)
    _update_grant_from_response(response, grant)
    channel = _get_channel_related_to_grant(response, session)
    if channel:
        channel.last_used_max_eirp = response.request.payload[OPERATION_PARAM]["maxEirp"]
        grant.channel = channel

    # Grant response codes worth considering here also are:
    # 400 - INTERFERENCE
    # 401 - GRANT_CONFLICT
    # Might need better processing, for now we set the state to IDLE in all cases other than 0
    if response.response_code == ResponseCodes.SUCCESS.value:
        new_state = obj.grant_states_map[GrantStates.GRANTED.value]
    else:
        new_state = obj.grant_states_map[GrantStates.IDLE.value]
    logger.info(f'process_grant_responses: Updating grant state from {grant.state} to {new_state}')
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
    grant = _get_or_create_grant_from_response(obj, response, session)
    logger.info(f'Processing grant: {grant}')
    _update_grant_from_response(response, grant)

    if response.response_code == ResponseCodes.SUCCESS.value:
        new_state = obj.grant_states_map[GrantStates.AUTHORIZED.value]
    elif response.response_code == ResponseCodes.SUSPENDED_GRANT.value:
        new_state = obj.grant_states_map[GrantStates.GRANTED.value]
    elif response.response_code in [ResponseCodes.TERMINATED_GRANT.value, ResponseCodes.UNSYNC_OP_PARAM.value]:
        new_state = obj.grant_states_map[GrantStates.IDLE.value]
    elif response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(obj, response, session)
        return
    else:
        new_state = grant.state
    logger.info(f'process_heartbeat_responses: Updating grant state from {grant.state} to {new_state}')
    grant.state = new_state
    grant.last_heartbeat_request_time = datetime.now(timezone.utc)  # TODO use different timezone?


def process_relinquishment_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    grant = _get_or_create_grant_from_response(obj, response, session)
    _update_grant_from_response(response, grant)

    if response.response_code == ResponseCodes.SUCCESS.value:
        new_state = obj.grant_states_map[GrantStates.IDLE.value]
    elif response.response_code == ResponseCodes.DEREGISTER.value:
        _terminate_all_grants_from_response(obj, response, session)
        return
    else:
        new_state = grant.state
    logger.info(f'process_relinquishment_responses: Updating grant state from {grant.state} to {new_state}')
    grant.state = new_state


def process_deregistration_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    # Is this a bug ?
    # according to documentation 8.8.2
    # cbsd should be unregistered regardless of response code
    if response.response_code in [ResponseCodes.SUCCESS.value, ResponseCodes.DEREGISTER.value]:
        _terminate_all_grants_from_response(obj, response, session)
    cbsd_id = response.payload.get("cbsdId", None)
    if cbsd_id:
        cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
        _change_cbsd_state(cbsd, session, CbsdStates.UNREGISTERED.value)


def _get_or_create_grant_from_response(obj: ResponseDBProcessor,
                                       response: DBResponse,
                                       session: Session) -> Optional[DBGrant]:
    cbsd_id = response.payload[CBSD_ID]
    grant_id = response.payload[GRANT_ID]
    cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
    logger.info(f'Getting grant by cbsd_id={cbsd_id} and grant_id={grant_id}')
    grant = session.query(DBGrant).filter(DBGrant.cbsd_id == cbsd.id, DBGrant.grant_id == grant_id).scalar()

    if not grant:
        grant_idle_state = obj.grant_states_map[GrantStates.IDLE.value]
        grant = DBGrant(cbsd=cbsd, grant_id=grant_id, state=grant_idle_state)
        session.add(grant)
        logger.info(f'Created new grant: {grant}')
    return grant


def _update_grant_from_response(response: DBResponse, grant: DBGrant) -> None:
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


def _terminate_all_grants_from_response(obj: ResponseDBProcessor, response: DBResponse, session: Session) -> None:
    cbsd_id = response.payload[CBSD_ID]
    cbsd = session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).scalar()
    logger.info(f'Terminating all grants for cbsd_id: {cbsd_id}')
    grant_idle_state = obj.grant_states_map[GrantStates.IDLE.value]
    for grant in session.query(DBGrant).filter(DBGrant.cbsd == cbsd).all():
        logger.info(f'Terminating grant {grant}')
        grant.state = grant_idle_state
