#!/usr/bin/env python3

"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import subprocess
import sys

import fire

from magma.health.health_service import AGWHealth


class AGWHealthCLI:
    """ Command line interface for Health-Checking specific to Access Gateway"""
    def __init__(self):
        self._health_checker = AGWHealth()

    def status(self):
        """ Access Gateway Health Status """
        print('Access Gateway health summary')
        print(str(self._health_checker.gateway_health_status()))
        checkin, error = subprocess.Popen(['health_cli.py'],
                                          stdout=subprocess.PIPE).communicate()
        print(str(checkin, 'utf-8'))

    def allocated_ips(self):
        """ List allocated IPs """
        print('\n'.join(self._health_checker.get_allocated_ips()))

    def subscriber_table(self):
        """ Get the subscriber table """
        print(str(self._health_checker.get_subscriber_table()))

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
        if len(sys.argv) == 1:
            fire.Fire(health_cli.status)
        else:
            fire.Fire(health_cli)
    except Exception as e:
        print('Error:', e)
