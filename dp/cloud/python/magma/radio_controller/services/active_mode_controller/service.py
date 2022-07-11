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
from typing import Any, List, Optional

import grpc
from dp.protos.active_mode_pb2 import (
    AcknowledgeCbsdUpdateRequest,
    Cbsd,
    Channel,
    DatabaseCbsd,
    DeleteCbsdRequest,
    EirpCapabilities,
    FrequencyPreferences,
    GetStateRequest,
    Grant,
    InstallationParams,
    SasSettings,
    State,
)
from dp.protos.active_mode_pb2_grpc import ActiveModeControllerServicer
from google.protobuf.empty_pb2 import Empty
from google.protobuf.wrappers_pb2 import FloatValue
from magma.db_service.models import (
    DBCbsd,
    DBChannel,
    DBGrant,
    DBGrantState,
    DBRequest,
)
from magma.db_service.session_manager import Session, SessionManager
from magma.mappings.cbsd_states import cbsd_state_mapping, grant_state_mapping
from magma.mappings.types import GrantStates
from magma.radio_controller.metrics import (
    ACKNOWLEDGE_UPDATE_PROCESSING_TIME,
    DELETE_CBSD_PROCESSING_TIME,
    GET_DB_STATE_PROCESSING_TIME,
)
from sqlalchemy import and_, or_
from sqlalchemy.orm import contains_eager, joinedload

logger = logging.getLogger(__name__)


class ActiveModeControllerService(ActiveModeControllerServicer):
    """
    Active Mode Controller gRPC Service class
    """

    def __init__(self, session_manager: SessionManager):
        self.session_manager = session_manager

    @GET_DB_STATE_PROCESSING_TIME.time()
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
            cbsds = _list_cbsds(session)
            state = _build_state(cbsds)
            session.commit()
            logger.debug(f"Sending state: {state}")
            return state

    @DELETE_CBSD_PROCESSING_TIME.time()
    def DeleteCbsd(self, request: DeleteCbsdRequest, context) -> Empty:
        """
        Delete CBSD from the Database

        Parameters:
            request: a DeleteCbsdRequest gRPC Message
            context: gRPC context

        Returns:
            Empty: an empty gRPC message
        """
        db_id = request.id
        logger.info(f"Deleting CBSD {db_id}")
        with self.session_manager.session_scope() as session:
            deleted = session.query(DBCbsd).filter(
                DBCbsd.id == db_id,
            ).delete()
            session.commit()
            if not deleted:
                context.set_code(grpc.StatusCode.NOT_FOUND)
        return Empty()

    @ACKNOWLEDGE_UPDATE_PROCESSING_TIME.time()
    def AcknowledgeCbsdUpdate(self, request: AcknowledgeCbsdUpdateRequest, context) -> Empty:
        """
        Mark CBSD in the Database as not updated

        Parameters:
            request: a AcknowledgeCbsdUpdateRequest gRPC Message
            context: gRPC context

        Returns:
            Empty: an empty gRPC message
        """
        db_id = request.id
        logger.info(f"Acknowledging CBSD update {db_id}")
        with self.session_manager.session_scope() as session:
            updated = session.query(DBCbsd).filter(
                DBCbsd.id == db_id,
            ).update({'should_deregister': False})
            session.commit()
            if not updated:
                context.set_code(grpc.StatusCode.NOT_FOUND)
        return Empty()

    def StoreAvailableFrequencies(self, request, context) -> Empty:
        """
        Store available frequencies in the database

        Parameters
            request: StoreAvailableFrequencies gRPC Message
            context: gRPC context

        Returns:
            Empty: an empty gRPC message
        """
        # Not implemented yet
        pass


def _list_cbsds(session: Session) -> State:
    # It might be possible to use join instead of nested queries
    # however it requires some serious investigation on how to use it
    # with eager_contains and filter (aliases)
    db_grant_idle_state_id = session.query(DBGrantState.id).filter(
        DBGrantState.name == GrantStates.IDLE.value,
    ).scalar_subquery()

    # Selectively load sqlalchemy object relations using a single query to avoid commit races.
    # We want to have CBSD entity "grants" relation only contain grants in a Non-IDLE state.
    # We want to have CBSD entity "requests" relation only contain PENDING requests.
    return (
        session.query(DBCbsd).
        outerjoin(
            DBGrant, and_(
                DBGrant.cbsd_id == DBCbsd.id,
                DBGrant.state_id != db_grant_idle_state_id,
            ),
        ).
        outerjoin(DBRequest).
        options(
            joinedload(DBCbsd.state),
            joinedload(DBCbsd.desired_state),
            joinedload(DBCbsd.channels),
            contains_eager(DBCbsd.grants).
            joinedload(DBGrant.state),
        ).
        filter(_build_filter()).
        populate_existing()
    )


