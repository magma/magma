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

import copy
import threading
from collections import OrderedDict
from contextlib import contextmanager

from magma.subscriberdb.sid import SIDUtils

from .base import BaseStore, DuplicateSubscriberError
from .onready import OnDataReady


class CachedStore(BaseStore):
    """
    A thread-safe cached persistent store of the subscriber database.
    Prerequisite: persistent_store need to be thread safe
    """

    def __init__(self, persistent_store, cache_capacity=512, loop=None):
        self._lock = threading.Lock()
        self._cache = OrderedDict()
        self._cache_capacity = cache_capacity
        self._persistent_store = persistent_store
        self._on_ready = OnDataReady(loop=loop)

    def add_subscriber(self, subscriber_data):
        """
        Method that adds the subscriber.
        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        with self._lock:
            if sid in self._cache:
                raise DuplicateSubscriberError(sid)

            self._persistent_store.add_subscriber(subscriber_data)
            self._cache_put(sid, subscriber_data)
        self._on_ready.add_subscriber(subscriber_data)

    @contextmanager
    def edit_subscriber(self, subscriber_id):
        """
        Context manager to modify the subscriber data.
        """
        with self._lock:
            if subscriber_id in self._cache:
                data = self._cache_get(subscriber_id)
                subscriber_data = copy.deepcopy(data)
            else:
                subscriber_data = \
                    self._persistent_store.get_subscriber_data(subscriber_id)
            yield subscriber_data
            self._persistent_store.update_subscriber(subscriber_data)
            self._cache_put(subscriber_id, subscriber_data)

    def delete_subscriber(self, subscriber_id):
        """
        Method that deletes a subscriber, if present.
        """
        with self._lock:
            if subscriber_id in self._cache:
                del self._cache[subscriber_id]

            self._persistent_store.delete_subscriber(subscriber_id)

    def delete_all_subscribers(self):
        """
        Method that removes all the subscribers from the store
        """
        with self._lock:
            self._cache_clear()
            self._persistent_store.delete_all_subscribers()

    def resync(self, subscribers):
        """
        Method that should resync the store with the mentioned list of
        subscribers. The resync leaves the current state of subscribers
        intact.

        Args:
            subscribers - list of subscribers to be in the store.
        """
        with self._lock:
            self._cache_clear()
            self._persistent_store.resync(subscribers)
        self._on_ready.resync(subscribers)

    def get_subscriber_data(self, subscriber_id):
        """
        Method that returns the subscriber data for the subscriber.
        """
        with self._lock:
            if subscriber_id in self._cache:
                return self._cache_get(subscriber_id)
            else:
                subscriber_data = \
                    self._persistent_store.get_subscriber_data(subscriber_id)
                self._cache_put(subscriber_id, subscriber_data)
                return subscriber_data

    def list_subscribers(self):
        """
        Method that returns the list of subscribers stored.
        Note: this method is not cached since it's easier to get the whole list
        from persistent store
        """
        return self._persistent_store.list_subscribers()

    async def on_ready(self):
        return await self._on_ready.event.wait()

    def _cache_get(self, k):
        """
        Get from the LRU cache. Move the last hit entry to the end.
        """
        self._cache.move_to_end(k)
        return self._cache[k]

    def _cache_put(self, k, v):
        """
        Put to the LRU cache. Evict the first item if full.
        """
        if self._cache_capacity == len(self._cache):
            self._cache.popitem(last=False)
        self._cache[k] = v

    def _cache_list(self):
        return list(self._cache.keys())

    def _cache_clear(self):
        self._cache.clear()
