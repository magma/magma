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
from typing import NamedTuple

import grpc
from lte.protos.subscriberdb_pb2 import M5GSUCIRegistrationRequest
from lte.protos.subscriberdb_pb2_grpc import (
    M5GSUCIRegistrationStub,
    SuciProfileDBStub,
)
from magma.common.rpc_utils import grpc_wrapper
from magma.subscriberdb.crypto.ECIES import ECIES_UE
from orc8r.protos.common_pb2 import Void

ue_encrypt_context = NamedTuple(
    'ue_encrypt_context', [
        ('pub_key', bytes),
        ('cipher_text', bytes),
        ('mac', bytes),
    ],
)


class UEmsin(object):
    """
    Class to protect IMSI
    """

    def __init__(self, msin: bytes, profile: str = "ProfileA"):
        """
        __init__ function to initialize
        """
        self.profile = profile
        self.msin = msin
        self.ue_encrypt_context = ue_encrypt_context(b'', b'', b'')

    def encrypt(self, home_network_public_key: bytes):
        """
        Encrypt
        """
        if self.profile == "ProfileA":
            ue = ECIES_UE(profile='A')
        elif self.profile == "ProfileB":
            ue = ECIES_UE(profile='B')

        ue.generate_sharedkey(home_network_public_key)
        ue_pubkey, ue_ciphertext, ue_mac = ue.protect(self.msin)
        self.ue_encrypt_context = ue_encrypt_context(
            ue_pubkey,
            ue_ciphertext, ue_mac,
        )

    def validate_msin(self, msin_in_home_network: bytes):
        """
        validate_msin
        """
        print("")
        print("\u0332".join("OUTPUT OF MSIN COMPARE"))
        print("UE MSIN        = ", self.msin)
        print("AMF recvd MSIN = ", msin_in_home_network)
        assert self.msin == msin_in_home_network        # noqa: S101


def get_suciprofile(home_net_public_key_id: int):
    """
    get_suciprofile
    """
    channel = grpc.insecure_channel('localhost:50051')
    stub = SuciProfileDBStub(channel)
    response = stub.ListSuciProfile(Void())
    if not response.suci_profiles:
        print("SuciProfileList is empty")
    else:
        for x in response.suci_profiles:
            if home_net_public_key_id == x.home_net_public_key_id:
                return x.home_net_public_key
        return None


@grpc_wrapper
def decrypt_msin(client, args):
    """
    decrypt_msin
    """
    if len(args.imsi) != 10:
        print("Length of the IMSI provided is incorrect")
        return

    if int(args.protection_scheme) == 0:
        profile = "ProfileA"
    elif int(args.protection_scheme) == 1:
        profile = "ProfileB"
    else:
        print("Invalid protection_scheme value")
        return

    if int(args.ue_pubkey_identifier) < 0 or int(args.ue_pubkey_identifier) > 255:
        print("Invalid home_net_public_key_id value:", args.ue_pubkey_identifier)
        return

    ue_msin = UEmsin(bytes.fromhex(args.imsi), profile)

    pubkey = get_suciprofile(int(args.ue_pubkey_identifier))

    if pubkey is not None:
        ue_msin.encrypt(pubkey)
        print("pubkey received:", pubkey.hex())
    else:
        print("Not a valid home_net_public_key")
        return

    print("")
    print(" ==== ENCRYPTION in UE ====")
    print("\u0332".join("UE : Using home network public key generate"))
    print("ue pub_key: ", ue_msin.ue_encrypt_context.pub_key.hex())
    print("cipher_tex: ", ue_msin.ue_encrypt_context.cipher_text.hex())
    print("mac       : ", ue_msin.ue_encrypt_context.mac.hex())

    request = M5GSUCIRegistrationRequest(
        ue_pubkey_identifier=int(args.ue_pubkey_identifier),
        ue_pubkey=ue_msin.ue_encrypt_context.pub_key,
        ue_ciphertext=ue_msin.ue_encrypt_context.cipher_text,
        ue_encrypted_mac=ue_msin.ue_encrypt_context.mac,
    )
    response = client.M5GDecryptImsiSUCIRegistration(request)
    print("Deconcealed IMSI received from subscriberdb:", response.ue_msin_recv.hex())


def main():
    """Creates the argparse parser with all the arguments."""          # noqa: D401

    parser = argparse.ArgumentParser(
        description='Management CLI for SuciProfile',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # decrypt_msin
    subparser = subparsers.add_parser('decrypt', help='Decrypt msin')
    subparser.add_argument("--imsi", help="Last 10 digit imsi" "  e.g: --imsi 1032547698")
    subparser.add_argument(
        "--protection_scheme", help="ProfileA(0), ProfileB(1)"
        "  e.g: --protection_scheme 0",
    )
    subparser.add_argument(
        "--ue_pubkey_identifier", help="Cached ue_pubkey_identifier"
        "  e.g: --ue_pubkey_identifier 0..255",
    )
    subparser.set_defaults(func=decrypt_msin)

    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, M5GSUCIRegistrationStub, 'subscriberdb')


if __name__ == "__main__":
    main()
