from typing import Any, List

from dp.protos.enodebd_dp_pb2 import CBSDStateRequest, CBSDStateResult, LteChannel
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import DBActiveModeConfig, DBCbsd, DBCbsdState, DBChannel, DBGrant, DBGrantState
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.mappings.cbsd_states import CbsdStates, GrantStates
from magma.radio_controller.services.dp.service import DPService

SOME_USER_ID = "some_user_id"
SOME_FCC_ID = "some_fcc_id"
SOME_SERIAL_NUMBER = "some_serial_number"


class DPTestCase(LocalDBTestCase):
    def setUp(self):
        super().setUp()
        self.dp_service = DPService(SessionManager(self.engine))
        DBInitializer(SessionManager(self.engine)).initialize()

        self.cbsd_states = {x.name: x.id for x in self.session.query(DBCbsdState).all()}
        self.grant_states = {x.name: x.id for x in self.session.query(DBGrantState).all()}

    def test_insert_cbsd_and_enabled_active_mode(self):
        request = self._build_request()
        self.dp_service.GetCBSDState(request, None)

        self.assertEqual(1, self._get_cbsd_count(request))

    def test_not_insert_cbsd_when_already_inserted(self):
        cbsd = self._build_cbsd()
        self.session.add(cbsd)
        self.session.commit()

        request = self._build_request()
        self.dp_service.GetCBSDState(request, None)

        self.assertEqual(1, self._get_cbsd_count(request))

    def test_update_active_mode_config_if_desired_state_is_unregistered(self):
        cbsd = self._build_cbsd()
        active_mode_config = DBActiveModeConfig(
            cbsd=cbsd,
            desired_state_id=self.cbsd_states[CbsdStates.UNREGISTERED.value],
        )
        self.session.add_all([cbsd, active_mode_config])
        self.session.commit()

        request = self._build_request()
        self.dp_service.GetCBSDState(request, None)

        self.assertEqual(1, self._get_cbsd_count(request))

    def test_fetch_state_with_authorized_grant(self):
        cbsd = self._build_cbsd()
        channel = self._build_channel(cbsd)
        grant = self._build_grant(cbsd, self.grant_states[GrantStates.AUTHORIZED.value], channel)
        self.session.add_all([cbsd, channel, grant])
        self.session.commit()

        request = self._build_request()
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_expected_result(channel), result)

    def test_fetch_state_with_unauthorized_grant(self):
        cbsd = self._build_cbsd()
        channel = self._build_channel(cbsd)
        grant = self._build_grant(cbsd, self.grant_states[GrantStates.IDLE.value], channel)
        self.session.add_all([cbsd, channel, grant])
        self.session.commit()

        request = self._build_request()
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_empty_result(), result)

    def test_fetch_state_without_channel(self):
        cbsd = self._build_cbsd()
        grant = self._build_grant(cbsd, self.grant_states[GrantStates.AUTHORIZED.value])
        self.session.add_all([cbsd, grant])
        self.session.commit()

        request = self._build_request()
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_empty_result(), result)

    @staticmethod
    def _build_request() -> CBSDStateRequest:
        return CBSDStateRequest(
            user_id=SOME_USER_ID,
            fcc_id=SOME_FCC_ID,
            serial_number=SOME_SERIAL_NUMBER,
        )

    def _build_cbsd(self) -> DBCbsd:
        return DBCbsd(
            user_id=SOME_USER_ID,
            fcc_id=SOME_FCC_ID,
            cbsd_serial_number=SOME_SERIAL_NUMBER,
            state_id=self.cbsd_states[CbsdStates.UNREGISTERED.value],
        )

    @staticmethod
    def _build_channel(cbsd: DBCbsd) -> DBChannel:
        return DBChannel(
            cbsd=cbsd,
            channel_type="some_channel_type",
            rule_applied="some_rule",
            low_frequency=3550_000_000,
            high_frequency=3570_000_000,
            last_used_max_eirp=20,
        )

    @staticmethod
    def _build_grant(cbsd: DBCbsd, state_id: str, channel: DBChannel = None) -> DBGrant:
        return DBGrant(
            cbsd=cbsd,
            channel=channel,
            grant_id="some_grant_id",
            state_id=state_id,
        )

    def _get_cbsd_count(self, request: CBSDStateRequest) -> int:
        return self.session.query(DBActiveModeConfig).join(DBCbsd). \
            filter(*self._build_expected_filter_from_request(request)).count()

    def _build_expected_filter_from_request(self, request: CBSDStateRequest) -> List[Any]:
        return [
            DBCbsd.cbsd_serial_number == request.serial_number,
            DBCbsd.user_id == request.user_id,
            DBCbsd.fcc_id == request.fcc_id,
            DBCbsd.state_id == self.cbsd_states[CbsdStates.UNREGISTERED.value],
            DBActiveModeConfig.desired_state_id == self.cbsd_states[CbsdStates.REGISTERED.value],
        ]

    @staticmethod
    def _build_expected_result(channel: DBChannel) -> CBSDStateResult:
        return CBSDStateResult(
            radio_enabled=True,
            channel=LteChannel(
                low_frequency_hz=channel.low_frequency,
                high_frequency_hz=channel.high_frequency,
                max_eirp_dbm_mhz=channel.last_used_max_eirp,
            ),
        )

    @staticmethod
    def _build_empty_result() -> CBSDStateResult:
        return CBSDStateResult(radio_enabled=False)
