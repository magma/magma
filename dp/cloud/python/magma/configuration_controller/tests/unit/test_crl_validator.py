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
from typing import Dict, List, Union
from unittest import TestCase, mock

import requests_mock
from freezegun import freeze_time
from magma.configuration_controller.crl_validator.crl_validator import (
    CRLValidator,
)
from magma.configuration_controller.request_router.request_router import (
    RequestRouter,
)
from magma.configuration_controller.tests.unit.fixtures.crl.certs import (
    CERT_SOCKET_EXCEPTION,
    CERTIFICATES_DATE,
    CRLS_CONNECTION_EXCEPTION,
    INVALID_CERT_CRLS_DATA,
    INVALID_CRL_CERT,
    NO_CRL_CERT,
    REVOKED_CERT_CRLS_DATA,
    REVOKED_CRL_CERT,
    VALID_CERT_CRLS_DATA,
    VALID_CRL_CERT,
)
from magma.mappings.request_mapping import request_mapping
from mocket import mocket
from parameterized import parameterized
from requests.exceptions import SSLError


@requests_mock.Mocker()
@freeze_time(CERTIFICATES_DATE)
class CRLValidatorTestCase(TestCase):
    sas_url = 'https://fake.sas.url/'

    @parameterized.expand([
        (NO_CRL_CERT, {}),
        (VALID_CRL_CERT, VALID_CERT_CRLS_DATA),
        (INVALID_CRL_CERT, {}),
        (VALID_CRL_CERT, INVALID_CERT_CRLS_DATA),
        (CERT_SOCKET_EXCEPTION, {}),
        (VALID_CRL_CERT, CRLS_CONNECTION_EXCEPTION),
    ])
    @mocket.mocketize
    def test_valid_certificates_passes(
            self,
            mocker: requests_mock.Mocker,
            certificate: bytes,
            crls: Dict[str, bytes],
    ) -> None:
        # Given
        self._set_mocket_cert(data=certificate)
        self._set_mocker_CRL(mocker=mocker, data=crls)
        validator = self._create_validator(urls=[self.sas_url])

        # When
        valid = validator.is_valid(url=self.sas_url)

        # Then
        self.assertIs(valid, True)

    @mocket.mocketize
    def test_exception_is_raised_on_revoked_certificate(self, mocker: requests_mock.Mocker) -> None:
        # Given
        self._set_mocket_cert(data=REVOKED_CRL_CERT)
        self._set_mocker_CRL(mocker=mocker, data=REVOKED_CERT_CRLS_DATA)
        validator = self._create_validator(urls=[self.sas_url])

        # When / Then
        with self.assertRaises(SSLError):
            validator.is_valid(url=self.sas_url)

    def test_not_prefetched_certificate(self, _) -> None:
        # Given
        validator = self._create_validator(urls=[])

        # When / Then
        with self.assertRaises(KeyError):
            validator.is_valid(url=self.sas_url)

    @mocket.mocketize
    def test_certificates_are_updated_in_the_background(self, mocker: requests_mock.Mocker) -> None:
        # Given
        self._set_mocket_cert(data=VALID_CRL_CERT)
        self._set_mocker_CRL(mocker=mocker, data=VALID_CERT_CRLS_DATA)

        validator = self._create_validator(urls=[self.sas_url])
        valid = validator.is_valid(url=self.sas_url)
        self.assertIs(valid, True)

        # When / Then
        self._set_mocket_cert(data=REVOKED_CRL_CERT)
        self._set_mocker_CRL(mocker=mocker, data=REVOKED_CERT_CRLS_DATA)

        # Ensure certs are not instant updated
        valid = validator.is_valid(url=self.sas_url)
        self.assertIs(valid, True)

        # Simulate BackgroundScheduler running update job.
        validator.update_certificates()

        # Then
        with self.assertRaises(SSLError):
            validator.is_valid(url=self.sas_url)

    @staticmethod
    def _create_validator(urls: List[str]) -> CRLValidator:
        """ Due to CRLValidator fetching certificates automatically in the background it is critical
        for it to be created after sockets and requests are mocked, and not on tests set up.

        Args:
            urls: list of urls that certificates will be fetched from

        Returns:
            CRLValidator instance
        """
        return CRLValidator(urls=urls)

    @staticmethod
    def _set_mocket_cert(data: Union[bytes, Exception]) -> None:
        def getpeercert(*args, **kwargs) -> bytes:
            if isinstance(data, Exception):
                raise data
            return data

        mocket.MocketSocket.getpeercert = getpeercert

    @staticmethod
    def _set_mocker_CRL(mocker: requests_mock.Mocker, data: Dict[str, Union[bytes, Exception]]) -> None:
        for url, response in data.items():
            if isinstance(response, Exception):
                mocker.register_uri('GET', url, exc=response)
            else:
                mocker.register_uri('GET', url, content=response)


@requests_mock.Mocker()
class RequestRouterWithCRLValidatorTestCase(TestCase):
    def setUp(self) -> None:
        super().setUp()
        self.sas_url = 'https://fake.sas.url'
        self.rc_ingest_url = 'https://fake.rc.url'
        self.crl_validator_mock = mock.MagicMock(spec=CRLValidator)
        self.router = RequestRouter(
            sas_url=self.sas_url,
            rc_ingest_url=self.rc_ingest_url,
            cert_path='fake/cert/path',
            ssl_key_path='fake/key/path',
            request_mapping=request_mapping,
            ssl_verify=False,
            crl_validator=self.crl_validator_mock,
        )

    @mocket.mocketize
    def test_request_router_uses_crl_validator(self, mocker: requests_mock.Mocker) -> None:
        # Given
        mocker.register_uri('POST', f'{self.sas_url}/registration')

        # When
        self.router.post_to_sas(request_dict={"registrationRequest": [{"foo": "bar"}]})

        # Then
        self.crl_validator_mock.is_valid.assert_called_once_with(url=self.sas_url)
