"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


import asyncio
import unittest

from lte.protos.mobilityd_pb2 import IPAddress, SubscriberIPTable
from magma.monitord.icmp_monitoring import ICMPMonitoring
from magma.subscriberdb.sid import SIDUtils


class ICMPMonitoringTests(unittest.TestCase):
    """
    Test class for the ICMPMonitor class
    """

    def _add_test_subscriber(self, imsi):
        ip = IPAddress(version=IPAddress.IPV4, address=b'127.0.0.1')
        sid = SIDUtils.to_pb(imsi)
        self.subscribers.entries.add(sid=sid, ip=ip, apn='test_apn')

    async def _ping_local(self):
        return await self._monitor._ping_subscribers(["127.0.0.1"],
                                                     self.subscribers.entries)

    def setUp(self):
        """
        Creates and sets up ICMP monitor
        """
        self.loop = asyncio.get_event_loop()
        self.subscribers = SubscriberIPTable()
        self._monitor = ICMPMonitoring(polling_interval=5,
                                       service_loop=self.loop,
                                       mtr_interface="test")

    def test_ping_subscriber_saves_response(self):
        imsi = 'IMSI00000000001'
        self._add_test_subscriber(imsi)
        self.loop.run_until_complete(self._ping_local())
        sub_states = self._monitor.get_subscriber_state()
        self.loop.close()
        self.assertTrue(imsi in sub_states)
