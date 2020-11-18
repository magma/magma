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
import ipaddress
import logging
from typing import Dict, List

import grpc
from lte.protos.mobilityd_pb2 import IPAddress, SubscriberIPTable
from magma.common.job import Job
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.service_registry import ServiceRegistry
from magma.magmad.check.network_check.ping import PingInterfaceCommandParams, ping_interface_async
from magma.magmad.check.network_check.ping import PingCommandResult
from magma.monitord.icmp_state import ICMPMonitoringResponse
from orc8r.protos.common_pb2 import Void


NUM_PACKETS = 4
DEFAULT_POLLING_INTERVAL = 60
TIMEOUT_SECS = 10
CHECKIN_INTERVAL = 10

class ICMPMonitoring(Job):
    """
    Class that handles main loop to send ICMP ping to valid subscribers.
    """

    def __init__(self, monitoring_module, polling_interval: int, service_loop,
                 mtr_interface: str):
        super().__init__(interval=CHECKIN_INTERVAL, loop=service_loop)
        self._MTR_PORT = mtr_interface
        # Matching response time output to get latency
        self._polling_interval = max(polling_interval,
                                     DEFAULT_POLLING_INTERVAL)
        self._loop = service_loop
        self._module = monitoring_module

    async def _ping_targets(self, hosts: List[str],
                                targets: {}):
        """
        Sends a count of ICMP pings to target IP address, returns response.
        Args:
            hosts: List of ip addresses to ping
            subscribers: List of valid subscribers to ping to

        Returns: (stdout, stderr)
        """
        ping_params = [
            PingInterfaceCommandParams(host, NUM_PACKETS, self._MTR_PORT,
                                            TIMEOUT_SECS) for host in hosts]
        ping_results = await ping_interface_async(ping_params, self._loop)
        ping_results_list = list(ping_results)
        for host, sub, result in zip(hosts, targets, ping_results_list):
            self._save_ping_response(sub, host, result)

    def _save_ping_response(self, target_id: str, ip_addr: str,
                            ping_resp: PingCommandResult) -> None:
        """
        Saves ping response to in-memory subscriber dict.
        Args:
            sid: subscriber ID
            ping_resp: response of ICMP ping command
        """
        if ping_resp.error:
            logging.debug('Failed to ping %s with error: %s',
                          target_id, ping_resp.error)
        else:
            self._module.save_ping_response(target_id, ip_addr, ping_resp)

    async def _run(self) -> None:
        logging.info("Running on interface %s..." % self._MTR_PORT)
        while True:
            targets, addresses = await self._module.get_ping_targets(self._loop)
            if len(targets) > 0:
                await self._ping_targets(addresses, targets)
                await asyncio.sleep(self._polling_interval, self._loop)
            else:
                logging.warning('No subscribers/ping targets found, retrying...')
                await asyncio.sleep(self._polling_interval, self._loop)
                continue
