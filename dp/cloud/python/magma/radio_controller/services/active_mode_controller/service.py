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

import json
import logging
from datetime import datetime
from typing import Optional

from dp.protos.active_mode_pb2 import (
    ActiveModeConfig,
    Cbsd,
    Channel,
    EirpCapabilities,
    FrequencyRange,
    GetStateRequest,
    Grant,
    State,
)
from dp.protos.active_mode_pb2_grpc import ActiveModeControllerServicer
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
from magma.mappings.cbsd_states import cbsd_state_mapping, grant_state_mapping
from magma.mappings.types import GrantStates, RequestStates
from sqlalchemy.orm import joinedload

logger = logging.getLogger(__name__)


class ActiveModeControllerService(ActiveModeControllerServicer):
    """
    Active Mode Controller gRPC Service class
    """

    def __init__(self, session_manager: SessionManager):
        self.session_manager = session_manager

    def GetState(self, request: GetStateRequest, context) -> State:
        """
        Get Active Mode Database state from the Database

        Parameters:
            request: a GetStateRequest gRPC Message
            context: gRPC context

        Returns:
            State: a gRPC State message
        """
        logger.info("Getting DB state")
        with self.session_manager.session_scope() as session:
            state = self._build_state(session)
            logger.info(f"Sending state: {state}")
            return state

    def _build_state(self, session: Session) -> State:
        db_grant_idle_state_id = session.query(DBGrantState.id).filter(
            DBGrantState.name == GrantStates.IDLE.value,
        ).scalar()
        db_request_pending_state_id = session.query(DBRequestState.id).filter(
            DBRequestState.name == RequestStates.PENDING.value,
        ).scalar()

        # Selectively load sqlalchemy object relations using a single query to avoid commit races.
        # We want to have CBSD entity "grants" relation only contain grants in a Non-IDLE state.
        # We want to have CBSD entity "requests" relation only contain PENDING requests.
        db_configs = session.query(DBActiveModeConfig).join(DBCbsd).options(
            joinedload(DBActiveModeConfig.cbsd).options(
                joinedload(DBCbsd.channels),
                joinedload(
                    DBCbsd.grants.and_(
                        DBGrant.state_id != db_grant_idle_state_id,
                    ),
                ),
                joinedload(
                    DBCbsd.requests.and_(
                        DBRequest.state_id == db_request_pending_state_id,
                    ),
                ),
            ),
        ).filter(*self._get_filter()).populate_existing()
        configs = [self._build_config(db_config) for db_config in db_configs]
        session.commit()
        return State(active_mode_configs=configs)

    def _get_filter(self):
        not_null_fields = [
            DBCbsd.fcc_id, DBCbsd.user_id, DBCbsd.number_of_ports,
            DBCbsd.antenna_gain, DBCbsd.min_power, DBCbsd.max_power,
        ]
        return [field != None for field in not_null_fields]  # noqa: E711

    def _build_config(self, config: DBActiveModeConfig) -> ActiveModeConfig:
        return ActiveModeConfig(
            desired_state=cbsd_state_mapping[config.desired_state.name],
            cbsd=self._build_cbsd(config.cbsd),
        )

    def _build_cbsd(self, cbsd: DBCbsd) -> Cbsd:
        pending_requests = [
            json.dumps(r.payload, separators=(',', ':')) for r in cbsd.requests
        ]

        # Application may not need those to be sorted.
        # Applying ordering mostly for easier assertions in testing
        cbsd_db_grants = sorted(cbsd.grants, key=lambda x: x.id)
        cbsd_db_channels = sorted(cbsd.channels, key=lambda x: x.id)

        grants = [self._build_grant(x) for x in cbsd_db_grants]
        channels = [self._build_channel(x) for x in cbsd_db_channels]

        last_seen = self._to_timestamp(cbsd.last_seen)
        eirp_capabilities = self._build_eirp_capabilities(cbsd)
        return Cbsd(
            id=cbsd.cbsd_id,
            user_id=cbsd.user_id,
            fcc_id=cbsd.fcc_id,
            serial_number=cbsd.cbsd_serial_number,
            state=cbsd_state_mapping[cbsd.state.name],
            grants=grants,
            channels=channels,
            pending_requests=pending_requests,
            last_seen_timestamp=last_seen,
            eirp_capabilities=eirp_capabilities,
        )

    def _build_grant(self, grant: DBGrant) -> Grant:
        last_heartbeat = self._to_timestamp(grant.last_heartbeat_request_time)
        return Grant(
            id=grant.grant_id,
            state=grant_state_mapping[grant.state.name],
            heartbeat_interval_sec=grant.heartbeat_interval,
            last_heartbeat_timestamp=last_heartbeat,
        )

    def _build_channel(self, channel: DBChannel) -> Channel:
        return Channel(
            frequency_range=FrequencyRange(
                low=channel.low_frequency,
                high=channel.high_frequency,
            ),
            max_eirp=channel.max_eirp,
            last_eirp=channel.last_used_max_eirp,
        )

    def _build_eirp_capabilities(self, cbsd: DBCbsd) -> EirpCapabilities:
        return EirpCapabilities(
            min_power=cbsd.min_power,
            max_power=cbsd.max_power,
            antenna_gain=cbsd.antenna_gain,
            number_of_ports=cbsd.number_of_ports,
        )

    def _to_timestamp(self, t: Optional[datetime]) -> int:
        return 0 if t is None else int(t.timestamp())
