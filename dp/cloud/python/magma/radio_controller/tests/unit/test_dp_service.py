from datetime import datetime, timedelta
from typing import Any, List

from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import (
    DBActiveModeConfig,
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBGrant,
    DBGrantState,
)
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.mappings.cbsd_states import CbsdStates, GrantStates
from magma.radio_controller.services.dp.service import DPService

test_cbsd_dict = {
    "serial_number": "some_serial_number",
    "fcc_id": "some_fcc_id",
    "user_id": "some_user_id",
    "min_power": 0,
    "max_power": 10,
    "antenna_gain": 5,
    "number_of_ports": 1,
}

SOME_TIMESTAMP = 1234


class DPTestCase(LocalDBTestCase):
    def setUp(self):
        super().setUp()
        self.dp_service = DPService(SessionManager(self.engine), self._now)
        DBInitializer(SessionManager(self.engine)).initialize()

        self.cbsd_states = {
            x.name: x.id for x in self.session.query(DBCbsdState).all()
        }
        self.grant_states = {
            x.name: x.id for x in self.session.query(DBGrantState).all()
        }

    def test_cbsd_registered_and_enabled_active_mode(self):
        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)

        self._then_exactly_one_cbsd_is_active(request)

    def test_cbsd_register_when_already_registered(self):
        cbsd = self._build_cbsd(**test_cbsd_dict)
        self.session.add(cbsd)
        self.session.commit()

        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)

        self._then_exactly_one_cbsd_is_active(request)

    def test_cbsd_register_updates_active_mode_config_if_desired_state_is_unregistered(self):
        cbsd = self._build_cbsd(**test_cbsd_dict)
        active_mode_config = DBActiveModeConfig(
            cbsd=cbsd,
            desired_state_id=self.cbsd_states[CbsdStates.UNREGISTERED.value],
        )
        self.session.add_all([cbsd, active_mode_config])
        self.session.commit()

        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)

        self.assertEqual(1, self._get_active_cbsd_count(request))

    def test_get_state_for_unknown_cbsd(self):
        request = self._build_request(**test_cbsd_dict)
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_empty_result(), result)

    def test_get_state_with_valid_authorized_grant(self):
        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)
        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_serial_number == request.serial_number,
        ).first()

        channel = self._build_channel(cbsd)
        grant = self._build_grant(
            cbsd, self.grant_states[GrantStates.AUTHORIZED.value], channel,
        )
        grant.transmit_expire_time = datetime.now() + timedelta(seconds=60)
        grant.grant_expire_time = datetime.now() + timedelta(days=7)
        self.session.add_all([cbsd, channel, grant])
        self.session.commit()

        request = self._build_request(**test_cbsd_dict)
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_expected_result(channel), result)

    def test_get_state_with_transmit_expired_authorized_grant(self):
        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)
        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_serial_number == request.serial_number,
        ).first()
        channel = self._build_channel(cbsd)
        grant = self._build_grant(
            cbsd, self.grant_states[GrantStates.AUTHORIZED.value], channel,
        )
        grant.transmit_expire_time = datetime.now() - timedelta(seconds=1)
        grant.grant_expire_time = datetime.now() + timedelta(days=7)
        self.session.add_all([cbsd, channel, grant])
        self.session.commit()

        request = self._build_request()
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_empty_result(), result)

    def test_get_state_with_grant_expired_authorized_grant(self):
        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)
        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_serial_number == request.serial_number,
        ).first()
        channel = self._build_channel(cbsd)
        grant = self._build_grant(
            cbsd, self.grant_states[GrantStates.AUTHORIZED.value], channel,
        )
        grant.transmit_expire_time = datetime.now() - timedelta(seconds=1)
        grant.grant_expire_time = datetime.now() - timedelta(seconds=1)
        self.session.add_all([cbsd, channel, grant])
        self.session.commit()

        request = self._build_request()
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_empty_result(), result)

    def test_get_state_with_unauthorized_grant(self):
        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)
        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_serial_number == request.serial_number,
        ).first()
        channel = self._build_channel(cbsd)
        grant = self._build_grant(
            cbsd, self.grant_states[GrantStates.IDLE.value], channel,
        )
        self.session.add_all([cbsd, channel, grant])
        self.session.commit()

        request = self._build_request()
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_empty_result(), result)

    def test_get_state_without_channel(self):
        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)
        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_serial_number == request.serial_number,
        ).first()
        grant = self._build_grant(
            cbsd, self.grant_states[GrantStates.AUTHORIZED.value],
        )
        self.session.add_all([cbsd, grant])
        self.session.commit()

        request = self._build_request()
        result = self.dp_service.GetCBSDState(request, None)

        self.assertEqual(self._build_empty_result(), result)

    def test_update_last_seen_on_request(self):
        request = self._build_request(**test_cbsd_dict)
        self.dp_service.CBSDRegister(request, None)
        self.dp_service.GetCBSDState(request, None)

        self.assertEqual(
            SOME_TIMESTAMP, self._get_last_seen_timestamp(request),
        )

    @staticmethod
    def _now() -> datetime:
        return datetime.fromtimestamp(SOME_TIMESTAMP)

    @staticmethod
    def _build_request(**kwargs) -> CBSDRequest:
        return CBSDRequest(**kwargs)

    def _build_cbsd(self, serial_number, **kwargs) -> DBCbsd:
        return DBCbsd(
            cbsd_serial_number=serial_number,
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

    def _then_exactly_one_cbsd_is_active(self, request: CBSDRequest):
        self.assertEqual(1, self._get_active_cbsd_count(request))
        self.assertEqual(
            1, self._get_cbsd_with_serial_count(request.serial_number),
        )

    def _get_cbsd_with_serial_count(self, serial_number: str) -> int:
        return self.session.query(DBCbsd).\
            filter(DBCbsd.cbsd_serial_number == serial_number).count()

    def _get_active_cbsd_count(self, request: CBSDRequest) -> int:
        return self._query_active_cbsd(request).count()

    def _get_last_seen_timestamp(self, request: CBSDRequest) -> int:
        return self._query_active_cbsd(request).first().cbsd.last_seen.timestamp()

    def _query_active_cbsd(self, request: CBSDRequest):
        return self.session.query(DBActiveModeConfig).join(DBCbsd). \
            filter(*self._build_expected_filter_from_request(request))

    def _build_expected_filter_from_request(self, request: CBSDRequest) -> List[Any]:
        return [
            DBCbsd.cbsd_serial_number == request.serial_number,
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
