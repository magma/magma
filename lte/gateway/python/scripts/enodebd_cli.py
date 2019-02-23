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
from orc8r.protos.common_pb2 import Void
from lte.protos.enodebd_pb2 import GetParameterRequest
from lte.protos.enodebd_pb2 import SetParameterRequest
from lte.protos.enodebd_pb2_grpc import EnodebdStub


@grpc_wrapper
def get_parameter(client, args):
    message = GetParameterRequest()
    message.parameter_name = args.parameter_name
    response = client.GetParameter(message)

    for name_value in response.parameters:
        print('%s = %s' % (name_value.name, name_value.value))


@grpc_wrapper
def set_parameter(client, args):
    message = SetParameterRequest()
    message.parameter_name = args.parameter_name
    if args.data_type == 'bool':
        if args.value.lower() == 'true':
            message.value_bool = True
        elif args.value.lower() == 'false':
            message.value_bool = False
        else:
            raise TypeError("Bool values should be True or False")
    elif args.data_type == 'string':
        message.value_string = str(args.value)
    elif args.data_type == 'int':
        message.value_int = int(args.value)
    else:
        raise TypeError('Unknown type %s' % args.data_type)
    client.SetParameter(message)


@grpc_wrapper
def configure_enodeb(client, args):
    client.Configure(Void())


@grpc_wrapper
def reboot_enodeb(client, args):
    client.Reboot(Void())


@grpc_wrapper
def get_status(client, args):
    """
    Get status information
    """
    def print_status_param(status, name, readable_name):
        """ Print parameter (of type BoolValue) if it exists in status message,
            otherwise print that the parameter is not known """
        meta = status.meta
        if name in meta:
            print('%s: %s' % (readable_name, meta[name]))
        else:
            print('%s: Unknown' % readable_name)

    status = client.GetStatus(Void())
    print_status_param(status, 'enodeb_connected', 'eNodeB connected')
    print_status_param(status, 'enodeb_configured', 'eNodeB configured')
    print_status_param(status, 'opstate_enabled', 'Opstate enabled')
    print_status_param(status, 'rf_tx_on', 'RF TX on')
    print_status_param(status, 'gps_connected', 'GPS connected')
    print_status_param(status, 'ptp_connected', 'PTP connected')
    print_status_param(status, 'mme_connected', 'MME connected')
    print_status_param(status, 'gps_longitude', 'GPS Longitude')
    print_status_param(status, 'gps_latitude', 'GPS Latitude')
    print_status_param(status, 'enodeb_state', 'FSM State')


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for Enodebd',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    parser_get_parameter = subparsers.add_parser(
        'get_parameter', help='Send GetParameterValues message')
    parser_get_parameter.add_argument(
        'parameter_name', help='Parameter Name')

    parser_set_parameter = subparsers.add_parser(
        'set_parameter', help='Send SetParameterValues message')
    parser_set_parameter.add_argument(
        'parameter_name', help='Parameter Name')
    parser_set_parameter.add_argument(
        'value', help='Parameter Value')
    parser_set_parameter.add_argument(
        'data_type', help='Parameter Data Type',
        choices=['int', 'bool', 'string'])
    parser_set_parameter.add_argument(
        '--parameter_key', default='', help='Parameter Key')

    parser_config_enodeb = subparsers.add_parser(
        'config_enodeb', help='Configure eNodeB')

    parser_reboot_enodeb = subparsers.add_parser(
        'reboot_enodeb', help='Reboot eNodeB')

    parser_get_status = subparsers.add_parser(
        'get_status', help='Get eNodeB status')

    # Add function callbacks
    parser_get_parameter.set_defaults(func=get_parameter)
    parser_set_parameter.set_defaults(func=set_parameter)
    parser_config_enodeb.set_defaults(func=configure_enodeb)
    parser_reboot_enodeb.set_defaults(func=reboot_enodeb)
    parser_get_status.set_defaults(func=get_status)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, EnodebdStub, 'enodebd')

if __name__ == "__main__":
    main()
