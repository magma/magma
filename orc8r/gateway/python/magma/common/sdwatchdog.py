"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import logging
import os
import systemd.daemon
import time

from typing import List, Optional, Set, cast


class SDWatchdogTask(object):
    """
    This is a base class that provides functions for a specific task to ensure
    regular completion of the loop.

    A coroutine would be implemented such that the loop will regularly call
    self.SetSDWatchdogAlive() on each iteration to mark this task as correctly
    running, such that SDWatchdog will check each loop to ensure the task
    continues to run.

    Use SetSDWatchdogTimeout() to specify MAX period of each loop iteration.
    """

    def __init__(self) -> None:
        """
        Classes that implement SDWatchdogTask must call
        super().__init__()
        in its __init__()
        """

        self._sdwatchdog = {
            # time.time() of last completed loop, track for watchdogging
            "time_last_completed_loop": cast(Optional[float], None),
            "timeout": 120,
        }

    def SetSDWatchdogTimeout(self, timeout: float) -> None:
        self._sdwatchdog["timeout"] = timeout

    def SetSDWatchdogAlive(self) -> None:
        self._sdwatchdog["time_last_completed_loop"] = time.time()

    def notCompleted(self, current_time: float) -> bool:
        last_time = self._sdwatchdog["time_last_completed_loop"]

        if last_time is None:
            return True
        if last_time < current_time - (self._sdwatchdog["timeout"] or 120):
            return True
        return False


class SDWatchdog(object):
    """
    This is a task that utilizes systemd watchdog functionality.

    SDWatchdog() task is started automatically in run in common/service.run(),
    where it will look at every task in the loop to see if it is a subclass
    of SDWatchdogTask

    To enable systemd watchdog, add "WatchdogSec=60" in the [Service] section
    of the systemd service file.
    """

    def __init__(
            self,
            tasks: Optional[List[SDWatchdogTask]],
            update_status: bool=False,  # update systemd status field
            period: float=30) -> None:
        """
        coroutine that will check each task's time_last_completed_loop to
        ensure that it was updated every in the last timeout_s seconds.

        Perform check of each service every period seconds.
        """

        self.tasks = cast(Set[SDWatchdogTask], set())
        self.update_status = update_status
        self.period = period

        if tasks:
            for t in tasks:
                if not issubclass(type(t), SDWatchdogTask):
                    logging.warning(
                        "'%s' is not a 'SDWatchdogTask', skipping", repr(t))
                else:
                    self.tasks.add(t)

    @staticmethod
    def has_notify() -> bool:
        return os.getenv("NOTIFY_SOCKET") is not None

    async def run(self) -> None:
        """
        check tasks every self.period seconds to see if they have completed
        a loop within the last 'timeout' seconds. If so, sd notify WATCHDOG=1
        """
        if not self.has_notify():
            logging.warning("Missing 'NOTIFY_SOCKET' for SDWatchdog, skipping")
            return
        logging.info("Starting SDWatchdog...")
        while True:
            current_time = time.time()
            anyStuck = False
            for task in self.tasks:
                if task.notCompleted(current_time):
                    errmsg = "SDWatchdog service '%s' has not completed %s" % (
                        repr(task), time.asctime(time.gmtime(current_time)))
                    if self.update_status:
                        systemd.daemon.notify("STATUS=%s\n" % errmsg)
                    logging.info(errmsg)
                    anyStuck = True

            if not anyStuck:
                systemd.daemon.notify(
                    'STATUS=SDWatchdog success %s\n' %
                    time.asctime(time.gmtime(current_time)))
                systemd.daemon.notify("WATCHDOG=1")
                systemd.daemon.notify("READY=1")  # only active if Type=notify

            await asyncio.sleep(self.period)
