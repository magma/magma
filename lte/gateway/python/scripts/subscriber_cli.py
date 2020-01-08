#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse

from lte.protos.subscriberdb_pb2 import GSMSubscription, LTESubscription, \
    SubscriberData, SubscriberState, SubscriberUpdate, Non3GPPUserProfile
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from orc8r.protos.common_pb2 import Void

from magma.common.rpc_utils import grpc_wrapper
from magma.subscriberdb.sid import SIDUtils


@grpc_wrapper
def add_subscriber(client, args):
    gsm = GSMSubscription()
    lte = LTESubscription()
    state = SubscriberState()
    apn = Non3GPPUserProfile()

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

    if args.apn_name is not None:
        apn.apn_config.service_selection = args.apn_name

    if args.qci is not None:
        print("qci", args.qci)
        apn.apn_config.qos_profile.class_id = args.qci

    data = SubscriberData(sid=SIDUtils.to_pb(args.sid), gsm=gsm,
            lte=lte, state=state, non_3gpp=apn)
    client.AddSubscriber(data)


@grpc_wrapper
def update_subscriber(client, args):
    update = SubscriberUpdate()
    update.data.sid.CopyFrom(SIDUtils.to_pb(args.sid))
    gsm = update.data.gsm
    lte = update.data.lte
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
        lte.auth_opc = bytes.fromhex(args.lte_auth_opc)
        fields.append('lte.state')
        fields.append('lte.auth_opc')


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
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_add = subparsers.add_parser('add', help='Add a new subscriber')
    parser_del = subparsers.add_parser('delete', help='Delete a subscriber')
    parser_update = subparsers.add_parser('update', help='Update a subscriber')
    parser_get = subparsers.add_parser('get', help='Get subscriber data')
    parser_list = subparsers.add_parser('list', help='List all subscriber ids')

    # Add arguments
    for cmd in [parser_add, parser_del, parser_update, parser_get]:
        cmd.add_argument('sid', help='Subscriber identifier')
    for cmd in [parser_add, parser_update]:
        cmd.add_argument('--gsm-auth-tuple', default=[], action='append',
                         help='GSM authentication tuple (hex digits)')
        cmd.add_argument('--lte-auth-key', help='LTE authentication key')
        cmd.add_argument('--lte-auth-opc', help='LTE authentication opc')
        cmd.add_argument('--lte-auth-next-seq', type=int,
                         help='LTE authentication seq number (hex digits)')
        cmd.add_argument('--apn-name',
                         help='Name of the APN (ims/internet)')
        cmd.add_argument('--qci', type=int,
                         help='QCI for the APN')


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
