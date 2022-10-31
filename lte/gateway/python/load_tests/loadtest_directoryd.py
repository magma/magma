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
import json
from typing import List

from google.protobuf import json_format
from load_tests.common import (
    benchmark_grpc_request,
    make_full_request_type,
    make_output_file_path,
)
from magma.common.service_registry import ServiceRegistry
from orc8r.protos.common_pb2 import Void
from orc8r.protos.directoryd_pb2 import (
    DeleteRecordRequest,
    DirectoryRecord,
    GetDirectoryFieldRequest,
    UpdateRecordRequest,
)
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub

DIRECTORYD_SERVICE_NAME = 'directoryd'
DIRECTORYD_SERVICE_RPC_PATH = 'magma.orc8r.GatewayDirectoryService'
DIRECTORYD_PORT = '127.0.0.1:50067'
PROTO_PATH = 'orc8r/protos/directoryd.proto'


def _load_subs(num_subs: int) -> List[DirectoryRecord]:
    """Load directory records"""
    client = GatewayDirectoryServiceStub(
        ServiceRegistry.get_rpc_channel(
            DIRECTORYD_SERVICE_NAME, ServiceRegistry.LOCAL,
        ),
    )
    sids = []
    for i in range(num_subs):
        mac_addr = (str(i) * 2 + ":") * 5 + (str(i) * 2)
        ipv4_addr = str(i) * 3 + "." + str(i) * 3 + "." + str(i) * 3 + "." + str(i) * 3
        fields = {"mac-addr": mac_addr, "ipv4_addr": ipv4_addr}
        sid = UpdateRecordRequest(
            fields=fields,
            id=str(i).zfill(15),
            location=str(i).zfill(15),
        )
        client.UpdateRecord(sid)
        sids.append(sid)
    return sids


def _cleanup_subs():
    """Clear directory records"""
    client = GatewayDirectoryServiceStub(
        ServiceRegistry.get_rpc_channel(
            DIRECTORYD_SERVICE_NAME, ServiceRegistry.LOCAL,
        ),
    )
    for record in client.GetAllDirectoryRecords(Void()).records:
        sid = DeleteRecordRequest(
            id=record.id,
        )
        client.DeleteRecord(sid)


def _build_update_records_data(num_requests: int, input_file: str):
    update_record_reqs = []
    for i in range(num_requests):
        id = str(i).zfill(15)
        location = str(i).zfill(15)
        request = UpdateRecordRequest(
            id=id,
            location=location,
        )
        request_dict = json_format.MessageToDict(request)
        update_record_reqs.append(request_dict)
    with open(input_file, 'w') as file:
        json.dump(update_record_reqs, file, separators=(',', ':'))


def _build_delete_records_data(record_list: list, input_file: str):
    delete_record_reqs = []
    for index, record in enumerate(record_list):
        request = DeleteRecordRequest(
            id=record.id,
        )
        request_dict = json_format.MessageToDict(request)
        delete_record_reqs.append(request_dict)
    with open(input_file, 'w') as file:
        json.dump(delete_record_reqs, file, separators=(',', ':'))


def _build_get_record_data(record_list: list, input_file: str):
    get_record_reqs = []
    for index, record in enumerate(record_list):
        request = GetDirectoryFieldRequest(
            id=record.id,
            field_key="mac-addr",
        )
        request_dict = json_format.MessageToDict(request)
        get_record_reqs.append(request_dict)
    with open(input_file, 'w') as file:
        json.dump(get_record_reqs, file, separators=(',', ':'))


def _build_get_all_record_data(record_list: list, input_file: str):
    request = Void()
    get_all_record_reqs = json_format.MessageToDict(request)
    with open(input_file, 'w') as file:
        json.dump(get_all_record_reqs, file, separators=(',', ':'))


def update_record_test(args):
    input_file = 'update_record.json'
    _build_update_records_data(args.num_of_requests, input_file)
    request_type = 'UpdateRecord'
    benchmark_grpc_request(
        proto_path=PROTO_PATH,
        full_request_type=make_full_request_type(
            DIRECTORYD_SERVICE_RPC_PATH, request_type,
        ),
        input_file=input_file,
        output_file=make_output_file_path(request_type),
        num_reqs=args.num_of_requests, address=DIRECTORYD_PORT,
        import_path=args.import_path,
    )
    _cleanup_subs()


def delete_record_test(args):
    input_file = 'delete_record.json'
    record_list = _load_subs(args.num_of_requests)
    _build_delete_records_data(record_list, input_file)

    request_type = 'DeleteRecord'
    benchmark_grpc_request(
        proto_path=PROTO_PATH,
        full_request_type=make_full_request_type(
            DIRECTORYD_SERVICE_RPC_PATH, request_type,
        ),
        input_file=input_file,
        output_file=make_output_file_path(request_type),
        num_reqs=args.num_of_requests, address=DIRECTORYD_PORT,
        import_path=args.import_path,
    )
    _cleanup_subs()


def get_record_test(args):
    input_file = 'get_record.json'
    record_list = _load_subs(args.num_of_requests)
    _build_get_record_data(record_list, input_file)
    request_type = 'GetDirectoryField'
    benchmark_grpc_request(
        proto_path=PROTO_PATH,
        full_request_type=make_full_request_type(
            DIRECTORYD_SERVICE_RPC_PATH, request_type,
        ),
        input_file=input_file,
        output_file=make_output_file_path(request_type),
        num_reqs=args.num_of_requests, address=DIRECTORYD_PORT,
        import_path=args.import_path,
    )
    _cleanup_subs()


def get_all_records_test(args):
    input_file = 'get_all_records.json'
    record_list = _load_subs(args.num_of_requests)
    _build_get_all_record_data(record_list, input_file)
    request_type = 'GetAllDirectoryRecords'
    benchmark_grpc_request(
        proto_path=PROTO_PATH,
        full_request_type=make_full_request_type(
            DIRECTORYD_SERVICE_RPC_PATH, request_type,
        ),
        input_file=input_file,
        output_file=make_output_file_path(request_type),
        num_reqs=2000, address=DIRECTORYD_PORT,
        import_path=args.import_path,
    )
    _cleanup_subs()


def create_parser():
    """
    Creates the argparse subparser for all args
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    parser_update_record = subparsers.add_parser(
        'update_record', help='Update record in directory',
    )

    parser_delete_record = subparsers.add_parser(
        'delete_record', help='Delete record in directory',
    )

    parser_get_record = subparsers.add_parser(
        'get_record', help='Get specific record in directory',
    )

    parser_get_all_records = subparsers.add_parser(
        'get_all_records', help='Get all records in directory',
    )

    for subcmd in [
        parser_update_record,
        parser_delete_record,
        parser_get_record,
        parser_get_all_records,
    ]:

        subcmd.add_argument(
            '--num_of_requests', help='Number of total records in directory',
            type=int, default=2000,
        )
        subcmd.add_argument(
            '--import_path', default=None, help='Protobuf import path directory',
        )

    parser_update_record.set_defaults(func=update_record_test)
    parser_delete_record.set_defaults(func=delete_record_test)
    parser_get_record.set_defaults(func=get_record_test)
    parser_get_all_records.set_defaults(func=get_all_records_test)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args)


if __name__ == "__main__":
    main()
