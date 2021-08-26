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
import unittest

from lte.protos.mobilityd_pb2 import IPAddress, SubscriberIPTable
from magma.monitord.cpe_monitoring import CpeMonitoringModule
from magma.monitord.icmp_monitoring import ICMPMonitoring

LOCALHOST = '127.0.0.1'


class ICMPMonitoringTests(unittest.TestCase):
    """
    Test class for the ICMPMonitor class
    """

    def _add_test_subscriber(self, imsi):
        ip = IPAddress(version=IPAddress.IPV4, address=b'127.0.0.1')
        self.subscribers[imsi] = ip

    async def _ping_local(self):
        return await self._monitor._ping_targets(
            [LOCALHOST],
            self.subscribers,
        )

    def setUp(self):
        """
        Creates and sets up ICMP monitor
        """
        self.loop = asyncio.get_event_loop()
        self.subscribers = {}
        self.obj = CpeMonitoringModule()
        self._monitor = ICMPMonitoring(
            self.obj, polling_interval=5,
            service_loop=self.loop,
            mtr_interface=LOCALHOST,
        )

    def test_ping_subscriber_saves_response(self):
        imsi = 'IMSI00000000001'
        self._add_test_subscriber(imsi)
        self.loop.run_until_complete(self._ping_local())
        sub_states = self.obj.get_subscriber_state()
        self.loop.close()
        self.assertTrue(imsi in sub_states)
