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

from lte.protos.enodebd_pb2 import (
    EnodebIdentity,
    GetParameterRequest,
    SetParameterRequest,
    SingleEnodebStatus,
)
from lte.protos.enodebd_pb2_grpc import EnodebdStub
from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos.common_pb2 import Void


@grpc_wrapper
def get_parameter(client, args):
    message = GetParameterRequest()
    message.device_serial = args.device_serial
    message.parameter_name = args.parameter_name
    response = client.GetParameter(message)

    for name_value in response.parameters:
        print('%s = %s' % (name_value.name, name_value.value))


@grpc_wrapper
def set_parameter(client, args):
    message = SetParameterRequest()
    message.device_serial = args.device_serial
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
    req = EnodebIdentity()
    req.device_serial = args.device_serial
    client.Configure(req)


@grpc_wrapper
def reboot_enodeb(client, args):
    req = EnodebIdentity()
    req.device_serial = args.device_serial
    client.Reboot(req)


@grpc_wrapper
def reboot_all_enodeb(client, args):
    client.RebootAll(Void())


@grpc_wrapper
def get_status(client, args):
    """ Get status information of enodebd service """
    def print_status_param(enb_status, name, readable_name):
        """
        Print parameter (of type BoolValue) if it exists in status message,
        otherwise print that the parameter is not known
        """
        if name in enb_status:
            _print_str_status_line(readable_name, enb_status[name])
        else:
            _print_str_status_line(readable_name, 'Unknown')

    status = client.GetStatus(Void())
    meta = status.meta
    print_status_param(meta, 'n_enodeb_connected', '# of eNodeB connected')
    print_status_param(meta, 'all_enodeb_configured', 'All eNodeB configured')
    print_status_param(
        meta, 'all_enodeb_opstate_enabled',
        'All eNB Opstate enabled',
    )
    print_status_param(
        meta, 'all_enodeb_rf_tx_configured',
        'All eNB RF TX configured to desired state',
    )
    print_status_param(
        meta, 'any_enodeb_gps_connected',
        'Any eNB GPS connected',
    )
    print_status_param(
        meta, 'all_enodeb_ptp_connected',
        'All eNB PTP connected',
    )
    print_status_param(
        meta, 'all_enodeb_mme_connected',
        'All eNB MME connected',
    )
    print_status_param(meta, 'gateway_gps_longitude', 'Gateway GPS Longitude')
    print_status_param(meta, 'gateway_gps_latitude', 'Gateway GPS Latitude')


@grpc_wrapper
def get_all_status(client, args):
    """ Get status information of each eNodeB """
    def print_enb_status(enb_status):
        print('--- eNodeB Serial:', enb_status.device_serial, '---')
        _print_str_status_line('IP Address', enb_status.ip_address)
        _print_prop_status_line(
            'eNodeB Connected via TR-069', enb_status.connected,
        )
        _print_prop_status_line('eNodeB Configured', enb_status.configured)
        _print_prop_status_line('Opstate Enabled', enb_status.opstate_enabled)
        _print_prop_status_line('RF TX on', enb_status.rf_tx_on)
        _print_prop_status_line('RF TX desired', enb_status.rf_tx_desired)
        _print_prop_status_line('GPS Connected', enb_status.gps_connected)
        _print_prop_status_line('PTP Connected', enb_status.ptp_connected)
        _print_prop_status_line('MME Connected', enb_status.mme_connected)
        _print_str_status_line('GPS Longitude', enb_status.gps_longitude)
        _print_str_status_line('GPS Latitude', enb_status.gps_latitude)
        _print_str_status_line('FSM State', enb_status.fsm_state)
        print('\n')

    status = client.GetAllEnodebStatus(Void())
    status_list = status.enb_status_list
    if len(status_list) == 0:
        print('No status information to report.')
        print(
            'Either there are no connected eNodeB devices, '
            'or no TR-069 messages have been received yet.',
        )
    else:
        for enb_status in status_list:
            print_enb_status(enb_status)


