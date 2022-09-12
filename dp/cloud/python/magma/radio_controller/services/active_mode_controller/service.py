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
from typing import Any, Dict, List, Optional

import grpc
from dp.protos.active_mode_pb2 import (
    AcknowledgeCbsdRelinquishRequest,
    AcknowledgeCbsdUpdateRequest,
    Cbsd,
    Channel,
    DatabaseCbsd,
    DeleteCbsdRequest,
    EirpCapabilities,
    FrequencyPreferences,
    GetStateRequest,
    Grant,
    GrantSettings,
    InstallationParams,
    RequestPayload,
    SasSettings,
    State,
    StoreAvailableFrequenciesRequest,
)
from dp.protos.active_mode_pb2_grpc import ActiveModeControllerServicer
from google.protobuf.empty_pb2 import Empty
from google.protobuf.wrappers_pb2 import FloatValue
from magma.db_service.models import DBCbsd, DBGrant, DBRequest
from magma.db_service.session_manager import Session, SessionManager
from magma.mappings.cbsd_states import cbsd_state_mapping, grant_state_mapping
from magma.radio_controller.metrics import (
    ACKNOWLEDGE_RELINQUISH_PROCESSING_TIME,
    ACKNOWLEDGE_UPDATE_PROCESSING_TIME,
    DELETE_CBSD_PROCESSING_TIME,
    GET_DB_STATE_PROCESSING_TIME,
    INSERT_TO_DB_PROCESSING_TIME,
    STORE_AVAILABLE_FREQUENCIES_PROCESSING_TIME,
)
from magma.radio_controller.services.active_mode_controller.strategies.strategies_mapping import (
    get_cbsd_filter_strategies,
)
from sqlalchemy import and_, or_
from sqlalchemy.orm import contains_eager, joinedload

logger = logging.getLogger(__name__)


class ActiveModeControllerService(ActiveModeControllerServicer):
    """
    Active Mode Controller gRPC Service class
    """

    def __init__(self, session_manager: SessionManager, request_types_map: Dict[str, int]):
        self.session_manager = session_manager
        self.request_types_map = request_types_map

    @INSERT_TO_DB_PROCESSING_TIME.time()
    def UploadRequests(self, request_payload: RequestPayload, context) -> Empty:
        """
        Insert uploaded requests to the database

        Parameters:
            request_payload: gRPC RequestPayload message
            context: gRPC context

        Returns:
            RequestDbIds: a list of IDs of inserted database records
        """
        logger.info("Storing requests in DB.")
        requests_map = json.loads(request_payload.payload)
        with self.session_manager.session_scope() as session:
            _store_requests_from_map_in_db(session, self.request_types_map, requests_map)
        return Empty()

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
        updated = self._update_cbsd(db_id, {'should_deregister': False})
        if not updated:
            context.set_code(grpc.StatusCode.NOT_FOUND)
        return Empty()

    @ACKNOWLEDGE_RELINQUISH_PROCESSING_TIME.time()
    def AcknowledgeCbsdRelinquish(self, request: AcknowledgeCbsdRelinquishRequest, context) -> Empty:
        """
        Mark CBSD in the Database as not relinquished

        Parameters:
            request: a AcknowledgeCbsdRelinquishRequest gRPC Message
            context: gRPC context

        Returns:
            Empty: an empty gRPC message
        """
        db_id = request.id
        logger.info(f"Acknowledging CBSD relinquish {db_id}")
        updated = self._update_cbsd(db_id, {'should_relinquish': False})
        if not updated:
            context.set_code(grpc.StatusCode.NOT_FOUND)
        return Empty()

    @STORE_AVAILABLE_FREQUENCIES_PROCESSING_TIME.time()
    def StoreAvailableFrequencies(self, request: StoreAvailableFrequenciesRequest, context) -> Empty:
        """
        Store available frequencies in the database

        Parameters:
            request: StoreAvailableFrequencies gRPC Message
            context: gRPC context

        Returns:
            Empty: an empty gRPC message
        """
        db_id = request.id
        logger.info(f"Storing available frequencies for {db_id}")
        updated = self._update_cbsd(db_id, {"available_frequencies": list(request.available_frequencies)})
        if not updated:
            context.set_code(grpc.StatusCode.NOT_FOUND)
        return Empty()

    def _update_cbsd(self, db_id: int, to_update: Dict) -> DBCbsd:
        with self.session_manager.session_scope() as session:
            updated = session.query(DBCbsd).filter(DBCbsd.id == db_id).update(to_update)
            session.commit()
        return updated


