"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

from magma.subscriberdb.crypto.gsm import UnsafePreComputedA3A8
from magma.subscriberdb.crypto.utils import CryptoError


class CryptoTests(unittest.TestCase):
    """
    Test class for the Crypto algorithms
    """

    def test_precomputed_a3a8(self):
        """
        Test if the UnsafePrecomputedA3A8 algo works as expected
        """
        crypto = UnsafePreComputedA3A8()
        rand = b'ni\x89\xbel\xeeqTT7p\xae\x80\xb1\xef\r'
        sres = b'\xd4\xac\x8bS'
        cipher_key = b'\x9f\xf54.\xb9]\x88\x00'

        # Tuple would be returned for a properly encoded input
        input_k = rand + sres + cipher_key
        self.assertEqual(
            crypto.generate_auth_tuple(input_k),
            (rand, sres, cipher_key))

        # If the length is not 28 bytes, CryptoError will be thrown
        input_k = rand + sres
        with self.assertRaises(CryptoError):
            crypto.generate_auth_tuple(input_k)


if __name__ == "__main__":
    unittest.main()
