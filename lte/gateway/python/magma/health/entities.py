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
import textwrap
from datetime import datetime

from dateutil import tz


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


class CoreDumps:
    def __init__(self, core_dump_files):
        self.core_dump_files = core_dump_files

    @property
    def earliest(self):
        timestamps = [int(f.split('-')[1]) for f in self.core_dump_files]
        if not timestamps:
            return '-'
        return datetime.utcfromtimestamp(min(timestamps))\
            .replace(tzinfo=tz.tzutc())\
            .astimezone(tz=tz.tzlocal())\
            .strftime('%Y-%m-%d %H:%M:%S')

    @property
    def latest(self):
        timestamps = [int(f.split('-')[1]) for f in self.core_dump_files]
        if not timestamps:
            return None
        return datetime.utcfromtimestamp(max(timestamps))\
            .replace(tzinfo=tz.tzutc())\
            .astimezone(tz=tz.tzlocal())\
            .strftime('%Y-%m-%d %H:%M:%S')

    def __len__(self):
        return len(self.core_dump_files)

    def __str__(self):
        return '#Core dumps:        {}      from: {}        to: {}'.format(
            len(self.core_dump_files), self.earliest, self.latest,
        )


class AGWHealthSummary:
    def __init__(
        self, hss_relay_enabled, nb_enbs_connected,
        allocated_ips, subscriber_table, core_dumps,
        registration_success_rate,
    ):
        self.hss_relay_enabled = hss_relay_enabled
        self.nb_enbs_connected = nb_enbs_connected
        self.allocated_ips = allocated_ips
        self.subscriber_table = subscriber_table
        self.core_dumps = core_dumps
        self.registration_success_rate = registration_success_rate

    def __str__(self):
        return textwrap.dedent("""
        {}
        #eNBs connected:    {} \t (run `enodebd_cli.py get_all_status` for more details)
        #IPs allocated:     {} \t (run `mobility_cli.py list_allocated_ips` for more details)
        #UEs connected:     {} \t (run `mobility_cli.py get_subscriber_table` for more details)
        #Core dumps:        {} \t (run `ls /var/core` to see core dumps)
        Earliest core-dump: {}, Latest core-dump: {}
        Registration success rate: {}
        """).format(
            'Using Feg' if self.hss_relay_enabled else 'Using SubscriberDB',
            self.nb_enbs_connected,
            len(self.allocated_ips),
            len(self.subscriber_table),
            len(self.core_dumps),
            self.core_dumps.earliest, self.core_dumps.latest,
            self.registration_success_rate,
        )
