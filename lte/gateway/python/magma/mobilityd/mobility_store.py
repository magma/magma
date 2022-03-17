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
from typing import Callable, Dict

import redis
from magma.common.redis.containers import RedisFlatDict, RedisHashDict, RedisSet
from magma.common.redis.serializers import (
    RedisSerde,
    get_json_deserializer,
    get_json_serializer,
)
from magma.mobilityd import serialize_utils
from magma.mobilityd.ip_descriptor import IPDesc, IPState
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
        self, client: redis.Redis,
    ):
        self.init_store(client)

    def init_store(
        self, client: redis.Redis,
    ):
        get_ip_states: Callable[[IPState], Dict[str, IPDesc]] = lambda key: ip_states(client, key)
        self.ip_state_map = IpDescriptorMap(
            defaultdict_key(get_ip_states),  # type: ignore[arg-type]
        )
        self.ipv6_state_map = IpDescriptorMap(
            defaultdict_key(get_ip_states),  # type: ignore[arg-type]
        )
        self.assigned_ip_blocks = AssignedIpBlocksSet(client)
        self.sid_ips_map = IPDescDict(client)
        self.dhcp_gw_info = UplinkGatewayInfo(GatewayInfoMap(client))
        self.dhcp_store = MacToIP(client)  # mac => DHCP_State
        self.allocated_iid = AllocatedIID(client)
        self.sid_session_prefix_allocated = AllocatedSessionPrefix(client)


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
        self[key] = self.default_factory(key)  # pylint: disable=E1102
        return self[key]


class MacToIP(RedisFlatDict):
    """
    Used for managing DHCP state of a Mac address.
    """

    def __init__(self, client):
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

    def __init__(self, client):
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

    def __init__(self, client):
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

    def __init__(self, client):
        serde = RedisSerde(
            ALLOCATED_SESSION_PREFIX_TYPE,
            get_json_serializer(), get_json_deserializer(),
        )
        super().__init__(client, serde)

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None
