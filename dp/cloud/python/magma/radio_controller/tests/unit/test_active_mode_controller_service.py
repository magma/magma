from __future__ import annotations

from concurrent import futures

import grpc
from dp.protos.active_mode_pb2 import (
    AcknowledgeCbsdUpdateRequest,
    Authorized,
    DeleteCbsdRequest,
    GetStateRequest,
    Granted,
    Registered,
    State,
    Unregistered,
)
from dp.protos.active_mode_pb2_grpc import (
    ActiveModeControllerStub,
    add_ActiveModeControllerServicer_to_server,
)
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBGrantState,
    DBRequestType,
)
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.mappings.types import CbsdStates, GrantStates, RequestTypes
from magma.radio_controller.services.active_mode_controller.service import (
    ActiveModeControllerService,
)
from magma.radio_controller.tests.test_utils.active_mode_cbsd_builder import (
    ActiveModeCbsdBuilder,
)
from magma.radio_controller.tests.test_utils.db_cbsd_builder import (
    DBCbsdBuilder,
)

SOME_ID = 123
OTHER_ID = 456


class ActiveModeControllerTestCase(LocalDBTestCase):
    def setUp(self):
        super().setUp()
        self.amc_service = ActiveModeControllerService(
            SessionManager(self.engine),
        )
        DBInitializer(SessionManager(self.engine)).initialize()

        grant_states = {
            x.name: x.id for x in self.session.query(DBGrantState).all()
        }
        cbsd_states = {
            x.name: x.id for x in self.session.query(DBCbsdState).all()
        }
        request_types = {
            x.name: x.id for x in self.session.query(DBRequestType).all()
        }

        self.unregistered = cbsd_states[CbsdStates.UNREGISTERED.value]
        self.registered = cbsd_states[CbsdStates.REGISTERED.value]

        self.idle = grant_states[GrantStates.IDLE.value]
        self.granted = grant_states[GrantStates.GRANTED.value]
        self.authorized = grant_states[GrantStates.AUTHORIZED.value]

        self.grant = request_types[RequestTypes.GRANT.value]

    def _prepare_base_cbsd(self) -> DBCbsdBuilder:
        return DBCbsdBuilder(). \
            with_id(SOME_ID). \
            with_category('a'). \
            with_state(self.unregistered). \
            with_registration('some'). \
            with_eirp_capabilities(0, 10, 1). \
            with_antenna_gain(20). \
            with_desired_state(self.registered)

    @staticmethod
    def _prepare_base_active_mode_config() -> ActiveModeCbsdBuilder:
        return ActiveModeCbsdBuilder(). \
            with_id(SOME_ID). \
            with_category('a'). \
            with_state(Unregistered). \
            with_registration('some'). \
            with_eirp_capabilities(0, 10, 1). \
            with_antenna_gain(20). \
            with_desired_state(Registered)


class ActiveModeControllerClientServerTestCase(ActiveModeControllerTestCase):
    def setUp(self):
        super().setUp()
        self.server = grpc.server(futures.ThreadPoolExecutor(max_workers=2))
        self.port = '50051'
        add_ActiveModeControllerServicer_to_server(
            self.amc_service, self.server,
        )
        self.server.add_insecure_port(f'[::]:{self.port}')
        self.server.start()
        channel = grpc.insecure_channel(f'localhost:{self.port}')
        self.stub = ActiveModeControllerStub(channel)

    def tearDown(self):
        self.server.stop(None)
        self.stub = None
        super().tearDown()

    def test_delete_cbsd(self):
        # Given
        cbsd1 = self._prepare_base_cbsd().\
            with_id(SOME_ID).\
            with_registration("some").\
            build()
        cbsd2 = self._prepare_base_cbsd(). \
            with_id(OTHER_ID). \
            with_registration("some_other").\
            build()
        self.session.add_all([cbsd1, cbsd2])
        self.session.commit()
        cbsds = self.session.query(DBCbsd)
        self.assertEqual(2, cbsds.count())

        # When
        self.stub.DeleteCbsd(DeleteCbsdRequest(id=cbsd1.id))

        # Then
        cbsds = self.session.query(DBCbsd)
        self.assertEqual(1, cbsds.count())
        self.assertEqual(cbsd2.id, cbsds.first().id)

    def test_delete_non_existent_cbsd(self):
        # Given cbsd is not in the database
        # Then
        with self.assertRaises(grpc.RpcError) as err:
            self.stub.DeleteCbsd(DeleteCbsdRequest(id=SOME_ID))
        self.assertEqual(grpc.StatusCode.NOT_FOUND, err.exception.code())

    def test_acknowledge_cbsd_update(self):
        cbsd = self._prepare_base_cbsd(). \
            updated(). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        self.stub.AcknowledgeCbsdUpdate(
            AcknowledgeCbsdUpdateRequest(id=SOME_ID),
        )

        self.assertFalse(cbsd.should_deregister)

    def test_acknowledge_non_existent_cbsd_update(self):
        with self.assertRaises(grpc.RpcError) as err:
            self.stub.AcknowledgeCbsdUpdate(
                AcknowledgeCbsdUpdateRequest(id=SOME_ID),
            )
        self.assertEqual(grpc.StatusCode.NOT_FOUND, err.exception.code())


