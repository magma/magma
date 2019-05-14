"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access
from unittest import TestCase
from magma.enodebd.data_models.transform_for_magma import gps_tr181, bandwidth
from magma.enodebd.exceptions import ConfigurationError


class TransformForMagmaTests(TestCase):
    def test_gps_tr181(self) -> None:
        # Negative longitude
        inp = '-122150583'
        out = gps_tr181(inp)
        expected = '-122.150583'
        self.assertEqual(out, expected, 'Should convert negative longitude')

        inp = '122150583'
        out = gps_tr181(inp)
        expected = '122.150583'
        self.assertEqual(out, expected, 'Should convert positive longitude')

        inp = '0'
        out = gps_tr181(inp)
        expected = '0.0'
        self.assertEqual(out, expected, 'Should leave zero as zero')

    def test_bandwidth(self) -> None:
        inp = 'n6'
        out = bandwidth(inp)
        expected = 1.4
        self.assertEqual(out, expected, 'Should convert RBs')

        inp = 1.4
        out = bandwidth(inp)
        expected = 1.4
        self.assertEqual(out, expected, 'Should accept MHz')

        with self.assertRaises(ConfigurationError):
            inp = 'asdf'
            bandwidth(inp)

        with self.assertRaises(ConfigurationError):
            inp = 1234
            bandwidth(inp)
