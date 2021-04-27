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
import errno
import random
import subprocess
import time
from collections import namedtuple
from datetime import datetime
from pprint import pprint

from lte.protos.pipelined_pb2 import (
    ActivateFlowsRequest,
    DeactivateFlowsRequest,
    DeactivateFlowsResult,
    RequestOriginType,
    RuleModResult,
    SubscriberQuotaUpdate,
    UEMacFlowRequest,
    UpdateSubscriberQuotaStateRequest,
    VersionedPolicy,
    VersionedPolicyID,
)
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    PolicyRule,
    RedirectInformation,
)
from lte.protos.subscriberdb_pb2 import AggregatedMaximumBitrate
from magma.common.rpc_utils import grpc_wrapper
from magma.configuration.service_configs import load_service_config
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from magma.pipelined.qos.common import QosManager
from magma.pipelined.service_manager import Tables
from magma.subscriberdb.sid import SIDUtils
from scripts.helpers.ng_set_session_msg import CreateSessionUtil
from scripts.helpers.pg_set_session_msg import CreateMMESessionUtils
from orc8r.protos.common_pb2 import Void

UEInfo = namedtuple('UEInfo', ['imsi_str', 'ipv4_src', 'ipv4_dst',
                               'uplink_tunnel', 'rule_id'])

def _gen_ue_set(num_of_ues):
    imsi = 123000000
    uplink_tunnel = 0x12345
    ue_set = set()
    for _ in range(0, num_of_ues):
        imsi_str = "IMSI" + str(imsi)
        ipv4_src = ".".join(str(random.randint(0, 255)) for _ in range(4))
        ipv4_dst = ".".join(str(random.randint(0, 255)) for _ in range(4))
        rule_id = "allow." + imsi_str
        uplink_tunnel = uplink_tunnel
        ue_set.add(UEInfo(imsi_str, ipv4_src, ipv4_dst, uplink_tunnel, rule_id))
        imsi += 1
        uplink_tunnel += 1
    return ue_set


@grpc_wrapper
def set_smf_session(client, args):
    cls_sess = CreateSessionUtil(args.subscriber_id, args.session_id, args.version)

    cls_sess.CreateSession(args.subscriber_id, args.pdr_state, args.in_teid, args.out_teid,
                           args.ue_ip_addr, args.gnb_ip_addr,
                           args.del_rule_id, args.add_rule_id, args.ipv4_dst,
                           args.allow, args.priority)

    print (cls_sess._set_session)
    response = client.SetSMFSessions(cls_sess._set_session)
    print (response)


@grpc_wrapper
def set_mme_session(client, args):
    cls_sess = CreateMMESessionUtils(args.imsi, args.priority, args.ue_ipv4_addr,
                                     args.ue_ipv6_addr, args.enb_ip_addr, args.apn,
                                     args.vlan, args.in_teid, args.out_teid,
                                     args.ue_state, args.flow_dl)

    print(cls_sess._set_pg_session)
    response = client.UpdateUEState(cls_sess._set_pg_session)

# --------------------------
# Enforcement App
# --------------------------

@grpc_wrapper
def deactivate_flows(client, args):
    policies = [VersionedPolicyID(rule_id=rule_id, version=1) for rule_id
                in args.rule_ids.split(',') if args.rule_ids]
    request = DeactivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        uplink_tunnel=args.uplink_tunnel,
        policies=policies,
        request_origin=RequestOriginType(type=RequestOriginType.GX))
    client.DeactivateFlows(request)


@grpc_wrapper
def activate_flows(client, args):
    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        uplink_tunnel=args.uplink_tunnel,
        policies=[VersionedPolicy(
            rule=PolicyRule(
                id=args.rule_id,
                priority=args.priority,
                hard_timeout=args.hard_timeout,
                flow_list=[
                    FlowDescription(match=FlowMatch(
                        ip_dst=convert_ipv4_str_to_ip_proto(args.ipv4_dst),
                        direction=FlowMatch.UPLINK)),
                    FlowDescription(match=FlowMatch(
                        ip_src=convert_ipv4_str_to_ip_proto(args.ipv4_dst),
                        direction=FlowMatch.DOWNLINK)),
                ],
            ),
            version=1)],
        request_origin=RequestOriginType(type=RequestOriginType.GX))
    response = client.ActivateFlows(request)
    _print_rule_mod_results(response.policy_results)


@grpc_wrapper
def activate_gy_redirect(client, args):
    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        uplink_tunnel=args.uplink_tunnel,
        policies=[VersionedPolicy(
            rule=PolicyRule(
                id=args.rule_id,
                priority=999,
                flow_list=[],
                redirect=RedirectInformation(
                    support=1,
                    address_type=2,
                    server_address=args.redirect_addr
                )
            ),
            version=1)],
        request_origin=RequestOriginType(type=RequestOriginType.GY))
    response = client.ActivateFlows(request)
    _print_rule_mod_results(response.policy_results)


