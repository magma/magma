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


import logging
import unittest
from collections import defaultdict
import subprocess

from magma.mobilityd.uplink_gw import UplinkGatewayInfo

LOG = logging.getLogger('mobilityd.def_gw.test')
LOG.isEnabledFor(logging.DEBUG)


class DefGwTest(unittest.TestCase):
    """
    Validate default router setting.
    """
    def setUp(self):
        self.gw_store = defaultdict(str)
        self.dhcp_gw_info = UplinkGatewayInfo(self.gw_store)

    def test_gw_ip_for_DHCP(self):
        self.assertEqual(self.dhcp_gw_info.getIP(), None)
        self.assertEqual(self.dhcp_gw_info.getMac(), None)

    def test_gw_ip_for_Ip_pool(self):
        self.dhcp_gw_info.read_default_gw()

        def_gw_cmd = "ip route show |grep default| awk '{print $3}'"
        p = subprocess.Popen([def_gw_cmd], stdout=subprocess.PIPE, shell=True)
        def_ip = p.stdout.read().decode("utf-8").strip()
        self.assertEqual(self.dhcp_gw_info.getIP(), str(def_ip))
        self.assertEqual(self.dhcp_gw_info.getMac(), None)
