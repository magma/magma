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
import os
from unittest import TestCase

import requests_mock
from magma.configuration_controller.request_router.exceptions import (
    RequestRouterError,
)
from magma.configuration_controller.request_router.request_router import (
    RequestRouter,
)
from magma.mappings.request_mapping import request_mapping


@requests_mock.Mocker()
class RequestRouterTestCase(TestCase):

    def setUp(self):
        super().setUp()
        self.sas_url = 'https://fake.sas.url'
        self.rc_ingest_url = 'https://fake.rc.url'
        self.router = RequestRouter(
            sas_url=self.sas_url,
            rc_ingest_url=self.rc_ingest_url,
            cert_path='fake/cert/path',
            ssl_key_path='fake/key/path',
            request_mapping=request_mapping,
            ssl_verify=False,
        )

    def test_router_forwarding_existing_sas_methods(self, mocker):
        # Given
        self._register_test_post_endpoints(
            mocker, f'{self.sas_url}/registration',
        )

        # When
        resp = self.router.post_to_sas(
            request_dict={"registrationRequest": [{"foo": "bar"}]},
        )

        # Then
        self.assertEqual(200, resp.status_code)

    def test_sas_response_is_forwarded_to_rc(self, mocker):
        # Given
        self._register_test_post_endpoints(
            mocker, f'{self.sas_url}/registration',
        )
        self._register_test_post_endpoints(mocker, self.rc_ingest_url)

        # When
        sas_resp = self.router.post_to_sas(
            request_dict={"registrationRequest": [{"foo": "bar"}]},
        )
        rc_resp = self.router.redirect_sas_response_to_radio_controller(
            sas_resp,
        )

        # Then
        self.assertEqual(200, rc_resp.status_code)

    def test_router_raises_exception_for_non_existing_sas_method(self, mocker):
        # Given
        self._register_test_post_endpoints(mocker, f'{self.sas_url}/test_post')

        # When / Then
        with self.assertRaises(RequestRouterError):
            self.router.post_to_sas(
                request_dict={"nonExistingSasMethod": [{}]},
            )

    def _register_test_post_endpoints(self, mocker, url):
        mocker.register_uri('POST', url, json=self._response_callback)

    @staticmethod
    def _get_from_fixture(fixture_rel_path):
        return os.path.join(os.path.dirname(__file__), "fixtures", fixture_rel_path)

    @staticmethod
    def _response_callback(request, context):
        context.status_code = 200
        return request.text
