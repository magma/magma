"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=broad-except

import asyncio
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


def _grpc_async_wrapper(f, gf):
    try:
        f.set_result(gf.result())
    except Exception as e:
        f.set_exception(e)


def grpc_async_wrapper(gf, loop=None):
    """
    Wraps a GRPC result in a future that can be yielded by asyncio

    Usage:

    async def my_fn(param):
        result =
            await grpc_async_wrapper(stub.function_name.future(param, timeout))

    Code taken and modified from:
        https://github.com/grpc/grpc/wiki/Integration-with-tornado-(python)
    """
    f = asyncio.Future()
    if loop is None:
        loop = asyncio.get_event_loop()
    gf.add_done_callback(
        lambda _: loop.call_soon_threadsafe(_grpc_async_wrapper, f, gf)
    )
    return f
