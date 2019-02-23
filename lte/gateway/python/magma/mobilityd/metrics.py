"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from prometheus_client import Counter


# Counters for IP address management
IP_ALLOCATED_TOTAL = Counter('ip_address_allocated',
                             'Total IP addresses allocated')
IP_RELEASED_TOTAL = Counter('ip_address_released',
                             'Total IP addresses released')
