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

import asyncio
from typing import Any

from magma.common.job import Job
from magma.common.log_count_handler import MsgCounterHandler
from magma.common.metrics import SERVICE_ERRORS

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
