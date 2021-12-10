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

import unittest
from binascii import unhexlify

from cryptography.hazmat.primitives.asymmetric.x25519 import X25519PrivateKey
from magma.subscriberdb.crypto.EC import ECDH_SECP256R1, X25519
from magma.subscriberdb.crypto.ECIES import ECIES_HN, ECIES_UE


class ECIES_test(unittest.TestCase):        # noqa: N801
    # annex C.4.3, ECIES Profile A test data
    def test_profileA(self):
        hn_privkey = unhexlify('c53c22208b61860b06c62e5406a7b330c2b577aa5558981510d128247d38bd1d')
        hn_pubkey = unhexlify('5a8d38864820197c3394b92613b20b91633cbd897119273bf8e4a6f4eec0a650')
        eph_privkey = unhexlify('c80949f13ebe61af4ebdbd293ea4f942696b9e815d7e8f0096bbf6ed7de62256')
        eph_pubkey = unhexlify('b2e92f836055a255837debf850b528997ce0201cb82adfe4be1f587d07d8457d')
        shared_key = unhexlify('028ddf890ec83cdf163947ce45f6ec1a0e3070ea5fe57e2b1f05139f3e82422a')
        #
        plaintext = bytes.fromhex('00012080f6')
        ciphertext = unhexlify('cb02352410')
        mactag = unhexlify('cddd9e730ef3fa87')

        x1 = X25519(eph_privkey)
        x2 = X25519(hn_privkey)

        ue = ECIES_UE(profile='A')
        ue.EC.priv_key = X25519PrivateKey.from_private_bytes(eph_privkey)
        hn = ECIES_HN(profile='A', hn_priv_key=hn_privkey)
        ue.generate_sharedkey(hn_pubkey, fresh=False)
        ue_pk, ue_ct, ue_mac = ue.protect(plaintext)
        hn_ct = hn.unprotect(ue_pk, ue_ct, ue_mac)

        self.assertEqual(x1.get_pubkey(), eph_pubkey)
        self.assertEqual(x1.generate_sharedkey(hn_pubkey), shared_key)
        self.assertEqual(x2.get_pubkey(), hn_pubkey)
        self.assertEqual(x2.generate_sharedkey(eph_pubkey), shared_key)
        self.assertEqual(ue_ct, ciphertext)
        self.assertEqual(ue_mac, mactag)
        self.assertEqual(hn_ct, plaintext)

# annex C.4.4, ECIES Profile A test data
    def test_profileB(self):
        hn_pubkey = unhexlify('0272DA71976234CE833A6907425867B82E074D44EF907DFB4B3E21C1C2256EBCD1')  # compressed
        hn_privkey = unhexlify('F1AB1074477EBCC7F554EA1C5FC368B1616730155E0041AC447D6301975FECDA')
        eph_pubkey = unhexlify('039AAB8376597021E855679A9778EA0B67396E68C66DF32C0F41E9ACCA2DA9B9D1')  # compressed
        eph_privkey = unhexlify('99798858A1DC6A2C68637149A4B1DBFD1FDFF5ADDD62A2142F06699ED7602529')
        shared_key = unhexlify('6C7E6518980025B982FBB2FF746E3C2E85A196D252099A7AD23EA7B4C0959CAE')
        #
        plaintext = bytes.fromhex('00012080F6')
        ciphertext = unhexlify('46A33FC271')
        mactag = unhexlify('6AC7DAE96AA30A4D')

        x1 = ECDH_SECP256R1(raw_keypair=(eph_privkey, eph_pubkey))
        x2 = ECDH_SECP256R1(raw_keypair=(hn_privkey, hn_pubkey))
        ue = ECIES_UE(profile='B')
        ue.EC._load_raw_keypair(eph_privkey, eph_pubkey)  # pylint: disable=W0212
        hn = ECIES_HN(None, profile='B', raw_keypair=(hn_privkey, hn_pubkey))
        ue.generate_sharedkey(hn_pubkey, fresh=False)
        ue_pk, ue_ct, ue_mac = ue.protect(plaintext)
        hn_ct = hn.unprotect(ue_pk, ue_ct, ue_mac)

        self.assertEqual(x1.generate_sharedkey(hn_pubkey), shared_key)
        self.assertEqual(x2.generate_sharedkey(eph_pubkey), shared_key)
        self.assertEqual(ue_ct, ciphertext)
        self.assertEqual(ue_mac, mactag)
        self.assertEqual(hn_ct, plaintext)


if __name__ == "__main__":
    unittest.main()
