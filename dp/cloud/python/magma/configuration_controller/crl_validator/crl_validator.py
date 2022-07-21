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
import threading
import time
from typing import Dict, List
from urllib.parse import urlparse

from cryptography import x509
from magma.configuration_controller.crl_validator.certificate import (
    get_certificate,
    get_certificate_crls,
    is_certificate_revoked,
)


def get_host(url: str) -> str:
    """ Get host from url.

    Args:
        url: url string

    Returns:
        str: host from given url
    """
    return urlparse(url=url).netloc


class CertificatesUpdaterThread(threading.Thread):
    """
    This class is responsible for periodical update of certificates and CRLs for SSLValidator
    using shared certificates and CRL dicts. Data is updated in the background to not delay the main application.
    """

    def __init__(
            self,
            certificates: Dict[str, x509.Certificate],
            crls: Dict[int, List[x509.CertificateRevocationList]],
            update_rate: int,
            *args,
            **kwargs,
    ) -> None:
        """
        Args:
            certificates: certificates dict that will be updated by this class
            crls: CRLs dict that will be updated by this class
            update_rate: time in seconds between each update
            *args: args
            **kwargs: kwargs
        """
        super().__init__(*args, **kwargs)
        self.daemon = True  # Die if main app exits.

        self._certificates = certificates
        self._crls = crls
        self._update_rate = update_rate

    def run(self) -> None:
        """ Start thread.

        Returns: None
        """
        while True:
            self._update_certificates()
            time.sleep(self._update_rate)

    def _update_certificates(self) -> None:
        """ Fetch new certificates and CRLs for all hosts
        and update self._certificates and self._crls dicts.

        Returns: None
        """
        for host in self._certificates.keys():
            certificate = get_certificate(hostname=host)
            crls = get_certificate_crls(certificate=certificate)

            self._certificates[host] = certificate
            self._crls[certificate.serial_number] = crls


class CRLValidator(object):
    """
    This class is responsible for validating given urls' SSL certs with their respective CRLs.
    """

    def __init__(
            self, urls: List[str], certificates_update_rate: int = 300,
    ) -> None:
        """
        Args:
            urls: list of urls that we should prefetch certificates and CRLs for
            certificates_update_rate: time in seconds between certificates and CRLs update
        """
        hosts = [get_host(url=url) for url in urls]

        # Certificates and CRLs dicts are shared with separate updater thread.
        # All we want is to have recent certs for given hosts, so usual thread related issues don't really bother us.
        self._certificates = dict.fromkeys(hosts)
        self._crls = {}

        # Start updater thread to update certificates and CRLs in the background.
        self._updater_thread = CertificatesUpdaterThread(
            certificates=self._certificates,
            crls=self._crls,
            update_rate=certificates_update_rate,
        )
        self._updater_thread.start()

    def is_valid(self, url: str) -> bool:
        """ Check if given url's SSL certificate is not revoked by its Certificate Revocation Lists.

        Args:
            url: url string

        Returns:
            bool: False if certificate is not revoked

        Raises:
            SSLError: if certificate is revoked
        """
        host = get_host(url=url)
        certificate = self._get_certificate(hostname=host)
        crls = self._get_certificate_crls(certificate=certificate)
        return not is_certificate_revoked(certificate=certificate, crls=crls)

    def _get_certificate(self, hostname: str) -> x509.Certificate:
        """ Get cached certificate for given host, fetch new if it does not exist.

        Args:
            hostname: host for which certificate was issued

        Returns:
            SSL certificate
        """
        if not self._certificates.get(hostname):
            self._certificates[hostname] = get_certificate(hostname=hostname)
        return self._certificates[hostname]

    def _get_certificate_crls(
            self, certificate: x509.Certificate,
    ) -> List[x509.CertificateRevocationList]:
        """ Get cached CertificateRevocationLists for given certificate, fetch new if they does not exist.

        Args:
            certificate: certificate that CRLs were attached to

        Returns:
            list of Certificate Revocation Lists
        """
        serial_number = certificate.serial_number
        if not self._crls.get(serial_number):
            self._crls[serial_number] = get_certificate_crls(certificate=certificate)
        return self._crls[serial_number]
