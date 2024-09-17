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
# pylint: disable=protected-access,unused-argument

import asyncio
import queue
import unittest.mock

from magma.common.service_registry import ServiceRegistry
from magma.magmad.proxy_client import ControlProxyHttpClient
from orc8r.protos.sync_rpc_service_pb2 import GatewayRequest


class MockUnaryClient(object):
    def __init__(self, payload, headers, trailers, expected_req):
        self._expected_payload = payload
        self._expected_headers = headers
        self._expected_trailers = trailers
        self._expected_req = expected_req
        self._event_handlers = {}
        self._num_calls_read_stream = 0

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

    async def read_stream(self, stream_id, num_bytes=None):
        self._num_calls_read_stream += 1
        if self._num_calls_read_stream > 1:
            return b''
        return self._expected_payload

    async def recv_trailers(self, stream_id):
        return self._expected_trailers

    def close_connection(self):
        return


class MockStreamingClient(object):
    def __init__(self, payload, headers, trailers, expected_req):
        self._expected_payload = payload
        self._expected_headers = headers
        self._expected_trailers = trailers
        self._expected_req = expected_req
        self._event_handlers = {}
        self._num_calls_read_stream = 0

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

    async def read_stream(self, stream_id, num_bytes=None):
        self._num_calls_read_stream += 1
        if self._num_calls_read_stream > 2:
            return b''
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
        ServiceRegistry._REGISTRY = {
            "services": {
                "mobilityd":
                {
                    "ip_address": "0.0.0.0",
                    "port": 3456,
                },
            },
        }
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        self._req_body = GatewayRequest(
            gwId="test id", authority='mobilityd',
            path='/magma.MobilityService'
                 '/ListAddedIPv4Blocks',
            headers={
                'te': 'trailers',
                'content-type':
                    'application/grpc',
                'user-agent':
                    'grpc-python/1.4.0',
                'grpc-accept-encoding':
                    'identity',
            },
            payload=bytes.fromhex('0000000000'),
        )

    @unittest.mock.patch('aioh2.open_connection')
    def test_http_client_unary(self, mock_conn):
        req_body = GatewayRequest(
            gwId="test id", authority='mobilityd',
            path='/magma.MobilityService'
                 '/ListAddedIPv4Blocks',
            headers={
                'te': 'trailers',
                'content-type': 'application/grpc',
                'user-agent': 'grpc-python/1.4.0',
                'grpc-accept-encoding': 'identity',
            },
            payload=bytes.fromhex('0000000000'),
        )
        expected_payload = \
            b'\x00\x00\x00\x00\n\n\x08\x12\x04\xc0\xa8\x80\x00\x18\x18'
        expected_header = [
            (':status', '200'),
            ('content-type', 'application/grpc'),
        ]
        expected_trailers = [('grpc-status', '0'), ('grpc-message', '')]

        mock_conn.side_effect = asyncio.coroutine(
            unittest.mock.MagicMock(
                return_value=MockUnaryClient(
                    expected_payload, expected_header,
                    expected_trailers, req_body,
                ),
            ),
        )

        request_queue = queue.Queue()
        conn_closed_table = {
            1234: False,
        }

        future = asyncio.ensure_future(
            self._proxy_client.send(
                self._req_body,
                1234,
                request_queue,
                conn_closed_table,
            ),
        )

        self._loop.run_until_complete(future)

        self.assertEqual(request_queue.qsize(), 1)
        res = request_queue.get()
        self.assertEqual(res.reqId, 1234)
        self.assertEqual(res.heartBeat, False)
        self.assertEqual(res.respBody.status, '200')
        self.assertEqual(res.respBody.payload, expected_payload)
        self.assertEqual(res.respBody.headers['grpc-status'], '0')
        self._loop.close()

    @unittest.mock.patch('aioh2.open_connection')
    def test_http_client_stream(self, mock_conn):
        req_body = GatewayRequest(
            gwId="test id", authority='mobilityd',
            path='/magma.MobilityService'
                 '/ListAddedIPv4Blocks',
            headers={
                'te': 'trailers',
                'content-type': 'application/grpc',
                'user-agent': 'grpc-python/1.4.0',
                'grpc-accept-encoding': 'identity',
            },
            payload=bytes.fromhex('0000000000'),
        )
        expected_payload = \
            b'\x00\x00\x00\x00\n\n\x08\x12\x04\xc0\xa8\x80\x00\x18\x18'
        expected_header = [
            (':status', '200'),
            ('content-type', 'application/grpc'),
        ]
        expected_trailers = [('grpc-status', '0'), ('grpc-message', '')]

        mock_conn.side_effect = asyncio.coroutine(
            unittest.mock.MagicMock(
                return_value=MockStreamingClient(
                    expected_payload,
                    expected_header,
                    expected_trailers,
                    req_body,
                ),
            ),
        )

        request_queue = queue.Queue()
        conn_closed_table = {
            1234: False,
        }

        future = asyncio.ensure_future(
            self._proxy_client.send(
                self._req_body,
                1234,
                request_queue,
                conn_closed_table,
            ),
        )

        self._loop.run_until_complete(future)

        self.assertEqual(request_queue.qsize(), 2)
        res_1 = request_queue.get(timeout=0)
        self.assertEqual(res_1.reqId, 1234)
        self.assertEqual(res_1.heartBeat, False)
        self.assertEqual(res_1.respBody.status, '200')
        self.assertEqual(res_1.respBody.payload, expected_payload)
        self.assertTrue('grpc-status' not in res_1.respBody.headers)
        res_2 = request_queue.get()
        self.assertEqual(res_2.reqId, 1234)
        self.assertEqual(res_2.heartBeat, False)
        self.assertEqual(res_2.respBody.status, '200')
        self.assertEqual(res_2.respBody.payload, expected_payload)
        self.assertEqual(res_2.respBody.headers['grpc-status'], '0')
        self._loop.close()


if __name__ == "__main__":
    unittest.main()
