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
import logging
from collections import defaultdict

import os
from time import time
from typing import List, NamedTuple, Optional

import shutil
from magma.magmad.check import subprocess_workflow
from redis.exceptions import ConnectionError

from magma.common.health.service_state_wrapper import ServiceStateWrapper
from magma.common.job import Job
from orc8r.protos.service_status_pb2 import ServiceExitStatus

SystemdServiceParams = NamedTuple('SystemdServiceParams', [('service', str)])


class StateRecoveryJob(Job):
    """
    Class that handles main loop to poll service status to check and restarts
    sctpd as recovery method for services crashing
    """

    def __init__(self, service_state: ServiceStateWrapper,
                 polling_interval: int, services_check: List[str],
                 restart_threshold: int, redis_dump_src: str,
                 snapshots_dir: str, service_loop):
        super().__init__(interval=polling_interval, loop=service_loop)
        self._state_wrapper = service_state
        self._services_check = services_check
        self._threshold = restart_threshold
        self._snapshots_dir = snapshots_dir
        self._redis_dump_src = redis_dump_src
        self._loop = service_loop
        self._polling_interval = polling_interval
        self._services_restarts_map = defaultdict(int)

    def _get_service_status(self, service_name: str) \
            -> Optional[ServiceExitStatus]:
        """
        Args:
            service_name: name of magma service to check status

        Returns: ServiceExitStatus obj for given service

        """
        try:
            service_status = self._state_wrapper.get_service_status(
                service_name)
            return service_status
        except (KeyError, ConnectionError) as err:
            logging.debug('Could not obtain service status: [%s]' % err)
            return None

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
            if result:
                last_num_restarts = self._services_restarts_map[service]
                current_restarts = result.num_fail_exits - last_num_restarts
                if current_restarts > self._threshold:
                    logging.info(
                        'Service %s has failed %s times over last %s seconds, '
                        'restarting sctpd to clean state...', service,
                        current_restarts,
                        self._polling_interval)

                    # Save RDB snapshot
                    os.makedirs(os.path.dirname(self._snapshots_dir),
                                exist_ok=True)
                    shutil.copy("%s/redis_dump.rdb" % self._redis_dump_src,
                                "%s/redis_dump_%s.rdb" % (
                                    self._snapshots_dir, time()))

                    # Restart sctpd
                    await self.restart_service_async('sctpd')

                # Update latest number of restarts for service
                self._services_restarts_map[service] = current_restarts
                logging.debug('Unexpected restarts for service %s: %s',
                              service,
                              result.num_fail_exits)


def _get_service_restart_args(param: SystemdServiceParams) -> List[str]:
    """

    Args:
        param: SystemdServiceParams tuple

    Returns: List of (str) parameters for systemctl service restart

    """
    return ['sudo', 'service', param.service, 'restart']
