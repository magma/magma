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

from lte.protos.policydb_pb2 import ChargingRuleNameSet
from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisHashDict
from magma.common.redis.serializers import (
    get_proto_deserializer,
    get_proto_serializer,
)


class BaseNameDict(RedisHashDict):
    """
    BaseNameDict uses the RedisHashDict collection to store basenames
    and their associated rule ids. Setting and deleting items in the dictionary
    syncs with Redis automatically.
    """
    _DICT_HASH = "policydb:basenames"
    _NOTIFY_CHANNEL = "policydb:basenames:stream_update"

    def __init__(self):
        client = get_default_client()
        super().__init__(
            client,
            self._DICT_HASH,
            get_proto_serializer(),
            get_proto_deserializer(ChargingRuleNameSet),
        )

    def send_update_notification(self):
        """
        Use Redis pub/sub channels to send notifications. Subscribers can listen
        to this channel to know when an update is done
        """
        self.redis.publish(self._NOTIFY_CHANNEL, "Stream Update")

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None
