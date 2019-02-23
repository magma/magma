"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


class BadConfigError(Exception):
    """
    Indicates that test config for launching ryu apps is invalid
    """
    pass


class ServiceRunningError(Exception):
    """
    Indicates that magma@pipelined service was running when trying to
    instantiate ryu apps
    """
    pass
