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

from lte.protos.policydb_pb2 import FlowQos, PolicyRule, QosArp
from lte.protos.spgw_service_pb2 import CreateBearerRequest, DeleteBearerRequest
from lte.protos.spgw_service_pb2_grpc import SpgwServiceStub
from magma.common.rpc_utils import grpc_wrapper
from magma.subscriberdb.sid import SIDUtils


@grpc_wrapper
def create_bearer(client, args):
    req = CreateBearerRequest(
        sid=SIDUtils.to_pb(args.imsi),
        link_bearer_id=args.lbi,
        policy_rules=[
            PolicyRule(
                qos=FlowQos(
                    qci=args.qci,
                    gbr_ul=args.gbr_ul,
                    gbr_dl=args.gbr_dl,
                    max_req_bw_ul=args.mbr_ul,
                    max_req_bw_dl=args.mbr_dl,
                    arp=QosArp(
                        priority_level=args.priority,
                        pre_capability=args.pre_cap,
                        pre_vulnerability=args.pre_vul,
                    ),
                ),
            ),
        ],
    )
    print("Creating dedicated bearer for : ", args.imsi)
    client.CreateBearer(req)


@grpc_wrapper
def delete_bearer(client, args):
    req = DeleteBearerRequest(
        sid=SIDUtils.to_pb(args.imsi),
        link_bearer_id=args.lbi,
        eps_bearer_ids=[args.ebi],
    )
    print("Deleting dedicated bearer for : ", args.imsi)
    client.DeleteBearer(req)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for SPGW service',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_create = subparsers.add_parser('create', help='Create Bearer')
    parser_delete = subparsers.add_parser('delete', help='Delete Bearer')

    for cmd in [parser_create, parser_delete]:
        cmd.add_argument('imsi', help='Subscriber identifier (IMSI00101..)')
        cmd.add_argument(
            '-lbi', type=int, required=True,
            help='Linked bearer id',
        )

    parser_create.add_argument(
        '--pre_cap', type=int, default=1,
        help='pre capability (0:ENABLE, 1:DISABLE)',
    )
    parser_create.add_argument(
        '--priority', type=int, default=1,
        help='priority level',
    )
    parser_create.add_argument(
        '--pre_vul', type=int, default=0,
        help='pre vulnerability (0:ENABLE, 1:DISABLE)',
    )
    parser_create.add_argument(
        '--qci', type=int, default=1,
        help='[0-9, 65 66, 67, 70, 75, 79]',
    )
    parser_create.add_argument(
        '--gbr_ul', type=int, default=1000000,
        help='UL guaranteed bit rate',
    )
    parser_create.add_argument(
        '--gbr_dl', type=int, default=1000000,
        help='DL guaranteed bit rate',
    )
    parser_create.add_argument(
        '--mbr_ul', type=int, default=1000000,
        help='UL maximum bit rate',
    )
    parser_create.add_argument(
        '--mbr_dl', type=int, default=1000000,
        help='DL maximum bit rate',
    )

    parser_delete.add_argument(
        '-ebi', type=int, required=True,
        help='ID of bearer to delete',
    )

    # Add function callbacks
    parser_create.set_defaults(func=create_bearer)
    parser_delete.set_defaults(func=delete_bearer)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, SpgwServiceStub, 'spgw_service')


if __name__ == "__main__":
    main()
