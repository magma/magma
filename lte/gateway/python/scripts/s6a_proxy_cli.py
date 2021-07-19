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
from feg.protos.s6a_proxy_pb2 import (
    AuthenticationInformationRequest,
    UpdateLocationRequest,
)
from feg.protos.s6a_proxy_pb2_grpc import S6aProxyStub
from magma.common.rpc_utils import cloud_grpc_wrapper

RESYNC_INFO_BYTES = 30


@cloud_grpc_wrapper
def send_air(client, args):
    """
    Sends an AIR (authentication information request) via gRPC to the S6aProxy.
    """
    req = AuthenticationInformationRequest()
    req.user_name = args.user_name
    req.visited_plmn = args.visited_plmn
    req.num_requested_eutran_vectors = args.num_requested_eutran_vectors
    req.immediate_response_preferred = args.immediate_response_preferred
    req.resync_info = args.resync_info

    print("sending AIR:\n %s" % req)
    try:
        answer = client.AuthenticationInformation(req)
        print('answer:\nerror code: %d' % answer.error_code)
        print('got %d E-UTRAN vector(s)' % len(answer.eutran_vectors))
        for i, vector in enumerate(answer.eutran_vectors):
            print('vector %d' % (i + 1))
            print("\tRAND = ", vector.rand)
            print("\tXRES = ", vector.xres)
            print("\tAUTN = ", vector.autn)
            print("\tKASME = ", vector.kasme)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


@cloud_grpc_wrapper
def send_ulr(client, args):
    """
    Sends an ULR (update location request) via gRPC to the S6aProxy.
    """
    req = UpdateLocationRequest()
    req.user_name = args.user_name
    req.skip_subscriber_data = False
    req.initial_attach = True

    print("sending ULR:\n %s" % req)
    try:
        answer = client.UpdateLocation(req)
        print('answer:')
        print(answer)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='CLI for S6A proxy to contact HSS',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Add AIR subcommand
    # usage: ./python/scripts/s6a_service_cli.py AIR user123
    parser_air = subparsers.add_parser(
        'AIR',
        help='Authentication Information'
             'Request',
    )
    parser_air.add_argument('user_name', help='subscriber identifier')
    parser_air.add_argument(
        'visited_plmn', help='visited site identifier',
        default=b'\x02\xf8\x59', nargs='?',
    )
    parser_air.add_argument(
        'num_requested_eutran_vectors',
        help='number of vectors to request in response',
        default=1, nargs='?',
    )
    parser_air.add_argument(
        'immediate_response_preferred',
        help='indicates to the HSS the values are'
             'requested for immediate attach',
        default=False, nargs='?',
    )
    parser_air.add_argument(
        'resync_info',
        help='concatenation of RAND and AUTS in the case'
             'of a resync attach case',
        default=b'\x00' * RESYNC_INFO_BYTES, nargs='?',
    )
    parser_air.set_defaults(func=send_air)

    # Add ULR subcommand
    # usage: ./python/scripts/s6a_service_cli.py ULR user123
    parser_ulr = subparsers.add_parser(
        'ULR',
        help='Update Location Request',
    )
    parser_ulr.add_argument('user_name', help='subscriber identifier')
    parser_ulr.add_argument(
        'visited_plmn', help='visited site identifier',
        default=b'\x02\xf8\x59', nargs='?',
    )
    parser_ulr.add_argument(
        'skip_subscriber_data',
        help='Skip subscription data in response',
        default=False, nargs='?',
    )
    parser_ulr.add_argument(
        'initial_attach',
        help='Send Cancel Location to other EPCs serving '
             'the UE',
        default=False, nargs='?',
    )
    parser_ulr.set_defaults(func=send_ulr)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, S6aProxyStub, 's6a_proxy')


if __name__ == "__main__":
    main()
