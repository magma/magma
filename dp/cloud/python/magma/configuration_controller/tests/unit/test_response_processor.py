from typing import Dict

import requests
import responses
from magma.configuration_controller.response_processor.response_db_processor import (
    ResponseDBProcessor,
)
from magma.configuration_controller.response_processor.strategies.strategies_mapping import (
    processor_strategies,
)
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBGrant,
    DBRequest,
    DBRequestState,
    DBRequestType,
    DBResponse,
)
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.fixtures.fake_requests.deregistration_requests import (
    deregistration_requests,
)
from magma.fixtures.fake_requests.grant_requests import grant_requests
from magma.fixtures.fake_requests.heartbeat_requests import (
    heartbeat_requests,
)
from magma.fixtures.fake_requests.registration_requests import (
    registration_requests,
)
from magma.fixtures.fake_requests.relinquishment_requests import (
    relinquishment_requests,
)
from magma.fixtures.fake_requests.spectrum_inquiry_requests import (
    spectrum_inquiry_requests,
)
from magma.fixtures.fake_responses.spectrum_inquiry_responses import (
    single_channel_for_one_cbsd,
    two_channels_for_one_cbsd,
    zero_channels_for_one_cbsd,
)
from magma.mappings.request_response_mapping import request_response
from magma.mappings.types import (
    CbsdStates,
    GrantStates,
    RequestStates,
    RequestTypes,
)
from parameterized import parameterized

CBSD_SERIAL_NR = "cbsdSerialNumber"
FCC_ID = "fccId"
USER_ID = "userId"
CBSD_ID = "cbsdId"
GRANT_ID = "grantId"
TEST_STATE = "test_state"


