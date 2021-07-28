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
import os
from pprint import PrettyPrinter

from magma.common.rpc_utils import grpc_wrapper
from magma.common.service_registry import ServiceRegistry
from magma.configuration import service_configs
from orc8r.protos.common_pb2 import Void
from orc8r.protos.magmad_pb2 import RestartServicesRequest
from orc8r.protos.magmad_pb2_grpc import MagmadStub
from orc8r.protos.service303_pb2 import ReloadConfigResponse
from orc8r.protos.service303_pb2_grpc import Service303Stub

pp = PrettyPrinter(width=1)


@grpc_wrapper
def load_override_config(_, args):
    cfg = service_configs.load_override_config(args.service)
    pp.pprint(cfg)


# Log level commands
LOG_LEVEL_KEY = 'log_level'


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
        print('Need to run as root to modify override config.')
        return
    print('New override config:')
    pp.pprint(cfg)
    _reload_service_config(client, args.service)


# Streamer commands
STREAMER_KEYS_BY_SERVICE = {
    'magmad': ['enable_config_streamer'],
    'subscriberdb': ['enable_streaming'],
    'policydb': ['enable_streaming'],
}


@grpc_wrapper
def get_streamer_status(_, args):
    if args.service not in STREAMER_KEYS_BY_SERVICE:
        print('No streamer config available for', args.service)
        return
    cfg = service_configs.load_service_config(args.service)
    if cfg is None:
        print('No config found!')
    for key in STREAMER_KEYS_BY_SERVICE[args.service]:
        enabled = cfg.get(key, False)
        print('%s: %s' % (key, enabled))


@grpc_wrapper
def set_streamer_status(_, args):
    service = args.service
    if args.keys is None:
        keys = STREAMER_KEYS_BY_SERVICE.get(service, [])
    else:
        keys = args.keys.split(',')

    invalid_keys = set(keys) - set(STREAMER_KEYS_BY_SERVICE.get(service, []))
    if invalid_keys:
        print(
            '%s does not have the following streamer config keys: %s' % (
            service, invalid_keys,
            ),
        )
        return

    cfg = service_configs.load_override_config(service)
    if cfg is None:
        cfg = {}

    for key in keys:
        if args.enabled == 'default':
            try:
                cfg.pop(key, None)
            except TypeError:
                # No log level set in the config
                pass
        elif args.enabled == 'True':
            cfg[key] = True
        elif args.enabled == 'False':
            cfg[key] = False
        else:
            print(
                'Invalid argument: %s. '
                'Expected one of "True", "False", "default"' % args.enabled,
            )
            return

    try:
        service_configs.save_override_config(service, cfg)
    except PermissionError:
        print('Need to run as root to modify override config.')
        return

    print('New override config:')
    pp.pprint(cfg)
    # Currently all streamers require service restart for config to take effect
    _restart_service(service)


def _reload_service_config(client, service):
    response = client.ReloadServiceConfig(Void())
    if response.result == ReloadConfigResponse.RELOAD_SUCCESS:
        print(
            'Config successfully reloaded. '
            'Some config may require a service restart to take effect',
        )
    else:
        print('Config cannot be reloaded. Service restart required.')
        should_restart = input('Restart %s? (y/N) ' % service).lower() == 'y'
        if should_restart:
            _restart_service(service)


def _restart_service(service):
    print('Restarting %s...' % service)
    if service == 'magmad':
        _restart_all()
        return

    chan = ServiceRegistry.get_rpc_channel('magmad', ServiceRegistry.LOCAL)
    client = MagmadStub(chan)
    client.RestartServices(RestartServicesRequest(services=[service]))


def _restart_all():
    if os.getuid() != 0:
        print('Need to run as root to restart magmad.')
        return
    os.system('service magma@* stop')
    os.system('service magma@magmad start')
    return


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for configs',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    load_override_cmd = subparsers.add_parser(
        'load_override_config',
        help='Load the config overrides '
             'of the specified service',
    )
    load_override_cmd.add_argument(
        'service', type=str,
        help='Name of the service',
    )

    get_log_level_cmd = subparsers.add_parser(
        'get_log_level',
        help='Get the log level of a '
             'service',
    )
    get_log_level_cmd.add_argument(
        'service', type=str,
        help='Name of the service',
    )

    set_log_level_cmd = subparsers.add_parser(
        'set_log_level',
        help='Set the log level of a '
             'service by updating the '
             'config overrides',
    )
    set_log_level_cmd.add_argument(
        'service', type=str,
        help='Name of the service',
    )
    set_log_level_cmd.add_argument(
        'log_level', type=str,
        help='Log level to be set. '
             'Specify "default" to use the default '
             'log level',
    )

    get_streamer_status_cmd = subparsers.add_parser(
        'get_streamer_status',
        help='Get the streamer '
             'status of a service',
    )
    get_streamer_status_cmd.add_argument(
        'service', type=str,
        help='Name of the service',
    )

    set_streamer_status_cmd = subparsers.add_parser(
        'set_streamer_status',
        help='Set the streamer '
             'status of a service',
    )
    set_streamer_status_cmd.add_argument(
        'service', type=str,
        help='Name of the service',
    )
    set_streamer_status_cmd.add_argument(
        '--keys', type=str,
        help='Comma separated config keys used to control specific streamers '
             'for the service. If this is not set, then all streamers for the '
             'service will be modified.',
    )
    set_streamer_status_cmd.add_argument(
        'enabled', type=str,
        help='"True" to enable the streamer. "False" to disable the streamer. '
             '"default" to use the default config',
    )

    # Add function callbacks
    load_override_cmd.set_defaults(func=load_override_config)
    get_log_level_cmd.set_defaults(func=get_log_level)
    set_log_level_cmd.set_defaults(func=set_log_level)
    get_streamer_status_cmd.set_defaults(func=get_streamer_status)
    set_streamer_status_cmd.set_defaults(func=set_streamer_status)
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
