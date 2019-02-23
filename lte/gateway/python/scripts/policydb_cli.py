#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse

from magma.policydb.rule_store import PolicyRuleDict
from lte.protos.policydb_pb2 import FlowMatch, FlowDescription, PolicyRule


def add_rule(args):
    rule_id = args.rule_id
    policy_dict = PolicyRuleDict()
    arg_list = {'ip_proto': args.ip_proto,
                'ipv4_dst': args.ipv4_dst,
                'ipv4_src': args.ipv4_src,
                'tcp_dst': args.tcp_dst,
                'tcp_src': args.tcp_src,
                'udp_dst': args.udp_dst,
                'udp_src': args.udp_src,
                'direction': args.direction}
    match = FlowMatch(**arg_list)
    flow = FlowDescription(match=match, action=args.action)
    rule = policy_dict.get(rule_id)
    if not rule or args.overwrite:
        action = 'add'
        rule = PolicyRule(id=rule_id, flow_list=[flow])
    else:
        action = 'edit'
        rule.flow_list.extend([flow])
    policy_dict[rule_id] = rule
    print("Rule '%s' successfully %sed!" % (rule_id, action))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for policydb',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_add = subparsers.add_parser('add_rule', help='Add rule')
    parser_add.add_argument('rule_id', help='rule id to add')
    parser_add.add_argument('-ipd', '--ipv4_dst', help='ipv4 dst for rule')
    parser_add.add_argument('-ips', '--ipv4_src', help='ipv4 src for rule')
    parser_add.add_argument('-p', '--ip_proto', help='ip proto for rule',
                            type=int)
    parser_add.add_argument('-td', '--tcp_dst', type=int,
                            help='tcp dst for rule, set -p to IPPROTO_TCP(6)')
    parser_add.add_argument('-ts', '--tcp_src', type=int,
                            help='tcp src for rule, set -p to IPPROTO_TCP(6)')
    parser_add.add_argument('-ud', '--udp_dst', type=int,
                            help='udp dst for rule, set -p to IPPROTO_UDP(17)')
    parser_add.add_argument('-us', '--udp_src', type=int,
                            help='udp src for rule, set -p to IPPROTO_UDP(17)')
    parser_add.add_argument('-d', '--direction', type=int,
                            help='0 == UPLINK, 1 == DOWNLINK')
    parser_add.add_argument('-a', '--action', type=int,
                            help='0 == PERMIT, 1 == DENY')
    parser_add.add_argument('-o', '--overwrite', action='store_true',
                            help='overwrite existing rule')

    parser_add.set_defaults(func=add_rule)
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
