#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse
import grpc

from orc8r.protos.common_pb2 import Void
from lte.protos.subscriberdb_pb2 import SubscriberID
from lte.protos.session_manager_pb2 import LocalCreateSessionRequest
from lte.protos.session_manager_pb2_grpc import LocalSessionManagerStub
from feg.protos.mock_core_pb2_grpc import MockPCRFStub, MockOCSStub
from magma.pipelined.tests.app.subscriber import SubContextConfig
from magma.common.rpc_utils import grpc_wrapper
from magma.common.service_registry import ServiceRegistry


@grpc_wrapper
def send_create_session(client, args):
    sub1 = SubContextConfig('IMSI' + args.imsi, '192.168.128.74', 4)

    try:
        create_account_in_PCRF(args.imsi)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))

    try:
        create_account_in_OCS(args.imsi)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))

    req = LocalCreateSessionRequest(
        sid=SubscriberID(id=sub1.imsi),
        ue_ipv4=sub1.ip,
    )
    print("Sending LocalCreateSessionRequest with following fields:\n %s" % req)
    try:
        client.CreateSession(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))

    req = SubscriberID(id=sub1.imsi)
    print("Sending EndSession with following fields:\n %s" % req)
    try:
        client.EndSession(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def create_account_in_PCRF(imsi):
    pcrf_chan = ServiceRegistry.get_rpc_channel('pcrf', ServiceRegistry.CLOUD)
    pcrf_client = MockPCRFStub(pcrf_chan)

    print("Clearing accounts in PCRF")
    pcrf_client.ClearSubscribers(Void())

    print("Creating account in PCRF")
    pcrf_client.CreateAccount(SubscriberID(id=imsi))


def create_account_in_OCS(imsi):
    ocs_chan = ServiceRegistry.get_rpc_channel('ocs', ServiceRegistry.CLOUD)
    ocs_client = MockOCSStub(ocs_chan)

    print("Clearing accounts in OCS")
    ocs_client.ClearSubscribers(Void())

    print("Creating account in OCS")
    ocs_client.CreateAccount(SubscriberID(id=imsi))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for testing session manager',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Create Session
    create_session_parser = subparsers.add_parser(
        'create_session',
        help='Send Create Session Request to session_proxy service in FeG',
    )
    create_session_parser.add_argument('imsi', help='e.g.001010000088888')
    create_session_parser.set_defaults(func=send_create_session)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, LocalSessionManagerStub, 'sessiond')


if __name__ == "__main__":
    main()
