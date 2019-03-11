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
from contextlib import suppress


class Job(abc.ABC):
    """
    This is a base class that provides functions for a specific task to
    ensure regular completion of the loop.

    A co-routine run must be implemented by a subclass.
    periodic() will call the co-routine at a regular interval set by
    self._interval.
    """

    def __init__(self, interval, loop=None) -> None:
        if loop is None:
            self._loop = asyncio.get_event_loop()
        else:
            self._loop = loop
        self._task = cast(Optional[asyncio.Task], None)
        self._interval = interval  # in seconds
        self._last_run = cast(Optional[float], None)
        self._timeout = cast(Optional[float], None)

    @abc.abstractmethod
    def _run(self):
        """
        Once implemented by a subclass, this function will contain the actual
        work of this Job.
        """
        pass

    def start(self) -> None:
        if self._task is None:
            self._task = self._loop.create_task(self._periodic())

    def stop(self) -> None:
        if self._task is not None:
            self._task.cancel()
            with suppress(asyncio.CancelledError):
                # Await task to execute it's cancellation
                self._loop.run_until_complete(self._task)
            self._task = None

    def set_timeout(self, timeout: float) -> None:
        self._timeout = timeout

    def set_interval(self, interval: int) -> None:
        self._interval = interval

    def heartbeat(self) -> None:
        # record time to keep track of iteration length
        self._last_run = time.time()

    def not_completed(self, current_time: float) -> bool:
        last_time = self._last_run

        if last_time is None:
            return True
        if last_time < current_time - (self._timeout or 120):
            return True
        return False

    async def _periodic(self) -> None:
        while True:
            self.heartbeat()
            try:
                await self._run()
            except Exception as exp:  # pylint: disable=broad-except
                logging.error("Exception from _run: %s", exp)
            await asyncio.sleep(self._interval)
