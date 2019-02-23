"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import unittest

from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows


class TestMetering(unittest.TestCase):
    METERING_TABLE = 3

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_default_metering_flows(self):
        datapath = get_datapath()

        print('Checking for default table 3 flows')
        flows = get_flows(datapath, {'table_id': self.METERING_TABLE,
                                     'priority': 0})
        self.assertEqual(len(flows), 2,
                         'There should be 2 default table 3 flows')


if __name__ == '__main__':
    unittest.main()
