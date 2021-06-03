"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import abc
import asyncio
import logging
import time
from contextlib import suppress
from typing import Optional, cast


class Job(abc.ABC):
    """
    This is a base class that provides functions for a specific task to
    ensure regular completion of the loop.

    A co-routine run must be implemented by a subclass.
    periodic() will call the co-routine at a regular interval set by
    self._interval.
    """

    def __init__(
            self,
            interval: int,
            loop: Optional[asyncio.AbstractEventLoop] = None,
    ) -> None:
        if loop is None:
            self._loop = asyncio.get_event_loop()
        else:
            self._loop = loop
        # Task in charge of periodically running the task
        self._periodic_task = cast(Optional[asyncio.Task], None)
        # Task in charge of deciding how long to wait until next run
        self._interval_wait_task = cast(Optional[asyncio.Task], None)
        self._interval = interval  # in seconds
        self._last_run = cast(Optional[float], None)
        self._timeout = cast(Optional[float], None)
        # Condition variable used to control how long the job waits until
        # executing its task again.
        self._cond = self._cond = asyncio.Condition(loop=self._loop)

    @abc.abstractmethod
    async def _run(self):
        """
        Once implemented by a subclass, this function will contain the actual
        work of this Job.
        """
        pass

    def start(self) -> None:
        """
        kicks off the _periodic while loop
        """
        if self._periodic_task is None:
            self._periodic_task = self._loop.create_task(self._periodic())

    def stop(self) -> None:
        """
        cancels the _periodic while loop
        """
        if self._periodic_task is not None:
            self._periodic_task.cancel()
            with suppress(asyncio.CancelledError):
                # Await task to execute it's cancellation
                self._loop.run_until_complete(self._periodic_task)
            self._periodic_task = None

    def set_timeout(self, timeout: float) -> None:
        self._timeout = timeout

    def set_interval(self, interval: int) -> None:
        """
        sets the interval used in _periodic to decide how long to sleep
        """
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

    async def _sleep_for_interval(self):
        await asyncio.sleep(self._interval)
        async with self._cond:
            self._cond.notify()

    async def wake_up(self):
        """
        Cancels the _sleep_for_interval task if it exists, and notifies the
        cond var so that the _periodic loop can continue.
        """
        if self._interval_wait_task is not None:
            self._interval_wait_task.cancel()

        async with self._cond:
            self._cond.notify()

    async def _periodic(self) -> None:
        while True:
            self.heartbeat()

            try:
                await self._run()
            except Exception as exp:  # pylint: disable=broad-except
                logging.exception("Exception from _run: %s", exp)

            # Wait for self._interval seconds or wake_up is explicitly called
            self._interval_wait_task = \
                self._loop.create_task(self._sleep_for_interval())
            async with self._cond:
                await self._cond.wait()
