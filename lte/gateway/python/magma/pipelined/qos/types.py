"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import json
from collections import namedtuple

QosInfo = namedtuple('QosInfo', 'gbr mbr')

# key_type - identifies the type of QosKey. This ideally should be in base type
# Will modify this when dataclasses are available
SubscriberKey = namedtuple('SubscriberKey', 'key_type imsi rule_num direction')


def get_subscriber_key(*args):
    keyType = "Subscriber"
    return SubscriberKey(keyType, *args)


def get_json(key):
    return json.dumps(key)


def get_key(json_val):
    return SubscriberKey(*json.loads(json_val))
