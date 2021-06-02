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
# pylint: disable=W0223

import asyncio
import logging
import os
import time
from typing import List, Optional, Set, cast

import systemd.daemon
from magma.common.job import Job


class SDWatchdogTask(Job):
    pass


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
            update_status: bool = False,  # update systemd status field
            period: float = 30,
    ) -> None:
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
                        "'%s' is not a 'SDWatchdogTask', skipping", repr(t),
                    )
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
                if task.not_completed(current_time):
                    errmsg = "SDWatchdog service '%s' has not completed %s" % (
                        repr(task), time.asctime(time.gmtime(current_time)),
                    )
                    if self.update_status:
                        systemd.daemon.notify("STATUS=%s\n" % errmsg)
                    logging.info(errmsg)
                    anyStuck = True

            if not anyStuck:
                systemd.daemon.notify(
                    'STATUS=SDWatchdog success %s\n' %
                    time.asctime(time.gmtime(current_time)),
                )
                systemd.daemon.notify("WATCHDOG=1")
                systemd.daemon.notify("READY=1")  # only active if Type=notify

            await asyncio.sleep(self.period)
