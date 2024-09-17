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
from feg.protos.mock_core_pb2_grpc import MockCoreConfiguratorStub
from magma.common.rpc_utils import cloud_grpc_wrapper
from orc8r.protos.common_pb2 import Void


@cloud_grpc_wrapper
def send_reset(client, args):
    print("Sending reset")
    try:
        client.Reset(Void())
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for mock PCRF',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Reset
    alert_ack_parser = subparsers.add_parser(
        'reset', help='Send Reset to mock PCRF hosted in FeG',
    )
    alert_ack_parser.set_defaults(func=send_reset)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, MockCoreConfiguratorStub, 'pcrf')


if __name__ == "__main__":
    main()
