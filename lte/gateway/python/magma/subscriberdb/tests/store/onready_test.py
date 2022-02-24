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

import asyncio
import tempfile
import unittest

from lte.protos.subscriberdb_pb2 import SubscriberData
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.sqlite import SqliteStore
from orc8r.protos.digest_pb2 import Digest, LeafDigest


class OnReadyMixinTests(unittest.TestCase):
    """
    Test class for the OnReady subscriber storage mixin
    """

    def setUp(self):
        self.loop = asyncio.new_event_loop()
        self._tmpfile = tempfile.TemporaryDirectory()
        self._store = SqliteStore(self._tmpfile.name + '/', loop=self.loop)

    def tearDown(self):
        self._tmpfile.cleanup()

    def _add_subscriber(self, sid):
        sub = SubscriberData(sid=SIDUtils.to_pb(sid))
        self._store.add_subscriber(sub)
        return (sid, sub)

    def test_subscriber_addition(self):
        """
        Test if subscriber addition triggers ready
        """
        self.assertEqual(self._store._on_ready.event.is_set(), False)
        self._add_subscriber('IMSI11111')

        async def defer():
            await self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)

    def test_resync(self):
        """
        Test if resync triggers ready
        """
        self.assertEqual(self._store._on_ready.event.is_set(), False)
        self._store.resync([])

        async def defer():
            await self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)

    def test_delete_subscriber(self):
        """
        Test if subscriber deletion triggers ready
        """
        self.assertEqual(self._store._on_ready.event.is_set(), False)
        self._store.delete_subscriber('IMSI11111')

        async def defer():
            await self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)

    def test_upsert_subscriber(self):
        """
        Test if subscriber upsertion triggers ready
        """
        self.assertEqual(self._store._on_ready.event.is_set(), False)
        self._store.upsert_subscriber(
            SubscriberData(sid=SIDUtils.to_pb('IMSI1111')),
        )

        async def defer():
            await self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)


class OnDigestsReadyMixinTests(unittest.TestCase):
    """
    Test class for the OnDigestsReady subscriber digests storage mixin
    """

    def setUp(self):
        self.loop = asyncio.new_event_loop()
        self._tmpfile = tempfile.TemporaryDirectory()
        self._store = SqliteStore(self._tmpfile.name + '/', loop=self.loop)

    def tearDown(self):
        self._tmpfile.cleanup()

    def test_leaf_digests_update(self):
        """
        Test if leaf digests update triggers ready
        """
        self.assertEqual(self._store._on_digests_ready.event.is_set(), False)
        self._store.update_leaf_digests([
            LeafDigest(
                id='IMSI11111',
                digest=Digest(md5_base64_digest='digest_cherry'),
            ),
        ])

        async def defer():
            await self._store.on_digests_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_digests_ready.event.is_set(), True)

    def test_root_digest_update(self):
        """
        Test if root digest update triggers ready
        """
        self.assertEqual(self._store._on_digests_ready.event.is_set(), False)
        self._store.update_root_digest("digest_apple")

        async def defer():
            await self._store.on_digests_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_digests_ready.event.is_set(), True)
