from __future__ import annotations

from concurrent import futures

import grpc
from dp.protos.active_mode_pb2 import (
    AcknowledgeCbsdUpdateRequest,
    Authorized,
    DeleteCbsdRequest,
    GetStateRequest,
    Granted,
    GrantRequest,
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
    DBRequestState,
    DBRequestType,
)
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.mappings.types import (
    CbsdStates,
    GrantStates,
    RequestStates,
    RequestTypes,
)
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
        request_states = {
            x.name: x.id for x in self.session.query(DBRequestState).all()
        }
        request_types = {
            x.name: x.id for x in self.session.query(DBRequestType).all()
        }

        self.unregistered = cbsd_states[CbsdStates.UNREGISTERED.value]
        self.registered = cbsd_states[CbsdStates.REGISTERED.value]

        self.idle = grant_states[GrantStates.IDLE.value]
        self.granted = grant_states[GrantStates.GRANTED.value]
        self.authorized = grant_states[GrantStates.AUTHORIZED.value]

        self.pending = request_states[RequestStates.PENDING.value]
        self.processed = request_states[RequestStates.PROCESSED.value]

        self.grant = request_types[RequestTypes.GRANT.value]

    def _prepare_base_cbsd(self) -> DBCbsdBuilder:
        return DBCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(self.unregistered). \
            with_registration('some'). \
            with_eirp_capabilities(0, 10, 20, 1). \
            with_active_mode_config(self.registered)

    @staticmethod
    def _prepare_base_active_mode_config() -> ActiveModeCbsdBuilder:
        return ActiveModeCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(Unregistered). \
            with_registration('some'). \
            with_eirp_capabilities(0, 10, 20, 1). \
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

        self.assertFalse(cbsd.is_updated)

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
            with_channel(1, 2, 3, 4). \
            with_channel(5, 6). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_channel(1, 2, 3, 4). \
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

    def test_get_state_with_requests(self):
        cbsd = self._prepare_base_cbsd(). \
            with_request(self.processed, self.grant, '{"key1":"value1"}'). \
            with_request(self.pending, self.grant, '{"key2":"value2"}'). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        config = self._prepare_base_active_mode_config(). \
            with_pending_request(GrantRequest, '{"key2":"value2"}'). \
            build()
        expected = State(cbsds=[config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_get_state_with_multiple_cbsds(self):
        some_cbsd = DBCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(self.unregistered). \
            with_registration('some'). \
            with_eirp_capabilities(0, 10, 20, 1). \
            with_active_mode_config(self.registered). \
            build()
        other_cbsd = DBCbsdBuilder(). \
            with_id(OTHER_ID). \
            with_state(self.registered). \
            with_registration('other'). \
            with_eirp_capabilities(5, 15, 25, 3). \
            with_active_mode_config(self.registered). \
            build()
        self.session.add_all([some_cbsd, other_cbsd])
        self.session.commit()

        some_config = ActiveModeCbsdBuilder(). \
            with_id(SOME_ID). \
            with_state(Unregistered). \
            with_registration('some'). \
            with_eirp_capabilities(0, 10, 20, 1). \
            with_desired_state(Registered). \
            build()
        other_config = ActiveModeCbsdBuilder(). \
            with_id(OTHER_ID). \
            with_state(Registered). \
            with_registration('other'). \
            with_eirp_capabilities(5, 15, 25, 3). \
            with_desired_state(Registered). \
            build()
        expected = State(cbsds=[some_config, other_config])

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_not_get_state_without_active_mode_config(self):
        cbsd = DBCbsdBuilder(). \
            with_state(self.unregistered). \
            with_registration('some'). \
            with_eirp_capabilities(0, 10, 20, 1). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        expected = State()

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)

    def test_not_get_state_without_registration_params(self):
        cbsd = DBCbsdBuilder(). \
            with_state(self.unregistered). \
            with_eirp_capabilities(0, 10, 20, 1). \
            with_active_mode_config(self.registered). \
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
            with_active_mode_config(self.registered). \
            build()
        self.session.add(cbsd)
        self.session.commit()

        expected = State()

        actual = self.amc_service.GetState(GetStateRequest(), None)
        self.assertEqual(expected, actual)
