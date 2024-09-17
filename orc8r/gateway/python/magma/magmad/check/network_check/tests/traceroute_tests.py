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

from magma.magmad.check.network_check import traceroute

# Allow access to protected variables for unit testing
# pylint: disable=protected-access


class TracerouteArgFactoryTests(unittest.TestCase):
    def test_function(self):
        actual = traceroute._get_traceroute_command_args_list(
            traceroute.TracerouteParams(
                host_or_ip='google.com',
                max_hops=10,
                bytes_per_packet=30,
            ),
        )
        self.assertEqual(
            ['traceroute', '-m', '10', 'google.com', '30'],
            actual,
        )

        actual = traceroute._get_traceroute_command_args_list(
            traceroute.TracerouteParams(
                host_or_ip='google.com', max_hops=None, bytes_per_packet=None,
            ),
        )
        self.assertEqual(
            ['traceroute', '-m', '30', 'google.com', '60'],
            actual,
        )


class TracerouteParsingTests(unittest.TestCase):
    def setUp(self):
        self.param = traceroute.TracerouteParams(
            host_or_ip='google.com',
            max_hops=30,
            bytes_per_packet=60,
        )

    @staticmethod
    def _configure_output(text):
        return textwrap.dedent(text).strip().encode('ascii')

    def test_parse_with_stderr(self):
        actual = traceroute.parse_traceroute_output('test', 'test', self.param)
        expected = traceroute.TracerouteResult(
            error='test',
            host_or_ip='google.com',
            stats=None,
        )
        self.assertEqual(expected, actual)

    def test_parse_good_output(self):
        """
        Some test cases with truncated output from `traceroute`
        run in vagrant VM
        """
        output = '''
        traceroute to google.com (216.58.195.238), 30 hops max, 60 byte packets
         1  10.0.2.2 (10.0.2.2)  0.332 ms  0.262 ms  0.223 ms
         2  mpk14-21-off-dgw1-vl301.corp.tfbnw.net (172.25.164.2)  0.926 ms * *
         7  prn1-sidf-cenet-cgw1-be6.corp.tfbnw.net (172.24.4.7)  31.895 ms  25.105 ms  31.396 ms
        10  prn1-sidf-cenet-fw1-ae4.corp.tfbnw.net (172.24.4.23)  30.686 ms prn1-sidf-cenet-fw1-ae8.corp.tfbnw.net (172.24.4.25)  23.317 ms prn1-sidf-cenet-fw1-ae4.corp.tfbnw.net (172.24.4.23)  30.546 ms
        '''
        actual = traceroute.TracerouteParser().parse(
            self._configure_output(output),
        )
        expected = traceroute.TracerouteStats(
            [
                traceroute.TracerouteHop(
                    1, [
                        traceroute.TracerouteProbe(
                            '10.0.2.2', '10.0.2.2', 0.332,
                        ),
                        traceroute.TracerouteProbe(
                            '10.0.2.2', '10.0.2.2', 0.262,
                        ),
                        traceroute.TracerouteProbe(
                            '10.0.2.2', '10.0.2.2', 0.223,
                        ),
                    ],
                ),
                traceroute.TracerouteHop(
                    2, [
                        traceroute.TracerouteProbe(
                            'mpk14-21-off-dgw1-vl301.corp.tfbnw.net',
                            '172.25.164.2', 0.926,
                        ),
                        traceroute.TracerouteProbe(
                            'mpk14-21-off-dgw1-vl301.corp.tfbnw.net',
                            '172.25.164.2', 0,
                        ),
                        traceroute.TracerouteProbe(
                            'mpk14-21-off-dgw1-vl301.corp.tfbnw.net',
                            '172.25.164.2', 0,
                        ),
                    ],
                ),
                traceroute.TracerouteHop(
                    7, [
                        traceroute.TracerouteProbe(
                            'prn1-sidf-cenet-cgw1-be6.corp.tfbnw.net',
                            '172.24.4.7', 31.895,
                        ),
                        traceroute.TracerouteProbe(
                            'prn1-sidf-cenet-cgw1-be6.corp.tfbnw.net',
                            '172.24.4.7', 25.105,
                        ),
                        traceroute.TracerouteProbe(
                            'prn1-sidf-cenet-cgw1-be6.corp.tfbnw.net',
                            '172.24.4.7', 31.396,
                        ),
                    ],
                ),
                traceroute.TracerouteHop(
                    10, [
                        traceroute.TracerouteProbe(
                            'prn1-sidf-cenet-fw1-ae4.corp.tfbnw.net',
                            '172.24.4.23', 30.686,
                        ),
                        traceroute.TracerouteProbe(
                            'prn1-sidf-cenet-fw1-ae8.corp.tfbnw.net',
                            '172.24.4.25', 23.317,
                        ),
                        traceroute.TracerouteProbe(
                            'prn1-sidf-cenet-fw1-ae4.corp.tfbnw.net',
                            '172.24.4.23', 30.546,
                        ),
                    ],
                ),
            ],
        )
        self.assertEqual(expected, actual)

    def test_parse_non_num_rtt(self):
        output = '''
        traceroute to google.com (216.58.195.238), 30 hops max, 60 byte packets
         2  mpk14-21-off-dgw1-vl301.corp.tfbnw.net (172.25.164.2)  abc ms * *
        '''
        with self.assertRaises(ValueError):
            traceroute.TracerouteParser().parse(self._configure_output(output))

    def test_parse_incomplete_lookahead(self):
        output = '''
        traceroute to google.com (216.58.195.238), 30 hops max, 60 byte packets
         7  prn1-sidf-cenet-cgw1-be6.corp.tfbnw.net (172.24.4.7)  31.895 ms  25.105
        '''
        with self.assertRaises(IndexError):
            traceroute.TracerouteParser().parse(self._configure_output(output))

    def test_parse_no_response(self):
        output = '''
        traceroute to google.com (216.58.195.238), 30 hops max, 60 byte packets
         4  * * *
         5  * * *
        '''
        actual = traceroute.TracerouteParser().parse(
            self._configure_output(output),
        )
        expected = traceroute.TracerouteStats([
            traceroute.TracerouteHop(
                4,
                [traceroute.TracerouteProbe(None, None, 0)] * 3,
            ),
            traceroute.TracerouteHop(
                5,
                [traceroute.TracerouteProbe(None, None, 0)] * 3,
            ),
        ])
        self.assertEqual(expected, actual)
