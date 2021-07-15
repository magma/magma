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
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.policydb_pb2 import (
    DisableStaticRuleRequest,
    EnableStaticRuleRequest,
    FlowDescription,
    FlowMatch,
    PolicyRule,
)
from lte.protos.policydb_pb2_grpc import PolicyDBStub
from magma.common.rpc_utils import grpc_wrapper
from magma.policydb.rule_store import PolicyRuleDict

DEBUG_MSG = 'You may want to check that a connection can be made to ' \
            'orc8r to update the assignments of rules/basenames to ' \
            'the subscriber.'


@grpc_wrapper
def add_rule(args):
    rule_id = args.rule_id
    policy_dict = PolicyRuleDict()
    arg_list = {
        'ip_proto': args.ip_proto,
        'ip_dst': IPAddress(
            version=IPAddress.IPV4,
            address=args.ipv4_dst.encode('utf-8'),
        ),
        'ip_src': IPAddress(
            version=IPAddress.IPV4,
            address=args.ipv4_src.encode('utf-8'),
        ),
        'tcp_dst': args.tcp_dst,
        'tcp_src': args.tcp_src,
        'udp_dst': args.udp_dst,
        'udp_src': args.udp_src,
        'direction': args.direction,
    }
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


@grpc_wrapper
def install_rules(client: PolicyDBStub, args):
    """
    Installs the specified static rules for the session.
    Also associates the static rules to the subscriber.
    """
    message = EnableStaticRuleRequest(
        imsi=args.imsi,
        rule_ids=args.rules,
        base_names=args.basenames,
    )
    try:
        client.EnableStaticRules(message)
    except grpc.RpcError as err:
        print('Failed to enable static rules and/or base names: %s' % str(err))
        print(DEBUG_MSG)


@grpc_wrapper
def uninstall_rules(client: PolicyDBStub, args):
    """
    Uninstalls the specified static rules for the session.
    Also unassociates the static rules from the subscriber.
    """
    message = DisableStaticRuleRequest(
        imsi=args.imsi,
        rule_ids=args.rules,
        base_names=args.basenames,
    )
    try:
        client.DisableStaticRules(message)
    except grpc.RpcError as err:
        print('Failed to disable static rules and/or base names: %s' % str(err))
        print(DEBUG_MSG)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for policydb',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_add = subparsers.add_parser('add_rule', help='Add rule')
    parser_add.add_argument('rule_id', help='rule id to add')
    parser_add.add_argument('-ipd', '--ipv4_dst', help='ipv4 dst for rule')
    parser_add.add_argument('-ips', '--ipv4_src', help='ipv4 src for rule')
    parser_add.add_argument(
        '-p', '--ip_proto', help='ip proto for rule',
        type=int,
    )
    parser_add.add_argument(
        '-td', '--tcp_dst', type=int,
        help='tcp dst for rule, set -p to IPPROTO_TCP(6)',
    )
    parser_add.add_argument(
        '-ts', '--tcp_src', type=int,
        help='tcp src for rule, set -p to IPPROTO_TCP(6)',
    )
    parser_add.add_argument(
        '-ud', '--udp_dst', type=int,
        help='udp dst for rule, set -p to IPPROTO_UDP(17)',
    )
    parser_add.add_argument(
        '-us', '--udp_src', type=int,
        help='udp src for rule, set -p to IPPROTO_UDP(17)',
    )
    parser_add.add_argument(
        '-d', '--direction', type=int,
        help='0 == UPLINK, 1 == DOWNLINK',
    )
    parser_add.add_argument(
        '-a', '--action', type=int,
        help='0 == PERMIT, 1 == DENY',
    )
    parser_add.add_argument(
        '-o', '--overwrite', action='store_true',
        help='overwrite existing rule',
    )

    parser_install = subparsers.add_parser(
        'install_rules', help='Install static rules for a subscriber',
    )
    parser_install.add_argument('-id', help='session id')
    parser_install.add_argument('-imsi', help='subscriber IMSI')
    parser_install.add_argument(
        '-rules', nargs='*',
        help='static rules to install',
    )
    parser_install.add_argument(
        '-basenames', nargs='*',
        help='basenames to install',
    )

    parser_uninstall = subparsers.add_parser(
        'uninstall_rules', help='Uninstall static rules for a subscriber',
    )
    parser_uninstall.add_argument('-id', help='session id')
    parser_uninstall.add_argument('-imsi', help='subscriber IMSI')
    parser_uninstall.add_argument(
        '-rules', nargs='*',
        help='static rules to uninstall',
    )
    parser_uninstall.add_argument(
        '-basenames', nargs='*',
        help='basenames to install',
    )

    # Add function callbacks
    parser_add.set_defaults(func=add_rule)
    parser_install.set_defaults(func=install_rules)
    parser_uninstall.set_defaults(func=uninstall_rules)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, PolicyDBStub, 'policydb')


if __name__ == "__main__":
    main()
