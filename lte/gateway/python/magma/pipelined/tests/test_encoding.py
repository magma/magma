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

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from magma.pipelined.encoding import (
    decrypt_str,
    encode_str,
    encrypt_str,
    get_hash,
)


class EncodingTest(unittest.TestCase):
    def test_rc4(self):
        """
        Example encoding:
            MD5(C14r0315v0x)=37ee40eecb484166d68c29930e48313c
            RC4(key=37ee40eecb484166d68c29930e48313c, msisdn=5521966054601))=>a82530f1f34cfcdba5569fb60f
            base64(a4cc8cb22393a667611f4939b6)=>qCUw8fNM/NulVp+2Dw==

        Check https://cryptii.com/pipes/rc4-encryption for encryption correctness
        """
        msisdn = '5521966054601'
        key = 'C14r0315v0x'
        mac = 'qqqqq12345'

        hash = get_hash(key, PipelineD.HEConfig.MD5)
        self.assertEqual(hash, b'7\xee@\xee\xcbHAf\xd6\x8c)\x93\x0eH1<')

        encrypted = encrypt_str(msisdn, hash, PipelineD.HEConfig.RC4)
        self.assertEqual(encrypted, 'a82530f1f34cfcdba5569fb60f')

        ret = encode_str(encrypted, PipelineD.HEConfig.BASE64)
        self.assertEqual(ret, 'qCUw8fNM/NulVp+2Dw==')

    def test_aes_encryption(self):
        msisdn = '5521966054601'
        key = 'C14r0315v0x'
        mac = 'qqqqq12345'

        hash = get_hash(key, PipelineD.HEConfig.SHA256)
        self.assertEqual(hash, b'\xe6\xc6?<\xa3\xa4*F\x19v%\xf5\xc6l\x89L\xb1w,TM}\xb7H \x19\xffL\x85\xb4\xac\x82')
        mac_hash = get_hash(mac, PipelineD.HEConfig.SHA256)
        self.assertEqual(mac_hash, b'v\xdf\x84\xf3#I\xb4\xd2\x10\x0c4\xaeQL`\xccDW\xec\x02\n\x9f\xa4\x1ft\x1e\x95:2\xb3-\xb1')

        encrypted = encrypt_str(
            msisdn, hash, PipelineD.HEConfig.AES256_CBC_HMAC_MD5, mac_hash)
        decrypted_str = decrypt_str(
            encrypted, hash, PipelineD.HEConfig.AES256_CBC_HMAC_MD5, mac_hash)
        self.assertEqual(decrypted_str, msisdn)

        encrypted = encrypt_str(
            msisdn, hash, PipelineD.HEConfig.AES256_ECB_HMAC_MD5, mac_hash)
        decrypted_str = decrypt_str(
            encrypted, hash, PipelineD.HEConfig.AES256_ECB_HMAC_MD5, mac_hash)
        self.assertEqual(decrypted_str, msisdn)

        ret = encode_str(encrypted, PipelineD.HEConfig.BASE64)
        self.assertEqual(ret, 'AcXxM8TPFg4bzzv5ya6GIimAf64whfeU+CcfBTnBNo8=')

        encrypted = encrypt_str(
            msisdn, hash, PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1, mac_hash)
        decrypted_str = decrypt_str(
            encrypted, hash, PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1, mac_hash)
        self.assertEqual(decrypted_str, msisdn)

    def test_encoding(self):
        text = '123abcd123'
        ret = encode_str(text, PipelineD.HEConfig.BASE64)
        self.assertEqual(ret, 'Ejq80SM=')

        ret = encode_str(text, PipelineD.HEConfig.HEX2BIN)
        self.assertEqual(ret, '0001001000111010101111001101000100100011')

    def test_hash(self):
        key = 'C14r0315v0x'
        hash = get_hash(key, PipelineD.HEConfig.MD5)
        self.assertEqual(hash, b'7\xee@\xee\xcbHAf\xd6\x8c)\x93\x0eH1<')

        hash = get_hash(key, PipelineD.HEConfig.HEX)
        self.assertEqual(hash, b'4331347230333135763078')

        hash = get_hash(key, PipelineD.HEConfig.SHA256)
        self.assertEqual(hash, b'\xe6\xc6?<\xa3\xa4*F\x19v%\xf5\xc6l\x89L\xb1w,TM}\xb7H \x19\xffL\x85\xb4\xac\x82')

if __name__ == "__main__":
    unittest.main()
