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
from pprint import pprint
import subprocess
from collections import namedtuple

from magma.common.rpc_utils import grpc_wrapper
from lte.protos.pipelined_pb2 import (
    SubscriberQuotaUpdate,
    UpdateSubscriberQuotaStateRequest,
)
from lte.protos.policydb_pb2 import RedirectInformation
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from magma.subscriberdb.sid import SIDUtils
from magma.configuration.service_configs import load_service_config
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.service_manager import Tables
from magma.pipelined.qos.common import QosManager
from orc8r.protos.common_pb2 import Void
from lte.protos.session_manager_pb2 import NodeID
from lte.protos.pipelined_pb2 import (
    SessionSet,
    SetGroupFAR,
    FwdParam,
    Action,
    OuterHeaderCreation,
    SetGroupPDR,
    PDI,
    OuterHeadRemoval,
    Fsm_state,
    PdrState,
    ActivateFlowsRequest,
    DeactivateFlowsRequest,
    RuleModResult,
    UEMacFlowRequest,
    RequestOriginType,
)
from magma.pipelined.ng_manager.session_state_manager_util import FARRuleEntry

from lte.protos.pipelined_pb2_grpc import PipelinedStub
from lte.protos.policydb_pb2 import FlowMatch, FlowDescription, PolicyRule

QoSEnforceRuleEntry = namedtuple(
                         'QoSEnforceRuleEntry',
                         ['imsi', 'rule_id', 'ipv4_dst', 'allow', 'priority', 'hard_timeout', 'direction'])

# --------------------------
# NG Servicer App
#---------------------------
class CreateSessionUtil:

    def __init__(self, subscriber_id:str, session_version, node_id="192.168.220.1"):
        self._set_session = \
                  SessionSet(subscriber_id=subscriber_id, session_version=session_version,\
                             node_id=NodeID(node_id_type=NodeID.IPv4, node_id=node_id),\
                             state=Fsm_state(state=Fsm_state.CREATED))


    def CreateQERinPDR(self, qos_enforce_rule: QoSEnforceRuleEntry, ue_ip_addr: str):
        if qos_enforce_rule.allow == 'YES':
           allow = FlowDescription.PERMIT
        else:
           allow = FlowDescription.DENY

        if qos_enforce_rule.direction == FlowMatch.UPLINK:
            flow_list =  [FlowDescription(match=FlowMatch(
                                         ip_dst=convert_ipv4_str_to_ip_proto(qos_enforce_rule.ipv4_dst),
                                         direction=qos_enforce_rule.direction),
                                         action=allow)]
        else:
            flow_list =  [FlowDescription(match=FlowMatch(
                                         ip_src=convert_ipv4_str_to_ip_proto(qos_enforce_rule.ipv4_dst),
                                         direction=qos_enforce_rule.direction),
                                         action=allow)]


        qos_enforce_rule = ActivateFlowsRequest(
                                  sid=SIDUtils.to_pb(qos_enforce_rule.imsi),
                                  ip_addr=ue_ip_addr,
                                  dynamic_rules=[PolicyRule(
                                  id=qos_enforce_rule.rule_id,
                                  priority=qos_enforce_rule.priority,
                                  hard_timeout=qos_enforce_rule.hard_timeout,
                                  flow_list=flow_list
                                )],
                                request_origin=RequestOriginType(type=RequestOriginType.N4))
        return  qos_enforce_rule

    def CreateDelQERinPDR(self, qos_enforce_rule: QoSEnforceRuleEntry, ue_ip_addr: str):

        qos_enforce_rule = DeactivateFlowsRequest(
                                  sid=SIDUtils.to_pb(qos_enforce_rule.imsi),
                                  ip_addr=ue_ip_addr,
                                  rule_ids=[qos_enforce_rule.rule_id],
                                  request_origin=RequestOriginType(type=RequestOriginType.N4))
        return  qos_enforce_rule


    def CreateSessionPDR(self, pdr_id:int, pdr_version:int, pdr_state,
                         precedence:int, local_f_teid:int, ue_ip_addr:str,
                         del_qos_enforcer=None, add_qos_enforcer=None,
                         far_tuple=None):

        del_qos_enforce_rule = None
        add_qos_enforce_rule = None
        far_gr_entry = None

        if far_tuple:
            if far_tuple.o_teid != 0:
                # For pdr_id=2 towards access
                far_gr_entry = SetGroupFAR(far_action_to_apply=[far_tuple.apply_action],\
                                           fwd_parm=FwdParam(dest_iface=0, \
                                           outr_head_cr=OuterHeaderCreation(\
                                             o_teid=far_tuple.o_teid, ipv4_adr=far_tuple.enodeb_ip_addr)))
            else:
                far_gr_entry = SetGroupFAR(far_action_to_apply=[far_tuple.apply_action])

        if del_qos_enforcer and len(del_qos_enforcer.rule_id):
            del_qos_enforce_rule = self.CreateDelQERinPDR(del_qos_enforcer, ue_ip_addr)

        if add_qos_enforcer and len(add_qos_enforcer.rule_id):
            add_qos_enforce_rule = self.CreateQERinPDR(add_qos_enforcer, ue_ip_addr)

        if local_f_teid != 0:
            self._set_session.set_gr_pdr.extend([\
                         SetGroupPDR(pdr_id=pdr_id, pdr_version=pdr_version,
                                     pdr_state=pdr_state,\
                                     precedence=precedence,\
                                     pdi=PDI(src_interface=0,\
                                              local_f_teid=local_f_teid,\
                                              ue_ip_adr=ue_ip_addr), \
                                      o_h_remo_desc=0, \
                                      set_gr_far=far_gr_entry,
                                      deactivate_flow_req=del_qos_enforce_rule,
                                      activate_flow_req=add_qos_enforce_rule)])


        else:
            self._set_session.set_gr_pdr.extend([\
                         SetGroupPDR(pdr_id=pdr_id, pdr_version=pdr_version,
                                     pdr_state=pdr_state,\
                                     precedence=precedence,\
                                     pdi=PDI(src_interface=1, ue_ip_adr=ue_ip_addr),\
                                     set_gr_far=far_gr_entry,
                                     deactivate_flow_req=del_qos_enforce_rule,
                                     activate_flow_req=add_qos_enforce_rule)])

