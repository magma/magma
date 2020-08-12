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
import unittest

from magma.mobilityd import subscriberdb_client

from lte.protos.subscriberdb_pb2 import (
    LTESubscription,
    SubscriberData,
    SubscriberState,
    SubscriberID,
    SubscriberUpdate,
    Non3GPPUserProfile,
    APNConfiguration,
)

from lte.protos.mconfig.mconfigs_pb2 import MobilityD
from magma.mobilityd.ip_descriptor import IPDesc, IPType
from magma.mobilityd.ip_address_man import IPAddressManager, \
    IPNotInUseError, MappingNotFoundError
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
    def add_sub(cls, sid: str, apn: str, ip: str):
        sub_db_sid = SIDUtils.to_pb(sid)
        lte = LTESubscription()
        lte.state = LTESubscription.ACTIVE
        state = SubscriberState()
        state.lte_auth_next_seq = 1
        non_3gpp = Non3GPPUserProfile()
        subs_data = SubscriberData(sid=sub_db_sid, lte=lte, state=state, non_3gpp=non_3gpp)

        cls.subs[str(sub_db_sid)] = subs_data
        cls.add_sub_ip(sid, apn, ip)

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
    def add_sub_ip(cls, sid: str, apn: str, ip: str):
        sub_db_sid = SIDUtils.to_pb(sid)
        apn_config = APNConfiguration()
        apn_config.context_id = 1
        apn_config.service_selection = apn
        apn_config.assigned_static_ip = ip

        subs_data = cls.subs[str(sub_db_sid)]
        subs_data.non_3gpp.apn_config.extend([apn_config])

    @classmethod
    def clear_subs(cls):
        cls.subs = {}


class StaticIPAllocationTests(unittest.TestCase):
    """
    Test class for the Mobilityd Static IP Allocator
    """
    RECYCLING_INTERVAL_SECONDS = 1

    def _new_ip_allocator(self, recycling_interval):
        """
        Creates and sets up an IPAllocator with the given recycling interval.
        """
        config = {
            'recycling_interval': recycling_interval,
            'persist_to_redis': False,
            'redis_port': 6379,
        }
        mconfig = MobilityD(ip_allocator_type=MobilityD.IP_POOL,
                            static_ip_enabled=True)

        self._allocator = IPAddressManager(recycling_interval=recycling_interval,
                                           subscriberdb_rpc_stub=MockedSubscriberDBStub(),
                                           config=config,
                                           mconfig=mconfig)
        self._allocator.add_ip_block(self._block)

    def setUp(self):
        self._block = ipaddress.ip_network('192.168.0.0/28')
        self._new_ip_allocator(self.RECYCLING_INTERVAL_SECONDS)

    def tearDown(self):
        MockedSubscriberDBStub.clear_subs()

    def check_type(self, sid: str, type: IPType):
        ip_desc = self._allocator.sid_ips_map[sid]
        self.assertEqual(ip_desc.type, type)

    def test_get_ip_for_subscriber(self):
        """ test get_ip_for_sid without any assignment """
        sid = 'IMSI11'
        ip0 = self._allocator.alloc_ip_address(sid)

        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.check_type(sid, IPType.IP_POOL)

    def test_get_ip_for_subscriber_with_apn(self):
        """ test get_ip_for_sid with static IP """
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        assigned_ip = '1.2.3.4'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn=apn, ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)

    def test_get_ip_for_subscriber_with_different_apn(self):
        """ test get_ip_for_sid with different APN assigned ip"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        assigned_ip = '1.2.3.4'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn="xyz", ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertNotEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.IP_POOL)

    def test_get_ip_for_subscriber_with_wildcard_apn(self):
        """ test wildcard apn"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        assigned_ip = '1.2.3.4'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn="*", ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)

    def test_get_ip_for_subscriber_with_wildcard_and_exact_apn(self):
        """ test IP assignement from multiple  APNs"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        assigned_ip = '1.2.3.4'
        assigned_ip_wild = '22.22.22.22'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn="*", ip=assigned_ip_wild)
        MockedSubscriberDBStub.add_sub_ip(sid=imsi, apn=apn, ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)

    def test_get_ip_for_subscriber_with_invalid_ip(self):
        """ test invalid data from DB """
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        assigned_ip = '1.2.3.hh'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn=apn, ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertNotEqual(str(ip0), assigned_ip)
        self.check_type(sid, IPType.IP_POOL)

    def test_get_ip_for_subscriber_with_multi_apn_but_no_match(self):
        """ test IP assignment from multiple  APNs"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        assigned_ip = '1.2.3.4'
        assigned_ip_wild = '22.22.22.22'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn="abc", ip=assigned_ip_wild)
        MockedSubscriberDBStub.add_sub_ip(sid=imsi, apn="xyz", ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertNotEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.IP_POOL)

    def test_get_ip_for_subscriber_with_incomplete_sub(self):
        """ test IP assignment from subscriber without non_3gpp config"""
        apn = 'magma'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        MockedSubscriberDBStub.add_incomplete_sub(sid=imsi)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.check_type(sid, IPType.IP_POOL)

    def test_get_ip_for_subscriber_with_wildcard_no_apn(self):
        """ test wildcard apn"""
        imsi = 'IMSI110'
        sid = imsi
        assigned_ip = '1.2.3.4'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn="*", ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)

    def test_get_ip_for_subscriber_with_apn_dot(self):
        """ test get_ip_for_sid with static IP """
        apn = 'magma.ipv4'
        imsi = 'IMSI110'
        sid = imsi + '.' + apn
        assigned_ip = '1.2.3.4'
        MockedSubscriberDBStub.add_sub(sid=imsi, apn=apn, ip=assigned_ip)

        ip0 = self._allocator.alloc_ip_address(sid)
        ip0_returned = self._allocator.get_ip_for_sid(sid)

        # check if retrieved ip is the same as the one allocated
        self.assertEqual(ip0, ip0_returned)
        self.assertEqual(ip0, ipaddress.ip_address(assigned_ip))
        self.check_type(sid, IPType.STATIC)

