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

from lte.protos.s6a_service_pb2 import DeleteSubscriberRequest
from lte.protos.s6a_service_pb2_grpc import S6aServiceStub
from magma.common.rpc_utils import grpc_wrapper


@grpc_wrapper
def delete_subscriber(client, args):
    req = DeleteSubscriberRequest()
    req.imsi_list.extend(args.imsi_list)
    print("deleting subs:", req.imsi_list)
    client.DeleteSubscriber(req)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for S6A Service',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_delete_subscriber = subparsers.add_parser(
        'delete',
        help='Delete Subscriber',
    )
    # usage: ./python/scripts/s6a_service_cli.py delete IMSI12345 IMSI23456 IMSI34567
    parser_delete_subscriber.add_argument(
        'imsi_list', nargs='+',
        help='imsis to delete',
    )

    # Add function callbacks
    parser_delete_subscriber.set_defaults(func=delete_subscriber)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, S6aServiceStub, 's6a_service')


if __name__ == "__main__":
    main()
