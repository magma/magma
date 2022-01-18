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

import grpc
from dp.protos import enodebd_dp_pb2 as enodebd_msgs
from dp.protos import enodebd_dp_pb2_grpc as enodebd_services

default_cbsd_dict = {
    "serial_number": "enodebd_client_serial_number",
    "fcc_id": "enodebd_client_fcc_id",
    "user_id": "enodebd_client_user_id",
    "min_power": 0,
    "max_power": 20,
    "antenna_gain": 15,
    "number_of_ports": 2,
}

DP_RC_SERVICE_ADDR = 'localhost:50053'

logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S',
)


def create_rc_service_channel(addr: str):
    """
    Create Radio Controller service gRPC channel

    Parameters:
        addr: Radio Controller gRPC service URL

    Returns:
        gRPC channel
    """
    channel = grpc.insecure_channel(addr)
    try:
        grpc.channel_ready_future(channel).result(timeout=10)
    except grpc.FutureTimeoutError:
        sys.exit('Error connecting to Radio Controller service')
    else:
        return channel


def create_grpc_dp_service(channel):
    """
    Create Radio Controller gRPC service

    Parameters:
        channel: Radio Controller gRPC service channel

    Returns:
        DPServiceStub
    """
    stub = enodebd_services.DPServiceStub(channel)
    return stub


def create_rc_grpc_cbsd_request(**kwargs):
    """
    Construct a CBSDRequest gRPC message

    Parameters:
        kwargs: dict where keys are CBSDRequest fields, and values are the values

    Returns:
        CBSDRequest message
    """
    msg = enodebd_msgs.CBSDRequest(**kwargs)
    return msg


def send_request(service, rpc, msg):
    """
    Send a CBSDRequest gRPC message to gRPC endpoint

    Parameters:
        service: DBService
        rpc: the gRPC method of the service
        msg: CBSDRequest message

    Returns:
        CBSDStateResult
    """
    resp = rpc(msg)
    return resp


def _get_cbsd_dict(json_file_path: str) -> dict:
    if not json_file_path:
        return default_cbsd_dict
    try:
        with open(json_file_path, 'r') as f:
            return json.load(f)
    except (ValueError, OSError) as err:
        logging.warning(
            f"Failed to read or parse CBSD file {json_file_path}. {err}",
        )
        raise


def _create_argparse(dp_service: enodebd_services.DPServiceStub):
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '-s', '--state', dest='rpc', action='store_const',
        default=dp_service.GetCBSDState, const=dp_service.GetCBSDState, help='CBSD Get State RPC',
    )
    parser.add_argument(
        '-r', '--register', dest='rpc', action='store_const',
        const=dp_service.CBSDRegister, help='CBSD Register RPC',
    )
    parser.add_argument(
        '-d', '--deregister', dest='rpc', action='store_const',
        const=dp_service.CBSDDeregister, help='CBSD Deregister RPC',
    )
    parser.add_argument(
        '-e', '--relinquish', dest='rpc', action='store_const',
        const=dp_service.CBSDRelinquish, help='CBSD Register RPC',
    )
    parser.add_argument(
        '-c', '--cbsd', dest='cbsd_json_file', action='store', type=str,
        default=None, help='Path to JSON file with CBSD config params',
    )
    return parser


def main():
    """
    Top level function of the module
    """
    channel = create_rc_service_channel(DP_RC_SERVICE_ADDR)
    dp_service = create_grpc_dp_service(channel)
    parser = _create_argparse(dp_service=dp_service)
    args = parser.parse_args()
    cbsd = _get_cbsd_dict(args.cbsd_json_file)
    msg = create_rc_grpc_cbsd_request(**cbsd)
    rpc_method_name = str(args.rpc._method)
    logging.info(f'Sending gRPC {rpc_method_name} request:\n{msg}')
    resp = send_request(dp_service, args.rpc, msg)
    logging.info(f'Received gRPC response:\n{resp}')


if __name__ == '__main__':
    main()
