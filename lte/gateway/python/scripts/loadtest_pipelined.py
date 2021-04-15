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

import json
import argparse
import random
import subprocess
import time
from collections import namedtuple
from datetime import datetime

from google.protobuf import json_format
from lte.protos.pipelined_pb2 import (
    ActivateFlowsRequest,
    DeactivateFlowsRequest,
    DeactivateFlowsResult,
    RequestOriginType,
)
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from lte.protos.policydb_pb2 import PolicyRule
from lte.protos.subscriberdb_pb2 import AggregatedMaximumBitrate
from magma.common.rpc_utils import grpc_wrapper
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from magma.pipelined.qos.common import QosManager
from magma.subscriberdb.sid import SIDUtils


UEInfo = namedtuple('UEInfo', ['imsi_str', 'ipv4_src', 'ipv4_dst',
                               'rule_id'])

def _gen_ue_set(num_of_ues):
    imsi = 123000000
    ue_set = set()
    for _ in range(0, num_of_ues):
        imsi_str = "IMSI" + str(imsi)
        ipv4_src = ".".join(str(random.randint(0, 255)) for _ in range(4))
        ipv4_dst = ".".join(str(random.randint(0, 255)) for _ in range(4))
        rule_id = "allow." + imsi_str
        ue_set.add(UEInfo(imsi_str, ipv4_src, ipv4_dst, rule_id))
        imsi = imsi + 1
    return ue_set

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
            dynamic_rules=[PolicyRule(
                id=ue.rule_id,
                priority=10,
                flow_list=[
                    FlowDescription(match=FlowMatch(
                        ip_dst=convert_ipv4_str_to_ip_proto(ue.ipv4_src),
                        direction=FlowMatch.UPLINK)),
                    FlowDescription(match=FlowMatch(
                        ip_src=convert_ipv4_str_to_ip_proto(ue.ipv4_dst),
                        direction=FlowMatch.DOWNLINK)),
                ],
            )],
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
            rule_ids=[ue.rule_id],
            request_origin=RequestOriginType(type=RequestOriginType.GX),
            remove_default_drop_flows=True)
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
    cmd_list = ['ghz', '--insecure', '--proto',
                '/home/vagrant/magma/lte/protos/pipelined.proto',
                '-i  /home/vagrant/magma/',
                '--total', str(num_reqs), '--call', req_name,
                '-D', file_name, '-O', 'html', '0.0.0.0:50063']

    subprocess.call(cmd_list)
    os.remove(file_name)

@grpc_wrapper
def ghz_attach_test(client, args):
    ue_dict = _gen_ue_set(args.num_of_ues)
    _build_activate_flows_data(ue_dict, args.disable_qos)
    _build_deactivate_flows_data(ue_dict)
    try:
        # call grpc GHZ load test tool
        _get_ghz_cmd_params(args.grpc_func_name, args.num_of_ues),
    except subprocess.CalledProcessError as e:
        print(e.output)
        print('Check if gRPC GHZ tool is installed')

def create_parser():
    """
    Creates the argparse subparser for all args
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    subcmd = subparsers.add_parser('ghz_attach_test',
        help='Sends a set of Activate grpc requests, followed by Deactivates')
    subcmd.add_argument('--num_of_ues', help='Number of total UEs to atach',
                        type=int, default=600)
    subcmd.add_argument('--disable_qos', help='If we want to disable QOS',
                        action="store_true")
    subcmd.add_argument('--grpc_func_name', help='Function name',
                        type=str, default='ActivateFlows')
    subcmd.set_defaults(func=ghz_attach_test)

def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, PipelinedStub, 'pipelined')


if __name__ == "__main__":
    main()
