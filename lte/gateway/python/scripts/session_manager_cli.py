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
from feg.protos.mock_core_pb2_grpc import MockOCSStub, MockPCRFStub
from lte.protos.abort_session_pb2 import AbortSessionRequest
from lte.protos.abort_session_pb2_grpc import AbortSessionResponderStub
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    FlowQos,
    PolicyRule,
    QosArp,
)
from lte.protos.session_manager_pb2 import (
    DynamicRuleInstall,
    LocalCreateSessionRequest,
    PolicyReAuthRequest,
    QoSInformation,
)
from lte.protos.session_manager_pb2_grpc import (
    LocalSessionManagerStub,
    SessionProxyResponderStub,
)
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.common.rpc_utils import grpc_wrapper
from magma.common.service_registry import ServiceRegistry
from magma.pipelined.tests.app.subscriber import SubContextConfig
from orc8r.protos.common_pb2 import Void


@grpc_wrapper
def send_create_session(client, args):
    sub1 = SubContextConfig("IMSI" + args.imsi, "192.168.128.74", 4)

    try:
        create_account_in_PCRF(args.imsi)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))

    try:
        create_account_in_OCS(args.imsi)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))

    req = LocalCreateSessionRequest(sid=SubscriberID(id=sub1.imsi), ue_ipv4=sub1.ip)
    print("Sending LocalCreateSessionRequest with following fields:\n %s" % req)
    try:
        client.CreateSession(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))

    req = SubscriberID(id=sub1.imsi)
    print("Sending EndSession with following fields:\n %s" % req)
    try:
        client.EndSession(req)
    except grpc.RpcError as e:
        print("gRPC failed with %s: %s" % (e.code(), e.details()))


def create_account_in_PCRF(imsi):
    pcrf_chan = ServiceRegistry.get_rpc_channel("pcrf", ServiceRegistry.CLOUD)
    pcrf_client = MockPCRFStub(pcrf_chan)

    print("Clearing accounts in PCRF")
    pcrf_client.ClearSubscribers(Void())

    print("Creating account in PCRF")
    pcrf_client.CreateAccount(SubscriberID(id=imsi))


def create_account_in_OCS(imsi):
    ocs_chan = ServiceRegistry.get_rpc_channel("ocs", ServiceRegistry.CLOUD)
    ocs_client = MockOCSStub(ocs_chan)

    print("Clearing accounts in OCS")
    ocs_client.ClearSubscribers(Void())

    print("Creating account in OCS")
    ocs_client.CreateAccount(SubscriberID(id=imsi))


@grpc_wrapper
def send_policy_rar(client, args):
    sessiond_chan = ServiceRegistry.get_rpc_channel("sessiond", ServiceRegistry.LOCAL)
    sessiond_client = SessionProxyResponderStub(sessiond_chan)
    flow_list_str = args.flow_rules.split(";")
    flow_match_list = []
    for i, flow_str in enumerate(flow_list_str):
        print("%d: %s" % (i, flow_str))
        flow_fields = flow_str.split(",")
        if flow_fields[0] == "UL":
            flow_direction = FlowMatch.UPLINK
        elif flow_fields[0] == "DL":
            flow_direction = FlowMatch.DOWNLINK
        else:
            print("%s is not valid" % flow_fields[0])
            raise ValueError(
                "UL or DL are the only valid"
                " values for first parameter of flow match",
            )
        ip_protocol = int(flow_fields[1])
        if flow_fields[1] == FlowMatch.IPPROTO_TCP:
            udp_src_port = 0
            udp_dst_port = 0
            if flow_fields[3]:
                tcp_src_port = int(flow_fields[3])
            else:
                tcp_src_port = 0
            if flow_fields[5]:
                tcp_dst_port = int(flow_fields[5])
            else:
                tcp_dst_port = 0
        elif flow_fields[1] == FlowMatch.IPPROTO_UDP:
            tcp_src_port = 0
            tcp_dst_port = 0
            if flow_fields[3]:
                udp_src_port = int(flow_fields[3])
            else:
                udp_src_port = 0
            if flow_fields[5]:
                udp_dst_port = int(flow_fields[5])
            else:
                udp_dst_port = 0
        else:
            udp_src_port = 0
            udp_dst_port = 0
            tcp_src_port = 0
            tcp_dst_port = 0

        flow_match_list.append(
            FlowDescription(
                match=FlowMatch(
                    direction=flow_direction,
                    ip_proto=ip_protocol,
                    ipv4_src=flow_fields[2],
                    ipv4_dst=flow_fields[4],
                    tcp_src=tcp_src_port,
                    tcp_dst=tcp_dst_port,
                    udp_src=udp_src_port,
                    udp_dst=udp_dst_port,
                ),
                action=FlowDescription.PERMIT,
            ),
        )

    qos_parameter_list = args.qos.split(",")
    if len(qos_parameter_list) == 7:
        # utilize user passed arguments
        policy_qos = FlowQos(
            qci=int(args.qci),
            max_req_bw_ul=int(qos_parameter_list[0]),
            max_req_bw_dl=int(qos_parameter_list[1]),
            gbr_ul=int(qos_parameter_list[2]),
            gbr_dl=int(qos_parameter_list[3]),
            arp=QosArp(
                priority_level=int(qos_parameter_list[4]),
                pre_capability=int(qos_parameter_list[5]),
                pre_vulnerability=int(qos_parameter_list[6]),
            ),
        )
    else:
        # parameter missing, use default values
        policy_qos = FlowQos(
            qci=int(args.qci),
            max_req_bw_ul=100000,
            max_req_bw_dl=100000,
            arp=QosArp(priority_level=1, pre_capability=1, pre_vulnerability=0),
        )

    policy_rule = PolicyRule(
        id=args.policy_id,
        priority=int(args.priority),
        flow_list=flow_match_list,
        tracking_type=PolicyRule.NO_TRACKING,
        rating_group=1,
        monitoring_key=None,
        qos=policy_qos,
    )

    qos = QoSInformation(qci=int(args.qci))

    reauth_result = sessiond_client.PolicyReAuth(
        PolicyReAuthRequest(
            session_id=args.session_id,
            imsi=args.imsi,
            rules_to_remove=[],
            rules_to_install=[],
            dynamic_rules_to_install=[DynamicRuleInstall(policy_rule=policy_rule)],
            event_triggers=[],
            revalidation_time=None,
            usage_monitoring_credits=[],
            qos_info=qos,
        ),
    )
    print(reauth_result)


