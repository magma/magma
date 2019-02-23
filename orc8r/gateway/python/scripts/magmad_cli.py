#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse

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
        magmad_pb2.RestartServicesRequest(services=args.services)
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
            ]
        )
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
            ]
        )
    )
    print(response)


@grpc_wrapper
def get_gateway_id(client, args):
    response = client.GetGatewayId(common_pb2.Void())
    print(response)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for Magmad',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_start = subparsers.add_parser('start_services',
                                         help='Start all magma services')
    parser_stop = subparsers.add_parser('stop_services',
                                        help='Stop all magma services')
    parser_reboot = subparsers.add_parser('reboot',
                                          help='Reboot the gateway device')
    parser_restart = subparsers.add_parser('restart_services',
                                           help='Restart specified magma services')
    parser_ping = subparsers.add_parser(
        'ping',
        help='Ping a host from the gateway')
    parser_traceroute = subparsers.add_parser(
        'traceroute',
        help='traceroute a host from the gateway')
    parser_get_id = subparsers.add_parser('get_gateway_id',
                                           help='Get gateway hardware ID')

    parser_ping.add_argument('hosts', nargs='+', type=str,
                             help='Hosts (URLs or IPs) to ping')
    parser_ping.add_argument('--packets', type=int, default=4,
                             help='Number of packets to send with each ping')

    parser_traceroute.add_argument('hosts', nargs='+', type=str,
                                   help='Hosts (URLs or IPs) to traceroute')
    parser_traceroute.add_argument('--max-hops', type=int, default=30,
                                   help='Max TTL for packets, defaults to 30')
    parser_traceroute.add_argument('--bytes', type=int, default=60,
                                   help='Bytes per packet, defaults to 60')
    parser_restart.add_argument('services', nargs='*', type=str,
                                help='Services to restart')
    # Add function callbacks
    parser_start.set_defaults(func=start_services)
    parser_stop.set_defaults(func=stop_services)
    parser_reboot.set_defaults(func=reboot)
    parser_restart.set_defaults(func=restart_services)
    parser_ping.set_defaults(func=ping)
    parser_traceroute.set_defaults(func=traceroute)
    parser_get_id.set_defaults(func=get_gateway_id)
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
