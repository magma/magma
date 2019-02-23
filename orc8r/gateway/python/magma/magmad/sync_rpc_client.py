"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import logging
import queue
import threading
import time

import grpc
from orc8r.protos.sync_rpc_service_pb2 import GatewayResponse, SyncRPCResponse
from orc8r.protos.sync_rpc_service_pb2_grpc import SyncRPCServiceStub

from magma.common.service_registry import ServiceRegistry
from magma.magmad.proxy_client import ControlProxyHttpClient


class SyncRPCClient(threading.Thread):
    """
    SyncRPCClient initiates a SyncRPCClient, and opens a bidirectional stream
    with the cloud.
    """

    RETRY_MAX_DELAY_SECS = 10  # seconds

    def __init__(self, loop, response_timeout):
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
                chan = ServiceRegistry.get_rpc_channel('dispatcher',
                                                       ServiceRegistry.CLOUD)
                client = SyncRPCServiceStub(chan)
                self._set_connect_time()
                self.process_streams(client)
            except Exception as exp:  # pylint: disable=broad-except
                conn_time = time.time() - start_time
                logging.error("[SyncRPC] Error after %ds: %s", conn_time, exp)
            # If the connection is terminated, wait for a period of time
            # before connecting back to the cloud.
            self._retry_connect_sleep()

    def process_streams(self, client):
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
            self.send_sync_rpc_response())
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
                resp = self._response_queue.get(block=True,
                                                timeout=self._response_timeout)
                logging.debug("[SyncRPC] Sending response")
                yield resp
            except queue.Empty:
                # response_queue is empty, send heartbeat
                # as the function itself has no knowledge on when it's
                # the first time it's called
                # this heartbeat response could be periodically called
                logging.debug("[SyncRPC] Sending heartbeat")
                yield SyncRPCResponse(heartBeat=True)

    def forward_requests(self, sync_rpc_requests):
        """
        Send requests in the sync_rpc_requests iterator.
        Args:
            sync_rpc_requests: an iterator of SyncRPCRequest from cloud

        Returns: Should only return when server shuts the stream (reaches
        end of the iterator sync_rpc_requests, or encounters an error)

        """
        try:
            while True:
                logging.debug("[SyncRPC] Waiting for requests")
                req = next(sync_rpc_requests)
                self.forward_request(req)
        except grpc.RpcError as err:
            # server end closed connection; retry rpc connection.
            raise Exception("Error when retrieving request: [%s] %s"
                            % (err.code(), err.details()))

    def forward_request(self, request):
        if request.heartBeat:
            logging.info("[SyncRPC] Got heartBeat from cloud")
            return
        try:
            logging.debug("[SyncRPC] Got a request")
            future = asyncio.run_coroutine_threadsafe(
                self._proxy_client.send(request.reqBody),
                self._loop)
            future.add_done_callback(lambda fut:
                                     self._loop.call_soon_threadsafe(
                                         self.send_request_done, request.reqId,
                                         fut))
        except Exception as exp:  # pylint: disable=broad-except
            logging.error("[SyncRPC] Error when forwarding request: %s", exp)

    def send_request_done(self, req_id, future):
        """
        A future that has a GatewayResponse is done. Check if a exception is
        raised. If so, log the error and enqueue an empty SyncRPCResponse.
        Else, enqueue a SyncRPCResponse that contains the GatewayResponse that
        became available in the future.
        Args:
            req_id: request id that's associated with the response
            future: A future that contains a GatewayResponse that is done.

        Returns: None

        """
        err = future.exception()
        if err:
            logging.error("[SyncRPC] Forward to control proxy error: %s", err)
            self._response_queue.put(
                SyncRPCResponse(heartBeat=False, reqId=req_id,
                                respBody=GatewayResponse(err=str(err))))
        else:
            res = future.result()
            self._response_queue.put(
                SyncRPCResponse(heartBeat=False, reqId=req_id, respBody=res))

    def _retry_connect_sleep(self):
        """
        Sleep for a current delay amount of time.
        If last connection time was over 60 seconds ago, sleep for 0 seconds.
        If current delay is less than RETRY_MAX_DELAY_SECS, exponentially
        increase current delay. If it exceeds RETRY_MAX_DELAY_SECS, sleep for
        RETRY_MAX_DELAY_SECS
        """
        # if last connect time was over 60 secs ago, reset current_delay to 0
        if time.time() - self._last_conn_time > 60:
            self._current_delay = 0
        elif self._current_delay == 0:
            self._current_delay = 1
        else:
            self._current_delay = min(2 * self._current_delay,
                                      self.RETRY_MAX_DELAY_SECS)
        time.sleep(self._current_delay)

    def _set_connect_time(self):
        logging.info("[SyncRPC] Opening stream to cloud")
        self._last_conn_time = time.time()
