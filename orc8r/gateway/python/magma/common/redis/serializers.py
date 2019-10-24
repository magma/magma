"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from orc8r.protos.redis_pb2 import RedisState

def get_proto_serializer():
    """
    Return a proto serializer that serializes the proto, adds the associated
    version, and then serializes the RedisState proto to a string
    """
    def _serialize_proto(proto, version):
        serialized_proto = proto.SerializeToString()
        redis_state = RedisState(
            serialized_msg=serialized_proto,
            version=version)
        return redis_state.SerializeToString()
    return _serialize_proto


def get_proto_deserializer(proto_class):
    """
    Return a proto deserializer that takes in a proto type to deserialize
    the serialized msg stored in the RedisState proto
    """
    def _deserialize_proto(serialized_rule):
        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(serialized_rule)
        serialized_proto = proto_wrapper.serialized_msg
        proto = proto_class()
        proto.ParseFromString(serialized_proto)
        return proto
    return _deserialize_proto
