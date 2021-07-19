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
import unittest
from typing import Optional

from lte.protos.apn_pb2 import APNConfiguration
from lte.protos.subscriberdb_pb2 import (
    LTESubscription,
    Non3GPPUserProfile,
    SubscriberData,
    SubscriberState,
)
from magma.common.redis.client import get_default_client
from magma.mobilityd.ip_address_man import IPAddressManager
from magma.mobilityd.ip_allocator_multi_apn import IPAllocatorMultiAPNWrapper
from magma.mobilityd.ip_allocator_pool import IpAllocatorPool
from magma.mobilityd.ip_allocator_static import IPAllocatorStaticWrapper
from magma.mobilityd.ip_descriptor import IPType
from magma.mobilityd.ipv6_allocator_pool import IPv6AllocatorPool
from magma.mobilityd.mobility_store import MobilityStore
from magma.subscriberdb.sid import SIDUtils


class MockedSubscriberDBStub:
    # subscriber map
    subs = {}

    def __init__(self):
        pass

    def GetSubscriberData(self, sid):
        cls = self.__class__
        return cls.subs.get(str(sid), None)

    @classmethod
    def add_sub(
        cls, sid: str, apn: str, ip: str, vlan: str = None,
        gw_ip=None, gw_mac=None,
    ):
        sub_db_sid = SIDUtils.to_pb(sid)
        lte = LTESubscription()
        lte.state = LTESubscription.ACTIVE
        state = SubscriberState()
        state.lte_auth_next_seq = 1
        non_3gpp = Non3GPPUserProfile()
        subs_data = SubscriberData(
            sid=sub_db_sid, lte=lte, state=state,
            non_3gpp=non_3gpp,
        )

        cls.subs[str(sub_db_sid)] = subs_data
        cls.add_sub_ip(sid, apn, ip, vlan, gw_ip, gw_mac)

    @classmethod
    def add_incomplete_sub(cls, sid: str):
        sub_db_sid = SIDUtils.to_pb(sid)
        lte = LTESubscription()
        lte.state = LTESubscription.ACTIVE
        state = SubscriberState()
        state.lte_auth_next_seq = 1
        subs_data = SubscriberData(sid=sub_db_sid, lte=lte, state=state)
        cls.subs[str(sub_db_sid)] = subs_data

    @classmethod
    def add_sub_ip(
        cls, sid: str, apn: str, ip: str, vlan: str = None,
        gw_ip=None, gw_mac=None,
    ):
        sub_db_sid = SIDUtils.to_pb(sid)
        apn_config = APNConfiguration()
        apn_config.context_id = 1
        apn_config.service_selection = apn
        if ip:
            apn_config.assigned_static_ip = ip
        if vlan:
            apn_config.resource.vlan_id = int(vlan)
        if gw_ip:
            apn_config.resource.gateway_ip = gw_ip
        if gw_mac:
            apn_config.resource.gateway_mac = gw_mac

        subs_data = cls.subs[str(sub_db_sid)]
        subs_data.non_3gpp.apn_config.extend([apn_config])

    @classmethod
    def clear_subs(cls):
        cls.subs = {}


