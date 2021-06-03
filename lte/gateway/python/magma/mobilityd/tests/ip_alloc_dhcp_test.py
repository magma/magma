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
import os
import subprocess
import sys
import threading
import unittest.mock
from ipaddress import ip_network

from magma.common.redis.client import get_default_client
from magma.mobilityd.ip_address_man import (
    IPAddressManager,
    MappingNotFoundError,
)
from magma.mobilityd.ip_allocator_dhcp import IPAllocatorDHCP
from magma.mobilityd.ip_descriptor import IPDesc, IPState, IPType
from magma.mobilityd.ipv6_allocator_pool import IPv6AllocatorPool
from magma.mobilityd.mac import create_mac_from_sid
from magma.mobilityd.mobility_store import MobilityStore
from magma.pipelined.bridge_util import BridgeTools

LOG = logging.getLogger('mobilityd.dhcp.test')
LOG.isEnabledFor(logging.DEBUG)

"""
This test needs to run on vagrant box due to dependency on
root access and redis.
command to run this specific test.
.  /home/vagrant/build/python/bin/activate; \
sudo /home/vagrant/build/python/bin/nosetests \
--with-coverage --cover-erase --cover-branches \
--cover-package=magma \
-s python/magma/mobilityd/tests/ip_alloc_dhcp_test.py

"""
logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)
SCRIPT_PATH = "/home/vagrant/magma/lte/gateway/python/magma/mobilityd/"


class DhcpIPAllocEndToEndTest(unittest.TestCase):
    def setUp(self):
        self._br = "t_up_br0"

        setup_dhcp_server = SCRIPT_PATH + "scripts/setup-test-dhcp-srv.sh"
        subprocess.check_call([setup_dhcp_server, "t0"])

        setup_uplink_br = [
            SCRIPT_PATH + "scripts/setup-uplink-br.sh",
            self._br,
            "t0uplink_p0",
            "t0_dhcp1",
        ]
        subprocess.check_call(setup_uplink_br)

        store = MobilityStore(get_default_client(), False, 3980)
        ipv4_allocator = IPAllocatorDHCP(
            store, iface='t0uplink_p0',
            retry_limit=50,
        )
        ipv6_allocator = IPv6AllocatorPool(
            store,
            session_prefix_alloc_mode='RANDOM',
        )
        self._dhcp_allocator = IPAddressManager(
            ipv4_allocator,
            ipv6_allocator,
            store,
            recycling_interval=2,
        )

    def tearDown(self):
        self._dhcp_allocator.ip_allocator.stop_dhcp_sniffer()
        BridgeTools.destroy_bridge(self._br)

    @unittest.skipIf(os.getuid(), reason="needs root user")
    def test_ip_alloc(self):
        sid1 = "IMSI02917"
        ip1, _ = self._dhcp_allocator.alloc_ip_address(sid1)
        threading.Event().wait(2)
        dhcp_gw_info = self._dhcp_allocator._store.dhcp_gw_info
        dhcp_store = self._dhcp_allocator._store.dhcp_store

        self.assertEqual(str(dhcp_gw_info.get_gw_ip()), "192.168.128.211")
        self._dhcp_allocator.release_ip_address(sid1, ip1)

        # wait for DHCP release
        threading.Event().wait(7)
        mac1 = create_mac_from_sid(sid1)
        dhcp_state1 = dhcp_store.get(mac1.as_redis_key(None))

        self.assertEqual(dhcp_state1, None)

        ip1_1, _ = self._dhcp_allocator.alloc_ip_address(sid1)
        threading.Event().wait(2)
        self.assertEqual(str(ip1), str(ip1_1))

        self._dhcp_allocator.release_ip_address(sid1, ip1_1)
        threading.Event().wait(5)
        self.assertEqual(self._dhcp_allocator.list_added_ip_blocks(), [])

        ip1, _ = self._dhcp_allocator.alloc_ip_address("IMSI02918")
        self.assertEqual(str(ip1), "192.168.128.146")
        self.assertEqual(
            self._dhcp_allocator.list_added_ip_blocks(),
            [ip_network('192.168.128.0/24')],
        )

        ip2, _ = self._dhcp_allocator.alloc_ip_address("IMSI029192")
        self.assertNotEqual(ip1, ip2)

        ip3, _ = self._dhcp_allocator.alloc_ip_address("IMSI0432")
        self.assertNotEqual(ip1, ip3)
        self.assertNotEqual(ip2, ip3)
        # release unallocated IP of SID
        ip_unallocated = IPDesc(
            ip=ip3, state=IPState.ALLOCATED,
            sid="IMSI033",
            ip_block=ip_network("1.1.1.0/24"),
            ip_type=IPType.DHCP,
        )
        self._dhcp_allocator.ip_allocator.release_ip(ip_unallocated)
        self.assertEqual(
            self._dhcp_allocator.list_added_ip_blocks(),
            [ip_network('192.168.128.0/24')],
        )

        sid4 = "IMSI54321"
        ip4, _ = self._dhcp_allocator.alloc_ip_address(sid4)
        threading.Event().wait(1)
        self._dhcp_allocator.release_ip_address(sid4, ip4)
        self.assertEqual(
            self._dhcp_allocator.list_added_ip_blocks(),
            [ip_network('192.168.128.0/24')],
        )

        # wait for DHCP release
        threading.Event().wait(7)
        mac4 = create_mac_from_sid(sid4)
        dhcp_state = dhcp_store.get(mac4.as_redis_key(None))

        self.assertEqual(dhcp_state, None)
        ip4_2, _ = self._dhcp_allocator.alloc_ip_address(sid4)
        self.assertEqual(ip4, ip4_2)

        try:
            self._dhcp_allocator.release_ip_address(sid1, ip1)
            self.assertEqual("should not", "reach here")
        except MappingNotFoundError:
            pass
