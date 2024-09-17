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
    CoreNetworkType,
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
    """
    Add a subscriber to the SubscriberDB service.

    Args:
        client (SubscriberDBStub): The gRPC client for the SubscriberDB service.
        args (argparse.Namespace): The command line arguments.

    Returns:
        None

    Raises:
        None

    Description:
        This function adds a subscriber to the SubscriberDB service. It takes a gRPC
        client and command-line arguments as input. The function initializes the GSM,
        LTE, state, and sub_network objects.

        If the --gsm_auth_tuple flag is provided, the function sets the GSM state to
        ACTIVE and appends the provided authentication tuples to the gsm.auth_tuples
        list.

        If the --lte_auth_key flag is provided, the function sets the LTE state to
        ACTIVE and sets the LTE authentication key.

        If the --lte_auth_next_seq flag is provided, the function sets the LTE next
        sequence number.

        If the --lte_auth_opc flag is provided, the function sets the LTE authentication
        OPC.

        If the --forbidden_network_types flag is provided, the function checks if the
        number of forbidden network types is greater than 2 and prints an error message.
        If the provided network types are valid, the function adds them to the
        sub_network.forbidden_network_types list.

        Finally, the function creates a SubscriberData object with the provided SID,
        GSM subscription details, LTE subscription details, subscriber state, and
        forbidden network types, and calls the AddSubscriber RPC on the client.

    Note:
        - The function assumes that the gRPC client and command-line arguments are valid.
        - The function does not perform any input validation.
        - The function does not handle any exceptions.
    """
    gsm = GSMSubscription()
    lte = LTESubscription()
    state = SubscriberState()
    sub_network = CoreNetworkType()

    if args.gsm_auth_tuple:
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

    if args.forbidden_network_types is not None:
        if (len(args.forbidden_network_types.split(",")) > 2):
            print("Forbidden Core Network Types are NT_5GC, NT_EPC")
            return
        for n in args.forbidden_network_types.split(","):
            if n == "NT_5GC":
                sub_network.forbidden_network_types.extend([CoreNetworkType.NT_5GC])
            elif n == "NT_EPC":
                sub_network.forbidden_network_types.extend([CoreNetworkType.NT_EPC])
            else:
                print("Invalid Network type, Forbidden Core Network Types are NT_5GC, NT_EPC")
                return

    data = SubscriberData(
        sid=SIDUtils.to_pb(args.sid), gsm=gsm, lte=lte, state=state, sub_network=sub_network,
    )
    client.AddSubscriber(data)


