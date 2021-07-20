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
import ipaddress
import re
import subprocess
from datetime import datetime
from time import sleep

import fire
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from magma.common.service_registry import ServiceRegistry
from orc8r.protos.common_pb2 import Void


class MonitoringCLI(object):
    """
    CLI for demonstrating simple CPE monitoring using ICMP ping
    """

    def __init__(self):
        # Default internal OVS port for CPE agent
        # TODO: Get mtr port name from yml config when service is created
        self.CPE_PORT_NAME = "mtr0"
        # Matching response time output to get latency
        self.matcher = re.compile(
            b"min/avg/max/mdev = (\\d+.\\d+)/(\\d+.\\d+)/(\\d+.\\d+)/("
            b"\\d+.\\d+)",
        )

    def _get_subscribers(self):
        chan = ServiceRegistry.get_rpc_channel(
            'mobilityd',
            ServiceRegistry.LOCAL,
        )
        client = MobilityServiceStub(chan)

        table = client.GetSubscriberIPTable(Void())
        return table.entries

    def _ping_subscriber(self, ip_addr):
        ping = subprocess.Popen(
            ["ping", "-c", "4,", "-I", self.CPE_PORT_NAME, ip_addr],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        return ping.communicate()

    def run(self, polling_interval):
        while True:
            try:
                subscribers = self._get_subscribers()
                # Ping all subscribers on SID => IP subscriber table
                for sub in subscribers:
                    ip_addr = ipaddress.IPv4Address(
                        sub.ip.address,
                    ) if sub.ip.version == 0 else \
                        ipaddress.IPv6Address(sub.ip.address)
                    print("SID => {} IP => {}".format(sub.sid, ip_addr))
                    sttime = datetime.now().strftime('%m/%d/%Y_%H:%M:%S')
                    out, error = self._ping_subscriber(str(ip_addr))
                    if error:
                        print('Error pinging {}'.format(ip_addr))
                        continue
                    avg_resp_time = self.matcher.search(out).groups()[1]
                    print(
                        "[{}] => Got response from {} in: {} ms".format(
                        sttime, ip_addr, avg_resp_time.decode('utf-8'),
                        ),
                    )
                sleep(polling_interval)
            except KeyboardInterrupt:
                break


if __name__ == "__main__":
    monitoring_cli = MonitoringCLI()
    try:
        fire.Fire(monitoring_cli)
    except Exception as e:
        print('Error: %', e)
