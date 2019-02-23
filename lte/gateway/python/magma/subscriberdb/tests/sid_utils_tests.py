"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

from lte.protos.subscriberdb_pb2 import SubscriberID

from magma.subscriberdb.sid import SIDUtils


class SIDTests(unittest.TestCase):
    """
    Tests for the SID utilities
    """

    def test_str_conversion(self):
        """
        Tests the string conversion utils
        """
        sid = SubscriberID(id='12345', type=SubscriberID.IMSI)
        self.assertEqual(SIDUtils.to_str(sid), 'IMSI12345')
        self.assertEqual(SIDUtils.to_pb('IMSI12345'), sid)

        # By default the id type is IMSI
        sid = SubscriberID(id='12345')
        self.assertEqual(SIDUtils.to_str(sid), 'IMSI12345')
        self.assertEqual(SIDUtils.to_pb('IMSI12345'), sid)

        # Raise ValueError if invalid strings are given
        with self.assertRaises(ValueError):
            SIDUtils.to_pb('IMS')

        with self.assertRaises(ValueError):
            SIDUtils.to_pb('IMSI12345a')

        with self.assertRaises(ValueError):
            SIDUtils.to_pb('')


if __name__ == "__main__":
    unittest.main()