class ActiveModeControllerServerTestCase(ActiveModeControllerTestCase):

    def test_get_basic_state(self):
        cbsd = self._prepare_base_cbsd().build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config().build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_with_frequency_preferences(self):
        cbsd = self._prepare_base_cbsd(). \
            with_preferences(15, [3600, 3580, 3620]). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_preferences(15, [3600, 3580, 3620]). \
            build()

        expected = State(cbsds=[config])
        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(actual, expected)

    def test_get_state_with_grants(self):
        cbsd = self._prepare_base_cbsd(). \
            with_grant("idle_grant", self.idle, 1, 2). \
            with_grant("granted_grant", self.granted, 3). \
            with_grant("authorized_grant", self.authorized, 5, 6). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_grant("granted_grant", Granted, 3, 0). \
            with_grant("authorized_grant", Authorized, 5, 6). \
            build()

        expected = State(cbsds=[config])
        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(actual, expected)

    def test_get_state_with_channels(self):
        cbsd = self._prepare_base_cbsd(). \
            with_channel(1, 2, 3). \
            with_channel(5, 6). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_channel(1, 2, 3). \
            with_channel(5, 6). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(actual, expected)

    def test_get_state_for_cbsd_marked_for_deletion(self):
        cbsd = self._prepare_base_cbsd(). \
            deleted(). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            deleted(). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_for_cbsd_marked_for_update(self):
        cbsd = self._prepare_base_cbsd(). \
            updated(). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            updated(). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_with_last_seen(self):
        cbsd = self._prepare_base_cbsd(). \
            with_last_seen(1). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_last_seen(1). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_with_grant_attempts(self):
        cbsd = self._prepare_base_cbsd(). \
            with_grant_attempts(1). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_grant_attempts(1). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_with_pending_requests(self):
        cbsd = self._prepare_base_cbsd(). \
            with_request(self.grant, '{"key2":"value2"}'). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        expected = State()

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_with_multiple_cbsds(self):
        some_cbsd = self._prepare_base_cbsd(). \
            with_id(SOME_ID). \
            with_registration('some'). \
            build()
        other_cbsd = self._prepare_base_cbsd(). \
            with_id(OTHER_ID). \
            with_registration('other'). \
            build()
        self.session.add_all([some_cbsd, other_cbsd])
        self.session.commit()

        some_config = self._prepare_base_active_mode_config(). \
            with_id(SOME_ID). \
            with_registration('some'). \
            build()
        other_config = self._prepare_base_active_mode_config(). \
            with_id(OTHER_ID). \
            with_registration('other'). \
            build()
        expected = State(cbsds=[some_config, other_config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_for_single_step(self):
        cbsd = self._prepare_base_cbsd(). \
            with_single_step_enabled(). \
            with_installation_params(1, 2, 3, 'AGL', True). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_single_step_enabled(). \
            with_installation_params(1, 2, 3, 'AGL', True). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_should_send_data_if_cbsd_needs_deregistration(self):
        cbsd = DBCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(self.registered). \
            with_desired_state(self.registered). \
            with_registration('some'). \
            updated(). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = ActiveModeCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(Registered). \
            with_desired_state(Registered). \
            with_registration('some'). \
            updated(). \
            with_category('b'). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_should_send_data_if_cbsd_needs_deletion(self):
        cbsd = DBCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(self.registered). \
            with_desired_state(self.registered). \
            with_registration('some'). \
            deleted(). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = ActiveModeCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(Registered). \
            with_desired_state(Registered). \
            with_registration('some'). \
            deleted(). \
            with_category('b'). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_not_get_state_without_registration_params(self):
        cbsd = DBCbsdBuilder(). \
            with_state(self.unregistered). \
            with_eirp_capabilities(0, 10, 1). \
            with_desired_state(self.registered). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        expected = State()

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_not_get_state_without_eirp_capabilities(self):
        cbsd = DBCbsdBuilder(). \
            with_state(self.unregistered). \
            with_registration('some'). \
            with_desired_state(self.registered). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        expected = State()

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_not_get_state_with_single_step_enabled_and_without_installation_params(self):
        cbsd = self._prepare_base_cbsd(). \
            with_single_step_enabled(). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        expected = State()

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_not_get_state_with_single_step_enabled_category_a_and_outdoor(self):
        cbsd = self._prepare_base_cbsd(). \
            with_category('a'). \
            with_installation_params(1, 2, 3, 'AGL', False). \
            with_single_step_enabled(). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        expected = State()

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)
