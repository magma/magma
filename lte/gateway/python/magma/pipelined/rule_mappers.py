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
import json
import threading
from collections import namedtuple

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisFlatDict, RedisHashDict
from magma.common.redis.serializers import (
    RedisSerde,
    get_json_deserializer,
    get_json_serializer,
)
from magma.pipelined.imsi import encode_imsi

SubscriberRuleKey = namedtuple('SubscriberRuleKey', 'key_type imsi teid rule_id')


class RuleIDToNumMapper:
    """
    Rule ID to Number Mapper

    This class assigns integers to rule ids so that they can be identified in
    an openflow register. The methods can be called from multiple threads
    """

    def __init__(self):
        self.redis_cli = get_default_client()
        self._curr_rule_num = 1
        self._rule_nums_by_rule = {}
        self._rules_by_rule_num = {}
        self._lock = threading.Lock()  # write lock

    def setup_redis(self):
        self._rule_nums_by_rule = RuleIDDict()
        self._rules_by_rule_num = RuleNameDict()

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

    def __init__(self):
        self._version_by_imsi_and_rule = {}
        self._lock = threading.Lock()  # write lock

    def _save_version_unsafe(self, imsi: str, teid_str: str, rule_id: str,
                             version):
        key = self._get_json_key(encode_imsi(imsi), teid_str, rule_id)
        self._version_by_imsi_and_rule[key] = version

    def update_all_ue_versions(self, imsi: str, teid: int):
        """
        Increment the version number for a given subscriber and rule. If the
        rule id is not specified, then all rules for the subscriber will be
        incremented.
        """
        encoded_imsi = encode_imsi(imsi)
        if teid is None or teid == 0:
            teid_str = ""
        else:
            teid_str = str(teid)
        with self._lock:
            for k, v in self._version_by_imsi_and_rule.items():
                _, imsi, teid_val, _ = SubscriberRuleKey(*json.loads(k))
                if imsi == encoded_imsi and teid_val == teid_str:
                    self._version_by_imsi_and_rule[k] = v + 1

    def save_version(self, imsi: str, teid: int, rule_id: [str], version: int):
        """
        Increment the version number for a given subscriber and rule. If the
        rule id is not specified, then all rules for the subscriber will be
        incremented.
        """
        if teid is None or teid == 0:
            teid_str = ""
        else:
            teid_str = str(teid)
        with self._lock:
            self._save_version_unsafe(imsi, teid_str, rule_id, version)

    def get_version(self, imsi: str, teid: int, rule_id: str) -> int:
        """
        Returns the version number given a subscriber and a rule.
        """
        if teid is None or teid == 0:
            teid_str = ""
        else:
            teid_str = str(teid)
        key = self._get_json_key(encode_imsi(imsi), teid_str, rule_id)
        with self._lock:
            version = self._version_by_imsi_and_rule.get(key)
            if version is None:
                version = 0
        return version

    def remove(self, imsi: str, teid: int, rule_id: str, version: int):
        """
        Removed the element from redis if the passed version matches the
        current one
        """
        if teid is None or teid == 0:
            teid_str = ""
        else:
            teid_str = str(teid)
        key = self._get_json_key(encode_imsi(imsi), teid_str, rule_id)
        with self._lock:
            cur_version = self._version_by_imsi_and_rule.get(key)
            if version is None:
                return
            if cur_version == version:
                del self._version_by_imsi_and_rule[key]

    def _get_json_key(self, imsi: str, teid_str: str, rule_id: str):
        return json.dumps(SubscriberRuleKey('imsi_rule', imsi, teid_str,
                                            rule_id))


class RuleIDDict(RedisFlatDict):
    """
    RuleIDDict uses the RedisHashDict collection to store a mapping of
    rule name to rule id.
    Setting and deleting items in the dictionary syncs with Redis automatically
    """
    _DICT_HASH = "pipelined:rule_ids"

    def __init__(self):
        client = get_default_client()
        serde = RedisSerde(self._DICT_HASH, get_json_serializer(),
                           get_json_deserializer())
        super().__init__(client, serde, writethrough=True)

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


class RuleVersionDict(RedisFlatDict):
    """
    RuleVersionDict uses the RedisHashDict collection to store a mapping of
    subscriber+rule_id to rule version.
    Setting and deleting items in the dictionary syncs with Redis automatically
    """
    _DICT_HASH = "pipelined:rule_versions"

    def __init__(self):
        client = get_default_client()
        serde = RedisSerde(self._DICT_HASH, get_json_serializer(),
                           get_json_deserializer())
        super().__init__(client, serde, writethrough=True)

    def __missing__(self, key):
        """Instead of throwing a key error, return None when key not found"""
        return None
