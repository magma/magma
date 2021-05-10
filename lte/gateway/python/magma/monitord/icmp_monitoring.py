"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import asyncio
import logging
from typing import Dict, List, Optional

from magma.common.job import Job
from magma.magmad.check.network_check.ping import (
    PingCommandResult,
    PingInterfaceCommandParams,
    ping_interface_async,
)

NUM_PACKETS = 4
DEFAULT_POLLING_INTERVAL = 60
TIMEOUT_SECS = 10
CHECKIN_INTERVAL = 10


class ICMPMonitoring(Job):
    """
    Class that handles main loop to send ICMP ping to valid subscribers.
    """

    def __init__(
        self, monitoring_module, polling_interval: int, service_loop,
        mtr_interface: str,
    ):
        super().__init__(interval=CHECKIN_INTERVAL, loop=service_loop)
        self._MTR_PORT = mtr_interface
        logging.info("Running on interface %s...", self._MTR_PORT)
        # Matching response time output to get latency
        self._polling_interval = max(
            polling_interval,
            DEFAULT_POLLING_INTERVAL,
        )
        self._loop = service_loop
        self._module = monitoring_module

    async def _ping_targets(
        self, hosts: List[str],
        targets: Optional[Dict] = None,
    ):
        """
        Sends a count of ICMP pings to target IP address, returns response.
        Args:
            hosts: List of ip addresses to ping
            targets: List of valid subscribers to ping to

        Returns: (stdout, stderr)
        """
        if targets:
            ping_params = [
                PingInterfaceCommandParams(
                    host, NUM_PACKETS, self._MTR_PORT,
                    TIMEOUT_SECS,
                ) for host in hosts
            ]
            ping_results = await ping_interface_async(ping_params, self._loop)
            ping_results_list = list(ping_results)
            for host, sub, result in zip(hosts, targets, ping_results_list):
                self._save_ping_response(sub, host, result)

    def _save_ping_response(
        self, target_id: str, ip_addr: str,
        ping_resp: PingCommandResult,
    ) -> None:
        """
        Saves ping response to in-memory subscriber dict.
        Args:
            target_id: target ID to ping
            ip_addr: IP Address to ping
            ping_resp: response of ICMP ping command
        """
        if ping_resp.error:
            logging.debug(
                'Failed to ping %s with error: %s',
                target_id, ping_resp.error,
            )
        else:
            self._module.save_ping_response(target_id, ip_addr, ping_resp)

    async def _run(self) -> None:
        targets, addresses = await self._module.get_ping_targets(self._loop)
        if len(targets) > 0:
            await self._ping_targets(addresses, targets)
        else:
            logging.warning('No subscribers/ping targets found')
        await asyncio.sleep(self._polling_interval, self._loop)
