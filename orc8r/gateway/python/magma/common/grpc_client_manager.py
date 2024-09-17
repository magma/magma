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

import logging
import sys

import grpc
import psutil
from magma.common.service_registry import ServiceRegistry


class GRPCClientManager:
    def __init__(
        self, service_name: str, service_stub,
        max_client_reuse: int = 60,
    ):
        self._client = None
        self._service_stub = service_stub
        self._service_name = service_name
        self._num_client_use = 0
        self._max_client_reuse = max_client_reuse

    def get_client(self):
        """
        get_client returns a grpc client of the specified service in the cloud.
        it will return a recycled client until the client fails or the number
        of recycling reaches the max_client_use.
        """
        if self._client is None or \
                self._num_client_use > self._max_client_reuse:
            chan = ServiceRegistry.get_rpc_channel(
                self._service_name,
                ServiceRegistry.CLOUD,
            )
            self._client = self._service_stub(chan)
            self._num_client_use = 0

        self._num_client_use += 1
        return self._client

    def on_grpc_fail(self, err_code):
        """
        Try to reuse the grpc client if possible. We are yet to fix a
        grpc behavior, where if DNS request blackholes then the DNS request
        is retried infinitely even after the channel is deleted. To prevent
        running out of fds, we try to reuse the channel during such failures
        as much as possible.
        """
        if err_code != grpc.StatusCode.DEADLINE_EXCEEDED:
            # Not related to the DNS issue
            self._reset_client()
        if self._num_client_use >= self._max_client_reuse:
            logging.info('Max client reuse reached. Cleaning up client')
            self._reset_client()

            # Sanity check if we are not leaking fds
            proc = psutil.Process()
            max_fds, _ = proc.rlimit(psutil.RLIMIT_NOFILE)
            open_fds = proc.num_fds()
            logging.info('Num open fds: %d', open_fds)
            if open_fds >= (max_fds * 0.8):
                logging.error("Reached 80% of allowed fds. Restarting process")
                sys.exit(1)

    def _reset_client(self):
        self._client = None
        self._num_client_use = 0
