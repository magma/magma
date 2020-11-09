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

from magma.pipelined.tunnel_id_store import TunnelToTunnelMapper


class InterfaceMappersTest(unittest.TestCase):
    def setUp(self):
        self._tunnel_mapper = TunnelToTunnelMapper()
        self._tunnel_mapper._tunnel_map = {}

    def test_prefix_mapper_test(self):
        uplink_tunnels = [0xf1231]
        downlink_tunnels = [0x111ef3, 0x21312]

        self._tunnel_mapper.save_tunnels(uplink_tunnels[0], downlink_tunnels[0])
        self.assertEqual(self._tunnel_mapper.get_tunnel(uplink_tunnels[0]),
                         downlink_tunnels[0])
        self.assertEqual(self._tunnel_mapper.get_tunnel(downlink_tunnels[0]),
                         uplink_tunnels[0])

        self._tunnel_mapper.save_tunnels(uplink_tunnels[0], downlink_tunnels[1])
        self.assertEqual(self._tunnel_mapper.get_tunnel(uplink_tunnels[0]),
                         downlink_tunnels[1])
        self.assertEqual(self._tunnel_mapper.get_tunnel(downlink_tunnels[1]),
                         uplink_tunnels[0])


if __name__ == "__main__":
    unittest.main()
