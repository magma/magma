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
import socket
import ssl
from typing import List

import requests
from cryptography import x509
from cryptography.hazmat.backends import default_backend
from requests.exceptions import SSLError

_x509_backend = default_backend()


def get_certificate(
        hostname: str,
        port: int = 443,
) -> x509.Certificate:
    """ Retrieve certificate of given host.

    Args:
        hostname: host from which certificate will be retrieved
        port: SSL port of the host

    Returns:
        x509.Certificate
    """
    context = ssl.SSLContext(protocol=ssl.PROTOCOL_SSLv23)
    context.minimum_version = ssl.TLSVersion.TLSv1_2

    with socket.create_connection((hostname, port)) as sock:
        sock = context.wrap_socket(sock, server_hostname=hostname)
        binary_certificate = sock.getpeercert(binary_form=True)
        sock.shutdown(socket.SHUT_RDWR)

    return x509.load_der_x509_certificate(data=binary_certificate, backend=_x509_backend)


def get_certificate_crls(
        certificate: x509.Certificate,
) -> List[x509.CertificateRevocationList]:
    """ Extract CRLs from given certificate.

    Args:
        certificate: SSL certificate of which CRLs will be retrieved from

    Returns:
        list of Certificate Revocation Lists
    """

    def get_crl(url: str) -> x509.CertificateRevocationList:
        data = requests.get(url=url).content
        return x509.load_der_x509_crl(data=data, backend=_x509_backend)

    try:
        distribution_points = certificate.extensions.get_extension_for_oid(
            oid=x509.ExtensionOID.CRL_DISTRIBUTION_POINTS,
        )
    except x509.ExtensionNotFound:
        return []

    crl_urls = []
    for distribution_point in distribution_points.value:
        for crl in distribution_point.full_name:
            crl_urls.append(crl.value)

    return [get_crl(url=url) for url in crl_urls]


def is_certificate_revoked(
        certificate: x509.Certificate,
        crls: List[x509.CertificateRevocationList],
) -> bool:
    """ Check if given certificate is revoked by any of given CRLs.

    Args:
        certificate: SSL certificate to be checked
        crls: list of Certificate Revocation Lists to check certificate against

    Returns:
        bool: False if certificate is not revoked

    Raises:
        SSLError: if certificate is revoked
    """
    for crl in crls:
        revoked = crl.get_revoked_certificate_by_serial_number(certificate.serial_number)
        if revoked:
            msg = f'[SSL: ERR_CERT_REVOKED] ' \
                  f'Certificate {revoked.serial_number} revoked on {revoked.revocation_date}'
            raise SSLError(msg)
    return False