@grpc_wrapper
def get_enb_status(client, args):
    """ Get status information for a particular eNodeB """
    req = EnodebIdentity()
    req.device_serial = args.device_serial
    enb_status = client.GetEnodebStatus(req)
    _print_prop_status_line('eNodeB Connected', enb_status.connected)
    _print_prop_status_line('eNodeB Configured', enb_status.configured)
    _print_prop_status_line('Opstate Enabled', enb_status.opstate_enabled)
    _print_prop_status_line('RF TX on', enb_status.rf_tx_on)
    _print_prop_status_line('RF TX desired', enb_status.rf_tx_desired)
    _print_prop_status_line('GPS Connected', enb_status.gps_connected)
    _print_prop_status_line('PTP Connected', enb_status.ptp_connected)
    _print_prop_status_line('MME Connected', enb_status.mme_connected)
    _print_str_status_line('GPS Longitude', enb_status.gps_longitude)
    _print_str_status_line('GPS Latitude', enb_status.gps_latitude)
    _print_str_status_line('FSM State', enb_status.fsm_state)


def _print_prop_status_line(header: str, value: int) -> None:
    """ Argument 'value' should be a StatusProperty enum """
    _print_str_status_line(
        header,
        SingleEnodebStatus.StatusProperty.Name(value),
    )


def _print_str_status_line(header: str, value: str) -> None:
    """
    Print a single line for status info.

    Example output:
    All eNB RF TX on...........False
    """
    formatted_str = '{:25}{:>15}'.format(header + ':', ':' + value)
    parts = formatted_str.split(':')
    parts[1] = str(parts[1]).replace(' ', '.')
    print(''.join(parts))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for Enodebd',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    parser_get_parameter = subparsers.add_parser(
        'get_parameter', help='Send GetParameterValues message',
    )
    parser_get_parameter.add_argument(
        'device_serial', help='eNodeB Serial ID',
    )
    parser_get_parameter.add_argument(
        'parameter_name', help='Parameter Name',
    )

    parser_set_parameter = subparsers.add_parser(
        'set_parameter', help='Send SetParameterValues message',
    )
    parser_set_parameter.add_argument(
        'device_serial', help='eNodeB Serial ID',
    )
    parser_set_parameter.add_argument(
        'parameter_name', help='Parameter Name',
    )
    parser_set_parameter.add_argument(
        'value', help='Parameter Value',
    )
    parser_set_parameter.add_argument(
        'data_type', help='Parameter Data Type',
        choices=['int', 'bool', 'string'],
    )
    parser_set_parameter.add_argument(
        '--parameter_key', default='', help='Parameter Key',
    )

    parser_config_enodeb = subparsers.add_parser(
        'config_enodeb', help='Configure eNodeB',
    )
    parser_config_enodeb.add_argument(
        'device_serial', help='eNodeB Serial ID',
    )

    parser_reboot_enodeb = subparsers.add_parser(
        'reboot_enodeb', help='Reboot eNodeB',
    )
    parser_reboot_enodeb.add_argument(
        'device_serial', help='eNodeB Serial ID',
    )

    parser_reboot_all_enodeb = subparsers.add_parser(
        'reboot_all_enodeb', help='Reboot all eNodeB',
    )

    parser_get_status = subparsers.add_parser(
        'get_status', help='Get enodebd status',
    )

    parser_get_all_status = subparsers.add_parser(
        'get_all_status', help='Get all attached eNodeB status',
    )

    parser_get_enb_status = subparsers.add_parser(
        'get_enb_status', help='Get eNodeB status',
    )
    parser_get_enb_status.add_argument(
        'device_serial', help='eNodeB Serial ID',
    )

    # Add function callbacks
    parser_get_parameter.set_defaults(func=get_parameter)
    parser_set_parameter.set_defaults(func=set_parameter)
    parser_config_enodeb.set_defaults(func=configure_enodeb)
    parser_reboot_enodeb.set_defaults(func=reboot_enodeb)
    parser_reboot_all_enodeb.set_defaults(func=reboot_all_enodeb)
    parser_get_status.set_defaults(func=get_status)
    parser_get_all_status.set_defaults(func=get_all_status)
    parser_get_enb_status.set_defaults(func=get_enb_status)
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
