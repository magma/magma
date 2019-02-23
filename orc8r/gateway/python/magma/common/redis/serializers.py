"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


def get_proto_serializer():
    """
    Return a basic proto serializer that simply serializes to a string with no
    extra field manipulation
    """
    def _serialize_proto(proto):
        return proto.SerializeToString()
    return _serialize_proto


def get_proto_deserializer(proto_class):
    """
    Return a basic proto deserializer that takes in a proto type to deserialize
    to
    """
    def _deserialize_proto(serialized_rule):
        proto = proto_class()
        proto.ParseFromString(serialized_rule)
        return proto
    return _deserialize_proto
