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
import textwrap
import unittest

from magma.magmad.check.network_check import routing_table


class RoutingTableParseTests(unittest.TestCase):
    def test_parse_bad_output(self):
        expected = routing_table.RouteCommandResult(
            error='err',
            routing_table=[],
        )
        actual = routing_table.parse_route_output('output', 'err', None)
        self.assertEqual(expected, actual)

        output = textwrap.dedent('''
        Kernel IP routing table
        bad
        ''').strip().encode('ascii')
        expected = routing_table.RouteCommandResult(
            error='Unexpected heading: bad',
            routing_table=[],
        )
        actual = routing_table.parse_route_output(output, '', None)
        self.assertEqual(expected, actual)

    def test_parse_good_output(self):
        output = textwrap.dedent('''
        Kernel IP routing table
        Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
        0.0.0.0         10.0.2.2        0.0.0.0         UG    0      0        0 eth0
        10.0.2.0        0.0.0.0         255.255.255.0   U     0      0        0 eth0
        169.254.0.0     0.0.0.0         255.255.0.0     U     1000   0        0 eth0
        192.168.60.0    0.0.0.0         255.255.255.0   U     0      0        0 eth1
        192.168.128.0   0.0.0.0         255.255.255.0   U     0      0        0 gtp_br0
        192.168.129.0   0.0.0.0         255.255.255.0   U     0      0        0 eth2
        ''').strip().encode('ascii')

        expected = routing_table.RouteCommandResult(
            error=None,
            routing_table=[
                routing_table.Route(
                    destination_ip='0.0.0.0',
                    gateway_ip='10.0.2.2',
                    genmask='0.0.0.0',
                    network_interface_id='eth0',
                )._asdict(),
                routing_table.Route(
                    destination_ip='10.0.2.0',
                    gateway_ip='0.0.0.0',
                    genmask='255.255.255.0',
                    network_interface_id='eth0',
                )._asdict(),
                routing_table.Route(
                    destination_ip='169.254.0.0',
                    gateway_ip='0.0.0.0',
                    genmask='255.255.0.0',
                    network_interface_id='eth0',
                )._asdict(),
                routing_table.Route(
                    destination_ip='192.168.60.0',
                    gateway_ip='0.0.0.0',
                    genmask='255.255.255.0',
                    network_interface_id='eth1',
                )._asdict(),
                routing_table.Route(
                    destination_ip='192.168.128.0',
                    gateway_ip='0.0.0.0',
                    genmask='255.255.255.0',
                    network_interface_id='gtp_br0',
                )._asdict(),
                routing_table.Route(
                    destination_ip='192.168.129.0',
                    gateway_ip='0.0.0.0',
                    genmask='255.255.255.0',
                    network_interface_id='eth2',
                )._asdict(),
            ],
        )
        actual = routing_table.parse_route_output(output, '', None)
        self.assertEqual(expected, actual)


if __name__ == '__main__':
    unittest.main()
