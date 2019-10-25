"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


class MockRedis(object):
    """
    MockRedis implements a mock Redis Server using an in-memory dictionary
    """

    redis = {}

    def __init__(self, host, port):
        self.host = host
        self.port = port

    def serialize_key(self, key):
        """ Serialize key to plaintext encoded as UTF-8 bytes. """
        return key.encode('utf-8')

    def deserialize_key(self, serialized):
        """ Deserialize key from plaintext encoded as UTF-8 bytes. """
        return serialized.decode('utf-8')  # Redis returns keys as bytes

    def delete(self, key):
        """Mock delete."""
        skey = self.serialize_key(key)
        if skey in self.redis:
            del self.redis[skey]
            return 1
        return 0

    def exists(self, key):
        """Mock exists."""
        return key in self.redis

    def get(self, key):
        """Mock get."""
        skey = self.serialize_key(key)
        return self.redis[skey] if skey in self.redis else None

    def set(self, key, value):
        """Mock set."""
        skey = self.serialize_key(key)
        self.redis[skey] = value

    def keys(self):
        """ Mock keys."""
        return list(self.redis.keys())

    def hget(self, hashkey, key):
        """Mock hget."""

        skey = self.serialize_key(key)
        if hashkey not in self.redis:
            return None
        return self.redis[hashkey][skey] if skey in self.redis[hashkey] \
            else None

    def hgetall(self, hashkey):
        """Mock hgetall."""

        if hashkey not in self.redis:
            return {}
        return self.redis[hashkey]

    def hlen(self, hashkey):
        """Mock hlen."""

        return 0 if hashkey not in self.redis else len(self.redis[hashkey])

    def hset(self, hashkey, key, value):
        """Mock hset."""

        skey = self.serialize_key(key)
        if hashkey not in self.redis:
            self.redis[hashkey] = {}
        self.redis[hashkey][skey] = value

    def hdel(self, hashkey, key):
        """ Mock hdel"""
        skey = self.serialize_key(key)
        if hashkey not in self.redis:
            return
        self.redis[hashkey].pop(skey)

    def pipeline(self):
        """ Mock pipline"""
        return MockRedisPipeline(self)

    # pylint: disable=unused-argument
    def transaction(self, func, *args, **kwargs):
        """ Mock transaction."""
        pipe = self.pipeline()
        func_value = func(pipe)
        pipe.execute()
        return func_value


class MockRedisPipeline(object):
    """Mock redis-python pipeline object"""

    def __init__(self, redis):
        """Initialize the object."""

        self.redis = redis

    def execute(self):
        """ no-op."""
        pass

    def delete(self, key):
        """ Mock delete."""

        self.redis.delete(key)
        return self

    def hget(self, hashkey, key):
        """Mock srem."""
        return self.redis.hget(hashkey, key)

    def hdel(self, hashkey, key):
        """ Mock hdel"""
        return self.redis.hdel(hashkey, key)
