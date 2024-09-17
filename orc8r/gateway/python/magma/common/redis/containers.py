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
from copy import deepcopy
from typing import Any, Iterator, List, MutableMapping, Optional, TypeVar

import redis
import redis_collections
import redis_lock
from magma.common.redis.serializers import RedisSerde
from orc8r.protos.redis_pb2 import RedisState
from redis.lock import Lock

# NOTE: these containers replace the serialization methods exposed by
# the redis-collection objects. Although the methods are hinted to be
# privately scoped, the method replacement is encouraged in the library's
# docs: http://redis-collections.readthedocs.io/en/stable/usage-notes.html

T = TypeVar('T')


class RedisList(redis_collections.List):
    """
    List-like interface serializing elements to a Redis datastore.

    Notes:
        - Provides persistence across sessions
        - Mutable elements handled correctly
        - Not expected to be thread safe, but could be extended
    """

    def __init__(self, client, key, serialize, deserialize):
        """
        Initialize instance.

        Args:
            client (redis.Redis): Redis client object
            key (str): key where this container's elements are stored in Redis
            serialize (function (any) -> bytes):
                function called to serialize an element
            deserialize (function (bytes) -> any):
                function called to deserialize an element
        Returns:
            redis_list (redis_collections.List): persistent list-like interface
        """
        self._pickle = serialize
        self._unpickle = deserialize
        super().__init__(redis=client, key=key, writeback=True)

    def __copy__(self):
        return [elt for elt in self]

    def __deepcopy__(self, memo):
        return [deepcopy(elt, memo) for elt in self]


class RedisSet(redis_collections.Set):
    """
    Set-like interface serializing elements to a Redis datastore.

    Notes:
        - Provides persistence across sessions
        - Mutable elements _not_ handled correctly:
            - Get/set mutable elements supported
            - Don't update the contents of a mutable element and
              expect things to go well
        - Expected to be thread safe, but not tested
    """

    def __init__(self, client, key, serialize, deserialize):
        """
        Initialize instance.

        Args:
            client (redis.Redis): Redis client object
            key (str): key where this container's elements are stored in Redis
            serialize (function (any) -> bytes):
                function called to serialize an element
            deserialize (function (bytes) -> any):
                function called to deserialize an element
        Returns:
            redis_set (redis_collections.Set): persistent set-like interface
        """
        # NOTE: redis_collections.Set doesn't have a writeback option, causing
        # issue when mutable elements are updated in-place.
        self._pickle = serialize
        self._unpickle = deserialize
        super().__init__(redis=client, key=key)

    def __copy__(self):
        return {elt for elt in self}

    def __deepcopy__(self, memo):
        return {deepcopy(elt, memo) for elt in self}


class RedisHashDict(redis_collections.DefaultDict):
    """
    Dict-like interface serializing elements to a Redis datastore. This dict
    utilizes Redis's hashmap functionality

    Notes:
        - Keys must be string-like and are serialized to plaintext (UTF-8)
        - Provides persistence across sessions
        - Mutable elements handled correctly
        - Not expected to be thread safe, but could be extended
        - Keys are serialized in plaintext
    """

    @staticmethod
    def serialize_key(key):
        """ Serialize key to plaintext. """
        return key

    @staticmethod
    def deserialize_key(serialized):
        """ Deserialize key from plaintext encoded as UTF-8 bytes. """
        return serialized.decode('utf-8')  # Redis returns bytes

    def __init__(
            self, client, key, serialize, deserialize,
            default_factory=None, writeback=False,
    ):
        """
        Initialize instance.

        Args:
            client (redis.Redis): Redis client object
            key (str): key where this container's elements are stored in Redis
            serialize (function (any) -> bytes):
                function called to serialize a value
            deserialize (function (bytes) -> any):
                function called to deserialize a value
            default_factory: function that provides default value for a
                non-existent key
            writeback (bool): if writeback is set to true, dict maintains a
                local cache of values and the `sync` method can be called to
                store these values. NOTE: only use this option if syncing
                between services is not important.

        Returns:
            redis_dict (redis_collections.Dict): persistent dict-like interface
        """
        # Key serialization (to/from plaintext)
        self._pickle_key = RedisHashDict.serialize_key
        self._unpickle_key = RedisHashDict.deserialize_key
        # Value serialization
        self._pickle_value = serialize
        self._unpickle = deserialize
        super().__init__(
            default_factory, redis=client, key=key, writeback=writeback,
        )

    def __setitem__(self, key, value):
        """Set ``d[key]`` to *value*.

        Override in order to increment version on each update
        """
        version = self.get_version(key)
        pickled_key = self._pickle_key(key)
        pickled_value = self._pickle_value(value, version + 1)
        self.redis.hset(self.key, pickled_key, pickled_value)

        if self.writeback:
            self.cache[key] = value

    def __copy__(self):
        return {key: self[key] for key in self}

    def __deepcopy__(self, memo):
        return {key: deepcopy(self[key], memo) for key in self}

    def get_version(self, key):
        """Return the version of the value for key *key*. Returns 0 if
        key is not in the map
        """
        try:
            value = self.cache[key]
        except KeyError:
            pickled_key = self._pickle_key(key)
            value = self.redis.hget(self.key, pickled_key)
            if value is None:
                return 0

        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(value)
        return proto_wrapper.version


