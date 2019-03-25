#!/usr/bin/env python3

"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import argparse
from pprint import PrettyPrinter

from magma.common.service_registry import ServiceRegistry
from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2 import ReloadConfigResponse
from orc8r.protos.service303_pb2_grpc import Service303Stub
from orc8r.protos.magmad_pb2 import RestartServicesRequest
from orc8r.protos.magmad_pb2_grpc import MagmadStub

from magma.common.rpc_utils import grpc_wrapper
from magma.configuration import service_configs

LOG_LEVEL_KEY = 'log_level'

pp = PrettyPrinter(width=1)


@grpc_wrapper
def load_override_config(_, args):
    cfg = service_configs.load_override_config(args.service)
    pp.pprint(cfg)


@grpc_wrapper
def get_log_level(_, args):
    cfg = service_configs.load_service_config(args.service)
    if cfg is None or LOG_LEVEL_KEY not in cfg:
        print('No log level set!')
        return
    print('Log level:', cfg[LOG_LEVEL_KEY])


@grpc_wrapper
def set_log_level(client, args):
    cfg = service_configs.load_override_config(args.service)
    if cfg is None:
        cfg = {}
    if args.log_level == 'default':
        try:
            cfg.pop(LOG_LEVEL_KEY, None)
        except TypeError:
            # No log level set in the config
            pass
    else:
        cfg[LOG_LEVEL_KEY] = args.log_level

    try:
        service_configs.save_override_config(args.service, cfg)
    except PermissionError:
        print('Need to run as root to modify override config. Use "venvsudo".')
        return
    print('New override config:')
    pp.pprint(cfg)
    _reload_service_config(client, args.service)


def _reload_service_config(client, service):
    response = client.ReloadServiceConfig(Void())
    if response.result == ReloadConfigResponse.RELOAD_SUCCESS:
        print('Config successfully reloaded. '
              'Some config may require a service restart to take effect')
    else:
        print('Config cannot be reloaded. Service restart required.')
        should_restart = input('Restart %s? (y/N) ' % service).lower() == 'y'
        if should_restart:
            _restart_service(service)


def _restart_service(service):
    print('Restarting %s...' % service)
    chan = ServiceRegistry.get_rpc_channel('magmad', ServiceRegistry.LOCAL)
    client = MagmadStub(chan)
    client.RestartServices(RestartServicesRequest(services=[service]))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for configs',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    load_override_cmd = subparsers.add_parser('load_override_config',
                                              help='Load the config overrides '
                                                   'of the specified service')
    load_override_cmd.add_argument('service', type=str,
                                   help='Name of the service')

    get_log_level_cmd = subparsers.add_parser('get_log_level',
                                              help='Get the log level of a '
                                                   'service')
    get_log_level_cmd.add_argument('service', type=str,
                                   help='Name of the service')

    set_log_level_cmd = subparsers.add_parser('set_log_level',
                                              help='Set the log level of a '
                                                   'service by updating the '
                                                   'config overrides')
    set_log_level_cmd.add_argument('service', type=str,
                                   help='Name of the service')
    set_log_level_cmd.add_argument('log_level', type=str,
                                   help='Log level to be set. '
                                        'Specify "default" to use the default '
                                        'log level')

    # Add function callbacks
    load_override_cmd.set_defaults(func=load_override_config)
    get_log_level_cmd.set_defaults(func=get_log_level)
    set_log_level_cmd.set_defaults(func=set_log_level)
    return parser


def main():
    parser = create_parser()

    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, Service303Stub, args.service)


if __name__ == "__main__":
    main()
