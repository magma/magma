"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


def decode_imsi(imsi64):
    """
    Convert from the compacted uint back to a string, using the second two bits
    to determine the padding
    Args:
        imsi64 - compacted representation of imsi with padding at end
    Returns:
        imsi string in the form IMSI00101...
    """
    prefix_len = (imsi64 >> 1) & 0x3
    return 'IMSI' + '0' * prefix_len + str(imsi64 >> 3)


def encode_imsi(imsi):
    """
    Convert a IMSI string to a uint + length. IMSI strings can contain two
    prefix zeros for test MCC and maximum fifteen digits. Bit 1 of the
    compacted uint is always 1, so that we can match on it set. Bits 2-3
    the compacted uint contain how many leading 0's are in the IMSI. For
    example, if the IMSI is 001010000000013, the first bit is 0b1, the second
    two bits would be 0b10 and the remaining bits would be 1010000000013 << 3
    Args:
        imsi - string representation of imsi
    Returns:
        int representation of imsi with padding amount at end
    """
    if imsi.startswith('IMSI'):
        imsi = imsi[4:]  # strip IMSI off of string
    prefix_len = len(imsi) - len(imsi.lstrip('0'))
    compacted = (int(imsi) << 2) | (prefix_len & 0x3)
    return compacted << 1 | 0x1