def _build_filter():
    multi_step = [
        DBCbsd.fcc_id, DBCbsd.user_id, DBCbsd.number_of_ports,
        DBCbsd.min_power, DBCbsd.max_power, DBCbsd.antenna_gain,
    ]
    single_step = [
        DBCbsd.latitude_deg, DBCbsd.longitude_deg, DBCbsd.height_m,
        DBCbsd.height_type, DBCbsd.indoor_deployment,
    ]
    return and_(
        DBRequest.id == None,  # noqa: E711
        or_(
            or_(
                DBCbsd.should_deregister == True,
                DBCbsd.is_deleted == True,
            ),
            and_(
                DBCbsd.single_step_enabled == False,
                not_null(multi_step),
            ),
            and_(
                DBCbsd.single_step_enabled == True,
                DBCbsd.cbsd_category == 'a',
                DBCbsd.indoor_deployment == True,
                not_null(multi_step + single_step),
            ),
        ),
    )


def not_null(fields: List[Any]):
    return and_(*[field != None for field in fields])  # noqa: E711


def _build_state(db_cbsds: List[DBCbsd]) -> State:
    cbsds = [_build_cbsd(db_cbsd) for db_cbsd in db_cbsds]
    return State(cbsds=cbsds)


def _build_cbsd(cbsd: DBCbsd) -> Cbsd:
    # Application may not need those to be sorted.
    # Applying ordering mostly for easier assertions in testing
    cbsd_db_grants = sorted(cbsd.grants, key=lambda x: x.id)
    cbsd_db_channels = sorted(cbsd.channels, key=lambda x: x.id)

    grants = [_build_grant(x) for x in cbsd_db_grants]
    channels = [_build_channel(x) for x in cbsd_db_channels]

    last_seen = _to_timestamp(cbsd.last_seen)
    eirp_capabilities = _build_eirp_capabilities(cbsd)
    preferences = _build_preferences(cbsd)
    sas_settings = _build_sas_settings(cbsd)
    installation_params = _build_installation_params(cbsd)
    db_data = _build_db_data(cbsd)
    return Cbsd(
        cbsd_id=cbsd.cbsd_id,
        state=cbsd_state_mapping[cbsd.state.name],
        desired_state=cbsd_state_mapping[cbsd.desired_state.name],
        grants=grants,
        channels=channels,
        last_seen_timestamp=last_seen,
        eirp_capabilities=eirp_capabilities,
        grant_attempts=cbsd.grant_attempts,
        db_data=db_data,
        preferences=preferences,
        sas_settings=sas_settings,
        installation_params=installation_params,
    )


def _build_grant(grant: DBGrant) -> Grant:
    last_heartbeat = _to_timestamp(grant.last_heartbeat_request_time)
    return Grant(
        id=grant.grant_id,
        state=grant_state_mapping[grant.state.name],
        heartbeat_interval_sec=grant.heartbeat_interval,
        last_heartbeat_timestamp=last_heartbeat,
    )


def _build_channel(channel: DBChannel) -> Channel:
    return Channel(
        low_frequency_hz=channel.low_frequency,
        high_frequency_hz=channel.high_frequency,
        max_eirp=_make_optional_float(channel.max_eirp),
    )


def _build_eirp_capabilities(cbsd: DBCbsd) -> EirpCapabilities:
    return EirpCapabilities(
        min_power=cbsd.min_power,
        max_power=cbsd.max_power,
        number_of_ports=cbsd.number_of_ports,
    )


def _build_preferences(cbsd: DBCbsd) -> FrequencyPreferences:
    return FrequencyPreferences(
        bandwidth_mhz=cbsd.preferred_bandwidth_mhz,
        frequencies_mhz=cbsd.preferred_frequencies_mhz,
    )


def _build_db_data(cbsd: DBCbsd) -> DatabaseCbsd:
    return DatabaseCbsd(
        id=cbsd.id,
        should_deregister=cbsd.should_deregister,
        is_deleted=cbsd.is_deleted,
    )


def _build_sas_settings(cbsd: DBCbsd) -> SasSettings:
    return SasSettings(
        single_step_enabled=cbsd.single_step_enabled,
        cbsd_category=cbsd.cbsd_category,
        serial_number=cbsd.cbsd_serial_number,
        fcc_id=cbsd.fcc_id,
        user_id=cbsd.user_id,
    )


def _build_installation_params(cbsd: DBCbsd) -> InstallationParams:
    # TODO do not send installation params to registered cbsd
    return InstallationParams(
        latitude_deg=cbsd.latitude_deg,
        longitude_deg=cbsd.longitude_deg,
        height_m=cbsd.height_m,
        height_type=cbsd.height_type,
        indoor_deployment=cbsd.indoor_deployment,
        antenna_gain_dbi=cbsd.antenna_gain,
    )


def _to_timestamp(t: Optional[datetime]) -> int:
    return 0 if t is None else int(t.timestamp())


def _make_optional_float(value: Optional[float]) -> FloatValue:
    return FloatValue(value=value) if value is not None else None
