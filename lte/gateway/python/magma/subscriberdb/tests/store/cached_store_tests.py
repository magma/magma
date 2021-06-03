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

# pylint: disable=protected-access
import tempfile
import unittest

from lte.protos.subscriberdb_pb2 import SubscriberData
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.base import (
    DuplicateSubscriberError,
    SubscriberNotFoundError,
)
from magma.subscriberdb.store.cached_store import CachedStore
from magma.subscriberdb.store.sqlite import SqliteStore


class StoreTests(unittest.TestCase):
    """
    Test class for the CachedStore subscriber storage
    """

    def setUp(self):
        cache_size = 3
        self._tmpfile = tempfile.TemporaryDirectory()
        sqlite = SqliteStore(self._tmpfile.name + '/')
        self._store = CachedStore(sqlite, cache_size)

    def tearDown(self):
        self._tmpfile.cleanup()

    def _add_subscriber(self, sid):
        sub = SubscriberData(sid=SIDUtils.to_pb(sid))
        self._store.add_subscriber(sub)
        return (sid, sub)

    def test_subscriber_addition(self):
        """
        Test if subscriber addition works as expected
        """
        self.assertEqual(self._store.list_subscribers(), [])
        (sid1, _) = self._add_subscriber('IMSI11111')
        self.assertEqual(self._store.list_subscribers(), [sid1])
        (sid2, sub2) = self._add_subscriber('IMSI22222')
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])

        # Check if adding an existing user throws an exception
        with self.assertRaises(DuplicateSubscriberError):
            self._store.add_subscriber(sub2)
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])

        self.assertEqual(self._store._cache_list(), [sid1, sid2])

        self._store.delete_all_subscribers()
        self.assertEqual(self._store.list_subscribers(), [])
        self.assertEqual(self._store._cache_list(), [])

    def test_subscriber_deletion(self):
        """
        Test if subscriber deletion works as expected
        """
        (sid1, _) = self._add_subscriber('IMSI11111')
        (sid2, _) = self._add_subscriber('IMSI22222')
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])
        self.assertEqual(self._store._cache_list(), [sid1, sid2])

        self._store.delete_subscriber(sid2)
        self.assertEqual(self._store.list_subscribers(), [sid1])
        self.assertEqual(self._store._cache_list(), [sid1])

        # Deleting a non-existent user would be ignored
        self._store.delete_subscriber(sid2)
        self.assertEqual(self._store.list_subscribers(), [sid1])
        self.assertEqual(self._store._cache_list(), [sid1])

        self._store.delete_subscriber(sid1)
        self.assertEqual(self._store.list_subscribers(), [])
        self.assertEqual(self._store._cache_list(), [])

    def test_subscriber_retrieval(self):
        """
        Test if subscriber retrieval works as expected
        """
        (sid1, sub1) = self._add_subscriber('IMSI11111')
        self.assertEqual(self._store.list_subscribers(), [sid1])
        self.assertEqual(self._store._cache_list(), [sid1])
        self.assertEqual(self._store.get_subscriber_data(sid1), sub1)

        with self.assertRaises(SubscriberNotFoundError):
            self._store.get_subscriber_data('IMSI30000')
        self.assertEqual(self._store._cache_list(), [sid1])

        self._store.delete_all_subscribers()
        self.assertEqual(self._store.list_subscribers(), [])
        self.assertEqual(self._store._cache_list(), [])

    def test_subscriber_edit(self):
        """
        Test if subscriber edit works as expected
        """
        (sid1, sub1) = self._add_subscriber('IMSI11111')
        self.assertEqual(self._store.get_subscriber_data(sid1), sub1)
        self.assertEqual(self._store._cache_list(), [sid1])

        # Update from cache
        with self._store.edit_subscriber(sid1) as subs:
            subs.lte.auth_key = b'5678'
        self.assertEqual(
            self._store.get_subscriber_data(sid1).lte.auth_key,
            b'5678',
        )
        self.assertEqual(self._store._cache_list(), [sid1])

        # Update from persistent store after eviction
        (sid2, _) = self._add_subscriber('IMSI22222')
        (sid3, _) = self._add_subscriber('IMSI33333')
        (sid4, _) = self._add_subscriber('IMSI44444')
        self.assertEqual(self._store._cache_list(), [sid2, sid3, sid4])
        with self._store.edit_subscriber(sid1) as subs:
            subs.lte.auth_key = b'2468'
        self.assertEqual(
            self._store.get_subscriber_data(sid1).lte.auth_key,
            b'2468',
        )
        self.assertEqual(self._store._cache_list(), [sid3, sid4, sid1])

        with self.assertRaises(SubscriberNotFoundError):
            with self._store.edit_subscriber('IMSI3000') as subs:
                pass

    def test_resync(self):
        """
        Test if resync works as expected
        """
        (sid1, sub1) = self._add_subscriber('IMSI11111')
        (sid2, _) = self._add_subscriber('IMSI11112')
        with self._store.edit_subscriber(sid1) as subs:
            subs.state.lte_auth_next_seq = 1000

        # Resync
        sub1.lte.auth_key = b'5678'
        sub1.state.lte_auth_next_seq = 2000
        self._store.resync([sub1])

        subs = self._store.get_subscriber_data(sid1)
        self.assertEqual(subs.lte.auth_key, b'5678')  # config updated
        self.assertEqual(
            subs.state.lte_auth_next_seq,
            1000,
        )  # state left intact

        with self.assertRaises(SubscriberNotFoundError):
            # sub2 was removed during resync
            self._store.get_subscriber_data(sid2)

    def test_lru_cache_invl(self):
        """
        Test if LRU eviction works as expected
        """
        (sid1, _) = self._add_subscriber('IMSI11111')
        (sid2, _) = self._add_subscriber('IMSI22222')
        (sid3, _) = self._add_subscriber('IMSI33333')
        (sid4, _) = self._add_subscriber('IMSI44444')
        (sid5, _) = self._add_subscriber('IMSI55555')
        (sid6, _) = self._add_subscriber('IMSI66666')

        self._store.get_subscriber_data(sid1)
        self.assertEqual(self._store._cache_list(), [sid5, sid6, sid1])
        self._store.get_subscriber_data(sid2)
        self.assertEqual(self._store._cache_list(), [sid6, sid1, sid2])
        self._store.get_subscriber_data(sid3)
        self.assertEqual(self._store._cache_list(), [sid1, sid2, sid3])

        self._store.get_subscriber_data(sid2)
        self.assertEqual(self._store._cache_list(), [sid1, sid3, sid2])

        self._store.get_subscriber_data(sid4)
        self.assertEqual(self._store._cache_list(), [sid3, sid2, sid4])

        self._store.get_subscriber_data(sid5)
        self.assertEqual(self._store._cache_list(), [sid2, sid4, sid5])

        self._store.get_subscriber_data(sid6)
        self.assertEqual(self._store._cache_list(), [sid4, sid5, sid6])

        self._store.delete_all_subscribers()
        self.assertEqual(self._store.list_subscribers(), [])
        self.assertEqual(self._store._cache_list(), [])
