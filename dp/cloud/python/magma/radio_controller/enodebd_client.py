"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import json
import logging
import sys
from collections import namedtuple
from typing import Optional, Type, Union

import grpc
from dp.protos.cbsd_pb2 import EnodebdUpdateCbsdRequest
from dp.protos.cbsd_pb2_grpc import CbsdManagementStub
from dp.protos.enodebd_dp_pb2 import CBSDRequest
from dp.protos.enodebd_dp_pb2_grpc import DPServiceStub
from google.protobuf import json_format

SERVICES_ADDRESSES = {
    DPServiceStub: 'localhost:50053',
    CbsdManagementStub: 'orc8r-dp:9180',
}
MAGMA_CERTS = (
    'x-magma-client-cert-serial',
    '7ZZXAF7CAETF241KL22B8YRR7B5UF401',
)

CommandConfig = namedtuple("CommandConfig", 'service_cls method_name message_cls request_kwargs')
COMMANDS_CONFIG = {
    'state': CommandConfig(DPServiceStub, 'GetCBSDState', CBSDRequest, {}),
    'register': CommandConfig(DPServiceStub, 'CBSDRegister', CBSDRequest, {}),
    'deregister': CommandConfig(DPServiceStub, 'CBSDDeregister', CBSDRequest, {}),
    'relinquish': CommandConfig(DPServiceStub, 'CBSDRelinquish', CBSDRequest, {}),
    'update': CommandConfig(
        CbsdManagementStub, 'EnodebdUpdateCbsd', EnodebdUpdateCbsdRequest, {'metadata': (MAGMA_CERTS,)},
    ),
}

ServiceType = Union[DPServiceStub, CbsdManagementStub]
MessageType = Union[CBSDRequest, EnodebdUpdateCbsdRequest]

default_cbsd_dict = {
    "serial_number": "enodebd_client_serial_number",
    "fcc_id": "enodebd_client_fcc_id",
    "user_id": "enodebd_client_user_id",
    "min_power": 0,
    "max_power": 20,
    "antenna_gain": 15,
    "number_of_ports": 2,
}
default_cbsd_update_dict = {
    "serial_number": "enodebd_client_serial_number",
    "installation_param": {
        "antenna_gain": 15,
        "latitude_deg": 10.6,
        "longitude_deg": 11.6,
        "indoor_deployment": True,
        "height_type": "agl",
        "height_m": 12.5,
    },
    "cbsd_category": "a",
}

logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S',
)


def _create_service(
        service_cls: Type[ServiceType], address: Optional[str] = None,
) -> ServiceType:
    """
    Create gRPC service.

    Parameters:
        service_cls: protobuf class of gRPC service
        address: gRPC service URL

    Returns:
        gRPC service instance
    """
    service_address = address or SERVICES_ADDRESSES[service_cls]
    channel = _create_service_channel(address=service_address)
    return service_cls(channel)


def _create_service_channel(address: str) -> grpc.Channel:
    """
    Create gRPC channel for the gRPC service

    Parameters:
        address: Radio Controller gRPC service URL

    Returns:
        gRPC channel
    """
    channel = grpc.insecure_channel(address)
    try:
        grpc.channel_ready_future(channel).result(timeout=10)
    except grpc.FutureTimeoutError:
        sys.exit('Error connecting to Radio Controller service')
    else:
        return channel


def _create_message(
        message_cls: Type[MessageType], json_file: Optional[str] = None,
) -> MessageType:
    """
    Create gRPC message

    Parameters:
        message_cls: protobuf class of the gRPC message
        json_file: path to json file containing message content

    Returns:
        gRPC message instance
    """
    if json_file:
        data = _load_json(path=json_file)
    else:
        data = _get_default_cbsd_message(message_cls=message_cls)

    return json_format.ParseDict(js_dict=data, message=message_cls())


def _dump_message(message: MessageType) -> str:
    """
    Dump message to json

    Parameters:
        message: gRPC message

    Returns:
        marshaled json
    """
    return json_format.MessageToJson(
        message=message,
        including_default_value_fields=True,
        preserving_proto_field_name=True,
        sort_keys=True,
    )


def _load_json(path: str) -> dict:
    """
    Read json file content

    Parameters:
        path: path to json file

    Returns:
        unmarshalled json
    """
    try:
        with open(path, 'r') as f:
            return json.load(f)
    except (ValueError, OSError) as e:
        raise ValueError(f"Failed to read the json file {path}: {e}")


def _get_default_cbsd_message(message_cls: Type[MessageType]) -> dict:
    """
    Get default message content corresponding to given message class

    Parameters:
        message_cls: protobuf class of the gRPC message

    Returns:
        default content for message class
    """
    if message_cls == EnodebdUpdateCbsdRequest:
        return default_cbsd_update_dict
    return default_cbsd_dict


def _create_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()

    commands_group = parser.add_argument_group('commands')
    commands = commands_group.add_mutually_exclusive_group()
    commands.add_argument(
        '-s', '--state', dest='command', action='store_const',
        default='state', const='state', help='RPC get CBSD state [default]',
    )
    commands.add_argument(
        '-r', '--register', dest='command', action='store_const',
        const='register', help='RPC register CBSD [deprecated]',
    )
    commands.add_argument(
        '-d', '--deregister', dest='command', action='store_const',
        const='deregister', help='RPC deregister CBSD [deprecated]',
    )
    commands.add_argument(
        '-e', '--relinquish', dest='command', action='store_const',
        const='relinquish', help='RPC relinquish CBSD',
    )
    commands.add_argument(
        '-u', '--update', dest='command', action='store_const',
        const='update', help='RPC update CBSD params',
    )

    parser.add_argument(
        '-a', '--address', dest='address', action='store', type=str,
        default=None, help='RPC call destination, if different than default',
    )
    parser.add_argument(
        '-c', '--cbsd', dest='json_file', action='store', type=str,
        default=None, help='Path to JSON file with CBSD config, if different than default',
    )
    return parser


def main() -> None:
    parser = _create_argparse()
    args = parser.parse_args()

    command_config = COMMANDS_CONFIG[args.command]
    message = _create_message(
        message_cls=command_config.message_cls, json_file=args.json_file,
    )
    service = _create_service(
        service_cls=command_config.service_cls, address=args.address,
    )
    method = getattr(service, command_config.method_name)

    logging.info(
        f'Sending gRPC {command_config.method_name} '
        f'request:\n{_dump_message(message=message)}',
    )
    response = method(message, **command_config.request_kwargs)
    logging.info(
        f'Received gRPC response:\n'
        f'{_dump_message(message=response)}\n',
    )


if __name__ == '__main__':
    main()
