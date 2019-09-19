"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
from typing import Any
from .metrics import SERVICE_ERRORS
from .log_count_handler import MsgCounterHandler
from magma.common.job import Job


# How frequently to poll systemd for error logs, in seconds
POLL_INTERVAL = 10


class ServiceLogErrorReporter(Job):
    """ Reports the number of logged errors for the service """

    def __init__(
        self,
        loop: asyncio.BaseEventLoop,
        service_config: Any,
        handler: MsgCounterHandler,
    ) -> None:
        super().__init__(interval=POLL_INTERVAL, loop=loop)
        self._service_config = service_config
        self._handler = handler

    async def _run(self):
        error_count = self._handler.pop_error_count()
        SERVICE_ERRORS.inc(error_count)