@grpc_wrapper
def deactivate_gy_flows(client, args):
    policies = [VersionedPolicyID(rule_id=rule_id, version=1) for rule_id
                in args.rule_ids.split(',') if args.rule_ids]
    request = DeactivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        uplink_tunnel=args.uplink_tunnel,
        policies=policies,
        request_origin=RequestOriginType(type=RequestOriginType.GY))
    client.DeactivateFlows(request)


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


@grpc_wrapper
def stress_test_grpc(client, args):
    print("WARNING: DO NOT USE ON PRODUCTION SETUPS")
    delta_time = 1/args.attaches_per_sec
    print("Attach every ~{0} seconds".format(delta_time))

    if args.disable_qos:
        print("QOS Disabled")
        apn_ambr = None
    else:
        print("QOS Enabled")
        apn_ambr = AggregatedMaximumBitrate(
            max_bandwidth_ul=1000000000,
            max_bandwidth_dl=1000000000,
        )

    for i in range (0, args.test_iterations):
        print("Starting iteration {0} of attach/detach requests".format(i))
        ue_dict = _gen_ue_set(args.num_of_ues)
        print("Starting attaches")

        timestamp = datetime.now()
        for ue in ue_dict:
            grpc_start_timestamp = datetime.now()
            request = ActivateFlowsRequest(
                sid=SIDUtils.to_pb(ue.imsi_str),
                ip_addr=ue.ipv4_src,
                uplink_tunnel=ue.uplink_tunnel,
                policies=[VersionedPolicy(
                    rule=PolicyRule(
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
                    ),
                    version=1)
                    ],
                request_origin=RequestOriginType(type=RequestOriginType.GX),
                apn_ambr=apn_ambr,
            )
            response = client.ActivateFlows(request)
            if any(r.result != RuleModResult.SUCCESS for
                   r in response.policy_results):
                _print_rule_mod_results(response.policy_results)

            grpc_end_timestamp = datetime.now()
            call_duration = (grpc_end_timestamp - grpc_start_timestamp).total_seconds()
            if call_duration < delta_time:
                time.sleep(delta_time - call_duration)

        duration = (datetime.now() - timestamp).total_seconds()
        print("Finished {0} attaches in {1} seconds".format(len(ue_dict),
                                                            duration))
        print("Actual attach rate = {0} UEs per sec".format(round(len(ue_dict)/duration)))

        time.sleep(args.time_between_detach)

        print("Starting detaches")
        timestamp = datetime.now()
        for ue in ue_dict:
            grpc_start_timestamp = datetime.now()
            request = DeactivateFlowsRequest(
                sid=SIDUtils.to_pb(ue.imsi_str),
                ip_addr=ue.ipv4_src,
                uplink_tunnel=ue.uplink_tunnel,
                policies=[
                    VersionedPolicyID(
                        rule_id=ue.rule_id,
                        version=1)
                ],
                request_origin=RequestOriginType(type=RequestOriginType.GX),
                remove_default_drop_flows=True)
            response = client.DeactivateFlows(request)
            if response.result != DeactivateFlowsResult.SUCCESS:
                _print_rule_mod_results(response.policy_results)

            grpc_end_timestamp = datetime.now()
            call_duration = (grpc_end_timestamp - grpc_start_timestamp).total_seconds()
            if call_duration < delta_time:
                time.sleep(delta_time - call_duration)

        duration = (datetime.now() - timestamp).total_seconds()
        print("Finished {0} detaches in {1} seconds".format(len(ue_dict),
                                                            duration))
        print("Actual detach rate = {0} UEs per sec",
              round(len(ue_dict)/duration))

def create_ng_services_parser(apps):
    """
    Creates the argparse subparser for the ng_services app
    """
    app = apps.add_parser('ng_services')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    subcmd = subparsers.add_parser('set_smf_session',
                                   help='SMF set Session Emulator')
    subcmd.add_argument('--subscriber_id', help='Subscriber Identity', default='IMSI12345')
    subcmd.add_argument('--session_id', help='Session Identity', type=int, default=100)
    subcmd.add_argument('--version', help='Session Version', type=int, default=2)
    subcmd.add_argument('--pdr_state', help='ADD / IDLE / REMOVE the PDR',
                        default="ADD")
    subcmd.add_argument('--in_teid', help='Match incoming teid from access',
                         type=int, default=0)
    subcmd.add_argument('--out_teid', help='Put outgoing teid towards access',
                         type=int, default=0)
    subcmd.add_argument('--ue_ip_addr', help='UE IP address ',
                         default='')
    subcmd.add_argument('--gnb_ip_addr', help='IP address of GNB Node',
                         default='')
    subcmd.add_argument('--del_rule_id', help='rule id to add', default='')
    subcmd.add_argument('--add_rule_id', help='rule id to add', default='')
    subcmd.add_argument('--ipv4_dst', help='ipv4 dst for rule', default='')
    subcmd.add_argument('--allow', help='YES/NO for allow and deny', default='YES')
    subcmd.add_argument('--priority', help='priority for rule',
                        type=int, default=0)
    subcmd.add_argument('--hard_timeout', help='hard timeout for rule',
                        type=int, default=0)

    subcmd.set_defaults(func=set_smf_session)

