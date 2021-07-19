#!/usr/bin/env python3

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
import ast
import json
import random
from json.decoder import JSONDecodeError
from typing import Union

import fire
import jsonpickle
from lte.protos.keyval_pb2 import IPDesc
from lte.protos.oai.mme_nas_state_pb2 import MmeNasState, UeContext
from lte.protos.oai.s1ap_state_pb2 import S1apImsiMap, S1apState, UeDescription
from lte.protos.oai.spgw_state_pb2 import SpgwState, SpgwUeContext
from lte.protos.policydb_pb2 import InstalledPolicies, PolicyRule
from magma.common.redis.client import get_default_client
from magma.common.redis.serializers import (
    get_json_deserializer,
    get_proto_deserializer,
    get_proto_version_deserializer,
)
from magma.mobilityd.serialize_utils import (
    deserialize_ip_block,
    deserialize_ip_desc,
)

NO_DESERIAL_MSG = "No deserializer exists for type '{}'"


def _deserialize_session_json(serialized_json_str: bytes) -> str:
    """
    Helper function to deserialize sessiond:sessions hash list values
    :param serialized_json_str
    """
    res = _deserialize_generic_json(str(serialized_json_str, 'utf-8', 'ignore'))
    dumped = json.dumps(res, indent=2, sort_keys=True)
    return dumped


def _deserialize_generic_json(
        element: Union[str, dict, list],
) -> Union[str, dict, list]:
    """
    Helper function to deserialize dictionaries or list with nested
    json strings
    :param element
    """
    if isinstance(element, str):
        # try to deserialize as json string
        try:
            element = ast.literal_eval(element)
        except:
            try:
                element = jsonpickle.decode(element)
            except:
                return element

    if isinstance(element, dict):
        keys = element.keys()
    elif isinstance(element, list):
        keys = range(len(element))
    else:
        # in case it is neither of the know elements, just return as is
        return element

    for k in keys:
        element[k] = _deserialize_generic_json(element[k])
    return element


class StateCLI(object):
    """
    CLI for debugging current Magma services state and displaying it
    in readable manner.
    """

    STATE_DESERIALIZERS = {
        'assigned_ip_blocks': deserialize_ip_block,
        'ip_states': deserialize_ip_desc,
        'sessions': _deserialize_session_json,
        'rule_names': get_json_deserializer(),
        'rule_ids': get_json_deserializer(),
        'rule_versions': get_json_deserializer(),
        'rules': get_proto_deserializer(PolicyRule),
    }

    STATE_PROTOS = {
        'mme_nas_state': MmeNasState,
        'spgw_state': SpgwState,
        's1ap_state': S1apState,
        's1ap_imsi_map': S1apImsiMap,
        'mme': UeContext,
        'spgw': SpgwUeContext,
        's1ap': UeDescription,
        'mobilityd_ipdesc_record': IPDesc,
        'rules': PolicyRule,
        'installed': InstalledPolicies,
    }

    def __init__(self):
        self.client = get_default_client()

    def keys(self, redis_key: str):
        """
        Get current keys on redis db that match the pattern

        Args:
            redis_key:pattern to match the redis keys

        """
        for k in self.client.keys(pattern="{}*".format(redis_key)):
            deserialized_key = k.decode('utf-8')
            print(deserialized_key)

    def parse(self, key: str):
        """
        Parse value of redis key on redis for encoded HASH, SET types, or
        JSON / Protobuf encoded state-wrapped types and prints it

        Args:
            key: key on redis

        """
        redis_type = self.client.type(key).decode('utf-8')
        key_type = key
        if ":" in key:
            key_type = key.split(":")[1]
        if redis_type == 'hash':
            deserializer = self.STATE_DESERIALIZERS.get(key_type)
            if not deserializer:
                raise AttributeError(NO_DESERIAL_MSG.format(key_type))
            self._parse_hash_type(deserializer, key)
        elif redis_type == 'set':
            deserializer = self.STATE_DESERIALIZERS.get(key_type)
            if not deserializer:
                raise AttributeError(NO_DESERIAL_MSG.format(key_type))
            self._parse_set_type(deserializer, key)
        else:
            value = self.client.get(key)
            # Try parsing as json first, if there's decoding error, parse proto
            try:
                self._parse_state_json(value)
            except (UnicodeDecodeError, JSONDecodeError, AttributeError):
                self._parse_state_proto(key_type, value)

    def corrupt(self, key):
        """
        Mostly used for debugging, purposely corrupts state encoded protobuf
        in redis, and writes it back to datastore

        Args:
            key: key on redis

        """
        rand_bytes = random.getrandbits(8)
        byte_str = bytes([rand_bytes])
        self.client[key] = byte_str
        print('Corrupted %s in redis' % key)

    def _parse_state_json(self, value):
        if value:
            deserializer = get_json_deserializer()
            value = json.loads(jsonpickle.encode(deserializer(value)))
            print(json.dumps(value, indent=2, sort_keys=True))
        else:
            raise AttributeError('Key not found on redis')

    def _parse_state_proto(self, key_type, value):
        proto = self.STATE_PROTOS.get(key_type.lower())
        if proto:
            deserializer = get_proto_deserializer(proto)
            version_deserializer = get_proto_version_deserializer()
            print(deserializer(value))
            print('==================')
            print('State version: %s' % version_deserializer(value))
        else:
            raise AttributeError('Key not found on redis')

    def _parse_set_type(self, deserializer, key):
        set_values = self.client.smembers(key)
        for value in set_values:
            print(deserializer(value))

    def _parse_hash_type(self, deserializer, key):
        value = self.client.hgetall(key)
        for key, val in value.items():
            print(key.decode('utf-8'))
            print(deserializer(val))


if __name__ == "__main__":
    state_cli = StateCLI()
    try:
        fire.Fire(state_cli)
    except Exception as e:
        print('Error: {}'.format(e))