class DefaultResponseDBProcessorTestCase(LocalDBTestCase):
    def setUp(self):
        super().setUp()
        DBInitializer(SessionManager(self.engine)).initialize()

    @parameterized.expand([
        (processor_strategies["registrationRequest"], registration_requests),
        (processor_strategies["spectrumInquiryRequest"], spectrum_inquiry_requests),
        (processor_strategies["grantRequest"], grant_requests),
        (processor_strategies["heartbeatRequest"], heartbeat_requests),
        (processor_strategies["relinquishmentRequest"], relinquishment_requests),
        (processor_strategies["deregistrationRequest"], deregistration_requests),
    ])
    @responses.activate
    def test_processor_splits_sas_response_into_separate_db_objects_and_links_them_with_requests(
            self, processor_strategy, requests_fixtures):

        # When
        cbsd_state = DBCbsdState(name=TEST_STATE)
        self.session.add(cbsd_state)
        self.session.commit()

        request_type_name = self._get_request_type_from_fixture(requests_fixtures)
        response_type_name = request_response[request_type_name]

        db_requests, response_payload = self._get_db_requests_and_response_payload(
            request_type_name, response_type_name, requests_fixtures, cbsd_state)

        response, processor = self._prepare_response_and_processor(
            response_payload, response_type_name, processor_strategy)

        processor.process_response(db_requests, response, self.session)
        self.session.commit()

        nr_of_requests = len(db_requests)

        # Then
        self.assertEqual(2, self.session.query(DBRequestState).count())
        self.assertEqual(1, self.session.query(DBRequestType).filter(DBRequestType.name == request_type_name).count())
        self.assertEqual(nr_of_requests, self.session.query(DBRequest).count())
        self.assertListEqual([r.id for r in db_requests], [_id for (_id,) in self.session.query(DBResponse.id).all()])
        self.assertListEqual(["processed"] * nr_of_requests,
                             [r.state.name for r in self.session.query(DBRequest).all()])

    @parameterized.expand([
        (processor_strategies["grantRequest"], grant_requests, 0, GrantStates.GRANTED.value),
        (processor_strategies["grantRequest"], grant_requests, 400, GrantStates.IDLE.value),
        (processor_strategies["grantRequest"], grant_requests, 401, GrantStates.IDLE.value),
        (processor_strategies["grantRequest"], grant_requests, 500, GrantStates.IDLE.value),
        (processor_strategies["heartbeatRequest"], heartbeat_requests, 0, GrantStates.AUTHORIZED.value),
        (processor_strategies["heartbeatRequest"], heartbeat_requests, 500, GrantStates.IDLE.value),
        (processor_strategies["heartbeatRequest"], heartbeat_requests, 501, GrantStates.GRANTED.value),
        (processor_strategies["heartbeatRequest"], heartbeat_requests, 502, GrantStates.IDLE.value),
        (processor_strategies["relinquishmentRequest"], relinquishment_requests, 0, GrantStates.IDLE.value),
    ])
    @responses.activate
    def test_grant_state_after_response(
            self, processor_strategy, requests_fixtures, response_code, expected_grant_state_name):

        # When
        cbsd_state = DBCbsdState(name=TEST_STATE)
        self.session.add(cbsd_state)
        self.session.commit()

        request_type_name = self._get_request_type_from_fixture(grant_requests)
        response_type_name = request_response[request_type_name]
        db_requests, response_payload = self._get_db_requests_and_response_payload(
            request_type_name, response_type_name, requests_fixtures, cbsd_state)

        for response_json in response_payload[response_type_name]:
            response_json["response"]["responseCode"] = response_code

        response, processor = self._prepare_response_and_processor(
            response_payload, response_type_name, processor_strategy)

        processor.process_response(db_requests, response, self.session)
        self.session.commit()

        nr_of_requests = len(db_requests)

        # Then
        self.assertEqual(nr_of_requests, self.session.query(DBRequest).count())
        self.assertListEqual([r.id for r in db_requests], [_id for (_id,) in self.session.query(DBResponse.id).all()])
        self.assertListEqual(["processed"] * nr_of_requests,
                             [r.state.name for r in self.session.query(DBRequest).all()])
        self.assertListEqual([expected_grant_state_name] * nr_of_requests,
                             [g.state.name for g in self.session.query(DBGrant).all()])

    @parameterized.expand([
        (0, CbsdStates.REGISTERED),
        (300, CbsdStates.UNREGISTERED),
        (400, CbsdStates.UNREGISTERED),
        (105, CbsdStates.UNREGISTERED),
        (104, CbsdStates.UNREGISTERED),
        (401, CbsdStates.UNREGISTERED),
        (500, CbsdStates.UNREGISTERED),
        (501, CbsdStates.UNREGISTERED),
    ])
    @responses.activate
    def test_cbsd_state_after_registration_response(self, sas_response_code, expected_cbsd_state):
        db_requests = self._create_db_requests_from_fixture(
            request_state=self._get_db_enum(DBRequestState, RequestStates.PENDING.value),
            request_type=self._get_db_enum(DBRequestType, RequestTypes.REGISTRATION.value),
            fixture=registration_requests,
            cbsd_state=self._get_db_enum(DBCbsdState, CbsdStates.UNREGISTERED.value),
        )

        self.session.add_all(db_requests)
        self.session.commit()

        response_payload = self._create_response_payload_from_db_requests(
            response_type_name="registrationResponse",
            db_requests=db_requests,
            sas_response_code=sas_response_code
        )

        response, processor = self._prepare_response_and_processor(
            response_payload, "registrationResponse", processor_strategies["registrationRequest"])

        processor.process_response(db_requests, response, self.session)
        self.session.commit()

        states = [req.cbsd.state for req in db_requests]
        [self.assertTrue(state.name == expected_cbsd_state.value) for state in states]

    @parameterized.expand([
        (0, CbsdStates.UNREGISTERED),
        (400, CbsdStates.UNREGISTERED),
        (500, CbsdStates.UNREGISTERED),
    ])
    @responses.activate
    def test_cbsd_state_after_deregistration_response(self, sas_response_code, expected_cbsd_state):
        db_requests = self._create_db_requests_from_fixture(
            request_state=self._get_db_enum(DBRequestState, RequestStates.PENDING.value),
            request_type=self._get_db_enum(DBRequestType, RequestTypes.DEREGISTRATION.value),
            fixture=deregistration_requests,
            cbsd_state=self._get_db_enum(DBCbsdState, CbsdStates.UNREGISTERED.value)
        )
        self.session.add_all(db_requests)
        self.session.commit()

        registered_state = self._get_db_enum(DBCbsdState, CbsdStates.REGISTERED.value)
        self.session.query(DBCbsd).update({DBCbsd.state_id: registered_state.id})
        self.session.commit()

        response_payload = self._create_response_payload_from_db_requests(
            response_type_name="deregistrationResponse",
            db_requests=db_requests,
            sas_response_code=sas_response_code
        )

        response, processor = self._prepare_response_and_processor(
            response_payload, "deregistrationResponse", processor_strategies["deregistrationRequest"])

        processor.process_response(db_requests, response, self.session)
        self.session.commit()

        states = [req.cbsd.state for req in db_requests]
        [self.assertTrue(state.name == expected_cbsd_state.value) for state in states]

    @parameterized.expand([
        (zero_channels_for_one_cbsd, 0),
        (single_channel_for_one_cbsd, 1),
        (two_channels_for_one_cbsd, 2)
    ])
    @responses.activate
    def test_channels_created_after_spectrum_inquiry_response(self, response_fixture_payload, expected_channels_count):
        # Given
        db_requests = self._create_db_requests_from_fixture(
            request_state=self._get_db_enum(DBRequestState, RequestStates.PENDING.value),
            request_type=self._get_db_enum(DBRequestType, RequestTypes.SPECTRUM_INQUIRY.value),
            fixture=spectrum_inquiry_requests,
            cbsd_state=self._get_db_enum(DBCbsdState, CbsdStates.REGISTERED.value),
        )

        self.session.add_all(db_requests)
        self.session.commit()

        response, processor = self._prepare_response_and_processor(
            response_fixture_payload, "spectrumInquiryResponse", processor_strategies["spectrumInquiryRequest"])

        # When
        processor.process_response(db_requests, response, self.session)
        self.session.commit()

        # Then
        cbsd = self.session.query(DBCbsd).filter(DBCbsd.cbsd_id == "foo").first()
        self.assertEqual(expected_channels_count, len(cbsd.channels))

    @responses.activate
    def test_old_channels_deleted_after_spectrum_inquiry_response(self):
        # TODO cleanup tests (currently a lot of duplications)
        # Given
        db_requests = self._create_db_requests_from_fixture(
            request_state=self._get_db_enum(DBRequestState, RequestStates.PENDING.value),
            request_type=self._get_db_enum(DBRequestType, RequestTypes.SPECTRUM_INQUIRY.value),
            fixture=spectrum_inquiry_requests,
            cbsd_state=self._get_db_enum(DBCbsdState, CbsdStates.REGISTERED.value),
        )

        self.session.add_all(db_requests)
        self.session.commit()

        cbsd = self.session.query(DBCbsd).filter(DBCbsd.cbsd_id == "foo").first()
        self._create_channel(cbsd, 1, 2)

        self.assertEqual(1, len(cbsd.channels))

        response, processor = self._prepare_response_and_processor(
            zero_channels_for_one_cbsd, "spectrumInquiryResponse", processor_strategies["spectrumInquiryRequest"])

        # When
        processor.process_response(db_requests, response, self.session)
        self.session.commit()

        # Then
        self.assertEqual(0, len(cbsd.channels))

    @responses.activate
    def test_max_eirp_set_on_channel_on_grant_response(self):
        # Given
        cbsd_id = "foo"
        low_frequency = 1
        high_frequency = 2
        max_eirp = 3

        fixture = self._build_grant_request(cbsd_id, low_frequency, high_frequency, max_eirp)
        db_requests = self._create_db_requests_from_fixture(
            request_state=self._get_db_enum(DBRequestState, RequestStates.PENDING.value),
            request_type=self._get_db_enum(DBRequestType, RequestTypes.GRANT.value),
            fixture=[fixture],
            cbsd_state=self._get_db_enum(DBCbsdState, CbsdStates.REGISTERED.value),
        )

        self.session.add_all(db_requests)
        self.session.commit()

        cbsd = self.session.query(DBCbsd).filter(DBCbsd.cbsd_id == cbsd_id).first()
        channel = self._create_channel(cbsd, low_frequency, high_frequency)

        response_payload = self._create_response_payload_from_db_requests(
            response_type_name="grantResponse",
            db_requests=db_requests)
        response, processor = self._prepare_response_and_processor(
            response_payload, "grantResponse", processor_strategies["grantRequest"])

        # When
        processor.process_response(db_requests, response, self.session)
        self.session.commit()

        # Then
        self.assertEqual(max_eirp, channel.last_used_max_eirp)

    def _get_db_requests_and_response_payload(
            self, request_type_name, response_type_name, requests_fixtures, cbsd_state):
        db_requests = self._create_db_requests_from_fixture(
            request_state=self._get_db_enum(DBRequestState, RequestStates.PENDING.value),
            request_type=self._get_db_enum(DBRequestType, request_type_name),
            fixture=requests_fixtures,
            cbsd_state=cbsd_state,
        )

        self.session.add_all(db_requests)
        self.session.commit()

        response_payload = self._create_response_payload_from_db_requests(
            response_type_name=response_type_name,
            db_requests=db_requests)

        return db_requests, response_payload

    def _get_db_enum(self, data_type, name):
        return self.session.query(data_type).filter(data_type.name == name).first()

    @staticmethod
    def _prepare_response_and_processor(response_payload, response_type_name, processor_strategy):
        any_url = 'https://foo.com/foobar'
        responses.add(responses.GET, any_url, json=response_payload, status=200)
        response = requests.get(any_url)  # url and method don't matter, I'm just crafting a qualified response here

        processor = ResponseDBProcessor(
            response_type_name,
            process_responses_func=processor_strategy["process_responses"],
        )
        return response, processor

    def _generate_cbsd_from_request_json(self, request_payload: Dict, cbsd_state: DBCbsdState):
        cbsd_id = request_payload.get(CBSD_ID)
        fcc_id = request_payload.get(FCC_ID)
        user_id = request_payload.get(USER_ID)
        serial_number = request_payload.get(CBSD_SERIAL_NR)

        cbsd = DBCbsd(
            cbsd_id=cbsd_id,
            fcc_id=fcc_id,
            cbsd_serial_number=serial_number,
            user_id=user_id,
            eirp_capability=1.1,
            state=cbsd_state,
        )

        self.session.add(cbsd)
        self.session.commit()

        return cbsd

    @staticmethod
    def _build_grant_request(cbsd_id: str, low_frequency: int, high_frequency: int, max_eirp: int) -> Dict:
        return {
            "grantRequest": [
                {
                    "cbsdId": cbsd_id,
                    "operationParam": {
                        "maxEirp": max_eirp,
                        "operationFrequencyRange": {
                            "lowFrequency": low_frequency,
                            "highFrequency": high_frequency,
                        }
                    },
                }
            ]
        }

    def _create_channel(self, cbsd: DBCbsd, low_frequency: int, high_frequency: int) -> DBChannel:
        channel = DBChannel(
            cbsd=cbsd,
            low_frequency=low_frequency,
            high_frequency=high_frequency,
            channel_type="some_type",
            rule_applied="some_rule",
        )
        self.session.add(channel)
        self.session.commit()
        return channel

    @staticmethod
    def _get_request_type_from_fixture(fixture):
        return next(iter(fixture[0].keys()))

    # FIXME this function assumes flat structure
    def _create_db_requests_from_fixture(self, request_state, request_type, fixture, cbsd_state):
        request_type_name = self._get_request_type_from_fixture(fixture)
        reqs = [
            DBRequest(
                cbsd=self._generate_cbsd_from_request_json(r[request_type_name][0], cbsd_state),
                state=request_state,
                type=request_type,
                payload=r[request_type_name][0])
            for r in fixture
        ]
        return reqs

    @staticmethod
    def _create_response_payload_from_db_requests(response_type_name, db_requests, sas_response_code=0):
        response_payload = {response_type_name: []}
        for i, db_request in enumerate(db_requests):
            cbsd_id = db_request.cbsd.cbsd_id or str(i)
            response_json = {"response": {"responseCode": sas_response_code}, "cbsdId": cbsd_id}
            if db_request.payload.get(GRANT_ID, ""):
                response_json[GRANT_ID] = db_request.payload.get(GRANT_ID)
            elif response_type_name == request_response[RequestTypes.GRANT.value]:
                response_json[GRANT_ID] = f'test_grant_id_for_{db_request.cbsd_id}'
            response_payload[response_type_name].append(response_json)

        return response_payload
