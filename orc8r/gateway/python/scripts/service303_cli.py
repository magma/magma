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

from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos.common_pb2 import LogLevel, Void
from orc8r.protos.service303_pb2 import LogLevelMessage
from orc8r.protos.service303_pb2_grpc import Service303Stub


@grpc_wrapper
def get_metrics(client, args):
    container = client.GetMetrics(Void())
    for family in container.family:
        print(family)


@grpc_wrapper
def get_info(client, args):
    print(client.GetServiceInfo(Void()))


@grpc_wrapper
def stop_service(client, args):
    client.StopService(Void())


@grpc_wrapper
def set_log_level(client, args):
    message = LogLevelMessage()
    message.level = LogLevel.Value(args.level)
    client.SetLogLevel(message)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for Service303',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_metrics = subparsers.add_parser(
        'metrics', help='Get service metrics',
    )
    parser_stop = subparsers.add_parser('stop', help='Stop service')
    parser_info = subparsers.add_parser('info', help='Get service info')

    parser_log_level = subparsers.add_parser('log_level', help='Set log level')
    parser_log_level.add_argument(
        'level', help='Log level', choices=LogLevel.keys(),
    )

    # Add arguments
    for cmd in [parser_metrics, parser_info, parser_stop, parser_log_level]:
        cmd.add_argument('service', help='Service identifier')

    # Add function callbacks
    parser_metrics.set_defaults(func=get_metrics)
    parser_stop.set_defaults(func=stop_service)
    parser_info.set_defaults(func=get_info)
    parser_log_level.set_defaults(func=set_log_level)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, Service303Stub, args.service)


if __name__ == "__main__":
    main()
