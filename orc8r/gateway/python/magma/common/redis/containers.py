"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import collections.abc as collections_abc
from copy import deepcopy
import redis
import redis_collections
from typing import Dict

from magma.common.redis.serializers import RedisSerde
from orc8r.protos.redis_pb2 import RedisState

# NOTE: these containers replace the serialization methods exposed by
# the redis-collection objects. Although the methods are hinted to be
# privately scoped, the method replacement is encouraged in the library's
# docs: http://redis-collections.readthedocs.io/en/stable/usage-notes.html


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
        default_factory=None,
        writeback=False,
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
            default_factory, redis=client, key=key, writeback=writeback)

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


class RedisFlatDict(collections_abc.MutableMapping):
    """
    Dict-like interface serializing elements to a Redis datastore. This
    dict stores key directly (i.e. without a hashmap).
    """

    def __init__(self, client: redis.Redis, serdes: Dict[str, RedisSerde]):
        """
        Args:
            client (redis.Redis): Redis client object
            serdes (): RedisSerdes for each type of object that can be stored
        """
        super().__init__()
        self.redis = client
        self.serdes = serdes

    def __len__(self):
        """Return the number of items in the dictionary."""
        return len(self.redis.keys())

    def __iter__(self):
        """Return an iterator over the keys of the dictionary."""
        for k in self.redis.keys():
            try:
                yield k.decode('utf-8')
            except AttributeError:
                yield k

    def __contains__(self, key):
        """Return ``True`` if *key* is present, else ``False``."""
        return bool(self.redis.exists(key))

    def __getitem__(self, key):
        """Return the item of dictionary with key *key*. Raises a
        :exc:`KeyError` if key is not in the map.
        """
        if ':' not in key:
            raise ValueError('key must be of format <id>:<type>')
        serde = self._get_serde(key)
        serialized_value = self.redis.get(key)
        if serialized_value is None:
            raise KeyError(key)

        value = serde.deserialize(serialized_value)
        return value

    def __setitem__(self, key, value):
        """Set ``d[key]`` to *value*."""
        if ':' not in key:
            raise ValueError('key must be of format <id>:<type>')
        serde = self._get_serde(key)
        version = self._get_version(key)
        serialized_value = serde.serialize(value, version + 1)

        self.redis.set(key, serialized_value)

        return self.redis.get(key)

    def __delitem__(self, key):
        """Remove ``d[key]`` from dictionary.
        Raises a :func:`KeyError` if *key* is not in the map.
        """
        deleted_count = self.redis.delete(key)
        if not deleted_count:
            raise KeyError(key)

    def clear(self):
        for key in self.keys():
            self.redis.delete(key)

    def get_version(self, idval, typeval):
        """Return the version of the value for key *key*. Returns 0 if
        key is not in the map
        """
        flat_key = idval + ":" + typeval
        return self._get_version(flat_key)

    def _get_version(self, key):
        value = self.redis.get(key)
        if value is None:
            return 0

        proto_wrapper = RedisState()
        proto_wrapper.ParseFromString(value)
        return proto_wrapper.version

    def _get_serde(self, key):
        parsed_key = key.split(':')
        if len(parsed_key) != 2:
            raise ValueError("Dictionary key must be of format <id>:<type>")
        typeval = parsed_key[1]

        if typeval not in self.serdes:
            raise ValueError("Dictionary is not configured for object type:"
                             " %s" % typeval)

        return self.serdes[typeval]
