"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisHashDict, RedisFlatDict
from magma.common.redis.mocks.mock_redis import MockRedis
from magma.common.redis.serializers import get_proto_deserializer, \
    get_proto_serializer, RedisSerde
from orc8r.protos.service303_pb2 import LogVerbosity
from unittest import TestCase, main, mock


class RedisDictTests(TestCase):
    """
    Tests for the RedisHashDict and RedisFlatDict containers
    """
    @mock.patch("redis.Redis", MockRedis)
    def setUp(self):
        client = get_default_client()
        # Use arbitrary orc8r proto to test with
        self._hash_dict = RedisHashDict(
            client,
            "unittest",
            get_proto_serializer(),
            get_proto_deserializer(LogVerbosity))

        serdes = {}
        serdes['log_verbosity'] = RedisSerde('log_verbosity',
                                get_proto_serializer(),
                                get_proto_deserializer(LogVerbosity))
        self._flat_dict = RedisFlatDict(client, serdes)

    @mock.patch("redis.Redis", MockRedis)
    def test_hash_insert(self):
        expected = LogVerbosity(verbosity=0)
        expected2 = LogVerbosity(verbosity=1)

        # insert proto
        self._hash_dict['key1'] = expected
        version = self._hash_dict.get_version("key1")
        actual = self._hash_dict['key1']
        self.assertEqual(1, version)
        self.assertEqual(expected, actual)

        # update proto
        self._hash_dict['key1'] = expected2
        version2 = self._hash_dict.get_version("key1")
        actual2 = self._hash_dict['key1']
        self.assertEqual(2, version2)
        self.assertEqual(expected2, actual2)

    @mock.patch("redis.Redis", MockRedis)
    def test_missing_version(self):
        missing_version = self._hash_dict.get_version("key2")
        self.assertEqual(0, missing_version)

    @mock.patch("redis.Redis", MockRedis)
    def test_hash_delete(self):
        expected = LogVerbosity(verbosity=2)
        self._hash_dict['key3'] = expected

        actual = self._hash_dict['key3']
        self.assertEqual(expected, actual)

        self._hash_dict.pop('key3')
        self.assertRaises(KeyError, self._hash_dict.__getitem__, 'key3')

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_insert(self):
        expected = LogVerbosity(verbosity=5)
        expected2 = LogVerbosity(verbosity=1)

        # insert proto
        self._flat_dict['key1:log_verbosity'] = expected
        version = self._flat_dict.get_version("key1", "log_verbosity")
        actual = self._flat_dict['key1:log_verbosity']
        self.assertEqual(1, version)
        self.assertEqual(expected, actual)

        # update proto
        self._flat_dict["key1:log_verbosity"] = expected2
        version2 = self._flat_dict.get_version("key1", "log_verbosity")
        actual2 = self._flat_dict["key1:log_verbosity"]
        self.assertEqual(2, version2)
        self.assertEqual(expected2, actual2)

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_missing_version(self):
        missing_version = self._flat_dict.get_version("key2", "log_verbosity")
        self.assertEqual(0, missing_version)

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_invalid_key(self):
        expected = LogVerbosity(verbosity=5)
        self.assertRaises(ValueError, self._flat_dict.__setitem__, 'key3',
                          expected)

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_invalid_serde(self):
        expected = LogVerbosity(verbosity=5)
        self.assertRaises(ValueError, self._flat_dict.__setitem__,
                          'key3:missing_serde', expected)

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_delete(self):
        expected = LogVerbosity(verbosity=2)
        self._flat_dict['key3:log_verbosity'] = expected

        actual = self._flat_dict['key3:log_verbosity']
        self.assertEqual(expected, actual)

        self._flat_dict.pop('key3:log_verbosity')
        self.assertRaises(KeyError, self._flat_dict.__getitem__,
                          'key3:log_verbosity')


if __name__ == "__main__":
    main()
