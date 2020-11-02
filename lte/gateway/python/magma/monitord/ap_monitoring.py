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
import asyncio
import logging
from collections import defaultdict
from datetime import datetime
from typing import Dict, List
from magma.configuration import load_service_config
from lte.protos.mobilityd_pb2 import IPAddress, SubscriberIPTable
from magma.magmad.check.network_check.ping import PingCommandResult
from magma.monitord.icmp_state import ICMPMonitoringResponse
from prometheus_client import Histogram


AP_ICMP_LATENCY_MS = Histogram('ap_icmp_latency_ms',
                                  'Reported latency for APs '
                                  'in milliseconds',
                                  ['ap_name'],
                                  buckets=[50, 100, 200, 500, 1000, 2000])

class ApMonitoring():
  def __init__(self):
    # TODO: Save to redis
    self._ap_state = defaultdict(ICMPMonitoringResponse)


  async def get_ping_targets(self):
      targets = {}
      addresses = []
      try:
        ap_list = load_service_config("monitord")["ping_targets"]
        for ap, data in ap_list.items():
            if "ip" in data:
                ip = IPAddress(version=IPAddress.IPV4, address=str.encode(data["ip"]))
                logging.debug('Adding {}:{}:{} to ping target'.format(ap, ip.version, ip.address))
                targets[ap] = ip
                addresses.append(data["ip"])
        return targets, addresses
      except KeyError:
          logging.error("No ping targets configured")
          return [], []

  def save_ping_response(self, sid: str, ip_addr: str,
                           ping_resp: PingCommandResult) -> None:
      reported_time = datetime.now().timestamp()
      self._ap_state[sid] = ICMPMonitoringResponse(
        last_reported_time=int(reported_time),
        latency_ms=ping_resp.stats.rtt_avg)
      AP_ICMP_LATENCY_MS.labels(sid).observe(ping_resp.stats.rtt_avg)
      logging.info(
        '{}:{} => {}ms'.format(sid, ip_addr,
                               self._ap_state[sid].latency_ms))

  def get_subscriber_state(self) -> Dict[str, ICMPMonitoringResponse]:
      return self._ap_state

