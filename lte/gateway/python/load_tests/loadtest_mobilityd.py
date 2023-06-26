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
import ipaddress
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
from lte.protos.mobilityd_pb2 import (
    AllocateIPRequest,
    IPBlock,
    ReleaseIPRequest,
    RemoveIPBlockRequest,
)
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from lte.protos.subscriberdb_pb2 import (
    Non3GPPUserProfile,
    SubscriberData,
    SubscriberID,
)
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.sid import SIDUtils
from orc8r.protos.common_pb2 import Void

TEST_APN = 'magma.ipv4'
MOBILITYD_SERVICE_NAME = 'mobilityd'
MOBILITYD_SERVICE_RPC_PATH = 'magma.lte.MobilityService'
SUBSCRIBERDB_SERVICE_NAME = 'subscriberdb'
MOBILITYD_PORT = '0.0.0.0:60051'
PROTO_PATH = PROTO_DIR + '/mobilityd.proto'

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


def _cleanup_subs():
    client = SubscriberDBStub(
        ServiceRegistry.get_rpc_channel(
            SUBSCRIBERDB_SERVICE_NAME, ServiceRegistry.LOCAL,
        ),
    )

    for sid in client.ListSubscribers(Void()).sids:
        client.DeleteSubscriber(SIDUtils.to_pb('IMSI%s' % sid.id))


def _build_allocate_ip_data(num_subs: int, input_file: str):
    active_sids = _load_subs(num_subs)
    allocate_ip_reqs = []
    for sid in active_sids:
        ip_req = AllocateIPRequest(
            sid=sid, version=AllocateIPRequest.IPV4,
            apn=TEST_APN,
        )  # hardcoding APN
        ip_req_dict = json_format.MessageToDict(ip_req)
        # Dumping AllocateIP request into json
        allocate_ip_reqs.append(ip_req_dict)
    with open(input_file, 'w') as file:
        json.dump(allocate_ip_reqs, file, separators=(',', ':'))


def _setup_ip_block(client):
    ip_blocks_rsp = client.ListAddedIPv4Blocks(Void())
    remove_blocks_req = RemoveIPBlockRequest(force=True)
    for block in ip_blocks_rsp.ip_block_list:
        remove_blocks_req.ip_blocks.append(block)
    client.RemoveIPBlock(remove_blocks_req)
    test_ip_block = '192.168.128.0/20'
    ip_block = ipaddress.ip_network(test_ip_block)
    client.AddIPBlock(
        IPBlock(
            version=IPBlock.IPV4,
            net_address=ip_block.network_address.packed,
            prefix_len=ip_block.prefixlen,
        ),
    )


def _build_release_ip_data(client: MobilityServiceStub, input_file: str):
    release_ip_reqs = []
    table = client.GetSubscriberIPTable(Void())
    if not table.entries:
        print('No IPs allocated to be freed, please run allocate test first')
        exit(1)
    for entry in table.entries:
        release_ip_req = ReleaseIPRequest(
            sid=entry.sid, ip=entry.ip,
            apn=entry.apn,
        )
        release_ip_dict = json_format.MessageToDict(release_ip_req)
        # Dumping ReleaseIP request into json
        release_ip_reqs.append(release_ip_dict)
    with open(input_file, 'w') as file:
        json.dump(release_ip_reqs, file)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title="subcommands", dest="cmd")
    parser_allocate = subparsers.add_parser(
        "allocate",
        help="Allocate IP load test",
    )
    parser_release = subparsers.add_parser(
        "release",
        help="Release IP load test",
    )

    # Add arguments
    for cmd in [
        parser_allocate,
        parser_release,
    ]:
        cmd.add_argument("--num", default=2000, help="Number of requests")
        cmd.add_argument("--import_path", help="Protobuf dir import path")

    # Add function callbacks
    parser_allocate.set_defaults(func=parser_allocate)
    parser_release.set_defaults(func=parser_release)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    client = MobilityServiceStub(
        ServiceRegistry.get_rpc_channel(
            MOBILITYD_SERVICE_NAME,
            ServiceRegistry.LOCAL,
        ),
    )

    if args.cmd == 'allocate':
        _cleanup_subs()
        _setup_ip_block(client)
        input_file = 'allocate_data.json'
        request_type = 'AllocateIPAddress'
        _build_allocate_ip_data(args.num, input_file)
        benchmark_grpc_request(
            proto_path=PROTO_PATH,
            full_request_type=make_full_request_type(
                MOBILITYD_SERVICE_RPC_PATH, request_type,
            ),
            input_file=input_file,
            output_file=make_output_file_path(request_type),
            num_reqs=args.num,
            address=MOBILITYD_PORT,
            import_path=args.import_path,
        )
        _cleanup_subs()

    elif args.cmd == 'release':
        input_file = 'release_data.json'
        request_type = 'ReleaseIPAddress'
        _build_release_ip_data(client, input_file)
        benchmark_grpc_request(
            proto_path=PROTO_PATH,
            full_request_type=make_full_request_type(
                MOBILITYD_SERVICE_RPC_PATH, request_type,
            ),
            input_file=input_file,
            output_file=make_output_file_path(request_type),
            num_reqs=args.num,
            address=MOBILITYD_PORT,
            import_path=args.import_path,
        )


if __name__ == "__main__":
    main()
