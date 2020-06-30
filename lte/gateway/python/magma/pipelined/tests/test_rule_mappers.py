"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

from magma.pipelined.rule_mappers import SessionRuleToVersionMapper


class RuleMappersTest(unittest.TestCase):
    def setUp(self):
        self._session_rule_version_mapper = SessionRuleToVersionMapper()
        self._session_rule_version_mapper._version_by_imsi_and_rule = {}

    def test_session_rule_version_mapper(self):
        rule_ids = ['rule1', 'rule2']
        imsi = 'IMSI12345'
        self._session_rule_version_mapper.update_version(imsi, rule_ids[0])
        self.assertEqual(
            self._session_rule_version_mapper.get_version(imsi, rule_ids[0]),
            1)

        self._session_rule_version_mapper.update_version(imsi, rule_ids[1])
        self.assertEqual(
            self._session_rule_version_mapper.get_version(imsi, rule_ids[1]),
            1)

        self._session_rule_version_mapper.update_version(imsi, rule_ids[0])
        self.assertEqual(
            self._session_rule_version_mapper.get_version(imsi, rule_ids[0]),
            2)

        # Test updating version for all rules of a subscriber
        self._session_rule_version_mapper.update_version(imsi, None)

        self.assertEqual(
            self._session_rule_version_mapper.get_version(imsi, rule_ids[0]),
            3)
        self.assertEqual(
            self._session_rule_version_mapper.get_version(imsi, rule_ids[1]),
            2)


if __name__ == "__main__":
    unittest.main()
