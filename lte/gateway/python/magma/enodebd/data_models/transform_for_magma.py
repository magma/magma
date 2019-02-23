"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import textwrap

DUPLEX_MAP = {
    '01': 'TDDMode',
    '02': 'FDDMode'
}


def duplex_mode(value):
    return DUPLEX_MAP.get(value)


def band_capability(value):
    return ','.join([str(int(b, 16)) for b in textwrap.wrap(value, 2)])


def gps_tr181(value):
    """Convert GPS value (lat or lng) to float

    Per TR-181 specification, coordinates are returned in degrees,
    multiplied by 1,000,000.

    Args:
        value (string): GPS value (latitude or longitude)
    Returns:
        str: GPS value (latitude/longitude) in degrees
    """
    if value.isnumeric():
        return str(float(value) / 1e6)
    return value
