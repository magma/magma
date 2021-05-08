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
import glob
import ipaddress
import math
from os import path

from lte.protos.enodebd_pb2_grpc import EnodebdStub
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from magma.common.service_registry import ServiceRegistry
from magma.configuration.mconfig_managers import load_service_mconfig_as_json
from magma.health.entities import (
    AGWHealthSummary,
    CoreDumps,
    RegistrationSuccessRate,
)
from orc8r.protos.common_pb2 import Void


class AGWHealth:

    def get_allocated_ips(self):
        chan = ServiceRegistry.get_rpc_channel(
            'mobilityd',
            ServiceRegistry.LOCAL,
        )
        client = MobilityServiceStub(chan)
        res = []

        list_blocks_resp = client.ListAddedIPv4Blocks(Void())
        for block_msg in list_blocks_resp.ip_block_list:

            list_ips_resp = client.ListAllocatedIPs(block_msg)
            for ip_msg in list_ips_resp.ip_list:
                if ip_msg.version == IPAddress.IPV4:
                    ip = ipaddress.IPv4Address(ip_msg.address)
                elif ip_msg.address == IPAddress.IPV6:
                    ip = ipaddress.IPv6Address(ip_msg.address)
                else:
                    continue
                res.append(ip)
        return res

    def get_subscriber_table(self):
        chan = ServiceRegistry.get_rpc_channel(
            'mobilityd',
            ServiceRegistry.LOCAL,
        )
        client = MobilityServiceStub(chan)

        table = client.GetSubscriberIPTable(Void())
        return table.entries

    def get_registration_success_rate(self, mme_log_file):
        with open(mme_log_file, 'r') as f:
            log = f.read()

        return RegistrationSuccessRate(
            attach_requests=log.count('Attach Request'),
            attach_accepts=log.count('Attach Accept'),
        )

    def get_core_dumps(
        self,
        directory='/var/core',
        start_timestamp=0,
        end_timestamp=math.inf,
    ):
        res = []
        for filename in glob.glob(path.join(directory, 'core-*')):
            # core-1565125801-python3-8042_bundle
            ts = int(filename.split('-')[1])
            if start_timestamp <= ts <= end_timestamp:
                res.append(filename)
        return CoreDumps(core_dump_files=res)

    def gateway_health_status(self):
        config = load_service_mconfig_as_json('mme')

        # eNB status for #eNBs connected
        chan = ServiceRegistry.get_rpc_channel(
            'enodebd', ServiceRegistry.LOCAL,
        )
        client = EnodebdStub(chan)
        status = client.GetStatus(Void())

        mme_log_path = '/var/log/mme.log'
        health_summary = AGWHealthSummary(
            hss_relay_enabled=config.get('hssRelayEnabled', False),
            nb_enbs_connected=status.meta.get('n_enodeb_connected', 0),
            allocated_ips=self.get_allocated_ips(),
            subscriber_table=self.get_subscriber_table(),
            core_dumps=self.get_core_dumps(),
            registration_success_rate=self.get_registration_success_rate(
                mme_log_path,
            ),
        )
        return health_summary
