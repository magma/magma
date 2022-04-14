#!/usr/bin/env python3

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
import os
import subprocess
from datetime import datetime

import apt
from dateutil import tz
from magma.common.health.entities import (
    ActiveState,
    Errors,
    HealthStatus,
    HealthSummary,
    RestartFrequency,
    ServiceHealth,
    Version,
)
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.configuration.mconfig_managers import load_service_mconfig_as_json
from magma.magmad.metrics import UNEXPECTED_SERVICE_RESTARTS
from magma.magmad.service_poller import ServicePoller
from orc8r.protos import common_pb2, magmad_pb2
from orc8r.protos.magmad_pb2_grpc import MagmadStub
from orc8r.protos.mconfig import mconfigs_pb2
from pystemd.systemd1 import Unit


class GenericHealthChecker:

    def ping(self, host, num_packets=4):
        chan = ServiceRegistry.get_rpc_channel('magmad', ServiceRegistry.LOCAL)
        client = MagmadStub(chan)

        response = client.RunNetworkTests(
            magmad_pb2.NetworkTestRequest(
                pings=[
                    magmad_pb2.PingParams(
                        host_or_ip=host,
                        num_packets=num_packets,
                    ),
                ],
            ),
        )
        return response.pings

    def ping_status(self, host):
        pings = self.ping(host=host, num_packets=4)[0]
        if pings.error:
            return HealthStatus.DOWN
        if pings.avg_response_ms:
            return HealthStatus.UP
        return HealthStatus.UNKNOWN

    def get_error_summary(self, service_names):
        """Get the list of services with the error count.

        Args:
            service_names: List of service names.

        Returns:
            A dictionary with service name as a key and the Errors object
            as a value.

        Raises:
            PermissionError: User has no permision to exectue the command
        """
        configs = {
            service_name: load_service_mconfig_as_json(service_name)
            for service_name in service_names
        }
        res = {
            service_name: Errors(
                log_level=configs[service_name].get('logLevel', 'INFO'),
                error_count=0,
            )
            for service_name in service_names
        }

        syslog_path = '/var/log/syslog'
        if not os.access(syslog_path, os.R_OK):
            raise PermissionError(
                'syslog is not readable. '
                'Try `sudo chmod a+r {}`. '
                'Or execute the command with sudo '
                'permissions: `venvsudo`'.format(syslog_path),
            )
        with open(syslog_path, 'r', encoding='utf-8') as f:
            for line in f:
                for service_name in service_names:
                    if service_name not in line:
                        continue
                    # Reset the counter for restart/start
                    if 'Starting {}...'.format(service_name) in line:
                        res[service_name].error_count = 0
                    elif 'ERROR' in line:
                        res[service_name].error_count += 1
        return res

    def get_magma_services_summary(self):
        """ Get health for all the running services """
        services_health_summary = []

        # DBus objects: https://www.freedesktop.org/wiki/Software/systemd/dbus/
        chan = ServiceRegistry.get_rpc_channel('magmad', ServiceRegistry.LOCAL)
        client = MagmadStub(chan)

        configs = client.GetConfigs(common_pb2.Void())

        service_names = [str(name) for name in configs.configs_by_key]
        services_errors = self.get_error_summary(service_names=service_names)

        for service_name in service_names:
            unit = Unit(
                'magma@{}.service'.format(service_name),
                _autoload=True,
            )
            active_state = ActiveState.dbus2state[unit.Unit.ActiveState]
            sub_state = str(unit.Unit.SubState, 'utf-8')
            if active_state == ActiveState.ACTIVE:
                pid = unit.Service.MainPID
                process = subprocess.Popen(
                    'ps -o etime= -p {}'.format(pid).split(),
                    stdout=subprocess.PIPE,
                )

                time_running, error = process.communicate()
                if error:
                    raise ValueError(
                        'Cannot get time running for the service '
                        '{} `ps -o etime= -p {}`'
                        .format(service_name, pid),
                    )
            else:
                time_running = b'00'

            services_health_summary.append(
                ServiceHealth(
                    service_name=service_name,
                    active_state=active_state, sub_state=sub_state,
                    time_running=str(time_running, 'utf-8').strip(),
                    errors=services_errors[service_name],
                ),
            )
        return services_health_summary

    def get_unexpected_restart_summary(self):
        service = MagmaService('magmad', mconfigs_pb2.MagmaD())
        service_poller = ServicePoller(service.loop, service.config)
        service_poller.start()

        asyncio.set_event_loop(service.loop)

        # noinspection PyProtectedMember
        # pylint: disable=protected-access
        async def fetch_info():
            restart_frequencies = {}
            await service_poller._get_service_info()
            for service_name in service_poller.service_info.keys():
                restarts = int(
                    UNEXPECTED_SERVICE_RESTARTS
                    .labels(service_name=service_name)
                    ._value.get(),
                )
                restart_frequencies[service_name] = RestartFrequency(
                    count=restarts,
                    time_interval='',
                )

            return restart_frequencies

        return service.loop.run_until_complete(fetch_info())

    def get_kernel_version(self):
        info, error = subprocess.Popen(
            'uname -a'.split(),
            stdout=subprocess.PIPE,
        ).communicate()

        if error:
            raise ValueError('Cannot get the kernel version')
        return str(info, 'utf-8')

    def get_magma_version(self):
        cache = apt.Cache()

        # Return the python version if magma is not there
        if 'magma' not in cache:
            return Version(
                version_code=cache['python3'].versions[0],
                last_update_time='-',
            )

        pkg = str(cache['magma'].versions[0])
        version = pkg.split('-')[0].split('=')[-1]
        timestamp = int(pkg.split('-')[1])

        return Version(
            version_code=version,
            last_update_time=datetime.utcfromtimestamp(timestamp)
            .replace(tzinfo=tz.tzutc())
            .astimezone(tz=tz.tzlocal())
            .strftime('%Y-%m-%d %H:%M:%S'),
        )

    def get_health_summary(self):

        return HealthSummary(
            version=self.get_magma_version(),
            platform=self.get_kernel_version(),
            services_health=self.get_magma_services_summary(),
            internet_health=self.ping_status(host='8.8.8.8'),
            dns_health=self.ping_status(host='google.com'),
            unexpected_restarts=self.get_unexpected_restart_summary(),
        )
