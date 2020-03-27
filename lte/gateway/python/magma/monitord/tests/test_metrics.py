"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

from magma.common import metrics_export
from magma.monitord.metrics import SUBSCRIBER_ICMP_LATENCY_MS


class MetricTests(unittest.TestCase):
    """
    Tests for the Service303 metrics interface
    """
    def test_metrics_defined(self):
        """ Test that all metrics are defined in proto enum """
        SUBSCRIBER_ICMP_LATENCY_MS.labels('IMSI00000001').observe(10.33)

        metrics_protos = list(metrics_export.get_metrics())
        for metrics_proto in metrics_protos:
            if metrics_proto.name == "subscriber_latency_ms":
                metric = metrics_proto.metric[0]
                self.assertEqual(metric.histogram.sample_sum, 10.33)
                self.assertEqual(metric.label[0].value, 'IMSI00000001')
