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
from typing import Callable, Generic, Type, TypeVar

import jsonpickle
from orc8r.protos.redis_pb2 import RedisState

T = TypeVar('T')


class RedisSerde(Generic[T]):
    """
    typeval (str): str representing the type of object the serde can
        de/serialize
    serializer (function (T, int) -> str):
                function called to serialize a value
    deserializer (function (str) -> T):
                function called to deserialize a value
    """

    def __init__(
        self,
        redis_type: str,
        serializer: Callable[[T, int], str],
        deserializer: Callable[[str], T],
    ):
        self.redis_type = redis_type
        self.serializer = serializer
        self.deserializer = deserializer

    def serialize(self, msg: T, version: int = 1) -> str:
        return self.serializer(msg, version)

    def deserialize(self, serialized_obj: str) -> T:
        return self.deserializer(serialized_obj)


def get_proto_serializer() -> Callable[[T, int], str]:
    """
    Return a proto serializer that serializes the proto, adds the associated
    version, and then serializes the RedisState proto to a string
    """
    def _serialize_proto(proto: T, version: int) -> str:
        serialized_proto = proto.SerializeToString()
        redis_state = RedisState(
            serialized_msg=serialized_proto,
            version=version,
            is_garbage=False,
        )
        return redis_state.SerializeToString()
    return _serialize_proto


def get_proto_deserializer(proto_class: Type[T]) -> Callable[[str], T]:
    """
    Return a proto deserializer that takes in a proto type to deserialize
    the serialized msg stored in the RedisState proto
    """
    def _deserialize_proto(serialized_rule: str) -> T:
        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(serialized_rule)
        serialized_proto = proto_wrapper.serialized_msg
        proto = proto_class()
        proto.ParseFromString(serialized_proto)
        return proto
    return _deserialize_proto


def get_json_serializer() -> Callable[[T, int], str]:
    """
    Return a json serializer that serializes the json msg, adds the
    associated version, and then serializes the RedisState proto to a string
    """
    def _serialize_json(msg: T, version: int) -> str:
        serialized_msg = jsonpickle.encode(msg)
        redis_state = RedisState(
            serialized_msg=serialized_msg.encode('utf-8'),
            version=version,
            is_garbage=False,
        )
        return redis_state.SerializeToString()

    return _serialize_json


def get_json_deserializer() -> Callable[[str], T]:
    """
    Returns a json deserializer that deserializes the RedisState proto and
    then deserializes the json msg
    """
    def _deserialize_json(serialized_rule: str) -> T:
        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(serialized_rule)
        serialized_msg = proto_wrapper.serialized_msg
        msg = jsonpickle.decode(serialized_msg.decode('utf-8'))
        return msg

    return _deserialize_json


def get_proto_version_deserializer() -> Callable[[str], T]:
    """
    Return a proto deserializer that takes in a proto type to deserialize
    the version number stored in the RedisState proto
    """
    def _deserialize_version(serialized_rule: str) -> T:
        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(serialized_rule)
        return proto_wrapper.version
    return _deserialize_version
