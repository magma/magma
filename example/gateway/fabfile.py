"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import sys

sys.path.append('../../orc8r')
import tools.fab.dev_utils as dev_utils


def register_vm():
    """ Provisions the gateway vm with the cloud vm """
    dev_utils.register_generic_gateway('test', 'example')
