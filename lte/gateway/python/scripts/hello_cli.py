#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.common.rpc_utils import grpc_wrapper
from feg.protos.hello_pb2 import HelloRequest
from feg.protos.hello_pb2_grpc import HelloStub


@grpc_wrapper
def say_hello(client, args):
    print(client.SayHello(HelloRequest(greeting='world')))


if __name__ == "__main__":
    args = None
    say_hello(args, HelloStub, 'hello')
