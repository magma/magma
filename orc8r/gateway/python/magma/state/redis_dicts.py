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

import importlib
import logging

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisFlatDict
from magma.common.redis.serializers import (
    RedisSerde,
    get_json_deserializer,
    get_json_serializer,
    get_proto_deserializer,
    get_proto_serializer,
)

PROTO_FORMAT = 0
JSON_FORMAT = 1


class StateDict(RedisFlatDict):
    """
    StateDict is a RedisFlatDict that holds state metadata and reads/writes
    state to Redis.
    """

    def __init__(self, serde: RedisSerde, state_scope: str, state_format: int):
        super().__init__(get_default_client(), serde)
        # Scope determines the deviceID to report the state with
        self.state_scope = state_scope
        self.state_format = state_format


def get_proto_redis_dicts(config):
    redis_dicts = []
    state_protos = config.get('state_protos', []) or []
    for proto_cfg in state_protos:
        is_invalid_cfg = 'proto_msg' not in proto_cfg or \
                         'proto_file' not in proto_cfg or \
                         'redis_key' not in proto_cfg or \
                         'state_scope' not in proto_cfg
        if is_invalid_cfg:
            logging.warning(
                "Invalid proto config found in state_protos "
                "configuration: %s", proto_cfg,
            )
            continue
        try:
            proto_module = importlib.import_module(proto_cfg['proto_file'])
            msg = getattr(proto_module, proto_cfg['proto_msg'])
            redis_key = proto_cfg['redis_key']
            logging.info(
                'Initializing RedisSerde for proto state %s',
                proto_cfg['redis_key'],
            )
            serde = RedisSerde(
                redis_key,
                get_proto_serializer(),
                get_proto_deserializer(msg),
            )
            redis_dict = StateDict(
                serde,
                proto_cfg['state_scope'],
                PROTO_FORMAT,
            )
            redis_dicts.append(redis_dict)

        except (ImportError, AttributeError) as err:
            logging.error(err)

    return redis_dicts


def get_json_redis_dicts(config):
    redis_dicts = []
    json_state = config.get('json_state', []) or []
    for json_cfg in json_state:
        is_invalid_cfg = 'redis_key' not in json_cfg or \
                         'state_scope' not in json_cfg
        if is_invalid_cfg:
            logging.warning(
                "Invalid json state config found in json_state"
                "configuration: %s", json_cfg,
            )
            continue

        logging.info(
            'Initializing RedisSerde for json state %s',
            json_cfg['redis_key'],
        )
        redis_key = json_cfg['redis_key']
        serde = RedisSerde(
            redis_key,
            get_json_serializer(),
            get_json_deserializer(),
        )
        redis_dict = StateDict(
            serde,
            json_cfg['state_scope'],
            JSON_FORMAT,
        )
        redis_dicts.append(redis_dict)

    return redis_dicts
