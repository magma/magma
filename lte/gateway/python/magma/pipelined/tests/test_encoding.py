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

from magma.pipelined.encoding import encrypt_str, get_hash, encode_str
from lte.protos.mconfig.mconfigs_pb2 import PipelineD


class EncodingTest(unittest.TestCase):
    def test_session_rule_version_mapper(self):
        """
        Example encoding:
            MD5(C14r0315v0x)=37ee40eecb484166d68c29930e48313c
            RC4(key=37ee40eecb484166d68c29930e48313c, msisdn=5521966054601))=>a4cc8cb22393a667611f4939b6
            base64(a4cc8cb22393a667611f4939b6)=>pMyMsiOTpmdhH0k5tg==
        """
        msisdn = '5521966054601'
        key = 'C14r0315v0x'

        hash = get_hash(key, PipelineD.HEConfig.MD5)
        self.assertEqual(hash, '37ee40eecb484166d68c29930e48313c')

        """
        encrypted = encrypt_str(msisdn, hash, PipelineD.HEConfig.RC4)
        self.assertEqual(encrypted, 'a4cc8cb22393a667611f4939b6')

        ret = encode_str(encrypted, PipelineD.HEConfig.BASE64)
        self.assertEqual(ret, 'pMyMsiOTpmdhH0k5tg==')
        """

        # Not implemented yet
        hash = get_hash(key, PipelineD.HEConfig.SHA256)
        self.assertEqual(hash, '')
        encrypted = encrypt_str(msisdn, hash,
                                PipelineD.HEConfig.AES256_CBC_HMAC_MD5)
        self.assertEqual(encrypted, '')
        ret = encode_str(encrypted, PipelineD.HEConfig.BASE64)
        self.assertEqual(ret, '')


if __name__ == "__main__":
    unittest.main()