class MultiAPNIPAllocationTests(unittest.TestCase):
    """
    Test class for the Mobilityd Multi APN Allocator
    """
    RECYCLING_INTERVAL_SECONDS = 1

    def _new_ip_allocator(self, recycling_interval):
        """
        Creates and sets up an IPAllocator with the given recycling interval.
        """
        store = MobilityStore(get_default_client(), False, 3980)
        ip_allocator = IpAllocatorPool(store)
        ip_allocator_static = IPAllocatorStaticWrapper(
            store, MockedSubscriberDBStub(), ip_allocator,
        )
        ipv4_allocator = IPAllocatorMultiAPNWrapper(
            store,
            subscriberdb_rpc_stub=MockedSubscriberDBStub(),
            ip_allocator=ip_allocator_static,
        )
        ipv6_allocator = IPv6AllocatorPool(
            store,
            session_prefix_alloc_mode='RANDOM',
        )
        self._allocator = IPAddressManager(
            ipv4_allocator,
            ipv6_allocator,
            store,
            recycling_interval,
        )
        self._allocator.add_ip_block(self._block)

    def setUp(self):
        self._block = ipaddress.ip_network('192.168.0.0/28')
        self._new_ip_allocator(self.RECYCLING_INTERVAL_SECONDS)

    def tearDown(self):
        MockedSubscriberDBStub.clear_subs()

    def check_type(self, sid: str, type1: IPType):
        ip_desc = self._allocator._store.sid_ips_map[sid]
        self.assertEqual(ip_desc.type, type1)

    def check_vlan(self, sid: str, vlan: str):
        ip_desc = self._allocator._store.sid_ips_map[sid]
        logging.info(
            "type ip_desc.vlan_id %s vlan %s", type(ip_desc.vlan_id),
            type(vlan),
        )
        self.assertEqual(ip_desc.vlan_id, vlan)

    def check_gw_info(
        self, vlan: Optional[int], gw_ip: str,
        gw_mac: Optional[str],
    ):
        gw_info_ip = self._allocator._store.dhcp_gw_info.get_gw_ip(vlan)
        self.assertEqual(gw_info_ip, gw_ip)
        gw_info_mac = self._allocator._store.dhcp_gw_info.get_gw_mac(vlan)
        self.assertEqual(gw_info_mac, gw_mac)

    def test_get_ip_vlan_for_subscriber(self):
        """ test get_ip_for_sid without any assignment """
        sid = 'IMSI11,ipv4'
        ip0, _ = self._allocator.alloc_ip_address(sid)

        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.check_type(sid, IPType.IP_POOL)

    def test_get_ip_vlan_for_subscriber_with_apn(self):
        """ test get_ip_for_sid with static IP """
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.4'
        vlan = 132
        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn=apn, ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)
        self.check_vlan(sid, vlan)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_different_apn(self):
        """ test get_ip_for_sid with different APN assigned ip"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.4'
        vlan = 188
        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn="xyz", ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertNotEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.IP_POOL)
        self.check_vlan(sid, 0)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_wildcard_apn(self):
        """ test wildcard apn"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.4'
        vlan = 166

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn="*", ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)
        self.check_vlan(sid, vlan)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_wildcard_and_exact_apn(self):
        """ test IP assignement from multiple  APNs"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.4'
        wild_assigned_ip = '44.2.3.11'

        vlan = 44
        vlan_wild = 66

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn="*", ip=wild_assigned_ip,
            vlan=vlan_wild,
        )
        MockedSubscriberDBStub.add_sub_ip(
            sid=imsi, apn=apn, ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)
        self.check_vlan(sid, vlan)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_invalid_ip(self):
        """ test invalid data from DB """
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.hh'
        vlan = 111

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn=apn, ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertNotEqual(str(ip0), assigned_ip)
        self.check_type(sid, IPType.IP_POOL)
        self.check_vlan(sid, 0)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_multi_apn_but_no_match(self):
        """ test IP assignment from multiple  APNs"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.4'
        vlan = 31
        vlan_wild = 552

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn="abc", ip=assigned_ip,
            vlan=vlan_wild,
        )
        MockedSubscriberDBStub.add_sub_ip(
            sid=imsi, apn="xyz", ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertNotEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.IP_POOL)
        self.check_vlan(sid, 0)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_incomplete_sub(self):
        """ test IP assignment from subscriber without non_3gpp config"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        MockedSubscriberDBStub.add_incomplete_sub(sid=imsi)

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.check_type(sid, IPType.IP_POOL)
        self.check_vlan(sid, 0)

    def test_get_ip_vlan_for_subscriber_with_wildcard_no_apn(self):
        """ test wildcard apn"""
        imsi = 'IMSI110'
        sid = imsi + ",ipv4"
        assigned_ip = '1.2.3.4'
        vlan = 122

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn="*", ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)
        self.check_vlan(sid, vlan)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_apn_dot(self):
        """ test get_ip_for_sid with static IP """
        apn = 'magma.ipv4'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.4'
        vlan = 165

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn=apn, ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)
        self.check_vlan(sid, vlan)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_wildcard_and_no_exact_apn(self):
        """ test IP assignement from multiple  APNs"""
        apn = 'dsddf'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        assigned_ip = '1.2.3.4'
        wild_assigned_ip = '44.2.3.11'

        vlan = 0
        vlan_wild = 66

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn="*", ip=wild_assigned_ip,
            vlan=vlan_wild,
        )
        MockedSubscriberDBStub.add_sub_ip(
            sid=imsi, apn="xyz", ip=assigned_ip,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(wild_assigned_ip))
        self.check_type(sid, IPType.STATIC)
        self.check_vlan(sid, vlan_wild)
        self.check_gw_info(vlan, None, None)

    def test_get_ip_vlan_for_subscriber_with_wildcard_and_exact_apn_no_ip(
            self,
    ):
        """ test IP assignement from multiple  APNs"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn + ",ipv4"
        wild_assigned_ip = '44.2.3.11'

        vlan = 44
        vlan_wild = 66

        MockedSubscriberDBStub.add_sub(
            sid=imsi, apn="*", ip=wild_assigned_ip,
            vlan=vlan_wild,
        )
        MockedSubscriberDBStub.add_sub_ip(
            sid=imsi, apn=apn, ip=None,
            vlan=vlan,
        )

        ip0, _ = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertNotEqual(ip0, ipaddress.ip_address(wild_assigned_ip))
        self.check_type(sid, IPType.IP_POOL)
        self.check_vlan(sid, 0)
        self.check_gw_info(vlan, None, None)
