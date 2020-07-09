"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import json
import threading
from collections import namedtuple
from typing import Optional

from magma.pipelined.imsi import encode_imsi
from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisHashDict
from magma.common.redis.serializers import get_json_deserializer, \
    get_json_serializer


SubscriberRuleKey = namedtuple('SubscriberRuleKey', 'key_type imsi rule_id')


class RuleIDToNumMapper:
    """
    Rule ID to Number Mapper

    This class assigns integers to rule ids so that they can be identified in
    an openflow register. The methods can be called from multiple threads
    """

    def __init__(self):
        self.redis_cli = get_default_client()
        self._curr_rule_num = 1
        self._rule_nums_by_rule = RuleIDDict()
        self._rules_by_rule_num = RuleNameDict()
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
        self._version_by_imsi_and_rule = RuleVersionDict()
        self._lock = threading.Lock()  # write lock

    def _update_version_unsafe(self, imsi: str, rule_id: str):
        key = self._get_json_key(encode_imsi(imsi), rule_id)
        version = self._version_by_imsi_and_rule.get(key)
        if not version:
            version = 0
        self._version_by_imsi_and_rule[key] = \
            (version % self.VERSION_LIMIT) + 1

    def update_version(self, imsi: str, rule_id: Optional[str] = None):
        """
        Increment the version number for a given subscriber and rule. If the
        rule id is not specified, then all rules for the subscriber will be
        incremented.
        """
        encoded_imsi = encode_imsi(imsi)
        with self._lock:
            if rule_id is None:
                for k, v in self._version_by_imsi_and_rule.items():
                    _, i, _ = SubscriberRuleKey(*json.loads(k))
                    if i == encoded_imsi:
                        self._version_by_imsi_and_rule[k] = v + 1
            else:
                self._update_version_unsafe(imsi, rule_id)

    def get_version(self, imsi: str, rule_id: str) -> int:
        """
        Returns the version number given a subscriber and a rule.
        """
        key = self._get_json_key(encode_imsi(imsi), rule_id)
        with self._lock:
            version = self._version_by_imsi_and_rule.get(key)
            if version is None:
                version = 0
        return version

    def _get_json_key(self, imsi: str, rule_id: str):
        return json.dumps(SubscriberRuleKey('imsi_rule', imsi, rule_id))


class RuleIDDict(RedisHashDict):
    """
    RuleIDDict uses the RedisHashDict collection to store a mapping of
    rule name to rule id.
    Setting and deleting items in the dictionary syncs with Redis automatically
    """
    _DICT_HASH = "pipelined:rule_ids"

    def __init__(self):
        client = get_default_client()
        super().__init__(
            client,
            self._DICT_HASH,
            get_json_serializer(), get_json_deserializer())

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None


class RuleNameDict(RedisHashDict):
    """
    RuleNameDict uses the RedisHashDict collection to store a mapping of
    rule id to rule name.
    Setting and deleting items in the dictionary syncs with Redis automatically
    """
    _DICT_HASH = "pipelined:rule_names"

    def __init__(self):
        client = get_default_client()
        super().__init__(
            client,
            self._DICT_HASH,
            get_json_serializer(), get_json_deserializer())

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None


class RuleVersionDict(RedisHashDict):
    """
    RuleVersionDict uses the RedisHashDict collection to store a mapping of
    subscriber+rule_id to rule version.
    Setting and deleting items in the dictionary syncs with Redis automatically
    """
    _DICT_HASH = "pipelined:rule_versions"

    def __init__(self):
        client = get_default_client()
        super().__init__(
            client,
            self._DICT_HASH,
            get_json_serializer(), get_json_deserializer())

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None
