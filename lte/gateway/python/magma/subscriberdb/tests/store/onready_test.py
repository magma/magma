"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access

import asyncio
import unittest

from lte.protos.subscriberdb_pb2 import SubscriberData
from magma.subscriberdb.store.cached_store import CachedStore
from magma.subscriberdb.store.sqlite import SqliteStore

from magma.subscriberdb.sid import SIDUtils


class OnReadyMixinTests(unittest.TestCase):
    """
    Test class for the OnReady subscriber storage mixin
    """

    def setUp(self):
        cache_size = 3
        self.loop = asyncio.new_event_loop()
        sqlite = SqliteStore("file::memory:", loop=self.loop)
        self._store = CachedStore(sqlite, cache_size, self.loop)

    def _add_subscriber(self, sid):
        sub = SubscriberData(sid=SIDUtils.to_pb(sid))
        self._store.add_subscriber(sub)
        return (sid, sub)

    def test_subscriber_addition(self):
        """
        Test if subscriber addition triggers ready
        """
        self.assertEqual(self._store._on_ready.event.is_set(), False)
        self.assertEqual(self._store._persistent_store._on_ready.event.is_set(), False)
        self._add_subscriber('IMSI11111')

        def defer():
            yield from self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)
        self.assertEqual(self._store._persistent_store._on_ready.event.is_set(), True)


    def test_resync(self):
        """
        Test if resync triggers ready
        """
        self.assertEqual(self._store._on_ready.event.is_set(), False)
        self.assertEqual(self._store._persistent_store._on_ready.event.is_set(), False)
        self._store.resync([])

        def defer():
            yield from self._store.on_ready()
        self.loop.run_until_complete(defer())

        self.assertEqual(self._store._on_ready.event.is_set(), True)
        self.assertEqual(self._store._persistent_store._on_ready.event.is_set(), True)
