#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse
import binascii
import errno
import subprocess

from magma.common.rpc_utils import grpc_wrapper
from magma.subscriberdb.sid import SIDUtils
from magma.configuration.service_configs import load_service_config
from magma.pipelined.bridge_util import BridgeTools
from orc8r.protos.common_pb2 import Void
from lte.protos.pipelined_pb2 import ActivateFlowsRequest, DeactivateFlowsRequest
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from lte.protos.policydb_pb2 import FlowMatch, FlowDescription, PolicyRule


# --------------------------
# Metering App
# --------------------------

@grpc_wrapper
def get_subscriber_metering_flows(client, _):
    flow_table = client.GetSubscriberMeteringFlows(Void())
    print(flow_table)


def create_metering_parser(apps):
    """
    Creates the argparse subparser for the metering app
    """
    app = apps.add_parser('meter')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('dump_flows',
                                   help='Prints all subscriber metering flows')
    subcmd.set_defaults(func=get_subscriber_metering_flows)


# --------------------------
# Enforcement App
# --------------------------

@grpc_wrapper
def activate_flows(client, args):
    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        rule_ids=args.rule_ids.split(','))
    client.ActivateFlows(request)


@grpc_wrapper
def deactivate_flows(client, args):
    request = DeactivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        rule_ids=args.rule_ids.split(','))
    client.DeactivateFlows(request)


@grpc_wrapper
def activate_dynamic_rule(client, args):
    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        dynamic_rules=[PolicyRule(
            id=args.rule_id,
            priority=args.priority,
            hard_timeout=args.hard_timeout,
            flow_list=[
                FlowDescription(match=FlowMatch(
                    ipv4_dst=args.ipv4_dst, direction=FlowMatch.UPLINK)),
                FlowDescription(match=FlowMatch(
                    ipv4_src=args.ipv4_dst, direction=FlowMatch.DOWNLINK)),
            ],
        )])
    client.ActivateFlows(request)


@grpc_wrapper
def display_flows(_unused, _):
    ENFORCEMENT_TABLE_NUM = 5
    pipelined_config = load_service_config('pipelined')
    bridge_name = pipelined_config['bridge_name']
    flows = []
    try:
        flows = BridgeTools.get_flows_for_bridge(bridge_name,
                                                 ENFORCEMENT_TABLE_NUM)
    except subprocess.CalledProcessError as e:
        if (e.returncode == errno.EPERM):
            print("Need to run as root to dump flows")
        return

    # Parse the flows and print it decoding note
    for flow in flows[1:]:
        flow = flow.replace('00', '').replace('.', '')
        # If there is a note, decode it otherwise just print the flow
        if 'note:' in flow:
            prefix = flow.split('note:')
            print(prefix[0] + "note:" + str(binascii.unhexlify(prefix[1])))
        else:
            print(flow)


def create_enforecement_parser(apps):
    """
    Creates the argparse subparser for the enforcement app
    """
    app = apps.add_parser('enforcement')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('activate_flows', help='Activate flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--rule_ids',
                        help='Comma separated rule ids', default='rule1,rule2')
    subcmd.set_defaults(func=activate_flows)

    subcmd = subparsers.add_parser('deactivate_flows', help='Deactivate flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--rule_ids',
                        help='Comma separated rule ids', default='rule1,rule2')
    subcmd.set_defaults(func=deactivate_flows)

    subcmd = subparsers.add_parser('activate_dynamic_rule',
                                   help='Activate dynamic flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--rule_id', help='rule id to add', default='rule1')
    subcmd.add_argument('--ipv4_dst', help='ipv4 dst for rule', default='')
    subcmd.add_argument('--priority', help='priority for rule',
                        type=int, default=0)
    subcmd.add_argument('--hard_timeout', help='hard timeout for rule',
                        type=int, default=0)
    subcmd.set_defaults(func=activate_dynamic_rule)

    subcmd = subparsers.add_parser('display_flows', help='Display flows')
    subcmd.set_defaults(func=display_flows)


# --------------------------
# Pipelined base CLI
# --------------------------

def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for pipelined',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    apps = parser.add_subparsers(title='apps', dest='cmd')
    create_metering_parser(apps)
    create_enforecement_parser(apps)
    return parser


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
