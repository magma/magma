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
import subprocess
from typing import List

from google.protobuf import json_format
from lte.protos.mobilityd_pb2 import AllocateIPRequest, ReleaseIPRequest
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from lte.protos.subscriberdb_pb2 import SubscriberID, SubscriberData, Non3GPPUserProfile, APNConfiguration
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from orc8r.protos.common_pb2 import Void

from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.sid import SIDUtils


# Helper functions to build input data for gRPC functions
def _load_subs(num_subs: int) -> List[SubscriberID]:
    client = SubscriberDBStub(
        ServiceRegistry.get_rpc_channel('subscriberdb', ServiceRegistry.LOCAL))
    sids = []

    for i in range(1, num_subs):
        sid = SubscriberID(id=str(i).zfill(15))
        config = Non3GPPUserProfile(apn_config=[APNConfiguration(service_selection="magma.ipv4")])
        data = SubscriberData(sid=sid, non_3gpp=config)
        client.AddSubscriber(data)
        sids.append(sid)
    return sids

def _cleanup_subs():
    client = SubscriberDBStub(
        ServiceRegistry.get_rpc_channel('subscriberdb', ServiceRegistry.LOCAL))

    for sid in client.ListSubscribers(Void()).sids:
        sid = 'IMSI'+sid.id
        client.DeleteSubscriber(SIDUtils.to_pb(sid))

def _build_allocate_ip_data(num_subs: int):
    active_sids = _load_subs(num_subs)
    allocate_ip_reqs = []
    for sid in active_sids:
        ip_req = AllocateIPRequest(sid=sid, version=AllocateIPRequest.IPV4,
                                   apn='magma.ipv4') # hardcoding APN
        ip_req_dict = json_format.MessageToDict(ip_req)
        # Dumping AllocateIP request into json
        allocate_ip_reqs.append(ip_req_dict)
    with open('allocate_data.json', 'w') as file:
        json.dump(allocate_ip_reqs, file, separators=(',', ':'))

def _build_release_ip_data():
    release_ip_reqs = []
    client = MobilityServiceStub(
        ServiceRegistry.get_rpc_channel('mobilityd',
                                        ServiceRegistry.LOCAL))
    table = client.GetSubscriberIPTable(Void())
    for entry in table.entries:
        release_ip_req = ReleaseIPRequest(sid=entry.sid, ip=entry.ip,
                                          apn=entry.apn)
        release_ip_dict = json_format.MessageToDict(release_ip_req)
        # Dumping ReleaseIP request into json
        release_ip_reqs.append(release_ip_dict)
    with open('release_data.json', 'w') as file:
        json.dump(release_ip_reqs, file)

# Building gHZ cmd and call subprocess with given params
def _get_ghz_cmd_params(req_type: str, num_reqs: int):
    req_name = 'magma.lte.MobilityService/%s' % req_type
    file_name = ''
    if req_type == 'AllocateIPAddress':
        file_name = 'allocate_data.json'
    elif req_type == 'ReleaseIPAddress':
        file_name = 'release_data.json'
    cmd_list = ['/home/vagrant/magma/lte/gateway/ghz/cmd/ghz/ghz',
                '--insecure', '--proto', './mobilityd.proto', '--total',
                str(num_reqs), '--call', req_name,
                '-D', file_name, '-O', 'csv', '0.0.0.0:60051']

    subprocess.call(cmd_list)


def _benchmark_grpc_request(args, req_name):
    try:
        # call grpc GHZ load test tool
        _get_ghz_cmd_params(req_name, args.num),
    except subprocess.CalledProcessError as e:
        print(e.output)
        print('Check if gRPC GHZ tool is installed')


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title="subcommands", dest="cmd")
    parser_allocate = subparsers.add_parser("allocate",
                                            help="Allocate IP load test")
    parser_release = subparsers.add_parser("release",
                                           help="Release IP load test")

    # Add arguments
    for cmd in [
        parser_allocate,
        parser_release,
    ]:
        cmd.add_argument("--num", default=2000, help="--num")

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

    _cleanup_subs()
    if args.cmd == 'allocate':
        print('Preparing AllocateIPRequest load test...')
        _build_allocate_ip_data(args.num)
        print('Launching AllocateIPRequest load test...')
        _benchmark_grpc_request(args, 'AllocateIPAddress')
        print('Done')

    elif args.cmd == 'release':
        print('Preparing ReleaseIPRequest load test...')
        _build_release_ip_data()
        print('Launching ReleaseIPRequest load test...')
        _benchmark_grpc_request(args, 'ReleaseIPAddress')
        print('Done')
    _cleanup_subs()


if __name__ == "__main__":
    main()
