#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse
import sys

import grpc
from grpc import StatusCode
from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos.directoryd_pb2 import GetLocationRequest, \
    UpdateDirectoryLocationRequest
from orc8r.protos.directoryd_pb2_grpc import DirectoryServiceStub


@grpc_wrapper
def get_location_handler(client, args):
    get_location_request_msg = GetLocationRequest()
    get_location_request_msg.table = int(args.table)
    get_location_request_msg.id = args.object_id
    try:
        location_record_msg = client.GetLocation(get_location_request_msg)
        print("%s => %s" % (args.object_id, location_record_msg.location))
    except grpc.RpcError as e:
        if e.code() == StatusCode.UNKNOWN:
            print("Object ID '%s' not found" % args.object_id)
        else:
            print("gRPC failed with %s: %s" % (e.code(), e.details()))


@grpc_wrapper
def update_location_handler(client, args):
    update_location_request_msg = UpdateDirectoryLocationRequest()
    update_location_request_msg.table = int(args.table)
    update_location_request_msg.id = args.object_id
    update_location_request_msg.record.location = args.location
    client.UpdateLocation(update_location_request_msg)
    print("Location Updated: %s => %s" % (args.object_id, args.location))

@grpc_wrapper
def delete_location_handler(client, args):
    delete_location_request_msg = GetLocationRequest()
    delete_location_request_msg.table = int(args.table)
    delete_location_request_msg.id = args.object_id
    try:
        client.DeleteLocation(delete_location_request_msg)
        print("Deleted %s" % args.object_id)
    except grpc.RpcError as e:
        if e.code() == StatusCode.UNKNOWN:
            print("Object ID '%s' not found" % args.object_id)
        else:
            print("gRPC failed with %s: %s" % (e.code(), e.details()))


def main():
    parser = argparse.ArgumentParser(
        description='Management CLI for DirectoryService',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # get_location
    subparser = subparsers.add_parser(
        'get_location', help='Get location of an object')
    subparser.add_argument('table',
        help='table ID. 0 = maps IMSI to HwId; 1 = maps HwId to HostName')
    subparser.add_argument('object_id', help='ID of the object')
    subparser.set_defaults(func=get_location_handler)

    # update_location
    subparser = subparsers.add_parser(
        'update_location', help='Get location of an object')
    subparser.add_argument('table',
        help='table ID. 0 = maps IMSI to HwId; 1 = maps HwId to HostName')
    subparser.add_argument('object_id', help='ID of the object')
    subparser.add_argument('location', help='Location of an object')
    subparser.set_defaults(func=update_location_handler)

    # get_location
    subparser = subparsers.add_parser(
        'delete_location', help='Delete location of an object')
    subparser.add_argument('table',
        help='table ID. 0 = maps IMSI to HwId; 1 = maps HwId to HostName')
    subparser.add_argument('object_id', help='ID of the object')
    subparser.set_defaults(func=delete_location_handler)

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        sys.exit(1)

    # Execute the subcommand function
    args.func(args, DirectoryServiceStub, 'directoryd')


if __name__ == "__main__":
    main()
