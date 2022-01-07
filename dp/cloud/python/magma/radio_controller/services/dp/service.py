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
from typing import Callable

from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from dp.protos.enodebd_dp_pb2_grpc import DPServiceServicer
from magma.db_service.models import (
    DBActiveModeConfig,
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBGrant,
    DBGrantState,
    DBLog,
    DBRequest,
    DBRequestState,
    DBRequestType,
)
from magma.db_service.session_manager import Session, SessionManager
from magma.mappings.types import (
    CbsdStates,
    GrantStates,
    RequestStates,
    RequestTypes,
)
from sqlalchemy.sql.functions import now

logger = logging.getLogger(__name__)


class DPService(DPServiceServicer):
    """
    DP gRPC service class
    """

    def __init__(self, session_manager: SessionManager, now_func: Callable[..., datetime]):
        self.session_manager = session_manager
        self.now = now_func

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
        with self.session_manager.session_scope() as session:
            cbsd = session.query(DBCbsd).filter(
                DBCbsd.cbsd_serial_number == request.serial_number,
            ).first()
            self._log_request(session, 'GetCBSDState', request, cbsd)
            if not cbsd:
                logger.warning(
                    f"State request from unknown CBSD: {request.serial_number}",
                )
                result = self._build_result(None)
            else:
                cbsd.last_seen = self.now()
                channel = self._get_channel_with_authorized_grant(
                    session, cbsd,
                )
                result = self._build_result(channel)
            self._log_result(
                session, 'GetCBSDState', result,
                cbsd, request.serial_number,
            )
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
        with self.session_manager.session_scope() as session:
            cbsd = self._get_or_create_cbsd(session, request)
            self._log_request(session, 'CBSDRegister', request, cbsd)
            self._create_or_update_active_mode_config(session, cbsd)
            channel = self._get_channel_with_authorized_grant(session, cbsd)
            result = self._build_result(channel)
            self._log_result(
                session, 'CBSDRegister', result,
                cbsd, request.serial_number,
            )
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
        with self.session_manager.session_scope() as session:
            cbsd = session.query(DBCbsd).filter(
                DBCbsd.cbsd_serial_number == request.serial_number,
            ).first()
            self._log_request(session, 'CBSDDeregister', request, cbsd)
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
            result = self._build_result(None)
            self._log_result(
                session, 'CBSDDeregister',
                result, cbsd, request.serial_number,
            )
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
        with self.session_manager.session_scope() as session:
            cbsd = session.query(DBCbsd).filter(
                DBCbsd.cbsd_serial_number == request.serial_number,
            ).first()
            self._log_request(session, 'CBSDRelinquish', request, cbsd)
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
            result = self._build_result(None)
            self._log_result(
                session, 'CBSDRelinquish',
                result, cbsd, request.serial_number,
            )
            session.commit()

        return result

    def _add_relinquish_requests(self, session: Session, cbsd: DBCbsd) -> None:
        request_pending_state = session.query(DBRequestState).filter(
            DBRequestState.name == RequestStates.PENDING.value,
        ).scalar()
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
                state=request_pending_state,
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

    def _create_or_update_active_mode_config(self, session: Session, cbsd: DBCbsd) -> DBActiveModeConfig:
        registered_state = session.query(DBCbsdState). \
            filter(DBCbsdState.name == CbsdStates.REGISTERED.value).first()
        active_mode_config = session.query(DBActiveModeConfig). \
            filter(DBActiveModeConfig.cbsd_id == cbsd.id).first()
        if active_mode_config:
            active_mode_config.desired_state = registered_state
            return None
        active_mode_config = DBActiveModeConfig(
            desired_state=registered_state,
            cbsd=cbsd,
        )
        session.add(active_mode_config)
        return active_mode_config

    def _get_channel_with_authorized_grant(self, session: Session, cbsd: DBCbsd) -> DBChannel:
        authorized_state = session.query(DBGrantState). \
            filter(DBGrantState.name == GrantStates.AUTHORIZED.value).first()
        channel = session.query(DBChannel).join(DBGrant).filter(
            DBChannel.cbsd_id == cbsd.id,
            DBGrant.state_id == authorized_state.id,
            (DBGrant.transmit_expire_time == None) | (  # noqa: WPS465,E711
                DBGrant.transmit_expire_time > now()
            ),
            (DBGrant.grant_expire_time == None) | (  # noqa: WPS465,E711
                DBGrant.grant_expire_time > now()
            ),
        ).first()
        return channel

    def _build_result(self, channel: DBChannel):
        if not channel:
            return CBSDStateResult(radio_enabled=False)
        return CBSDStateResult(
            radio_enabled=True,
            channel=LteChannel(
                low_frequency_hz=channel.low_frequency,
                high_frequency_hz=channel.high_frequency,
                max_eirp_dbm_mhz=channel.last_used_max_eirp,
            ),
        )

    def _log_request(self, session: Session, method_name: str, request: CBSDRequest, cbsd: DBCbsd):
        cbsd_serial_number = request.serial_number
        network_id = ''
        fcc_id = ''
        if cbsd:
            network_id = cbsd.network_id or ''
            fcc_id = cbsd.fcc_id or ''
        log = DBLog(
            log_from='CBSD',
            log_to='DP',
            log_name=method_name + '_request',
            log_message=f'{request}',
            cbsd_serial_number=f'{cbsd_serial_number}',
            network_id=f'{network_id}',
            fcc_id=f'{fcc_id}',
        )
        session.add(log)

    def _log_result(self, session: Session, method_name: str, result: CBSDStateResult, cbsd: DBCbsd, serial_number: str):
        network_id = ''
        fcc_id = ''
        if cbsd:
            network_id = cbsd.network_id or ''
            fcc_id = cbsd.fcc_id or ''
        log = DBLog(
            log_from='DP',
            log_to='CBSD',
            log_name=method_name + '_response',
            log_message=f'{result}',
            cbsd_serial_number=f'{serial_number}',
            network_id=f'{network_id}',
            fcc_id=f'{fcc_id}',
        )
        session.add(log)
