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

from magma.pipelined.imsi import decode_imsi, encode_imsi


class IMSIEncodingTest(unittest.TestCase):
    """
    Test the encoding and decoding of IMSI values used in flow metadata in
    pipelined
    """

    def test_leading_zeros(self):
        imsi_list = [
            "IMSI001010000000013",
            "IMSI011010000000013",
            "IMSI111010000000013",
        ]
        for imsi in imsi_list:
            self._assert_conversion(imsi)

    def test_different_lengths(self):
        imsi_list = [
            "IMSI001010000000013",
            "IMSI01010000000013",
            "IMSI28950000000013",
        ]
        for imsi in imsi_list:
            self._assert_conversion(imsi)

    def _assert_conversion(self, imsi):
        """
        To test the validity of the conversion, we can encode the imsi and
        decode to see if the two values match
        """
        test_val = decode_imsi(encode_imsi(imsi))
        self.assertEqual(imsi, test_val)


if __name__ == "__main__":
    unittest.main()
