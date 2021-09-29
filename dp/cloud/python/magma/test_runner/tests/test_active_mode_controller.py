import json
from time import sleep

import grpc
from retrying import retry
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import DBCbsd
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.db_testcase import DBTestCase
from magma.fixtures.fake_requests.registration_requests import (
    registration_requests,
)
from magma.mappings.cbsd_states import cbsd_state_mapping
from magma.mappings.types import CbsdStates, Switch
from magma.test_runner.config import TestConfig
from dp.protos.active_mode_pb2 import ToggleActiveModeParams
from dp.protos.active_mode_pb2_grpc import ActiveModeControllerStub
from dp.protos.enodebd_dp_pb2 import CBSDStateRequest, CBSDStateResult, LteChannel
from dp.protos.enodebd_dp_pb2_grpc import DPServiceStub
from dp.protos.requests_pb2 import RequestPayload
from dp.protos.requests_pb2_grpc import RadioControllerStub

config = TestConfig()


class ActiveModeControllerTestCase(DBTestCase):

    def setUp(self):
        super().setUp()
        self.grpc_channel = grpc.insecure_channel(f"{config.GRPC_SERVICE}:{config.GRPC_PORT}")
        DBInitializer(SessionManager(self.engine)).initialize()

    # retrying is needed because of a possible deadlock
    # with cc locking request table
    @retry(stop_max_attempt_number=5, wait_fixed=100)
    def drop_all(self):
        super().drop_all()

    def test_cbsd_auto_registered(self):
        # Given
        amc_client = ActiveModeControllerStub(self.grpc_channel)
        rc_client = RadioControllerStub(self.grpc_channel)
        dp_client = DPServiceStub(self.grpc_channel)

        rc_client.UploadRequests(RequestPayload(payload=json.dumps(registration_requests[0])), wait_for_ready=True)

        cbsd = self.session.query(DBCbsd).first()
        self.session.commit()

        # When
        amc_client.ToggleActiveMode(ToggleActiveModeParams(
            cbsd_id=cbsd.id,
            switch=Switch.ON.value,
            desired_state=cbsd_state_mapping[CbsdStates.REGISTERED.value]),
            wait_for_ready=True,
        )

        # Then
        self.then_cbsd_eventually_acquires_grant()

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def then_cbsd_eventually_acquires_grant(self):
        self.session.commit()
        cbsd = self.session.query(DBCbsd).first()

        self.assertEqual(CbsdStates.REGISTERED.value, cbsd.state.name)
        self.assertEqual(1, len(cbsd.channels))
        self.assertEqual(1, len(cbsd.grants))

    def test_provision_cbsd_in_sas_requested_by_dp_client(self):
        dp_client = DPServiceStub(self.grpc_channel)
        dp_client.GetCBSDState(self._build_get_state_request(), wait_for_ready=True)

        self.then_cbsd_is_eventually_provisioned_in_sas(dp_client)

    @staticmethod
    def _build_get_state_request() -> CBSDStateRequest:
        return CBSDStateRequest(
            user_id="some_user_id",
            fcc_id="some_fcc_id",
            serial_number="some_serial_number",
        )

    @staticmethod
    def _build_get_state_result() -> CBSDStateResult:
        return CBSDStateResult(
            radio_enabled=True,
            channel=LteChannel(
                low_freq_hz=3620_000_000,
                high_freq_hz=3630_000_000,
                max_eirp_dbm_mhz=37.0,
            ),
        )

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def then_cbsd_is_eventually_provisioned_in_sas(self, dp_client: DPServiceStub):
        state = dp_client.GetCBSDState(self._build_get_state_request())
        self.assertEqual(self._build_get_state_result(), state)
