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
import math
import subprocess

import fire
from magma.health.health_service import AGWHealth
from termcolor import colored


class AGWHealthCLI:
    """ Command line interface for Health-Checking specific to Access Gateway"""

    def __init__(self):
        self._health_checker = AGWHealth()

    def __call__(self, *args, **kwargs):
        self.status()

    def status(self):
        """ Access Gateway Health Status """
        print('Access Gateway health summary')
        print(str(self._health_checker.gateway_health_status()))
        checkin, error = subprocess.Popen(
            ['health_cli.py'],
            stdout=subprocess.PIPE,
        ).communicate()
        print(str(checkin, 'utf-8'))

    def allocated_ips(self):
        """ List allocated IPs """
        print('\n'.join(self._health_checker.get_allocated_ips()))

    def subscriber_table(self):
        """ Get the subscriber table """
        print(str(self._health_checker.get_subscriber_table()))

    def core_dumps(
        self,
        directory='/var/core',
        start_timestamp=0,
        end_timestamp=math.inf,
    ):
        """
        Get number of core dumps created during the specified time range \n
        The core dump file is expected to have the format:
        core-<TIMESTAMP>-*_bundle (like core-1566010183-python3-31982_bundle)
        Example:
            `agw_health_cli.py core_dumps
                --directory /tmp/
                --start_timestamp 1565125801
                --end_timestamp 1565447811`
        :param directory: directory where to look for the core dumps
        :param start_timestamp: timestamp integer from where to start counting
        :param end_timestamp: timestamp integer from where to finish counting
        """
        print(
            str(
                self._health_checker.get_core_dumps(
                    directory,
                    start_timestamp,
                    end_timestamp,
                ),
            ),
        )

    def registration_success_rate(self, log_file):
        """
        Get the registration success rate from a mme log file. \n
        RegistrationSuccessRate = #AttachAccepts / #AttachRequests
        Example:
            `agw_health_cli.py registration_success_rate /var/log/mme.log`
        :param log_file: path to the mme log file
        """
        print(str(self._health_checker.get_registration_success_rate(log_file)))


if __name__ == '__main__':
    health_cli = AGWHealthCLI()
    try:
        fire.Fire(health_cli)
    except Exception as e:
        print(colored('Error:', 'red'), e)
