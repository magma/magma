"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import threading
from collections import defaultdict
from typing import Optional

from magma.pipelined.imsi import encode_imsi


class RuleIDToNumMapper:
    """
    Rule ID to Number Mapper

    This class assigns integers to rule ids so that they can be identified in
    an openflow register. The methods can be called from multiple threads
    """

    def __init__(self):
        self._curr_rule_num = 1
        self._rule_nums_by_rule = {}
        self._rules_by_rule_num = {}
        self._lock = threading.Lock()  # write lock

    def _register_rule(self, rule_id):
        """ NOT thread safe """
        rule_num = self._rule_nums_by_rule.get(rule_id)
        if rule_num is not None:
            return rule_num
        rule_num = self._curr_rule_num
        self._rule_nums_by_rule[rule_id] = rule_num
        self._rules_by_rule_num[rule_num] = rule_id
        self._curr_rule_num += 1
        return rule_num

    def get_rule_num(self, rule_id):
        with self._lock:
            return self._rule_nums_by_rule[rule_id]

    def get_or_create_rule_num(self, rule_id):
        with self._lock:
            rule_num = self._rule_nums_by_rule.get(rule_id)
            if rule_num is None:
                return self._register_rule(rule_id)
            return rule_num

    def get_rule_id(self, rule_num):
        with self._lock:
            return self._rules_by_rule_num[rule_num]


class SessionRuleToVersionMapper:
    """
    Session & Rule to Version Mapper

    This class assigns version numbers to rule id & subscriber id combinations
    that can be used in an openflow register. The methods can be called from
    multiple threads.
    """

    VERSION_LIMIT = 0xFFFFFFFF  # 32 bit unsigned int limit (inclusive)

    def __init__(self):
        self._version_by_imsi_and_rule = defaultdict(lambda: defaultdict(int))
        self._lock = threading.Lock()  # write lock

    def _update_version_unsafe(self, imsi: str, rule_id: str):
        encoded_imsi = encode_imsi(imsi)
        version = self._version_by_imsi_and_rule[encoded_imsi][
            rule_id]
        self._version_by_imsi_and_rule[encoded_imsi][rule_id] = \
            (version % self.VERSION_LIMIT) + 1

    def update_version(self, imsi: str, rule_id: Optional[str] = None):
        """
        Increment the version number for a given subscriber and rule. If the
        rule id is not specified, then all rules for the subscriber will be
        incremented.
        """
        with self._lock:
            if rule_id is None:
                for rule in self._version_by_imsi_and_rule[encode_imsi(imsi)]:
                    self._update_version_unsafe(imsi, rule)
            else:
                self._update_version_unsafe(imsi, rule_id)

    def get_version(self, imsi: str, rule_id: str) -> int:
        """
        Returns the version number given a subscriber and a rule.
        """
        with self._lock:
            return self._version_by_imsi_and_rule[encode_imsi(imsi)][rule_id]
