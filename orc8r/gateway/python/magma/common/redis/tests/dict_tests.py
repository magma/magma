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
from unittest import TestCase, main

import fakeredis
from magma.common.redis.containers import RedisFlatDict, RedisHashDict
from magma.common.redis.serializers import (
    RedisSerde,
    get_proto_deserializer,
    get_proto_serializer,
)
from orc8r.protos.service303_pb2 import LogVerbosity


class RedisDictTests(TestCase):
    """
    Tests for the RedisHashDict and RedisFlatDict containers
    """

    def setUp(self):
        client = fakeredis.FakeStrictRedis()
        # Use arbitrary orc8r proto to test with
        self._hash_dict = RedisHashDict(
            client,
            "unittest",
            get_proto_serializer(),
            get_proto_deserializer(LogVerbosity),
        )

        serde = RedisSerde(
            'log_verbosity',
            get_proto_serializer(),
            get_proto_deserializer(LogVerbosity),
        )
        self._flat_dict = RedisFlatDict(client, serde)

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

    def test_missing_version(self):
        missing_version = self._hash_dict.get_version("key2")
        self.assertEqual(0, missing_version)

    def test_hash_delete(self):
        expected = LogVerbosity(verbosity=2)
        self._hash_dict['key3'] = expected

        actual = self._hash_dict['key3']
        self.assertEqual(expected, actual)

        self._hash_dict.pop('key3')
        self.assertRaises(KeyError, self._hash_dict.__getitem__, 'key3')

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

    def test_flat_missing_version(self):
        missing_version = self._flat_dict.get_version("key2")
        self.assertEqual(0, missing_version)

    def test_flat_bad_key(self):
        expected = LogVerbosity(verbosity=2)
        self.assertRaises(
            ValueError, self._flat_dict.__setitem__,
            'bad:key', expected,
        )
        self.assertRaises(
            ValueError, self._flat_dict.__getitem__,
            'bad:key',
        )
        self.assertRaises(
            ValueError, self._flat_dict.__delitem__,
            'bad:key',
        )

    def test_flat_delete(self):
        expected = LogVerbosity(verbosity=2)
        self._flat_dict['key3'] = expected

        actual = self._flat_dict['key3']
        self.assertEqual(expected, actual)

        del self._flat_dict['key3']
        self.assertRaises(
            KeyError, self._flat_dict.__getitem__,
            'key3',
        )
        self.assertEqual(None, self._flat_dict.get('key3'))

    def test_flat_clear(self):
        expected = LogVerbosity(verbosity=2)
        self._flat_dict['key3'] = expected

        actual = self._flat_dict['key3']
        self.assertEqual(expected, actual)

        self._flat_dict.clear()
        self.assertEqual(0, len(self._flat_dict.keys()))

    def test_flat_garbage_methods(self):
        expected = LogVerbosity(verbosity=2)
        expected2 = LogVerbosity(verbosity=3)

        key = "k1"
        key2 = "k2"
        bad_key = "bad_key"
        self._flat_dict[key] = expected
        self._flat_dict[key2] = expected2

        self._flat_dict.mark_as_garbage(key)
        is_garbage = self._flat_dict.is_garbage(key)
        self.assertTrue(is_garbage)
        is_garbage2 = self._flat_dict.is_garbage(key2)
        self.assertFalse(is_garbage2)

        self.assertEqual([key], self._flat_dict.garbage_keys())
        self.assertEqual([key2], self._flat_dict.keys())

        self.assertIsNone(self._flat_dict.get(key))
        self.assertEqual(expected2, self._flat_dict.get(key2))

        deleted = self._flat_dict.delete_garbage(key)
        not_deleted = self._flat_dict.delete_garbage(key2)
        self.assertTrue(deleted)
        self.assertFalse(not_deleted)

        self.assertIsNone(self._flat_dict.get(key))
        self.assertEqual(expected2, self._flat_dict.get(key2))

        with self.assertRaises(KeyError):
            self._flat_dict.is_garbage(bad_key)
        with self.assertRaises(KeyError):
            self._flat_dict.mark_as_garbage(bad_key)


if __name__ == "__main__":
    main()