def _store_requests_from_map_in_db(session: Session, request_types_map: Dict[str, int], request_map: Dict[str, List[Dict]]) -> None:
    request_type = next(iter(request_map))
    for request_json in request_map[request_type]:
        filters = get_cbsd_filter_strategies[request_type](request_json)
        cbsd_id = session.query(DBCbsd.id).filter(*filters).first()
        if not cbsd_id:
            logger.error(
                f"Could not obtain cbsd to bind to the request: {request_json}",
            )
            continue
        db_request = DBRequest(
            type_id=request_types_map[request_type],
            cbsd_id=cbsd_id[0],
            payload=request_json,
        )
        logger.info(f"Adding request {db_request}.")
        session.add(db_request)
    session.commit()


def _list_cbsds(session: Session) -> State:
    # Selectively load sqlalchemy object relations using a single query to avoid commit races.
    return (
        session.query(DBCbsd).
        outerjoin(DBGrant).
        outerjoin(DBRequest).
        options(
            joinedload(DBCbsd.state),
            joinedload(DBCbsd.desired_state),
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
                DBCbsd.should_relinquish == True,
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
    """
    Check that the fields are not null

    Parameters:
        fields (List[Any]): db fields

    Returns:
        None
    """
    return and_(*[field != None for field in fields])  # noqa: E711


def _build_state(db_cbsds: List[DBCbsd]) -> State:
    cbsds = [_build_cbsd(db_cbsd) for db_cbsd in db_cbsds]
    return State(cbsds=cbsds)


def _build_cbsd(cbsd: DBCbsd) -> Cbsd:
    # Application may not need those to be sorted.
    # Applying ordering mostly for easier assertions in testing
    cbsd_db_grants = sorted(cbsd.grants, key=lambda x: x.id)
    grants = [_build_grant(x) for x in cbsd_db_grants]
    channels = [_build_channel(x) for x in cbsd.channels]

    last_seen = _to_timestamp(cbsd.last_seen)
    eirp_capabilities = _build_eirp_capabilities(cbsd)
    preferences = _build_preferences(cbsd)
    sas_settings = _build_sas_settings(cbsd)
    installation_params = _build_installation_params(cbsd)
    db_data = _build_db_data(cbsd)
    grant_settings = _build_grant_settings(cbsd)
    return Cbsd(
        cbsd_id=cbsd.cbsd_id,
        state=cbsd_state_mapping[cbsd.state.name],
        desired_state=cbsd_state_mapping[cbsd.desired_state.name],
        grants=grants,
        channels=channels,
        last_seen_timestamp=last_seen,
        eirp_capabilities=eirp_capabilities,
        db_data=db_data,
        preferences=preferences,
        sas_settings=sas_settings,
        installation_params=installation_params,
        grant_settings=grant_settings,
    )


def _build_grant(grant: DBGrant) -> Grant:
    last_heartbeat = _to_timestamp(grant.last_heartbeat_request_time)
    return Grant(
        id=grant.grant_id,
        state=grant_state_mapping[grant.state.name],
        heartbeat_interval_sec=grant.heartbeat_interval,
        last_heartbeat_timestamp=last_heartbeat,
        low_frequency_hz=grant.low_frequency,
        high_frequency_hz=grant.high_frequency,
    )


def _build_channel(channel: dict) -> Channel:
    return Channel(
        low_frequency_hz=channel.get('low_frequency'),
        high_frequency_hz=channel.get('high_frequency'),
        max_eirp=_make_optional_float(channel.get('max_eirp')),
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
        should_relinquish=cbsd.should_relinquish,
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


def _build_grant_settings(cbsd: DBCbsd) -> GrantSettings:
    return GrantSettings(
        grant_redundancy_enabled=cbsd.grant_redundancy,
        carrier_aggregation_enabled=cbsd.carrier_aggregation_enabled,
        max_ibw_mhz=cbsd.max_ibw_mhz,
        available_frequencies=cbsd.available_frequencies,
    )


def _to_timestamp(t: Optional[datetime]) -> int:
    return 0 if t is None else int(t.timestamp())


def _make_optional_float(value: Optional[float]) -> FloatValue:
    return FloatValue(value=value) if value is not None else None
