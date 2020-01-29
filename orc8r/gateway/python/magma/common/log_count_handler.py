"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

import logging


class MsgCounterHandler(logging.Handler):
    """ Register this handler to logging to count the logs by level """

    count_by_level = None

    def __init__(self, *args, **kwargs):
        super(MsgCounterHandler, self).__init__(*args, **kwargs)
        self.count_by_level = {}

    def emit(self, record: logging.LogRecord):
        level = record.levelname
        if (level not in self.count_by_level):
            self.count_by_level[level] = 0
        self.count_by_level[level] += 1

    def pop_error_count(self) -> int:
        error_count = self.count_by_level.get('ERROR', 0)
        self.count_by_level['ERROR'] = 0
        return error_count
