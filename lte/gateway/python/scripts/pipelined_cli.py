#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse
import errno
from pprint import pprint
import subprocess

from magma.common.rpc_utils import grpc_wrapper
from lte.protos.pipelined_pb2 import (
    SubscriberQuotaUpdate,
    UpdateSubscriberQuotaStateRequest,
)
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.subscriberdb.sid import SIDUtils
from magma.configuration.service_configs import load_service_config
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.service_manager import Tables
from orc8r.protos.common_pb2 import Void
from lte.protos.pipelined_pb2 import (
    ActivateFlowsRequest,
    DeactivateFlowsRequest,
    RuleModResult,
    UEMacFlowRequest,
    ActivateAPNTaggingFlowsForUserRequest,
    DeactivateAPNTaggingFlowsForUserRequest,
)
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from lte.protos.policydb_pb2 import FlowMatch, FlowDescription, PolicyRule


# --------------------------
# Enforcement App
# --------------------------

@grpc_wrapper
def activate_flows(client, args):
    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        rule_ids=args.rule_ids.split(','))
    response = client.ActivateFlows(request)
    _print_rule_mod_results(response.static_rule_results)


@grpc_wrapper
def deactivate_flows(client, args):
    request = DeactivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        rule_ids=args.rule_ids.split(',') if args.rule_ids else [])
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
    response = client.ActivateFlows(request)
    _print_rule_mod_results(response.dynamic_rule_results)


def _print_rule_mod_results(results):
    # The message cannot be directly printed because SUCCESS is mapped to 0,
    # which is ignored in the printing by default.
    for result in results:
        print(result.rule_id,
              RuleModResult.Result.Name(result.result))


@grpc_wrapper
def display_enforcement_flows(client, _):
    _display_flows(client, [EnforcementController.APP_NAME,
                            EnforcementStatsController.APP_NAME])


@grpc_wrapper
def get_policy_usage(client, _):
    rule_table = client.GetPolicyUsage(Void())
    pprint(rule_table)


def create_enforcement_parser(apps):
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
    subcmd.add_argument('--rule_ids', help='Comma separated rule ids')
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

    subcmd = subparsers.add_parser('display_flows',
                                   help='Display flows related to policy '
                                        'enforcement')
    subcmd.set_defaults(func=display_enforcement_flows)

    subcmd = subparsers.add_parser('get_policy_usage',
                                   help='Get policy usage stats')
    subcmd.set_defaults(func=get_policy_usage)


# -------------
# UE MAC APP
# -------------

@grpc_wrapper
def add_ue_mac_flow(client, args):
    request = UEMacFlowRequest(
        sid=SIDUtils.to_pb(args.imsi),
        mac_addr=args.mac
    )
    res = client.AddUEMacFlow(request)
    if res is None:
        print("Error associating MAC to IMSI")


@grpc_wrapper
def delete_ue_mac_flow(client, args):
    request = UEMacFlowRequest(
        sid=SIDUtils.to_pb(args.imsi),
        mac_addr=args.mac
    )
    res = client.DeleteUEMacFlow(request)
    if res is None:
        print("Error associating MAC to IMSI")


def create_ue_mac_parser(apps):
    """
    Creates the argparse subparser for the MAC App
    """
    app = apps.add_parser('ue_mac')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('add_ue_mac_flow',
                                   help='Add flow to match UE MAC \
                                   with a subscriber')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--mac', help='UE MAC address',
                        default='5e:cc:cc:b1:49:ff')
    subcmd.set_defaults(func=add_ue_mac_flow)
    # Delete subcommands
    subcmd = subparsers.add_parser('delete_ue_mac_flow',
                                   help='Delete flow to match UE MAC \
                                   with a subscriber')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--mac', help='UE MAC address',
                        default='5e:cc:cc:b1:49:ff')
    subcmd.set_defaults(func=delete_ue_mac_flow)


# -------------
# Check Quota APP
# -------------

@grpc_wrapper
def update_quota(client, args):
    update = SubscriberQuotaUpdate(
        sid=SIDUtils.to_pb(args.imsi),
        mac_addr=args.mac,
        update_type=args.update_type
    )
    request = UpdateSubscriberQuotaStateRequest(updates=[update],)
    res = client.UpdateSubscriberQuotaState(request)
    if res is None:
        print("Error updating check quota flows")


def create_check_flows_parser(apps):
    """
    Creates the argparse subparser for the MAC App
    """
    app = apps.add_parser('check_quota')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('update_quota',
                                   help='Add flow to match UE MAC \
                                   with a subscriber')
    subcmd.add_argument('imsi', help='Subscriber ID')
    subcmd.add_argument('mac', help='Subscriber mac')
    subcmd.add_argument('update_type', type=int,
                        help='0 - valid quota, 1 -no quota, 2 - terminate')
    subcmd.set_defaults(func=update_quota)


