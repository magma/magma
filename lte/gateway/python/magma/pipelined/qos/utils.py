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

import logging
from collections import deque

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisHashDict
from magma.common.redis.serializers import (
    get_json_deserializer,
    get_json_serializer,
)

LOG = logging.getLogger('pipelined.qos.id_manager')


class IdManager(object):
    """
    Simple utility class to manage IDs
    """
    def __init__(self, start_idx, max_idx):
        self._start_idx = start_idx
        self._max_idx = max_idx
        self._counter = start_idx
        self._free_idx_list = deque()
        self._restore_done = False

    def allocate_idx(self,) -> int:
        idx = self._get_free_idx()
        if idx is None:
            idx = self._counter
            if idx == self._max_idx:
                raise ValueError("maximum id allocation exceeded")
            self._counter += 1
        LOG.debug("allocating idx %d ", idx)
        return idx

    def release_idx(self, idx):
        LOG.debug("releasing idx %d ", idx)
        if idx < self._start_idx or idx > (self._max_idx - 1):
            LOG.error("attempting to release invalid idx %d", idx)
            return

        self._free_idx_list.append(idx)

    def restore_state(self, id_set):
        if self._restore_done:
            return
        if not id_set:
            return

        self._counter = min(self._max_idx, max(id_set) + 1)
        for idx in range(self._start_idx, self._counter):
            if idx not in id_set:
                self._free_idx_list.append(idx)
        self._restore_done = True

    def _get_free_idx(self) -> int:
        if self._free_idx_list:
            return self._free_idx_list.popleft()
        return None


class QosStore(RedisHashDict):
    def __init__(self, redis_type):
        self.client = get_default_client()
        super().__init__(self.client, redis_type,
                         get_json_serializer(), get_json_deserializer())
