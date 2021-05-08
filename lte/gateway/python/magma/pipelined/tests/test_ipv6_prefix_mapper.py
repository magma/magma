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

from magma.pipelined.ipv6_prefix_store import (
    InterfaceIDToPrefixMapper,
    get_ipv6_interface_id,
    get_ipv6_prefix,
)


class InterfaceMappersTest(unittest.TestCase):
    def setUp(self):
        self._interface_to_prefix_mapper = InterfaceIDToPrefixMapper()
        self._interface_to_prefix_mapper._prefix_by_interface = {}

    def test_prefix_mapper_test(self):
        ipv6_addrs = ['ba10:5:6c:9:9d21:4407:d337:1928',
                      '321b:534:6c:9:999:0:d337:1928',
                      '222b:5334:111c:111::d337:1928']
        prefixes = [get_ipv6_prefix(ipv6_addrs[0]),
                    get_ipv6_prefix(ipv6_addrs[1])]
        interfaces = [get_ipv6_interface_id(ipv6_addrs[0]),
                      get_ipv6_interface_id(ipv6_addrs[1]),
                      get_ipv6_interface_id(ipv6_addrs[2])]
        self._interface_to_prefix_mapper.save_prefix(
            interfaces[0], prefixes[0])
        self.assertEqual(
            self._interface_to_prefix_mapper.get_prefix(
                interfaces[0]),
            'ba10:5:6c:9::')

        self._interface_to_prefix_mapper.save_prefix(
            interfaces[1], prefixes[1])
        self.assertEqual(interfaces[1], '::999:0:d337:1928')
        self.assertEqual(
            self._interface_to_prefix_mapper.get_prefix(
                interfaces[1]),
            prefixes[1])

        self._interface_to_prefix_mapper.save_prefix(
            interfaces[0], prefixes[1])
        self.assertEqual(
            self._interface_to_prefix_mapper.get_prefix(
                interfaces[0]),
            '321b:534:6c:9::')

        self.assertEqual(
            self._interface_to_prefix_mapper.get_prefix(
                interfaces[2]),
            None)


if __name__ == "__main__":
    unittest.main()
