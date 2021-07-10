#!/usr/bin/env python3

"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import sys

import grpc
from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos.common_pb2 import Void
from orc8r.protos.directoryd_pb2 import (
    DeleteRecordRequest,
    GetDirectoryFieldRequest,
    UpdateRecordRequest,
)
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub


@grpc_wrapper
def update_record_handler(client, args):
    update_record_request = UpdateRecordRequest()
    update_record_request.id = args.id
    if args.field_key is not None and args.field_value is not None:
        update_record_request.fields[args.field_key] = args.field_value

    try:
        client.UpdateRecord(update_record_request)
        print("Successfully updated record for ID: %s" % args.id)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


@grpc_wrapper
def delete_record_handler(client, args):
    delete_record_request = DeleteRecordRequest()
    delete_record_request.id = args.id
    try:
        client.DeleteRecord(delete_record_request)
        print("Successfully deleted record for %s" % args.id)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


@grpc_wrapper
def get_record_field_handler(client, args):
    get_request = GetDirectoryFieldRequest()
    get_request.id = args.id
    get_request.field_key = args.field_key

    try:
        res = client.GetDirectoryField(get_request)
        print(
            "Successfully got field: (%s, %s) for ID: %s" % (
                res.key,
                res.value,
                args.id,
            ),
        )
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


@grpc_wrapper
def get_all_records_handler(client, args):
    void_request = Void()
    try:
        res = client.GetAllDirectoryRecords(void_request)
        for record in res.records:
            print(record)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def main():
    parser = argparse.ArgumentParser(
        description='Management CLI for DirectoryService',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # update_record
    subparser = subparsers.add_parser(
        'update_record', help='Update record of an object',
    )
    subparser.add_argument('id', help='ID')
    subparser.add_argument('--field_key', required=False)
    subparser.add_argument('--field_value', required=False)
    subparser.set_defaults(func=update_record_handler)

    # delete_record
    subparser = subparsers.add_parser(
        'delete_record', help='Delete record of an object',
    )
    subparser.add_argument('id', help='ID')
    subparser.set_defaults(func=delete_record_handler)

    # get_record_field
    subparser = subparsers.add_parser(
        'get_record_field', help='Get field of a record object',
    )
    subparser.add_argument('id', help='ID')
    subparser.add_argument('field_key', help='Field key to lookup')
    subparser.set_defaults(func=get_record_field_handler)

    # get_all_records
    subparser = subparsers.add_parser(
        'get_all_records', help='Get all records',
    )
    subparser.set_defaults(func=get_all_records_handler)

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        sys.exit(1)

    # Execute the subcommand function
    args.func(args, GatewayDirectoryServiceStub, 'directoryd')


if __name__ == "__main__":
    main()
