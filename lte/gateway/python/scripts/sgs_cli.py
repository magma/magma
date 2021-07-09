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

import grpc
from feg.protos.csfb_pb2 import (
    AlertAck,
    AlertReject,
    EPSDetachIndication,
    IMSIDetachIndication,
)
from feg.protos.csfb_pb2_grpc import CSFBFedGWServiceStub
from magma.common.rpc_utils import cloud_grpc_wrapper


@cloud_grpc_wrapper
def send_alert_ack(client, args):
    req = AlertAck(imsi=args.imsi)
    print("Sending Alert Ack with following fields:\n %s" % req)
    try:
        client.AlertAc(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


@cloud_grpc_wrapper
def send_alert_reject(client, args):
    req = AlertReject(imsi=args.imsi, sgs_cause=b'\x01')
    print("Sending Alert Reject with following fields:\n %s" % req)
    try:
        client.AlertRej(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


@cloud_grpc_wrapper
def send_eps_detach_indication(client, args):
    req = EPSDetachIndication(
        imsi=args.imsi,
        mme_name=args.mme_name,
        imsi_detach_from_eps_service_type=bytes(
            [args.imsi_detach_from_eps_service_type],
        ),
    )
    print("Sending EPS Detach Indication with following fields:\n %s" % req)
    try:
        client.EPSDetachInd(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


@cloud_grpc_wrapper
def send_imsi_detach_indication(client, args):
    req = IMSIDetachIndication(
        imsi=args.imsi,
        mme_name=args.mme_name,
        imsi_detach_from_non_eps_service_type=b'\x11',
    )
    print("Sending IMSI Detach Indication with following fields:\n %s" % req)
    try:
        client.IMSIDetachInd(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for CSFB',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Alert Ack
    alert_ack_parser = subparsers.add_parser(
        'AA', help='Send Alert Ack to CSFB service in FeG',
    )
    alert_ack_parser.add_argument('imsi', help='e.g.123456')
    alert_ack_parser.set_defaults(func=send_alert_ack)

    # Alert Reject
    alert_reject_parser = subparsers.add_parser(
        'AR', help='Send Alert Reject to csfb in feg',
    )
    alert_reject_parser.add_argument('imsi', help='e.g. 123456')
    alert_reject_parser.set_defaults(func=send_alert_reject)

    # EPS Detach Indication
    eps_detach_indication_parser = subparsers.add_parser(
        'EDI', help='Send EPS Detach Indication to CSFB service in FeG',
    )
    eps_detach_indication_parser.add_argument('imsi', help='e.g. 123456')
    eps_detach_indication_parser.add_argument(
        'mme_name',
        help='MME name is a 55-character FQDN, specified in 3GPP TS 23.003',
    )
    eps_detach_indication_parser.add_argument(
        'imsi_detach_from_eps_service_type',
        help='Enter either 1, 2 or 3', choices=[1, 2, 3], type=int,
    )
    eps_detach_indication_parser.set_defaults(func=send_eps_detach_indication)

    # IMSI Detach Indication
    imsi_detach_indication_parser = subparsers.add_parser(
        'IDI', help='Send IMSI Detach Indication to CSFB service in FeG',
    )
    imsi_detach_indication_parser.add_argument('imsi', help='e.g. 123456')
    imsi_detach_indication_parser.add_argument(
        'mme_name',
        help='MME name is a 55-character FQDN, specified in 3GPP TS 23.003',
    )
    imsi_detach_indication_parser.set_defaults(
        func=send_imsi_detach_indication,
    )

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, CSFBFedGWServiceStub, 'csfb')


if __name__ == "__main__":
    main()
