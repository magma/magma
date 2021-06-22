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

from lte.protos.subscriberdb_pb2 import (
    GSMSubscription,
    LTESubscription,
    SubscriberData,
    SubscriberState,
    SubscriberUpdate,
)
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.common.rpc_utils import grpc_wrapper
from magma.subscriberdb.sid import SIDUtils
from orc8r.protos.common_pb2 import Void


@grpc_wrapper
def add_subscriber(client, args):
    gsm = GSMSubscription()
    lte = LTESubscription()
    state = SubscriberState()

    if len(args.gsm_auth_tuple) != 0:
        gsm.state = GSMSubscription.ACTIVE
        for auth_tuple in args.gsm_auth_tuple:
            gsm.auth_tuples.append(bytes.fromhex(auth_tuple))

    if args.lte_auth_key is not None:
        lte.state = LTESubscription.ACTIVE
        lte.auth_key = bytes.fromhex(args.lte_auth_key)

    if args.lte_auth_next_seq is not None:
        state.lte_auth_next_seq = args.lte_auth_next_seq

    if args.lte_auth_opc is not None:
        lte.auth_opc = bytes.fromhex(args.lte_auth_opc)

    data = SubscriberData(
        sid=SIDUtils.to_pb(args.sid), gsm=gsm, lte=lte, state=state,
    )
    client.AddSubscriber(data)


@grpc_wrapper
def update_subscriber(client, args):
    update = SubscriberUpdate()
    update.data.sid.CopyFrom(SIDUtils.to_pb(args.sid))
    gsm = update.data.gsm
    lte = update.data.lte
    non_3gpp = update.data.non_3gpp
    fields = update.mask.paths

    if len(args.gsm_auth_tuple) != 0:
        gsm.state = GSMSubscription.ACTIVE
        for auth_tuple in args.gsm_auth_tuple:
            gsm.auth_tuples.append(bytes.fromhex(auth_tuple))
        fields.append('gsm.state')
        fields.append('gsm.auth_tuples')

    if args.lte_auth_key is not None:
        lte.state = LTESubscription.ACTIVE
        lte.auth_key = bytes.fromhex(args.lte_auth_key)
        fields.append('lte.state')
        fields.append('lte.auth_key')

    if args.lte_auth_next_seq is not None:
        update.data.state.lte_auth_next_seq = args.lte_auth_next_seq
        fields.append('state.lte_auth_next_seq')

    if args.lte_auth_opc is not None:
        lte.state = LTESubscription.ACTIVE
        lte.auth_opc = bytes.fromhex(args.lte_auth_opc)
        fields.append('lte.state')
        fields.append('lte.auth_opc')

    if args.apn_config is not None:
        apn_name = "apn_name"
        qci = "qci"
        priority = "priority"
        pre_cap = "preemption_capability"
        pre_vul = "preemption_vulnerability"
        ul = "mbr_uplink"
        dl = "mbr_downlink"
        pdn_type = "pdn_type"
        static_ip = "static_ip"
        vlan_id = "vlan"
        gw_ip = "gw_ip"
        gw_mac = "gw_mac"

        apn_keys = (
            apn_name,
            qci,
            priority,
            pre_cap,
            pre_vul,
            ul,
            dl,
            pdn_type,
            static_ip,
            vlan_id,
            gw_ip,
            gw_mac,
        )
        apn_data = args.apn_config
        for apn_d in apn_data:
            apn_val = apn_d.split(",")
            if len(apn_val) != 12:
                print(
                    "Incorrect APN parameters."
                    "Please check: subscriber_cli.py update -h",
                )
                return
            apn_dict = dict(zip(apn_keys, apn_val))
            apn_config = non_3gpp.apn_config.add()
            apn_config.service_selection = apn_dict[apn_name]
            apn_config.qos_profile.class_id = int(apn_dict[qci])
            apn_config.qos_profile.priority_level = int(apn_dict[priority])
            apn_config.qos_profile.preemption_capability = int(
                apn_dict[pre_cap],
            )
            apn_config.qos_profile.preemption_vulnerability = int(
                apn_dict[pre_vul],
            )
            apn_config.ambr.max_bandwidth_ul = int(apn_dict[ul])
            apn_config.ambr.max_bandwidth_dl = int(apn_dict[dl])
            apn_config.pdn = int(apn_dict[pdn_type])
            apn_config.assigned_static_ip = apn_dict[static_ip]

            if apn_dict[vlan_id]:
                apn_config.resource.vlan_id = int(apn_dict[vlan_id])
            if apn_dict[gw_ip]:
                apn_config.resource.gateway_ip = apn_dict[gw_ip]
                # allow mac address if gw-ip is specified
                if apn_dict[gw_mac]:
                    apn_config.resource.gateway_mac = apn_dict[gw_mac]

        fields.append("non_3gpp")

    client.UpdateSubscriber(update)


