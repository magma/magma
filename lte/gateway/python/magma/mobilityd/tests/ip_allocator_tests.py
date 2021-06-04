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

from __future__ import (
    absolute_import,
    division,
    print_function,
    unicode_literals,
)

import ipaddress
import time
import unittest

from magma.common.redis.client import get_default_client
from magma.mobilityd.ip_address_man import (
    IPAddressManager,
    IPNotInUseError,
    MappingNotFoundError,
)
from magma.mobilityd.ip_allocator_pool import (
    IpAllocatorPool,
    IPBlockNotFoundError,
    NoAvailableIPError,
)
from magma.mobilityd.ipv6_allocator_pool import IPv6AllocatorPool
from magma.mobilityd.mobility_store import MobilityStore

# If the preallocated IP addresses in ip_address_man.py code
# changes, then the test cases must be updated


class IPAllocatorTests(unittest.TestCase):
    """
    Test class for the Mobilityd IP Allocator
    """

    RECYCLING_INTERVAL_SECONDS = 1

    def _new_ip_allocator(self, recycling_interval):
        """
        Creates and sets up an IPAllocator with the given recycling interval.
        """
        store = MobilityStore(get_default_client(), False, 3980)
        store.dhcp_gw_info.read_default_gw()
        ip_allocator = IpAllocatorPool(store)
        ipv6_allocator = IPv6AllocatorPool(
            store,
            session_prefix_alloc_mode='RANDOM',
        )
        self._allocator = IPAddressManager(
            ip_allocator,
            ipv6_allocator,
            store,
            recycling_interval,
        )
        self._allocator.add_ip_block(self._block)

    def setUp(self):
        #  need to allocate at least 4 bits, as 13 addresses
        #  are either preallocated or not valid for hosts
        self._block = ipaddress.ip_network('192.168.0.0/28')
        #  192.168.0.0 is not valid host IP
        #  192.168.0.1 to 192.168.0.11 are preallocated
        self._ip0 = ipaddress.ip_address('192.168.0.12')
        self._ip1 = ipaddress.ip_address('192.168.0.13')
        self._ip2 = ipaddress.ip_address('192.168.0.14')
        self._new_ip_allocator(self.RECYCLING_INTERVAL_SECONDS)

    def test_list_added_ip_blocks(self):
        """ test list assigned IP blocks """
        ip_block_list = self._allocator.list_added_ip_blocks()
        self.assertEqual(ip_block_list, [self._block])

    def test_list_empty_ip_block(self):
        """ test list empty ip block """
        ip_list = self._allocator.list_allocated_ips(self._block)
        self.assertEqual(len(ip_list), 0)

    def test_list_unknown_ip_block(self):
        """ test list unknown ip block """
        block = ipaddress.ip_network('10.0.0.0/28')
        with self.assertRaises(IPBlockNotFoundError):
            self._allocator.list_allocated_ips(block)

    def test_alloc_ip_address(self):
        """ test alloc_ip_address """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        self.assertTrue(ip0 in [self._ip0, self._ip1, self._ip2])
        self.assertTrue(ip0 in self._allocator.list_allocated_ips(self._block))
        self.assertEqual(self._allocator.get_sid_ip_table(), [('SID0', ip0)])

        ip1, _ = self._allocator.alloc_ip_address('SID1')
        self.assertTrue(ip1 in [self._ip0, self._ip1, self._ip2])
        self.assertNotEqual(ip1, ip0)
        self.assertEqual(
            {ip0, ip1},
            set(self._allocator.list_allocated_ips(self._block)),
        )
        self.assertEqual(
            set(self._allocator.get_sid_ip_table()),
            {('SID0', ip0), ('SID1', ip1)},
        )

        ip2, _ = self._allocator.alloc_ip_address('SID2')
        self.assertTrue(ip2 in [self._ip0, self._ip1, self._ip2])
        self.assertNotEqual(ip2, ip0)
        self.assertNotEqual(ip2, ip1)
        self.assertEqual(
            {ip0, ip1, ip2},
            set(self._allocator.list_allocated_ips(self._block)),
        )
        self.assertEqual(
            set(self._allocator.get_sid_ip_table()),
            {('SID0', ip0), ('SID1', ip1), ('SID2', ip2)},
        )

        # allocate from empty free set
        with self.assertRaises(NoAvailableIPError):
            self._allocator.alloc_ip_address('SID3')

    def test_release_ip_address(self):
        """ test release_ip_address """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')
        ip2, _ = self._allocator.alloc_ip_address('SID2')

        # release ip
        self._allocator.release_ip_address('SID0', ip0)
        self.assertFalse(
            ip0 in self._allocator.list_allocated_ips(self._block),
        )

        # check not recyled
        self.assertEqual(
            set(self._allocator.get_sid_ip_table()),
            {('SID0', ip0), ('SID1', ip1), ('SID2', ip2)},
        )
        with self.assertRaises(NoAvailableIPError):
            self._allocator.alloc_ip_address('SID3')

        # double release
        with self.assertRaises(IPNotInUseError):
            self._allocator.release_ip_address('SID0', ip0)

        # ip does not exist
        with self.assertRaises(MappingNotFoundError):
            non_existing_ip = ipaddress.ip_address('192.168.1.16')
            self._allocator.release_ip_address('SID0', non_existing_ip)

    def test_get_ip_for_subscriber(self):
        """ test get_ip_for_sid """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')

        ip0_returned = self._allocator.get_ip_for_sid('SID0')
        ip1_returned = self._allocator.get_ip_for_sid('SID1')

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip1, ip1_returned)

    def test_get_ip_for_unknown_subscriber(self):
        """ Getting ip for non existent subscriber should return None """
        self.assertIsNone(self._allocator.get_ip_for_sid('SID0'))

    def test_get_sid_for_ip(self):
        """ test get_sid_for_ip """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')

        sid0_returned = self._allocator.get_sid_for_ip(ip0)
        sid1_returned = self._allocator.get_sid_for_ip(ip1)

        self.assertEqual('SID0', sid0_returned)
        self.assertEqual('SID1', sid1_returned)

    def test_get_sid_for_unknown_ip(self):
        """ Getting sid for non allocated ip address should return None """
        self.assertIsNone(
            self._allocator.get_sid_for_ip(ipaddress.ip_address('1.1.1.1')),
        )

    def test_allocate_allocate(self):
        """ Duplicated IP requests for the same UE returns same IP """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID0')
        self.assertEqual(ip0, ip1)

    def test_allocated_release_allocate(self):
        """ Immediate allocation after releasing get the same IP """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        self._allocator.release_ip_address('SID0', ip0)
        ip2, _ = self._allocator.alloc_ip_address('SID0')
        self.assertEqual(ip0, ip2)

    def test_allocate_release_recycle_allocate(self):
        """ Allocation after recycling should get different IPs """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        self._allocator.release_ip_address('SID0', ip0)

        # Wait for auto-recycler to kick in
        time.sleep(1.2 * self.RECYCLING_INTERVAL_SECONDS)

        ip1, _ = self._allocator.alloc_ip_address('SID1')
        ip2, _ = self._allocator.alloc_ip_address('SID0')
        self.assertNotEqual(ip1, ip2)

    def test_recycle_tombstone_ip_on_timer(self):
        """ test recycle tombstone IP on interval loop """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')
        ip2, _ = self._allocator.alloc_ip_address('SID2')
        self._allocator.release_ip_address('SID0', ip0)

        # Wait for auto-recycler to kick in
        time.sleep(2 * self.RECYCLING_INTERVAL_SECONDS)

        ip3, _ = self._allocator.alloc_ip_address('SID3')
        self.assertEqual(ip0, ip3)

        self._allocator.release_ip_address('SID1', ip1)

        # Wait for auto-recycler to kick in
        time.sleep(2 * self.RECYCLING_INTERVAL_SECONDS)

        ip4, _ = self._allocator.alloc_ip_address('SID4')
        self.assertEqual(ip1, ip4)

        self._allocator.release_ip_address('SID2', ip2)

        # Wait for auto-recycler to kick in
        time.sleep(2 * self.RECYCLING_INTERVAL_SECONDS)

        ip5, _ = self._allocator.alloc_ip_address('SID5')
        self.assertEqual(ip2, ip5)

    def test_allocate_unrecycled_IP(self):
        """ Allocation should fail before IP recycling """
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        self._allocator.alloc_ip_address('SID1')
        self._allocator.alloc_ip_address('SID2')
        self._allocator.release_ip_address('SID0', ip0)
        with self.assertRaises(NoAvailableIPError):
            self._allocator.alloc_ip_address('SID3')

    def test_recycle_tombstone_ip(self):
        """ test recycle tombstone IP """
        self._new_ip_allocator(0)

        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')
        ip2, _ = self._allocator.alloc_ip_address('SID2')
        self._allocator.release_ip_address('SID0', ip0)
        ip3, _ = self._allocator.alloc_ip_address('SID3')
        self.assertEqual(ip0, ip3)

        self._allocator.release_ip_address('SID1', ip1)
        ip4, _ = self._allocator.alloc_ip_address('SID4')
        self.assertEqual(ip1, ip4)

        self._allocator.release_ip_address('SID2', ip2)
        ip5, _ = self._allocator.alloc_ip_address('SID5')
        self.assertEqual(ip2, ip5)

    def test_remove_unallocated_block(self):
        """ test removing the allocator for an unallocated block """
        self.assertEqual(
            [self._block], self._allocator.remove_ip_blocks(self._block),
        )

    def test_remove_allocated_block_without_force(self):
        """ test removing the allocator for an allocated block unforcibly """
        self._allocator.alloc_ip_address('SID0')
        self.assertEqual(
            [], self._allocator.remove_ip_blocks(self._block, force=False),
        )

    def test_remove_unforcible_is_default_behavior(self):
        """ test that removing by default is unforcible remove """
        self._allocator.alloc_ip_address('SID0')
        self.assertEqual([], self._allocator.remove_ip_blocks(self._block))

    def test_remove_allocated_block_with_force(self):
        """ test removing the allocator for an allocated block forcibly """
        self._allocator.alloc_ip_address('SID0')
        self.assertEqual(
            [self._block],
            self._allocator.remove_ip_blocks(self._block, force=True),
        )

    def test_remove_after_releasing_all_addresses(self):
        """ removing after releasing all allocated addresses """
        self._new_ip_allocator(0)  # Immediately recycle

        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')

        self.assertEqual(
            [], self._allocator.remove_ip_blocks(self._block, force=False),
        )

        self._allocator.release_ip_address('SID0', ip0)
        self._allocator.release_ip_address('SID1', ip1)

        self.assertEqual(
            [self._block],
            self._allocator.remove_ip_blocks(self._block, force=False),
        )

    def test_remove_after_releasing_some_addresses(self):
        """ removing after releasing all allocated addresses """
        self._new_ip_allocator(0)  # Immediately recycle

        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')

        self.assertEqual(
            [], self._allocator.remove_ip_blocks(self._block, force=False),
        )

        self._allocator.release_ip_address('SID0', ip0)

        self.assertEqual(
            [], self._allocator.remove_ip_blocks(self._block, force=False),
        )

        self.assertTrue(
            ip0 not in self._allocator.list_allocated_ips(self._block),
        )
        self.assertTrue(
            ip1 in self._allocator.list_allocated_ips(self._block),
        )

    def test_reap_after_forced_remove(self):
        """
        test reaping after a forced remove and readding the reaped ips doesn't
        free them
        """
        recycling_interval_seconds = 1  # plenty of time to set up
        self._new_ip_allocator(recycling_interval_seconds)

        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')
        self._allocator.release_ip_address('SID0', ip0)
        self.assertEqual(
            [self._block],
            self._allocator.remove_ip_blocks(self._block, force=True),
        )
        self._allocator.add_ip_block(self._block)
        ip0, _ = self._allocator.alloc_ip_address('SID0')
        ip1, _ = self._allocator.alloc_ip_address('SID1')

        # Wait for auto-recycler to kick in
        time.sleep(recycling_interval_seconds)

        # Ensure that released-then-allocated address doesn't get reaped
        self.assertTrue(
            ip0 in self._allocator.list_allocated_ips(self._block),
        )
        self.assertTrue(
            ip1 in self._allocator.list_allocated_ips(self._block),
        )