@grpc_wrapper
def send_abort_session(client, args):
    abort_session_chan = ServiceRegistry.get_rpc_channel("abort_session_service", ServiceRegistry.LOCAL)
    abort_session_client = AbortSessionResponderStub(abort_session_chan)

    asr = AbortSessionRequest(
        session_id=args.session_id,
        user_name=args.user_name,
    )

    asa = abort_session_client.AbortSession(asr)
    print(asa)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description="Management CLI for testing session manager",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title="subcommands", dest="cmd")

    # Create Session
    create_session_parser = subparsers.add_parser(
        "create_session",
        help="Send Create Session Request to session_proxy service in FeG",
    )
    create_session_parser.add_argument("imsi", help="e.g., 001010000088888")
    create_session_parser.set_defaults(func=send_create_session)

    # PolicyReAuth
    create_session_parser = subparsers.add_parser(
        "policy_rar", help="Send Policy Reauthorization Request to sessiond",
    )
    create_session_parser.add_argument("imsi", help="e.g., IMSI001010000088888")
    create_session_parser.add_argument(
        "session_id", help="e.g., IMSI001010000088888-910385",
    )
    create_session_parser.add_argument("policy_id", help="e.g., ims-voice")
    create_session_parser.add_argument(
        "priority", help="e.g., precedence value in the range [0-255]",
    )
    create_session_parser.add_argument(
        "qci", help="e.g., 9 for default, 1 for VoIP data, 5 for IMS signaling",
    )
    create_session_parser.add_argument(
        "flow_rules",
        help="List of 6-tuples: "
             "[direction,protocol,src_ip,src_port,dst_ip,dst_port] "
             "separated by ';',e.g., "
             "UL,6,192.168.50.1,0,192.168.40.2,12345;DL,1,8.8.8.8,0,192.168.50.1,0",
    )
    create_session_parser.add_argument(
        "qos",
        help="QoS-tuple: [max_req_bw_ul,max_req_bw_dl,gbr_ul,gbr_dl,arp_prio,"
             "pre_cap,pre_vul] e.g., 10000000,10000000,0,0,15,1,0",
    )
    create_session_parser.set_defaults(func=send_policy_rar)

    # ASR
    abort_session_parser = subparsers.add_parser(
        "abort_session_service",
        help="Send an AbortSessionRequest to SessionD service",
    )
    abort_session_parser.add_argument("session_id", help="e.g., 001010000088888")
    abort_session_parser.add_argument("user_name", help="e.g., IMSI010000088888")
    abort_session_parser.set_defaults(func=send_abort_session)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, LocalSessionManagerStub, "sessiond")


if __name__ == "__main__":
    main()