@grpc_wrapper
def delete_subscriber(client, args):
    client.DeleteSubscriber(SIDUtils.to_pb(args.sid))


@grpc_wrapper
def get_subscriber(client, args):
    data = client.GetSubscriberData(SIDUtils.to_pb(args.sid))
    print(data)


@grpc_wrapper
def list_subscribers(client, args):
    for sid in client.ListSubscribers(Void()).sids:
        print(SIDUtils.to_str(sid))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for SubscriberDB',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title="subcommands", dest="cmd")
    parser_add = subparsers.add_parser("add", help="Add a new subscriber")
    parser_del = subparsers.add_parser("delete", help="Delete a subscriber")
    parser_update = subparsers.add_parser("update", help="Update a subscriber")
    parser_get = subparsers.add_parser("get", help="Get subscriber data")
    parser_list = subparsers.add_parser("list", help="List all subscriber ids")

    # Add arguments
    for cmd in [
        parser_add,
        parser_del,
        parser_update,
        parser_get,
    ]:
        cmd.add_argument("sid", help="Subscriber identifier")
    for cmd in [parser_add]:
        cmd.add_argument(
            "--gsm-auth-tuple",
            default=[],
            action="append",
            help="GSM authentication tuple (hex digits)",
        )
        cmd.add_argument("--lte-auth-key", help="LTE authentication key")
        cmd.add_argument("--lte-auth-opc", help="LTE authentication opc")
        cmd.add_argument(
            "--lte-auth-next-seq",
            type=int,
            help="LTE authentication seq number (hex digits)",
        )

    for cmd in [parser_update]:
        cmd.add_argument(
            "--gsm-auth-tuple",
            default=[],
            action="append",
            help="GSM authentication tuple (hex digits)",
        )
        cmd.add_argument("--lte-auth-key", help="LTE authentication key")
        cmd.add_argument("--lte-auth-opc", help="LTE authentication opc")
        cmd.add_argument(
            "--lte-auth-next-seq",
            type=int,
            help="LTE authentication seq number (hex digits)",
        )
        cmd.add_argument(
            "--apn-config",
            action="append",
            help="APN parameters to add/update in the order :"
            " [apn-name, qci, priority, preemption-capability,"
            " preemption-vulnerability, mbr-ul, mbr-dl, pdn-type,"
            " [0-IPv4, 1-IPv6, 2-IPv4v6]"
            " static-ip, vlan_id, internet_gw_ip, internet_gw_mac]"
            " [e.g --apn-config ims,5,15,1,1,1000,2000,1,,,,"
            " --apn-config internet,9,1,0,0,3000,4000,0,1.2.3.4,,,"
            " --apn-config internet,9,1,0,0,3000,4000,2,"
            "1.2.3.4,1,2.2.2.2,11:22:33:44:55:66]",
        )

# Add function callbacks
    parser_add.set_defaults(func=add_subscriber)
    parser_del.set_defaults(func=delete_subscriber)
    parser_update.set_defaults(func=update_subscriber)
    parser_get.set_defaults(func=get_subscriber)
    parser_list.set_defaults(func=list_subscribers)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, SubscriberDBStub, 'subscriberdb')


if __name__ == "__main__":
    main()
