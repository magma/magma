"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from collections import defaultdict

from magma.common.redis.containers import RedisFlatDict, RedisHashDict, \
    RedisSet
from magma.common.redis.serializers import RedisSerde
from magma.mobilityd import serialize_utils

IPDESC_REDIS_TYPE = "mobilityd_ipdesc_record"
IPSTATES_REDIS_TYPE = "mobilityd:ip_states:{}"
IPBLOCKS_REDIS_TYPE = "mobilityd:assigned_ip_blocks"


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
        serde = RedisSerde(IPDESC_REDIS_TYPE,
                           serialize_utils.serialize_ip_desc,
                           serialize_utils.deserialize_ip_desc,
                           )
        super().__init__(client, serde)


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
