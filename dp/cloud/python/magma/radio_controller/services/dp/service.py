import logging

from dp.protos.enodebd_dp_pb2_grpc import DPServiceServicer
from dp.protos.enodebd_dp_pb2 import CBSDStateRequest, CBSDStateResult, LteChannel
from magma.db_service.session_manager import Session, SessionManager
from magma.db_service.models import DBActiveModeConfig, DBCbsd, DBCbsdState, DBChannel, DBGrant, DBGrantState
from magma.mappings.types import CbsdStates, GrantStates

logger = logging.getLogger(__name__)


class DPService(DPServiceServicer):
    def __init__(self, session_manager: SessionManager):
        self.session_manager = session_manager

    def GetCBSDState(self, request: CBSDStateRequest, context) -> CBSDStateResult:
        logger.info(f"Getting CBSD state for {request.serial_number}")
        with self.session_manager.session_scope() as session:
            cbsd = self._get_or_create_cbsd(session, request)
            self._create_or_update_active_mode_config(session, cbsd)
            grant = self._get_channel_with_authorized_grant(session, cbsd)
            result = self._build_result(grant)
            session.commit()
        return result

    @staticmethod
    def _get_or_create_cbsd(session: Session, request: CBSDStateRequest) -> DBCbsd:
        cbsd = session.query(DBCbsd).filter(
            DBCbsd.user_id == request.user_id,
            DBCbsd.fcc_id == request.fcc_id,
            DBCbsd.cbsd_serial_number == request.serial_number,
        ).first()
        if cbsd:
            return cbsd
        unregistered_state = session.query(DBCbsdState). \
            filter(DBCbsdState.name == CbsdStates.UNREGISTERED.value).first()
        cbsd = DBCbsd(
            cbsd_serial_number=request.serial_number,
            fcc_id=request.fcc_id,
            user_id=request.user_id,
            state=unregistered_state,
        )
        session.add(cbsd)
        return cbsd

    @staticmethod
    def _create_or_update_active_mode_config(session: Session, cbsd: DBCbsd) -> None:
        registered_state = session.query(DBCbsdState). \
            filter(DBCbsdState.name == CbsdStates.REGISTERED.value).first()
        active_mode_config = session.query(DBActiveModeConfig). \
            filter(DBActiveModeConfig.cbsd_id == cbsd.id).first()
        if active_mode_config:
            active_mode_config.desired_state = registered_state
            return
        active_mode_config = DBActiveModeConfig(
            desired_state=registered_state,
            cbsd=cbsd,
        )
        session.add(active_mode_config)

    @staticmethod
    def _get_channel_with_authorized_grant(session: Session, cbsd: DBCbsd) -> DBChannel:
        authorized_state = session.query(DBGrantState). \
            filter(DBGrantState.name == GrantStates.AUTHORIZED.value).first()
        return session.query(DBChannel).join(DBGrant).filter(
            DBChannel.cbsd_id == cbsd.id,
            DBGrant.state_id == authorized_state.id,
        ).first()

    @staticmethod
    def _build_result(channel: DBChannel):
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
