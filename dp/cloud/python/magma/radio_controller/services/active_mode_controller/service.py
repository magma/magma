import json
import logging
from datetime import datetime
from typing import Optional

from magma.db_service.models import (
    DBActiveModeConfig,
    DBCbsd,
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
)
from magma.mappings.types import GrantStates, RequestStates
from dp.protos.active_mode_pb2 import (
    ActiveModeConfig,
    Cbsd,
    Channel,
    FrequencyRange,
    GetStateRequest,
    Grant,
    State,
)
from dp.protos.active_mode_pb2_grpc import ActiveModeControllerServicer
from sqlalchemy.orm import joinedload

logger = logging.getLogger(__name__)


class ActiveModeControllerService(ActiveModeControllerServicer):
    def __init__(self, session_manager: SessionManager):
        self.session_manager = session_manager

    def GetState(self, request: GetStateRequest, context) -> State:
        logger.info("Getting DB state")
        with self.session_manager.session_scope() as session:
            return self._build_state(session)

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
        last_seen = self._to_timestamp(cbsd.last_seen)
        return Cbsd(
            id=cbsd.cbsd_id,
            user_id=cbsd.user_id,
            fcc_id=cbsd.fcc_id,
            serial_number=cbsd.cbsd_serial_number,
            state=cbsd_state_mapping[cbsd.state.name],
            eirp_capability=cbsd.eirp_capability,
            grants=grants,
            channels=channels,
            pending_requests=pending_requests,
            last_seen_timestamp=last_seen,
        )

    def _build_grant(self, grant: DBGrant) -> Grant:
        last_heartbeat = self._to_timestamp(grant.last_heartbeat_request_time)
        return Grant(
            id=grant.grant_id,
            state=grant_state_mapping[grant.state.name],
            heartbeat_interval_sec=grant.heartbeat_interval,
            last_heartbeat_timestamp=last_heartbeat,
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

    @staticmethod
    def _to_timestamp(t: Optional[datetime]) -> int:
        return 0 if t is None else int(t.timestamp())
