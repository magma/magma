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
import threading

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisHashDict
from magma.common.redis.serializers import (
    get_json_deserializer,
    get_json_serializer,
)


class InterfaceIDToPrefixMapper:
    """
    Interface ID to Number Mapper

    This class maintains dictionary from interface id to ipv6 prefixes
    """

    def __init__(self):
        self._prefix_by_interface = {}
        self._lock = threading.Lock()  # write lock

    def setup_redis(self):
        self._prefix_by_interface = PrefixDict()

    def get_prefix(self, interface):
        with self._lock:
            if interface not in self._prefix_by_interface:
                return None
            return self._prefix_by_interface[interface]

    def save_prefix(self, interface, prefix):
        with self._lock:
            self._prefix_by_interface[interface] = prefix


class PrefixDict(RedisHashDict):
    """
    PrefixDict uses the RedisHashDict collection to store a mapping of
    interface id to ipv6 prefix.
    Setting and deleting items in the dictionary syncs with Redis automatically
    """
    _DICT_HASH = "pipelined:ipv6_prefixes"

    def __init__(self):
        client = get_default_client()
        super().__init__(
            client,
            self._DICT_HASH,
            get_json_serializer(), get_json_deserializer())

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None


def get_ipv6_interface_id(ipv6: str) -> str:
    """
    Retrieve the interface id out of the lower 64 bits
    """
    ipv6_block = ipaddress.ip_address(ipv6)
    interface = ipaddress.ip_address(int(ipv6_block) & 0xffffffffffffffff)

    return str(interface)


def get_ipv6_prefix(ipv6: str) -> str:
    """
    Retrieve the prefix out of the higher 64 bits
    """
    ipv6_block = ipaddress.ip_address(ipv6)
    interface = ipaddress.ip_address(
        int(ipv6_block) & 0xffffffffffffffff0000000000000000)

    return str(interface)
