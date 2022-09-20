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
from typing import Dict

import requests
import responses
from magma.configuration_controller.response_processor.response_db_processor import (
    ResponseDBProcessor,
)
from magma.configuration_controller.response_processor.strategies.response_processing import (
    unset_frequency,
)
from magma.configuration_controller.response_processor.strategies.strategies_mapping import (
    processor_strategies,
)
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBGrant,
    DBRequest,
    DBRequestType,
)
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.fixtures.fake_requests.deregistration_requests import (
    deregistration_requests,
)
from magma.fixtures.fake_requests.grant_requests import grant_requests
from magma.fixtures.fake_requests.heartbeat_requests import heartbeat_requests
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
    single_channel_for_one_cbsd_with_no_max_eirp,
    two_channels_for_one_cbsd,
    zero_channels_for_one_cbsd,
)
from magma.fluentd_client.client import FluentdClient
from magma.mappings.request_response_mapping import request_response
from magma.mappings.types import (
    CbsdStates,
    GrantStates,
    RequestTypes,
    ResponseCodes,
)
from parameterized import parameterized

CBSD_SERIAL_NR = "cbsdSerialNumber"
FCC_ID = "fccId"
USER_ID = "userId"
CBSD_ID = "cbsdId"
GRANT_ID = "grantId"

REGISTRATION_REQ = RequestTypes.REGISTRATION.value
DEREGISTRATION_REQ = RequestTypes.DEREGISTRATION.value
RELINQUISHMENT_REQ = RequestTypes.RELINQUISHMENT.value
HEARTBEAT_REQ = RequestTypes.HEARTBEAT.value
GRANT_REQ = RequestTypes.GRANT.value
SPECTRUM_INQ_REQ = RequestTypes.SPECTRUM_INQUIRY.value
DEFAULT_MAX_EIRP = 37

_fake_requests_map = {
    REGISTRATION_REQ: registration_requests,
    SPECTRUM_INQ_REQ: spectrum_inquiry_requests,
    GRANT_REQ: grant_requests,
    HEARTBEAT_REQ: heartbeat_requests,
    RELINQUISHMENT_REQ: relinquishment_requests,
    DEREGISTRATION_REQ: deregistration_requests,
}


