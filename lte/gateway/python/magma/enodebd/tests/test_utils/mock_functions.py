"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from typing import Any


GET_IP_FROM_IF_PATH = \
    'magma.enodebd.device_config.configuration_init.get_ip_from_if'


def mock_get_ip_from_if(
    _iface_name: str,
    _preference: Any = None,
) -> str:
    return '192.168.60.142'
