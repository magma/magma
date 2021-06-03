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

from magma.mobilityd.mac import create_mac_from_sid


class MacAddressGenTests(unittest.TestCase):
    def test_mac_address_generation(self):
        mac = create_mac_from_sid("IMSI00000000001")
        self.assertEqual(str(mac), "8A:00:00:00:00:01")

        mac = create_mac_from_sid("IMSI980000002918")
        self.assertEqual(str(mac), "8A:E4:2C:8D:53:66")

        mac = create_mac_from_sid("IMSI311980000002918")
        self.assertEqual(str(mac), "8A:E4:2C:8D:53:66")

        mac = create_mac_from_sid("IMSI02918")
        self.assertEqual(str(mac), "8A:00:00:00:0B:66")

    def test_mac_from_hex(self):
        mac = create_mac_from_sid("123456789012")
        self.assertEqual(str(mac), "12:34:56:78:90:12")

    def test_mac_to_hex(self):
        mac = create_mac_from_sid("12:34:56:78:90:12")
        self.assertEqual(mac.as_hex(), b'\x124Vx\x90\x12')


if __name__ == '__main__':
    unittest.main()
