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

import aioh2
import h2.events
from magma.common.sentry import SEND_TO_ERROR_MONITORING
from magma.common.service_registry import ServiceRegistry
from orc8r.protos.sync_rpc_service_pb2 import GatewayResponse, SyncRPCResponse


class ControlProxyHttpClient(object):
    """
    ControlProxyHttpClient is a httpclient sending request
    to the control proxy local port. It's used in SyncRPCClient
    for forwarding GatewayRequests from the cloud, and gets a GatewayResponse.
    """

    def __init__(self):
        self._connection_table = {}  # map req id -> client

    async def send(
        self, gateway_request, req_id, sync_rpc_response_queue,
        conn_closed_table,
    ):
        """
        Forwards the given request to the service provided
        in :authority and awaits a response. If a exception is
        raised, log the error and enqueue an empty SyncRPCResponse.
        Else, enqueue SyncRPCResponse(s) that contains the GatewayResponse.

        Args:
            gateway_request: gateway_request: A GatewayRequest that is
        defined in the sync_rpc_service.proto. It has fields gwId, authority,
        path, headers, and payload.
            req_id: request id that's associated with the response
            sync_rpc_response_queue: the response queue that responses
        will be put in
            conn_closed_table: table that maps req ids to if the conn is closed

        Returns: None.

        """
        client = await self._get_client(gateway_request.authority)

        # Small hack to set PingReceived to no-op because the log gets spammed
        # with KeyError messages since aioh2 doesn't have a handler for
        # PingReceived. Remove if future versions support it.
        # pylint: disable=protected-access
        if hasattr(h2.events, "PingReceived"):
            # Need the hasattr here because some older versions of h2 may not
            # have the PingReceived event
            client._event_handlers[h2.events.PingReceived] = lambda _: None
        # pylint: enable=protected-access

        if req_id in self._connection_table:
            logging.error(
                "[SyncRPC] proxy_client is already handling "
                "request ID %s", req_id,
            )
            sync_rpc_response_queue.put(
                SyncRPCResponse(
                    heartBeat=False,
                    reqId=req_id,
                    respBody=GatewayResponse(
                        err=str(
                            "request ID {} is already being handled"
                            .format(req_id),
                        ),
                    ),
                ),
            )
            client.close_connection()
            return
        self._connection_table[req_id] = client

        try:
            await client.wait_functional()
            req_headers = self._get_req_headers(
                gateway_request.headers,
                gateway_request.path,
                gateway_request.authority,
            )
            body = gateway_request.payload
            stream_id = await client.start_request(req_headers)
            await self._await_gateway_response(
                client, stream_id, body,
                req_id, sync_rpc_response_queue,
                conn_closed_table,
            )
        except ConnectionAbortedError:
            logging.error(
                "[SyncRPC] proxy_client connection "
                "terminated by cloud",
            )
        except Exception as e:  # pylint: disable=broad-except
            logging.error(
                "[SyncRPC] Exception in proxy_client: %s", e,
                extra=SEND_TO_ERROR_MONITORING,
            )
            sync_rpc_response_queue.put(
                SyncRPCResponse(
                    heartBeat=False, reqId=req_id,
                    respBody=GatewayResponse(err=str(e)),
                ),
            )
        finally:
            del self._connection_table[req_id]
            try:
                client.close_connection()
            except AttributeError as e:
                logging.error(
                    '[SyncRPC] Error while trying to close conn: %s',
                    str(e),
                )

    def close_all_connections(self):
        connections = list(self._connection_table.values())
        for client in connections:
            try:
                client.close_connection()
            except (ConnectionAbortedError, AttributeError) as e:
                logging.error(
                    '[SyncRPC] Error while trying to close conn: %s',
                    str(e),
                )
        self._connection_table.clear()

    @staticmethod
    async def _get_client(service):
        (ip, port) = ServiceRegistry.get_service_address(service)
        return await aioh2.open_connection(ip, port)

    async def _await_gateway_response(
        self, client, stream_id, body,
        req_id, response_queue,
        conn_closed_table,
    ):
        await client.send_data(stream_id, body, end_stream=True)

        resp_headers = await client.recv_response(stream_id)
        status = self._get_resp_status(resp_headers)

        curr_payload = await self._read_stream(
            client, stream_id, req_id,
            response_queue,
            conn_closed_table,
        )
        next_payload = await self._read_stream(
            client, stream_id, req_id,
            response_queue,
            conn_closed_table,
        )

        while True:
            trailers = await client.recv_trailers(stream_id) \
                if not next_payload else []
            headers = self._get_resp_headers(resp_headers, trailers)
            res = GatewayResponse(
                status=status, headers=headers,
                payload=curr_payload,
            )
            response_queue.put(
                SyncRPCResponse(heartBeat=False, reqId=req_id, respBody=res),
            )
            if not next_payload:
                break

            curr_payload = next_payload
            next_payload = await self._read_stream(
                client, stream_id, req_id,
                response_queue,
                conn_closed_table,
            )

    @staticmethod
    def _get_req_headers(raw_req_headers, path, authority):
        headers = [
            (":method", "POST"),
            (":scheme", "http"),
            (":path", path),
            (":authority", authority),
        ]
        for key, val in raw_req_headers.items():
            headers.append((key, val))
        return headers

    @staticmethod
    def _get_resp_status(raw_headers):
        return dict(raw_headers)[":status"]

    @staticmethod
    def _get_resp_headers(raw_headers, raw_trailers):
        """
        Concatenate raw_headers and raw_tailers into a new dict
        raw_headers: a list of headers
        raw_trailers: a dict of trailers

        Return: a dict of headers and trailers
        """
        headers_dict = dict(raw_headers)
        headers_dict.update(raw_trailers)
        return headers_dict

    @staticmethod
    async def _read_stream(
        client, stream_id, req_id, response_queue,
        conn_closed_table,
    ):
        """
        Attempt to read from the stream. If it times out, send a keepConnActive
        response to the response queue. If it continues to time out after a
        very long period of time, raise asyncio.TimeoutError. If the connection
        is closed by the client, raise ConnectionAbortedError.
        """

        async def try_read_stream():
            while True:
                try:
                    payload = await asyncio.wait_for(
                        client.read_stream(stream_id), timeout=10.0,
                    )
                    if conn_closed_table.get(req_id, False):
                        raise ConnectionAbortedError
                    return payload
                except asyncio.TimeoutError:
                    if conn_closed_table.get(req_id, False):
                        raise ConnectionAbortedError
                    response_queue.put(
                        SyncRPCResponse(
                            heartBeat=False,
                            reqId=req_id,
                            respBody=GatewayResponse(keepConnActive=True),
                        ),
                    )

        return await asyncio.wait_for(try_read_stream(), timeout=120.0)
