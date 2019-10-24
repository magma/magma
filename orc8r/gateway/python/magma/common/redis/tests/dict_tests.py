"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisDict
from magma.common.redis.mocks.mock_redis import MockRedis
from magma.common.redis.serializers import get_proto_deserializer, \
    get_proto_serializer
from orc8r.protos.service303_pb2 import LogVerbosity
from unittest import TestCase, main, mock


class RedisTests(TestCase):
    """
    Tests for the RedisDict container
    """
    @mock.patch("redis.Redis", MockRedis)
    def setUp(self):
        # Use arbitrary orc8r proto to test with
        self._dict = RedisDict(
            get_default_client(),
            "unittest",
            get_proto_serializer(),
            get_proto_deserializer(LogVerbosity))

    @mock.patch("redis.Redis", MockRedis)
    def test_insert(self):
        expected = LogVerbosity(verbosity=0)
        expected2 = LogVerbosity(verbosity=1)

        # insert proto
        self._dict['key1'] = expected
        version = self._dict.get_version("key1")
        actual = self._dict['key1']
        self.assertEqual(1, version)
        self.assertEqual(expected, actual)

        # update proto
        self._dict['key1'] = expected2
        version2 = self._dict.get_version("key1")
        actual2 = self._dict['key1']
        self.assertEqual(2, version2)
        self.assertEqual(expected2, actual2)

    @mock.patch("redis.Redis", MockRedis)
    def test_missing_version(self):
        missing_version = self._dict.get_version("key2")
        self.assertEqual(0, missing_version)

    @mock.patch("redis.Redis", MockRedis)
    def test_delete(self):
        expected = LogVerbosity(verbosity=2)
        self._dict['key3'] = expected

        actual = self._dict['key3']
        self.assertEqual(expected, actual)

        self._dict.pop('key3')
        self.assertRaises(KeyError, self._dict.__getitem__, 'key3')


if __name__ == "__main__":
    main()
