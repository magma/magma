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
from collections import namedtuple

QosInfo = namedtuple('QosInfo', 'gbr mbr')

# key_type - identifies the type of QosKey. This ideally should be in base type
# Will modify this when dataclasses are available
SubscriberRuleKey = namedtuple('SubscriberRuleKey', 'key_type imsi ip_addr rule_num direction')


def get_subscriber_key(*args):
    keyType = "Subscriber"
    return SubscriberRuleKey(keyType, *args)


def get_key_json(key):
    return json.dumps(key)


def get_key(json_val):
    return SubscriberRuleKey(*json.loads(json_val))


SubscriberRuleData = namedtuple('SubscriberRuleData', 'data_type qid ambr leaf')


def get_subscriber_data(*args):
    keyType = "qos_data"
    return SubscriberRuleData(keyType, *args)


def get_data_json(key):
    return json.dumps(key)


def get_data(json_val):
    return SubscriberRuleData(*json.loads(json_val))

