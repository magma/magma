"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

class CloudInitException(Exception):
    def __init__(self, message):
        super().__init__(message)


class SshRemoteCommandException(Exception):
    def __init__(self, message):
        super().__init__(message)


class MagmaRequestException(Exception):
    def __init__(self, message):
        super().__init__(message)
