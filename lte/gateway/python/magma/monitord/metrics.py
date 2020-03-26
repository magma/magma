#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.

from prometheus_client import Histogram

SUBSCRIBER_ICMP_LATENCY_MS = Histogram('subscriber_icmp_latency_ms',
                                  'Reported latency for subscriber '
                                  'in milliseconds',
                                  ['imsi'],
                                  buckets=[50, 100, 200, 500, 1000, 2000])
