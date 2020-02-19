#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse

from lte.protos.subscriberdb_pb2 import (
    SubscriberData,
    Non3GPPUserProfile,
)
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub

from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos.common_pb2 import Void


@grpc_wrapper
def add_apn(client, args):
    non_3gpp = Non3GPPUserProfile()

    print("Adding APN : ", args.apn)
    apn_config = non_3gpp.apn_config.add()
    apn_config.service_selection = args.apn
    apn_config.qos_profile.class_id = args.qci
    apn_config.qos_profile.priority_level = args.priority
    apn_config.qos_profile.preemption_capability = args.preemptionCapability
    apn_config.qos_profile.preemption_vulnerability = (
        args.preemptionVulnerability
    )
    apn_config.ambr.max_bandwidth_ul = args.mbrUL
    apn_config.ambr.max_bandwidth_dl = args.mbrDL

    data = SubscriberData(non_3gpp=non_3gpp)
    client.AddApn(data)


@grpc_wrapper
def update_apn(client, args):
    non_3gpp = Non3GPPUserProfile()

    print("Updating APN : ", args.apn)
    apn = non_3gpp.apn_config.add()
    apn.service_selection = args.apn
    apn.qos_profile.preemption_capability = args.preemptionCapability
    apn.qos_profile.preemption_vulnerability = args.preemptionVulnerability
    if args.qci is not None:
        apn.qos_profile.class_id = args.qci
    if args.priority is not None:
        apn.qos_profile.priority_level = args.priority
    if args.mbrUL is not None:
        apn.ambr.max_bandwidth_ul = args.mbrUL
    if args.mbrDL is not None:
        apn.ambr.max_bandwidth_dl = args.mbrDL
    update = SubscriberData(non_3gpp=non_3gpp)
    client.UpdateApn(update)


@grpc_wrapper
def delete_apn(client, args):
    print("Deleting APN : ", args.apn)
    non_3gpp = Non3GPPUserProfile()
    apn_config = non_3gpp.apn_config.add()
    apn_config.service_selection = args.apn
    data = SubscriberData(non_3gpp=non_3gpp)
    client.DeleteApn(data)


@grpc_wrapper
def get_apn(client, args):
    print("Retrieving APN : ", args.apn)
    non_3gpp = Non3GPPUserProfile()
    apn_config = non_3gpp.apn_config.add()
    apn_config.service_selection = args.apn
    data = SubscriberData(non_3gpp=non_3gpp)
    apn_data = client.GetApnData(data)
    print(apn_data)


@grpc_wrapper
def list_apns(client, args):
    print("Retrieving APN list")
    for apn in client.ListApns(Void()).apn_name:
        print(apn)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description="Management CLI for APN configuration",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title="subcommands", dest="cmd")
    parser_add = subparsers.add_parser("add", help="Add a new apn")
    parser_del = subparsers.add_parser("delete", help="Delete an apn")
    parser_update = subparsers.add_parser("update", help="Update an apn")
    parser_get = subparsers.add_parser("get", help="Get apn data")
    parser_list = subparsers.add_parser("list", help="List all APNs")

    # Add arguments
    for cmd in [parser_add, parser_del, parser_update, parser_get]:
        cmd.add_argument("apn", help="Name of the APN (ims/internet)")
    for cmd in [parser_add]:
        cmd.add_argument("qci", type=int, help="QCI for APN [1-9]")
        cmd.add_argument(
            "priority", type=int, help="Priority of the APN [1-15]")
        cmd.add_argument(
            "preemptionCapability", type=int, help="Enabled/Disabled [0/1]"
        )
        cmd.add_argument(
            "preemptionVulnerability", type=int, help="Enabled/Disabled [0/1]"
        )
        cmd.add_argument("mbrUL", type=int, help="Max bit rate UL")
        cmd.add_argument("mbrDL", type=int, help="Max bit rate DL")
    for cmd in [parser_update]:
        # preemption_capability and preemption_vulnerability are bool type
        # and cannot be checked for non-zero. Hence they are
        # mandatory parameters

        cmd.add_argument(
            "preemptionCapability", type=int, help="Enabled/Disabled [0/1]"
        )
        cmd.add_argument(
            "preemptionVulnerability", type=int, help="Enabled/Disabled [0/1]"
        )
        cmd.add_argument("--qci", type=int, help="QCI for APN [1-9]")
        cmd.add_argument(
            "--priority", type=int, help="Priority of the APN vaules [1-15]"
        )
        cmd.add_argument("--mbrUL", type=int, help="Max bit rate UL")
        cmd.add_argument("--mbrDL", type=int, help="Max bit rate DL")

    # Add function callbacks
    parser_add.set_defaults(func=add_apn)
    parser_del.set_defaults(func=delete_apn)
    parser_update.set_defaults(func=update_apn)
    parser_get.set_defaults(func=get_apn)
    parser_list.set_defaults(func=list_apns)
    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, SubscriberDBStub, "subscriberdb")


if __name__ == "__main__":
    main()
