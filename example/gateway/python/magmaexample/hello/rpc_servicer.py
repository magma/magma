"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from protos import hello_pb2, hello_pb2_grpc

from magmaexample.hello import metrics


class HelloRpcServicer(hello_pb2_grpc.HelloServicer):
    """
    gRPC based server for Hello service
    """

    def __init__(self):
        pass

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        hello_pb2_grpc.add_HelloServicer_to_server(self, server)

    def SayHello(self, request, context):
        """
        Echo the message as the response
        """
        metrics.NUM_REQUESTS.inc()
        return hello_pb2.HelloReply(greeting=request.greeting)
