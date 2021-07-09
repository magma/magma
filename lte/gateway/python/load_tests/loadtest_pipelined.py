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
import json
import os
import random
import subprocess
from collections import namedtuple

from google.protobuf import json_format
from lte.protos.pipelined_pb2 import (
    ActivateFlowsRequest,
    DeactivateFlowsRequest,
    RequestOriginType,
    VersionedPolicy,
    VersionedPolicyID,
)
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule
from lte.protos.apn_pb2 import AggregatedMaximumBitrate
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from magma.pipelined.qos.common import QosManager
from magma.subscriberdb.sid import SIDUtils
from scripts.pipelined_cli import UEInfo, _gen_ue_set

PROTO_DIR = 'lte/protos'
IMPORT_PATH = '/home/vagrant/magma'
RESULTS_PATH = '/var/tmp'


def _build_activate_flows_data(ue_dict, disable_qos):
    activate_flow_reqs = []

    if disable_qos:
        print("QOS Disabled")
        apn_ambr = None
    else:
        print("QOS Enabled")
        apn_ambr = AggregatedMaximumBitrate(
            max_bandwidth_ul=1000000000,
            max_bandwidth_dl=1000000000,
        )
    for ue in ue_dict:
        request = ActivateFlowsRequest(
            sid=SIDUtils.to_pb(ue.imsi_str),
            ip_addr=ue.ipv4_src,
            policies=[
                VersionedPolicy(
                rule=PolicyRule(
                    id=ue.rule_id,
                    priority=10,
                    flow_list=[
                        FlowDescription(
                            match=FlowMatch(
                            ip_dst=convert_ipv4_str_to_ip_proto(ue.ipv4_src),
                            direction=FlowMatch.UPLINK,
                            ),
                        ),
                        FlowDescription(
                            match=FlowMatch(
                            ip_src=convert_ipv4_str_to_ip_proto(ue.ipv4_dst),
                            direction=FlowMatch.DOWNLINK,
                            ),
                        ),
                    ],
                ),
                version=1,
                ),
            ],
            request_origin=RequestOriginType(type=RequestOriginType.GX),
            apn_ambr=apn_ambr,
        )
        request_dict = json_format.MessageToDict(request)
        # Dumping ActivateFlows request into json
        activate_flow_reqs.append(request_dict)
    with open('activate_flows.json', 'w') as file:
        json.dump(activate_flow_reqs, file, separators=(',', ':'))


def _build_deactivate_flows_data(ue_dict):
    deactivate_flow_reqs = []

    for ue in ue_dict:
        request = DeactivateFlowsRequest(
            sid=SIDUtils.to_pb(ue.imsi_str),
            ip_addr=ue.ipv4_src,
            policies=[
                VersionedPolicyID(
                    rule_id=ue.rule_id,
                    version=1,
                ),
            ],
            request_origin=RequestOriginType(type=RequestOriginType.GX),
            remove_default_drop_flows=True,
        )
        request_dict = json_format.MessageToDict(request)
        # Dumping ActivateFlows request into json
        deactivate_flow_reqs.append(request_dict)
    with open('deactivate_flows.json', 'w') as file:
        json.dump(deactivate_flow_reqs, file, separators=(',', ':'))


# Building gHZ cmd and call subprocess with given params
def _get_ghz_cmd_params(req_type: str, num_reqs: int):
    req_name = 'magma.lte.Pipelined/%s' % req_type
    file_name = ''
    if req_type == 'ActivateFlows':
        file_name = 'activate_flows.json'
    elif req_type == 'DeactivateFlows':
        file_name = 'deactivate_flows.json'
    else:
        print('Use valid request type (ActivateFlows/DeactivateFlows)')
        return
    cmd_list = [
        'ghz', '--insecure', '--proto',
        '%s/pipelined.proto' % PROTO_DIR,
        '-i', IMPORT_PATH, '--total', str(num_reqs),
        '--call', req_name, '-D', file_name, '-O', 'json',
        '-o', '%s/result_%s.json' % (RESULTS_PATH, req_type),
        '0.0.0.0:50063',
    ]

    subprocess.call(cmd_list)
    os.remove(file_name)


def activate_flows_test(args):
    ue_dict = _gen_ue_set(args.num_of_ues)
    _build_activate_flows_data(ue_dict, args.disable_qos)
    _benchmark_grpc_request(args, 'ActivateFlows')


def deactivate_flows_test(args):
    ue_dict = _gen_ue_set(args.num_of_ues)
    _build_deactivate_flows_data(ue_dict)
    _benchmark_grpc_request(args, 'DeactivateFlows')


def _benchmark_grpc_request(args, req_name):
    try:
        # call grpc GHZ load test tool
        _get_ghz_cmd_params(req_name, args.num_of_ues),
    except subprocess.CalledProcessError as e:
        print(e.output)
        print('Check if gRPC GHZ tool is installed')


def create_parser():
    """
    Creates the argparse subparser for all args
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    parser_activate = subparsers.add_parser(
        "activate_flows",
        help="ActivateFlows load test",
    )
    parser_deactivate = subparsers.add_parser(
        "deactivate_flows",
        help="DeactivateFlows load test",
    )
    for subcmd in [
        parser_activate,
        parser_deactivate,
    ]:
        subcmd.add_argument(
            '--num_of_ues', help='Number of total UEs to atach',
            type=int, default=600,
        )
        subcmd.add_argument(
            '--disable_qos', help='If we want to disable QOS',
            action="store_true",
        )
    parser_activate.set_defaults(func=activate_flows_test)
    parser_deactivate.set_defaults(func=deactivate_flows_test)

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