@grpc_wrapper
def set_smf_session(client, args):

    cls_sess = CreateSessionUtil(args.subscriber_id, args.version)

    imsi_val = args.subscriber_id

    if args.pdr_state == "ADD":
        # From Access towards Core pdrid, version, precedence, teid, ue-ip, farid
        cls_sess.CreateSessionPDR(1, 1, PdrState.Value('INSTALL'), 32, args.in_teid, args.ue_ip_addr,\
                                  #For deletion
                                  QoSEnforceRuleEntry(imsi_val, args.del_rule_id, None, \
                                                      None, None, None, None),\
                                  #For addition
                                  QoSEnforceRuleEntry(imsi_val, args.add_rule_id, args.ipv4_dst,\
                                                      args.allow, args.priority, args.hard_timeout,
                                                      FlowMatch.UPLINK),
                                  FARRuleEntry(Action.Value('FORW'), 0, ''))

        # From Core towards access, pdrid, version, precedence, teid=0, ue-ip, farid
        cls_sess.CreateSessionPDR(2, 1, PdrState.Value('INSTALL'), 32, 0, args.ue_ip_addr,\
                                  #For deleteion
                                  QoSEnforceRuleEntry(imsi_val, args.del_rule_id, None,\
                                                      None, None, None, None),\
                                  #For addition
                                  QoSEnforceRuleEntry(imsi_val, args.add_rule_id, args.ipv4_dst,\
                                                      args.allow, args.priority, args.hard_timeout,
                                                      FlowMatch.DOWNLINK),
                                  FARRuleEntry(Action.Value('FORW'), args.out_teid, args.gnb_ip_addr))

    else:
        # Use ADD-Rules to delete the QoS along with PDR
        if args.add_rule_id:
            # Delete PDR from Acces toward CORE
            cls_sess.CreateSessionPDR(1, 1, PdrState.Value('REMOVE'), 32, args.in_teid, args.ue_ip_addr,\
                                      None,
                                      QoSEnforceRuleEntry(imsi_val, args.add_rule_id, args.ipv4_dst,\
                                                          args.allow, args.priority, args.hard_timeout,
                                                          FlowMatch.UPLINK),
                                      FARRuleEntry(Action.Value('FORW'), 0, ''))

            # Delete PDR from Core toward ACCESS
            cls_sess.CreateSessionPDR(2, 1, PdrState.Value('REMOVE'), 32, 0, args.ue_ip_addr,\
                                      None,
                                      QoSEnforceRuleEntry(imsi_val, args.add_rule_id, args.ipv4_dst,\
                                                          args.allow, args.priority, args.hard_timeout,
                                                          FlowMatch.DOWNLINK),
                                      FARRuleEntry(Action.Value('FORW'), args.out_teid, args.gnb_ip_addr))

        else:
            # Delete PDR from Acces toward CORE
            cls_sess.CreateSessionPDR(1, 1, PdrState.Value('REMOVE'), 32, args.in_teid, args.ue_ip_addr,\
                                      QoSEnforceRuleEntry(imsi_val, args.del_rule_id, None, None, None,\
                                                          None, None),
                                      None, FARRuleEntry(Action.Value('FORW'), 0, ''))
        
            # Delete PDR from Core toward ACCESS
            cls_sess.CreateSessionPDR(2, 1, PdrState.Value('REMOVE'), 32, 0, args.ue_ip_addr,\
                                      QoSEnforceRuleEntry(imsi_val, args.del_rule_id, None, None, None,\
                                                          None, None),
                                      None, FARRuleEntry(Action.Value('FORW'), args.out_teid, args.gnb_ip_addr))

    print (cls_sess._set_session)
    response = client.SetSMFSessions(cls_sess._set_session)
    print (response)

# --------------------------
# Enforcement App
# --------------------------

