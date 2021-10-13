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

import tempfile
import unittest

from lte.protos.subscriberdb_pb2 import SubscriberData
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.base import (
    DuplicateSubscriberError,
    SubscriberNotFoundError,
)
from magma.subscriberdb.store.sqlite import SqliteStore
from orc8r.protos.digest_pb2 import Digest, LeafDigest


class StoreTests(unittest.TestCase):
    """
    Test class for subscriber storage
    """

    def setUp(self):
        # Create sqlite3 database for testing
        self._tmpfile = tempfile.TemporaryDirectory()
        self._store = SqliteStore(self._tmpfile.name + '/')

    def tearDown(self):
        self._tmpfile.cleanup()

    def _add_subscriber(self, sid):
        sub = SubscriberData(sid=SIDUtils.to_pb(sid))
        self._store.add_subscriber(sub)
        return (sid, sub)

    def _upsert_subscriber(self, sid):
        sub = SubscriberData(sid=SIDUtils.to_pb(sid))
        self._store.upsert_subscriber(sub)
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

        self._store.delete_all_subscribers()
        self.assertEqual(self._store.list_subscribers(), [])

    def test_subscriber_deletion(self):
        """
        Test if subscriber deletion works as expected
        """
        (sid1, _) = self._add_subscriber('IMSI11111')
        (sid2, _) = self._add_subscriber('IMSI22222')
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])

        self._store.delete_subscriber(sid2)
        self.assertEqual(self._store.list_subscribers(), [sid1])

        # Deleting a non-existent user would be ignored
        self._store.delete_subscriber(sid2)
        self.assertEqual(self._store.list_subscribers(), [sid1])

        self._store.delete_subscriber(sid1)
        self.assertEqual(self._store.list_subscribers(), [])

    def test_subscriber_deletion_digests(self):
        """
        Test if subscriber deletion also unconditionally removes digest info.

        Regression test for #9029.
        """
        (sid1, _) = self._add_subscriber('IMSI11111')
        (sid2, _) = self._add_subscriber('IMSI22222')
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])

        root_digest = "apple"
        leaf_digest = LeafDigest(
            id='IMSI11111',
            digest=Digest(md5_base64_digest="digest_apple"),
        )
        self._store.update_root_digest(root_digest)
        self._store.update_leaf_digests([leaf_digest])
        self.assertNotEqual(self._store.get_current_root_digest(), "")
        self.assertNotEqual(self._store.get_current_leaf_digests(), [])

        self._store.delete_subscriber(sid2)
        self.assertEqual(self._store.list_subscribers(), [sid1])

        # Deleting a subscriber also deletes all digest info
        self.assertEqual(self._store.get_current_root_digest(), "")
        self.assertEqual(self._store.get_current_leaf_digests(), [])

    def test_subscriber_retrieval(self):
        """
        Test if subscriber retrieval works as expected
        """
        (sid1, sub1) = self._add_subscriber('IMSI11111')
        self.assertEqual(self._store.list_subscribers(), [sid1])
        self.assertEqual(self._store.get_subscriber_data(sid1), sub1)

        with self.assertRaises(SubscriberNotFoundError):
            self._store.get_subscriber_data('IMSI30000')

        self._store.delete_all_subscribers()
        self.assertEqual(self._store.list_subscribers(), [])

    def test_subscriber_edit(self):
        """
        Test if subscriber edit works as expected
        """
        (sid1, sub1) = self._add_subscriber('IMSI11111')
        self.assertEqual(self._store.get_subscriber_data(sid1), sub1)

        sub1.lte.auth_key = b'1234'
        self._store.update_subscriber(sub1)
        self.assertEqual(
            self._store.get_subscriber_data(sid1).lte.auth_key,
            b'1234',
        )

        with self._store.edit_subscriber(sid1) as subs:
            subs.lte.auth_key = b'5678'
        self.assertEqual(
            self._store.get_subscriber_data(sid1).lte.auth_key,
            b'5678',
        )

        with self.assertRaises(SubscriberNotFoundError):
            sub1.sid.id = '30000'
            self._store.update_subscriber(sub1)
        with self.assertRaises(SubscriberNotFoundError):
            with self._store.edit_subscriber('IMSI3000') as subs:
                pass

    def test_subscriber_upsert(self):
        """
        Test if subscriber upsertion works as expected
        """
        self.assertEqual(self._store.list_subscribers(), [])
        (sid1, _) = self._upsert_subscriber('IMSI11111')
        self.assertEqual(self._store.list_subscribers(), [sid1])
        (sid2, _) = self._add_subscriber('IMSI22222')
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])

        self._upsert_subscriber('IMSI11111')
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])
        self._upsert_subscriber('IMSI22222')
        self.assertEqual(self._store.list_subscribers(), [sid1, sid2])

        self._store.delete_all_subscribers()
        self.assertEqual(self._store.list_subscribers(), [])

    def test_digest(self):
        """
        Test if digest gets & updates work as expected
        """
        self.assertEqual(self._store.get_current_root_digest(), "")
        self._store.update_root_digest("digest_apple")
        self.assertEqual(self._store.get_current_root_digest(), "digest_apple")
        self._store.update_root_digest("digest_banana")
        self.assertEqual(self._store.get_current_root_digest(), "digest_banana")

    def test_leaf_digests(self):
        """
        Test if leaf digests gets & updates work as expected
        """
        self.assertEqual(self._store.get_current_leaf_digests(), [])
        digests1 = [
            LeafDigest(
                id='IMSI11111',
                digest=Digest(md5_base64_digest='digest_apple'),
            ),
            LeafDigest(
                id='IMSI22222',
                digest=Digest(md5_base64_digest='digest_banana'),
            ),
        ]
        self._store.update_leaf_digests(digests1)
        self.assertEqual(self._store.get_current_leaf_digests(), digests1)

        digests2 = [
            LeafDigest(
                id='IMSI11111',
                digest=Digest(md5_base64_digest='digest_apple'),
            ),
            LeafDigest(
                id='IMSI33333',
                digest=Digest(md5_base64_digest='digest_cherry'),
            ),
            LeafDigest(
                id='IMSI44444',
                digest=Digest(md5_base64_digest='digest_dragonfruit'),
            ),
        ]
        self._store.update_leaf_digests(digests2)
        self.assertEqual(self._store.get_current_leaf_digests(), digests2)


if __name__ == "__main__":
    unittest.main()
