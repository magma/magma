#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.
import asyncio
import ipaddress
import logging
from collections import defaultdict
from datetime import datetime
from typing import Dict, List

import grpc
from lte.protos.mobilityd_pb2 import IPAddress, SubscriberIPTable
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from magma.common.job import Job
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.service_registry import ServiceRegistry
from magma.magmad.check.network_check import ping
from magma.magmad.check.network_check.ping import PingCommandResult
from magma.monitord.icmp_state import ICMPMonitoringResponse
from magma.monitord.metrics import SUBSCRIBER_ICMP_LATENCY_MS
from orc8r.protos.common_pb2 import Void

NUM_PACKETS = 4
DEFAULT_POLLING_INTERVAL = 60
TIMEOUT_SECS = 10
CHECKIN_INTERVAL = 10


def _get_addr_from_subscriber(sub) -> str:
    return str(ipaddress.IPv4Address(
        sub.ip.address) if sub.ip.version == 0 else \
                   ipaddress.IPv6Address(sub.ip.address))


class ICMPMonitoring(Job):
    """
    Class that handles main loop to send ICMP ping to valid subscribers.
    """

    def __init__(self, polling_interval: int, service_loop,
                 mtr_interface: str):
        super().__init__(interval=CHECKIN_INTERVAL, loop=service_loop)
        self._MTR_PORT = mtr_interface
        # Matching response time output to get latency
        self._polling_interval = max(polling_interval,
                                     DEFAULT_POLLING_INTERVAL)
        # TODO: Save to redis
        self._subscriber_state = defaultdict(ICMPMonitoringResponse)
        self._loop = service_loop

    async def _get_subscribers(self) -> List[IPAddress]:
        """
        Sends gRPC call to mobilityd to get all subscribers table.

        Returns: List of [Subscriber ID => IP address, APN] entries
        """
        try:
            mobilityd_chan = ServiceRegistry.get_rpc_channel('mobilityd',
                                                             ServiceRegistry.LOCAL)
            mobilityd_stub = MobilityServiceStub(mobilityd_chan)
            response = await grpc_async_wrapper(
                mobilityd_stub.GetSubscriberIPTable.future(Void(),
                                                           TIMEOUT_SECS),
                self._loop)
            return response.entries
        except grpc.RpcError as err:
            logging.error(
                "GetSubscribers Error for %s! %s", err.code(), err.details())
            return []

    async def _ping_subscribers(self, hosts: List[str],
                                subscribers: SubscriberIPTable):
        """
        Sends a count of ICMP pings to target IP address, returns response.
        Args:
            hosts: List of ip addresses to ping
            subscribers: List of valid subscribers to ping to

        Returns: (stdout, stderr)
        """
        ping_params = [
            ping.PingInterfaceCommandParams(host, NUM_PACKETS, self._MTR_PORT,
                                            TIMEOUT_SECS) for host in hosts]
        ping_results = await ping.ping_interface_async(ping_params, self._loop)
        ping_results_list = list(ping_results)
        for host, sub, result in zip(hosts, subscribers, ping_results_list):
            sid = "IMSI%s" % sub.sid.id
            self._save_ping_response(sid, host, result)

    def _save_ping_response(self, sid: str, ip_addr: str,
                            ping_resp: PingCommandResult) -> None:
        """
        Saves ping response to in-memory subscriber dict.
        Args:
            sid: subscriber ID
            ping_resp: response of ICMP ping command
        """
        if ping_resp.error:
            logging.debug('Failed to ping %s with error: %s',
                          sid, ping_resp.error)
        reported_time = datetime.now().timestamp()
        self._subscriber_state[sid] = ICMPMonitoringResponse(
            last_reported_time=int(reported_time),
            latency_ms=ping_resp.stats.rtt_avg)
        SUBSCRIBER_ICMP_LATENCY_MS.labels(sid).observe(ping_resp.stats.rtt_avg)
        logging.info(
            '{}:{} => {}ms'.format(sid, ip_addr,
                                   self._subscriber_state[sid].latency_ms))

    def get_subscriber_state(self) -> Dict[str, ICMPMonitoringResponse]:
        return self._subscriber_state

    async def _run(self) -> None:
        logging.info("Running on interface %s..." % self._MTR_PORT)
        while True:
            try:
                subscribers = await self._get_subscribers()
                addresses = [_get_addr_from_subscriber(sub) for sub in
                             subscribers]
                await self._ping_subscribers(addresses, subscribers)
                await asyncio.sleep(self._polling_interval, self._loop)
            except AttributeError:
                logging.warning('No subscribers found, retrying...')
                await asyncio.sleep(self._polling_interval, self._loop)
                continue
