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

from magma.magmad.check.network_check import ping

# Allow access to protected variables for unit testing
# pylint: disable=protected-access


class PingArgFactoryTests(unittest.TestCase):
    def test_function(self):
        actual = ping._get_ping_command_args_list(
            ping.PingCommandParams(
                host_or_ip='google.com',
                num_packets=4,
                timeout_secs=5,
            ),
        )
        self.assertEqual(
            ['ping', 'google.com', '-c', '4', '-w', '5'],
            actual,
        )

        actual = ping._get_ping_command_args_list(
            ping.PingCommandParams(
                host_or_ip='google.com',
                num_packets=None,
                timeout_secs=None,
            ),
        )
        self.assertEqual(
            ['ping', 'google.com', '-c', '4', '-w', '20'],
            actual,
        )


class PingParseTests(unittest.TestCase):

    def setUp(self):
        self.param = ping.PingCommandParams(
            host_or_ip='google.com',
            num_packets=4,
            timeout_secs=None,
        )

    def test_parse_with_errors(self):
        param = ping.PingCommandParams(
            host_or_ip='localhost',
            num_packets=None,
            timeout_secs=None,
        )
        actual = ping.parse_ping_output('test', 'test', param)
        expected = ping.PingCommandResult(
            error='test',
            host_or_ip='localhost',
            num_packets=ping.DEFAULT_NUM_PACKETS,
            stats=None,
        )
        self.assertEqual(expected, actual)

    def test_parse_good_output(self):
        output = textwrap.dedent('''
        PING google.com (172.217.3.206) 56(84) bytes of data.
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=1 ttl=63 time=27.3 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=2 ttl=63 time=26.7 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=3 ttl=63 time=27.8 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=4 ttl=63 time=29.0 ms

        --- google.com ping statistics ---
        4 packets transmitted, 4 received, 0% packet loss, time 3008ms
        rtt min/avg/max/mdev = 26.780/27.731/29.014/0.832 ms
        ''').strip().encode('ascii')

        actual = ping.parse_ping_output(output, '', self.param)
        self.assertEqual(4, actual.stats.packets_transmitted)
        self.assertEqual(4, actual.stats.packets_received)
        self.assertEqual(0, actual.stats.packet_loss_pct)
        self.assertEqual(26.780, actual.stats.rtt_min)
        self.assertEqual(27.731, actual.stats.rtt_avg)
        self.assertEqual(29.014, actual.stats.rtt_max)
        self.assertEqual(0.832, actual.stats.rtt_mdev)

    def test_parse_no_header_line(self):
        output = textwrap.dedent('''
        PING google.com (172.217.3.206) 56(84) bytes of data.
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=1 ttl=63 time=27.3 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=2 ttl=63 time=26.7 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=3 ttl=63 time=27.8 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=4 ttl=63 time=29.0 ms

        4 packets transmitted, 4 received, 0% packet loss, time 3008ms
        rtt min/avg/max/mdev = 26.780/27.731/29.014/0.832 ms
        ''').strip().encode('ascii')

        expected = ping.PingCommandResult(
            error='Could not find statistics header in ping output',
            host_or_ip='google.com',
            num_packets=4,
            stats=None,
        )
        actual = ping.parse_ping_output(output, '', self.param)
        self.assertEqual(expected, actual)

    def test_parse_no_packet_line_match(self):
        output = textwrap.dedent('''
        PING google.com (172.217.3.206) 56(84) bytes of data.
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=1 ttl=63 time=27.3 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=2 ttl=63 time=26.7 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=3 ttl=63 time=27.8 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=4 ttl=63 time=29.0 ms

        --- google.com ping statistics ---
        4 packets transmitted, b received, 0% packet loss, time 3008ms
        rtt min/avg/max/mdev = 26.780/27.731/29.014/0.832 ms
        ''').strip().encode('ascii')
        packet_line = '4 packets transmitted, b received, ' \
                      '0% packet loss, time 3008ms'
        expected_error_msg = 'Could not parse packet line:' \
                             '\n{packet_line}'.format(packet_line=packet_line)

        expected = ping.PingCommandResult(
            error=expected_error_msg,
            host_or_ip='google.com',
            num_packets=4,
            stats=None,
        )
        actual = ping.parse_ping_output(output, '', self.param)
        self.assertEqual(expected, actual)

    def test_parse_no_rtt_line_match(self):
        output = textwrap.dedent('''
        PING google.com (172.217.3.206) 56(84) bytes of data.
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=1 ttl=63 time=27.3 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=2 ttl=63 time=26.7 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=3 ttl=63 time=27.8 ms
        64 bytes from sea15s12-in-f206.1e100.net (172.217.3.206): icmp_seq=4 ttl=63 time=29.0 ms

        --- google.com ping statistics ---
        4 packets transmitted, 4 received, 0% packet loss, time 3008ms
        rtt min/avg/max/mdev = a/27.731/29.014/0.832 ms
        ''').strip().encode('ascii')
        rtt_line = 'rtt min/avg/max/mdev = a/27.731/29.014/0.832 ms'

        expected_error_msg = 'Could not parse rtt line:\n{rtt_line}'\
            .format(rtt_line=rtt_line)
        expected = ping.PingCommandResult(
            error=expected_error_msg,
            host_or_ip='google.com',
            num_packets=4,
            stats=None,
        )
        actual = ping.parse_ping_output(output, '', self.param)
        self.assertEqual(expected, actual)

    def test_parse_deadline_reached_no_results(self):
        output = textwrap.dedent('''
        PING google.com (172.217.3.206) 56(84) bytes of data.

        --- google.com ping statistics ---
        1 packets transmitted, 0 received, 100% packet loss, time 0ms
        ''').strip().encode('ascii')
        expected_error_msg = 'Not enough output lines in ping output. ' \
                             'The ping may have timed out.'

        expected = ping.PingCommandResult(
            error=expected_error_msg,
            host_or_ip='google.com',
            num_packets=4,
            stats=None,
        )
        actual = ping.parse_ping_output(output, '', self.param)
        self.assertEqual(expected, actual)


if __name__ == '__main__':
    unittest.main()
