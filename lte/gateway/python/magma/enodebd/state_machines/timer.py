"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from datetime import datetime, timedelta


class StateMachineTimer():
    def __init__(self, seconds_remaining: int) -> None:
        self.start_time = datetime.now()
        self.seconds = seconds_remaining

    def is_done(self) -> bool:
        time_elapsed = datetime.now() - self.start_time
        if time_elapsed > timedelta(seconds=self.seconds):
            return True
        return False

    def seconds_elapsed(self) -> int:
        time_elapsed = datetime.now() - self.start_time
        return int(time_elapsed.total_seconds())

    def seconds_remaining(self) -> int:
        return max(0, self.seconds - self.seconds_elapsed())