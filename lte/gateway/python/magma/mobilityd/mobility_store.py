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
from collections import defaultdict

import redis
from lte.protos.mobilityd_pb2 import GWInfo
from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisFlatDict, RedisHashDict, RedisSet
from magma.common.redis.serializers import (
    RedisSerde,
    get_json_deserializer,
    get_json_serializer,
)
from magma.mobilityd import serialize_utils
from magma.mobilityd.ip_descriptor import IPDesc
from magma.mobilityd.ip_descriptor_map import IpDescriptorMap
from magma.mobilityd.uplink_gw import UplinkGatewayInfo

IPDESC_REDIS_TYPE = "mobilityd_ipdesc_record"
IPSTATES_REDIS_TYPE = "mobilityd:ip_states:{}"
IPBLOCKS_REDIS_TYPE = "mobilityd:assigned_ip_blocks"
MAC_TO_IP_REDIS_TYPE = "mobilityd_mac_to_ip"
DHCP_GW_INFO_REDIS_TYPE = "mobilityd_gw_info"
ALLOCATED_IID_REDIS_TYPE = "mobilityd_allocated_iid"
ALLOCATED_SESSION_PREFIX_TYPE = "mobilityd_allocated_session_prefix"


class MobilityStore(object):
    def __init__(
        self, client: redis.Redis, persist_to_redis: bool,
        redis_port: int,
    ):
        self.init_store(client, persist_to_redis, redis_port)

    def init_store(
        self, client: redis.Redis, persist_to_redis: bool,
        redis_port: int,
    ):
        if not persist_to_redis:
            self.ip_state_map = IpDescriptorMap(defaultdict(dict))
            self.ipv6_state_map = IpDescriptorMap(defaultdict(dict))
            self.assigned_ip_blocks = set()  # {ip_block}
            self.sid_ips_map = defaultdict(IPDesc)  # {SID=>IPDesc}
            self.dhcp_gw_info = UplinkGatewayInfo(defaultdict(GWInfo))
            self.dhcp_store = {}  # mac => DHCP_State
            self.allocated_iid = {}  # {ipv6 interface identifiers}
            self.sid_session_prefix_allocated = {}  # SID => session prefix
        else:
            if not redis_port:
                raise ValueError(
                    'Must specify a redis_port in mobilityd config.',
                )
            self.ip_state_map = IpDescriptorMap(
                defaultdict_key(lambda key: ip_states(client, key)),
            )
            self.ipv6_state_map = IpDescriptorMap(
                defaultdict_key(lambda key: ip_states(client, key)),
            )
            self.assigned_ip_blocks = AssignedIpBlocksSet(client)
            self.sid_ips_map = IPDescDict(client)
            self.dhcp_gw_info = UplinkGatewayInfo(GatewayInfoMap())
            self.dhcp_store = MacToIP()  # mac => DHCP_State
            self.allocated_iid = AllocatedIID()
            self.sid_session_prefix_allocated = AllocatedSessionPrefix()


class AssignedIpBlocksSet(RedisSet):
    def __init__(self, client):
        super().__init__(
            client,
            IPBLOCKS_REDIS_TYPE,
            serialize_utils.serialize_ip_block,
            serialize_utils.deserialize_ip_block,
        )


class IPDescDict(RedisFlatDict):
    def __init__(self, client):
        serde = RedisSerde(
            IPDESC_REDIS_TYPE,
            serialize_utils.serialize_ip_desc,
            serialize_utils.deserialize_ip_desc,
        )
        super().__init__(client, serde, writethrough=True)


def ip_states(client, key):
    """ Get Redis view of IP states. """
    redis_dict = RedisHashDict(
        client,
        IPSTATES_REDIS_TYPE.format(key),
        serialize_utils.serialize_ip_desc,
        serialize_utils.deserialize_ip_desc,
    )
    return redis_dict


class defaultdict_key(defaultdict):
    """
    Same as standard lib's defaultdict, but takes the key as a parameter.
    """

    def __missing__(self, key):
        # Follow defaultdict pattern in raising KeyError
        if self.default_factory is None:
            raise KeyError(key)
        self[key] = self.default_factory(key)
        return self[key]


class MacToIP(RedisFlatDict):
    """
    Used for managing DHCP state of a Mac address.
    """

    def __init__(self):
        client = get_default_client()
        serde = RedisSerde(
            MAC_TO_IP_REDIS_TYPE,
            get_json_serializer(), get_json_deserializer(),
        )
        super().__init__(client, serde)

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None


class GatewayInfoMap(RedisFlatDict):
    """
    Used for mainatining uplink GW info
    """

    def __init__(self):
        client = get_default_client()
        serde = RedisSerde(
            DHCP_GW_INFO_REDIS_TYPE,
            get_json_serializer(), get_json_deserializer(),
        )
        super().__init__(client, serde)

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None


class AllocatedIID(RedisFlatDict):
    """
    Used for tracking already allocated Interface identifiers for IPv6
    allocation
    """

    def __init__(self):
        client = get_default_client()
        serde = RedisSerde(
            ALLOCATED_IID_REDIS_TYPE,
            get_json_serializer(), get_json_deserializer(),
        )
        super().__init__(client, serde)

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None


class AllocatedSessionPrefix(RedisFlatDict):
    """
    Used for tracking already allocated session prefix for IPv6 allocation
    """

    def __init__(self):
        client = get_default_client()
        serde = RedisSerde(
            ALLOCATED_SESSION_PREFIX_TYPE,
            get_json_serializer(), get_json_deserializer(),
        )
        super().__init__(client, serde)

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None
