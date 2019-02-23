"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import abc
import logging
import asyncio
import time
from typing import Optional, cast


class Job(abc.ABC):
    """
    This is a base class that provides functions for a specific task to
    ensure regular completion of the loop.

    A co-routine run must be implemented by a subclass.
    periodic() will call the co-routine at a regular interval set by
    self._interval.
    """

    def __init__(self, interval) -> None:
        self._loop = asyncio.get_event_loop()
        self._task = cast(Optional[asyncio.Task], None)
        self._interval = interval  # in seconds
        self._last_run = cast(Optional[float], None)

    @abc.abstractmethod
    def run(self):
        """
        Once implemented by a subclass, this function will contain the actual
        work of this Job.
        """
        pass

    async def start(self) -> None:
        if self._task is None:
            self._task = self._loop.create_task(self._periodic())

    async def stop(self) -> None:
        if self._task is not None:
            self._task.cancel()
            self._task = None

    def set_interval(self, interval: int) -> None:
        self._interval = interval

    def _heartbeat(self) -> None:
        # record time to keep track of iteration length
        self._last_run = time.time()

    async def _periodic(self) -> None:
        while True:
            try:
                self._heartbeat()
                self.run()
            except Exception as exp:  # pylint: disable=broad-except
                logging.error("Error getting service status: %s", exp)
            await asyncio.sleep(self._interval)
