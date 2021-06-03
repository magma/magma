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
import unittest

import fakeredis
from magma.mobilityd.ip_address_man import IPNotInUseError
from magma.mobilityd.ip_descriptor import IPDesc
from magma.mobilityd.ipv6_allocator_pool import IPv6AllocatorPool
from magma.mobilityd.mobility_store import MobilityStore


class TestIPV6Allocator(unittest.TestCase):
    """
    Test class for the Mobilityd IPv6 Allocator
    """

    def _new_ip_allocator(self, block):
        """
        Creates and sets up an IPAllocator with the given IPv6 block.
        """
        store = MobilityStore(fakeredis.FakeStrictRedis(), False, 3980)
        allocator = IPv6AllocatorPool(store, 'RANDOM')
        allocator.add_ip_block(block)
        return allocator

    def setUp(self):
        self._block = ipaddress.ip_network('fedd:5:6c::/48')
        self._allocator = self._new_ip_allocator(self._block)

    def test_alloc_ipv6_address(self):
        """ test alloc_ip_address """
        ip0 = self._allocator.alloc_ip_address('SID0', 0).ip
        self.assertTrue(ip0 in self._block)

        ip1 = self._allocator.alloc_ip_address('SID1', 0).ip
        self.assertTrue(ip1 in self._block)
        self.assertNotEqual(ip1, ip0)

        ip2 = self._allocator.alloc_ip_address('SID2', 0).ip
        self.assertTrue(ip2 in self._block)
        self.assertNotEqual(ip2, ip0)
        self.assertNotEqual(ip2, ip1)

    def test_release_ipv6_address(self):
        """ test release_ip_address """
        ip0 = self._allocator.alloc_ip_address('SID0', 0).ip

        # release ip
        ip_desc = IPDesc(ip=ip0, sid='SID0')
        self._allocator.release_ip(ip_desc)

        # double release
        with self.assertRaises(IPNotInUseError):
            self._allocator.release_ip(ip_desc)

    def test_remove_unallocated_ipv6_block(self):
        """ test removing the allocator for an unallocated block """
        self.assertEqual([], self._allocator.remove_ip_blocks(self._block))

    def test_remove_after_releasing_some_addresses(self):
        """ removing after releasing all allocated addresses """
        self._new_ip_allocator(self._block)

        ip0 = self._allocator.alloc_ip_address('SID0', 0).ip
        ip1 = self._allocator.alloc_ip_address('SID1', 0).ip
        ip2 = self._allocator.alloc_ip_address('SID2', 0).ip

        ip_desc0 = IPDesc(ip=ip0, sid='SID0')
        ip_desc1 = IPDesc(ip=ip1, sid='SID1')
        ip_desc2 = IPDesc(ip=ip2, sid='SID2')
        self._allocator.release_ip(ip_desc0)
        self._allocator.release_ip(ip_desc1)
        self._allocator.release_ip(ip_desc2)

        self.assertEqual(
            {},
            self._allocator._store.sid_session_prefix_allocated,
        )
