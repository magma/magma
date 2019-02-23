"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import unittest

from spyne import ComplexModelBase

from magma.enodebd.tr069.models import DeviceIdStruct


class DeviceIdStructTests(unittest.TestCase):
    def test_as_dict_memory_leak(self):
        """
        Test to ensure as_dict() doesn't leak model instances
        """
        thing = DeviceIdStruct(Manufacturer='abc',
                               OUI='def',
                               ProductClass='ghi',
                               SerialNumber='jkl')
        res = thing.as_dict()
        self.assertEqual(
            {'Manufacturer': 'abc',
             'OUI': 'def',
             'ProductClass': 'ghi',
             'SerialNumber': 'jkl'},
            res,
        )
        # inspect the spyne.util.memoize object that wraps the staticmethod
        self.assertEqual(1, len(ComplexModelBase.get_flat_type_info.memo))

        # should produce a different result and not grow the size of memo
        thing.OUI = 'aaaa'
        res = thing.as_dict()
        self.assertEqual(
            {'Manufacturer': 'abc',
             'OUI': 'aaaa',
             'ProductClass': 'ghi',
             'SerialNumber': 'jkl'},
            res,
        )
        self.assertEqual(1, len(ComplexModelBase.get_flat_type_info.memo))

        # use a different object this time. Again should not grow memo
        thing = DeviceIdStruct(Manufacturer='abc',
                               OUI='def',
                               ProductClass='ghi',
                               SerialNumber='jkl')
        res = thing.as_dict()
        self.assertEqual(
            {'Manufacturer': 'abc',
             'OUI': 'def',
             'ProductClass': 'ghi',
             'SerialNumber': 'jkl'},
            res,
        )
        self.assertEqual(1, len(ComplexModelBase.get_flat_type_info.memo))
