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
# Disable pylint warning as ConnectionError overlaps with built-in
# pylint: disable=redefined-builtin

import logging
import os
import shutil
from time import time
from typing import Dict, List, NamedTuple, Optional

from magma.common.health.service_state_wrapper import ServiceStateWrapper
from magma.common.job import Job
from magma.magmad.check import subprocess_workflow
from orc8r.protos.service_status_pb2 import ServiceExitStatus
from redis.exceptions import ConnectionError

SystemdServiceParams = NamedTuple('SystemdServiceParams', [('service', str)])


class StateRecoveryJob(Job):
    """
    Class that handles main loop to poll service status to identify whether
    services are in a crash loop
    If crash loop is detected, restart sctpd to clean Redis state and
    restart all services
    """

    def __init__(
        self, service_state: ServiceStateWrapper,
        polling_interval: int, services_check: List[str],
        restart_threshold: int, redis_dump_src: str,
        snapshots_dir: str, service_loop,
    ):
        super().__init__(interval=polling_interval, loop=service_loop)
        self._state_wrapper = service_state
        self._services_check = services_check
        self._threshold = restart_threshold
        self._snapshots_dir = snapshots_dir
        self._redis_dump_src = redis_dump_src
        self._loop = service_loop
        self._polling_interval = polling_interval
        self._services_restarts_map = self._get_last_service_restarts()

    def _get_service_status(self, service_name: str) \
            -> Optional[ServiceExitStatus]:
        """
        Args:
            service_name: name of magma service to check status

        Returns: ServiceExitStatus obj for given service

        """
        try:
            service_status = self._state_wrapper.get_service_status(
                service_name,
            )
            return service_status
        except (KeyError, ConnectionError) as err:
            logging.debug('Could not obtain service status: [%s]', err)
            return None

    def _get_last_service_restarts(self) -> Dict[str, int]:
        """
        Gets last unexpected restarts for each service from systemd_status
        redis db

        Returns: Dictionary containing each service and the num of unexpected
        service restarts
        """
        last_services_restarts = {}
        for service in self._services_check:
            last_result = self._get_service_status(service)
            if last_result:
                last_services_restarts[service] = last_result.num_fail_exits
            else:
                last_services_restarts[service] = 0
        return last_services_restarts

    async def restart_service_async(self, service: str):
        """
        Execute service restart commands asynchronously.
        """
        await subprocess_workflow.exec_and_parse_subprocesses_async(
            [SystemdServiceParams(service)],
            _get_service_restart_args,
            None,
            self._loop,
        )

    async def _run(self):
        for service in self._services_check:
            result = self._get_service_status(service)
            if not result:
                continue
            last_num_restarts = self._services_restarts_map[service]
            current_restarts = result.num_fail_exits - last_num_restarts
            if current_restarts > self._threshold:
                logging.info(
                    'Service %s has failed %s times over last %s seconds, '
                    'restarting sctpd to clean state...', service,
                    current_restarts,
                    self._polling_interval,
                )

                # Save RDB snapshot
                os.makedirs(self._snapshots_dir, exist_ok=True)
                shutil.copy(
                    "%s/redis_dump.rdb" % self._redis_dump_src,
                    "%s/redis_dump_%s.rdb" % (
                        self._snapshots_dir, time(),
                    ),
                )

                # Restart sctpd
                await self.restart_service_async('sctpd')

            # Update latest number of restarts for services
            self._services_restarts_map[
                service
            ] = result.num_fail_exits
            logging.debug(
                'Unexpected restarts for service %s: %s',
                service, current_restarts,
            )


def _get_service_restart_args(param: SystemdServiceParams) -> List[str]:
    """

    Args:
        param: SystemdServiceParams tuple

    Returns: List of (str) parameters for systemctl service restart

    """
    return ['sudo', 'service', param.service, 'restart']
