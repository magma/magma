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
from magma.subscriberdb.store.cached_store import CachedStore
from magma.subscriberdb.store.sqlite import SqliteStore


class OnReadyMixinTests(unittest.TestCase):
    """
    Test class for the OnReady subscriber storage mixin
    """

    def setUp(self):
        cache_size = 3
        self.loop = asyncio.new_event_loop()
        self._tmpfile = tempfile.TemporaryDirectory()
        sqlite = SqliteStore(self._tmpfile.name + '/', loop=self.loop)
        self._store = CachedStore(sqlite, cache_size, self.loop)

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
        self.assertEqual(
            self._store._persistent_store._on_ready.event.is_set(), False,
        )
        self._add_subscriber('IMSI11111')

        async def defer():
            await self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)
        self.assertEqual(
            self._store._persistent_store._on_ready.event.is_set(), True,
        )

    def test_resync(self):
        """
        Test if resync triggers ready
        """
        self.assertEqual(self._store._on_ready.event.is_set(), False)
        self.assertEqual(
            self._store._persistent_store._on_ready.event.is_set(), False,
        )
        self._store.resync([])

        async def defer():
            await self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)
        self.assertEqual(
            self._store._persistent_store._on_ready.event.is_set(), True,
        )
