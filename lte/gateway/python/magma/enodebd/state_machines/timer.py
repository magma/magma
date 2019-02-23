"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from datetime import datetime, timedelta


class StateMachineTimer():
    def __init__(self, seconds_remaining):
        self.start_time = datetime.now()
        self.seconds = seconds_remaining

    def is_done(self):
        time_elapsed = datetime.now() - self.start_time
        if time_elapsed > timedelta(seconds=self.seconds):
            return True
        return False
