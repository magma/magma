"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging

import aioh2
import h2.events
from orc8r.protos.sync_rpc_service_pb2 import GatewayResponse

from magma.common.service_registry import ServiceRegistry


class ControlProxyHttpClient(object):
    """
    ControlProxyHttpClient is a httpclient sending request
    to the control proxy local port. It's used in SyncRPCClient
    for forwarding GatewayRequests from the cloud, and gets a GatewayResponse.
    """

    async def send(self, gateway_request):
        """
        Forwards the given request to the service provided
        in :authority and awaits a response

        Args:
            gateway_request: gateway_request: A GatewayRequest that is
        defined in the sync_rpc_service.proto. It has fields gwId, authority,
        path, headers, and payload.

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

        try:
            await client.wait_functional()
            req_headers = self._get_req_headers(gateway_request.headers,
                                                gateway_request.path,
                                                gateway_request.authority)
            body = gateway_request.payload
            stream_id = await client.start_request(req_headers)
            gw_resp = await self._await_gateway_response(client, stream_id,
                                                         body)
            return gw_resp
        except Exception as e:  # pylint: disable=broad-except
            logging.error("[SyncRPC] Exception in proxy_client: %s", e)
            raise e
        finally:
            client.close_connection()

    @staticmethod
    async def _get_client(service):
        (ip, port) = ServiceRegistry.get_service_address(service)
        return await aioh2.open_connection(ip, port)

    async def _await_gateway_response(self, client, stream_id, body):
        await client.send_data(stream_id, body, end_stream=True)

        resp_headers = await client.recv_response(stream_id)
        status = self._get_resp_status(resp_headers)
        payload = await client.read_stream(stream_id, -1)
        trailers = await client.recv_trailers(stream_id)
        headers = self._get_resp_headers(resp_headers, trailers)
        return GatewayResponse(status=status, headers=headers,
                               payload=payload)

    @staticmethod
    def _get_req_headers(raw_req_headers, path, authority):
        headers = [(":method", "POST"),
                   (":scheme", "http"),
                   (":path", path),
                   (":authority", authority)]
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
