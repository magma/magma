"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from copy import deepcopy
import redis_collections

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


class RedisDict(redis_collections.DefaultDict):
    """
    Dict-like interface serializing elements to a Redis datastore.

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
        self._pickle_key = RedisDict.serialize_key
        self._unpickle_key = RedisDict.deserialize_key
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
