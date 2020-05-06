"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.

Describes protobuf-based serialize and deserialize functions for mobilityd.
"""

from ipaddress import ip_address, ip_network

import magma.mobilityd.ip_descriptor as ip_descriptor
from lte.protos.keyval_pb2 import IPDesc
from lte.protos.mobilityd_pb2 import IPAddress, IPBlock
from lte.protos.subscriberdb_pb2 import SubscriberID
from orc8r.protos.redis_pb2 import RedisState


def _ip_version_int_to_proto(version):
    proto = {4: IPBlock.IPV4, 6: IPBlock.IPV6}[version]
    return proto


def _desc_state_str_to_proto(state):
    proto_map = {
        ip_descriptor.IPState.FREE: IPDesc.FREE,
        ip_descriptor.IPState.ALLOCATED: IPDesc.ALLOCATED,
        ip_descriptor.IPState.RELEASED: IPDesc.RELEASED,
        ip_descriptor.IPState.REAPED: IPDesc.REAPED,
        ip_descriptor.IPState.RESERVED: IPDesc.RESERVED,
    }
    proto = proto_map[state]
    return proto


def _desc_state_proto_to_str(proto):
    state_map = {
        IPDesc.FREE: ip_descriptor.IPState.FREE,
        IPDesc.ALLOCATED: ip_descriptor.IPState.ALLOCATED,
        IPDesc.RELEASED: ip_descriptor.IPState.RELEASED,
        IPDesc.REAPED: ip_descriptor.IPState.REAPED,
        IPDesc.RESERVED: ip_descriptor.IPState.RESERVED,
    }
    state = state_map[proto]
    return state


def _ip_desc_to_proto(desc):
    """
    Convert an IP descriptor to protobuf.

    Args:
        desc (magma.mobilityd.IPDesc): IP descriptor
    Returns:
        proto (protos.keyval_pb2.IPDesc): protobuf of :desc:
    """
    ip = IPAddress(
        version=_ip_version_int_to_proto(desc.ip_block.version),
        address=desc.ip.packed,
    )
    ip_block = IPBlock(
        version=_ip_version_int_to_proto(desc.ip_block.version),
        net_address=desc.ip_block.network_address.packed,
        prefix_len=desc.ip_block.prefixlen,
    )
    state = _desc_state_str_to_proto(desc.state)
    sid = SubscriberID(
        id=desc.sid,
        type=SubscriberID.IMSI,
    )
    proto = IPDesc(ip=ip, ip_block=ip_block, state=state, sid=sid)
    return proto


def _ip_desc_from_proto(proto):
    """
    Convert protobuf to an IP descriptor.

    Args:
        proto (protos.keyval_pb2.IPDesc): protobuf of an IP descriptor
    Returns:
        desc (magma.mobilityd.IPDesc): IP descriptor from :proto:
    """
    ip = ip_address(proto.ip.address)
    ip_block_addr = ip_address(proto.ip_block.net_address).exploded
    ip_block = ip_network('{}/{}'.format(
        ip_block_addr, proto.ip_block.prefix_len))
    state = _desc_state_proto_to_str(proto.state)
    sid = proto.sid.id
    desc = ip_descriptor.IPDesc(ip=ip, ip_block=ip_block, state=state, sid=sid)
    return desc


def serialize_ip_block(block):
    """
    Serialize an IP block to protobuf string.

    Args:
        block (ipaddress.ip_network): object to serialize
    Returns:
        serialized (bytes): serialized object
    """
    proto = IPBlock(
        version=_ip_version_int_to_proto(block.version),
        net_address=block.network_address.packed,
        prefix_len=block.prefixlen,
    )
    serialized = proto.SerializeToString()
    return serialized


def deserialize_ip_block(serialized):
    """
    Deserialize protobuf string to an IP block.

    Args:
        serialized (bytes): object to deserialize
    Returns:
        block (ipaddress.ip_network): deserialized object
    """
    proto = IPBlock()
    proto.ParseFromString(serialized)
    address_str = ip_address(proto.net_address).exploded
    network_str = '{}/{}'.format(address_str, proto.prefix_len)
    block = ip_network(network_str)
    return block


def serialize_ip_desc(desc, version):
    """
    Serialize an IP descriptor to protobuf string.

    Args:
        desc (magma.mobilityd.IPDesc): object to serialize
    Returns:
        serialized (bytes): serialized object
    """
    proto = _ip_desc_to_proto(desc)
    serialized = proto.SerializeToString()
    redis_state = RedisState(
        serialized_msg=serialized,
        version=version)
    return redis_state.SerializeToString()


def deserialize_ip_desc(serialized):
    """
    Deserialize protobuf string to an IP descriptor.

    Args:
        serialized (bytes): object to deserialize
    Returns:
        block (magma.mobilityd.IPDesc): deserialized object
    """
    proto_wrapper = RedisState()
    proto_wrapper.ParseFromString(serialized)
    serialized_proto = proto_wrapper.serialized_msg
    proto = IPDesc()
    proto.ParseFromString(serialized_proto)
    desc = _ip_desc_from_proto(proto)
    return desc
