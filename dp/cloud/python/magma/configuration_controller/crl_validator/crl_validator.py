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
import logging
from typing import List
from urllib.parse import urlparse

from magma.configuration_controller.crl_validator.certificate import (
    get_certificate,
    get_certificate_crls,
    is_certificate_revoked,
)
from rwmutex import RWLock

logger = logging.getLogger(__name__)

CERTIFICATE_UNAVAILABLE = 'certificate unavailable'


def get_host(url: str) -> str:
    """ Get host from url.

    Args:
        url: url string

    Returns:
        str: host from given url
    """
    return urlparse(url=url).netloc


class CRLValidator(object):
    """
    This class is responsible for validating given urls' SSL certs with their respective CRLs.
    """

    def __init__(self, urls: List[str]) -> None:
        """
        Args:
            urls: list of urls that we should download certificates and CRLs for
        """
        hosts = [get_host(url=url) for url in urls]
        self._certificates = dict.fromkeys(hosts)
        self._crls = {}
        self._lock = RWLock()

        self.update_certificates()

    @property
    def _hosts(self) -> list:
        return list(self._certificates.keys())

    def is_valid(self, url: str) -> bool:
        """ Check if given url SSL certificate is not revoked by its Certificate Revocation Lists.

        Args:
            url: url string

        Returns:
            bool: False if certificate is not revoked

        Raises:
            SSLError: if certificate is revoked
            KeyError: if asked for validity of a host that was not given on initialization
        """
        host = get_host(url=url)

        try:
            with self._lock.read:
                certificate = self._certificates[host]
                if certificate == CERTIFICATE_UNAVAILABLE:
                    # If certificate is not available now it should be considered valid.
                    # We can't get CRLs without having the certificate, but CRLs are not mandatory in any way
                    # and certificate is needed here only to validate it against its CRLs. If we don't have it
                    # now we'll get it on next update, and if there was anything wrong with the SSL, but it
                    # wasn't related to CRLs then requests will catch that while making the call to that host.
                    return True

                crls = self._crls.get(certificate.serial_number, [])
        except KeyError:
            raise KeyError(
                f'{host} CRL validation unavailable, host was not given upon initialization, '
                f'so certificates were not prefetched. Available hosts: {self._hosts}',
            )

        return not is_certificate_revoked(certificate=certificate, crls=crls)

    def update_certificates(self) -> None:
        """ Download new certificates and CRLs for all hosts and update self._certificates
        and self._crls dicts. This method should be called periodically by some
        cron or scheduler to ensure that recent certificates are always fetched.

        Returns: None
        """
        for host in self._hosts:
            try:
                certificate = get_certificate(hostname=host)
            except (ConnectionError, ValueError) as e:
                logger.warning(e)

                with self._lock.write:
                    self._certificates[host] = CERTIFICATE_UNAVAILABLE
                continue

            crls = get_certificate_crls(certificate=certificate)

            with self._lock.write:
                self._certificates[host] = certificate
                self._crls[certificate.serial_number] = crls
