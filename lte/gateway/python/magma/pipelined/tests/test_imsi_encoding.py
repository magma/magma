"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

from magma.pipelined.imsi import encode_imsi, decode_imsi


class IMSIEncodingTest(unittest.TestCase):
    """
    Test the encoding and decoding of IMSI values used in flow metadata in
    pipelined
    """

    def test_leading_zeros(self):
        imsi_list = [
            "IMSI001010000000013",
            "IMSI011010000000013",
            "IMSI111010000000013",
        ]
        for imsi in imsi_list:
            self._assert_conversion(imsi)

    def test_different_lengths(self):
        imsi_list = [
            "IMSI001010000000013",
            "IMSI01010000000013",
            "IMSI28950000000013",
        ]
        for imsi in imsi_list:
            self._assert_conversion(imsi)

    def _assert_conversion(self, imsi):
        """
        To test the validity of the conversion, we can encode the imsi and
        decode to see if the two values match
        """
        test_val = decode_imsi(encode_imsi(imsi))
        self.assertEqual(imsi, test_val)


if __name__ == "__main__":
    unittest.main()
