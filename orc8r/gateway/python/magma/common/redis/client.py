"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import redis

from magma.configuration.service_configs import get_service_config_value


def get_default_client():
    """
    Return a default redis client using the configured port in redis.yml
    """
    redis_port = get_service_config_value('redis', 'port', 6379)
    redis_addr = get_service_config_value('redis', 'bind', 'localhost')
    return redis.Redis(host=redis_addr, port=redis_port)
