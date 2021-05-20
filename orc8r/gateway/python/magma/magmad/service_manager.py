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
import json
import logging
import subprocess
from enum import Enum
from typing import List, Tuple

import magma.magmad.events as magmad_events
from magma.magmad.service_poller import ServicePoller


class ServiceState(Enum):
    """
    Enum for process status return code
    """
    Active = 1
    Activating = 2
    Deactivating = 3
    Inactive = 4
    Failed = 5
    Unknown = 6
    Error = 7


class CommandReturn(object):
    """ Return from _run_command. Intended to mimic return from envoy.run() """

    def __init__(self, status_code, std_out):
        self.status_code = status_code
        self.std_out = std_out


class ServiceManager(object):
    """
    Manages a set of systemd services, performing bulk operations across all
    of them.
    """
    _service_control = {}
    _services = []
    _registered_dynamic_services = []

    def __init__(
            self,
            services: List[str],
            init_system,
            service_poller: ServicePoller,
            registered_dynamic_services: List[str] = None,
            dynamic_services: List[str] = None,
    ):
        if registered_dynamic_services is None:
            registered_dynamic_services = []
        if dynamic_services is None:
            dynamic_services = []
        self._services = services
        self._service_poller = service_poller
        self._registered_dynamic_services = registered_dynamic_services

        init_system_spec = self._get_init_system_spec(init_system)
        for s in self._services:
            self._service_control[s] = self.ServiceControl(s, init_system_spec)
        for s in self._registered_dynamic_services:
            self._service_control[s] = self.ServiceControl(s, init_system_spec)
        self._services.extend(self._clean_dynamic_services(dynamic_services))

    def _get_init_system_spec(self, init_system):
        """
        Returns the approriate ServiceControl subclass based on init_system.
        """
        init_systems = {
            'systemd': self.SystemdInitSystem,
            'runit': self.RunitInitSystem,
            'docker': self.DockerInitSystem,
        }

        try:
            return init_systems[init_system]
        except KeyError:
            raise ValueError(
                'Init system {} not found'.format(init_system),
            )

    async def start_services(self):
        await asyncio.gather(
            *[
                self._service_control[s].start_process()
                for s in self._services
            ],
        )

    async def stop_services(self):
        await asyncio.gather(
            *[self._service_control[s].stop_process() for s in self._services],
        )

    async def restart_services(self, services=None):
        """ Restart some or all services """
        if not services:
            services = self._services
        for service in services:
            self._service_poller.process_service_restart(service)
        await asyncio.gather(
            *[self._service_control[s].restart_process() for s in services],
        )
        magmad_events.restarted_services(services)

    async def update_dynamic_services(self, dynamic_services: List[str]):
        """
        Start/Stop dynamic services, after running this the only dynamic
        services left running are the ones passed in (dynamic_services).
        """
        start, stop = self._parse_dynamic_services(dynamic_services)
        self._service_poller.update_dynamic_services(start, stop)
        self._services = [s for s in self._services + start if s not in stop]

        await asyncio.gather(
            *[self._service_control[s].stop_process() for s in stop],
            *[self._service_control[s].start_process() for s in start],
        )

    def _parse_dynamic_services(
            self,
            dynamic_services: List[str],
    ) -> Tuple[List[str], List[str]]:
        """
        Figure out what dynamic services to stop/start

        Return Tuple:
            start [string]: services which need to be started
            stop  [string]: services which need to be stopped
        """
        clean_dynamic_services = self._clean_dynamic_services(dynamic_services)

        start, stop = [], []
        for s in self._registered_dynamic_services:
            if s not in self._services and s in clean_dynamic_services:
                start.append(s)
            elif s in self._services and s not in clean_dynamic_services:
                stop.append(s)
        return start, stop

    def _clean_dynamic_services(
            self, dynamic_services: List[str],
    ) -> List[str]:
        """Filter out any dynamic services that are not registered

        Args:
            dynamic_services (List[str]): list of services specified in mconfig

        Returns:
            List[str]: intersection of dynamic_services and
            registered_dynamic_services
        """
        clean_dynamic_services = []
        for service_name in dynamic_services:
            if service_name not in self._registered_dynamic_services:
                logging.error(
                    "Not enabling %s as it is not listed as a registered dynamic service in magmad.yml",
                    service_name,
                )
            else:
                clean_dynamic_services.append(service_name)
        return clean_dynamic_services

    class ServiceControl(object):
        """
        Control for managing a systemd service. Issues shell commands and
        blocks until completion. Unit files should be copied to
        /etc/systemd/system/ for systemd managed services. Runit service files
        should be copied to /etc/sv/.
        """

        def __init__(self, name, init_system_spec):
            self._init_system_spec = init_system_spec(name)

        def _run_command(self, command):
            try:
                std_out = subprocess.check_output([
                    self._init_system_spec.init_cmd,
                    command,
                    self._init_system_spec.name,
                ])
                status_code = 0
            except subprocess.CalledProcessError as err:
                status_code = err.returncode
                std_out = err.output
            return CommandReturn(status_code, std_out)

        async def _run_command_async(self, command):
            args = [
                self._init_system_spec.init_cmd,
                command,
                self._init_system_spec.name,
            ]
            proc = await asyncio.create_subprocess_exec(*args)
            std_out = await proc.communicate()
            # returncode of proc is only set when proc.communicate() returns
            status_code = proc.returncode
            return CommandReturn(status_code, std_out)

        async def start_process(self):
            """Starts a process by name."""
            ret = await self._run_command_async(
                self._init_system_spec.start_cmd,
            )
            return ret.status_code

        async def stop_process(self):
            """Stops a process by name."""
            ret = await self._run_command_async(
                self._init_system_spec.stop_cmd,
            )
            return ret.status_code

        async def restart_process(self):
            """Restarts a process by name."""
            ret = await self._run_command_async(
                self._init_system_spec.restart_cmd,
            )
            return ret.status_code

        def status(self):
            """Gets process status.

                Executes the status command and parses response to return an
                approriate ServiceState.
            """
            ret = self._run_command(self._init_system_spec.status_cmd)
            return self._init_system_spec.parse_status(ret.std_out)

    class InitSystemSpec(object):
        # Command to interact with the init system (i.e. 'systemctl', 'sv')
        _init_cmd = None

        # Commands used by the init system to control processes.
        # start/stop/restart seem to be consistent among most init systems so
        # some defaults are set here
        _start_cmd = 'start'
        _stop_cmd = 'stop'
        _restart_cmd = 'reload-or-restart'
        _status_cmd = None

        # Dictionary mapping the status response to a ServiceState
        _statuses = None

        def __init__(self, name):
            assert isinstance(name, str), "Process name is not a string"

        @property
        def name(self):
            if self._name is None:
                raise NotImplementedError()

            return self._name

        @property
        def init_cmd(self):
            """
            Returns command to interact with init system
            (i.e. 'systemctl', 'sv')
            """
            if self._init_cmd is None:
                raise NotImplementedError()

            return self._init_cmd

        @property
        def start_cmd(self):
            """Returns command used by init system to start process"""
            if self._start_cmd is None:
                raise NotImplementedError()

            return self._start_cmd

        @property
        def stop_cmd(self):
            """Returns command used by init system to stop process"""
            if self._stop_cmd is None:
                raise NotImplementedError()

            return self._stop_cmd

        @property
        def restart_cmd(self):
            """Returns command used by init system to restart process"""
            if self._restart_cmd is None:
                raise NotImplementedError()

            return self._restart_cmd

        @property
        def status_cmd(self):
            """Returns command used by init system to get process status"""
            if self._status_cmd is None:
                raise NotImplementedError()

            return self._status_cmd

        def parse_status(self, status):
            """Transforms status returned by init system into a ServiceState"""
            raise NotImplementedError()

    class SystemdInitSystem(InitSystemSpec):
        _init_cmd = 'systemctl'
        _statuses = {
            'active': ServiceState.Active,
            'activating': ServiceState.Activating,
            'inactive': ServiceState.Inactive,
            'deactivating': ServiceState.Deactivating,
            'failed': ServiceState.Failed,
            'unknown': ServiceState.Unknown,
        }
        _status_cmd = 'is-active'

        def __init__(self, name):
            super().__init__(name)
            self._name = 'magma@%s' % name

        def parse_status(self, status):
            """Transforms status returned by init system into a ServiceState"""
            statuses = self._statuses
            std_out_formatted = status.strip().decode()
            if std_out_formatted in statuses:
                return statuses[std_out_formatted]
            else:
                return ServiceState.Error

    class RunitInitSystem(InitSystemSpec):
        _init_cmd = 'sv'
        _statuses = {
            'run': ServiceState.Active,
            'down': ServiceState.Inactive,
            'fail': ServiceState.Failed,
        }
        _status_cmd = 'status'

        def __init__(self, name):
            super().__init__(name)
            self._name = name

        def parse_status(self, status):
            """Transforms status returned by init system into a ServiceState"""
            statuses = self._statuses
            std_out_formatted = status.split(':')[0]
            if std_out_formatted in statuses:
                return statuses[std_out_formatted]
            else:
                return ServiceState.Error

    class DockerInitSystem(InitSystemSpec):
        _init_cmd = 'docker'
        _statuses = {
            'created': ServiceState.Inactive,
            'restarting': ServiceState.Activating,
            'running': ServiceState.Active,
            'paused': ServiceState.Inactive,
            'removing': ServiceState.Deactivating,
            'exited': ServiceState.Inactive,
            'dead': ServiceState.Failed,
        }
        _status_cmd = 'inspect'
        _restart_cmd = 'restart'

        def __init__(self, name):
            super().__init__(name)
            self._name = name

        def parse_status(self, status):
            """Transforms status returned by init system into a ServiceState"""
            statuses = self._statuses

            try:
                inspect_data = json.loads(status.decode())
                return statuses[inspect_data[0]['State']['Status']]
            except (json.decoder.JSONDecodeError, IndexError, KeyError):
                return ServiceState.Error