@grpc_wrapper
def activate_flows(client, args):
    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        rule_ids=args.rule_ids.split(','),
        request_origin=RequestOriginType(type=RequestOriginType.GX))
    response = client.ActivateFlows(request)
    _print_rule_mod_results(response.static_rule_results)


@grpc_wrapper
def deactivate_flows(client, args):
    request = DeactivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        rule_ids=args.rule_ids.split(',') if args.rule_ids else [],
        request_origin=RequestOriginType(type=RequestOriginType.GX))
    client.DeactivateFlows(request)


@grpc_wrapper
def activate_dynamic_rule(client, args):

    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        dynamic_rules=[PolicyRule(
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
        )],
        request_origin=RequestOriginType(type=RequestOriginType.GX))
    response = client.ActivateFlows(request)
    _print_rule_mod_results(response.dynamic_rule_results)


@grpc_wrapper
def activate_gy_redirect(client, args):
    request = ActivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        dynamic_rules=[PolicyRule(
            id=args.rule_id,
            priority=999,
            flow_list=[],
            redirect=RedirectInformation(
                support=1,
                address_type=2,
                server_address=args.redirect_addr
            )
        )],
        request_origin=RequestOriginType(type=RequestOriginType.GY))
    response = client.ActivateFlows(request)
    _print_rule_mod_results(response.dynamic_rule_results)


@grpc_wrapper
def deactivate_gy_flows(client, args):
    request = DeactivateFlowsRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ip_addr=args.ipv4,
        rule_ids=args.rule_ids.split(',') if args.rule_ids else [],
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


def create_ng_services_parser(apps):
    """
    Creates the argparse subparser for the ng_services app
    """
    app = apps.add_parser('ng_services')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    subcmd = subparsers.add_parser('set_smf_session',
                                   help='SMF set Session Emulator')
    subcmd.add_argument('--subscriber_id', help='Subscriber Identity', default='IMSI12345-100')
    subcmd.add_argument('--version', help='Session Version', type=int, default=2)
    subcmd.add_argument('--pdr_state', help='ADD/DEL the PDR',
                        default="ADD")
    subcmd.add_argument('--in_teid', help='Match incoming teid from access',
                         type=int, default=98)
    subcmd.add_argument('--out_teid', help='Put outgoing teid towards access',
                         type=int, default=1)
    subcmd.add_argument('--ue_ip_addr', help='UE IP address ',
                         default='60.60.60.1')
    subcmd.add_argument('--gnb_ip_addr', help='IP address of GNB Node',
                         default='100.200.200.1')
    subcmd.add_argument('--del_rule_id', help='rule id to add', default='')
    subcmd.add_argument('--add_rule_id', help='rule id to add', default='rule1')
    subcmd.add_argument('--ipv4_dst', help='ipv4 dst for rule', default='')
    subcmd.add_argument('--allow', help='YES/NO for allow and deny', default='YES')
    subcmd.add_argument('--priority', help='priority for rule',
                        type=int, default=0)
    subcmd.add_argument('--hard_timeout', help='hard timeout for rule',
                        type=int, default=0)

    subcmd.set_defaults(func=set_smf_session)


def create_enforcement_parser(apps):
    """
    Creates the argparse subparser for the enforcement app
    """
    app = apps.add_parser('enforcement')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    # Add subcommands
    subcmd = subparsers.add_parser('activate_flows', help='Activate flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--rule_ids',
                        help='Comma separated rule ids', default='rule1,rule2')
    subcmd.set_defaults(func=activate_flows)

    subcmd = subparsers.add_parser('deactivate_flows', help='Deactivate flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--rule_ids', help='Comma separated rule ids')
    subcmd.set_defaults(func=deactivate_flows)

    subcmd = subparsers.add_parser('activate_dynamic_rule',
                                   help='Activate dynamic flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--rule_id', help='rule id to add', default='rule1')
    subcmd.add_argument('--ipv4_dst', help='ipv4 dst for rule', default='')
    subcmd.add_argument('--priority', help='priority for rule',
                        type=int, default=0)
    subcmd.add_argument('--hard_timeout', help='hard timeout for rule',
                        type=int, default=0)
    subcmd.set_defaults(func=activate_dynamic_rule)

    subcmd = subparsers.add_parser('activate_gy_redirect',
                                   help='Activate gy final action redirect')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--rule_id', help='rule id to add', default='redirect')
    subcmd.add_argument('--redirect_addr', help='Webpage to redirect to',
                        default='http://about.sha.ddih.org/')
    subcmd.set_defaults(func=activate_gy_redirect)

    subcmd = subparsers.add_parser('deactivate_gy_flows',
                                   help='Deactivate gy flows')
    subcmd.add_argument('--imsi', help='Subscriber ID', default='IMSI12345')
    subcmd.add_argument('--ipv4', help='Subscriber IPv4', default='120.12.1.9')
    subcmd.add_argument('--rule_ids', help='Comma separated rule ids')
    subcmd.set_defaults(func=deactivate_gy_flows)

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
