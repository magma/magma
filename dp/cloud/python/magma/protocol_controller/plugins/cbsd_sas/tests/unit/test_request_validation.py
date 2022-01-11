from unittest import TestCase

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
from magma.protocol_controller.plugins.cbsd_sas.validators.deregistration_request import (
    DeregistrationRequestSchema,
)
from magma.protocol_controller.plugins.cbsd_sas.validators.grant_request import (
    GrantRequestSchema,
)
from magma.protocol_controller.plugins.cbsd_sas.validators.heartbeat_request import (
    HeartbeatRequestSchema,
)
from magma.protocol_controller.plugins.cbsd_sas.validators.registration_request import (
    RegistrationRequestSchema,
)
from magma.protocol_controller.plugins.cbsd_sas.validators.relinquishment_request import (
    RelinquishmentRequestSchema,
)
from magma.protocol_controller.plugins.cbsd_sas.validators.spectrum_inquiry_request import (
    SpectrumInquiryRequestSchema,
)
from marshmallow.exceptions import MarshmallowError
from parameterized import parameterized


class RequestValidationTestCase(TestCase):

    @parameterized.expand([
        (registration_requests[0], RegistrationRequestSchema),
        (deregistration_requests[0], DeregistrationRequestSchema),
        (relinquishment_requests[0], RelinquishmentRequestSchema),
        (grant_requests[0], GrantRequestSchema),
        (heartbeat_requests[0], HeartbeatRequestSchema),
        (spectrum_inquiry_requests[0], SpectrumInquiryRequestSchema),
    ])
    def test_req_passes_validation_if_fields_required_to_create_request_identifier_are_present(self, request, schema):
        # Given / When
        validated_request = schema().load(request)

        # Then
        self.assertEqual(1, len(list(validated_request.keys())))

    def test_request_with_more_than_one_key_fails_validation(self):
        # Given
        request_json = registration_requests[0].update(
            {"deregistrationRequest": [{"cbsdId": "foo"}]},
        )

        # When / Then
        with self.assertRaises(MarshmallowError):
            RegistrationRequestSchema().load(request_json)

    @parameterized.expand([
        ({"registrationRequest": [{}]}, RegistrationRequestSchema),
        ({"deregistrationRequest": [{}]}, DeregistrationRequestSchema),
        ({"relinquishmentRequest": [{}]}, RelinquishmentRequestSchema),
        ({"grantRequest": [{}]}, GrantRequestSchema),
        ({"heartbeatRequest": [{}]}, HeartbeatRequestSchema),
        ({"spectrumInquiryRequest": [{}]}, SpectrumInquiryRequestSchema),
    ])
    def test_request_without_fields_required_to_create_request_identifier_fails_validation(self, request_json, schema):
        # Given / When / Then
        self.assertRaisesMarshmallowError(request_json, schema)

    @parameterized.expand([
        (
            {
                "bad_registrationRequest": [
                    {"fccId": "foo", "cbsdSerialNumber": "bar"},
                ],
            }, RegistrationRequestSchema,
        ),
        (
            {
                "registrationRequest": [
                    {"bad_fccId": "foo", "cbsdSerialNumber": "bar"},
                ],
            }, RegistrationRequestSchema,
        ),
        (
            {
                "registrationRequest": [
                    {"fccId": "foo", "bad_cbsdSerialNumber": "bar"},
                ],
            }, RegistrationRequestSchema,
        ),
        (
            {
                "bad_deregistrationRequest": [
                    {"cbsdId": "foo"},
                ],
            }, DeregistrationRequestSchema,
        ),
        (
            {
                "deregistrationRequest": [
                    {"bad_cbsdId": "foo"},
                ],
            }, DeregistrationRequestSchema,
        ),
        (
            {
                "bad_relinquishmentRequest": [
                    {"cbsdId": "foo"},
                ],
            }, RelinquishmentRequestSchema,
        ),
        (
            {
                "relinquishmentRequest": [
                    {"bad_cbsdId": "foo"},
                ],
            }, RelinquishmentRequestSchema,
        ),
        ({"bad_grantRequest": [{"cbsdId": "foo"}]}, GrantRequestSchema),
        ({"grantRequest": [{"bad_cbsdId": "foo"}]}, GrantRequestSchema),
        (
            {
                "bad_heartbeatRequest": [
                    {"cbsdId": "foo", "grantId": "bar"},
                ],
            }, HeartbeatRequestSchema,
        ),
        (
            {
                "heartbeatRequest": [{
                    "bad_cbsdId": "foo",
                    "grantId": "bar",
                }],
            }, HeartbeatRequestSchema,
        ),
        (
            {
                "heartbeatRequest": [
                    {"cbsdId": "foo", "bad_grantId": "bar"},
                ],
            }, HeartbeatRequestSchema,
        ),
        (
            {
                "bad_spectrumInquiryRequest": [
                    {"cbsdId": "foo"},
                ],
            }, SpectrumInquiryRequestSchema,
        ),
        (
            {
                "spectrumInquiryRequest": [
                    {"bad_cbsdId": "foo"},
                ],
            }, SpectrumInquiryRequestSchema,
        ),
    ])
    def test_request_fails_validation_when_its_keys_are_malformed(self, request_json, schema):
        # Given / When / Then
        self.assertRaisesMarshmallowError(request_json, schema)

    @parameterized.expand([
        ("", RegistrationRequestSchema),
        ("", DeregistrationRequestSchema),
        ("", RelinquishmentRequestSchema),
        ("", GrantRequestSchema),
        ("", HeartbeatRequestSchema),
        ("", SpectrumInquiryRequestSchema),
    ])
    def test_request_in_a_non_json_serializable_format_fails_validation(self, request_json, schema):
        # Given / When / Then
        self.assertRaisesMarshmallowError(request_json, schema)

    def assertRaisesMarshmallowError(self, request_json, schema):
        with self.assertRaises(MarshmallowError):
            schema().load(request_json)
