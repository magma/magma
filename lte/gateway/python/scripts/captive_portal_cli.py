#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse
from magma.common.rpc_utils import grpc_wrapper
from lte.protos.session_manager_pb2 import CreateSessionRequest, \
    UpdateSessionRequest, SessionTerminateRequest
from lte.protos.session_manager_pb2_grpc import CentralSessionControllerStub


@grpc_wrapper
def create_session(client, args):
    message = CreateSessionRequest()
    client.CreateSession(message)
    print('Successfully got a response from CreateSession')


@grpc_wrapper
def update_session(client, args):
    message = UpdateSessionRequest()
    client.UpdateSession(message)
    print('Successfully got a response from UpdateSession')


@grpc_wrapper
def terminate_session(client, args):
    message = SessionTerminateRequest()
    client.TerminateSession(message)
    print('Successfully got a response from TerminateSession')


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for Captive Portal',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    parser_create_session = subparsers.add_parser(
        'create_session', help='Send CreateSessionRequest message')

    parser_update_session = subparsers.add_parser(
        'update_session', help='Send UpdateSessionRequest message')

    parser_terminate_session = subparsers.add_parser(
        'terminate_session', help='Send SessionTerminateRequest message')

    # Add function callbacks
    parser_create_session.set_defaults(func=create_session)
    parser_update_session.set_defaults(func=update_session)
    parser_terminate_session.set_defaults(func=terminate_session)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, CentralSessionControllerStub, 'captive_portal')


if __name__ == "__main__":
    main()
