"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access
from unittest import TestCase
from magma.enodebd.data_models.transform_for_enb import bandwidth


class TransformForMagmaTests(TestCase):
    def test_bandwidth(self) -> None:
        inp = 1.4
        out = bandwidth(inp)
        expected = 'n6'
        self.assertEqual(out, expected, 'Should work with a float')

        inp = 20
        out = bandwidth(inp)
        expected = 'n100'
        self.assertEqual(out, expected, 'Should work with an int')

        inp = 10
        out = bandwidth(inp)
        expected = 'n50'
        self.assertEqual(out, expected, 'Should work with int 10')
