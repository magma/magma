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
# pylint: disable=broad-except

import os
import shutil
import logging
import time
from typing import Any
import asyncio

from magma.common.job import Job

from magma.magmad.service_poller import ServicePoller
from magma.magmad.service_manager import ServiceManager


class ServiceHealthWatchdog(Job):
    """
    Periodically collects service grpc timeout stats and restarts the services
    if needed
    """

    # Periodicity for checking services
    CHECK_STATUS_INTERVAL = 15
    # Default number of continuous timeouts before doing a service restart
    DEFAULT_RESTART_TIMEOUT_THRESHOLD = 15

    RDB_DBFILE_FMT = 'redis_dump_%s.rdb'
    DEFAULT_SNAPSHOTS_DIR = '/var/opt/magma/snapshots'
    RDB_DUMP_PATH = '/var/opt/magma/redis_dump.rdb'

    def __init__(self, config: Any,
                 loop: asyncio.AbstractEventLoop,
                 service_poller: ServicePoller,
                 service_manager: ServiceManager):
        super().__init__(
            interval=self.CHECK_STATUS_INTERVAL,
            loop=loop
        )
        self._loop = loop
        self._config = config
        self._services_to_restart = []
        if 'services_to_restart' in config:
            self._services_to_restart = config['services_to_restart']
        self._restart_timeout_threshold = self.DEFAULT_RESTART_TIMEOUT_THRESHOLD
        if 'restart_timeout_threshold' in config:
            self._restart_timeout_threshold = \
                config['restart_timeout_threshold']
        self._enable_state_recovery = False
        self._snapshots_dir = self.DEFAULT_SNAPSHOTS_DIR
        if 'enable_state_recovery' in config:
            self._enable_state_recovery = self._config['enable_state_recovery']
            self._snapshots_dir = self._config.get('snapshots_dir',
                                                   self.DEFAULT_SNAPSHOTS_DIR)
        self._service_poller = service_poller
        self._service_manager = service_manager

    async def _run(self):
        await self._check_service_timeouts()

    async def _check_service_timeouts(self):
        """
        Make RPC calls to 'GetServiceInfo' functions of other services, to
        get current status.
        """
        services_to_restart = []
        timeout_dict = self._service_poller.get_service_timeouts()
        for service, timeouts in timeout_dict.items():
            if service not in self._services_to_restart:
                continue
            if timeouts >= self._restart_timeout_threshold:
                logging.info('Adding service %s to restart list', service)
                services_to_restart.append(service)
                self._service_poller.reset_timeout_counter(service)
        if services_to_restart:
            await asyncio.gather(
                self._service_manager.restart_services(services_to_restart)
            )
            if self._enable_state_recovery:
                await self._save_rdb_snapshot()
                logging.info('Restarting sctpd')

    async def _save_rdb_snapshot(self):
        rdb_file_name = self.RDB_DBFILE_FMT % int(time.time())
        dest_path = '%s/%s' % (self._snapshots_dir, rdb_file_name)
        logging.info('Saving rdb snapshot in: %s', dest_path)
        os.makedirs(os.path.dirname(dest_path), exist_ok=True)
        try:
            shutil.copyfile(self.RDB_DUMP_PATH, dest_path)
        except FileNotFoundError:
            logging.warning('Error saving rdb file from %s',
                            self.RDB_DUMP_PATH)
