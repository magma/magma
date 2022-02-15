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
    DBGrantState,
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
    two_channels_for_one_cbsd,
    zero_channels_for_one_cbsd,
)
from magma.mappings.request_response_mapping import request_response
from magma.mappings.types import (
    CbsdStates,
    GrantStates,
    RequestStates,
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


class DefaultResponseDBProcessorTestCase(LocalDBTestCase):
    def setUp(self):
        super().setUp()
        DBInitializer(SessionManager(self.engine)).initialize()

    @parameterized.expand([
        (REGISTRATION_REQ, registration_requests),
        (SPECTRUM_INQ_REQ, spectrum_inquiry_requests),
        (GRANT_REQ, grant_requests),
        (HEARTBEAT_REQ, heartbeat_requests),
        (RELINQUISHMENT_REQ, relinquishment_requests),
        (DEREGISTRATION_REQ, deregistration_requests),
    ])
    @responses.activate
    def test_processor_splits_sas_response_into_separate_db_objects_and_links_them_with_requests(
            self, request_type_name, requests_fixtures,
    ):
        # Given
        db_requests = self._create_db_requests(
            request_type_name, requests_fixtures,
        )
        response = self._prepare_response_from_db_requests(db_requests)

        # When
        self._process_response(
            request_type_name=request_type_name, response=response, db_requests=db_requests,
        )
        nr_of_requests = len(db_requests)

        # Then
        self._verify_requests_number_and_state(db_requests, nr_of_requests)
        self.assertEqual(2, self.session.query(DBRequestState).count())
        self.assertEqual(
            1, self.session.query(DBRequestType).filter(
                DBRequestType.name == request_type_name,
            ).count(),
        )

    @parameterized.expand([
        (
            GRANT_REQ, grant_requests, ResponseCodes.SUCCESS.value,
            GrantStates.GRANTED.value,
        ),
        (
            GRANT_REQ, grant_requests,
            ResponseCodes.INTERFERENCE.value, GrantStates.IDLE.value,
        ),
        (
            GRANT_REQ, grant_requests,
            ResponseCodes.GRANT_CONFLICT.value, GrantStates.IDLE.value,
        ),
        (
            GRANT_REQ, grant_requests,
            ResponseCodes.TERMINATED_GRANT.value, GrantStates.IDLE.value,
        ),
        (
            HEARTBEAT_REQ, heartbeat_requests,
            ResponseCodes.SUCCESS.value, GrantStates.AUTHORIZED.value,
        ),
        (
            HEARTBEAT_REQ, heartbeat_requests,
            ResponseCodes.TERMINATED_GRANT.value, GrantStates.IDLE.value,
        ),
        (
            HEARTBEAT_REQ, heartbeat_requests,
            ResponseCodes.SUSPENDED_GRANT.value, GrantStates.GRANTED.value,
        ),
        (
            HEARTBEAT_REQ, heartbeat_requests,
            ResponseCodes.UNSYNC_OP_PARAM.value, GrantStates.UNSYNC.value,
        ),
        (
            RELINQUISHMENT_REQ, relinquishment_requests,
            ResponseCodes.SUCCESS.value, GrantStates.IDLE.value,
        ),
    ])
    @responses.activate
    def test_grant_state_after_response(
            self, request_type_name, requests_fixtures, response_code, expected_grant_state_name,
    ):
        # Given
        db_requests = self._create_db_requests(
            request_type_name, requests_fixtures,
        )
        response = self._prepare_response_from_db_requests(
            db_requests, response_code=response_code,
        )

        # When
        self._process_response(
            request_type_name=request_type_name, response=response, db_requests=db_requests,
        )
        nr_of_requests = len(db_requests)

        # Then
        self._verify_requests_number_and_state(db_requests, nr_of_requests)
        self.assertListEqual(
            [expected_grant_state_name] * nr_of_requests,
            [g.state.name for g in self.session.query(DBGrant).all()],
        )

    @parameterized.expand([
        (GRANT_REQ, grant_requests, ResponseCodes.SUCCESS.value, 5, 5),
        (GRANT_REQ, grant_requests, ResponseCodes.INTERFERENCE.value, 5, 5),
        (GRANT_REQ, grant_requests, ResponseCodes.GRANT_CONFLICT.value, 5, 5),
        (GRANT_REQ, grant_requests, ResponseCodes.TERMINATED_GRANT.value, 5, 5),
        (HEARTBEAT_REQ, heartbeat_requests, ResponseCodes.SUCCESS.value, 5, 5),
        (
            HEARTBEAT_REQ, heartbeat_requests,
            ResponseCodes.TERMINATED_GRANT.value, 5, None,
        ),
        (
            HEARTBEAT_REQ, heartbeat_requests,
            ResponseCodes.SUSPENDED_GRANT.value, 5, 5,
        ),
        (
            HEARTBEAT_REQ, heartbeat_requests,
            ResponseCodes.UNSYNC_OP_PARAM.value, 5, 5,
        ),
        (RELINQUISHMENT_REQ, relinquishment_requests, 0, 5, None),
    ])
    @responses.activate
    def test_last_used_eirp_is_reset(
            self, request_type_name, requests_fixtures, response_code, last_used_eirp_value, final_last_used_eirp_value,
    ):
        """
        last_used_max_eirp should only be reset for an existing channel when grant goes into IDLE state due to
        relinquishmentRequest or heartbeatRequest. Grant requests should not affect this value
        """
        # Given
        db_requests = self._create_db_requests(
            request_type_name, requests_fixtures,
        )
        granted_state = self.session.query(DBGrantState).filter(
            DBGrantState.name == GrantStates.GRANTED.value,
        ).one()

        response = self._prepare_response_from_db_requests(
            db_requests, response_code=response_code,
        )
        response_payload = response.json()[request_response[request_type_name]]

        for resp in response_payload:
            cbsd = self.session.query(DBCbsd).filter(
                DBCbsd.cbsd_id == resp["cbsdId"],
            ).first()
            channel = self._create_channel(
                cbsd, 1, 2, last_used_max_eirp=last_used_eirp_value,
            )
            self._create_grant(
                resp["grantId"], channel=channel, cbsd=cbsd, state=granted_state,
            )

        grants_query = self.session.query(DBGrant)

        for grant in grants_query.all():
            self.assertEqual(
                last_used_eirp_value,
                grant.channel.last_used_max_eirp,
            )

        # When
        self._process_response(request_type_name, response, db_requests)

        # Then
        for grant in grants_query.all():
            self.assertEqual(
                final_last_used_eirp_value,
                grant.channel.last_used_max_eirp,
            )

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
            db_requests, sas_response_code,
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
            db_requests, sas_response_code,
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
            SPECTRUM_INQ_REQ, response_fixture_payload,
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
            SPECTRUM_INQ_REQ, zero_channels_for_one_cbsd,
        )

        # When
        self._process_response(SPECTRUM_INQ_REQ, response, db_requests)

        # Then
        self.assertEqual(0, len(cbsd.channels))

    @responses.activate
    def test_max_eirp_set_on_channel_on_grant_response(self):
        max_eirp = 3
        channel = self._setup_channel_and_process_grant_response(max_eirp)

        # Then
        self.assertEqual(max_eirp, channel.last_used_max_eirp)

    @responses.activate
    def test_assign_channel_to_grant_on_grant_response(self):
        # Given / When
        channel = self._setup_channel_and_process_grant_response()

        # Then
        count = self.session.query(DBGrant).filter(
            DBGrant.channel_id == channel.id,
        ).count()
        self.assertEqual(1, count)

    def _process_response(self, request_type_name, response, db_requests):
        processor = self._get_response_processor(request_type_name)

        processor.process_response(db_requests, response, self.session)
        self.session.commit()

    @staticmethod
    def _get_response_processor(req_type):
        return ResponseDBProcessor(
            request_response[req_type],
            process_responses_func=processor_strategies[req_type]["process_responses"],
        )

    def _verify_requests_number_and_state(self, db_requests, nr_of_requests, desired_state="processed"):
        self.assertEqual(nr_of_requests, self.session.query(DBRequest).count())
        self.assertListEqual(
            [r.id for r in db_requests], [
                _id for (
                    _id,
                ) in self.session.query(DBResponse.id).all()
            ],
        )
        self.assertListEqual(
            [desired_state] * nr_of_requests,
            [r.state.name for r in self.session.query(DBRequest).all()],
        )

    def _setup_channel_and_process_grant_response(self, max_eirp=None):
        # Given
        cbsd_id = "foo"
        low_frequency = 1
        high_frequency = 2

        fixture = self._build_grant_request(
            cbsd_id, low_frequency, high_frequency, max_eirp,
        )
        db_requests = self._create_db_requests(GRANT_REQ, [fixture])

        cbsd = self.session.query(DBCbsd).filter(
            DBCbsd.cbsd_id == cbsd_id,
        ).first()
        channel = self._create_channel(cbsd, low_frequency, high_frequency)

        response = self._prepare_response_from_db_requests(db_requests)

        self._process_response(
            request_type_name=GRANT_REQ,
            db_requests=db_requests, response=response,
        )

        return channel

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
            request_state=RequestStates.PENDING.value,
    ):
        db_requests = self._create_db_requests_from_fixture(
            request_state=request_state,
            request_type=request_type_name,
            fixture=requests_fixtures,
            cbsd_state=cbsd_state,
        )

        self.session.add_all(db_requests)
        self.session.commit()

        return db_requests

    def _get_db_enum(self, data_type, name):
        return self.session.query(data_type).filter(data_type.name == name).first()

    def _prepare_response_from_db_requests(self, db_requests, response_code=None):
        req_type = db_requests[0].type.name
        response_payload = self._create_response_payload_from_db_requests(
            response_type_name=request_response[req_type],
            db_requests=db_requests,
        )
        return self._prepare_response_from_payload(req_type, response_payload, response_code)

    @staticmethod
    def _prepare_response_from_payload(req_type, response_payload, response_code=None):
        if response_code is not None:
            for response_json in response_payload[request_response[req_type]]:
                response_json["response"]["responseCode"] = response_code
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
        last_used_max_eirp=None,
    ) -> DBChannel:
        channel = DBChannel(
            cbsd=cbsd,
            low_frequency=low_frequency,
            high_frequency=high_frequency,
            channel_type="some_type",
            rule_applied="some_rule",
            last_used_max_eirp=last_used_max_eirp,
        )
        self.session.add(channel)
        self.session.commit()
        return channel

    def _create_grant(self, grant_id, channel, cbsd, state):
        grant = DBGrant(
            channel=channel,
            cbsd=cbsd,
            state=state,
            grant_id=grant_id,
        )
        self.session.add(grant)
        self.session.commit()
        return grant

    def _create_db_requests_from_fixture(self, request_state, request_type, fixture, cbsd_state):
        db_requests = []
        for reqs in fixture:
            for req in reqs[request_type]:
                db_requests.append(
                    DBRequest(
                        cbsd=self._generate_cbsd_from_request_json(
                            req, self._get_db_enum(DBCbsdState, cbsd_state),
                        ),
                        state=self._get_db_enum(DBRequestState, request_state),
                        type=self._get_db_enum(DBRequestType, request_type),
                        payload=req,
                    ),
                )
        return db_requests

    @staticmethod
    def _create_response_payload_from_db_requests(response_type_name, db_requests, sas_response_code=0):
        response_payload = {response_type_name: []}
        for i, db_request in enumerate(db_requests):
            cbsd_id = db_request.cbsd.cbsd_id or str(i)
            response_json = {
                "response": {
                    "responseCode": sas_response_code,
                }, "cbsdId": cbsd_id,
            }
            if db_request.payload.get(GRANT_ID, ""):
                response_json[GRANT_ID] = db_request.payload.get(GRANT_ID)
            elif response_type_name == request_response[GRANT_REQ]:
                response_json[GRANT_ID] = f'test_grant_id_for_{db_request.cbsd_id}'
            response_payload[response_type_name].append(response_json)

        return response_payload
