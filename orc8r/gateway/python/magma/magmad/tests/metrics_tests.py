"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

from magma.common import metrics_export


class MetricTests(unittest.TestCase):
    """
    Tests for the Service303 metrics interface
    """
    def test_metrics_defined(self):
        """ Test that all metrics are defined in proto enum """
        import magma.magmad.metrics
        # Avoid lint error about unused imports
        magma.magmad.metrics.CPU_PERCENT.set(1)

        metrics_protos = metrics_export.get_metrics()
        for metrics_proto in metrics_protos:
            # Check that all proto names have been mapped to numbers. Will
            # raise ValueError if not
            int(metrics_proto.name)
