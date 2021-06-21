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
from unittest import mock
from unittest.mock import MagicMock

import fakeredis
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from magma.pipelined.rule_mappers import SessionRuleToVersionMapper


class RuleMappersTest(unittest.TestCase):
    def setUp(self):
        # mock the get_default_client function used to return a fakeredis object
        func_mock = MagicMock(return_value=fakeredis.FakeStrictRedis())
        with mock.patch(
                'magma.pipelined.rule_mappers.get_default_client',
                func_mock):
            self._session_rule_version_mapper = SessionRuleToVersionMapper()
        self._session_rule_version_mapper._version_by_imsi_and_rule = {}

    def test_session_rule_version_mapper(self):
        rule_ids = ['rule1', 'rule2']
        imsi = 'IMSI12345'
        ip_addr = '1.2.3.4'
        self._session_rule_version_mapper.save_version(
            imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[0], 1)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[0]),
            1)

        self._session_rule_version_mapper.save_version(
            imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[1], 1)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[1]),
            1)

        self._session_rule_version_mapper.save_version(
            imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[0], 2)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[0]),
            2)

        # Test updating version for all rules of a subscriber
        self._session_rule_version_mapper.remove_all_ue_versions(imsi,
            convert_ipv4_str_to_ip_proto(ip_addr))

        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[0]),
            -1)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, convert_ipv4_str_to_ip_proto(ip_addr), rule_ids[1]),
            -1)

    def test_session_rule_version_mapper_cwf(self):
        rule_ids = ['rule1', 'rule2']
        imsi = 'IMSI12345'
        self._session_rule_version_mapper.save_version(
            imsi, None, rule_ids[0], 1)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, None, rule_ids[0]),
            1)

        self._session_rule_version_mapper.save_version(
            imsi, None, rule_ids[1], 1)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, None, rule_ids[1]),
            1)

        self._session_rule_version_mapper.save_version(
            imsi, None, rule_ids[0], 2)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, None, rule_ids[0]),
            2)

        # Test updating version for all rules of a subscriber
        self._session_rule_version_mapper.remove_all_ue_versions(imsi, None)

        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, None, rule_ids[0]),
            -1)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(
                imsi, None, rule_ids[1]),
            -1)

if __name__ == "__main__":
    unittest.main()
