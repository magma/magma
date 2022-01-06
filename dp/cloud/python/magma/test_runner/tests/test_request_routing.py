"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from concurrent.futures import ThreadPoolExecutor

import requests
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.db_testcase import DBTestCase
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
from magma.test_runner.config import TestConfig
from parameterized import parameterized

incorrect_request_payload = {"incorrect": ["payload"]}

config = TestConfig()


class RequestRoutingTestCase(DBTestCase):
    def setUp(self):
        super().setUp()
        DBInitializer(SessionManager(self.engine)).initialize()

    @parameterized.expand([
        ('registration', registration_requests[0], 1),
        ('registration', registration_requests[1], 2),
        ('deregistration', deregistration_requests[0], 1),
        ('relinquishment', relinquishment_requests[0], 1),
        ('grant', grant_requests[0], 1),
        ('heartbeat', heartbeat_requests[0], 1),
        ('spectrumInquiry', spectrum_inquiry_requests[0], 1),
    ])
    def test_cbsd_gets_response_from_sas(self, route, request_payload, expected_resp_len):
        # Given / When
        resp = requests.post(
            f'{config.CBSD_SAS_PROTOCOL_CONTROLLER_API_PREFIX}/{route}',
            json=request_payload,
            verify=False,  # noqa: S501
        )

        # Then
        sas_responses = resp.json()[f"{route}Response"]
        self.assertEqual(
            expected_resp_len, len(
                resp.json()[f"{route}Response"],
            ),
        )
        self.assertTrue(all([bool(resp) for resp in sas_responses]))

    @parameterized.expand([
        ('registration', incorrect_request_payload),
        ('deregistration', incorrect_request_payload),
        ('relinquishment', incorrect_request_payload),
        ('grant', incorrect_request_payload),
        ('heartbeat', incorrect_request_payload),
        ('spectrumInquiry', incorrect_request_payload),
    ])
    def test_dp_raises_400_when_payload_doesnt_pass_validation(self, route, payload):
        # Given / When
        resp = requests.post(
            f'{config.CBSD_SAS_PROTOCOL_CONTROLLER_API_PREFIX}/{route}',
            json=payload,
        )

        # Then
        self.assertEqual(400, resp.status_code)

    def test_cbsd_only_gets_response_from_sas_for_the_request_it_sent(self):
        # Given / When
        response_name = "registrationResponse"
        with ThreadPoolExecutor() as executor:
            resps = executor.map(
                lambda payload:
                requests.post(
                    f'{config.CBSD_SAS_PROTOCOL_CONTROLLER_API_PREFIX}/registration',
                    json=payload,
                ),
                registration_requests,
            )

        resps = list(resps)
        resps.sort(key=lambda resp: len(resp.json()[response_name]))

        # Then
        self.assertEqual(2, len(resps))
        self.assertEqual(1, len(resps[0].json()[response_name]))
        self.assertEqual(2, len(resps[1].json()[response_name]))