# --------------------------
# Debugging
# --------------------------

@grpc_wrapper
def get_table_assignment(client, args):
    response = client.GetAllTableAssignments(Void())
    table_assignments = response.table_assignments
    if args.apps:
        app_filter = args.apps.split(',')
        table_assignments = [table_assignment for table_assignment in
                             table_assignments if
                             table_assignment.app_name in app_filter]

    table_template = '{:<25}{:<20}{:<25}'
    print(table_template.format('App', 'Main Table', 'Scratch Tables'))
    print('-' * 70)
    for table_assignment in table_assignments:
        print(table_template.format(
            table_assignment.app_name,
            table_assignment.main_table,
            str([table for table in table_assignment.scratch_tables])))


@grpc_wrapper
def display_raw_flows(_unused, args):
    pipelined_config = load_service_config('pipelined')
    bridge_name = pipelined_config['bridge_name']
    try:
        flows = BridgeTools.get_flows_for_bridge(bridge_name, args.table_num)
    except subprocess.CalledProcessError as e:
        if e.returncode == errno.EPERM:
            print("Need to run as root to dump flows")
        return

    for flow in flows:
        print(flow)


def _display_flows(client, apps=None):
    pipelined_config = load_service_config('pipelined')
    bridge_name = pipelined_config['bridge_name']
    response = client.GetAllTableAssignments(Void())
    table_assignments = {
        table_assignment.app_name:
            Tables(main_table=table_assignment.main_table, type=None,
                   scratch_tables=table_assignment.scratch_tables)
        for table_assignment in response.table_assignments}
    try:
        flows = BridgeTools.get_annotated_flows_for_bridge(
            bridge_name, table_assignments, apps)
    except subprocess.CalledProcessError as e:
        if e.returncode == errno.EPERM:
            print("Need to run as root to dump flows")
        return

    for flow in flows:
        print(flow)


@grpc_wrapper
def display_flows(client, args):
    if args.apps is None:
        _display_flows(client)
        return
    _display_flows(client, args.apps.split(','))


def create_debug_parser(apps):
    """
    Creates the argparse subparser for the debugging commands
    """
    app = apps.add_parser('debug')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('table_assignment',
                                   help='Get the table assignment for apps.')
    subcmd.add_argument('--apps',
                        help='Comma separated list of app names. If not set, '
                             'all table assignments will be printed.')
    subcmd.set_defaults(func=get_table_assignment)

    subcmd = subparsers.add_parser('display_raw_flows',
                                   help='Display raw flows from ovs dump')
    subcmd.add_argument('--table_num', help='Table number to filter the flows.'
                                            'If not set, all flows will be '
                                            'printed')
    subcmd.set_defaults(func=display_raw_flows)

    subcmd = subparsers.add_parser('display_flows', help='Display flows')
    subcmd.add_argument('--apps',
                        help='Comma separated list of app names to filter the'
                             'flows. If not set, all flows will be printed.')
    subcmd.set_defaults(func=display_flows)



# -------------
# APN APP
# -------------

@grpc_wrapper
def add_tagging_for_user(client, args):
    request = ActivateAPNTaggingFlowsForUserRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ip,
        ap_name=args.apn
    )
    res = client.AddAPNFlowForUser(request)
    if res is None:
        print("Error creating APN tagging flows")

@grpc_wrapper
def delete_tagging_for_user(client, args):
    request = DeactivateAPNTaggingFlowsForUserRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ip,
        ap_name=args.apn
    )
    res = client.DeleteAPNFlowForUser(request)
    if res is None:
        print("Error deleting APN tagging flows")

def create_apn_app_parser(apps):
    """
    Creates the argparse subparser for the APN App
    """
    app = apps.add_parser('apn')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('add_tagging_for_user',
                                   help='Add flow to match to mark each packet with APN')
    subcmd.add_argument('imsi', help='Subscriber ID')
    subcmd.add_argument('ip', help='IP address allocated for the UE')
    subcmd.add_argument('apn', help='APN UE connected to')
    subcmd.set_defaults(func=add_tagging_for_user)

    subcmd = subparsers.add_parser('delete_tagging_for_user',
                                   help='Add flow to match to mark each packet with APN')
    subcmd.add_argument('imsi', help='Subscriber ID')
    subcmd.add_argument('ip', help='IP address allocated for the UE')
    subcmd.add_argument('apn', help='APN UE connected to')
    subcmd.set_defaults(func=delete_tagging_for_user)


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
    create_enforcement_parser(apps)
    create_ue_mac_parser(apps)
    create_check_flows_parser(apps)
    create_debug_parser(apps)
    create_apn_app_parser(apps)
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
