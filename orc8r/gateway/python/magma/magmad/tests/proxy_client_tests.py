"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=protected-access,unused-argument

import asyncio
import unittest.mock

from orc8r.protos.sync_rpc_service_pb2 import GatewayRequest

from magma.common.service_registry import ServiceRegistry
from magma.magmad.proxy_client import ControlProxyHttpClient


class MockClient(object):
    def __init__(self, payload, headers, trailers, expected_req):
        self._expected_payload = payload
        self._expected_headers = headers
        self._expected_trailers = trailers
        self._expected_req = expected_req

    async def start_request(self, headers):
        return 3

    async def end_stream(self, stream_id):
        return

    async def wait_functional(self):
        return

    async def send_data(self, stream_id, body, end_stream=False):
        return

    async def recv_response(self, stream_id):
        return self._expected_headers

    async def read_stream(self, stream_id, num_bytes):
        return self._expected_payload

    async def recv_trailers(self, stream_id):
        return self._expected_trailers

    def close_connection(self):
        return


class ProxyClientTests(unittest.TestCase):
    """
    Tests for the ProxyClient.
    """

    def setUp(self):
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        self._loop = loop
        self._proxy_client = ControlProxyHttpClient()
        ServiceRegistry._REGISTRY = {"services": {"mobilityd":
                                                  {"ip_address": "0.0.0.0",
                                                   "port": 3456}}
                                     }
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        self._req_body = GatewayRequest(gwId="test id", authority='mobilityd',
                                        path='/magma.MobilityService'
                                             '/ListAddedIPv4Blocks',
                                        headers={'te': 'trailers',
                                                 'content-type':
                                                     'application/grpc',
                                                 'user-agent':
                                                     'grpc-python/1.4.0',
                                                 'grpc-accept-encoding':
                                                     'identity'},
                                        payload=bytes.fromhex('0000000000'))

    @unittest.mock.patch('aioh2.open_connection')
    def test_http_client(self, mock_conn):
        req_body = GatewayRequest(gwId="test id", authority='mobilityd',
                                  path='/magma.MobilityService'
                                       '/ListAddedIPv4Blocks',
                                  headers={'te': 'trailers',
                                           'content-type': 'application/grpc',
                                           'user-agent': 'grpc-python/1.4.0',
                                           'grpc-accept-encoding': 'identity'},
                                  payload=bytes.fromhex('0000000000'))
        expected_payload = \
            b'\x00\x00\x00\x00\n\n\x08\x12\x04\xc0\xa8\x80\x00\x18\x18'
        expected_header = [(':status', '200'),
                           ('content-type', 'application/grpc')]
        expected_trailers = [('grpc-status', '0'), ('grpc-message', '')]

        mock_conn.side_effect = asyncio.coroutine(
            unittest.mock.MagicMock(
                return_value=MockClient(expected_payload, expected_header,
                                        expected_trailers, req_body)))
        future = asyncio.ensure_future(
            self._proxy_client.send(self._req_body))

        try:
            self._loop.run_until_complete(future)
            res = future.result()
            self.assertEqual(res.status,
                             '200')
            self.assertEqual(res.payload, expected_payload)
        except Exception as e:  # pylint: disable=broad-except
            self.fail(e)
        self._loop.close()


if __name__ == "__main__":
    unittest.main()