class RedisFlatDict(MutableMapping[str, T]):
    """
    Dict-like interface serializing elements to a Redis datastore. This
    dict stores key directly (i.e. without a hashmap).
    """

    def __init__(
        self, client: redis.Redis, serde: RedisSerde[T],
        writethrough: bool = False,
    ):
        """
        Args:
            client (redis.Redis): Redis client object
            serde (): RedisSerde for de/serializing the object stored
            writethrough (bool): if writethrough is set to true,
            RedisFlatDict maintains a local write-through cache of values.
        """
        super().__init__()
        self._writethrough = writethrough
        self.redis = client
        self.serde = serde
        self.redis_type = serde.redis_type
        self.cache = {}
        if self._writethrough:
            self._sync_cache()

    def __len__(self) -> int:
        """Return the number of items in the dictionary."""
        if self._writethrough:
            return len(self.cache)

        return len(self.keys())

    def __iter__(self) -> Iterator[str]:
        """Return an iterator over the keys of the dictionary."""
        type_pattern = self._get_redis_type_pattern()

        if self._writethrough:
            for k in self.cache:
                split_key, _ = k.split(":", 1)
                yield split_key
        else:
            for k in self.redis.keys(pattern=type_pattern):
                try:
                    deserialized_key = k.decode('utf-8')
                    split_key = deserialized_key.split(":", 1)
                except AttributeError:
                    split_key = k.split(":", 1)
                # There could be a delete key in between KEYS and GET, so ignore
                # invalid values for now
                try:
                    if self.is_garbage(split_key[0]):
                        continue
                except KeyError:
                    continue
                yield split_key[0]

    def __contains__(self, key: str) -> bool:
        """Return ``True`` if *key* is present and not garbage,
        else ``False``.
        """
        composite_key = self._make_composite_key(key)

        if self._writethrough:
            return composite_key in self.cache

        return bool(self.redis.exists(composite_key)) and \
            not self.is_garbage(key)

    def __getitem__(self, key: str) -> T:
        """Return the item of dictionary with key *key:type*. Raises a
        :exc:`KeyError` if *key:type* is not in the map or the object is
        garbage
        """
        if ':' in key:
            raise ValueError("Key %s cannot contain ':' char" % key)
        composite_key = self._make_composite_key(key)

        if self._writethrough:
            cached_value = self.cache.get(composite_key)
            if cached_value:
                return cached_value

        serialized_value = self.redis.get(composite_key)
        if serialized_value is None:
            raise KeyError(composite_key)

        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(serialized_value)
        if proto_wrapper.is_garbage:
            raise KeyError("Key %s is garbage" % key)

        return self.serde.deserialize(serialized_value)

    def __setitem__(self, key: str, value: T) -> Any:
        """Set ``d[key:type]`` to *value*."""
        if ':' in key:
            raise ValueError("Key %s cannot contain ':' char" % key)
        version = self.get_version(key)
        serialized_value = self.serde.serialize(value, version + 1)
        composite_key = self._make_composite_key(key)
        if self._writethrough:
            self.cache[composite_key] = value
        return self.redis.set(composite_key, serialized_value)

    def __delitem__(self, key: str) -> int:
        """Remove ``d[key:type]`` from dictionary.
        Raises a :func:`KeyError` if *key:type* is not in the map.
        """
        if ':' in key:
            raise ValueError("Key %s cannot contain ':' char" % key)
        composite_key = self._make_composite_key(key)
        if self._writethrough:
            del self.cache[composite_key]
        deleted_count = self.redis.delete(composite_key)
        if not deleted_count:
            raise KeyError(composite_key)
        return deleted_count

    def get(self, key: str, default=None) -> Optional[T]:
        """Get ``d[key:type]`` from dictionary.
        Returns None if *key:type* is not in the map
        """
        try:
            return self.__getitem__(key)
        except (KeyError, ValueError):
            return default

    def clear(self) -> None:
        """
        Clear all keys in the dictionary. Objects are immediately deleted
        (i.e. not garbage collected)
        """
        if self._writethrough:
            self.cache.clear()
        for key in self.keys():
            composite_key = self._make_composite_key(key)
            self.redis.delete(composite_key)

    def get_version(self, key: str) -> int:
        """Return the version of the value for key *key:type*. Returns 0 if
        key is not in the map
        """
        composite_key = self._make_composite_key(key)
        value = self.redis.get(composite_key)
        if value is None:
            return 0

        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(value)
        return proto_wrapper.version

    def keys(self) -> List[str]:
        """Return a copy of the dictionary's list of keys
        Note: for redis *key:type* key is returned
        """
        if self._writethrough:
            return list(self.cache.keys())

        return list(self.__iter__())

    def mark_as_garbage(self, key: str) -> Any:
        """Mark ``d[key:type]`` for garbage collection
        Raises a KeyError if *key:type* is not in the map.
        """
        composite_key = self._make_composite_key(key)
        value = self.redis.get(composite_key)
        if value is None:
            raise KeyError(composite_key)

        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(value)
        proto_wrapper.is_garbage = True
        garbage_serialized = proto_wrapper.SerializeToString()
        return self.redis.set(composite_key, garbage_serialized)

    def is_garbage(self, key: str) -> bool:
        """Return if d[key:type] has been marked for garbage collection.
        Raises a KeyError if *key:type* is not in the map.
        """
        composite_key = self._make_composite_key(key)
        value = self.redis.get(composite_key)
        if value is None:
            raise KeyError(composite_key)

        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(value)
        return proto_wrapper.is_garbage

    def garbage_keys(self) -> List[str]:
        """Return a copy of the dictionary's list of keys that are garbage
        Note: for redis *key:type* key is returned
        """
        garbage_keys = []
        type_pattern = self._get_redis_type_pattern()
        for k in self.redis.keys(pattern=type_pattern):
            try:
                deserialized_key = k.decode('utf-8')
                split_key = deserialized_key.split(":", 1)
            except AttributeError:
                split_key = k.split(":", 1)
            # There could be a delete key in between KEYS and GET, so ignore
            # invalid values for now
            try:
                if not self.is_garbage(split_key[0]):
                    continue
            except KeyError:
                continue
            garbage_keys.append(split_key[0])
        return garbage_keys

    def delete_garbage(self, key) -> bool:
        """Remove ``d[key:type]`` from dictionary iff the object is garbage
        Returns False if *key:type* is not in the map
        """
        if not self.is_garbage(key):
            return False
        count = self.__delitem__(key)
        return count > 0

    def lock(self, key: str) -> Lock:
        """Lock the dictionary for key *key*"""
        return redis_lock.Lock(
            self.redis,
            name=self._make_composite_key(key) + ":lock",
            expire=60,
            auto_renewal=True,
            strict=False,
        )

    def _sync_cache(self):
        """
        Syncs write-through cache with redis data on store.
        """
        type_pattern = self._get_redis_type_pattern()
        for k in self.redis.keys(pattern=type_pattern):
            composite_key = k.decode('utf-8')
            serialized_value = self.redis.get(composite_key)
            value = self.serde.deserialize(serialized_value)
            self.cache[composite_key] = value

    def _get_redis_type_pattern(self):
        return "*:" + self.redis_type

    def _make_composite_key(self, key):
        return key + ":" + self.redis_type
