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

import grpc
from lte.protos.ha_service_pb2 import StartAgwOffloadRequest
from lte.protos.ha_service_pb2_grpc import HaServiceStub
from magma.common.rpc_utils import grpc_wrapper


@grpc_wrapper
def send_offload_trigger(client, args):
    req = StartAgwOffloadRequest(
            enb_id=args.enb_id,
            imsi=args.imsi,
    )
    print("Sending offload trigger with following fields:\n %s" % req)
    try:
        client.StartAgwOffload(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='CLI for High Availability',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Downlink Unitdata
    ha_parser = subparsers.add_parser(
            'offload', help="Send downlink unitdata to SMSOrc8rGW service",
    )
    ha_parser.add_argument('--imsi', help='e.g. 001010000090122 (no prefix required)')
    ha_parser.add_argument("--enb-id", type=int, help="Cell ID to offload")
    ha_parser.set_defaults(func=send_offload_trigger)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, HaServiceStub, 'ha_service')


if __name__ == "__main__":
    main()