def create_pg_services_parser(apps):
    """
    Creates the argparse subparser for the pg_services app
    pg refers to services from MME to PIPELINED
    """
    app = apps.add_parser('pg_services')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    subcmd = subparsers.add_parser('set_mme_session',
                                   help='MME set Session Emulator')
    subcmd.add_argument('--imsi', help='Subscriber Identity', default='IMSI12345')
    subcmd.add_argument('--priority', help='priority for rule',
                        type=int, default=10)
    subcmd.add_argument('--ue_ipv4_addr', help='UE IPv4 address ',
                         default='192.168.128.11')
    subcmd.add_argument('--ue_ipv6_addr', help='UE IPv6 address ',
                         default='')
    subcmd.add_argument('--enb_ip_addr', help='IP address of ENB Node',
                         default='192.168.60.141')
    subcmd.add_argument('--apn', help='APN for accessing net',
                        default="magma.com")
    subcmd.add_argument('--vlan', help='Vlan Configuration for out ports',
                         type=int, default=0)
    subcmd.add_argument('--in_teid', help='Match incoming teid from access',
                         type=int, default=100)
    subcmd.add_argument('--out_teid', help='Put outgoing teid towards access',
                         type=int, default=200)
    subcmd.add_argument('--ue_state', help='ADD/DEL/ADD_IDLE/DEL_IDLE/SUSPENDED/RESUME',
                         default='ACTIVE')
    subcmd.add_argument('--flow_dl', help='ENABLE/DISABLE flow dl', default='DISABLE')
    subcmd.set_defaults(func=set_mme_session)

def create_enforcement_parser(apps):
    """
    Creates the argparse subparser for the enforcement app
    """
    app = apps.add_parser('enforcement')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('activate_flows',
                                   help='Activate flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--uplink_tunnel', help='Subscriber Uplink Tunnel ID',
                        default=0x12345)
    subcmd.add_argument('--rule_id', help='rule id to add', default='rule1')
    subcmd.add_argument('--ipv4_dst', help='ipv4 dst for rule', default='')
    subcmd.add_argument('--priority', help='priority for rule',
                        type=int, default=0)
    subcmd.add_argument('--hard_timeout', help='hard timeout for rule',
                        type=int, default=0)
    subcmd.set_defaults(func=activate_flows)

    subcmd = subparsers.add_parser('deactivate_flows', help='Deactivate flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--uplink_tunnel', help='Subscriber Uplink Tunnel ID',
                        default=0x12345)
    subcmd.add_argument('--rule_ids', help='Comma separated rule ids',
                        default="")
    subcmd.set_defaults(func=deactivate_flows)

    subcmd = subparsers.add_parser('activate_gy_redirect',
                                   help='Activate gy final action redirect')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--uplink_tunnel', help='Subscriber Uplink Tunnel ID',
                        default=0x12345)
    subcmd.add_argument('--rule_id', help='rule id to add', default='redirect')
    subcmd.add_argument('--redirect_addr', help='Webpage to redirect to',
                        default='http://about.sha.ddih.org/')
    subcmd.set_defaults(func=activate_gy_redirect)

    subcmd = subparsers.add_parser('deactivate_gy_flows',
                                   help='Deactivate gy flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--uplink_tunnel', help='Subscriber Uplink Tunnel ID',
                        default=0x12345)
    subcmd.add_argument('--rule_ids', help='Comma separated rule ids',
                        default="")
    subcmd.set_defaults(func=deactivate_gy_flows)

    subcmd = subparsers.add_parser('display_flows',
                                   help='Display flows related to policy '
                                        'enforcement')
    subcmd.set_defaults(func=display_enforcement_flows)

    subcmd = subparsers.add_parser('get_policy_usage',
                                   help='Get policy usage stats')
    subcmd.set_defaults(func=get_policy_usage)

    subcmd = subparsers.add_parser('stress_test_grpc',
        help='Sends a set of Activate grpc requests, followed by Deactivates')
    subcmd.add_argument('--attaches_per_sec',
                        help='Number of grpc Attach requests per second',
                        type=int, default=10)
    subcmd.add_argument('--num_of_ues', help='Number of total UEs to atach',
                        type=int, default=600)
    subcmd.add_argument('--time_between_detach',
                        help='Time between attaches and detaches in seconds',
                        type=int, default=10)
    subcmd.add_argument('--test_iterations', help='Test duration in seconds',
                        type=int, default=5)
    subcmd.add_argument('--disable_qos', help='If we want to disable QOS',
                        action="store_true")
    subcmd.set_defaults(func=stress_test_grpc)

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

    subcmd = subparsers.add_parser('qos', help='Debug Qos')
    subcmd.set_defaults(func=QosManager.debug)

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
    create_pg_services_parser(apps)
    create_ng_services_parser(apps)
    create_enforcement_parser(apps)
    create_ue_mac_parser(apps)
    create_check_flows_parser(apps)
    create_debug_parser(apps)
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
