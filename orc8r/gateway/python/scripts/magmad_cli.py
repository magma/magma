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

from google.protobuf import json_format
from google.protobuf.struct_pb2 import Struct
from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos import common_pb2, magmad_pb2, magmad_pb2_grpc


@grpc_wrapper
def start_services(client, args):
    client.StartServices(common_pb2.Void())


@grpc_wrapper
def stop_services(client, args):
    client.StopServices(common_pb2.Void())


@grpc_wrapper
def reboot(client, args):
    client.Reboot(common_pb2.Void())


@grpc_wrapper
def restart_services(client, args):
    client.RestartServices(
        magmad_pb2.RestartServicesRequest(services=args.services),
    )


@grpc_wrapper
def ping(client, args):
    response = client.RunNetworkTests(
        magmad_pb2.NetworkTestRequest(
            pings=[
                magmad_pb2.PingParams(
                    host_or_ip=host,
                    num_packets=args.packets,
                ) for host in args.hosts
            ],
        ),
    )
    print(response)


@grpc_wrapper
def traceroute(client, args):
    response = client.RunNetworkTests(
        magmad_pb2.NetworkTestRequest(
            traceroutes=[
                magmad_pb2.TracerouteParams(
                    host_or_ip=host,
                    max_hops=args.max_hops,
                    bytes_per_packet=args.bytes,
                ) for host in args.hosts
            ],
        ),
    )
    print(response)


@grpc_wrapper
def get_gateway_id(client, args):
    response = client.GetGatewayId(common_pb2.Void())
    print(response)


@grpc_wrapper
def generic_command(client, args):
    params = json_format.Parse(args.params, Struct())
    response = client.GenericCommand(
        magmad_pb2.GenericCommandParams(command=args.command, params=params),
    )
    print(response)


@grpc_wrapper
def tail_logs(client, args):
    stream = client.TailLogs(magmad_pb2.TailLogsRequest(service=args.service))
    for log_line in stream:
        print(log_line.line, end='')


@grpc_wrapper
def check_stateless(client, args):
    response = client.CheckStateless(common_pb2.Void())
    print(
        "AGW Mode:",
        magmad_pb2.CheckStatelessResponse.AGWMode.Name(response.agw_mode),
    )


@grpc_wrapper
def config_stateless(client, args):
    if args.switch == "enable":
        print("Enable switch")
        config_arg = magmad_pb2.ConfigureStatelessRequest.ENABLE
    elif args.switch == "disable":
        print("Disable switch")
        config_arg = magmad_pb2.ConfigureStatelessRequest.DISABLE
    client.ConfigureStateless(
        magmad_pb2.ConfigureStatelessRequest(config_cmd=config_arg),
    )


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for Magmad',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_start = subparsers.add_parser(
        'start_services',
        help='Start all magma services',
    )
    parser_stop = subparsers.add_parser(
        'stop_services',
        help='Stop all magma services',
    )
    parser_reboot = subparsers.add_parser(
        'reboot',
        help='Reboot the gateway device',
    )
    parser_restart = subparsers.add_parser(
        'restart_services',
        help='Restart specified magma services',
    )
    parser_ping = subparsers.add_parser(
        'ping',
        help='Ping a host from the gateway',
    )
    parser_traceroute = subparsers.add_parser(
        'traceroute',
        help='traceroute a host from the gateway',
    )
    parser_get_id = subparsers.add_parser(
        'get_gateway_id',
        help='Get gateway hardware ID',
    )
    parser_generic_command = subparsers.add_parser(
        'generic_command',
        help='Execute generic command',
    )
    parser_tail_logs = subparsers.add_parser(
        'tail_logs',
        help='Tail logs',
    )
    parser_stateless_check = subparsers.add_parser(
        'check_stateless',
        help='Check AGW stateless mode',
    )
    parser_stateless_config = subparsers.add_parser(
        'config_stateless',
        help='Change AGW stateless mode',
    )

    parser_ping.add_argument(
        'hosts', nargs='+', type=str,
        help='Hosts (URLs or IPs) to ping',
    )
    parser_ping.add_argument(
        '--packets', type=int, default=4,
        help='Number of packets to send with each ping',
    )

    parser_traceroute.add_argument(
        'hosts', nargs='+', type=str,
        help='Hosts (URLs or IPs) to traceroute',
    )
    parser_traceroute.add_argument(
        '--max-hops', type=int, default=30,
        help='Max TTL for packets, defaults to 30',
    )
    parser_traceroute.add_argument(
        '--bytes', type=int, default=60,
        help='Bytes per packet, defaults to 60',
    )
    parser_restart.add_argument(
        'services', nargs='*', type=str,
        help='Services to restart',
    )
    parser_generic_command.add_argument(
        'command', type=str,
        help='Command name',
    )
    parser_generic_command.add_argument(
        'params', type=str,
        help='Params (string)',
    )
    parser_tail_logs.add_argument(
        'service', type=str, nargs='?',
        help='Service',
    )
    parser_stateless_config.add_argument(
        'switch', type=str,
        help='Enable/Disable',
    )

    # Add function callbacks
    parser_start.set_defaults(func=start_services)
    parser_stop.set_defaults(func=stop_services)
    parser_reboot.set_defaults(func=reboot)
    parser_restart.set_defaults(func=restart_services)
    parser_ping.set_defaults(func=ping)
    parser_traceroute.set_defaults(func=traceroute)
    parser_get_id.set_defaults(func=get_gateway_id)
    parser_generic_command.set_defaults(func=generic_command)
    parser_tail_logs.set_defaults(func=tail_logs)
    parser_stateless_check.set_defaults(func=check_stateless)
    parser_stateless_config.set_defaults(func=config_stateless)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, magmad_pb2_grpc.MagmadStub, 'magmad')


if __name__ == "__main__":
    main()
