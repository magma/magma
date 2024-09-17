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
    PROTO_DIR,
    benchmark_grpc_request,
    make_full_request_type,
    make_output_file_path,
)
from lte.protos.apn_pb2 import APNConfiguration
from lte.protos.subscriberdb_pb2 import (
    Non3GPPUserProfile,
    SubscriberData,
    SubscriberID,
    SubscriberUpdate,
)
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.sid import SIDUtils
from orc8r.protos.common_pb2 import Void

TEST_APN = 'magma.ipv4'
TEST_APN_UPDATE = 'magma.ipv6'
SUBSCRIBERDB_SERVICE_RPC_PATH = 'magma.lte.SubscriberDB'
SUBSCRIBERDB_SERVICE_NAME = 'subscriberdb'
SUBSCRIBERDB_PORT = '0.0.0.0:50051'
PROTO_PATH = PROTO_DIR + '/subscriberdb.proto'


# Helper functions to build input data for gRPC functions

def _load_subs(num_subs: int) -> List[SubscriberID]:
    client = SubscriberDBStub(
        ServiceRegistry.get_rpc_channel(
            SUBSCRIBERDB_SERVICE_NAME, ServiceRegistry.LOCAL,
        ),
    )
    sids = []

    for i in range(1, num_subs):
        sid = SubscriberID(id=str(i).zfill(15))
        config = Non3GPPUserProfile(
            apn_config=[APNConfiguration(service_selection=TEST_APN)],
        )
        data = SubscriberData(sid=sid, non_3gpp=config)
        client.AddSubscriber(data)
        sids.append(sid)
    return sids


def _list_subs():
    client = SubscriberDBStub(
        ServiceRegistry.get_rpc_channel(
            SUBSCRIBERDB_SERVICE_NAME, ServiceRegistry.LOCAL,
        ),
    )
    client.ListSubscribers(Void())


def _cleanup_subs():
    client = SubscriberDBStub(
        ServiceRegistry.get_rpc_channel(
            SUBSCRIBERDB_SERVICE_NAME, ServiceRegistry.LOCAL,
        ),
    )

    for sid in client.ListSubscribers(Void()).sids:
        client.DeleteSubscriber(SIDUtils.to_pb('IMSI%s' % sid.id))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title="subcommands", dest="cmd")
    parser_add = subparsers.add_parser(
        "add",
        help="Add Subscriber load test",
    )
    parser_get = subparsers.add_parser(
        "get",
        help="Get Subscriber load test",
    )
    parser_update = subparsers.add_parser(
        "update",
        help="Update Subscriber load test",
    )
    parser_list = subparsers.add_parser(
        "list",
        help="List Subscriber load test",
    )
    parser_delete = subparsers.add_parser(
        "delete",
        help="Delete Subscriber load test",
    )

    # Add arguments
    for cmd in [
        parser_add,
        parser_get,
        parser_update,
        parser_list,
        parser_delete,
    ]:
        cmd.add_argument("--num", default=2000, help="Number of requests")
        cmd.add_argument("--import_path", help="Protobuf dir import path")

    # Add function callbacks
    parser_add.set_defaults(func=parser_add)
    return parser


def _build_add_subs_data(num_subs: int, input_file: str):
    add_subs_reqs = []
    for i in range(1, num_subs):
        sid = SubscriberID(id=str(i).zfill(15))
        config = Non3GPPUserProfile(
            apn_config=[APNConfiguration(service_selection=TEST_APN)],
        )
        data = SubscriberData(sid=sid, non_3gpp=config)
        add_sub_req_dict = json_format.MessageToDict(data)
        add_subs_reqs.append(add_sub_req_dict)

    with open(input_file, 'w') as file:
        json.dump(add_subs_reqs, file, separators=(',', ':'))


def _build_list_subs_data(num_subs: int, input_file: str):
    _load_subs(num_subs)
    list_subs_reqs = []
    for i in range(1, num_subs):
        list_sub_req_dict = json_format.MessageToDict(Void())
        list_subs_reqs.append(list_sub_req_dict)

    with open(input_file, 'w') as file:
        json.dump(list_subs_reqs, file, separators=(',', ':'))


def _build_delete_subs_data(num_subs: int, input_file: str):
    active_sids = _load_subs(num_subs)
    delete_subs_reqs = []
    for sid in active_sids:
        delete_sub_req_dict = json_format.MessageToDict(sid)
        delete_subs_reqs.append(delete_sub_req_dict)

    with open(input_file, 'w') as file:
        json.dump(delete_subs_reqs, file, separators=(',', ':'))


def _build_get_subs_data(num_subs: int, input_file: str):
    active_sids = _load_subs(num_subs)
    get_subs_reqs = []
    for sid in active_sids:
        get_sub_req_dict = json_format.MessageToDict(sid)
        get_subs_reqs.append(get_sub_req_dict)

    with open(input_file, 'w') as file:
        json.dump(get_subs_reqs, file, separators=(',', ':'))


def _build_update_subs_data(num_subs: int, input_file: str):
    active_sids = _load_subs(num_subs)
    update_subs_reqs = []
    for sid in active_sids:
        config = Non3GPPUserProfile(
            apn_config=[APNConfiguration(service_selection=TEST_APN_UPDATE)],
        )
        data = SubscriberData(sid=sid, non_3gpp=config)
        update = SubscriberUpdate(data=data)
        update_sub_req_dict = json_format.MessageToDict(update)
        update_subs_reqs.append(update_sub_req_dict)

    with open(input_file, 'w') as file:
        json.dump(update_subs_reqs, file, separators=(',', ':'))


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    if args.cmd == 'add':
        _cleanup_subs()
        input_file = 'add_subs.json'
        request_type = 'AddSubscriber'
        _build_add_subs_data(args.num, input_file)

    if args.cmd == 'list':
        _cleanup_subs()
        input_file = 'list_subs.json'
        request_type = 'ListSubscribers'
        _build_list_subs_data(args.num, input_file)

    if args.cmd == 'delete':
        _cleanup_subs()
        input_file = 'delete_subs.json'
        request_type = 'DeleteSubscriber'
        _build_delete_subs_data(args.num, input_file)

    if args.cmd == 'get':
        _cleanup_subs()
        input_file = 'get_subs.json'
        request_type = 'GetSubscriberData'
        _build_get_subs_data(args.num, input_file)

    if args.cmd == 'update':
        _cleanup_subs()
        input_file = 'update_subs.json'
        request_type = 'UpdateSubscriber'
        _build_update_subs_data(args.num, input_file)

    benchmark_grpc_request(
        proto_path=PROTO_PATH,
        full_request_type=make_full_request_type(
            SUBSCRIBERDB_SERVICE_RPC_PATH, request_type,
        ),
        input_file=input_file,
        output_file=make_output_file_path(request_type),
        num_reqs=args.num,
        address=SUBSCRIBERDB_PORT,
        import_path=args.import_path,
    )
    _cleanup_subs()


if __name__ == "__main__":
    main()
