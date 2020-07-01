"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


import logging
import os
import subprocess
import sys
import threading
import unittest

from ipaddress import ip_network
from magma.mobilityd.ip_address_man import IPAddressManager, MappingNotFoundError
from unittest import mock
from magma.common.redis.mocks.mock_redis import MockRedis
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
    @mock.patch("redis.Redis", MockRedis)
    def setUp(self):
        self._br = "t_up_br0"

        subprocess.check_call(["redis-cli", "flushall"])

        setup_dhcp_server = SCRIPT_PATH + "scripts/setup-test-dhcp-srv.sh"
        subprocess.check_call([setup_dhcp_server, "t0"])

        setup_uplink_br = [SCRIPT_PATH + "scripts/setup-uplink-br.sh",
                           self._br,
                           "t0uplink_p0",
                           "8A:00:00:00:00:01"]
        subprocess.check_call(setup_uplink_br)

        config = {
            'dhcp_iface': 't_dhcp0',
            'retry_limit': 50,
            'allocator_type': 'dhcp',
            'persist_to_redis': False,
        }
        self._dhcp_allocator = IPAddressManager(recycling_interval=2,
                                                config=config)
        print("dhcp allocator created")

    def tearDown(self):
        self._dhcp_allocator.ip_allocator.stop_dhcp_sniffer()
        BridgeTools.destroy_bridge(self._br)

    @unittest.skipIf(os.getuid(), reason="needs root user")
    def test_ip_alloc(self):
        sid1 = "IMSI02917"
        ip1 = self._dhcp_allocator.alloc_ip_address(sid1)

        threading.Event().wait(2)
        self._dhcp_allocator.release_ip_address(sid1, ip1)
        threading.Event().wait(2)
        ip1_1 = self._dhcp_allocator.alloc_ip_address(sid1)
        threading.Event().wait(2)
        self.assertEqual(str(ip1), str(ip1_1))

        self._dhcp_allocator.release_ip_address(sid1, ip1_1)
        threading.Event().wait(5)
        self.assertEqual(self._dhcp_allocator.list_added_ip_blocks(), [])

        ip1 = self._dhcp_allocator.alloc_ip_address("IMSI02918")
        self.assertEqual(str(ip1), "192.168.128.146")
        self.assertEqual(self._dhcp_allocator.list_added_ip_blocks(),
                         [ip_network('192.168.128.0/24')])

        ip2 = self._dhcp_allocator.alloc_ip_address("IMSI029192")
        self.assertNotEqual(ip1, ip2)

        ip3 = self._dhcp_allocator.alloc_ip_address("IMSI0432")
        self.assertNotEqual(ip1, ip3)
        self.assertNotEqual(ip2, ip3)
        # release unallocated IP of SID
        self._dhcp_allocator.ip_allocator.release_ip("IMSI033",
                                                     ip3,
                                                     ip_network("1.1.1.0/24"))
        self.assertEqual(self._dhcp_allocator.list_added_ip_blocks(),
                         [ip_network('192.168.128.0/24')])

        sid4 = "IMSI54321"
        ip4 = self._dhcp_allocator.alloc_ip_address(sid4)
        threading.Event().wait(1)
        self._dhcp_allocator.release_ip_address(sid4, ip4)
        self.assertEqual(self._dhcp_allocator.list_added_ip_blocks(),
                         [ip_network('192.168.128.0/24')])

        # wait for DHCP release
        threading.Event().wait(20)
        ip4_2 = self._dhcp_allocator.alloc_ip_address(sid4)
        self.assertEqual(ip4, ip4_2)

        try:
            self._dhcp_allocator.release_ip_address(sid1, ip1)
            self.assertEqual("should not", "reach here")
        except MappingNotFoundError:
            pass
