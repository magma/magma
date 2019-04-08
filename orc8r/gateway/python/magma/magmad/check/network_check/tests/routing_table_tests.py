"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import textwrap
import unittest

from magma.magmad.check.network_check import routing_table


class RoutingTableParseTests(unittest.TestCase):
    def test_parse_bad_output(self):
        expected = routing_table.RouteCommandResult(error='err',
                                                    routing_table=None)
        actual = routing_table.parse_route_output('output', 'err', None)
        self.assertEqual(expected, actual)

        output = textwrap.dedent('''
        Kernel IP routing table
        bad
        ''').strip().encode('ascii')
        expected = routing_table.RouteCommandResult(
            error='Unexpected heading: bad',
            routing_table=None)
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
                    destination='0.0.0.0',
                    gateway='10.0.2.2',
                    genmask='0.0.0.0',
                    flags='UG',
                    metric='0',
                    ref='0',
                    use='0',
                    interface='eth0',
                ),
                routing_table.Route(
                    destination='10.0.2.0',
                    gateway='0.0.0.0',
                    genmask='255.255.255.0',
                    flags='U',
                    metric='0',
                    ref='0',
                    use='0',
                    interface='eth0',
                ),
                routing_table.Route(
                    destination='169.254.0.0',
                    gateway='0.0.0.0',
                    genmask='255.255.0.0',
                    flags='U',
                    metric='1000',
                    ref='0',
                    use='0',
                    interface='eth0',
                ),
                routing_table.Route(
                    destination='192.168.60.0',
                    gateway='0.0.0.0',
                    genmask='255.255.255.0',
                    flags='U',
                    metric='0',
                    ref='0',
                    use='0',
                    interface='eth1',
                ),
                routing_table.Route(
                    destination='192.168.128.0',
                    gateway='0.0.0.0',
                    genmask='255.255.255.0',
                    flags='U',
                    metric='0',
                    ref='0',
                    use='0',
                    interface='gtp_br0',
                ),
                routing_table.Route(
                    destination='192.168.129.0',
                    gateway='0.0.0.0',
                    genmask='255.255.255.0',
                    flags='U',
                    metric='0',
                    ref='0',
                    use='0',
                    interface='eth2',
                ),
            ])
        actual = routing_table.parse_route_output(output, '', None)
        self.assertEqual(expected, actual)


if __name__ == '__main__':
    unittest.main()
