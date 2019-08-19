#!/usr/bin/env python3

"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import textwrap


class RegistrationSuccessRate:
    def __init__(self, attach_requests, attach_accepts):
        self.attach_requests = attach_requests
        self.attach_accepts = attach_accepts

    @property
    def rate(self):
        if self.attach_requests == 0:
            return float('Inf')
        return 100. * self.attach_accepts / self.attach_requests

    def __str__(self):
        rate = str(self.rate) + '%' if self.attach_requests != 0 \
            else 'No Attach requests'

        return '{} ({} requests, {} accepted)'.format(
            rate,
            self.attach_requests,
            self.attach_accepts,
        )


class AGWHealthSummary:
    def __init__(self, relay_enabled, nb_enbs_connected,
                 allocated_ips, subscriber_table,
                 registration_success_rate):
        self.relay_enabled = relay_enabled
        self.nb_enbs_connected = nb_enbs_connected
        self.allocated_ips = allocated_ips
        self.subscriber_table = subscriber_table
        self.registration_success_rate = registration_success_rate

    def __str__(self):
        return textwrap.dedent("""
        {}
        #eNBs connected: {}\t\t(run `enodebd_cli.py get_all_status` for more details)
        #IPs allocated: {}\t\t(run `mobilityd_cli.py list_allocated_ips` for more details)
        #UEs connected: {}\t\t(run `mobilityd_cli.py get_subscriber_table` for more details)
        Registration success rate: {}
        """).format(
            'Using Feg' if self.relay_enabled else 'Using subscriberdb',
            self.nb_enbs_connected,
            len(self.allocated_ips),
            len(self.subscriber_table),
            self.registration_success_rate,
        )
