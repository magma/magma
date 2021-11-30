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
import logging
from collections import defaultdict
from time import time
from typing import Dict, List, NamedTuple, Optional

import grpc
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.sentry import SEND_TO_ERROR_MONITORING
from magma.common.service_registry import ServiceRegistry
from magma.magmad.check.network_check.ping import PingCommandResult
from magma.monitord.icmp_state import ICMPMonitoringResponse
from orc8r.protos.common_pb2 import Void
from prometheus_client import Histogram

subscriber_icmp_latency_ms = Histogram(
    'subscriber_icmp_latency_ms',
    'Reported latency for subscriber '
    'in milliseconds',
    ['imsi'],
    buckets=[50, 100, 200, 500, 1000, 2000],
)

PingedTargets = NamedTuple(
    'PingedTargets',
    [
        ('ping_targets', Dict['str', IPAddress]),
        ('ping_addresses', List[IPAddress]),
    ],
)


def _get_addr_from_subscribers(sub_ip) -> str:
    if sub_ip.version == IPAddress.IPV4:
        return str(ipaddress.IPv4Address(sub_ip.address))
    else:
        return str(ipaddress.IPv6Address(sub_ip.address))


class CpeMonitoringModule:
    def __init__(self):
        self._subscriber_state = defaultdict(ICMPMonitoringResponse)
        self.ping_addresses = []
        self.ping_targets = {}

    def set_manually_configured_targets(
            self, configured_ping_targets: Optional[Dict] = None,
    ):
        if configured_ping_targets:
            self.ping_targets = configured_ping_targets.copy()
            for value in self.ping_targets.values():
                ip = _get_addr_from_subscribers(value)
                self.ping_addresses.append(ip)

    async def get_ping_targets(self, service_loop) -> PingedTargets:
        """
        Sends gRPC call to mobilityd to get all subscribers table.

        Returns: List of [Subscriber ID => IP address, APN] entries
        """

        ping_addresses = self.ping_addresses.copy()
        ping_targets = self.ping_targets.copy()
        try:
            mobilityd_chan = ServiceRegistry.get_rpc_channel(
                'mobilityd',
                ServiceRegistry.LOCAL,
            )
            mobilityd_stub = MobilityServiceStub(mobilityd_chan)
            response = await grpc_async_wrapper(
                mobilityd_stub.GetSubscriberIPTable.future(
                    Void(),
                    10,
                ), service_loop,
            )
            for sub in response.entries:
                ip = _get_addr_from_subscribers(sub.ip)
                ping_addresses.append(ip)
                ping_targets[sub.sid.id] = ip
        except grpc.RpcError as err:
            logging.error(
                "GetSubscribers Error for %s! %s",
                err.code(),
                err.details(),
                extra=SEND_TO_ERROR_MONITORING,
            )
        return PingedTargets(ping_targets, ping_addresses)

    def save_ping_response(
        self, sid: str, ip_addr: str,
        ping_resp: PingCommandResult,
    ) -> None:
        reported_time = int(round(time() * 1000))
        self._subscriber_state[sid] = ICMPMonitoringResponse(
            last_reported_time=reported_time,
            latency_ms=ping_resp.stats.rtt_avg,
        )
        subscriber_icmp_latency_ms.labels(sid).observe(ping_resp.stats.rtt_avg)
        logging.info(
            '%s:%s => %sms', sid, ip_addr,
            self._subscriber_state[sid].latency_ms,
        )

    def get_subscriber_state(self) -> Dict[str, ICMPMonitoringResponse]:
        return self._subscriber_state
