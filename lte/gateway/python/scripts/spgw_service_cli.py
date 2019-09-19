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
from lte.protos.spgw_service_pb2 import CreateBearerRequest, \
    DeleteBearerRequest
from lte.protos.spgw_service_pb2_grpc import SpgwServiceStub


@grpc_wrapper
def create_bearer(client, args):
    req = CreateBearerRequest()
    # TODO populate request
    print("Creating bearer: ")
    client.CreateBearer(req)


@grpc_wrapper
def delete_bearer(client, args):
    req = DeleteBearerRequest()
    # TODO populate request
    print("Deleting bearer: ")
    client.DeleteBearer(req)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for SPGW service',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_create_bearer = subparsers.add_parser('create',
                                                 help='Create Bearer')
    parser_delete_bearer = subparsers.add_parser('delete',
                                                 help='Delete Bearer')

    # Add function callbacks
    parser_create_bearer.set_defaults(func=create_bearer)
    parser_delete_bearer.set_defaults(func=delete_bearer)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, SpgwServiceStub, 'spgw_service')


if __name__ == "__main__":
    main()
