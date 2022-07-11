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
import logging
from datetime import datetime
from typing import Callable, List, Optional

from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from dp.protos.enodebd_dp_pb2_grpc import DPServiceServicer
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBGrant,
    DBGrantState,
    DBRequest,
    DBRequestType,
)
from magma.db_service.session_manager import Session, SessionManager
from magma.fluentd_client.client import FluentdClient, FluentdClientException
from magma.fluentd_client.dp_logs import make_dp_log
from magma.mappings.types import CbsdStates, GrantStates, RequestTypes
from magma.radio_controller.config import get_config
from magma.radio_controller.metrics import GET_CBSD_STATE_PROCESSING_TIME
from sqlalchemy.sql.functions import now

logger = logging.getLogger(__name__)


config = get_config()


class DPService(DPServiceServicer):
    """
    DP gRPC service class
    """

    def __init__(
            self,
            session_manager: SessionManager,
            now_func: Callable[..., datetime],
            fluentd_client: FluentdClient,
    ):
        self.session_manager = session_manager
        self.now = now_func
        self.fluentd_client = fluentd_client

    @GET_CBSD_STATE_PROCESSING_TIME.time()
    def GetCBSDState(self, request: CBSDRequest, context) -> CBSDStateResult:
        """
        Get CBSD SAS state for a given CBSD

        Parameters:
            request: CBSDRequest gRPC message
            context: gRPC context

        Returns:
            CBSDStateResult: State result with RF data
        """
        logger.info(f"Getting CBSD state for {request.serial_number=}")
        message_type = 'GetCBSDState'
        with self.session_manager.session_scope() as session:
            cbsd = session.query(DBCbsd).filter(
                DBCbsd.cbsd_serial_number == request.serial_number,
            ).first()
            self._log_request(message_type, request, cbsd)
            if not cbsd:
                logger.warning(
                    f"State request from unknown CBSD: {request.serial_number}",
                )
                result = self._build_result()
            else:
                cbsd.last_seen = self.now()
                result = self._build_result(cbsd, session)
            self._log_result(message_type, result, cbsd, request.serial_number)
            session.commit()
        logger.info(
            f"Returning CBSD state {result=} for {request.serial_number=}",
        )
        return result

    def CBSDRegister(self, request: CBSDRequest, context) -> CBSDStateResult:
        """
        Register CBSD in Domain Proxy using params in the Request message

        Parameters:
            request: CBSDRequest gRPC message
            context: gRPC context

        Returns:
            CBSDStateResult: State result with RF data
        """

        logger.info(f"Registering CBSD {request.serial_number=}")
        message_type = 'CBSDRegister'
        with self.session_manager.session_scope() as session:
            cbsd = self._get_or_create_cbsd(session, request)
            self._log_request(message_type, request, cbsd)
            self._set_desired_state_to_registered(session, cbsd)
            result = self._build_result(cbsd, session)
            self._log_result(message_type, result, cbsd, request.serial_number)
            session.commit()
        return result

    def CBSDDeregister(self, request: CBSDRequest, context) -> CBSDStateResult:
        """
        Deregister CBSD in Domain Proxy

        Parameters:
            request: CBSDRequest gRPC message
            context: gRPC context

        Returns:
            CBSDStateResult: State result with RF data
        """

        logger.info(f"Deregistering CBSD {request.serial_number=}")
        message_type = 'CBSDDeregister'
        with self.session_manager.session_scope() as session:
            cbsd = session.query(DBCbsd).filter(
                DBCbsd.cbsd_serial_number == request.serial_number,
            ).first()
            self._log_request(message_type, request, cbsd)
            if not cbsd:
                logger.info(
                    f"CBSD with serial number {request.serial_number} does not exist.",
                )
            elif not cbsd.active_mode_config:
                logger.info(
                    f"CBSD with serial number {request.serial_number} does not have active mode config.",
                )
            else:
                deregistered_state = session.query(DBCbsdState).filter(
                    DBCbsdState.name == CbsdStates.UNREGISTERED.value,
                ).first()
                cbsd.active_mode_config[0].desired_state = deregistered_state
                logger.info(
                    f"{cbsd.active_mode_config=} set for {cbsd.cbsd_serial_number=}.",
                )
            result = self._build_result()
            self._log_result(message_type, result, cbsd, request.serial_number)
            session.commit()
        return result

    def CBSDRelinquish(self, request: CBSDRequest, context) -> CBSDStateResult:
        """
        Relinquish all CBSD grants in Domain Proxy

        Parameters:
            request: CBSDRequest gRPC message
            context: gRPC context

        Returns:
            CBSDStateResult: State result with RF data
        """

        logger.info(f"Relinquishing grants for CBSD {request.serial_number=}")
        message_type = 'CBSDRelinquish'
        with self.session_manager.session_scope() as session:
            cbsd = session.query(DBCbsd).filter(
                DBCbsd.cbsd_serial_number == request.serial_number,
            ).first()
            self._log_request(message_type, request, cbsd)
            if not cbsd:
                logger.info(
                    f"CBSD with serial number {request.serial_number} does not exist.",
                )
            elif not cbsd.cbsd_id:
                logger.info(
                    f"CBSD with serial number {request.serial_number} does not have SAS CBSD ID.",
                )
            else:
                self._add_relinquish_requests(session, cbsd)
            result = self._build_result()
            self._log_result(message_type, result, cbsd, request.serial_number)
            session.commit()

        return result

    def _add_relinquish_requests(self, session: Session, cbsd: DBCbsd) -> None:
        deregister_request_type = session.query(DBRequestType).filter(
            DBRequestType.name == RequestTypes.RELINQUISHMENT.value,
        ).scalar()
        grants = session.query(DBGrant).join(DBGrantState).filter(
            DBGrant.cbsd_id == cbsd.id, DBGrantState.name != GrantStates.IDLE.value,
        )
        for grant in grants:
            request_dict = {"cbsdId": cbsd.cbsd_id, "grantId": grant.grant_id}
            db_request = DBRequest(
                type=deregister_request_type,
                cbsd=cbsd,
                payload=request_dict,
            )
            session.add(db_request)
            logger.debug(f"Added {db_request=}.")
        pass

    def _get_or_create_cbsd(self, session: Session, request: CBSDRequest) -> DBCbsd:
        cbsd = session.query(DBCbsd).filter(
            DBCbsd.cbsd_serial_number == request.serial_number,
        ).first()
        if cbsd:
            self._update_fields_from_request(cbsd, request)
            return cbsd
        unregistered_state = session.query(DBCbsdState). \
            filter(DBCbsdState.name == CbsdStates.UNREGISTERED.value).first()
        cbsd = DBCbsd(
            cbsd_serial_number=request.serial_number,
            state=unregistered_state,
            desired_state=unregistered_state,
        )
        self._update_fields_from_request(cbsd, request)
        session.add(cbsd)
        return cbsd

    def _update_fields_from_request(self, cbsd: DBCbsd, request: CBSDRequest):
        cbsd.fcc_id = request.fcc_id
        cbsd.user_id = request.user_id
        cbsd.min_power = request.min_power
        cbsd.max_power = request.max_power
        cbsd.antenna_gain = request.antenna_gain
        cbsd.number_of_ports = request.number_of_ports

    def _set_desired_state_to_registered(self, session: Session, cbsd: DBCbsd) -> None:
        registered_state = session.query(DBCbsdState). \
            filter(DBCbsdState.name == CbsdStates.REGISTERED.value).first()
        cbsd.desired_state = registered_state

    def _get_authorized_grants(self, session: Session, cbsd: DBCbsd) -> List[DBGrant]:
        authorized_state = session.query(DBGrantState). \
            filter(DBGrantState.name == GrantStates.AUTHORIZED.value).first()
        grants = session.query(DBGrant).filter(
            DBGrant.cbsd_id == cbsd.id,
            DBGrant.state_id == authorized_state.id,
            (DBGrant.transmit_expire_time == None) | (  # noqa: WPS465,E711
                DBGrant.transmit_expire_time > now()
            ),
            (DBGrant.grant_expire_time == None) | (  # noqa: WPS465,E711
                DBGrant.grant_expire_time > now()
            ),
        ).all()
        return grants

    def _build_result(self, cbsd: Optional[DBCbsd] = None, session: Optional[Session] = None):
        logger.debug("Building GetCbsdResult")
        if not cbsd:
            return CBSDStateResult(radio_enabled=False)
        grants = self._get_authorized_grants(session, cbsd)
        if not grants or cbsd.is_deleted:
            return CBSDStateResult(radio_enabled=False)
        channels = self._build_lte_channels(grants)
        return CBSDStateResult(
            radio_enabled=True,
            carrier_aggregation_enabled=cbsd.carrier_aggregation_enabled,
            channel=channels[0],
            channels=channels,
        )

    def _build_lte_channels(self, grants: List[DBChannel]) -> List[LteChannel]:
        channels = []
        for g in grants:
            channels.append(
                LteChannel(
                    low_frequency_hz=g.low_frequency,
                    high_frequency_hz=g.high_frequency,
                    max_eirp_dbm_mhz=g.max_eirp,
                ),
            )
        return channels

    def _log_request(self, message_type: str, request: CBSDRequest, cbsd: DBCbsd):
        try:
            log = make_dp_log(request, message_type, cbsd)
            self.fluentd_client.send_dp_log(log)
        except (FluentdClientException, TypeError) as err:
            logging.error(f"Failed to log {message_type} request. {err}")

    def _log_result(self, message_type: str, result: CBSDStateResult, cbsd: DBCbsd, serial_number: str):
        try:
            log = make_dp_log(result, message_type, cbsd, serial_number)
            self.fluentd_client.send_dp_log(log)
        except (FluentdClientException, TypeError) as err:
            logging.error(f"Failed to log {message_type} result. {err}")
