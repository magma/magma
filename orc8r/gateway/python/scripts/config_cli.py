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

from magma.configuration import service_configs

LOG_LEVEL_KEY = 'log_level'

pp = PrettyPrinter(width=1)


def load_override_config(args):
    cfg = service_configs.load_override_config(args.service)
    if cfg is None:
        print('No override config exists!')
        return
    pp.pprint(cfg)


def get_log_level(args):
    cfg = service_configs.load_service_config(args.service)
    if cfg is None or LOG_LEVEL_KEY not in cfg:
        print('No log level set!')
        return
    print('Log level:', cfg[LOG_LEVEL_KEY])


def set_log_level(args):
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

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args)


if __name__ == "__main__":
    main()
