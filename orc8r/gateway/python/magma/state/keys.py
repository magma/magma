"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from magma.common.misc_utils import get_gateway_hwid

def make_mem_key(device_id: str, state_type: str) -> str:
    """
    Create a key of the format <id>:<type>
    """
    return device_id + ":" + state_type


def make_scoped_device_id(idval: str, scope: str) -> str:
    """
    Create a deviceID of the format <id> for scope 'network'
    Otherwise create a key of the format <hwid>:<id> for 'gateway' or
    unrecognized scope.
    """
    if scope == "network":
        return idval
    else:
        return get_gateway_hwid() + ":" + idval
