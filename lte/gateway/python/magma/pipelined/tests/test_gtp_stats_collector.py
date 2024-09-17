"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import unittest

from magma.pipelined.gtp_stats_collector import _match_gtp_lines


class GtpStatsCollecotrTest(unittest.TestCase):

    def test_match_gtp_lines(self):
        """
        Verifies that correct lines are matched and groups are captured.
        """

        # test input taken from https://github.com/magma/magma/issues/8926
        test_lines = [
            "dhcp0        {}                                 {collisions=0, rx_bytes=2639990, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=62842, tx_bytes=25926, tx_dropped=0, tx_errors=0, tx_packets=369}",
            "eth1.1506    {}                                 {collisions=0, rx_bytes=586074454, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=1123201, tx_bytes=1654735320, tx_dropped=0, tx_errors=0, tx_packets=1681111}",
            "g_2509330a   {key=flow, remote_ip=\"10.51.9.37\"} {rx_bytes=737220973, rx_packets=552257, tx_bytes=16652981, tx_packets=277380}",
            "gtp0         {key=flow, remote_ip=flow}         {rx_bytes=0, rx_packets=0, tx_bytes=0, tx_packets=0}",
            "gtp_br0      {}                                 {collisions=0, rx_bytes=300, rx_crc_err=0, rx_dropped=5, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=4, tx_bytes=1076, tx_dropped=0, tx_errors=0, tx_packets=14}",
            "ipfix0       {}                                 {collisions=0, rx_bytes=376, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=5, tx_bytes=1076, tx_dropped=0, tx_errors=0, tx_packets=14}",
            "li_port      {}                                 {collisions=0, rx_bytes=300, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=4, tx_bytes=1076, tx_dropped=0, tx_errors=0, tx_packets=14}",
            "mtr0         {}                                 {collisions=0, rx_bytes=600, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=8, tx_bytes=1076, tx_dropped=0, tx_errors=0, tx_packets=14}",
            "patch-agw    {peer=patch-up}                    {rx_bytes=1629631587, rx_packets=1620630, tx_bytes=601730199, tx_packets=1122219}",
            "patch-up     {peer=patch-agw}                   {rx_bytes=16657013, rx_packets=277443, tx_bytes=737205499, tx_packets=551224}",
            "proxy_port   {}                                 {collisions=0, rx_bytes=4540665, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=54539, tx_bytes=1727478, tx_dropped=13, tx_errors=0, tx_packets=19937}",
            "uplink_br0   {}                                 {collisions=0, rx_bytes=2639914, rx_crc_err=0, rx_dropped=1, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=62841, tx_bytes=2664422, tx_dropped=0, tx_errors=0, tx_packets=63191}",
            "vlan_pop_in  {}                                 {collisions=0, rx_bytes=25996, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=370, tx_bytes=25996, tx_dropped=0, tx_errors=0, tx_packets=370}",
            "vlan_pop_out {}                                 {collisions=0, rx_bytes=25996, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=370, tx_bytes=25996, tx_dropped=0, tx_errors=0, tx_packets=370}",
        ]

        result = _match_gtp_lines(test_lines)

        expected = [
            {'Interface': 'g_2509330a', 'remote_ip': '10.51.9.37', 'rx_bytes': '737220973', 'tx_bytes': '16652981'},
            {'Interface': 'gtp0', 'remote_ip': 'flow', 'rx_bytes': '0', 'tx_bytes': '0'},
        ]

        self.assertEqual(expected, result)
