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
from typing import List
from urllib.parse import urlparse

from magma.configuration_controller.crl_validator.certificate import (
    get_certificate,
    get_certificate_crls,
    is_certificate_revoked,
)
from rwmutex import RWLock


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
                crls = self._crls[certificate.serial_number]
        except KeyError:
            raise KeyError(
                f'{host} certificates unavailable, host was not given upon initialization'
                f'so it was not prefetched. Available certificates: {self._hosts}',
            )

        return not is_certificate_revoked(certificate=certificate, crls=crls)

    def update_certificates(self) -> None:
        """ Download new certificates and CRLs for all hosts and update self._certificates
        and self._crls dicts. This method should be called periodically by some
        cron or scheduler to ensure that recent certificates are always fetched.

        Returns: None
        """
        for host in self._hosts:
            certificate = get_certificate(hostname=host)
            crls = get_certificate_crls(certificate=certificate)

            with self._lock.write:
                self._certificates[host] = certificate
                self._crls[certificate.serial_number] = crls
