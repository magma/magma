"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from prometheus_client import Counter

# Define a new counter, which can be incremented in the rpc servicer
NUM_REQUESTS = Counter('num_requests', 'Total requests')
