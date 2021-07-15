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
from lte.protos.sms_orc8r_pb2 import SMODownlinkUnitdata
from lte.protos.sms_orc8r_pb2_grpc import SMSOrc8rGatewayServiceStub
from magma.common.rpc_utils import grpc_wrapper


@grpc_wrapper
def send_downlink_unitdata(client, args):
    req = SMODownlinkUnitdata(
            imsi=args.imsi,
            nas_message_container=bytes.fromhex(args.data),
    )
    print("Sending Downlink Unitdata with following fields:\n %s" % req)
    try:
        client.SMODownlink(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='CLI for SMS',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Downlink Unitdata
    downlink_unitdata_parser = subparsers.add_parser(
            'DU', help="Send downlink unitdata to SMSOrc8rGW service",
    )
    downlink_unitdata_parser.add_argument('imsi', help='e.g. 001010000090122 (no prefix required)')
    downlink_unitdata_parser.add_argument('data', help='Data as a hex string e.g. 1fc13a00')
    downlink_unitdata_parser.set_defaults(func=send_downlink_unitdata)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, SMSOrc8rGatewayServiceStub, 'sms_mme_service')


if __name__ == "__main__":
    main()
