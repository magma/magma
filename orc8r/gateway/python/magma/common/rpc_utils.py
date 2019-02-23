"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import grpc
from orc8r.protos import common_pb2

from .service_registry import ServiceRegistry


def return_void(func):
    """
    Reusable decorator for returning common_pb2.Void() message.
    """
    def wrapper(*args, **kwargs):
        func(*args, **kwargs)
        return common_pb2.Void()
    return wrapper


def grpc_wrapper(func):
    """
    Wraps a function with a gRPC wrapper which creates a RPC client to
    the service and handles any RPC Exceptions.

    Usage:
    @grpc_wrapper
    def func(client, args):
        pass
    func(args, ProtoStubClass, 'service')
    """
    def wrapper(*alist):
        args = alist[0]
        stub_cls = alist[1]
        service = alist[2]
        chan = ServiceRegistry.get_rpc_channel(service, ServiceRegistry.LOCAL)
        client = stub_cls(chan)
        try:
            func(client, args)
        except grpc.RpcError as err:
            print("Error! [%s] %s" % (err.code(), err.details()))
            exit(1)
    return wrapper


def cloud_grpc_wrapper(func):
    """
    Wraps a function with a gRPC wrapper which creates a RPC client to
    the service and handles any RPC Exceptions.

    Usage:
    @cloud_grpc_wrapper
    def func(client, args):
        pass
    func(args, ProtoStubClass, 'service')
    """
    def wrapper(*alist):
        args = alist[0]
        stub_cls = alist[1]
        service = alist[2]
        chan = ServiceRegistry.get_rpc_channel(service, ServiceRegistry.CLOUD)
        client = stub_cls(chan)
        try:
            func(client, args)
        except grpc.RpcError as err:
            print("Error! [%s] %s" % (err.code(), err.details()))
            exit(1)
    return wrapper
