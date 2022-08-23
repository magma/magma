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

from __future__ import annotations

import logging
from dataclasses import asdict, dataclass
from typing import Optional

import requests
from magma.fluentd_client.config import Config, get_config

logging.basicConfig(
    level=logging.DEBUG,
    datefmt='%Y-%m-%d %H:%M:%S',
    format='%(asctime)s %(levelname)-8s %(message)s',
)
logger = logging.getLogger("fluentd_client.client")


class FluentdClientException(Exception):
    """Generic Fluentd Client Exception"""
    pass  # noqa:WPS604


@dataclass
class DPLog(object):
    """
    Class representation of DPLog. An abstraction for messages sent to and from DP.
    """
    event_timestamp: int
    log_from: str
    log_to: str
    log_name: str
    log_message: str
    cbsd_serial_number: str
    network_id: str
    fcc_id: str
    response_code: Optional[str] = None


class FluentdClient(object):
    """
    Client class for Fluentd communication.
    """

    def __init__(self, config: Optional[Config] = None):
        self.config = config or get_config()

    def send_dp_log(self, log: DPLog):
        """
        Send DP log to Fluentd

        Args:
            log (DPLog): DP Log

        Raises:
            FluentdClientException: Generic Fluentd Client Exception
        """
        logger.debug(f"Sending {log.log_name=} to Fluentd")
        try:
            resp = requests.post(
                url=self.config.FLUENTD_URL,
                json=asdict(log),
                verify=self.config.FLUENTD_TLS_ENABLED,
                cert=(self.config.FLUENTD_CERT_PATH, self.config.FLUENTD_CERT_PATH),
            )
        except (requests.HTTPError, requests.RequestException) as err:
            msg = f"Failed to log {log.log_name} response. {err}"
            logging.error(msg)
            raise FluentdClientException(msg)
        logger.debug(f"Sent {log.log_name=} to Fluentd. Response code = {resp.status_code}")
