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
import os
import subprocess

import fire as fire
from magma.common.health.docker_health_service import DockerHealthChecker
from magma.common.health.health_service import GenericHealthChecker
from termcolor import colored


def is_docker():
    """ Checks if the current script is executed in a docker container """
    path = '/proc/self/cgroup'
    return (
        os.path.exists('/.dockerenv') or
        os.path.isfile(path) and any('docker' in line for line in open(path))
    )


class HealthCLI:
    """
    Command line interface for generic Health-Checking.
    """

    def __init__(self):
        self._health_checker = DockerHealthChecker() \
            if is_docker() \
            else GenericHealthChecker()

    def __call__(self, *args, **kwargs):
        self.status()

    def status(self):
        """
        Global health status \n
        Example:
            `health_cli.py status`
            `health_cli.py`
            `venvsudo health_cli.py` (if running without sufficient permissions)
        """
        print('Health Summary')
        # Check connection to the orchestrator
        # This part is implemented in the checkin_cli.py - we'll just execute it
        print('\nGateway <-> Controller connectivity')
        checkin, error = subprocess.Popen(
            ['checkin_cli.py'],
            stdout=subprocess.PIPE,
        ).communicate()
        print(str(checkin, 'utf-8'))
        print(str(self._health_checker.get_health_summary()))

    def magma_version(self):
        """
        Get the installed magma version
        """
        print(str(self._health_checker.get_magma_version()))

    def kernel_version(self):
        """
        Get kernel version of the VM
        """
        print(str(self._health_checker.get_kernel_version()))

    def internet_status(self, host):
        """
        Checks if it's possible to connect to the specified host \n
        Examples:
            `health_cli.py internet_status --host 8.8.8.8`
            `health_cli.py internet_status --host google.com`
        """
        print(str(self._health_checker.ping_status(host)))

    def services_status(self):
        """
        Get status summary for all the magma services
        """
        print(str(self._health_checker.get_magma_services_summary()))

    def restarts_status(self):
        """
        How many times each services was restarting since the whole system start
        """
        return str(self._health_checker.get_unexpected_restart_summary())

    def error_status(self, service_names):
        """
        How many errors have each service had since the last restart \n
        Examples:
            `health_cli.py error_status --service_names mme,dnsd`
            `health_cli.py error_status --service_names '[pipelined,mme]'`
        :param service_names: list or tuple of service names
        """
        print(
            '\n'.join([
                '{}:\t{}'.format(name, errors) for name, errors in
                 self._health_checker
                .get_error_summary(service_names)
                .items()
            ]),
        )


if __name__ == '__main__':
    health_cli = HealthCLI()
    try:
        fire.Fire(health_cli)
    except Exception as e:
        print(colored('Error:', 'red'), e)
