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
import socket

import requests


def _load_certificate_file(file: str) -> bytes:
    dir_path = os.path.join(os.path.dirname(__file__))
    file_path = f'{dir_path}/{file}'
    with open(file_path, 'rb') as f:
        return f.read()


CERTIFICATES_DATE = '2022-07-20 12:00:00'

NO_CRL_CERT = _load_certificate_file(file='no_crl.crt')

VALID_CRL_CERT = _load_certificate_file(file='valid.crt')
VALID_CERT_CRLS_DATA = {
    'http://crl3.digicert.com/DigiCertTLSHybridECCSHA3842020CA1-1.crl': _load_certificate_file(file='valid_1.crl'),
    'http://crl4.digicert.com/DigiCertTLSHybridECCSHA3842020CA1-1.crl': _load_certificate_file(file='valid_2.crl'),
}

REVOKED_CRL_CERT = _load_certificate_file(file='revoked.crt')
REVOKED_CERT_CRLS_DATA = {
    'http://crl3.digicert.com/RapidSSLTLSDVRSAMixedSHA2562020CA-1.crl': _load_certificate_file(file='revoked_1.crl'),
    'http://crl4.digicert.com/RapidSSLTLSDVRSAMixedSHA2562020CA-1.crl': _load_certificate_file(file='revoked_2.crl'),
}

INVALID_CRL_CERT = b'invalid crl data'
INVALID_CERT_CRLS_DATA = {
    'http://crl3.digicert.com/DigiCertTLSHybridECCSHA3842020CA1-1.crl': b'invalid crl 1',
    'http://crl4.digicert.com/DigiCertTLSHybridECCSHA3842020CA1-1.crl': b'invalid crl 2',
}

CERT_SOCKET_EXCEPTION = socket.error('Socket error')
CRLS_CONNECTION_EXCEPTION = {
    'http://crl3.digicert.com/DigiCertTLSHybridECCSHA3842020CA1-1.crl': requests.exceptions.RequestException('Request error'),
    'http://crl4.digicert.com/DigiCertTLSHybridECCSHA3842020CA1-1.crl': requests.exceptions.RequestException('Request error'),
}
