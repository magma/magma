"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import hashlib


def encode_apn(apn):
    """Converts string definition of APN with random length in to fixed length byte array.

    Args:
        apn: APN string
    """
    md_hasher = hashlib.md5()
    md_hasher.update(bytearray(apn))
    return md_hasher.digest()
