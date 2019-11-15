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

        serde = RedisSerde('log_verbosity',
                           get_proto_serializer(),
                           get_proto_deserializer(LogVerbosity))
        self._flat_dict = RedisFlatDict(client, serde)

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
        self._flat_dict['key1'] = expected
        version = self._flat_dict.get_version("key1")
        actual = self._flat_dict['key1']
        self.assertEqual(1, version)
        self.assertEqual(expected, actual)

        # update proto
        self._flat_dict["key1"] = expected2
        version2 = self._flat_dict.get_version("key1")
        actual2 = self._flat_dict["key1"]
        actual3 = self._flat_dict.get("key1")
        self.assertEqual(2, version2)
        self.assertEqual(expected2, actual2)
        self.assertEqual(expected2, actual3)

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_missing_version(self):
        missing_version = self._flat_dict.get_version("key2")
        self.assertEqual(0, missing_version)

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_bad_key(self):
        expected = LogVerbosity(verbosity=2)
        self.assertRaises(ValueError, self._flat_dict.__setitem__,
                          'bad:key', expected)
        self.assertRaises(ValueError, self._flat_dict.__getitem__,
                          'bad:key')
        self.assertRaises(ValueError, self._flat_dict.__delitem__,
                          'bad:key')

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_delete(self):
        expected = LogVerbosity(verbosity=2)
        self._flat_dict['key3'] = expected

        actual = self._flat_dict['key3']
        self.assertEqual(expected, actual)

        del self._flat_dict['key3']
        self.assertRaises(KeyError, self._flat_dict.__getitem__,
                          'key3')
        self.assertEqual(None, self._flat_dict.get('key3'))

    @mock.patch("redis.Redis", MockRedis)
    def test_flat_clear(self):
        expected = LogVerbosity(verbosity=2)
        self._flat_dict['key3'] = expected

        actual = self._flat_dict['key3']
        self.assertEqual(expected, actual)

        self._flat_dict.clear()
        self.assertEqual(0, len(self._flat_dict.keys()))


if __name__ == "__main__":
    main()
