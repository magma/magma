"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import os


def is_dev_mode() -> bool:
    """
    Returns whether the environment is set for dev mode
    """
    return os.environ.get('MAGMA_DEV_MODE') == '1'


def is_docker_network_mode() -> bool:
    """
    Returns whether the environment is set for dev mode
    """
    return os.environ.get('DOCKER_NETWORK_MODE') == '1'
