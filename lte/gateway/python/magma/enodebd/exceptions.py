"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

class ConfigurationError(Exception):
    """ Indicates that the eNodeB could not be configured correctly. """
    pass


class Tr069Error(Exception):
    pass


class IncorrectDeviceHandlerError(Exception):
    """ Indicates that we're using the wrong data model for configuration. """
    def __init__(self, device_name: str):
        """
        device_name: What device we actually are dealing with
        """
        super().__init__()
        self.device_name = device_name


class UnrecognizedEnodebError(Exception):
    """
    Indicates that the Access Gateway does not recognize the eNodeB.
    The Access Gateway will not interact with the eNodeB in question.
    """
    pass
