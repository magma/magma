import json
import logging

from magma.db_service.models import (
    DBActiveModeConfig,
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBGrant,
    DBGrantState,
    DBRequest,
    DBRequestState,
)
from magma.db_service.session_manager import Session, SessionManager
from magma.mappings.cbsd_states import (
    cbsd_state_mapping,
    grant_state_mapping,
    switch_mapping,
)
from magma.mappings.types import GrantStates, RequestStates, Switch
from dp.protos.active_mode_pb2 import (
    _CBSDSTATE,
    ActiveModeConfig,
    Cbsd,
    Channel,
    GetStateRequest,
    Grant,
    State,
    ToggleActiveModeParams,
    ToggleActiveModeResponse,
)
from dp.protos.active_mode_pb2_grpc import ActiveModeControllerServicer
from dp.protos.common_pb2 import FrequencyRange
from sqlalchemy.orm import joinedload

logger = logging.getLogger(__name__)


class ActiveModeControllerService(ActiveModeControllerServicer):
    def __init__(self, session_manager: SessionManager):
        self.session_manager = session_manager

    def GetState(self, request: GetStateRequest, context) -> State:
        logger.info("Getting DB state")
        with self.session_manager.session_scope() as session:
            return self._build_state(session)

    def ToggleActiveMode(self, request: ToggleActiveModeParams, context) -> ToggleActiveModeResponse:
        with self.session_manager.session_scope() as session:
            return self._toggle_active_mode(session, request)

    def _build_state(self, session: Session) -> State:
        db_configs = session.query(DBActiveModeConfig) \
            .join(DBCbsd).options(
            joinedload(DBActiveModeConfig.cbsd).options(
                joinedload(DBCbsd.channels),
                joinedload(DBCbsd.grants).options(joinedload(DBGrant.state)),
            )
        )
        configs = [self._build_config(session, x) for x in db_configs]
        session.commit()
        return State(active_mode_configs=configs)

    def _build_config(self, session: Session, config: DBActiveModeConfig) -> ActiveModeConfig:
        return ActiveModeConfig(
            desired_state=cbsd_state_mapping[config.desired_state.name],
            cbsd=self._build_cbsd(session, config.cbsd),
        )

    def _build_cbsd(self, session: Session, cbsd: DBCbsd) -> Cbsd:
        db_grants = session.query(DBGrant) \
            .join(DBGrantState) \
            .filter(
                DBGrant.cbsd_id == cbsd.id,
                DBGrantState.name != GrantStates.IDLE.value,
            )
        pending_requests_payloads = session.query(DBRequest.payload) \
            .join(DBRequestState) \
            .filter(DBRequestState.name == RequestStates.PENDING.value, DBRequest.cbsd_id == cbsd.id)
        grants = [self._build_grant(x) for x in db_grants]
        channels = [self._build_channel(x) for x in cbsd.channels]
        pending_requests = [json.dumps(payload, separators=(',', ':')) for (payload,) in pending_requests_payloads]
        return Cbsd(
            id=cbsd.cbsd_id,
            user_id=cbsd.user_id,
            fcc_id=cbsd.fcc_id,
            serial_number=cbsd.cbsd_serial_number,
            state=cbsd_state_mapping[cbsd.state.name],
            eirp_capability=cbsd.eirp_capability,
            grants=grants,
            channels=channels,
            pending_requests=pending_requests
        )

    @staticmethod
    def _build_grant(grant: DBGrant) -> Grant:
        last_heartbeat = grant.last_heartbeat_request_time
        if last_heartbeat:
            last_heartbeat_timestamp = last_heartbeat.timestamp()
        else:
            last_heartbeat_timestamp = 0

        return Grant(
            id=grant.grant_id,
            state=grant_state_mapping[grant.state.name],
            heartbeat_interval_sec=grant.heartbeat_interval,
            last_heartbeat_timestamp=int(last_heartbeat_timestamp),
        )

    @staticmethod
    def _build_channel(channel: DBChannel) -> Channel:
        return Channel(
            frequency_range=FrequencyRange(
                low=channel.low_frequency,
                high=channel.high_frequency,
            ),
            max_eirp=channel.max_eirp,
            last_eirp=channel.last_used_max_eirp,
        )

    def _toggle_active_mode(self, session: Session, params: ToggleActiveModeParams):
        cbsd_id = params.cbsd_id
        cbsd = session.query(DBCbsd).filter(DBCbsd.id == cbsd_id).first()
        if params.switch == switch_mapping[Switch.OFF.value]:
            logger.info(f"Switching active mode off for {cbsd}")
            return self._switch_active_mode_off(session, cbsd_id)
        else:
            return self._create_or_update_active_mode_config(
                session,
                cbsd_id,
                _CBSDSTATE.values_by_number[params.desired_state].name.lower())

    @staticmethod
    def _switch_active_mode_off(session: Session, cbsd_id: int):
        session.query(DBActiveModeConfig).filter(DBActiveModeConfig.cbsd_id == cbsd_id).delete()
        session.commit()
        return ToggleActiveModeResponse(
            cbsd_db_id=cbsd_id,
            active_mode_enabled=False,
        )

    @staticmethod
    def _create_or_update_active_mode_config(
            session: Session, cbsd_id: int, new_state_name: str) -> ToggleActiveModeResponse:

        cbsd_state = session.query(DBCbsdState).filter(DBCbsdState.name == new_state_name).first()

        active_mode_config = session.query(DBActiveModeConfig).options(
            joinedload(DBActiveModeConfig.cbsd).options(
                joinedload(DBCbsd.state)
            )
        ).filter(DBActiveModeConfig.cbsd_id == cbsd_id) \
            .first()

        if active_mode_config:
            active_mode_config.cbsd.state = cbsd_state
        else:
            active_mode_config = DBActiveModeConfig(cbsd_id=cbsd_id, desired_state=cbsd_state)
            session.add(active_mode_config)
        session.commit()
        return ToggleActiveModeResponse(
            cbsd_db_id=cbsd_id,
            active_mode_enabled=True,
            cbsd_state=new_state_name.title()
        )
