#!/usr/bin/env python3

"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

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
        return str(self._health_checker.gateway_health_status())

    def allocated_ips(self):
        """ List allocated IPs """
        return '\n'.join(self._health_checker.get_allocated_ips())

    def subscriber_table(self):
        """ Get the subscriber table """
        return str(self._health_checker.get_subscriber_table())

    def registration_success_rate(self, log_file):
        """
        Get the registration success rate from a mme log file.
        RegistrationSuccessRate = #AttachAccepts / #AttachRequests

        example: `agw_health_cli.py registration_success_rate /var/log/mme.log`
        """
        return str(self._health_checker.get_registration_success_rate(log_file))



if __name__ == '__main__':
    health_cli = AGWHealthCLI()
    if len(sys.argv) == 1:
        fire.Fire(health_cli.status)
    else:
        fire.Fire(health_cli)
