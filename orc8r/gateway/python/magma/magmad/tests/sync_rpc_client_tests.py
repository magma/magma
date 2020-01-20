"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=protected-access

import asyncio
import unittest.mock

from orc8r.protos.sync_rpc_service_pb2 import GatewayRequest, GatewayResponse, \
    SyncRPCResponse, SyncRPCRequest

from magma.common.service_registry import ServiceRegistry
from magma.magmad.sync_rpc_client import SyncRPCClient


class SyncRPCClientTests(unittest.TestCase):
    """
    Tests for the SyncRPCClient
    """
    def setUp(self):
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        self._loop = loop
        self._sync_rpc_client = SyncRPCClient(loop=loop, response_timeout=3)
        self._sync_rpc_client._conn_closed_table = {
            12345: False
        }
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        ServiceRegistry._PROXY_CONFIG = {'local_port': 2345,
                                         'cloud_address': 'test',
                                         'proxy_cloud_connections': True}
        self._req_body = GatewayRequest(gwId="test id", authority='mobility',
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
        self._expected_resp = GatewayResponse(status="400",
                                              headers={"test_key": "test_val"},
                                              payload=b'\x00'
                                                      b'\x00\x00\x00\n\n\x08')
        self._expected_err_msg = "test error"

    def test_forward_request_conn_closed(self):
        self._sync_rpc_client.forward_request(
            SyncRPCRequest(reqId=12345, connClosed=True))
        self.assertEqual(self._sync_rpc_client._conn_closed_table[12345], True)

    def test_send_sync_rpc_response(self):
        expected = SyncRPCResponse(reqId=123, respBody=self._expected_resp)
        self._sync_rpc_client._response_queue.put(expected)
        res = self._sync_rpc_client.send_sync_rpc_response()
        actual = next(res)
        self.assertEqual(expected, actual)
        expected = SyncRPCResponse(heartBeat=True)
        actual = next(res)
        self.assertEqual(expected, actual)

    def test_retry_connect_sleep(self):
        self._sync_rpc_client._current_delay = 0
        for i in range(5):
            self._sync_rpc_client._retry_connect_sleep()
            if i == 4:
                self.assertEqual(self._sync_rpc_client.RETRY_MAX_DELAY_SECS,
                                 self._sync_rpc_client._current_delay)
            else:
                self.assertEqual(2 ** i, self._sync_rpc_client._current_delay)


if __name__ == "__main__":
    unittest.main()
