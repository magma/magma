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
import threading

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisHashDict
from magma.common.redis.serializers import (
    get_json_deserializer,
    get_json_serializer,
)


class TunnelToTunnelMapper:
    """
    Interface ID to Number Mapper

    This class maintains dictionary from tunnel id to tunnel id
    """

    def __init__(self):
        self._tunnel_map = {}
        self._lock = threading.Lock()  # write lock

    def setup_redis(self):
        self._tunnel_map = TunnelDict()

    def get_tunnel(self, tunnel: int):
        with self._lock:
            if tunnel not in self._tunnel_map:
                return None
            return self._tunnel_map[tunnel]

    def save_tunnels(self, uplink_tunnel: int, downlink_tunnel: int):
        with self._lock:
            self._tunnel_map[uplink_tunnel] = downlink_tunnel
            self._tunnel_map[downlink_tunnel] = uplink_tunnel


class TunnelDict(RedisHashDict):
    """
    TunnelDict uses the RedisHashDict collection to store a mapping of
    tunnel id to tunnel id.
    Setting and deleting items in the dictionary syncs with Redis automatically
    """
    _DICT_HASH = "pipelined:tunnel_map"

    def __init__(self):
        client = get_default_client()
        super().__init__(
            client,
            self._DICT_HASH,
            get_json_serializer(), get_json_deserializer())

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None