@grpc_wrapper
def update_subscriber(client, args):
    """
    Update a subscriber's information.

    Args:
        client (SubscriberDBStub): The gRPC client for the SubscriberDB service.
        args (Namespace): The parsed command-line arguments.

    Returns:
        None

    Raises:
        None

    Description:
        This function updates a subscriber's information in the SubscriberDB service.
        It takes a gRPC client and command-line arguments as input. The function creates
        a SubscriberUpdate object and populates its data field with the subscriber's
        information. It then updates the fields parameter with the paths of the fields
        that need to be updated.

        If the gsm_auth_tuple argument is provided, the function sets the gsm.state field
        to ACTIVE and appends the provided authentication tuples to the gsm.auth_tuples
        field. It also appends the 'gsm.state' and 'gsm.auth_tuples' paths to the fields
        parameter.

        If the lte_auth_key argument is provided, the function sets the lte.state field
        to ACTIVE and sets the lte.auth_key field to the provided authentication key. It
        also appends the 'lte.state' and 'lte.auth_key' paths to the fields parameter.

        If the lte_auth_next_seq argument is provided, the function sets the
        state.lte_auth_next_seq field to the provided value. It appends the
        'state.lte_auth_next_seq' path to the fields parameter.

        If the lte_auth_opc argument is provided, the function sets the lte.state field
        to ACTIVE and sets the lte.auth_opc field to the provided authentication OPC. It
        also appends the 'lte.state' and 'lte.auth_opc' paths to the fields parameter.

        If the apn_config argument is provided, the function creates a dictionary of
        APN configuration parameters and populates the non_3gpp.apn_config field of the
        data object with the provided APN configurations. It appends the 'non_3gpp' path
        to the fields parameter.

        Finally, the function calls the UpdateSubscriber method of the gRPC client with
        the SubscriberUpdate object.

    Note:
        - The function assumes that the gRPC client and command-line arguments are valid.
        - The function does not perform any input validation.
        - The function does not handle any exceptions.
    """
    update = SubscriberUpdate()
    data = update.data
    data.sid.CopyFrom(SIDUtils.to_pb(args.sid))
    fields = update.mask.paths

    if args.gsm_auth_tuple:
        data.gsm.state = GSMSubscription.ACTIVE
        for auth_tuple in args.gsm_auth_tuple:
            data.gsm.auth_tuples.append(bytes.fromhex(auth_tuple))
        fields.append('gsm.state')
        fields.append('gsm.auth_tuples')

    if args.lte_auth_key is not None:
        data.lte.state = LTESubscription.ACTIVE
        data.lte.auth_key = bytes.fromhex(args.lte_auth_key)
        fields.append('lte.state')
        fields.append('lte.auth_key')

    if args.lte_auth_next_seq is not None:
        data.state.lte_auth_next_seq = args.lte_auth_next_seq
        fields.append('state.lte_auth_next_seq')

    if args.lte_auth_opc is not None:
        data.lte.state = LTESubscription.ACTIVE
        data.lte.auth_opc = bytes.fromhex(args.lte_auth_opc)
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
        is_default = "is_default"

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
            is_default,
        )
        apn_data = args.apn_config
        for apn_d in apn_data:
            apn_val = apn_d.split(",")
            if len(apn_val) != 13:
                print(
                    "Incorrect APN parameters."
                    "Please check: subscriber_cli.py update -h",
                )
                return
            apn_dict = dict(zip(apn_keys, apn_val))
            apn_config = data.non_3gpp.apn_config.add()
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
            apn_config.is_default = apn_dict[is_default]

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
    """
    Delete a subscriber using the provided client and arguments.

    Args:
        client (SubscriberDBServicer): The gRPC client for the SubscriberDBServicer.
        args: The command line arguments passed to the script.
    """
    client.DeleteSubscriber(SIDUtils.to_pb(args.sid))


@grpc_wrapper
def get_subscriber(client, args):
    """
    Fetch subscriber data based on the provided client and arguments.
    """
    data = client.GetSubscriberData(SIDUtils.to_pb(args.sid))
    print(data)


@grpc_wrapper
def list_subscribers(client, args):
    """
    List all subscribers by making a gRPC request to the client.

    Args:
        client (SubscriberLookupServicer): The gRPC client for the SubscriberLookupServicer.
        args: The command line arguments passed to the script.
    """
    for sid in client.ListSubscribers(Void()).sids:
        print(SIDUtils.to_str(sid))


def create_parser():
    """
    Create the argparse parser with all the arguments.
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
    for cmd in list(
        parser_add,
        parser_del,
        parser_update,
        parser_get,
    ):
        cmd.add_argument("sid", help="Subscriber identifier")

    # Add subcommand arguments
    parser_add.add_argument(
        "--gsm-auth-tuple",
        default=[],
        action="append",
        help="GSM authentication tuple (hex digits)",
    )
    parser_add.add_argument("--lte-auth-key", help="LTE authentication key")
    parser_add.add_argument("--lte-auth-opc", help="LTE authentication opc")
    parser_add.add_argument(
        "--lte-auth-next-seq",
        type=int,
        help="LTE authentication seq number (hex digits)",
    )
    parser_add.add_argument("--forbidden-network-types", help="Core NetworkType Restriction")

    # Update subcommand arguments
    parser_update.add_argument(
        "--gsm-auth-tuple",
        default=[],
        action="append",
        help="GSM authentication tuple (hex digits)",
    )
    parser_update.add_argument("--lte-auth-key", help="LTE authentication key")
    parser_update.add_argument("--lte-auth-opc", help="LTE authentication opc")
    parser_update.add_argument(
        "--lte-auth-next-seq",
        type=int,
        help="LTE authentication seq number (hex digits)",
    )
    parser_update.add_argument(
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
    """
    Create a parser, parses arguments, check for a command, and execute a subcommand function.
    """
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