class DefaultResponseDBProcessorTestCase(LocalDBTestCase):
    def setUp(self):
        super().setUp()
        DBInitializer(SessionManager(self.engine)).initialize()

    @parameterized.expand([
        (REGISTRATION_REQ,),
        (SPECTRUM_INQ_REQ,),
        (GRANT_REQ,),
        (HEARTBEAT_REQ,),
        (RELINQUISHMENT_REQ,),
        (DEREGISTRATION_REQ,),
    ])
    @responses.activate
    def test_processor_splits_sas_response_into_separate_db_objects_and_links_them_with_requests(
            self, request_type_name,
    ):
        # Given
        requests_fixtures = _fake_requests_map[request_type_name]
        db_requests = self._create_db_requests(
            request_type_name, requests_fixtures,
        )
        response = self._prepare_response_from_db_requests(db_requests=db_requests)

        # When
        self._process_response(
            request_type_name=request_type_name, response=response, db_requests=db_requests,
        )

        # Then
        self._verify_processed_requests_were_deleted()
        self.assertEqual(
            1, self.session.query(DBRequestType).filter(
                DBRequestType.name == request_type_name,
            ).count(),
        )

    @parameterized.expand([
        (GRANT_REQ, ResponseCodes.SUCCESS.value, None, GrantStates.GRANTED.value),
        (GRANT_REQ, ResponseCodes.INTERFERENCE.value, None, None),
        (GRANT_REQ, ResponseCodes.GRANT_CONFLICT.value, ['grant1', 'grant2'], GrantStates.UNSYNC.value),
        (GRANT_REQ, ResponseCodes.TERMINATED_GRANT.value, None, None),
        (HEARTBEAT_REQ, ResponseCodes.SUCCESS.value, None, GrantStates.AUTHORIZED.value),
        (HEARTBEAT_REQ, ResponseCodes.TERMINATED_GRANT.value, None, None),
        (HEARTBEAT_REQ, ResponseCodes.SUSPENDED_GRANT.value, None, GrantStates.GRANTED.value),
        (HEARTBEAT_REQ, ResponseCodes.UNSYNC_OP_PARAM.value, None, GrantStates.UNSYNC.value),
        (RELINQUISHMENT_REQ, ResponseCodes.SUCCESS.value, None, None),
    ])
    @responses.activate
    def test_grant_state_after_response(
            self, request_type_name, response_code, response_data, expected_grant_state_name,
    ):
        # Given
        requests_fixtures = _fake_requests_map[request_type_name]
        db_requests = self._create_db_requests(
            request_type_name, requests_fixtures,
        )
        response = self._prepare_response_from_db_requests(
            db_requests=db_requests, response_code=response_code, response_data=response_data,
        )

        # When
        self._process_response(
            request_type_name=request_type_name, response=response, db_requests=db_requests,
        )
        expected_grant_state = [expected_grant_state_name] if expected_grant_state_name else []
        nr_of_requests = len(db_requests)

        # Then
        self._verify_processed_requests_were_deleted()
        self.assertListEqual(
            expected_grant_state * nr_of_requests,
            [g.state.name for g in self.session.query(DBGrant).all()],
        )

    @parameterized.expand([
        (GRANT_REQ, ResponseCodes.SUCCESS.value, None, [GrantStates.GRANTED.value]),
        (GRANT_REQ, ResponseCodes.INTERFERENCE.value, None, []),
        (
            GRANT_REQ, ResponseCodes.GRANT_CONFLICT.value,
            ['test_grant_id_for_1', 'test_grant_id_for_2'],
            [GrantStates.GRANTED.value, GrantStates.UNSYNC.value],
        ),
        (GRANT_REQ, ResponseCodes.TERMINATED_GRANT.value, None, []),
        (HEARTBEAT_REQ, ResponseCodes.SUCCESS.value, None, [GrantStates.AUTHORIZED.value]),
        (HEARTBEAT_REQ, ResponseCodes.TERMINATED_GRANT.value, None, []),
        (HEARTBEAT_REQ, ResponseCodes.SUSPENDED_GRANT.value, None, [GrantStates.GRANTED.value]),
        (HEARTBEAT_REQ, ResponseCodes.UNSYNC_OP_PARAM.value, None, [GrantStates.UNSYNC.value]),
        (RELINQUISHMENT_REQ, ResponseCodes.SUCCESS.value, None, []),
    ])
    @responses.activate
    def test_preexisting_grant_state_after_response(
            self, request_type_name, response_code, response_data, expected_grants_states,
    ):
        # Given
        requests_fixtures = [_fake_requests_map[request_type_name][0]]
        db_requests = self._create_db_requests(request_type_name, requests_fixtures)

        grant_id = db_requests[0].payload.get('grantId') or "test_grant_id_for_1"
        grant = DBGrant(
            cbsd=db_requests[0].cbsd,
            state_id=1,
            grant_id=grant_id,
            low_frequency=3560000000,
            high_frequency=3580000000,
            max_eirp=15,
        )
        self.session.add(grant)
        self.session.commit()

        response = self._prepare_response_from_db_requests(
            db_requests=db_requests, response_code=response_code, response_data=response_data,
        )

        # When
        self._process_response(
            request_type_name=request_type_name, response=response, db_requests=db_requests,
        )

        # Then
        self.assertListEqual(
            expected_grants_states,
            [g.state.name for g in self.session.query(DBGrant).order_by(DBGrant.id).all()],
        )
        self.assertEqual(self.session.query(DBCbsd).count(), 1)

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
        # Given
        db_requests = self._create_db_requests(
            REGISTRATION_REQ, registration_requests,
        )
        response = self._prepare_response_from_db_requests(
            db_requests=db_requests, response_code=sas_response_code,
        )

        # When
        self._process_response(
            request_type_name=REGISTRATION_REQ, response=response, db_requests=db_requests,
        )
        states = [req.cbsd.state for req in db_requests]

        # Then
        [
            self.assertTrue(state.name == expected_cbsd_state.value)
            for state in states
        ]

    @parameterized.expand([
        (0, CbsdStates.UNREGISTERED),
        (400, CbsdStates.UNREGISTERED),
        (500, CbsdStates.UNREGISTERED),
    ])
    @responses.activate
    def test_cbsd_state_after_deregistration_response(self, sas_response_code, expected_cbsd_state):
        # Given
        db_requests = self._create_db_requests(
            DEREGISTRATION_REQ, deregistration_requests,
        )
        self._set_cbsds_to_state(CbsdStates.REGISTERED.value)
        response = self._prepare_response_from_db_requests(
            db_requests=db_requests, response_code=sas_response_code,
        )

        # When
        self._process_response(
            request_type_name=DEREGISTRATION_REQ, response=response, db_requests=db_requests,
        )
        states = [req.cbsd.state for req in db_requests]

        # Then
        [
            self.assertTrue(state.name == expected_cbsd_state.value)
            for state in states
        ]

    @parameterized.expand([
        (zero_channels_for_one_cbsd, 0),
        (single_channel_for_one_cbsd, 1),
        (two_channels_for_one_cbsd, 2),
    ])
    @responses.activate
    def test_channels_created_after_spectrum_inquiry_response(self, response_fixture_payload, expected_channels_count):
        # Given
        db_requests = self._create_db_requests(
            SPECTRUM_INQ_REQ, spectrum_inquiry_requests,
        )
        response = self._prepare_response_from_payload(
            response_fixture_payload,
        )

        # When
        self._process_response(
            request_type_name=SPECTRUM_INQ_REQ, response=response, db_requests=db_requests,
        )

        # Then
        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_id == "foo",
        ).first()
        self.assertEqual(expected_channels_count, len(cbsd.channels))

    @responses.activate
    def test_channel_created_with_default_max_eirp(self):
        # Given
        db_requests = self._create_db_requests(SPECTRUM_INQ_REQ, spectrum_inquiry_requests)
        response = self._prepare_response_from_payload(single_channel_for_one_cbsd_with_no_max_eirp)

        # When
        self._process_response(request_type_name=SPECTRUM_INQ_REQ, response=response, db_requests=db_requests)

        # Then
        cbsd = self.session.query(DBCbsd).filter(DBCbsd.cbsd_id == "foo").first()
        self.assertEqual(DEFAULT_MAX_EIRP, cbsd.channels[0]["max_eirp"])

    @responses.activate
    def test_old_channels_deleted_after_spectrum_inquiry_response(self):
        # Given
        db_requests = self._create_db_requests(
            SPECTRUM_INQ_REQ, spectrum_inquiry_requests,
        )
        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_id == "foo",
        ).first()
        self._create_channel(cbsd, 1, 2)

        self.assertEqual(1, len(cbsd.channels))

        response = self._prepare_response_from_payload(
            zero_channels_for_one_cbsd,
        )

        # When
        self._process_response(SPECTRUM_INQ_REQ, response, db_requests)

        # Then
        self.assertEqual(0, len(cbsd.channels))

    @responses.activate
    def test_available_frequencies_deleted_after_spectrum_inquiry_response(self):
        # Given
        db_requests = self._create_db_requests(SPECTRUM_INQ_REQ, spectrum_inquiry_requests)
        cbsd = self.session.query(DBCbsd).filter(DBCbsd.cbsd_id == "foo").first()

        response = self._prepare_response_from_payload(zero_channels_for_one_cbsd)

        # When
        self._process_response(SPECTRUM_INQ_REQ, response, db_requests)

        # Then
        self.assertIsNone(cbsd.available_frequencies)

    @responses.activate
    def test_channel_params_set_on_grant_response(self):
        # Given
        cbsd_id = "foo"
        low_frequency = 3560000000
        high_frequency = 3561000000
        max_eirp = 3

        fixture = self._build_grant_request(
            cbsd_id, low_frequency, high_frequency, max_eirp,
        )
        db_requests = self._create_db_requests(GRANT_REQ, [fixture])

        response = self._prepare_response_from_db_requests(db_requests=db_requests)

        # When
        self._process_response(
            request_type_name=GRANT_REQ,
            db_requests=db_requests, response=response,
        )

        # Then
        grant = self.session.query(DBGrant).first()
        self.assertEqual(low_frequency, grant.low_frequency)
        self.assertEqual(high_frequency, grant.high_frequency)
        self.assertEqual(max_eirp, grant.max_eirp)

    @parameterized.expand([
        (REGISTRATION_REQ, ResponseCodes.DEREGISTER.value, None, CbsdStates.UNREGISTERED.value),
        (SPECTRUM_INQ_REQ, ResponseCodes.DEREGISTER.value, None, CbsdStates.UNREGISTERED.value),
        (GRANT_REQ, ResponseCodes.DEREGISTER.value, None, CbsdStates.UNREGISTERED.value),
        (HEARTBEAT_REQ, ResponseCodes.DEREGISTER.value, None, CbsdStates.UNREGISTERED.value),
        (RELINQUISHMENT_REQ, ResponseCodes.DEREGISTER.value, None, CbsdStates.UNREGISTERED.value),
        (DEREGISTRATION_REQ, ResponseCodes.DEREGISTER.value, None, CbsdStates.UNREGISTERED.value),
        (SPECTRUM_INQ_REQ, ResponseCodes.INVALID_VALUE.value, [CBSD_ID], CbsdStates.UNREGISTERED.value),
        (GRANT_REQ, ResponseCodes.INVALID_VALUE.value, [CBSD_ID], CbsdStates.UNREGISTERED.value),
        (HEARTBEAT_REQ, ResponseCodes.INVALID_VALUE.value, [CBSD_ID], CbsdStates.UNREGISTERED.value),
        (RELINQUISHMENT_REQ, ResponseCodes.INVALID_VALUE.value, [CBSD_ID], CbsdStates.UNREGISTERED.value),
        (DEREGISTRATION_REQ, ResponseCodes.INVALID_VALUE.value, [CBSD_ID], CbsdStates.UNREGISTERED.value),
        (SPECTRUM_INQ_REQ, ResponseCodes.INVALID_VALUE.value, [GRANT_ID], CbsdStates.UNREGISTERED.value),
        (GRANT_REQ, ResponseCodes.INVALID_VALUE.value, [GRANT_ID], CbsdStates.UNREGISTERED.value),
        (HEARTBEAT_REQ, ResponseCodes.INVALID_VALUE.value, [GRANT_ID], CbsdStates.UNREGISTERED.value),
        (RELINQUISHMENT_REQ, ResponseCodes.INVALID_VALUE.value, [GRANT_ID], CbsdStates.UNREGISTERED.value),
        (DEREGISTRATION_REQ, ResponseCodes.INVALID_VALUE.value, [GRANT_ID], CbsdStates.UNREGISTERED.value),
        (SPECTRUM_INQ_REQ, ResponseCodes.INVALID_VALUE.value, None, CbsdStates.UNREGISTERED.value),
        (GRANT_REQ, ResponseCodes.INVALID_VALUE.value, None, CbsdStates.UNREGISTERED.value),
        (HEARTBEAT_REQ, ResponseCodes.INVALID_VALUE.value, None, CbsdStates.UNREGISTERED.value),
        (RELINQUISHMENT_REQ, ResponseCodes.INVALID_VALUE.value, None, CbsdStates.UNREGISTERED.value),
    ])
    @responses.activate
    def test_cbsd_state_after_unsuccessful_response_code(self, request_type, response_code, response_data, expected_cbsd_sate):
        # Given
        request_fixture = _fake_requests_map[request_type]
        db_requests = self._create_db_requests(
            request_type, request_fixture,
        )
        self._set_cbsds_to_state(CbsdStates.REGISTERED.value)
        response = self._prepare_response_from_db_requests(
            db_requests, response_code=response_code, response_data=response_data,
        )

        # When
        self._process_response(
            request_type_name=request_type, response=response, db_requests=db_requests,
        )
        states = [req.cbsd.state for req in db_requests]

        # Then
        [
            self.assertTrue(state.name == expected_cbsd_sate)
            for state in states
        ]

    @parameterized.expand([
        (None, 3560000000, 3580000000, None),
        ([0b1111, 0b110, 0b1100, 0b1010], 3562500000, 3567500000, [0b0111, 0b110, 0b1100, 0b1010]),  # 5 MHz bw
        ([0b0, 0b110, 0b1100, 0b1010], 3550000000, 3560000000, [0b0, 0b100, 0b1100, 0b1010]),  # 10 MHz bw
        ([0b0, 0b110, 0b1111100, 0b1010], 3572500000, 3587500000, [0b0, 0b110, 0b0111100, 0b1010]),  # 15 MHz bw
        ([0b0, 0b110, 0b1111100, 0b10101], 3560000000, 3580000000, [0b0, 0b110, 0b1111100, 0b00101]),  # 20 MHz bw
    ])
    def test_unset_frequency(self, orig_avail_freqs, low_freq, high_freq, expected_avail_freq):
        # Given
        cbsd = DBCbsd(
            cbsd_id="some_cbsd_id",
            fcc_id="some_fcc_id",
            cbsd_serial_number="some_serial_number",
            user_id="some_user_id",
            state_id=1,
            desired_state_id=1,
            available_frequencies=orig_avail_freqs,
        )

        grant = DBGrant(
            cbsd=cbsd,
            state_id=1,
            grant_id="some_grant_id",
            low_frequency=low_freq,
            high_frequency=high_freq,
            max_eirp=15,
        )

        self.session.add_all([cbsd, grant])
        self.session.commit()

        # When
        unset_frequency(grant)
        cbsd = self.session.query(DBCbsd).filter(DBCbsd.id == cbsd.id).scalar()

        # Then
        self.assertEqual(expected_avail_freq, cbsd.available_frequencies)

    def _process_response(self, request_type_name, response, db_requests):
        processor = self._get_response_processor(request_type_name)

        processor.process_response(db_requests, response, self.session)
        self.session.commit()

    @staticmethod
    def _get_response_processor(req_type):
        return ResponseDBProcessor(
            request_response[req_type],
            process_responses_func=processor_strategies[req_type]["process_responses"],
            fluentd_client=FluentdClient(),
        )

    def _verify_processed_requests_were_deleted(self):
        self.assertEqual(0, self.session.query(DBRequest).count())

    def _set_cbsds_to_state(self, state_name):
        registered_state = self._get_db_enum(DBCbsdState, state_name)
        self.session.query(DBCbsd).update(
            {DBCbsd.state_id: registered_state.id},
        )
        self.session.commit()

    def _create_db_requests(
            self,
            request_type_name,
            requests_fixtures,
            cbsd_state=CbsdStates.UNREGISTERED.value,
    ):
        db_requests = self._create_db_requests_from_fixture(
            request_type=request_type_name,
            fixture=requests_fixtures,
            cbsd_state=cbsd_state,
        )

        self.session.add_all(db_requests)
        self.session.commit()

        return db_requests

    def _get_db_enum(self, data_type, name):
        return self.session.query(data_type).filter(data_type.name == name).first()

    def _prepare_response_from_db_requests(self, db_requests, response_code=ResponseCodes.SUCCESS.value, response_data=None):
        req_type = db_requests[0].type.name
        response_payload = self._create_response_payload_from_db_requests(
            response_type_name=request_response[req_type],
            db_requests=db_requests,
            sas_response_code=response_code,
            sas_response_data=response_data,
        )
        return self._prepare_response_from_payload(response_payload)

    @staticmethod
    def _prepare_response_from_payload(response_payload):
        any_url = 'https://foo.com/foobar'
        responses.add(
            responses.GET, any_url,
            json=response_payload, status=200,
        )
        # url and method don't matter, I'm just crafting a qualified response here
        return requests.get(any_url)

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
            state=cbsd_state,
            desired_state=cbsd_state,
            available_frequencies=[0b11111100, 0b11110111, 0b11001111, 0b11110001],
        )

        self.session.add(cbsd)
        self.session.commit()

        return cbsd

    @staticmethod
    def _build_grant_request(cbsd_id: str, low_frequency: int, high_frequency: int, max_eirp: int) -> Dict:
        return {
            GRANT_REQ: [
                {
                    "cbsdId": cbsd_id,
                    "operationParam": {
                        "maxEirp": max_eirp,
                        "operationFrequencyRange": {
                            "lowFrequency": low_frequency,
                            "highFrequency": high_frequency,
                        },
                    },
                },
            ],
        }

    def _create_channel(
        self,
        cbsd: DBCbsd,
        low_frequency: int,
        high_frequency: int,
    ) -> None:
        channels = [{
            "low_frequency": low_frequency,
            "high_frequency": high_frequency,
        }]
        cbsd.channels = channels
        self.session.commit()

    def _create_db_requests_from_fixture(self, request_type, fixture, cbsd_state):
        db_requests = []
        for reqs in fixture:
            for req in reqs[request_type]:
                db_requests.append(
                    DBRequest(
                        cbsd=self._generate_cbsd_from_request_json(
                            req, self._get_db_enum(DBCbsdState, cbsd_state),
                        ),
                        type=self._get_db_enum(DBRequestType, request_type),
                        payload=req,
                    ),
                )
        return db_requests

    @staticmethod
    def _create_response_payload_from_db_requests(response_type_name, db_requests, sas_response_code=0, sas_response_data=None):
        response_payload = {response_type_name: []}
        for i, db_request in enumerate(db_requests):
            response_json = {
                "response": {
                    "responseCode": sas_response_code,
                },
            }

            if sas_response_data:
                response_json["response"]["responseData"] = sas_response_data
            else:
                cbsd_id = db_request.cbsd.cbsd_id or str(i)
                response_json["cbsdId"] = cbsd_id

            if db_request.payload.get(GRANT_ID, ""):
                response_json[GRANT_ID] = db_request.payload.get(GRANT_ID)
            elif response_type_name == request_response[GRANT_REQ]:
                response_json[GRANT_ID] = f'test_grant_id_for_{db_request.cbsd_id}'
            response_payload[response_type_name].append(response_json)

        return response_payload
