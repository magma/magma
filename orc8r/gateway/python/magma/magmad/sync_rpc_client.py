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
import logging
import queue
import random
import threading
import time
from typing import List

import grpc
import magma.magmad.events as magmad_events
from google.protobuf.json_format import MessageToJson
from magma.common.rpc_utils import is_grpc_error_retryable
from magma.common.service_registry import ServiceRegistry
from magma.magmad.proxy_client import ControlProxyHttpClient
from orc8r.protos.sync_rpc_service_pb2 import SyncRPCRequest, SyncRPCResponse
from orc8r.protos.sync_rpc_service_pb2_grpc import SyncRPCServiceStub


class SyncRPCClient(threading.Thread):
    """
    SyncRPCClient initiates a SyncRPCClient, and opens a bidirectional stream
    with the cloud.
    """

    RETRY_MAX_DELAY_SECS = 10  # seconds

    def __init__(
        self, loop, response_timeout: int,
        print_grpc_payload: bool = False,
    ):
        threading.Thread.__init__(self)
        # a synchronized queue
        self._response_queue = queue.Queue()
        self._loop = loop
        asyncio.set_event_loop(self._loop)
        # seconds to wait for an actual SyncRPCResponse to become available
        # before sending out a heartBeat
        self._response_timeout = response_timeout
        self._proxy_client = ControlProxyHttpClient()
        self.daemon = True
        self._current_delay = 0
        self._last_conn_time = 0
        self._conn_closed_table = {}  # mapping of req id -> conn closed
        self._print_grpc_payload = print_grpc_payload

    def run(self):
        """
        This is executed when the thread is started. It gets a connection to
        the cloud dispatcher, and calls its bidirectional streaming rpc
        EstablishSyncRPCStream(). process_streams should never return, and
        if it did, exception will be logged, and new connection to dispatcher
        will be attempted after RETRY_DELAY_SECS seconds.
        """
        while True:
            try:
                start_time = time.time()
                chan = ServiceRegistry.get_rpc_channel(
                    'dispatcher',
                    ServiceRegistry.CLOUD,
                )
                client = SyncRPCServiceStub(chan)
                self._set_connect_time()
                self.process_streams(client)
            except grpc.RpcError as err:
                if is_grpc_error_retryable(err):
                    logging.error(
                        "[SyncRPC] Transient gRPC error, retrying: %s",
                        err.details(),
                    )
                    self._retry_connect_sleep()
                    continue
                else:
                    logging.error(
                        "[SyncRPC] gRPC error: %s, reconnecting to "
                        "cloud.", err.details(),
                    )
                    self._cleanup_and_reconnect()
            except Exception as exp:  # pylint: disable=broad-except
                conn_time = time.time() - start_time
                logging.error("[SyncRPC] Error after %ds: %s", conn_time, exp)
                self._cleanup_and_reconnect()

    def process_streams(self, client: SyncRPCServiceStub) -> None:
        """
        Calls rpc function EstablishSyncRPCStream on the client to establish
        a stream with dispatcher in the cloud, processes all requests from
        the stream, and writes all responses to the stream.
        Args:
            client: a grpc client to dispatcher in the cloud.
        Returns:
            Should only return when an exception is encountered.
        """

        # call to bidirectional streaming grpc takes in an iterator,
        # and returns an iterator
        sync_rpc_requests = client.EstablishSyncRPCStream(
            self.send_sync_rpc_response(),
        )
        magmad_events.established_sync_rpc_stream()
        # forward incoming requests from cloud to control_proxy
        self.forward_requests(sync_rpc_requests)

    def send_sync_rpc_response(self):
        """
        Retrieve SyncRPCResponse from queue. If no response is available yet,
        block for at most response_timeout seconds, and send a heartBeat if
        timeout.
        Returns: A generator of SyncRPCResponse
        """
        while True:
            try:
                resp = self._response_queue.get(
                    block=True,
                    timeout=self._response_timeout,
                )
                yield resp
            except queue.Empty:
                # response_queue is empty, send heartbeat
                # as the function itself has no knowledge on when it's
                # the first time it's called
                # this heartbeat response could be periodically called
                logging.debug("[SyncRPC] Sending heartbeat")
                yield SyncRPCResponse(heartBeat=True)

    def forward_requests(
        self,
        sync_rpc_requests: List[SyncRPCRequest],
    ) -> None:
        """
        Send requests in the sync_rpc_requests iterator.
        Args:
            sync_rpc_requests: an iterator of SyncRPCRequest from cloud

        Returns: Should only return when server shuts the stream (reaches
        end of the iterator sync_rpc_requests, or encounters an error)

        """
        logging.info("[SyncRPC] Waiting for requests")
        while True:
            try:
                req = next(sync_rpc_requests)
                self.forward_request(req)
            except grpc.RpcError as err:
                logging.error(
                    "[SyncRPC] Failing to forward request, err: %s",
                    err.details(),
                )
                raise err

    def forward_request(self, request: SyncRPCRequest) -> None:
        self._print_grpc(request)
        if request.heartBeat:
            logging.info("[SyncRPC] Got heartBeat from cloud")
            return

        if request.connClosed:
            logging.debug("[SyncRPC] Got connClosed from cloud")
            self._conn_closed_table[request.reqId] = True
            return

        logging.debug("[SyncRPC] Got a request")
        asyncio.run_coroutine_threadsafe(
            self._proxy_client.send(
                request.reqBody,
                request.reqId,
                self._response_queue,
                self._conn_closed_table,
            ),
            self._loop,
        )

    def _retry_connect_sleep(self) -> None:
        """
        Sleep for a current delay amount of time, with random backoff
        If current delay is less than RETRY_MAX_DELAY_SECS, exponentially
        increase current delay. If it exceeds RETRY_MAX_DELAY_SECS, sleep for
        RETRY_MAX_DELAY_SECS
        """
        sleep_time = self._current_delay + (random.randint(0, 1000) / 1000)
        self._current_delay = min(
            2 * self._current_delay,
            self.RETRY_MAX_DELAY_SECS,
        )
        self._current_delay = max(self._current_delay, 1)
        time.sleep(sleep_time)

    def _set_connect_time(self) -> None:
        logging.info("[SyncRPC] Opening stream to cloud")
        self._current_delay = 0
        self._last_conn_time = time.time()

    def _cleanup_and_reconnect(self):
        """
        If the connection is terminated, wait for a period of time
        before connecting back to the cloud. Also clear the conn
        closed table since cloud may reuse req IDs, and clear
        current proxy client connections
        """
        self._conn_closed_table.clear()
        self._proxy_client.close_all_connections()
        self._retry_connect_sleep()
        magmad_events.disconnected_sync_rpc_stream()

    def _print_grpc(self, message):
        if self._print_grpc_payload:
            try:
                log_msg = "{} {}".format(
                    message.DESCRIPTOR.full_name,
                    MessageToJson(message),
                )
                # add indentation
                padding = 2 * ' '
                log_msg = ''.join(
                    "{}{}".format(padding, line)
                    for line in log_msg.splitlines(True)
                )

                log_msg = "GRPC message:\n{}".format(log_msg)
                logging.info(log_msg)
            except Exception as e:  # pylint: disable=broad-except
                logging.debug("Exception while trying to log GRPC: %s", e)
