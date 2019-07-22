#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import subprocess

import os
import sys
import fire as fire

from magma.common.health.docker_health_service import DockerHealthChecker
from magma.common.health.health_service import GenericHealthChecker


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

    def status(self):
        """
        Global health status

        example: `health_cli.py status` or just `health_cli.py`
        """
        print('Health Summary')
        # Check connection to the orchestrator
        # This part is implemented in the checkin_cli.py - we'll just execute it
        print('\nGateway <-> Controller connectivity')
        checkin, error = subprocess.Popen(['checkin_cli.py'],
                                          stdout=subprocess.PIPE).communicate()
        print(str(checkin, 'utf-8'))
        return str(self._health_checker.get_health_summary())

    def magma_version(self):
        """
        Get the installed magma version
        """
        return str(self._health_checker.get_magma_version())

    def kernel_version(self):
        """
        Get kernel version of the VM
        """
        return str(self._health_checker.get_kernel_version())

    def internet_status(self, host):
        """
        Checks if it's possible to connect to the specified host

        examples:
            `health_cli.py internet_status --host 8.8.8.8`
            `health_cli.py internet_status --host google.com`
        """
        return str(self._health_checker.ping_status(host))

    def services_status(self):
        """
        Get status summary for all the magma services
        """
        return str(self._health_checker.get_magma_services_summary())

    def restarts_status(self):
        """
        How many times each services was restarting since the whole system start
        """
        return str(self._health_checker.get_unexpected_restart_summary())

    def error_status(self, service_names):
        """
        How many errors have each service had since the last restart

        examples:
            `health_cli.py error_status --service_names mme,dnsd`
            `health_cli.py error_status --service_names '[pipelined,mme]'`
        """
        return '\n'.join(['{}:\t{}'.format(name, errors) for name, errors in
                          self._health_checker
                              .get_error_summary(service_names)
                              .items()
                          ])

if __name__ == '__main__':
    health_cli = HealthCLI()

    if len(sys.argv) == 1:
        fire.Fire(health_cli.status)
    else:
        fire.Fire(health_cli)
