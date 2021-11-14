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
    SuciProfile,
)
from lte.protos.subscriberdb_pb2_grpc import SuciProfileDBStub
from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos.common_pb2 import Void
from magma.subscriberdb.crypto.EC import *
from magma.subscriberdb.crypto.ECIES import *
from typing import NamedTuple

home_network_key_pair = NamedTuple(
    'home_network_key_pair', [
        ('home_network_public_key', bytes),
        ('home_network_private_key', bytes),
    ],
)

class HomeNetworkKeyPairGen:

    def __init__(self, profile: str):
        self.profile=profile
        self.home_network_key_pair = home_network_key_pair(b'', b'')

    def core_home_network_key_gen(self):

        if self.profile == "ProfileA" :
            ec = X25519()
        elif self.profile == "ProfileB" :
            ec = ECDH_SECP256R1()

        if ec:
           ec.generate_keypair()
        else:
           return None

        self.home_network_key_pair = home_network_key_pair(ec.get_pubkey(),
                                                           ec.get_privkey())
    def get_home_network_public_key(self):
        return self.home_network_key_pair.home_network_public_key

    def get_home_network_private_key(self):
        return self.home_network_key_pair.home_network_private_key

    def print_key_pair(self):
        print(self.profile)
        print(self.home_network_key_pair.home_network_public_key)
        print(self.home_network_key_pair.home_network_private_key)


@grpc_wrapper
def add_suciprofile(client, args):

    if int(args.protection_scheme) == 0 :
        hnp_gen=HomeNetworkKeyPairGen("ProfileA")
    elif int(args.protection_scheme) == 1:
        hnp_gen=HomeNetworkKeyPairGen("ProfileB")
    else:
        print("Invalid protection_scheme value:", args.protection_scheme)
        return

    if int(args.home_net_public_key_id) < 0 or int(args.home_net_public_key_id) > 255:
        print("Invalid home_net_public_key_id value:", args.home_net_public_key_id)
        return

    hnp_gen.core_home_network_key_gen()

    data = SuciProfile()

    data.home_net_public_key_id = int(args.home_net_public_key_id)
    data.protection_scheme = int(args.protection_scheme)
    data.home_net_public_key = hnp_gen.get_home_network_public_key()
    data.home_net_private_key = hnp_gen.get_home_network_private_key()
    client.AddSuciProfile(data)
          
    
@grpc_wrapper
def delete_suciprofile(client, args):
    if int(args.home_net_public_key_id) < 0 or int(args.home_net_public_key_id) > 255:
        print("Invalid home_net_public_key_id value:", args.home_net_public_key_id)
        return
    client.DeleteSuciProfile(args.home_net_public_key_id)

@grpc_wrapper
def list_suciprofile(client, args):
    #for sid in client.ListSubscribers(Void()).sids:
    for x in range(len(client.ListSuciProfile(Void()).SuciProfileList)):
        print(x)

               
def main():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for SuciProfile',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # add_suciprofile
    subparser = subparsers.add_parser('add', help='Add SuciProfile record')
    subparser.add_argument("--home_net_public_key_id", help="home_network_public_key_id" 
            "  e.g: --home_net_public_key_id 0..255")
    subparser.add_argument("--protection_scheme",      help="ECIESProtectionScheme"
            "  e.g: --protection_scheme 0 or 1")

    subparser.set_defaults(func=add_suciprofile)

    #delete_suciprofile
    subparser = subparsers.add_parser('delete', help='Delete SuciProfile record')
    subparser.add_argument("--home_net_public_key_id", help="home_network_public_key_id"
                            "  e.g: --home_net_public_key_id 0..255")
    subparser.set_defaults(func=delete_suciprofile)


    #list_suciprofile
    subparser = subparsers.add_parser('list', help='List SuciProfile records')
    subparser.set_defaults(func=list_suciprofile)

     # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, SuciProfileDBStub, 'subscriberdb')


if __name__ == "__main__":
    main()

